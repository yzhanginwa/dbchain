package cli

import (
    "fmt"
    "github.com/dbchaincloud/cosmos-sdk/client/context"
    "github.com/dbchaincloud/cosmos-sdk/codec"
    "github.com/yzhanginwa/dbchain-sm/x/dbchain/internal/types"
    "github.com/spf13/cobra"
)

func GetCmdShowFriends(queryRoute string, cdc *codec.Codec) *cobra.Command {
    return &cobra.Command{
        Use: "show-friends [accessCode]",
        Short: "show friends",
        Args: cobra.ExactArgs(1),
        RunE: func(cmd *cobra.Command, args []string) error {
            cliCtx := context.NewCLIContext().WithCodec(cdc)

            accessCode := args[0]
            res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/friends/%s", queryRoute, accessCode), nil)
            if err != nil {
                fmt.Printf("could not show friend")
                return nil
            }

            var out types.QueryOfFriends
            cdc.MustUnmarshalJSON(res, &out)
            return cliCtx.PrintOutput(out)
        },
    }
}

func GetCmdShowPendingFriends(queryRoute string, cdc *codec.Codec) *cobra.Command {
    return &cobra.Command{
        Use: "show-pending-friends [accessCode]",
        Short: "show pending friends",
        Args: cobra.ExactArgs(1),
        RunE: func(cmd *cobra.Command, args []string) error {
            cliCtx := context.NewCLIContext().WithCodec(cdc)

            accessCode := args[0]
            res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/pending_friends/%s", queryRoute, accessCode), nil)
            if err != nil {
                fmt.Printf("could not show friend")
                return nil
            }

            var out types.QueryOfFriends
            cdc.MustUnmarshalJSON(res, &out)
            return cliCtx.PrintOutput(out)
        },
    }
}

