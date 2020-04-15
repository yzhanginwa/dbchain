package cli

import (
    "fmt"
    "github.com/cosmos/cosmos-sdk/client/context"
    "github.com/cosmos/cosmos-sdk/codec"
    "github.com/yzhanginwa/dbchain/x/dbchain/internal/types"
    "github.com/spf13/cobra"
)

func GetCmdShowGroup(queryRoute string, cdc *codec.Codec) *cobra.Command {
    return &cobra.Command{
        Use: "show-group",
        Short: "show group",
        Args: cobra.MinimumNArgs(2),
        RunE: func(cmd *cobra.Command, args []string) error {
            cliCtx := context.NewCLIContext().WithCodec(cdc)

            accessCode := args[0]
            appCode    := args[1]
            if len(args) > 2 {
                groupName  := args[2]
                res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/group/%s/%s/%s", queryRoute, accessCode, appCode, groupName), nil)
                if err != nil {
                    fmt.Printf("could not show members of %s %s", appCode, groupName)
                    return nil
                }
                var out types.QueryGroup
                cdc.MustUnmarshalJSON(res, &out)
                return cliCtx.PrintOutput(out)
            } else {
                res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/groups/%s/%s", queryRoute, accessCode, appCode), nil)
                if err != nil {
                    fmt.Printf("could not show groups of %s", appCode)
                    return nil
                }
                var out types.QuerySliceOfString
                cdc.MustUnmarshalJSON(res, &out)
                return cliCtx.PrintOutput(out)
            }

        },
    }
}

