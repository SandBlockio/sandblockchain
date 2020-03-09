package cli

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/sandblockio/sandblockchain/x/surprise/internal/types"
)

// GetQueryCmd returns the cli query commands for this module
func GetQueryCmd(queryRoute string, cdc *codec.Codec) *cobra.Command {
	// Group surprise queries under a subcommand
	surpriseQueryCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      fmt.Sprintf("Querying commands for the %s module", types.ModuleName),
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	surpriseQueryCmd.AddCommand(
		flags.GetCommands(
			GetCmdListBrandedTokens(queryRoute, cdc),
			GetCmdGetBrandedToken(queryRoute, cdc),
			GetCmdGetTotalSupply(queryRoute, cdc),
		)...,
	)

	return surpriseQueryCmd
}

func GetCmdListBrandedTokens(queryRoute string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List the branded tokens",
		// Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/"+types.QueryListBrandedTokens, queryRoute), nil)
			if err != nil {
				fmt.Printf("could not get branded tokens\n%s\n", err.Error())
				return nil
			}

			var out types.QueryResFetch
			cdc.MustUnmarshalJSON(res, &out)
			return cliCtx.PrintOutput(out)
		},
	}
}

func GetCmdGetBrandedToken(queryRoute string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "get [name]",
		Short: "Get the informations about a branded token",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			name := args[0]

			res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s/%s", queryRoute, types.QueryGetBrandedToken, name), nil)
			if err != nil {
				fmt.Printf("could not resolve branded token\n%s\n", err.Error())
				return nil
			}

			var out types.BrandedToken
			cdc.MustUnmarshalJSON(res, &out)
			return cliCtx.PrintOutput(out)
		},
	}
}

func GetCmdGetTotalSupply(queryRoute string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "supply",
		Short: "Get the total supply across all branded tokens",
		// Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s", queryRoute, types.QueryGetTotalSupply), nil)
			if err != nil {
				fmt.Printf("could not get branded tokens\n%s\n", err.Error())
				return nil
			}

			var out int64
			cdc.MustUnmarshalJSON(res, &out)
			return cliCtx.PrintOutput(out)
		},
	}
}
