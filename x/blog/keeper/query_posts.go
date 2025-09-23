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

	// Use the new collections API with proper pagination
	var posts []types.Post
	var pageRes *query.PageResponse

	if req.Pagination != nil {
		// Use standard Cosmos SDK pagination
		storeAdapter := runtime.KVStoreAdapter(q.k.storeService.OpenKVStore(ctx))
		postStore := prefix.NewStore(storeAdapter, []byte(types.PostKey))

		var err error
		pageRes, err = query.Paginate(postStore, req.Pagination, func(key []byte, value []byte) error {
			var post types.Post
			if err := q.k.cdc.Unmarshal(value, &post); err != nil {
				return err
			}
			// Skip deleted posts
			if !post.Deleted {
				posts = append(posts, post)
			}
			return nil
		})
		if err != nil {
			return nil, status.Error(codes.Internal, err.Error())
		}
	} else {
		// No pagination requested, return all posts
		var err error
		posts, err = q.k.GetAllPosts(ctx)
		if err != nil {
			return nil, status.Error(codes.Internal, err.Error())
		}
	}

	return &types.QueryPostsResponse{Posts: posts, Pagination: pageRes}, nil
}

func (q queryServer) Post(ctx context.Context, req *types.QueryPostRequest) (*types.QueryPostResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	post, err := q.k.GetPost(ctx, req.Id)
	if err != nil {
		return nil, status.Error(codes.NotFound, err.Error())
	}

	return &types.QueryPostResponse{Post: post}, nil
}