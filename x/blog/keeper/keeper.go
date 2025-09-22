package keeper

import (
	"fmt"

	"cosmossdk.io/collections"
	"cosmossdk.io/core/address"
	corestore "cosmossdk.io/core/store"
	"github.com/cosmos/cosmos-sdk/codec"

	"blogchain/x/blog/types"
)

type Keeper struct {
	storeService corestore.KVStoreService
	cdc          codec.Codec
	addressCodec address.Codec
	// Address capable of executing a MsgUpdateParams message.
	// Typically, this should be the x/gov module account.
	authority []byte

	Schema    collections.Schema
	Params    collections.Item[types.Params]
	Posts     collections.Map[uint64, types.Post]
	PostCount collections.Sequence
	LikedBy   collections.Map[collections.Pair[uint64, string], bool] // (postID, userAddr) -> liked
}

func NewKeeper(
	storeService corestore.KVStoreService,
	cdc codec.Codec,
	addressCodec address.Codec,
	authority []byte,
) (Keeper, error) {
	if _, err := addressCodec.BytesToString(authority); err != nil {
		return Keeper{}, fmt.Errorf("invalid authority address %s: %w", authority, err)
	}

	sb := collections.NewSchemaBuilder(storeService)

	k := Keeper{
		storeService: storeService,
		cdc:          cdc,
		addressCodec: addressCodec,
		authority:    authority,

		Params:    collections.NewItem(sb, types.ParamsKey, "params", codec.CollValue[types.Params](cdc)),
		Posts:     collections.NewMap(sb, collections.NewPrefix([]byte(types.PostKey)), "posts", collections.Uint64Key, codec.CollValue[types.Post](cdc)),
		PostCount: collections.NewSequence(sb, collections.NewPrefix([]byte(types.PostCountKey)), "post_count"),
		LikedBy:   collections.NewMap(sb, collections.NewPrefix([]byte("liked_by/")), "liked_by", collections.PairKeyCodec(collections.Uint64Key, collections.StringKey), collections.BoolValue),
	}

	schema, err := sb.Build()
	if err != nil {
		return Keeper{}, fmt.Errorf("failed to build schema: %w", err)
	}
	k.Schema = schema

	return k, nil
}

// GetAuthority returns the module's authority.
func (k Keeper) GetAuthority() []byte {
	return k.authority
}
