package keeper_test

import (
	"testing"

	"blogchain/x/blog/keeper"
	"blogchain/x/blog/types"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/stretchr/testify/require"
)

func setupKeeper(t *testing.T) (keeper.Keeper, sdk.Context) {
	f := initFixture(t)
	ctx := sdk.UnwrapSDKContext(f.ctx)
	return f.keeper, ctx
}

func TestPostOperations(t *testing.T) {
	k, ctx := setupKeeper(t)

	t.Run("AppendPost", func(t *testing.T) {
		post := types.Post{
			Creator:   "cosmos1abcdef",
			Title:     "Test Post",
			Body:      "This is a test post body",
			Tags:      []string{"test", "cosmos"},
			CreatedAt: 1234567890,
			Likes:     0,
		}

		// Test creating first post
		id, err := k.AppendPost(ctx, post)
		require.NoError(t, err)
		require.Equal(t, uint64(0), id) // Sequence starts at 0

		// Test creating second post
		post2 := post
		post2.Title = "Second Post"
		id2, err := k.AppendPost(ctx, post2)
		require.NoError(t, err)
		require.Equal(t, uint64(1), id2)
	})

	t.Run("GetPost", func(t *testing.T) {
		// Create a post first
		post := types.Post{
			Creator:   "cosmos1xyz",
			Title:     "Get Test",
			Body:      "Body for get test",
			Tags:      []string{"get"},
			CreatedAt: 1234567890,
			Likes:     0,
		}
		id, err := k.AppendPost(ctx, post)
		require.NoError(t, err)

		// Test getting existing post
		retrievedPost, err := k.GetPost(ctx, id)
		require.NoError(t, err)
		require.Equal(t, id, retrievedPost.Id)
		require.Equal(t, post.Title, retrievedPost.Title)
		require.Equal(t, post.Creator, retrievedPost.Creator)

		// Test getting non-existent post
		_, err = k.GetPost(ctx, 999999)
		require.Error(t, err)
		require.Contains(t, err.Error(), "post not found")
	})

	t.Run("SetPost", func(t *testing.T) {
		// Create a post first
		post := types.Post{
			Creator:   "cosmos1update",
			Title:     "Original Title",
			Body:      "Original body",
			Tags:      []string{"original"},
			CreatedAt: 1234567890,
			Likes:     0,
		}
		id, err := k.AppendPost(ctx, post)
		require.NoError(t, err)

		// Retrieve and update the post
		updatedPost, err := k.GetPost(ctx, id)
		require.NoError(t, err)
		updatedPost.Title = "Updated Title"
		updatedPost.Body = "Updated body"
		updatedPost.Likes = 10

		err = k.SetPost(ctx, updatedPost)
		require.NoError(t, err)

		// Verify the update
		retrievedPost, err := k.GetPost(ctx, id)
		require.NoError(t, err)
		require.Equal(t, "Updated Title", retrievedPost.Title)
		require.Equal(t, "Updated body", retrievedPost.Body)
		require.Equal(t, uint64(10), retrievedPost.Likes)

		// Test updating non-existent post
		nonExistentPost := types.Post{
			Id:    999999,
			Title: "Non-existent",
		}
		err = k.SetPost(ctx, nonExistentPost)
		require.Error(t, err)
		require.Contains(t, err.Error(), "post not found")
	})

	t.Run("GetPostCount", func(t *testing.T) {
		// Get initial count
		count, err := k.GetPostCount(ctx)
		require.NoError(t, err)
		initialCount := count

		// Add posts and verify count increases
		post := types.Post{
			Creator: "cosmos1count",
			Title:   "Count Test",
		}
		
		_, err = k.AppendPost(ctx, post)
		require.NoError(t, err)
		
		count, err = k.GetPostCount(ctx)
		require.NoError(t, err)
		require.Equal(t, initialCount+1, count)
		
		_, err = k.AppendPost(ctx, post)
		require.NoError(t, err)
		
		count, err = k.GetPostCount(ctx)
		require.NoError(t, err)
		require.Equal(t, initialCount+2, count)
	})

	t.Run("GetAllPosts", func(t *testing.T) {
		// Clear and add specific posts for this test
		k, ctx := setupKeeper(t) // Fresh keeper
		
		// Add multiple posts
		posts := []types.Post{
			{Creator: "user1", Title: "Post 1", Body: "Body 1"},
			{Creator: "user2", Title: "Post 2", Body: "Body 2"},
			{Creator: "user3", Title: "Post 3", Body: "Body 3"},
		}
		
		for _, post := range posts {
			_, err := k.AppendPost(ctx, post)
			require.NoError(t, err)
		}
		
		// Get all posts
		allPosts, err := k.GetAllPosts(ctx)
		require.NoError(t, err)
		require.Len(t, allPosts, 3)
		
		// Verify posts are correct
		for i, post := range allPosts {
			require.Equal(t, posts[i].Title, post.Title)
			require.Equal(t, posts[i].Creator, post.Creator)
			require.Equal(t, posts[i].Body, post.Body)
		}
	})

	t.Run("GetAllPostsPaginated", func(t *testing.T) {
		k, ctx := setupKeeper(t) // Fresh keeper
		
		// Add 10 posts for pagination testing
		for i := 0; i < 10; i++ {
			post := types.Post{
				Creator: "cosmos1page",
				Title:   "Post " + string(rune('A'+i)),
				Body:    "Body for pagination test",
			}
			_, err := k.AppendPost(ctx, post)
			require.NoError(t, err)
		}
		
		// Test first page
		pageReq := &types.PageRequest{
			Limit:  3,
			Offset: 0,
		}
		posts, pageRes, err := k.GetAllPostsPaginated(ctx, pageReq)
		require.NoError(t, err)
		require.Len(t, posts, 3)
		require.Equal(t, uint64(10), pageRes.Total)
		require.Equal(t, uint64(3), pageRes.Limit)
		require.Equal(t, uint64(0), pageRes.Offset)
		
		// Test second page
		pageReq = &types.PageRequest{
			Limit:  3,
			Offset: 3,
		}
		posts, pageRes, err = k.GetAllPostsPaginated(ctx, pageReq)
		require.NoError(t, err)
		require.Len(t, posts, 3)
		require.Equal(t, uint64(10), pageRes.Total)
		
		// Test last page (partial)
		pageReq = &types.PageRequest{
			Limit:  3,
			Offset: 9,
		}
		posts, pageRes, err = k.GetAllPostsPaginated(ctx, pageReq)
		require.NoError(t, err)
		require.Len(t, posts, 1) // Only 1 post remaining
		
		// Test nil pagination (should use defaults)
		posts, pageRes, err = k.GetAllPostsPaginated(ctx, nil)
		require.NoError(t, err)
		require.Len(t, posts, 10) // Default limit is 50, so all 10 posts returned
		
		// Test limit exceeding 100 (should be capped)
		pageReq = &types.PageRequest{
			Limit:  200,
			Offset: 0,
		}
		posts, pageRes, err = k.GetAllPostsPaginated(ctx, pageReq)
		require.NoError(t, err)
		require.Equal(t, uint64(50), pageRes.Limit) // Capped at 50
	})

	t.Run("UserLikedPost", func(t *testing.T) {
		k, ctx := setupKeeper(t)
		
		// Create a post
		post := types.Post{
			Creator: "cosmos1like",
			Title:   "Like Test",
		}
		postID, err := k.AppendPost(ctx, post)
		require.NoError(t, err)
		
		userAddr := "cosmos1user"
		
		// Check initial state - user hasn't liked
		liked, err := k.HasUserLikedPost(ctx, postID, userAddr)
		require.NoError(t, err)
		require.False(t, liked)
		
		// Set user liked post
		err = k.SetUserLikedPost(ctx, postID, userAddr)
		require.NoError(t, err)
		
		// Check user has now liked
		liked, err = k.HasUserLikedPost(ctx, postID, userAddr)
		require.NoError(t, err)
		require.True(t, liked)
		
		// Different user hasn't liked
		liked, err = k.HasUserLikedPost(ctx, postID, "cosmos1other")
		require.NoError(t, err)
		require.False(t, liked)
		
		// Non-existent post returns false (not error)
		liked, err = k.HasUserLikedPost(ctx, 999999, userAddr)
		require.NoError(t, err)
		require.False(t, liked)
	})
}

func TestPostSequentialIDs(t *testing.T) {
	k, ctx := setupKeeper(t)
	
	// Create multiple posts and verify IDs are sequential
	ids := make([]uint64, 5)
	for i := range ids {
		post := types.Post{
			Creator: "cosmos1seq",
			Title:   "Sequential Test",
		}
		id, err := k.AppendPost(ctx, post)
		require.NoError(t, err)
		ids[i] = id
	}
	
	// Verify IDs are sequential starting from 0
	for i, id := range ids {
		require.Equal(t, uint64(i), id)
	}
}