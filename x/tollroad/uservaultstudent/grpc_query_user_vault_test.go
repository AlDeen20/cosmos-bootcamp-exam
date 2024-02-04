package uservaultstudent_test

import (
	"strconv"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	keepertest "github.com/b9lab/toll-road/testutil/keeper"
	"github.com/b9lab/toll-road/testutil/nullify"
	"github.com/b9lab/toll-road/x/tollroad/keeper"
	"github.com/b9lab/toll-road/x/tollroad/types"
)

func createNUserVault(keeper *keeper.Keeper, ctx sdk.Context, n int) []types.UserVault {
	items := make([]types.UserVault, n)
	for i := range items {
		items[i].Owner = strconv.Itoa(i)
		items[i].RoadOperatorIndex = strconv.Itoa(i)
		items[i].Token = strconv.Itoa(i)

		keeper.SetUserVault(ctx, items[i])
	}
	return items
}

func TestUserVaultQuerySingle(t *testing.T) {
	keeper, ctx := keepertest.TollroadKeeper(t)
	wctx := sdk.WrapSDKContext(ctx)
	msgs := createNUserVault(keeper, ctx, 2)
	for _, tc := range []struct {
		desc     string
		request  *types.QueryGetUserVaultRequest
		response *types.QueryGetUserVaultResponse
		err      error
	}{
		{
			desc: "First",
			request: &types.QueryGetUserVaultRequest{
				Owner:             msgs[0].Owner,
				RoadOperatorIndex: msgs[0].RoadOperatorIndex,
				Token:             msgs[0].Token,
			},
			response: &types.QueryGetUserVaultResponse{UserVault: msgs[0]},
		},
		{
			desc: "Second",
			request: &types.QueryGetUserVaultRequest{
				Owner:             msgs[1].Owner,
				RoadOperatorIndex: msgs[1].RoadOperatorIndex,
				Token:             msgs[1].Token,
			},
			response: &types.QueryGetUserVaultResponse{UserVault: msgs[1]},
		},
		{
			desc: "KeyNotFound",
			request: &types.QueryGetUserVaultRequest{
				Owner:             strconv.Itoa(100000),
				RoadOperatorIndex: strconv.Itoa(100000),
				Token:             strconv.Itoa(100000),
			},
			err: status.Error(codes.NotFound, "not found"),
		},
		{
			desc: "InvalidRequest",
			err:  status.Error(codes.InvalidArgument, "invalid request"),
		},
	} {
		t.Run(tc.desc, func(t *testing.T) {
			response, err := keeper.UserVault(wctx, tc.request)
			if tc.err != nil {
				require.ErrorIs(t, err, tc.err)
			} else {
				require.NoError(t, err)
				require.Equal(t,
					nullify.Fill(tc.response),
					nullify.Fill(response),
				)
			}
		})
	}
}
