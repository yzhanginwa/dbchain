package bank

import (
	sdk "github.com/dbchaincloud/cosmos-sdk/types"
	sdkerrors "github.com/dbchaincloud/cosmos-sdk/types/errors"
	"github.com/yzhanginwa/dbchain/x/bank/internal/keeper"
	"github.com/yzhanginwa/dbchain/x/bank/internal/types"
)

var (
	ExistentialDeposit int64
)

// NewHandler returns a handler for "bank" type messages.
func NewHandler(k keeper.Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
		ctx = ctx.WithEventManager(sdk.NewEventManager())

		switch msg := msg.(type) {
		case types.MsgSend:
			return handleMsgSend(ctx, k, msg)

		case types.MsgMultiSend:
			return handleMsgMultiSend(ctx, k, msg)

		default:
			return nil, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "unrecognized bank message type: %T", msg)
		}
	}
}

// Handle MsgSend.
func handleMsgSend(ctx sdk.Context, k keeper.Keeper, msg types.MsgSend) (*sdk.Result, error) {
	if !k.GetSendEnabled(ctx) {
		return nil, types.ErrSendDisabled
	}

	//check coin amount
	if !hasEnoughDeposits(ctx, k, msg.FromAddress, msg.Amount) {
		return nil, sdkerrors.Wrapf(types.ErrSendDisabled, "%s does not have enough coins", msg.FromAddress)
	}

	if k.BlacklistedAddr(msg.ToAddress) {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrUnauthorized, "%s is not allowed to receive transactions", msg.ToAddress)
	}

	//check can send coins
	baseKeeper ,ok := k.(keeper.BaseKeeper)
	if ok {
		if !baseKeeper.CheckCanSendCoins(ctx, msg.FromAddress, msg.ToAddress) {
			return nil , sdkerrors.Wrap(sdkerrors.ErrUnknownRequest,"permission forbidden")
		}
	}
	err := k.SendCoins(ctx, msg.FromAddress, msg.ToAddress, msg.Amount)
	if err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
		),
	)

	return &sdk.Result{Events: ctx.EventManager().Events()}, nil
}

func hasEnoughDeposits(ctx sdk.Context, k keeper.Keeper, FromAddress sdk.AccAddress, Amount sdk.Coins) bool {
	if ExistentialDeposit <= 0 {
		return true
	}
	ownCoins := k.GetCoins(ctx, FromAddress)
	for _, sendCoin := range Amount {
		hasCoin := ownCoins.AmountOf(sendCoin.Denom)
		if hasCoin.Sub(sendCoin.Amount).Int64() < ExistentialDeposit {
			return false
		}
	}

	return true
}

// Handle MsgMultiSend.
func handleMsgMultiSend(ctx sdk.Context, k keeper.Keeper, msg types.MsgMultiSend) (*sdk.Result, error) {
	// NOTE: totalIn == totalOut should already have been checked
	if !k.GetSendEnabled(ctx) {
		return nil, types.ErrSendDisabled
	}

	baseKeeper ,ok := k.(keeper.BaseKeeper)
	if ok {
		if !baseKeeper.CheckCanMultiSend(ctx, msg.Inputs, msg.Outputs) {
			return nil , sdkerrors.Wrap(sdkerrors.ErrUnknownRequest,"permission forbidden")
		}
	}

	for _, out := range msg.Outputs {
		if k.BlacklistedAddr(out.Address) {
			return nil, sdkerrors.Wrapf(sdkerrors.ErrUnauthorized, "%s is not allowed to receive transactions", out.Address)
		}
	}

	err := k.InputOutputCoins(ctx, msg.Inputs, msg.Outputs)
	if err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
		),
	)

	return &sdk.Result{Events: ctx.EventManager().Events()}, nil
}
