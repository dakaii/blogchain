package keeper_test

import (
	"testing"

	"blogchain/x/blog/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

func TestMsgServerDeletePost(t *testing.T) {
	f := initFixture(t)
	msgServer := NewMsgServerImpl(f.keeper)
	ctx := sdk.UnwrapSDKContext(f.ctx)
	
	// Create posts for testing
	creator := GetTestAddress(t)
	otherUser := GetTestAddress(t)
	
	// Create first post
	createMsg1 := &types.MsgCreatePost{
		Creator: creator,
		Title:   "Post to Delete",
		Body:    "This post will be deleted",
		Tags:    []string{"delete", "test"},
	}
	createResp1, err := msgServer.CreatePost(f.ctx, createMsg1)
	require.NoError(t, err)
	postId1 := createResp1.Id
	
	// Create second post
	createMsg2 := &types.MsgCreatePost{
		Creator: creator,
		Title:   "Post to Keep",
		Body:    "This post will remain",
		Tags:    []string{"keep", "test"},
	}
	createResp2, err := msgServer.CreatePost(f.ctx, createMsg2)
	require.NoError(t, err)
	postId2 := createResp2.Id
	
	t.Run("Valid delete by creator", func(t *testing.T) {
		deleteMsg := &types.MsgDeletePost{
			Creator: creator,
			Id:      postId1,
		}
		
		resp, err := msgServer.DeletePost(f.ctx, deleteMsg)
		require.NoError(t, err)
		require.NotNil(t, resp)
		
		// Verify post is soft deleted
		post, err := f.keeper.GetPost(ctx, postId1)
		require.NoError(t, err)
		require.True(t, post.Deleted)
		require.Greater(t, post.DeletedAt, int64(0))
		
		// Original data should still be intact
		require.Equal(t, "Post to Delete", post.Title)
		require.Equal(t, creator, post.Creator)
	})
	
	t.Run("Delete by non-creator fails", func(t *testing.T) {
		deleteMsg := &types.MsgDeletePost{
			Creator: otherUser,
			Id:      postId2,
		}
		
		resp, err := msgServer.DeletePost(f.ctx, deleteMsg)
		require.Error(t, err)
		require.Nil(t, resp)
		require.Contains(t, err.Error(), "only the post creator can delete")
		
		// Verify post was not deleted
		post, err := f.keeper.GetPost(ctx, postId2)
		require.NoError(t, err)
		require.False(t, post.Deleted)
	})
	
	t.Run("Delete non-existent post fails", func(t *testing.T) {
		deleteMsg := &types.MsgDeletePost{
			Creator: creator,
			Id:      999999,
		}
		
		resp, err := msgServer.DeletePost(f.ctx, deleteMsg)
		require.Error(t, err)
		require.Nil(t, resp)
		require.Contains(t, err.Error(), "post not found")
	})
	
	t.Run("Double delete fails", func(t *testing.T) {
		// Try to delete already deleted post
		deleteMsg := &types.MsgDeletePost{
			Creator: creator,
			Id:      postId1, // Already deleted in first test
		}
		
		resp, err := msgServer.DeletePost(f.ctx, deleteMsg)
		require.Error(t, err)
		require.Nil(t, resp)
		require.Contains(t, err.Error(), "post is already deleted")
	})
	
	t.Run("Cannot update deleted post", func(t *testing.T) {
		// Try to update the deleted post
		updateMsg := &types.MsgUpdatePost{
			Creator: creator,
			Id:      postId1, // Deleted post
			Title:   "Trying to update deleted",
			Body:    "This should fail",
			Tags:    []string{"fail"},
		}
		
		resp, err := msgServer.UpdatePost(f.ctx, updateMsg)
		require.Error(t, err)
		require.Nil(t, resp)
		require.Contains(t, err.Error(), "cannot update a deleted post")
	})
	
	t.Run("Deleted posts excluded from queries", func(t *testing.T) {
		// Get all posts - deleted ones should be excluded
		posts, err := f.keeper.GetAllPosts(ctx)
		require.NoError(t, err)
		
		// Check that deleted post is not in the list
		for _, post := range posts {
			require.NotEqual(t, postId1, post.Id, "Deleted post should not appear in GetAllPosts")
			require.False(t, post.Deleted)
		}
		
		// Check that non-deleted post is still there
		found := false
		for _, post := range posts {
			if post.Id == postId2 {
				found = true
				break
			}
		}
		require.True(t, found, "Non-deleted post should appear in GetAllPosts")
	})
	
	t.Run("Active post count excludes deleted", func(t *testing.T) {
		// Create one more post
		createMsg := &types.MsgCreatePost{
			Creator: creator,
			Title:   "Another Post",
			Body:    "Another body",
			Tags:    []string{"another"},
		}
		_, err := msgServer.CreatePost(f.ctx, createMsg)
		require.NoError(t, err)
		
		// Get active count (should exclude deleted posts)
		activeCount, err := f.keeper.GetActivePostCount(ctx)
		require.NoError(t, err)
		
		// Get all posts to verify count
		posts, err := f.keeper.GetAllPosts(ctx)
		require.NoError(t, err)
		
		require.Equal(t, uint64(len(posts)), activeCount)
		
		// Total count should include deleted
		totalCount, err := f.keeper.GetPostCount(ctx)
		require.NoError(t, err)
		require.Greater(t, totalCount, activeCount)
	})
}