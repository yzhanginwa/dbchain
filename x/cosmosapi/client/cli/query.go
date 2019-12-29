package cli

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/yzhanginwa/rcv-chain/x/rcvchain/internal/types"
	"github.com/spf13/cobra"
)

func GetQueryCmd(storeKey string, cdc *codec.Codec) *cobra.Command {
	rcvchainQueryCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "Querying commands for the rcvchain module",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}
	rcvchainQueryCmd.AddCommand(client.GetCommands(
		GetCmdPoll(storeKey, cdc),
		GetCmdPollStatus(storeKey, cdc),
		GetCmdTitles(storeKey, cdc),
		GetCmdBallot(storeKey, cdc),
		GetCmdUserPolls(storeKey, cdc),
	)...)
	return rcvchainQueryCmd
}

// GetCmdPoll queries information about a poll
func GetCmdPoll(queryRoute string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use: "poll [id]",
		Short: "query poll",
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			id := args[0]

			res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/poll/%s", queryRoute, id), nil)
			if err != nil {
				fmt.Printf("could not get poll %s \n", id)
				return nil
			}

			var out types.QueryResPoll
			cdc.MustUnmarshalJSON(res, &out)
			return cliCtx.PrintOutput(out)
		},
	}
}

// GetCmdPollStatus queries information about a poll's status
func GetCmdPollStatus(queryRoute string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use: "status [name]",
		Short: "query status",
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			name := args[0]

			res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/status/%s", queryRoute, name), nil)
			if err != nil {
				fmt.Printf("could not get status of poll %s \n", name)
				return nil
			}

			var out types.QueryResStatus
			cdc.MustUnmarshalJSON(res, &out)
			return cliCtx.PrintOutput(out)
		},
	}
}

// GetCmdTitles queries information about all titles
func GetCmdTitles(queryRoute string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use: "titles",
		Short: "query titles",
		Args: cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/titles", queryRoute), nil)
			if err != nil {
				fmt.Print("could not get titles")
				return nil
			}

			var out types.QueryResTitles
			cdc.MustUnmarshalJSON(res, &out)
			return cliCtx.PrintOutput(out)
		},
	}
}

// GetCmdBallot queries information about Ballot of a vote
func GetCmdBallot(queryRoute string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use: "ballot",
		Short: "query ballot",
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/ballot/%s/%s", queryRoute, args[0], args[1]), nil)
			if err != nil {
				fmt.Print("could not get ballot")
				return nil
			}

			var out types.QueryResTitles
			cdc.MustUnmarshalJSON(res, &out)
			return cliCtx.PrintOutput(out)
		},
	}
}

// GetCmdUserPoll queries information about polls related to a user
func GetCmdUserPolls(queryRoute string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use: "user-polls",
		Short: "query user's polls",
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/userpolls/%s", queryRoute, args[0]), nil)
			if err != nil {
				fmt.Print("could not get user-polls")
				return nil
			}

			var out types.QueryResUserPolls
			cdc.MustUnmarshalJSON(res, &out)
			return cliCtx.PrintOutput(out)
		},
	}
}
