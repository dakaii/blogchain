package keeper

import (
	"context"

	"blogchain/x/blog/types"

	errorsmod "cosmossdk.io/errors"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

func (k msgServer) LikePost(ctx context.Context, msg *types.MsgLikePost) (*types.MsgLikePostResponse, error) {
	if _, err := k.addressCodec.StringToBytes(msg.Liker); err != nil {
		return nil, errorsmod.Wrap(err, "invalid liker address")
	}

	// Check if user has already liked this post
	alreadyLiked, err := k.HasUserLikedPost(ctx, msg.PostId, msg.Liker)
	if err != nil {
		return nil, errorsmod.Wrap(err, "failed to check if user liked post")
	}
	if alreadyLiked {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "user has already liked this post")
	}

	post, err := k.GetPost(ctx, msg.PostId)
	if err != nil {
		return nil, errorsmod.Wrap(err, "failed to get post")
	}

	post.Likes++
	if err := k.SetPost(ctx, post); err != nil {
		return nil, errorsmod.Wrap(err, "failed to update post")
	}

	// Mark that user has liked the post
	if err := k.SetUserLikedPost(ctx, msg.PostId, msg.Liker); err != nil {
		return nil, errorsmod.Wrap(err, "failed to record user like")
	}

	return &types.MsgLikePostResponse{}, nil
}