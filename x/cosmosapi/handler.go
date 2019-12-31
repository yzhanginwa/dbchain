package cosmosapi

import (
    "fmt"
    sdk "github.com/cosmos/cosmos-sdk/types"
)

// NewHandler returns a handler for "nameservice" type messages.
func NewHandler(keeper Keeper) sdk.Handler {
    return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
        switch msg := msg.(type) {
        case MsgCreateTable:
            return handleMsgCreateTable(ctx, keeper, msg)
        default:
            errMsg := fmt.Sprintf("Unrecognized cosmosapi Msg type: %v", msg.Type())
            return sdk.ErrUnknownRequest(errMsg).Result()
        }
    }
}

// Handle a message to create poll
func handleMsgCreateTable(ctx sdk.Context, keeper Keeper, msg MsgCreateTable) sdk.Result {
    if keeper.IsTablePresent(ctx, msg.Name) {
        return sdk.ErrUnknownRequest("Poll name existed already!").Result()
    }
    keeper.CreateTable(ctx, msg.Owner, msg.Name, msg.Fields)
    return sdk.Result{}
}

