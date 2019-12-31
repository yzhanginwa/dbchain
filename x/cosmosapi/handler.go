package cosmosapi

import (
    "fmt"
    sdk "github.com/cosmos/cosmos-sdk/types"
    "github.com/yzhanginwa/cosmos-api/x/cosmosapi/internal/types"
)

// NewHandler returns a handler for "nameservice" type messages.
func NewHandler(keeper Keeper) sdk.Handler {
    return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
        switch msg := msg.(type) {
        case MsgCreateTable:
            return handleMsgCreateTable(ctx, keeper, msg)
        case MsgInsertRow:
            return handleMsgInsertRow(ctx, keeper, msg)
        default:
            errMsg := fmt.Sprintf("Unrecognized cosmosapi Msg type: %v", msg.Type())
            return sdk.ErrUnknownRequest(errMsg).Result()
        }
    }
}

// Handle a message to create table 
func handleMsgCreateTable(ctx sdk.Context, keeper Keeper, msg MsgCreateTable) sdk.Result {
    if keeper.IsTablePresent(ctx, msg.TableName) {
        return sdk.ErrUnknownRequest("Poll name existed already!").Result()
    }
    keeper.CreateTable(ctx, msg.Owner, msg.TableName, msg.Fields)
    return sdk.Result{}
}

// TODO
func handleMsgInsertRow(ctx sdk.Context, keeper Keeper, msg types.MsgInsertRow) sdk.Result {
    if keeper.IsTablePresent(ctx, msg.TableName) {
        return sdk.ErrUnknownRequest("Poll name existed already!").Result()
    }
    keeper.Insert(ctx, msg.TableName, msg.Fields)
    return sdk.Result{}
}
