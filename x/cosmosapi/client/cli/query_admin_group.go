package cli

import (
    "fmt"
    "github.com/cosmos/cosmos-sdk/client/context"
    "github.com/cosmos/cosmos-sdk/codec"
    "github.com/yzhanginwa/cosmos-api/x/dbchain/internal/types"
    "github.com/spf13/cobra"
)

func GetCmdShowAdminGroup(queryRoute string, cdc *codec.Codec) *cobra.Command {
    return &cobra.Command{
        Use: "admin-group",
        Short: "show admin group",
        Args: cobra.ExactArgs(2),
        RunE: func(cmd *cobra.Command, args []string) error {
            cliCtx := context.NewCLIContext().WithCodec(cdc)

            accessCode := args[0]
            appCode    := args[1]
            res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/admin_group/%s/%s", queryRoute, accessCode, appCode), nil)
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

