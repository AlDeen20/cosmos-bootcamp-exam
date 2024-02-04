package keeper

import (
	"context"
	"fmt"
	"strconv"

	"github.com/b9lab/toll-road/x/tollroad/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

func (k msgServer) CreateRoadOperator(goCtx context.Context, msg *types.MsgCreateRoadOperator) (*types.MsgCreateRoadOperatorResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	systemInfo, found := k.GetSystemInfo(ctx)
	if !found {
		panic("system info not found")
	}
	newIndex := strconv.FormatUint(systemInfo.NextOperatorId, 10)

	// Check if the value already exists
	_, isFound := k.GetRoadOperator(
		ctx,
		newIndex,
	)
	if isFound {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "index already set")
	}

	var roadOperator = types.RoadOperator{
		Creator: msg.Creator,
		Index:   newIndex,
		Name:    msg.Name,
		Token:   msg.Token,
		Active:  msg.Active,
	}

	k.SetRoadOperator(
		ctx,
		roadOperator,
	)

	systemInfo.NextOperatorId++
	k.SetSystemInfo(ctx, systemInfo)

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			"new-road-operator-created",
			sdk.NewAttribute("creator", msg.Creator),
			sdk.NewAttribute("road-operator-index", newIndex),
			sdk.NewAttribute("name", msg.Name),
			sdk.NewAttribute("token", msg.Token),
			sdk.NewAttribute("active", fmt.Sprintf("%t", msg.Active)),
		),
	)

	return &types.MsgCreateRoadOperatorResponse{
		Index: newIndex,
	}, nil
}

func (k msgServer) UpdateRoadOperator(goCtx context.Context, msg *types.MsgUpdateRoadOperator) (*types.MsgUpdateRoadOperatorResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Check if the value exists
	valFound, isFound := k.GetRoadOperator(
		ctx,
		msg.Index,
	)
	if !isFound {
		return nil, sdkerrors.Wrap(sdkerrors.ErrKeyNotFound, "index not set")
	}

	// Checks if the the msg creator is the same as the current owner
	if msg.Creator != valFound.Creator {
		return nil, sdkerrors.Wrap(sdkerrors.ErrUnauthorized, "incorrect owner")
	}

	var roadOperator = types.RoadOperator{
		Creator: msg.Creator,
		Index:   msg.Index,
		Name:    msg.Name,
		Token:   msg.Token,
		Active:  msg.Active,
	}

	k.SetRoadOperator(ctx, roadOperator)

	return &types.MsgUpdateRoadOperatorResponse{}, nil
}

func (k msgServer) DeleteRoadOperator(goCtx context.Context, msg *types.MsgDeleteRoadOperator) (*types.MsgDeleteRoadOperatorResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Check if the value exists
	valFound, isFound := k.GetRoadOperator(
		ctx,
		msg.Index,
	)
	if !isFound {
		return nil, sdkerrors.Wrap(sdkerrors.ErrKeyNotFound, "index not set")
	}

	// Checks if the the msg creator is the same as the current owner
	if msg.Creator != valFound.Creator {
		return nil, sdkerrors.Wrap(sdkerrors.ErrUnauthorized, "incorrect owner")
	}

	k.RemoveRoadOperator(
		ctx,
		msg.Index,
	)

	return &types.MsgDeleteRoadOperatorResponse{}, nil
}
