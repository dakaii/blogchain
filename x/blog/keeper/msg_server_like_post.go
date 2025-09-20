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

	post, found := k.GetPost(ctx, msg.PostId)
	if !found {
		return nil, errorsmod.Wrap(sdkerrors.ErrNotFound, "post not found")
	}

	post.Likes++
	k.SetPost(ctx, post)

	return &types.MsgLikePostResponse{}, nil
}