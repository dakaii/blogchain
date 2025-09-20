package keeper

import (
	"context"

	"blogchain/x/blog/types"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k msgServer) CreatePost(ctx context.Context, msg *types.MsgCreatePost) (*types.MsgCreatePostResponse, error) {
	if _, err := k.addressCodec.StringToBytes(msg.Creator); err != nil {
		return nil, errorsmod.Wrap(err, "invalid creator address")
	}

	sdkCtx := sdk.UnwrapSDKContext(ctx)
	
	post := types.Post{
		Creator:   msg.Creator,
		Title:     msg.Title,
		Body:      msg.Body,
		Tags:      msg.Tags,
		CreatedAt: sdkCtx.BlockTime().Unix(),
		Likes:     0,
	}

	id := k.AppendPost(ctx, post)

	return &types.MsgCreatePostResponse{
		Id: id,
	}, nil
}
