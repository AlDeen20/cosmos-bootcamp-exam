package cli

import (
	"github.com/b9lab/toll-road/x/tollroad/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/spf13/cast"
	"github.com/spf13/cobra"
)

func CmdCreateUserVault() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create-user-vault [road-operator-index] [token] [balance]",
		Short: "Create a new userVault",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			// Get indexes
			indexRoadOperatorIndex := args[0]
			indexToken := args[1]

			// Get value arguments
			argBalance, err := cast.ToUint64E(args[2])
			if err != nil {
				return err
			}

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := types.NewMsgCreateUserVault(
				clientCtx.GetFromAddress().String(),
				indexRoadOperatorIndex,
				indexToken,
				argBalance,
			)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

func CmdUpdateUserVault() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update-user-vault [road-operator-index] [token] [balance]",
		Short: "Update a userVault",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			// Get indexes
			indexRoadOperatorIndex := args[0]
			indexToken := args[1]

			// Get value arguments
			argBalance, err := cast.ToUint64E(args[2])
			if err != nil {
				return err
			}

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := types.NewMsgUpdateUserVault(
				clientCtx.GetFromAddress().String(),
				indexRoadOperatorIndex,
				indexToken,
				argBalance,
			)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

func CmdDeleteUserVault() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete-user-vault [road-operator-index] [token]",
		Short: "Delete a userVault",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			indexRoadOperatorIndex := args[0]
			indexToken := args[1]

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := types.NewMsgDeleteUserVault(
				clientCtx.GetFromAddress().String(),
				indexRoadOperatorIndex,
				indexToken,
			)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}
