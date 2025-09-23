package keeper

import (
	"context"

	"blogchain/x/blog/types"
)

// InitGenesis initializes the module's state from a provided genesis state.
func (k Keeper) InitGenesis(ctx context.Context, genState types.GenesisState) error {
	// Set all the posts and update indexes
	for _, post := range genState.Posts {
		if err := k.Posts.Set(ctx, post.Id, post); err != nil {
			return err
		}
		
		// Update indexes based on deletion status
		if post.Deleted {
			if err := k.DeletedPosts.Set(ctx, post.Id, true); err != nil {
				return err
			}
		} else {
			if err := k.ActivePosts.Set(ctx, post.Id, true); err != nil {
				return err
			}
		}
	}
	
	// Set the post count to match the highest ID
	if genState.PostCount > 0 {
		for i := uint64(0); i < genState.PostCount; i++ {
			if _, err := k.PostCount.Next(ctx); err != nil {
				return err
			}
		}
	}
	
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
	genesis.Posts, err = k.GetAllPosts(ctx)
	if err != nil {
		return nil, err
	}
	
	// Get the post count
	genesis.PostCount, err = k.GetPostCount(ctx)
	if err != nil {
		return nil, err
	}

	return genesis, nil
}
