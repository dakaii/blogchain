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

func (q queryServer) Comments(ctx context.Context, req *types.QueryCommentsRequest) (*types.QueryCommentsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	// Verify post exists
	if _, err := q.k.GetPost(ctx, req.PostId); err != nil {
		return nil, status.Error(codes.NotFound, "post not found")
	}

	var comments []types.Comment

	if req.Pagination != nil {
		// Use standard pagination for compatibility
		storeAdapter := runtime.KVStoreAdapter(q.k.storeService.OpenKVStore(ctx))
		commentStore := prefix.NewStore(storeAdapter, []byte("comments/"))

		pageRes, err := query.Paginate(commentStore, req.Pagination, func(key []byte, value []byte) error {
			var comment types.Comment
			if err := q.k.cdc.Unmarshal(value, &comment); err != nil {
				return err
			}
			// Filter by post and parent
			if comment.PostId == req.PostId && comment.ParentId == req.ParentId && !comment.Deleted {
				comments = append(comments, comment)
			}
			return nil
		})
		if err != nil {
			return nil, status.Error(codes.Internal, err.Error())
		}
		return &types.QueryCommentsResponse{
			Comments:   comments,
			Pagination: pageRes,
		}, nil
	}

	// No pagination - get all comments for the post/parent
	fetchedComments, err := q.k.GetPostComments(ctx, req.PostId, req.ParentId)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryCommentsResponse{
		Comments: fetchedComments,
	}, nil
}

func (q queryServer) Comment(ctx context.Context, req *types.QueryCommentRequest) (*types.QueryCommentResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	comment, err := q.k.GetComment(ctx, req.Id)
	if err != nil {
		return nil, status.Error(codes.NotFound, err.Error())
	}

	if comment.Deleted {
		return nil, status.Error(codes.NotFound, "comment not found")
	}

	return &types.QueryCommentResponse{Comment: comment}, nil
}

func (q queryServer) CommentThread(ctx context.Context, req *types.QueryCommentThreadRequest) (*types.QueryCommentThreadResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	// Default to unlimited depth if not specified
	maxDepth := req.MaxDepth
	if maxDepth == 0 {
		maxDepth = 10 // Reasonable default
	}

	thread, err := q.k.GetCommentThread(ctx, req.Id, maxDepth)
	if err != nil {
		return nil, status.Error(codes.NotFound, err.Error())
	}

	if thread.Comment.Deleted {
		return nil, status.Error(codes.NotFound, "comment not found")
	}

	return &types.QueryCommentThreadResponse{Thread: *thread}, nil
}