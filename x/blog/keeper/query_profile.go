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

func (q queryServer) Profile(ctx context.Context, req *types.QueryProfileRequest) (*types.QueryProfileResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	if _, err := q.k.addressCodec.StringToBytes(req.Address); err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid address")
	}

	profile, err := q.k.GetProfile(ctx, req.Address)
	if err != nil {
		return nil, status.Error(codes.NotFound, err.Error())
	}

	return &types.QueryProfileResponse{Profile: profile}, nil
}

func (q queryServer) ProfileByUsername(ctx context.Context, req *types.QueryProfileByUsernameRequest) (*types.QueryProfileByUsernameResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	if req.Username == "" {
		return nil, status.Error(codes.InvalidArgument, "username cannot be empty")
	}

	profile, err := q.k.GetProfileByUsername(ctx, req.Username)
	if err != nil {
		return nil, status.Error(codes.NotFound, err.Error())
	}

	return &types.QueryProfileByUsernameResponse{Profile: profile}, nil
}

func (q queryServer) Profiles(ctx context.Context, req *types.QueryProfilesRequest) (*types.QueryProfilesResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	var profiles []types.Profile

	if req.Pagination != nil {
		// Use standard pagination
		storeAdapter := runtime.KVStoreAdapter(q.k.storeService.OpenKVStore(ctx))
		profileStore := prefix.NewStore(storeAdapter, []byte("profiles/"))

		pageRes, err := query.Paginate(profileStore, req.Pagination, func(key []byte, value []byte) error {
			var profile types.Profile
			if err := q.k.cdc.Unmarshal(value, &profile); err != nil {
				return err
			}
			profiles = append(profiles, profile)
			return nil
		})
		if err != nil {
			return nil, status.Error(codes.Internal, err.Error())
		}
		return &types.QueryProfilesResponse{
			Profiles:   profiles,
			Pagination: pageRes,
		}, nil
	}

	// No pagination - get all profiles
	allProfiles, err := q.k.GetAllProfiles(ctx)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryProfilesResponse{
		Profiles: allProfiles,
	}, nil
}

func (q queryServer) Followers(ctx context.Context, req *types.QueryFollowersRequest) (*types.QueryFollowersResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	if _, err := q.k.addressCodec.StringToBytes(req.Address); err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid address")
	}

	// Verify profile exists
	if _, err := q.k.GetProfile(ctx, req.Address); err != nil {
		return nil, status.Error(codes.NotFound, "profile not found")
	}

	followers, err := q.k.GetFollowers(ctx, req.Address)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	// TODO: Implement pagination for followers
	// For now, return all followers
	return &types.QueryFollowersResponse{
		Followers: followers,
	}, nil
}

func (q queryServer) Following(ctx context.Context, req *types.QueryFollowingRequest) (*types.QueryFollowingResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	if _, err := q.k.addressCodec.StringToBytes(req.Address); err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid address")
	}

	// Verify profile exists
	if _, err := q.k.GetProfile(ctx, req.Address); err != nil {
		return nil, status.Error(codes.NotFound, "profile not found")
	}

	following, err := q.k.GetFollowing(ctx, req.Address)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	// TODO: Implement pagination for following
	// For now, return all following
	return &types.QueryFollowingResponse{
		Following: following,
	}, nil
}

func (q queryServer) IsFollowing(ctx context.Context, req *types.QueryIsFollowingRequest) (*types.QueryIsFollowingResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	if _, err := q.k.addressCodec.StringToBytes(req.Follower); err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid follower address")
	}

	if _, err := q.k.addressCodec.StringToBytes(req.Following); err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid following address")
	}

	isFollowing, err := q.k.IsFollowing(ctx, req.Follower, req.Following)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryIsFollowingResponse{IsFollowing: isFollowing}, nil
}