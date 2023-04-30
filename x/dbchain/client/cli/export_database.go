package cli

import (
    "fmt"
    "github.com/cosmos/cosmos-sdk/client/context"
    "github.com/cosmos/cosmos-sdk/codec"
    "github.com/spf13/cobra"
)

func GetCmdExportDatabase (queryRoute string, cdc *codec.Codec) *cobra.Command {
    return &cobra.Command{
        Use: "export-database",
        Short: "export database schema",
        Args: cobra.ExactArgs(1),
        RunE: func(cmd *cobra.Command, args []string) error {
            cliCtx := context.NewCLIContext().WithCodec(cdc)

            appCode    := args[0]
            res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/export_database/%s", queryRoute, appCode), nil)
            if err != nil {
                fmt.Printf("Failed to export database")
                return nil
            }

            fmt.Print(string(res))
            return nil 
        },
    }
}
