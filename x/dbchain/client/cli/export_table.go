package cli

import (
    "fmt"
    "github.com/cosmos/cosmos-sdk/client/context"
    "github.com/cosmos/cosmos-sdk/codec"
    "github.com/spf13/cobra"
)

func GetCmdExportTableRows(queryRoute string, cdc *codec.Codec) *cobra.Command {
    return &cobra.Command{
        Use: "export-table-rows <access-code> <app-code> <table-name>",
        Short: "export database schema",
        Args: cobra.ExactArgs(3),
        RunE: func(cmd *cobra.Command, args []string) error {
            cliCtx := context.NewCLIContext().WithCodec(cdc)

            accessCode := args[0]
            appCode    := args[1]
            tableName  := args[2]
            res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/export_table_rows/%s/%s/%s", queryRoute, accessCode, appCode, tableName), nil)
            if err != nil {
                fmt.Printf("Failed to export database")
                return nil
            }

            fmt.Print(string(res))
            return nil
        },
    }
}
