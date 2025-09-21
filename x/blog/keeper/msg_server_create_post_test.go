package keeper_test

import (
	"fmt"
	"testing"

	"blogchain/x/blog/types"
	"github.com/stretchr/testify/require"
)

func TestMsgServerCreatePost(t *testing.T) {
	f := initFixture(t)
	require := require.New(t)
	
	creator := GetTestAddress(t)
	
	tests := []struct {
		name      string
		msg       *types.MsgCreatePost
		shouldErr bool
		errContains string
	}{
		{
			name: "valid post creation",
			msg: &types.MsgCreatePost{
				Creator: creator,
				Title:   "Test Post",
				Body:    "This is a test post body",
				Tags:    []string{"test", "blockchain"},
			},
			shouldErr: false,
		},
		{
			name: "empty title",
			msg: &types.MsgCreatePost{
				Creator: creator,
				Title:   "",
				Body:    "This is a test post body",
				Tags:    []string{"test"},
			},
			shouldErr: false, // Currently allows empty titles
		},
		{
			name: "empty body",
			msg: &types.MsgCreatePost{
				Creator: creator,
				Title:   "Test Post",
				Body:    "",
				Tags:    []string{"test"},
			},
			shouldErr: false, // Currently allows empty body
		},
		{
			name: "no tags",
			msg: &types.MsgCreatePost{
				Creator: creator,
				Title:   "Test Post",
				Body:    "This is a test post body",
				Tags:    []string{},
			},
			shouldErr: false,
		},
		{
			name: "multiple tags",
			msg: &types.MsgCreatePost{
				Creator: creator,
				Title:   "Test Post",
				Body:    "This is a test post body",
				Tags:    []string{"blockchain", "cosmos", "web3", "decentralized"},
			},
			shouldErr: false,
		},
	}
	
	msgServer := NewMsgServerImpl(f.keeper)
	
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			resp, err := msgServer.CreatePost(f.ctx, tc.msg)
			
			if tc.shouldErr {
				require.Error(err)
				if tc.errContains != "" {
					require.Contains(err.Error(), tc.errContains)
				}
			} else {
				require.NoError(err)
				require.NotNil(resp)
				require.NotZero(resp.Id)
				
				// Verify post was actually created
				post, found := f.keeper.GetPost(f.ctx, resp.Id)
				require.True(found)
				require.Equal(tc.msg.Creator, post.Creator)
				require.Equal(tc.msg.Title, post.Title)
				require.Equal(tc.msg.Body, post.Body)
				if len(tc.msg.Tags) == 0 {
					// Empty slice might be stored as nil
					require.Empty(post.Tags)
				} else {
					require.Equal(tc.msg.Tags, post.Tags)
				}
				require.NotZero(post.CreatedAt)
				require.Equal(uint64(0), post.Likes)
			}
		})
	}
}

func TestMsgServerCreatePostValidation(t *testing.T) {
	f := initFixture(t)
	require := require.New(t)
	
	msgServer := NewMsgServerImpl(f.keeper)
	
	// Test invalid creator address
	msg := &types.MsgCreatePost{
		Creator: "invalid-address",
		Title:   "Test Post",
		Body:    "Test body",
		Tags:    []string{"test"},
	}
	
	_, err := msgServer.CreatePost(f.ctx, msg)
	require.Error(err)
	require.Contains(err.Error(), "invalid creator address")
}

func TestMsgServerCreateMultiplePosts(t *testing.T) {
	f := initFixture(t)
	require := require.New(t)
	
	creator := GetTestAddress(t)
	msgServer := NewMsgServerImpl(f.keeper)
	
	// Create multiple posts
	postCount := 5
	postIds := make([]uint64, postCount)
	
	for i := 0; i < postCount; i++ {
		msg := &types.MsgCreatePost{
			Creator: creator,
			Title:   fmt.Sprintf("Post %d", i+1),
			Body:    fmt.Sprintf("This is post number %d", i+1),
			Tags:    []string{fmt.Sprintf("tag%d", i)},
		}
		
		resp, err := msgServer.CreatePost(f.ctx, msg)
		require.NoError(err)
		require.NotNil(resp)
		postIds[i] = resp.Id
	}
	
	// Verify all posts exist and have sequential IDs
	for i, id := range postIds {
		post, found := f.keeper.GetPost(f.ctx, id)
		require.True(found)
		require.Equal(fmt.Sprintf("Post %d", i+1), post.Title)
		
		// IDs should be sequential
		if i > 0 {
			require.Equal(postIds[i-1]+1, id)
		}
	}
	
	// Verify all posts can be retrieved
	// Note: GetAllPost method doesn't exist, we verified individual posts above
}