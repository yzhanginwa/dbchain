package cli

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/yzhanginwa/cosmos-api/x/cosmosapi/internal/types"
	"github.com/spf13/cobra"
)

func GetQueryCmd(storeKey string, cdc *codec.Codec) *cobra.Command {
	cosmosapiQueryCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "Querying commands for the cosmosapi module",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}
	cosmosapiQueryCmd.AddCommand(client.GetCommands(
		GetCmdTable(storeKey, cdc),
		GetCmdPollStatus(storeKey, cdc),
		GetCmdTitles(storeKey, cdc),
		GetCmdBallot(storeKey, cdc),
		GetCmdUserPolls(storeKey, cdc),
	)...)
	return cosmosapiQueryCmd
}

// GetCmdTable queries information about a table
func GetCmdTable(queryRoute string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use: "table [name]",
		Short: "query table",
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			name := args[0]

			res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/tables/%s", queryRoute, name), nil)
			if err != nil {
				fmt.Printf("could not get table %s \n", name)
				return nil
			}

			var out types.Table
			cdc.MustUnmarshalJSON(res, &out)
			return cliCtx.PrintOutput(out)
		},
	}
}

