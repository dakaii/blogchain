package keeper

import (
	"context"
	"errors"
	"fmt"

	"blogchain/x/blog/types"

	"cosmossdk.io/collections"
)

const MAX_COMMENT_DEPTH = 5 // Maximum nesting depth for comments

// CreateComment creates a new comment
func (k Keeper) CreateComment(ctx context.Context, comment types.Comment) (uint64, error) {
	// Verify post exists and is not deleted
	post, err := k.GetPost(ctx, comment.PostId)
	if err != nil {
		return 0, fmt.Errorf("post not found: %w", err)
	}
	if post.Deleted {
		return 0, fmt.Errorf("cannot comment on deleted post")
	}

	// If it's a reply, verify parent comment exists
	if comment.ParentId > 0 {
		parentComment, err := k.GetComment(ctx, comment.ParentId)
		if err != nil {
			return 0, fmt.Errorf("parent comment not found: %w", err)
		}
		if parentComment.Deleted {
			return 0, fmt.Errorf("cannot reply to deleted comment")
		}
		// Check nesting depth
		if parentComment.Depth >= MAX_COMMENT_DEPTH {
			return 0, fmt.Errorf("maximum comment depth (%d) reached", MAX_COMMENT_DEPTH)
		}
		comment.Depth = parentComment.Depth + 1
		// Ensure reply is on same post
		if parentComment.PostId != comment.PostId {
			return 0, fmt.Errorf("reply must be on the same post as parent")
		}
	} else {
		comment.Depth = 0 // Root comment
	}

	// Get next ID from sequence
	commentID, err := k.CommentCount.Next(ctx)
	if err != nil {
		return 0, fmt.Errorf("failed to get next comment ID: %w", err)
	}
	comment.Id = commentID

	// Save the comment
	if err := k.Comments.Set(ctx, commentID, comment); err != nil {
		return 0, fmt.Errorf("failed to save comment: %w", err)
	}

	// Add to active comments index
	if err := k.ActiveComments.Set(ctx, commentID, true); err != nil {
		return 0, fmt.Errorf("failed to add comment to active index: %w", err)
	}

	// Add to post comments index
	postCommentKey := collections.Join(comment.PostId, commentID)
	if err := k.PostComments.Set(ctx, postCommentKey, true); err != nil {
		return 0, fmt.Errorf("failed to add comment to post index: %w", err)
	}

	// If it's a reply, add to child comments index
	if comment.ParentId > 0 {
		childKey := collections.Join(comment.ParentId, commentID)
		if err := k.ChildComments.Set(ctx, childKey, true); err != nil {
			return 0, fmt.Errorf("failed to add comment to child index: %w", err)
		}
	}

	// Update post comment count
	post.CommentCount++
	if err := k.SetPost(ctx, post); err != nil {
		return 0, fmt.Errorf("failed to update post comment count: %w", err)
	}

	return commentID, nil
}

// GetComment retrieves a comment by ID
func (k Keeper) GetComment(ctx context.Context, id uint64) (types.Comment, error) {
	comment, err := k.Comments.Get(ctx, id)
	if err != nil {
		if errors.Is(err, collections.ErrNotFound) {
			return types.Comment{}, fmt.Errorf("comment not found: %d", id)
		}
		return types.Comment{}, fmt.Errorf("failed to get comment: %w", err)
	}
	return comment, nil
}

// UpdateComment updates an existing comment
func (k Keeper) UpdateComment(ctx context.Context, comment types.Comment) error {
	// Verify comment exists
	oldComment, err := k.Comments.Get(ctx, comment.Id)
	if err != nil {
		if errors.Is(err, collections.ErrNotFound) {
			return fmt.Errorf("comment not found: %d", comment.Id)
		}
		return fmt.Errorf("failed to verify comment exists: %w", err)
	}

	// Check if comment is deleted
	if oldComment.Deleted {
		return fmt.Errorf("cannot update deleted comment")
	}

	// Preserve immutable fields
	comment.PostId = oldComment.PostId
	comment.ParentId = oldComment.ParentId
	comment.CreatedAt = oldComment.CreatedAt
	comment.Depth = oldComment.Depth
	comment.Likes = oldComment.Likes

	// Save updated comment
	if err := k.Comments.Set(ctx, comment.Id, comment); err != nil {
		return fmt.Errorf("failed to update comment: %w", err)
	}

	return nil
}

// DeleteComment soft deletes a comment
func (k Keeper) DeleteComment(ctx context.Context, commentID uint64, deletedAt int64) error {
	comment, err := k.GetComment(ctx, commentID)
	if err != nil {
		return err
	}

	if comment.Deleted {
		return fmt.Errorf("comment already deleted")
	}

	// Soft delete the comment
	comment.Deleted = true
	comment.DeletedAt = deletedAt

	if err := k.Comments.Set(ctx, commentID, comment); err != nil {
		return fmt.Errorf("failed to delete comment: %w", err)
	}

	// Remove from active index, add to deleted
	if err := k.ActiveComments.Remove(ctx, commentID); err != nil && !errors.Is(err, collections.ErrNotFound) {
		return fmt.Errorf("failed to remove from active index: %w", err)
	}

	// Update post comment count
	post, err := k.GetPost(ctx, comment.PostId)
	if err == nil && !post.Deleted {
		if post.CommentCount > 0 {
			post.CommentCount--
			if err := k.SetPost(ctx, post); err != nil {
				return fmt.Errorf("failed to update post comment count: %w", err)
			}
		}
	}

	return nil
}

// GetPostComments retrieves comments for a post
func (k Keeper) GetPostComments(ctx context.Context, postID uint64, parentID uint64) ([]types.Comment, error) {
	var comments []types.Comment

	if parentID == 0 {
		// Get root comments for the post
		err := k.PostComments.Walk(ctx, collections.NewPrefixedPairRange[uint64, uint64](postID), 
			func(key collections.Pair[uint64, uint64], _ bool) (stop bool, err error) {
				commentID := key.K2()
				comment, err := k.GetComment(ctx, commentID)
				if err != nil {
					return false, nil // Skip if not found
				}
				// Only include root, non-deleted comments
				if !comment.Deleted && comment.ParentId == 0 {
					comments = append(comments, comment)
				}
				return false, nil
			})
		if err != nil {
			return nil, fmt.Errorf("failed to get post comments: %w", err)
		}
	} else {
		// Get replies to a specific comment
		err := k.ChildComments.Walk(ctx, collections.NewPrefixedPairRange[uint64, uint64](parentID),
			func(key collections.Pair[uint64, uint64], _ bool) (stop bool, err error) {
				commentID := key.K2()
				comment, err := k.GetComment(ctx, commentID)
				if err != nil {
					return false, nil // Skip if not found
				}
				if !comment.Deleted {
					comments = append(comments, comment)
				}
				return false, nil
			})
		if err != nil {
			return nil, fmt.Errorf("failed to get child comments: %w", err)
		}
	}

	return comments, nil
}

// GetCommentThread recursively builds a comment thread
func (k Keeper) GetCommentThread(ctx context.Context, commentID uint64, maxDepth uint32) (*types.CommentThread, error) {
	comment, err := k.GetComment(ctx, commentID)
	if err != nil {
		return nil, err
	}

	thread := &types.CommentThread{
		Comment: comment,
		Replies: make([]*types.CommentThread, 0),
	}

	// Don't fetch replies if we've reached max depth
	if maxDepth > 0 && comment.Depth >= maxDepth {
		return thread, nil
	}

	// Get child comments
	err = k.ChildComments.Walk(ctx, collections.NewPrefixedPairRange[uint64, uint64](commentID),
		func(key collections.Pair[uint64, uint64], _ bool) (stop bool, err error) {
			childID := key.K2()
			childThread, err := k.GetCommentThread(ctx, childID, maxDepth)
			if err != nil {
				return false, nil // Skip on error
			}
			if !childThread.Comment.Deleted {
				thread.Replies = append(thread.Replies, childThread)
			}
			return false, nil
		})

	if err != nil {
		return nil, fmt.Errorf("failed to get comment thread: %w", err)
	}

	return thread, nil
}

// LikeComment adds a like to a comment
func (k Keeper) LikeComment(ctx context.Context, commentID uint64, liker string) error {
	comment, err := k.GetComment(ctx, commentID)
	if err != nil {
		return err
	}

	if comment.Deleted {
		return fmt.Errorf("cannot like deleted comment")
	}

	// Check if already liked
	likeKey := collections.Join(commentID, liker)
	liked, _ := k.LikedComments.Get(ctx, likeKey)
	if liked {
		return fmt.Errorf("comment already liked by user")
	}

	// Add like
	comment.Likes++
	if err := k.Comments.Set(ctx, commentID, comment); err != nil {
		return fmt.Errorf("failed to update comment likes: %w", err)
	}

	// Mark as liked
	if err := k.LikedComments.Set(ctx, likeKey, true); err != nil {
		return fmt.Errorf("failed to record comment like: %w", err)
	}

	return nil
}