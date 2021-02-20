package cli

import (
    "bufio"
    "github.com/cosmos/cosmos-sdk/client/context"
    "github.com/cosmos/cosmos-sdk/codec"
    sdk "github.com/cosmos/cosmos-sdk/types"
    "github.com/cosmos/cosmos-sdk/x/auth"
    "github.com/cosmos/cosmos-sdk/x/auth/client/utils"
    "github.com/spf13/cobra"
    "github.com/yzhanginwa/dbchain/x/dbchain/internal/types"
)

func GetCmdAddFunction(cdc *codec.Codec) *cobra.Command {
    return &cobra.Command{
        Use:   "add-function [appCode] [name] [parameters] [code]",
        Short: "add a function",
        Args:  cobra.ExactArgs(4),
        RunE: func(cmd *cobra.Command, args []string) error {
            cliCtx := context.NewCLIContext().WithCodec(cdc)
            inBuf := bufio.NewReader(cmd.InOrStdin())
            txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(cdc))

            appCode   := args[0]
            funcName  := args[1]
            parameter := args[2]
            body      := args[3]

            msg := types.NewMsgAddFunction(cliCtx.GetFromAddress(), appCode, funcName, parameter, body)
            err := msg.ValidateBasic()
            if err != nil {
                return err
            }

            return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
        },
    }
}

func GetCmdCallFunction(cdc *codec.Codec) *cobra.Command {
    return &cobra.Command{
        Use:   "call-function [appCode] [name] [parameters]",
        Short: "call a function",
        Args:  cobra.ExactArgs(3),
        RunE: func(cmd *cobra.Command, args []string) error {
            cliCtx := context.NewCLIContext().WithCodec(cdc)
            inBuf := bufio.NewReader(cmd.InOrStdin())
            txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(cdc))

            appCode  := args[0]
            funcName := args[1]
            argument := args[2]
            msg := types.NewMsgCallFunction(cliCtx.GetFromAddress(), appCode, funcName, argument)
            err := msg.ValidateBasic()
            if err != nil {
                return err
            }

            return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
        },
    }
}

func GetCmdAddCustomQuerier(cdc *codec.Codec) *cobra.Command {
    return &cobra.Command{
        Use:   "add-custom-querier [appCode] [name] [parameters] [code]",
        Short: "add a function",
        Args:  cobra.ExactArgs(4),
        RunE: func(cmd *cobra.Command, args []string) error {
            cliCtx := context.NewCLIContext().WithCodec(cdc)
            inBuf := bufio.NewReader(cmd.InOrStdin())
            txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(cdc))

            appCode      := args[0]
            querierName  := args[1]
            parameter    := args[2]
            body         := args[3]


            msg := types.NewMsgAddCustomQuerier(cliCtx.GetFromAddress(), appCode, querierName, parameter, body)
            err := msg.ValidateBasic()
            if err != nil {
                return err
            }

            return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
        },
    }
}
