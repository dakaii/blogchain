package keeper

import (
	"context"
	"errors"
	"fmt"

	"blogchain/x/blog/types"

	"cosmossdk.io/collections"
)

// CreateProfile creates a new user profile
func (k Keeper) CreateProfile(ctx context.Context, profile types.Profile) error {
	// Check if profile already exists for this address
	_, err := k.Profiles.Get(ctx, profile.Address)
	if err == nil {
		return fmt.Errorf("profile already exists for address %s", profile.Address)
	}
	if !errors.Is(err, collections.ErrNotFound) {
		return fmt.Errorf("failed to check existing profile: %w", err)
	}

	// Check username uniqueness
	taken, _ := k.UsernameToAddress.Get(ctx, profile.Username)
	if taken != "" {
		return fmt.Errorf("username %s is already taken", profile.Username)
	}

	// Save the profile
	if err := k.Profiles.Set(ctx, profile.Address, profile); err != nil {
		return fmt.Errorf("failed to save profile: %w", err)
	}

	// Save username mapping
	if err := k.UsernameToAddress.Set(ctx, profile.Username, profile.Address); err != nil {
		return fmt.Errorf("failed to save username mapping: %w", err)
	}

	return nil
}

// GetProfile retrieves a profile by address
func (k Keeper) GetProfile(ctx context.Context, address string) (types.Profile, error) {
	profile, err := k.Profiles.Get(ctx, address)
	if err != nil {
		if errors.Is(err, collections.ErrNotFound) {
			return types.Profile{}, fmt.Errorf("profile not found for address %s", address)
		}
		return types.Profile{}, fmt.Errorf("failed to get profile: %w", err)
	}
	return profile, nil
}

// GetProfileByUsername retrieves a profile by username
func (k Keeper) GetProfileByUsername(ctx context.Context, username string) (types.Profile, error) {
	address, err := k.UsernameToAddress.Get(ctx, username)
	if err != nil {
		if errors.Is(err, collections.ErrNotFound) {
			return types.Profile{}, fmt.Errorf("profile not found for username %s", username)
		}
		return types.Profile{}, fmt.Errorf("failed to get username mapping: %w", err)
	}
	return k.GetProfile(ctx, address)
}

// UpdateProfile updates an existing profile
func (k Keeper) UpdateProfile(ctx context.Context, address string, displayName, bio, avatarURL, website string, updatedAt int64) error {
	profile, err := k.GetProfile(ctx, address)
	if err != nil {
		return err
	}

	// Update fields
	profile.DisplayName = displayName
	profile.Bio = bio
	profile.AvatarUrl = avatarURL
	profile.Website = website
	profile.UpdatedAt = updatedAt

	// Save updated profile
	if err := k.Profiles.Set(ctx, address, profile); err != nil {
		return fmt.Errorf("failed to update profile: %w", err)
	}

	return nil
}

// Follow creates a follow relationship
func (k Keeper) Follow(ctx context.Context, followerAddr, followingAddr string) error {
	if followerAddr == followingAddr {
		return fmt.Errorf("cannot follow yourself")
	}

	// Verify both profiles exist
	followerProfile, err := k.GetProfile(ctx, followerAddr)
	if err != nil {
		return fmt.Errorf("follower profile not found: %w", err)
	}

	followingProfile, err := k.GetProfile(ctx, followingAddr)
	if err != nil {
		return fmt.Errorf("following profile not found: %w", err)
	}

	// Check if already following
	followKey := collections.Join(followerAddr, followingAddr)
	isFollowing, _ := k.Follows.Get(ctx, followKey)
	if isFollowing {
		return fmt.Errorf("already following this user")
	}

	// Create follow relationship
	if err := k.Follows.Set(ctx, followKey, true); err != nil {
		return fmt.Errorf("failed to create follow relationship: %w", err)
	}

	// Update follower count for the followed user
	followingProfile.Followers++
	if err := k.Profiles.Set(ctx, followingAddr, followingProfile); err != nil {
		return fmt.Errorf("failed to update follower count: %w", err)
	}

	// Update following count for the follower
	followerProfile.Following++
	if err := k.Profiles.Set(ctx, followerAddr, followerProfile); err != nil {
		return fmt.Errorf("failed to update following count: %w", err)
	}

	return nil
}

// Unfollow removes a follow relationship
func (k Keeper) Unfollow(ctx context.Context, followerAddr, followingAddr string) error {
	if followerAddr == followingAddr {
		return fmt.Errorf("cannot unfollow yourself")
	}

	// Check if following
	followKey := collections.Join(followerAddr, followingAddr)
	isFollowing, err := k.Follows.Get(ctx, followKey)
	if err != nil || !isFollowing {
		return fmt.Errorf("not following this user")
	}

	// Remove follow relationship
	if err := k.Follows.Remove(ctx, followKey); err != nil {
		return fmt.Errorf("failed to remove follow relationship: %w", err)
	}

	// Update follower count for the unfollowed user
	followingProfile, err := k.GetProfile(ctx, followingAddr)
	if err == nil && followingProfile.Followers > 0 {
		followingProfile.Followers--
		if err := k.Profiles.Set(ctx, followingAddr, followingProfile); err != nil {
			return fmt.Errorf("failed to update follower count: %w", err)
		}
	}

	// Update following count for the unfollower
	followerProfile, err := k.GetProfile(ctx, followerAddr)
	if err == nil && followerProfile.Following > 0 {
		followerProfile.Following--
		if err := k.Profiles.Set(ctx, followerAddr, followerProfile); err != nil {
			return fmt.Errorf("failed to update following count: %w", err)
		}
	}

	return nil
}

// IsFollowing checks if one user is following another
func (k Keeper) IsFollowing(ctx context.Context, followerAddr, followingAddr string) (bool, error) {
	followKey := collections.Join(followerAddr, followingAddr)
	isFollowing, err := k.Follows.Get(ctx, followKey)
	if err != nil {
		if errors.Is(err, collections.ErrNotFound) {
			return false, nil
		}
		return false, fmt.Errorf("failed to check follow relationship: %w", err)
	}
	return isFollowing, nil
}

// GetFollowers retrieves all followers of a user
func (k Keeper) GetFollowers(ctx context.Context, address string) ([]string, error) {
	var followers []string
	
	// Walk through all follows to find followers
	err := k.Follows.Walk(ctx, nil, func(key collections.Pair[string, string], _ bool) (stop bool, err error) {
		follower := key.K1()
		following := key.K2()
		if following == address {
			followers = append(followers, follower)
		}
		return false, nil
	})
	
	if err != nil {
		return nil, fmt.Errorf("failed to get followers: %w", err)
	}
	
	return followers, nil
}

// GetFollowing retrieves all users that a user is following
func (k Keeper) GetFollowing(ctx context.Context, address string) ([]string, error) {
	var following []string
	
	// Walk through follows with prefix
	err := k.Follows.Walk(ctx, collections.NewPrefixedPairRange[string, string](address), 
		func(key collections.Pair[string, string], _ bool) (stop bool, err error) {
			followingAddr := key.K2()
			following = append(following, followingAddr)
			return false, nil
		})
	
	if err != nil {
		return nil, fmt.Errorf("failed to get following: %w", err)
	}
	
	return following, nil
}

// GetAllProfiles retrieves all profiles
func (k Keeper) GetAllProfiles(ctx context.Context) ([]types.Profile, error) {
	var profiles []types.Profile
	
	err := k.Profiles.Walk(ctx, nil, func(address string, profile types.Profile) (stop bool, err error) {
		profiles = append(profiles, profile)
		return false, nil
	})
	
	if err != nil {
		return nil, fmt.Errorf("failed to get all profiles: %w", err)
	}
	
	return profiles, nil
}

// IncrementPostCount increments the post count for a user profile
func (k Keeper) IncrementPostCount(ctx context.Context, address string) error {
	profile, err := k.GetProfile(ctx, address)
	if err != nil {
		// Profile doesn't exist yet, that's okay
		return nil
	}
	
	profile.PostCount++
	if err := k.Profiles.Set(ctx, address, profile); err != nil {
		return fmt.Errorf("failed to update post count: %w", err)
	}
	
	return nil
}

// DecrementPostCount decrements the post count for a user profile
func (k Keeper) DecrementPostCount(ctx context.Context, address string) error {
	profile, err := k.GetProfile(ctx, address)
	if err != nil {
		// Profile doesn't exist, that's okay
		return nil
	}
	
	if profile.PostCount > 0 {
		profile.PostCount--
		if err := k.Profiles.Set(ctx, address, profile); err != nil {
			return fmt.Errorf("failed to update post count: %w", err)
		}
	}
	
	return nil
}