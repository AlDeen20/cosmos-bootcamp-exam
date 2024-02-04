package keeper

import (
	"context"
	"strconv"

	"github.com/b9lab/toll-road/x/tollroad/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

func (k msgServer) CreateUserVault(goCtx context.Context, msg *types.MsgCreateUserVault) (*types.MsgCreateUserVaultResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	owner := msg.Creator

	// Check if the value already exists
	_, isFound := k.GetUserVault(
		ctx,
		owner,
		msg.RoadOperatorIndex,
		msg.Token,
	)
	if isFound {
		return nil, sdkerrors.Wrap(types.ErrIndexSet, "index already set")
	}

	// check if the balance is not zero
	if msg.Balance == 0 {
		return nil, sdkerrors.Wrapf(types.ErrZeroTokens, "invalid balance (%d)", msg.Balance)
	}

	// Banking
	ownerAddr, _ := sdk.AccAddressFromBech32(owner)
	coins := sdk.NewCoins(sdk.NewCoin(msg.Token, sdk.NewIntFromUint64(msg.Balance)))
	err := k.bank.SendCoinsFromAccountToModule(ctx, ownerAddr, types.ModuleName, coins)
	if err != nil {
		return nil, err
	}

	var userVault = types.UserVault{
		Owner:             owner,
		RoadOperatorIndex: msg.RoadOperatorIndex,
		Token:             msg.Token,
		Balance:           msg.Balance,
	}

	k.SetUserVault(
		ctx,
		userVault,
	)

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			"new-user-vault-created",
			sdk.NewAttribute("owner", msg.Creator),
			sdk.NewAttribute("road-operator-index", msg.RoadOperatorIndex),
			sdk.NewAttribute("token", msg.Token),
			sdk.NewAttribute("balance", strconv.FormatUint(msg.Balance, 10)),
		),
	)

	return &types.MsgCreateUserVaultResponse{}, nil
}

func (k msgServer) UpdateUserVault(goCtx context.Context, msg *types.MsgUpdateUserVault) (*types.MsgUpdateUserVaultResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	owner := msg.Creator

	// Check if the value exists
	valFound, isFound := k.GetUserVault(
		ctx,
		owner,
		msg.RoadOperatorIndex,
		msg.Token,
	)
	if !isFound {
		return nil, sdkerrors.Wrap(sdkerrors.ErrKeyNotFound, "index not set")
	}

	// Checks if the the msg creator is the same as the current owner
	if owner != valFound.Owner {
		return nil, sdkerrors.Wrap(sdkerrors.ErrUnauthorized, "incorrect owner")
	}

	// check if the balance is not zero
	if msg.Balance == 0 {
		return nil, sdkerrors.Wrapf(types.ErrZeroTokens, "invalid balance (%d)", msg.Balance)
	}

	// Banking
	ownerAddr, _ := sdk.AccAddressFromBech32(owner)
	if msg.Balance >= valFound.Balance {
		amount := msg.Balance - valFound.Balance
		coins := sdk.NewCoins(sdk.NewCoin(msg.Token, sdk.NewIntFromUint64(amount)))
		err := k.bank.SendCoinsFromAccountToModule(ctx, ownerAddr, types.ModuleName, coins)
		if err != nil {
			return nil, err
		}
	}

	if msg.Balance < valFound.Balance {
		amount := valFound.Balance - msg.Balance
		coins := sdk.NewCoins(sdk.NewCoin(msg.Token, sdk.NewIntFromUint64(amount)))
		err := k.bank.SendCoinsFromModuleToAccount(ctx, types.ModuleName, ownerAddr, coins)
		if err != nil {
			panic("bank error")
		}
	}

	var userVault = types.UserVault{
		Owner:             owner,
		RoadOperatorIndex: msg.RoadOperatorIndex,
		Token:             msg.Token,
		Balance:           msg.Balance,
	}

	k.SetUserVault(ctx, userVault)

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			"new-user-vault-updated",
			sdk.NewAttribute("owner", msg.Creator),
			sdk.NewAttribute("road-operator-index", msg.RoadOperatorIndex),
			sdk.NewAttribute("token", msg.Token),
			sdk.NewAttribute("balance", strconv.FormatUint(msg.Balance, 10)),
		),
	)

	return &types.MsgUpdateUserVaultResponse{}, nil
}

func (k msgServer) DeleteUserVault(goCtx context.Context, msg *types.MsgDeleteUserVault) (*types.MsgDeleteUserVaultResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	owner := msg.Creator

	// Check if the value exists
	valFound, isFound := k.GetUserVault(
		ctx,
		owner,
		msg.RoadOperatorIndex,
		msg.Token,
	)
	if !isFound {
		return nil, sdkerrors.Wrap(sdkerrors.ErrKeyNotFound, "index not set")
	}

	// Checks if the the msg creator is the same as the current owner
	if owner != valFound.Owner {
		return nil, sdkerrors.Wrap(sdkerrors.ErrUnauthorized, "incorrect owner")
	}

	// Banking
	ownerAddr, _ := sdk.AccAddressFromBech32(owner)
	coins := sdk.NewCoins(sdk.NewCoin(msg.Token, sdk.NewIntFromUint64(valFound.Balance)))
	err := k.bank.SendCoinsFromModuleToAccount(ctx, types.ModuleName, ownerAddr, coins)
	if err != nil {
		panic("bank error")
	}

	k.RemoveUserVault(
		ctx,
		owner,
		msg.RoadOperatorIndex,
		msg.Token,
	)

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			"new-user-vault-deleted",
			sdk.NewAttribute("owner", msg.Creator),
			sdk.NewAttribute("road-operator-index", msg.RoadOperatorIndex),
			sdk.NewAttribute("token", msg.Token),
		),
	)

	return &types.MsgDeleteUserVaultResponse{}, nil
}
