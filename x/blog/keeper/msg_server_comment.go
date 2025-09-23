package keeper

import (
	"context"
	"fmt"

	"blogchain/x/blog/types"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

func (k msgServer) CreateComment(ctx context.Context, msg *types.MsgCreateComment) (*types.MsgCreateCommentResponse, error) {
	if _, err := k.addressCodec.StringToBytes(msg.Creator); err != nil {
		return nil, errorsmod.Wrap(err, "invalid creator address")
	}

	// Validate content
	if len(msg.Content) == 0 {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "comment content cannot be empty")
	}
	if len(msg.Content) > 5000 { // Max 5000 characters
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "comment content too long (max 5000 characters)")
	}

	sdkCtx := sdk.UnwrapSDKContext(ctx)

	comment := types.Comment{
		PostId:    msg.PostId,
		ParentId:  msg.ParentId,
		Creator:   msg.Creator,
		Content:   msg.Content,
		CreatedAt: sdkCtx.BlockTime().Unix(),
		UpdatedAt: sdkCtx.BlockTime().Unix(),
		Likes:     0,
		Deleted:   false,
	}

	id, err := k.Keeper.CreateComment(ctx, comment)
	if err != nil {
		return nil, errorsmod.Wrap(err, "failed to create comment")
	}

	// Emit event
	sdkCtx.EventManager().EmitEvent(
		sdk.NewEvent(
			"comment_created",
			sdk.NewAttribute("id", fmt.Sprintf("%d", id)),
			sdk.NewAttribute("post_id", fmt.Sprintf("%d", msg.PostId)),
			sdk.NewAttribute("parent_id", fmt.Sprintf("%d", msg.ParentId)),
			sdk.NewAttribute("creator", msg.Creator),
		),
	)

	return &types.MsgCreateCommentResponse{Id: id}, nil
}

func (k msgServer) UpdateComment(ctx context.Context, msg *types.MsgUpdateComment) (*types.MsgUpdateCommentResponse, error) {
	if _, err := k.addressCodec.StringToBytes(msg.Creator); err != nil {
		return nil, errorsmod.Wrap(err, "invalid creator address")
	}

	// Validate content
	if len(msg.Content) == 0 {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "comment content cannot be empty")
	}
	if len(msg.Content) > 5000 {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "comment content too long (max 5000 characters)")
	}

	// Get existing comment
	comment, err := k.Keeper.GetComment(ctx, msg.Id)
	if err != nil {
		return nil, errorsmod.Wrap(err, "comment not found")
	}

	// Check ownership
	if comment.Creator != msg.Creator {
		return nil, errorsmod.Wrap(sdkerrors.ErrUnauthorized, "only comment creator can update")
	}

	if comment.Deleted {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "cannot update deleted comment")
	}

	sdkCtx := sdk.UnwrapSDKContext(ctx)

	// Update comment
	comment.Content = msg.Content
	comment.UpdatedAt = sdkCtx.BlockTime().Unix()

	if err := k.Keeper.UpdateComment(ctx, comment); err != nil {
		return nil, errorsmod.Wrap(err, "failed to update comment")
	}

	// Emit event
	sdkCtx.EventManager().EmitEvent(
		sdk.NewEvent(
			"comment_updated",
			sdk.NewAttribute("id", fmt.Sprintf("%d", msg.Id)),
			sdk.NewAttribute("creator", msg.Creator),
			sdk.NewAttribute("updated_at", fmt.Sprintf("%d", comment.UpdatedAt)),
		),
	)

	return &types.MsgUpdateCommentResponse{}, nil
}

func (k msgServer) DeleteComment(ctx context.Context, msg *types.MsgDeleteComment) (*types.MsgDeleteCommentResponse, error) {
	if _, err := k.addressCodec.StringToBytes(msg.Creator); err != nil {
		return nil, errorsmod.Wrap(err, "invalid creator address")
	}

	// Get existing comment
	comment, err := k.Keeper.GetComment(ctx, msg.Id)
	if err != nil {
		return nil, errorsmod.Wrap(err, "comment not found")
	}

	// Check ownership
	if comment.Creator != msg.Creator {
		return nil, errorsmod.Wrap(sdkerrors.ErrUnauthorized, "only comment creator can delete")
	}

	if comment.Deleted {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "comment already deleted")
	}

	sdkCtx := sdk.UnwrapSDKContext(ctx)

	// Soft delete
	if err := k.Keeper.DeleteComment(ctx, msg.Id, sdkCtx.BlockTime().Unix()); err != nil {
		return nil, errorsmod.Wrap(err, "failed to delete comment")
	}

	// Emit event
	sdkCtx.EventManager().EmitEvent(
		sdk.NewEvent(
			"comment_deleted",
			sdk.NewAttribute("id", fmt.Sprintf("%d", msg.Id)),
			sdk.NewAttribute("creator", msg.Creator),
		),
	)

	return &types.MsgDeleteCommentResponse{}, nil
}

func (k msgServer) LikeComment(ctx context.Context, msg *types.MsgLikeComment) (*types.MsgLikeCommentResponse, error) {
	if _, err := k.addressCodec.StringToBytes(msg.Liker); err != nil {
		return nil, errorsmod.Wrap(err, "invalid liker address")
	}

	if err := k.Keeper.LikeComment(ctx, msg.CommentId, msg.Liker); err != nil {
		return nil, errorsmod.Wrap(err, "failed to like comment")
	}

	// Emit event
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	sdkCtx.EventManager().EmitEvent(
		sdk.NewEvent(
			"comment_liked",
			sdk.NewAttribute("comment_id", fmt.Sprintf("%d", msg.CommentId)),
			sdk.NewAttribute("liker", msg.Liker),
		),
	)

	return &types.MsgLikeCommentResponse{}, nil
}