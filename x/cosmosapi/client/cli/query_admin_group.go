package cli

import (
    "fmt"
    "github.com/cosmos/cosmos-sdk/client/context"
    "github.com/cosmos/cosmos-sdk/codec"
    "github.com/yzhanginwa/cosmos-api/x/cosmosapi/internal/types"
    "github.com/spf13/cobra"
)

func GetCmdShowAdminGroup(queryRoute string, cdc *codec.Codec) *cobra.Command {
    return &cobra.Command{
        Use: "admin-group",
        Short: "show admin group",
        Args: cobra.ExactArgs(1),
        RunE: func(cmd *cobra.Command, args []string) error {
            cliCtx := context.NewCLIContext().WithCodec(cdc)

            appCode := args[0]
            res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/admin_group/%s", queryRoute, appCode), nil)
            if err != nil {
                fmt.Printf("could not show admin group")
                return nil
            }

            var out types.QueryAdminGroup
            cdc.MustUnmarshalJSON(res, &out)
            return cliCtx.PrintOutput(out)
        },
    }
}

