package cli

import (
	"bufio"
	"fmt"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/auth/client/utils"
	"github.com/spf13/cobra"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"strconv"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/sandblockio/sandblockchain/x/surprise/internal/types"
)

// GetTxCmd returns the transaction commands for this module
func GetTxCmd(cdc *codec.Codec) *cobra.Command {
	surpriseTxCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      fmt.Sprintf("%s transactions subcommands", types.ModuleName),
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	surpriseTxCmd.AddCommand(flags.PostCommands(
		GetCmdCreateBrandedToken(cdc),
		GetCmdTransferBrandedTokenOwnership(cdc),
		GetCmdMintBrandedToken(cdc),
		GetCmdBurnBrandedToken(cdc),
	)...)

	return surpriseTxCmd
}

func GetCmdCreateBrandedToken(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use: "create-token [name] [initial-supply]",
		Short: "Create a new branded token with an initial supply",
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			// Acquire instances
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			inBuf := bufio.NewReader(cmd.InOrStdin())
			txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(cdc))

			// Extract params
			coins, err := strconv.ParseInt(args[1], 10, 64)
			if err != nil {
				return err
			}

			// Construct and validate the payload
			msg := types.NewMsgCreateBrandedToken(args[0], sdk.NewInt(coins), cliCtx.GetFromAddress())
			err = msg.ValidateBasic()
			if err != nil {
				return err
			}

			// Dispatch and return
			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}
}

func GetCmdTransferBrandedTokenOwnership(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use: "transfer-token-ownership [name] [owner]",
		Short: "Transfer the ownership over a token to a new address",
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			// Acquire instances
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			inBuf := bufio.NewReader(cmd.InOrStdin())
			txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(cdc))

			// Extract params
			destinationAddress, err := sdk.AccAddressFromBech32(args[1])
			if err != nil {
				return err
			}

			// Construct and validate the payload
			msg := types.NewMsgTransferBrandedTokenOwnership(args[0], cliCtx.GetFromAddress(), destinationAddress)
			err = msg.ValidateBasic()
			if err != nil {
				return err
			}

			// Dispatch and return
			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}
}

func GetCmdMintBrandedToken(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use: "mint-token [name] [amount]",
		Short: "Mint new units of that Branded Token",
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			// Acquire instances
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			inBuf := bufio.NewReader(cmd.InOrStdin())
			txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(cdc))

			// Extract params
			coins, err := strconv.ParseInt(args[1], 10, 64)
			if err != nil {
				return err
			}

			// Construct and validate the payload
			msg := types.NewMsgMintBrandedToken(cliCtx.GetFromAddress(), args[0], sdk.NewInt(coins))
			err = msg.ValidateBasic()
			if err != nil {
				return err
			}

			// Dispatch and return
			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}
}

func GetCmdBurnBrandedToken(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use: "burn-token [name] [amount]",
		Short: "Burn units of a given Branded Token",
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			// Acquire instances
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			inBuf := bufio.NewReader(cmd.InOrStdin())
			txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(cdc))

			// Extract params
			coins, err := strconv.ParseInt(args[1], 10, 64)
			if err != nil {
				return err
			}

			// Construct and validate the payload
			msg := types.NewMsgBurnBrandedToken(cliCtx.GetFromAddress(), args[0], sdk.NewInt(coins))
			err = msg.ValidateBasic()
			if err != nil {
				return err
			}

			// Dispatch and return
			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}
}