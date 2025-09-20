package keeper

import (
	"context"

	"blogchain/x/blog/types"

	"cosmossdk.io/store/prefix"
	"github.com/cosmos/cosmos-sdk/runtime"
	"github.com/cosmos/cosmos-sdk/types/query"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (q queryServer) Posts(ctx context.Context, req *types.QueryPostsRequest) (*types.QueryPostsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	posts := []types.Post{}
	storeAdapter := runtime.KVStoreAdapter(q.k.storeService.OpenKVStore(ctx))
	postStore := prefix.NewStore(storeAdapter, []byte(types.PostKey))

	pageRes, err := query.Paginate(postStore, req.Pagination, func(key []byte, value []byte) error {
		var post types.Post
		if err := q.k.cdc.Unmarshal(value, &post); err != nil {
			return err
		}

		posts = append(posts, post)
		return nil
	})

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryPostsResponse{Posts: posts, Pagination: pageRes}, nil
}

func (q queryServer) Post(ctx context.Context, req *types.QueryPostRequest) (*types.QueryPostResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	post, found := q.k.GetPost(ctx, req.Id)
	if !found {
		return nil, status.Error(codes.NotFound, "post not found")
	}

	return &types.QueryPostResponse{Post: post}, nil
}