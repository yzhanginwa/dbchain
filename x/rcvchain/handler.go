package rcvchain

import (
	"fmt"

	"github.com/yzhanginwa/rcv-chain/x/rcvchain/internal/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// NewHandler returns a handler for "nameservice" type messages.
func NewHandler(keeper Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		switch msg := msg.(type) {
		case MsgCreatePoll:
			return handleMsgCreatePoll(ctx, keeper, msg)
		case MsgAddChoice:
			return handleMsgAddChoice(ctx, keeper, msg)
		case MsgInviteVoter:
			return handleMsgInviteVoter(ctx, keeper, msg)
		case MsgBeginVoting:
			return handleMsgBeginVoting(ctx, keeper, msg)
		case MsgCreateBallot:
			return handleMsgCreateBallot(ctx, keeper, msg)
		case MsgEndVoting:
			return handleMsgEndVoting(ctx, keeper, msg)
		default:
			errMsg := fmt.Sprintf("Unrecognized rcvchain Msg type: %v", msg.Type())
			return sdk.ErrUnknownRequest(errMsg).Result()
		}
	}
}

// Handle a message to create poll
func handleMsgCreatePoll(ctx sdk.Context, keeper Keeper, msg MsgCreatePoll) sdk.Result {
        if keeper.IsPollPresent(ctx, types.NewPollId(msg.Title, msg.Owner)) {
		return sdk.ErrUnknownRequest("Poll name existed already!").Result()
	}
        keeper.CreatePoll(ctx, msg.Title, msg.Owner)
	keeper.SaveUserPolls(ctx, types.NewPollId(msg.Title, msg.Owner), msg.Owner, true)
	return sdk.Result{}
}

func handleMsgAddChoice(ctx sdk.Context, keeper Keeper, msg MsgAddChoice) sdk.Result {
        if !keeper.IsPollPresent(ctx, msg.Id) {
		return sdk.ErrUnknownRequest("Poll doesn't exist!").Result()
	}
	if keeper.GetPollStatus(ctx, msg.Id) != types.PollStatusNew {
		return sdk.ErrUnknownRequest("too late to add choice").Result()
	}

	keeper.AddChoice(ctx, msg.Id, msg.Choice, msg.Owner)
	return sdk.Result{}
}

func handleMsgInviteVoter(ctx sdk.Context, keeper Keeper, msg MsgInviteVoter) sdk.Result {
        if !keeper.IsPollPresent(ctx, msg.Id) {
		return sdk.ErrUnknownRequest("Poll doesn't exist!").Result()
	}
	if keeper.GetPollStatus(ctx, msg.Id) != types.PollStatusNew {
		return sdk.ErrUnknownRequest("too late to invite voters").Result()
	}

	keeper.InviteVoter(ctx, msg.Id, msg.Voter, msg.Owner)
	keeper.SaveUserPolls(ctx, msg.Id, msg.Voter, false)
	return sdk.Result{}
}

func handleMsgBeginVoting(ctx sdk.Context, keeper Keeper, msg MsgBeginVoting) sdk.Result {
        if !keeper.IsPollPresent(ctx, msg.Id) {
		return sdk.ErrUnknownRequest("Poll doesn't exist!").Result()
	}
	if keeper.GetPollStatus(ctx, msg.Id) != types.PollStatusNew {
		return sdk.ErrUnknownRequest("status new is expected").Result()
	}

	keeper.BeginVoting(ctx, msg.Id, msg.Owner)
	return sdk.Result{}
}

func handleMsgCreateBallot(ctx sdk.Context, keeper Keeper, msg MsgCreateBallot) sdk.Result {
        if !keeper.IsPollPresent(ctx, msg.Id) {
		return sdk.ErrUnknownRequest("Poll doesn't exist!").Result()
	}
	poll, _ := keeper.GetPoll(ctx, msg.Id)
	if poll.Status != types.PollStatusReady {
		return sdk.ErrUnknownRequest("Poll is not ready for voting!").Result()
	}

	for _, addr := range poll.Voters {
		if addr.Equals(msg.Voter) {
			keeper.CreateBallot(ctx, msg.Id, msg.Votes, msg.Voter)
			return sdk.Result{}
		}
	}
	return sdk.ErrUnknownRequest("Unauthorized voter!").Result()
}

func handleMsgEndVoting(ctx sdk.Context, keeper Keeper, msg MsgEndVoting) sdk.Result {
        if !keeper.IsPollPresent(ctx, msg.Id) {
		return sdk.ErrUnknownRequest("Poll doesn't exist!").Result()
	}

	if keeper.GetPollStatus(ctx, msg.Id) != types.PollStatusReady {
		return sdk.ErrUnknownRequest("Poll is not ready for voting!").Result()
	}

	keeper.EndVoting(ctx, msg.Id, msg.Owner)
	return sdk.Result{}
}

