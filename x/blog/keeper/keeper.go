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

	Schema        collections.Schema
	Params        collections.Item[types.Params]
	Posts         collections.Map[uint64, types.Post]
	ActivePosts   collections.Map[uint64, bool] // Index of non-deleted posts
	DeletedPosts  collections.Map[uint64, bool] // Index of deleted posts
	PostCount     collections.Sequence
	LikedBy       collections.Map[collections.Pair[uint64, string], bool] // (postID, userAddr) -> liked
	
	// Comment collections
	Comments        collections.Map[uint64, types.Comment]
	ActiveComments  collections.Map[uint64, bool] // Index of non-deleted comments
	PostComments    collections.Map[collections.Pair[uint64, uint64], bool] // (postID, commentID) -> exists
	ChildComments   collections.Map[collections.Pair[uint64, uint64], bool] // (parentID, commentID) -> exists
	CommentCount    collections.Sequence
	LikedComments   collections.Map[collections.Pair[uint64, string], bool] // (commentID, userAddr) -> liked
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

		Params:       collections.NewItem(sb, types.ParamsKey, "params", codec.CollValue[types.Params](cdc)),
		Posts:        collections.NewMap(sb, collections.NewPrefix([]byte(types.PostKey)), "posts", collections.Uint64Key, codec.CollValue[types.Post](cdc)),
		ActivePosts:  collections.NewMap(sb, collections.NewPrefix([]byte("active_posts/")), "active_posts", collections.Uint64Key, collections.BoolValue),
		DeletedPosts: collections.NewMap(sb, collections.NewPrefix([]byte("deleted_posts/")), "deleted_posts", collections.Uint64Key, collections.BoolValue),
		PostCount:    collections.NewSequence(sb, collections.NewPrefix([]byte(types.PostCountKey)), "post_count"),
		LikedBy:      collections.NewMap(sb, collections.NewPrefix([]byte("liked_by/")), "liked_by", collections.PairKeyCodec(collections.Uint64Key, collections.StringKey), collections.BoolValue),
		
		// Comment collections
		Comments:       collections.NewMap(sb, collections.NewPrefix([]byte("comments/")), "comments", collections.Uint64Key, codec.CollValue[types.Comment](cdc)),
		ActiveComments: collections.NewMap(sb, collections.NewPrefix([]byte("active_comments/")), "active_comments", collections.Uint64Key, collections.BoolValue),
		PostComments:   collections.NewMap(sb, collections.NewPrefix([]byte("post_comments/")), "post_comments", collections.PairKeyCodec(collections.Uint64Key, collections.Uint64Key), collections.BoolValue),
		ChildComments:  collections.NewMap(sb, collections.NewPrefix([]byte("child_comments/")), "child_comments", collections.PairKeyCodec(collections.Uint64Key, collections.Uint64Key), collections.BoolValue),
		CommentCount:   collections.NewSequence(sb, collections.NewPrefix([]byte("comment_count/")), "comment_count"),
		LikedComments:  collections.NewMap(sb, collections.NewPrefix([]byte("liked_comments/")), "liked_comments", collections.PairKeyCodec(collections.Uint64Key, collections.StringKey), collections.BoolValue),
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
