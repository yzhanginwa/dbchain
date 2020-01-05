package cosmosapi

import (
    "fmt"
    "encoding/json"
    sdk "github.com/cosmos/cosmos-sdk/types"
    "github.com/yzhanginwa/cosmos-api/x/cosmosapi/internal/types"
)

// NewHandler returns a handler for "nameservice" type messages.
func NewHandler(keeper Keeper) sdk.Handler {
    return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
        switch msg := msg.(type) {
        case MsgCreateTable:
            return handleMsgCreateTable(ctx, keeper, msg)
        case MsgCreateIndex:
            return handleMsgCreateIndex(ctx, keeper, msg)
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
        return sdk.ErrUnknownRequest("Table name existed already!").Result()
    }
    keeper.CreateTable(ctx, msg.Owner, msg.TableName, msg.Fields)
    return sdk.Result{}
}

func handleMsgCreateIndex(ctx sdk.Context, keeper Keeper, msg MsgCreateIndex) sdk.Result {
    if ! keeper.IsFieldPresent(ctx, msg.TableName, msg.Field) {
        return sdk.ErrUnknownRequest(fmt.Sprintf("Field %s of table %s does not exist yet!", msg.Field, msg.TableName)).Result()
    }
    keeper.CreateIndex(ctx, msg.Owner, msg.TableName, msg.Field)
    return sdk.Result{}
}

func handleMsgInsertRow(ctx sdk.Context, keeper Keeper, msg types.MsgInsertRow) sdk.Result {
    if !keeper.IsTablePresent(ctx, msg.TableName) {
        return sdk.ErrUnknownRequest(fmt.Sprintf("Table % does not exist!", msg.TableName)).Result()
    }
    
    var rowFields types.RowFields
    if err := json.Unmarshal(msg.Fields, &rowFields); err != nil {
        return sdk.ErrUnknownRequest("Failed to parse row fields!").Result()
    }

    keeper.Insert(ctx, msg.TableName, rowFields)
    return sdk.Result{}
}
