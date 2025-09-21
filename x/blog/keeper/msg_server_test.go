package keeper_test

import (
	"blogchain/x/blog/keeper"
	"blogchain/x/blog/types"
)

// NewMsgServerImpl returns an implementation of the MsgServer interface
// for the provided Keeper.
func NewMsgServerImpl(k keeper.Keeper) types.MsgServer {
	return keeper.NewMsgServerImpl(k)
}