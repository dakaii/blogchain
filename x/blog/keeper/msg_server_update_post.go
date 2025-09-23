package keeper

import (
	"context"
	"fmt"

	"blogchain/x/blog/types"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

func (k msgServer) UpdatePost(ctx context.Context, msg *types.MsgUpdatePost) (*types.MsgUpdatePostResponse, error) {
	if _, err := k.addressCodec.StringToBytes(msg.Creator); err != nil {
		return nil, errorsmod.Wrap(err, "invalid creator address")
	}

	// Get the existing post
	post, err := k.GetPost(ctx, msg.Id)
	if err != nil {
		return nil, errorsmod.Wrap(err, "post not found")
	}

	// Check if the message sender is the post creator
	if post.Creator != msg.Creator {
		return nil, errorsmod.Wrap(sdkerrors.ErrUnauthorized, "only the post creator can update the post")
	}

	// Check if post is deleted
	if post.Deleted {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "cannot update a deleted post")
	}

	sdkCtx := sdk.UnwrapSDKContext(ctx)

	// Update the post fields
	post.Title = msg.Title
	post.Body = msg.Body
	post.Tags = msg.Tags
	post.UpdatedAt = sdkCtx.BlockTime().Unix()
	post.MediaBlobIds = msg.MediaBlobIds
	post.ContentBlobId = msg.ContentBlobId

	// Save the updated post
	if err := k.SetPost(ctx, post); err != nil {
		return nil, errorsmod.Wrap(err, "failed to update post")
	}

	// Emit event
	sdkCtx.EventManager().EmitEvent(
		sdk.NewEvent(
			"post_updated",
			sdk.NewAttribute("id", fmt.Sprintf("%d", post.Id)),
			sdk.NewAttribute("creator", post.Creator),
			sdk.NewAttribute("updated_at", fmt.Sprintf("%d", post.UpdatedAt)),
		),
	)

	return &types.MsgUpdatePostResponse{}, nil
}