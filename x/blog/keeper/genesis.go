package keeper

import (
	"context"

	"blogchain/x/blog/types"
)

// InitGenesis initializes the module's state from a provided genesis state.
func (k Keeper) InitGenesis(ctx context.Context, genState types.GenesisState) error {
	// Set all the posts
	for _, post := range genState.Posts {
		k.SetPost(ctx, post)
	}
	
	// Set the post count
	k.SetPostCount(ctx, genState.PostCount)
	
	// Set the params
	return k.Params.Set(ctx, genState.Params)
}

// ExportGenesis returns the module's exported genesis.
func (k Keeper) ExportGenesis(ctx context.Context) (*types.GenesisState, error) {
	var err error

	genesis := types.DefaultGenesis()
	genesis.Params, err = k.Params.Get(ctx)
	if err != nil {
		return nil, err
	}
	
	// Get all posts
	genesis.Posts = k.GetAllPosts(ctx)
	
	// Get the post count
	genesis.PostCount = k.GetPostCount(ctx)

	return genesis, nil
}
