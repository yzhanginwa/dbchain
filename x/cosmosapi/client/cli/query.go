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
        GetCmdTables(storeKey, cdc),
        GetCmdFindRow(storeKey, cdc),
    )...)
    return cosmosapiQueryCmd
}

// GetCmdTables lists all table names
func GetCmdTables(queryRoute string, cdc *codec.Codec) *cobra.Command {
    return &cobra.Command{
        Use: "tables",
        Short: "query tables",
        Args: cobra.MaximumNArgs(1),
        RunE: func(cmd *cobra.Command, args []string) error {
            cliCtx := context.NewCLIContext().WithCodec(cdc)
            var path string
            if len(args) == 1 {
                path = fmt.Sprintf("custom/%s/tables/%s", queryRoute, args[0])
            } else {
                path = fmt.Sprintf("custom/%s/tables", queryRoute)
            }

            res, _, err := cliCtx.QueryWithData(path, nil)
            if err != nil {
                fmt.Printf("could not get table names")
                return nil
            }

            if len(args) == 1 {
                var out types.Table
                cdc.MustUnmarshalJSON(res, &out)
                return cliCtx.PrintOutput(out)
            } else {
                var out types.QueryTables
                cdc.MustUnmarshalJSON(res, &out)
                return cliCtx.PrintOutput(out)
            }
        },
    }
}

func GetCmdFindRow(queryRoute string, cdc *codec.Codec) *cobra.Command {
    return &cobra.Command{
        Use: "find",
        Short: "find row",
        Args: cobra.ExactArgs(2),
        RunE: func(cmd *cobra.Command, args []string) error {
            cliCtx := context.NewCLIContext().WithCodec(cdc)

            res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/find/%s/%s", queryRoute, args[0], args[1]), nil)
            if err != nil {
                fmt.Printf("could not find row")
                return nil
            }

            var out types.QueryRowFields
            cdc.MustUnmarshalJSON(res, &out)
            return cliCtx.PrintOutput(out)
        },
    }
}
