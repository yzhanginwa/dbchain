package cli

import (
	"strings"
	"strconv"
	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/auth/client/utils"
	"github.com/yzhanginwa/rcv-chain/x/rcvchain/internal/types"
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
		GetCmdCreatePoll(cdc),
		GetCmdAddChoice(cdc),
		GetCmdInviteVoter(cdc),
		GetCmdBeginVoting(cdc),
		GetCmdVote(cdc),
		GetCmdEndVoting(cdc),
	)...)

	return nameserviceTxCmd
}

// GetCmdCreatePoll is the CLI command for sending a CreatePoll transaction
func GetCmdCreatePoll(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "create-poll [name] [amount]",
		Short: "create a new poll",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))

			msg := types.NewMsgCreatePoll(args[0], cliCtx.GetFromAddress())
			err := msg.ValidateBasic()
			if err != nil {
				return err
			}

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}
}

// GetCmdAddChoice is the CLI command for adding a choice to poll
func  GetCmdAddChoice(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command {
		Use:   "add-choice [id] [choice]",
		Short: "add a choice",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))

			msg := types.NewMsgAddChoice(args[0], args[1], cliCtx.GetFromAddress())
			err := msg.ValidateBasic()
			if err != nil {
				return err
			}

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}
}

// GetCmdInviteVote is the CLI command for adding a voter to poll
func  GetCmdInviteVoter(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command {
		Use:   "invite-voter [id] [voter]",
		Short: "invite a voter",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))

			voterAddr, err := sdk.AccAddressFromBech32(args[1])
			if err != nil {
				return err
			}

			msg := types.NewMsgInviteVoter(args[0], voterAddr, cliCtx.GetFromAddress())
			err = msg.ValidateBasic()
			if err != nil {
				return err
			}

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}
}

// GetCmdVote is the CLI command for voters to vote
func GetCmdVote(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command {
		Use:   "vote [id] [ballod]",
		Short: "vote on a poll",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))

			votesUint, err := convertCLIBallotToUintSlice(args[1])
			if err != nil {
				return err
			}

			msg := types.NewMsgCreateBallot(args[0], votesUint, cliCtx.GetFromAddress())
			err = msg.ValidateBasic()
			if err != nil {
				return err
			}

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}
}

// GetCmdBeginVoting is the CLI command to set status to "ready"
func GetCmdBeginVoting(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command {
		Use:   "begin-voting [id]",
		Short: "begin voting",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))

			msg := types.NewMsgBeginVoting(args[0], cliCtx.GetFromAddress())
			err := msg.ValidateBasic()
			if err != nil {
				return err
			}

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}
}

// GetCmdEndVoting is the CLI command to set status to "ready"
func GetCmdEndVoting(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command {
		Use:   "end-voting [id]",
		Short: "end voting",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))

			msg := types.NewMsgEndVoting(args[0], cliCtx.GetFromAddress())
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

func convertCLIBallotToUintSlice(votes string) ([]uint16, error) {
        votesSlice := strings.Split(votes, ",")
        var votesUint []uint16
        for _, vote := range votesSlice {
                u16, err := strconv.ParseUint(vote, 10, 16)    
                        if err != nil {
                                return nil, err
                        }
                votesUint = append(votesUint, uint16(u16))  // without casting, the u16 would be type of uint64
        }
        return votesUint, nil
}

