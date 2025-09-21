package keeper_test

import (
	"testing"
	
	"blogchain/x/blog/types"
	"github.com/stretchr/testify/require"
)

func TestMsgServerLikePost(t *testing.T) {
	f := initFixture(t)
	require := require.New(t)
	
	addresses := GenerateTestAddresses(t, 2)
	creator := addresses[0]
	liker := addresses[1]
	
	msgServer := NewMsgServerImpl(f.keeper)
	
	// First create a post
	createMsg := &types.MsgCreatePost{
		Creator: creator,
		Title:   "Post to Like",
		Body:    "This post will be liked",
		Tags:    []string{"likeable"},
	}
	
	createResp, err := msgServer.CreatePost(f.ctx, createMsg)
	require.NoError(err)
	postId := createResp.Id
	
	// Get initial state
	post, found := f.keeper.GetPost(f.ctx, postId)
	require.True(found)
	require.Equal(uint64(0), post.Likes)
	
	// Like the post
	likeMsg := &types.MsgLikePost{
		Liker:  liker,
		PostId: postId,
	}
	
	resp, err := msgServer.LikePost(f.ctx, likeMsg)
	require.NoError(err)
	require.NotNil(resp)
	
	// Verify like count increased
	post, found = f.keeper.GetPost(f.ctx, postId)
	require.True(found)
	require.Equal(uint64(1), post.Likes)
}

func TestMsgServerLikePostMultipleTimes(t *testing.T) {
	f := initFixture(t)
	require := require.New(t)
	
	creator := GetTestAddress(t)
	msgServer := NewMsgServerImpl(f.keeper)
	
	// Create a post
	createMsg := &types.MsgCreatePost{
		Creator: creator,
		Title:   "Popular Post",
		Body:    "This post will get many likes",
		Tags:    []string{"popular"},
	}
	
	createResp, err := msgServer.CreatePost(f.ctx, createMsg)
	require.NoError(err)
	postId := createResp.Id
	
	// Multiple users like the post
	likers := GenerateTestAddresses(t, 5)
	
	for _, liker := range likers {
		likeMsg := &types.MsgLikePost{
			Liker:  liker,
			PostId: postId,
		}
		
		resp, err := msgServer.LikePost(f.ctx, likeMsg)
		require.NoError(err)
		require.NotNil(resp)
	}
	
	// Verify final like count
	post, found := f.keeper.GetPost(f.ctx, postId)
	require.True(found)
	require.Equal(uint64(len(likers)), post.Likes)
}

func TestMsgServerLikePostNonExistent(t *testing.T) {
	f := initFixture(t)
	require := require.New(t)
	
	liker := GetTestAddress(t)
	msgServer := NewMsgServerImpl(f.keeper)
	
	// Try to like non-existent post
	likeMsg := &types.MsgLikePost{
		Liker:  liker,
		PostId: 999999,
	}
	
	_, err := msgServer.LikePost(f.ctx, likeMsg)
	require.Error(err)
	require.Contains(err.Error(), "post not found")
}

func TestMsgServerLikePostInvalidPostId(t *testing.T) {
	f := initFixture(t)
	require := require.New(t)
	
	liker := GetTestAddress(t)
	msgServer := NewMsgServerImpl(f.keeper)
	
	// Try to like with invalid post ID  
	// Note: PostId is uint64, so we can't test non-numeric strings
	// ID 0 is never valid since we start at 1
	likeMsg := &types.MsgLikePost{
		Liker:  liker,
		PostId: 0,
	}
	
	_, err := msgServer.LikePost(f.ctx, likeMsg)
	require.Error(err)
	require.Contains(err.Error(), "post not found")
}

func TestMsgServerLikePostInvalidLiker(t *testing.T) {
	f := initFixture(t)
	require := require.New(t)
	
	creator := GetTestAddress(t)
	msgServer := NewMsgServerImpl(f.keeper)
	
	// Create a post
	createMsg := &types.MsgCreatePost{
		Creator: creator,
		Title:   "Post",
		Body:    "Body",
		Tags:    []string{"test"},
	}
	
	createResp, err := msgServer.CreatePost(f.ctx, createMsg)
	require.NoError(err)
	postId := createResp.Id
	
	// Try to like with invalid liker address
	likeMsg := &types.MsgLikePost{
		Liker:  "invalid-address",
		PostId: postId,
	}
	
	_, err = msgServer.LikePost(f.ctx, likeMsg)
	require.Error(err)
	require.Contains(err.Error(), "invalid liker address")
}

func TestMsgServerLikeSamePostTwice(t *testing.T) {
	f := initFixture(t)
	require := require.New(t)
	
	addresses := GenerateTestAddresses(t, 2)
	creator := addresses[0]
	liker := addresses[1]
	msgServer := NewMsgServerImpl(f.keeper)
	
	// Create a post
	createMsg := &types.MsgCreatePost{
		Creator: creator,
		Title:   "Post",
		Body:    "Body",
		Tags:    []string{"test"},
	}
	
	createResp, err := msgServer.CreatePost(f.ctx, createMsg)
	require.NoError(err)
	postId := createResp.Id
	
	// Like the post
	likeMsg := &types.MsgLikePost{
		Liker:  liker,
		PostId: postId,
	}
	
	resp, err := msgServer.LikePost(f.ctx, likeMsg)
	require.NoError(err)
	require.NotNil(resp)
	
	// Try to like the same post again with same user
	// This should ideally fail, but depends on your business logic
	// For now, it seems to increment the count
	resp, err = msgServer.LikePost(f.ctx, likeMsg)
	require.NoError(err)
	require.NotNil(resp)
	
	// Verify like count
	post, found := f.keeper.GetPost(f.ctx, postId)
	require.True(found)
	require.Equal(uint64(2), post.Likes) // Currently allows duplicate likes
}