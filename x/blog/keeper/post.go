package keeper

import (
	"context"
	"errors"
	"fmt"

	"blogchain/x/blog/types"

	"cosmossdk.io/collections"
)

// AppendPost creates a new post with sequential ID
func (k Keeper) AppendPost(ctx context.Context, post types.Post) (uint64, error) {
	// Get next ID from sequence
	postID, err := k.PostCount.Next(ctx)
	if err != nil {
		return 0, fmt.Errorf("failed to get next post ID: %w", err)
	}
	
	post.Id = postID
	
	// Save the post using collections API
	if err := k.Posts.Set(ctx, postID, post); err != nil {
		return 0, fmt.Errorf("failed to save post: %w", err)
	}
	
	// Add to active posts index
	if err := k.ActivePosts.Set(ctx, postID, true); err != nil {
		return 0, fmt.Errorf("failed to add post to active index: %w", err)
	}
	
	return postID, nil
}

// GetPost retrieves a post by ID
func (k Keeper) GetPost(ctx context.Context, id uint64) (types.Post, error) {
	post, err := k.Posts.Get(ctx, id)
	if err != nil {
		if errors.Is(err, collections.ErrNotFound) {
			return types.Post{}, fmt.Errorf("post not found: %d", id)
		}
		return types.Post{}, fmt.Errorf("failed to get post: %w", err)
	}
	return post, nil
}

// SetPost updates an existing post
func (k Keeper) SetPost(ctx context.Context, post types.Post) error {
	// Verify post exists
	oldPost, err := k.Posts.Get(ctx, post.Id)
	if err != nil {
		if errors.Is(err, collections.ErrNotFound) {
			return fmt.Errorf("post not found: %d", post.Id)
		}
		return fmt.Errorf("failed to verify post exists: %w", err)
	}
	
	if err := k.Posts.Set(ctx, post.Id, post); err != nil {
		return fmt.Errorf("failed to update post: %w", err)
	}
	
	// Update indexes if deletion status changed
	if oldPost.Deleted != post.Deleted {
		if post.Deleted {
			// Move from active to deleted
			if err := k.ActivePosts.Remove(ctx, post.Id); err != nil && !errors.Is(err, collections.ErrNotFound) {
				return fmt.Errorf("failed to remove from active index: %w", err)
			}
			if err := k.DeletedPosts.Set(ctx, post.Id, true); err != nil {
				return fmt.Errorf("failed to add to deleted index: %w", err)
			}
		} else {
			// Move from deleted to active (restore)
			if err := k.DeletedPosts.Remove(ctx, post.Id); err != nil && !errors.Is(err, collections.ErrNotFound) {
				return fmt.Errorf("failed to remove from deleted index: %w", err)
			}
			if err := k.ActivePosts.Set(ctx, post.Id, true); err != nil {
				return fmt.Errorf("failed to add to active index: %w", err)
			}
		}
	}
	
	return nil
}

// GetPostCount returns the total number of posts
func (k Keeper) GetPostCount(ctx context.Context) (uint64, error) {
	seq, err := k.PostCount.Peek(ctx)
	if err != nil {
		if errors.Is(err, collections.ErrNotFound) {
			return 0, nil
		}
		return 0, fmt.Errorf("failed to get post count: %w", err)
	}
	return seq, nil
}

// GetActivePostCount returns the number of non-deleted posts efficiently
func (k Keeper) GetActivePostCount(ctx context.Context) (uint64, error) {
	count := uint64(0)
	err := k.ActivePosts.Walk(ctx, nil, func(_ uint64, _ bool) (stop bool, err error) {
		count++
		return false, nil
	})
	
	if err != nil {
		return 0, fmt.Errorf("failed to count active posts: %w", err)
	}
	
	return count, nil
}

// GetAllPostsPaginated retrieves non-deleted posts with pagination support
func (k Keeper) GetAllPostsPaginated(ctx context.Context, pageReq *types.PageRequest) ([]types.Post, *types.PageResponse, error) {
	var posts []types.Post
	
	// Default pagination if not provided
	if pageReq == nil {
		pageReq = &types.PageRequest{
			Limit:  50,
			Offset: 0,
		}
	}
	
	// Ensure reasonable limits
	if pageReq.Limit == 0 || pageReq.Limit > 100 {
		pageReq.Limit = 50
	}
	
	skipped := uint64(0)
	activePostCount := uint64(0)
	
	// Walk through active posts index only
	err := k.ActivePosts.Walk(ctx, nil, func(postID uint64, _ bool) (stop bool, err error) {
		activePostCount++
		
		// Skip entries based on offset
		if skipped < pageReq.Offset {
			skipped++
			return false, nil
		}
		
		// Stop if we've collected enough entries
		if uint64(len(posts)) >= pageReq.Limit {
			return true, nil
		}
		
		// Fetch the actual post
		post, err := k.Posts.Get(ctx, postID)
		if err != nil {
			if errors.Is(err, collections.ErrNotFound) {
				return false, nil // Skip if not found
			}
			return false, err
		}
		
		posts = append(posts, post)
		return false, nil
	})
	
	if err != nil {
		return nil, nil, fmt.Errorf("failed to paginate posts: %w", err)
	}
	
	pageRes := &types.PageResponse{
		Total:  activePostCount, // Only count non-deleted posts
		Limit:  pageReq.Limit,
		Offset: pageReq.Offset,
	}
	
	return posts, pageRes, nil
}

// GetAllPosts retrieves all non-deleted posts efficiently using index
func (k Keeper) GetAllPosts(ctx context.Context) ([]types.Post, error) {
	var posts []types.Post
	
	// Walk through active posts index only
	err := k.ActivePosts.Walk(ctx, nil, func(postID uint64, _ bool) (stop bool, err error) {
		post, err := k.Posts.Get(ctx, postID)
		if err != nil {
			// Handle case where index is out of sync
			if errors.Is(err, collections.ErrNotFound) {
				return false, nil // Skip this entry
			}
			return false, err
		}
		posts = append(posts, post)
		return false, nil
	})
	
	if err != nil {
		return nil, fmt.Errorf("failed to get all posts: %w", err)
	}
	
	return posts, nil
}

// HasUserLikedPost checks if a user has already liked a post
func (k Keeper) HasUserLikedPost(ctx context.Context, postID uint64, userAddr string) (bool, error) {
	key := collections.Join(postID, userAddr)
	liked, err := k.LikedBy.Get(ctx, key)
	if err != nil {
		if errors.Is(err, collections.ErrNotFound) {
			return false, nil
		}
		return false, fmt.Errorf("failed to check if user liked post: %w", err)
	}
	return liked, nil
}

// SetUserLikedPost marks that a user has liked a post
func (k Keeper) SetUserLikedPost(ctx context.Context, postID uint64, userAddr string) error {
	key := collections.Join(postID, userAddr)
	if err := k.LikedBy.Set(ctx, key, true); err != nil {
		return fmt.Errorf("failed to set user liked post: %w", err)
	}
	return nil
}