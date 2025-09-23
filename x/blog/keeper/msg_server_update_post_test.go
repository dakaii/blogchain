package keeper_test

import (
	"testing"

	"blogchain/x/blog/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

func TestMsgServerUpdatePost(t *testing.T) {
	f := initFixture(t)
	msgServer := NewMsgServerImpl(f.keeper)
	ctx := sdk.UnwrapSDKContext(f.ctx)
	
	// Create a post first
	creator := GetTestAddress(t)
	createMsg := &types.MsgCreatePost{
		Creator: creator,
		Title:   "Original Title",
		Body:    "Original Body",
		Tags:    []string{"original", "test"},
	}
	
	createResp, err := msgServer.CreatePost(f.ctx, createMsg)
	require.NoError(t, err)
	postId := createResp.Id
	
	t.Run("Valid update by creator", func(t *testing.T) {
		updateMsg := &types.MsgUpdatePost{
			Creator:       creator,
			Id:            postId,
			Title:         "Updated Title",
			Body:          "Updated Body with more content",
			Tags:          []string{"updated", "modified", "test"},
			MediaBlobIds:  []string{"blob123", "blob456"},
			ContentBlobId: "content789",
		}
		
		resp, err := msgServer.UpdatePost(f.ctx, updateMsg)
		require.NoError(t, err)
		require.NotNil(t, resp)
		
		// Verify post was updated
		post, err := f.keeper.GetPost(ctx, postId)
		require.NoError(t, err)
		require.Equal(t, "Updated Title", post.Title)
		require.Equal(t, "Updated Body with more content", post.Body)
		require.Equal(t, []string{"updated", "modified", "test"}, post.Tags)
		require.Equal(t, []string{"blob123", "blob456"}, post.MediaBlobIds)
		require.Equal(t, "content789", post.ContentBlobId)
		require.Greater(t, post.UpdatedAt, int64(0))
		// In tests, block time may be the same, so UpdatedAt could equal CreatedAt
		require.GreaterOrEqual(t, post.UpdatedAt, post.CreatedAt)
	})
	
	t.Run("Update by non-creator fails", func(t *testing.T) {
		otherUser := GetTestAddress(t)
		updateMsg := &types.MsgUpdatePost{
			Creator: otherUser,
			Id:      postId,
			Title:   "Hacker Title",
			Body:    "Trying to hack",
			Tags:    []string{"hack"},
		}
		
		resp, err := msgServer.UpdatePost(f.ctx, updateMsg)
		require.Error(t, err)
		require.Nil(t, resp)
		require.Contains(t, err.Error(), "only the post creator can update")
		
		// Verify post was not changed
		post, err := f.keeper.GetPost(ctx, postId)
		require.NoError(t, err)
		require.NotEqual(t, "Hacker Title", post.Title)
	})
	
	t.Run("Update non-existent post fails", func(t *testing.T) {
		updateMsg := &types.MsgUpdatePost{
			Creator: creator,
			Id:      999999,
			Title:   "Ghost Post",
			Body:    "This should fail",
			Tags:    []string{"ghost"},
		}
		
		resp, err := msgServer.UpdatePost(f.ctx, updateMsg)
		require.Error(t, err)
		require.Nil(t, resp)
		require.Contains(t, err.Error(), "post not found")
	})
	
	t.Run("Multiple updates allowed", func(t *testing.T) {
		// First update
		updateMsg1 := &types.MsgUpdatePost{
			Creator: creator,
			Id:      postId,
			Title:   "First Update",
			Body:    "First update body",
			Tags:    []string{"first"},
		}
		
		resp, err := msgServer.UpdatePost(f.ctx, updateMsg1)
		require.NoError(t, err)
		require.NotNil(t, resp)
		
		post, err := f.keeper.GetPost(ctx, postId)
		require.NoError(t, err)
		firstUpdateTime := post.UpdatedAt
		
		// Second update
		updateMsg2 := &types.MsgUpdatePost{
			Creator: creator,
			Id:      postId,
			Title:   "Second Update",
			Body:    "Second update body",
			Tags:    []string{"second"},
		}
		
		resp, err = msgServer.UpdatePost(f.ctx, updateMsg2)
		require.NoError(t, err)
		require.NotNil(t, resp)
		
		post, err = f.keeper.GetPost(ctx, postId)
		require.NoError(t, err)
		require.Equal(t, "Second Update", post.Title)
		require.GreaterOrEqual(t, post.UpdatedAt, firstUpdateTime)
	})
}