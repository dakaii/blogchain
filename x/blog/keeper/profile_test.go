package keeper_test

import (
	"testing"

	"blogchain/x/blog/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

func TestProfileOperations(t *testing.T) {
	f := initFixture(t)
	k := f.keeper
	ctx := sdk.UnwrapSDKContext(f.ctx)
	msgServer := NewMsgServerImpl(k)

	user1 := GetTestAddress(t)
	user2 := GetTestAddress(t)
	user3 := GetTestAddress(t)

	t.Run("Create profile", func(t *testing.T) {
		msg := &types.MsgCreateProfile{
			Creator:     user1,
			Username:    "alice_crypto",
			DisplayName: "Alice",
			Bio:         "Blockchain enthusiast",
			AvatarUrl:   "https://example.com/avatar.png",
			Website:     "https://alice.example.com",
		}

		_, err := msgServer.CreateProfile(f.ctx, msg)
		require.NoError(t, err)

		// Verify profile was created
		profile, err := k.GetProfile(ctx, user1)
		require.NoError(t, err)
		require.Equal(t, user1, profile.Address)
		require.Equal(t, "alice_crypto", profile.Username) // Should be lowercase
		require.Equal(t, "Alice", profile.DisplayName)
		require.Equal(t, "Blockchain enthusiast", profile.Bio)
		require.Equal(t, uint64(0), profile.Followers)
		require.Equal(t, uint64(0), profile.Following)
		require.Equal(t, uint64(0), profile.PostCount)
		require.False(t, profile.Verified)

		// Verify username lookup works
		profileByUsername, err := k.GetProfileByUsername(ctx, "alice_crypto")
		require.NoError(t, err)
		require.Equal(t, profile.Address, profileByUsername.Address)
	})

	t.Run("Username must be unique", func(t *testing.T) {
		// First user creates profile
		msg1 := &types.MsgCreateProfile{
			Creator:     user2,
			Username:    "bob_unique",
			DisplayName: "Bob",
			Bio:         "First Bob",
		}
		_, err := msgServer.CreateProfile(f.ctx, msg1)
		require.NoError(t, err)

		// Second user tries same username
		msg2 := &types.MsgCreateProfile{
			Creator:     user3,
			Username:    "bob_unique",
			DisplayName: "Bob 2",
			Bio:         "Second Bob",
		}
		_, err = msgServer.CreateProfile(f.ctx, msg2)
		require.Error(t, err)
		require.Contains(t, err.Error(), "username bob_unique is already taken")
	})

	t.Run("Username validation", func(t *testing.T) {
		testCases := []struct {
			name     string
			username string
			valid    bool
		}{
			{"too short", "ab", false},
			{"too long", "this_username_is_way_too_long_for_validation", false},
			{"special chars", "alice@crypto", false},
			{"spaces", "alice crypto", false},
			{"valid", "alice_123", true},
			{"valid underscore", "alice_crypto_", true},
			{"valid numbers", "alice123", true},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				// Generate unique address for each test
				testAddr := GetTestAddress(t)
				msg := &types.MsgCreateProfile{
					Creator:     testAddr,
					Username:    tc.username,
					DisplayName: "Test User",
					Bio:         "Test bio",
				}
				_, err := msgServer.CreateProfile(f.ctx, msg)
				if tc.valid {
					require.NoError(t, err)
				} else {
					require.Error(t, err)
				}
			})
		}
	})

	t.Run("Cannot create duplicate profile", func(t *testing.T) {
		// Create first profile
		addr := GetTestAddress(t)
		msg1 := &types.MsgCreateProfile{
			Creator:     addr,
			Username:    "unique_user",
			DisplayName: "First",
		}
		_, err := msgServer.CreateProfile(f.ctx, msg1)
		require.NoError(t, err)

		// Try to create another profile for same address
		msg2 := &types.MsgCreateProfile{
			Creator:     addr,
			Username:    "different_username",
			DisplayName: "Second",
		}
		_, err = msgServer.CreateProfile(f.ctx, msg2)
		require.Error(t, err)
		require.Contains(t, err.Error(), "profile already exists")
	})

	t.Run("Update profile", func(t *testing.T) {
		// Create profile
		addr := GetTestAddress(t)
		createMsg := &types.MsgCreateProfile{
			Creator:     addr,
			Username:    "updatable_user",
			DisplayName: "Original Name",
			Bio:         "Original bio",
		}
		_, err := msgServer.CreateProfile(f.ctx, createMsg)
		require.NoError(t, err)

		// Update profile
		updateMsg := &types.MsgUpdateProfile{
			Creator:     addr,
			DisplayName: "Updated Name",
			Bio:         "Updated bio with more details",
			AvatarUrl:   "https://new-avatar.com/image.png",
			Website:     "https://updated.example.com",
		}
		_, err = msgServer.UpdateProfile(f.ctx, updateMsg)
		require.NoError(t, err)

		// Verify update
		profile, err := k.GetProfile(ctx, addr)
		require.NoError(t, err)
		require.Equal(t, "Updated Name", profile.DisplayName)
		require.Equal(t, "Updated bio with more details", profile.Bio)
		require.Equal(t, "https://new-avatar.com/image.png", profile.AvatarUrl)
		require.Equal(t, "https://updated.example.com", profile.Website)
		// Username should remain unchanged
		require.Equal(t, "updatable_user", profile.Username)
		require.GreaterOrEqual(t, profile.UpdatedAt, profile.CreatedAt)
	})

	t.Run("Follow and unfollow", func(t *testing.T) {
		// Create profiles for follower and following
		followerAddr := GetTestAddress(t)
		followingAddr := GetTestAddress(t)

		createFollower := &types.MsgCreateProfile{
			Creator:     followerAddr,
			Username:    "follower_user",
			DisplayName: "Follower",
		}
		_, err := msgServer.CreateProfile(f.ctx, createFollower)
		require.NoError(t, err)

		createFollowing := &types.MsgCreateProfile{
			Creator:     followingAddr,
			Username:    "following_user",
			DisplayName: "Following",
		}
		_, err = msgServer.CreateProfile(f.ctx, createFollowing)
		require.NoError(t, err)

		// Follow
		followMsg := &types.MsgFollow{
			Follower:  followerAddr,
			Following: followingAddr,
		}
		_, err = msgServer.Follow(f.ctx, followMsg)
		require.NoError(t, err)

		// Verify follow relationship
		isFollowing, err := k.IsFollowing(ctx, followerAddr, followingAddr)
		require.NoError(t, err)
		require.True(t, isFollowing)

		// Verify counts updated
		followerProfile, err := k.GetProfile(ctx, followerAddr)
		require.NoError(t, err)
		require.Equal(t, uint64(1), followerProfile.Following)

		followingProfile, err := k.GetProfile(ctx, followingAddr)
		require.NoError(t, err)
		require.Equal(t, uint64(1), followingProfile.Followers)

		// Cannot follow twice
		_, err = msgServer.Follow(f.ctx, followMsg)
		require.Error(t, err)
		require.Contains(t, err.Error(), "already following")

		// Unfollow
		unfollowMsg := &types.MsgUnfollow{
			Follower:  followerAddr,
			Following: followingAddr,
		}
		_, err = msgServer.Unfollow(f.ctx, unfollowMsg)
		require.NoError(t, err)

		// Verify unfollowed
		isFollowing, err = k.IsFollowing(ctx, followerAddr, followingAddr)
		require.NoError(t, err)
		require.False(t, isFollowing)

		// Verify counts updated
		followerProfile, err = k.GetProfile(ctx, followerAddr)
		require.NoError(t, err)
		require.Equal(t, uint64(0), followerProfile.Following)

		followingProfile, err = k.GetProfile(ctx, followingAddr)
		require.NoError(t, err)
		require.Equal(t, uint64(0), followingProfile.Followers)

		// Cannot unfollow if not following
		_, err = msgServer.Unfollow(f.ctx, unfollowMsg)
		require.Error(t, err)
		require.Contains(t, err.Error(), "not following")
	})

	t.Run("Cannot follow yourself", func(t *testing.T) {
		addr := GetTestAddress(t)
		createMsg := &types.MsgCreateProfile{
			Creator:     addr,
			Username:    "self_follower",
			DisplayName: "Self",
		}
		_, err := msgServer.CreateProfile(f.ctx, createMsg)
		require.NoError(t, err)

		followMsg := &types.MsgFollow{
			Follower:  addr,
			Following: addr,
		}
		_, err = msgServer.Follow(f.ctx, followMsg)
		require.Error(t, err)
		require.Contains(t, err.Error(), "cannot follow yourself")
	})

	t.Run("Get followers and following", func(t *testing.T) {
		// Create a central user
		centralAddr := GetTestAddress(t)
		createCentral := &types.MsgCreateProfile{
			Creator:     centralAddr,
			Username:    "central_user",
			DisplayName: "Central",
		}
		_, err := msgServer.CreateProfile(f.ctx, createCentral)
		require.NoError(t, err)

		// Create followers
		var followerAddrs []string
		for i := 0; i < 3; i++ {
			addr := GetTestAddress(t)
			followerAddrs = append(followerAddrs, addr)
			
			createMsg := &types.MsgCreateProfile{
				Creator:     addr,
				Username:    GetUniqueUsername(t),
				DisplayName: "Follower",
			}
			_, err := msgServer.CreateProfile(f.ctx, createMsg)
			require.NoError(t, err)

			followMsg := &types.MsgFollow{
				Follower:  addr,
				Following: centralAddr,
			}
			_, err = msgServer.Follow(f.ctx, followMsg)
			require.NoError(t, err)
		}

		// Create users that central follows
		var followingAddrs []string
		for i := 0; i < 2; i++ {
			addr := GetTestAddress(t)
			followingAddrs = append(followingAddrs, addr)
			
			createMsg := &types.MsgCreateProfile{
				Creator:     addr,
				Username:    GetUniqueUsername(t),
				DisplayName: "Following",
			}
			_, err := msgServer.CreateProfile(f.ctx, createMsg)
			require.NoError(t, err)

			followMsg := &types.MsgFollow{
				Follower:  centralAddr,
				Following: addr,
			}
			_, err = msgServer.Follow(f.ctx, followMsg)
			require.NoError(t, err)
		}

		// Get followers
		followers, err := k.GetFollowers(ctx, centralAddr)
		require.NoError(t, err)
		require.Len(t, followers, 3)

		// Get following
		following, err := k.GetFollowing(ctx, centralAddr)
		require.NoError(t, err)
		require.Len(t, following, 2)
	})

	t.Run("Post count tracking", func(t *testing.T) {
		// Create profile
		addr := GetTestAddress(t)
		createProfileMsg := &types.MsgCreateProfile{
			Creator:     addr,
			Username:    "post_counter",
			DisplayName: "Post Counter",
		}
		_, err := msgServer.CreateProfile(f.ctx, createProfileMsg)
		require.NoError(t, err)

		// Create posts
		for i := 0; i < 3; i++ {
			createPostMsg := &types.MsgCreatePost{
				Creator: addr,
				Title:   "Test Post",
				Body:    "Test content",
				Tags:    []string{"test"},
			}
			_, err := msgServer.CreatePost(f.ctx, createPostMsg)
			require.NoError(t, err)
		}

		// Verify post count
		profile, err := k.GetProfile(ctx, addr)
		require.NoError(t, err)
		require.Equal(t, uint64(3), profile.PostCount)

		// Delete a post
		deletePostMsg := &types.MsgDeletePost{
			Creator: addr,
			Id:      0, // First post
		}
		_, err = msgServer.DeletePost(f.ctx, deletePostMsg)
		require.NoError(t, err)

		// Verify post count decreased
		profile, err = k.GetProfile(ctx, addr)
		require.NoError(t, err)
		require.Equal(t, uint64(2), profile.PostCount)
	})
}

// Helper function to generate unique usernames
var usernameCounter int

func GetUniqueUsername(t *testing.T) string {
	t.Helper()
	usernameCounter++
	// Use counter to ensure uniqueness
	return "user_test_" + string(rune('a'+usernameCounter))
}