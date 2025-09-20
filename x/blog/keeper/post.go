package keeper

import (
	"context"
	"encoding/binary"

	"blogchain/x/blog/types"

	"cosmossdk.io/store/prefix"
	storetypes "cosmossdk.io/store/types"
	"github.com/cosmos/cosmos-sdk/runtime"
)

func (k Keeper) GetPostCount(ctx context.Context) uint64 {
	store := k.storeService.OpenKVStore(ctx)
	
	bz, err := store.Get([]byte(types.PostCountKey))
	if err != nil || bz == nil {
		return 0
	}
	
	return binary.BigEndian.Uint64(bz)
}

func (k Keeper) SetPostCount(ctx context.Context, count uint64) {
	store := k.storeService.OpenKVStore(ctx)
	bz := make([]byte, 8)
	binary.BigEndian.PutUint64(bz, count)
	store.Set([]byte(types.PostCountKey), bz)
}

func (k Keeper) AppendPost(ctx context.Context, post types.Post) uint64 {
	count := k.GetPostCount(ctx)
	post.Id = count
	
	store := k.storeService.OpenKVStore(ctx)
	key := GetPostKey(post.Id)
	
	appendedValue := k.cdc.MustMarshal(&post)
	store.Set(key, appendedValue)
	
	k.SetPostCount(ctx, count+1)
	return count
}

func (k Keeper) GetPost(ctx context.Context, id uint64) (types.Post, bool) {
	store := k.storeService.OpenKVStore(ctx)
	key := GetPostKey(id)
	
	bz, err := store.Get(key)
	if err != nil || bz == nil {
		return types.Post{}, false
	}
	
	var post types.Post
	k.cdc.MustUnmarshal(bz, &post)
	return post, true
}

func (k Keeper) SetPost(ctx context.Context, post types.Post) {
	store := k.storeService.OpenKVStore(ctx)
	key := GetPostKey(post.Id)
	
	bz := k.cdc.MustMarshal(&post)
	store.Set(key, bz)
}

func (k Keeper) GetAllPosts(ctx context.Context) []types.Post {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, []byte(types.PostKey))
	
	iterator := storetypes.KVStorePrefixIterator(store, []byte{})
	defer iterator.Close()
	
	var posts []types.Post
	for ; iterator.Valid(); iterator.Next() {
		var post types.Post
		k.cdc.MustUnmarshal(iterator.Value(), &post)
		posts = append(posts, post)
	}
	
	return posts
}

func GetPostKey(id uint64) []byte {
	bz := make([]byte, 8)
	binary.BigEndian.PutUint64(bz, id)
	return append([]byte(types.PostKey), bz...)
}