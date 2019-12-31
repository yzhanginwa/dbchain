package cli

import (
	"strings"
	//"strconv"
	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/auth/client/utils"
	"github.com/yzhanginwa/cosmos-api/x/cosmosapi/internal/types"
)

func GetTxCmd(storeKey string, cdc *codec.Codec) *cobra.Command {
	nameserviceTxCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "Nameservice transaction subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	nameserviceTxCmd.AddCommand(client.PostCommands(
		GetCmdCreateTable(cdc),
	)...)

	return nameserviceTxCmd
}

// GetCmdCreatePoll is the CLI command for sending a CreatePoll transaction
func GetCmdCreateTable(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "create-table [name] [fields]",
		Short: "create a new poll",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))

			name := args[0]
			fields := strings.Split(args[1], ",")
			msg := types.NewMsgCreateTable(cliCtx.GetFromAddress(), name, fields)
			err := msg.ValidateBasic()
			if err != nil {
				return err
			}

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}
}


//////////////////////
//                  //
// helper functions //
//                  //
//////////////////////

