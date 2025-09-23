package keeper

import (
	"context"
	"regexp"
	"strings"

	"blogchain/x/blog/types"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var (
	// Username must be 3-20 characters, alphanumeric and underscore only
	usernameRegex = regexp.MustCompile(`^[a-zA-Z0-9_]{3,20}$`)
	// Simple URL validation
	urlRegex = regexp.MustCompile(`^(https?://)?([\da-z\.-]+)\.([a-z\.]{2,6})([/\w \.-]*)*/?$`)
)

func (k msgServer) CreateProfile(ctx context.Context, msg *types.MsgCreateProfile) (*types.MsgCreateProfileResponse, error) {
	if _, err := k.addressCodec.StringToBytes(msg.Creator); err != nil {
		return nil, errorsmod.Wrap(err, "invalid creator address")
	}

	// Validate username
	if !usernameRegex.MatchString(msg.Username) {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "invalid username format (3-20 chars, alphanumeric and underscore only)")
	}

	// Validate display name length
	if len(msg.DisplayName) > 50 {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "display name too long (max 50 characters)")
	}

	// Validate bio length
	if len(msg.Bio) > 500 {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "bio too long (max 500 characters)")
	}

	// Validate URLs if provided
	if msg.AvatarUrl != "" && len(msg.AvatarUrl) > 500 {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "avatar URL too long (max 500 characters)")
	}

	if msg.Website != "" {
		if len(msg.Website) > 200 {
			return nil, errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "website URL too long (max 200 characters)")
		}
		if !urlRegex.MatchString(msg.Website) {
			return nil, errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "invalid website URL format")
		}
	}

	sdkCtx := sdk.UnwrapSDKContext(ctx)

	profile := types.Profile{
		Address:     msg.Creator,
		Username:    strings.ToLower(msg.Username), // Store usernames in lowercase
		DisplayName: msg.DisplayName,
		Bio:         msg.Bio,
		AvatarUrl:   msg.AvatarUrl,
		Website:     msg.Website,
		CreatedAt:   sdkCtx.BlockTime().Unix(),
		UpdatedAt:   sdkCtx.BlockTime().Unix(),
		Followers:   0,
		Following:   0,
		PostCount:   0,
		Verified:    false,
	}

	if err := k.Keeper.CreateProfile(ctx, profile); err != nil {
		return nil, errorsmod.Wrap(err, "failed to create profile")
	}

	// Emit event
	sdkCtx.EventManager().EmitEvent(
		sdk.NewEvent(
			"profile_created",
			sdk.NewAttribute("address", msg.Creator),
			sdk.NewAttribute("username", profile.Username),
		),
	)

	return &types.MsgCreateProfileResponse{}, nil
}

func (k msgServer) UpdateProfile(ctx context.Context, msg *types.MsgUpdateProfile) (*types.MsgUpdateProfileResponse, error) {
	if _, err := k.addressCodec.StringToBytes(msg.Creator); err != nil {
		return nil, errorsmod.Wrap(err, "invalid creator address")
	}

	// Validate display name length
	if len(msg.DisplayName) > 50 {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "display name too long (max 50 characters)")
	}

	// Validate bio length
	if len(msg.Bio) > 500 {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "bio too long (max 500 characters)")
	}

	// Validate URLs if provided
	if msg.AvatarUrl != "" && len(msg.AvatarUrl) > 500 {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "avatar URL too long (max 500 characters)")
	}

	if msg.Website != "" {
		if len(msg.Website) > 200 {
			return nil, errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "website URL too long (max 200 characters)")
		}
		if !urlRegex.MatchString(msg.Website) {
			return nil, errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "invalid website URL format")
		}
	}

	sdkCtx := sdk.UnwrapSDKContext(ctx)

	if err := k.Keeper.UpdateProfile(ctx, msg.Creator, msg.DisplayName, msg.Bio, 
		msg.AvatarUrl, msg.Website, sdkCtx.BlockTime().Unix()); err != nil {
		return nil, errorsmod.Wrap(err, "failed to update profile")
	}

	// Emit event
	sdkCtx.EventManager().EmitEvent(
		sdk.NewEvent(
			"profile_updated",
			sdk.NewAttribute("address", msg.Creator),
		),
	)

	return &types.MsgUpdateProfileResponse{}, nil
}

func (k msgServer) Follow(ctx context.Context, msg *types.MsgFollow) (*types.MsgFollowResponse, error) {
	if _, err := k.addressCodec.StringToBytes(msg.Follower); err != nil {
		return nil, errorsmod.Wrap(err, "invalid follower address")
	}

	if _, err := k.addressCodec.StringToBytes(msg.Following); err != nil {
		return nil, errorsmod.Wrap(err, "invalid following address")
	}

	if err := k.Keeper.Follow(ctx, msg.Follower, msg.Following); err != nil {
		return nil, errorsmod.Wrap(err, "failed to follow user")
	}

	// Emit event
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	sdkCtx.EventManager().EmitEvent(
		sdk.NewEvent(
			"user_followed",
			sdk.NewAttribute("follower", msg.Follower),
			sdk.NewAttribute("following", msg.Following),
		),
	)

	return &types.MsgFollowResponse{}, nil
}

func (k msgServer) Unfollow(ctx context.Context, msg *types.MsgUnfollow) (*types.MsgUnfollowResponse, error) {
	if _, err := k.addressCodec.StringToBytes(msg.Follower); err != nil {
		return nil, errorsmod.Wrap(err, "invalid follower address")
	}

	if _, err := k.addressCodec.StringToBytes(msg.Following); err != nil {
		return nil, errorsmod.Wrap(err, "invalid following address")
	}

	if err := k.Keeper.Unfollow(ctx, msg.Follower, msg.Following); err != nil {
		return nil, errorsmod.Wrap(err, "failed to unfollow user")
	}

	// Emit event
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	sdkCtx.EventManager().EmitEvent(
		sdk.NewEvent(
			"user_unfollowed",
			sdk.NewAttribute("follower", msg.Follower),
			sdk.NewAttribute("following", msg.Following),
		),
	)

	return &types.MsgUnfollowResponse{}, nil
}