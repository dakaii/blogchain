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
	if _, err := k.Posts.Get(ctx, post.Id); err != nil {
		if errors.Is(err, collections.ErrNotFound) {
			return fmt.Errorf("post not found: %d", post.Id)
		}
		return fmt.Errorf("failed to verify post exists: %w", err)
	}
	
	if err := k.Posts.Set(ctx, post.Id, post); err != nil {
		return fmt.Errorf("failed to update post: %w", err)
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

// GetAllPostsPaginated retrieves posts with pagination support
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
	
	err := k.Posts.Walk(ctx, nil, func(key uint64, post types.Post) (stop bool, err error) {
		// Skip entries based on offset
		if skipped < pageReq.Offset {
			skipped++
			return false, nil
		}
		
		// Stop if we've collected enough entries
		if uint64(len(posts)) >= pageReq.Limit {
			return true, nil
		}
		
		posts = append(posts, post)
		return false, nil
	})
	
	if err != nil {
		return nil, nil, fmt.Errorf("failed to paginate posts: %w", err)
	}
	
	// Get total count for pagination response
	totalCount, err := k.GetPostCount(ctx)
	if err != nil {
		return nil, nil, err
	}
	
	pageRes := &types.PageResponse{
		Total:  totalCount,
		Limit:  pageReq.Limit,
		Offset: pageReq.Offset,
	}
	
	return posts, pageRes, nil
}

// GetAllPosts retrieves all posts (backward compatibility)
func (k Keeper) GetAllPosts(ctx context.Context) ([]types.Post, error) {
	var posts []types.Post
	
	err := k.Posts.Walk(ctx, nil, func(key uint64, post types.Post) (stop bool, err error) {
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