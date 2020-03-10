package cosmosapi

import (
    "fmt"
    "strings"
    "bytes"
    "encoding/json"
    sdk "github.com/cosmos/cosmos-sdk/types"
    "github.com/yzhanginwa/cosmos-api/x/cosmosapi/internal/types"
    "github.com/yzhanginwa/cosmos-api/x/cosmosapi/internal/utils"
)

// NewHandler returns a handler for "nameservice" type messages.
func NewHandler(keeper Keeper) sdk.Handler {
    return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
        switch msg := msg.(type) {
        case MsgCreateApplication:
            return handleMsgCreateApplication(ctx, keeper, msg)
        case MsgCreateTable:
            return handleMsgCreateTable(ctx, keeper, msg)
        case MsgDropTable:
            return handleMsgDropTable(ctx, keeper, msg)
        case MsgAddColumn:
            return handleMsgAddColumn(ctx, keeper, msg)
        case MsgDropColumn:
            return handleMsgDropColumn(ctx, keeper, msg)
        case MsgRenameColumn:
            return handleMsgRenameColumn(ctx, keeper, msg)
        case MsgCreateIndex:
            return handleMsgCreateIndex(ctx, keeper, msg)
        case MsgDropIndex:
            return handleMsgDropIndex(ctx, keeper, msg)
        case MsgModifyOption:
            return handleMsgModifyOption(ctx, keeper, msg)
        case MsgModifyColumnOption:
            return handleMsgModifyColumnOption(ctx, keeper, msg)
        case MsgInsertRow:
            return handleMsgInsertRow(ctx, keeper, msg)
        case MsgUpdateRow:
            return handleMsgUpdateRow(ctx, keeper, msg)
        case MsgDeleteRow:
            return handleMsgDeleteRow(ctx, keeper, msg)
        case MsgAddAdminAccount:
            return handleMsgAddAdminAccount(ctx, keeper, msg)
        default:
            errMsg := fmt.Sprintf("Unrecognized cosmosapi Msg type: %v", msg.Type())
            return sdk.ErrUnknownRequest(errMsg).Result()
        }
    }
}

// Handle a message to create application
func handleMsgCreateApplication(ctx sdk.Context, keeper Keeper, msg MsgCreateApplication) sdk.Result {
    // for now, we allow anybody to create application
    // TODO: Add a system paramter "allow-creating-application", which is controlled by genesis admin
    //       If it's false, nobody can create application

    // We use the term database for internal use. To outside we use application to make users understand easily
    keeper.CreateDatabase(ctx, msg.Owner, msg.Name, msg.Description)
    return sdk.Result{}
}

// Handle a message to create table 
func handleMsgCreateTable(ctx sdk.Context, keeper Keeper, msg MsgCreateTable) sdk.Result {
    appId, err := keeper.GetDatabaseId(ctx, msg.AppCode)
    if err != nil {
        return sdk.ErrUnknownRequest("Invalid app code").Result()
    }

    if !isAdmin(ctx, keeper, msg.AppCode, msg.Owner) {
        return sdk.ErrUnknownRequest("Not authorized").Result()
    }
 
    if keeper.IsTablePresent(ctx, appId, msg.TableName) {
        return sdk.ErrUnknownRequest("Table name existed already!").Result()
    }
    keeper.CreateTable(ctx, appId, msg.Owner, msg.TableName, msg.Fields)
    return sdk.Result{}
}

func handleMsgDropTable(ctx sdk.Context, keeper Keeper, msg MsgDropTable) sdk.Result {
    appId, err := keeper.GetDatabaseId(ctx, msg.AppCode)
    if err != nil {
        return sdk.ErrUnknownRequest("Invalid app code").Result()
    }

    if !isAdmin(ctx, keeper, msg.AppCode, msg.Owner) {
        return sdk.ErrUnknownRequest("Not authorized").Result()
    }

    if !keeper.IsTablePresent(ctx, appId, msg.TableName) {
        return sdk.ErrUnknownRequest("Table name does not exist!").Result()
    }
    keeper.DropTable(ctx, appId, msg.Owner, msg.TableName)
    return sdk.Result{}
}

func handleMsgAddColumn(ctx sdk.Context, keeper Keeper, msg MsgAddColumn) sdk.Result {
    appId, err := keeper.GetDatabaseId(ctx, msg.AppCode)
    if err != nil {
        return sdk.ErrUnknownRequest("Invalid app code").Result()
    }

    if !isAdmin(ctx, keeper, msg.AppCode, msg.Owner) {
        return sdk.ErrUnknownRequest("Not authorized").Result()
    }

    field := strings.ToLower(msg.Field)
    if keeper.IsFieldPresent(ctx, appId, msg.TableName, field) {
        return sdk.ErrUnknownRequest(fmt.Sprintf("Field %s of table %s exists already!", msg.Field, msg.TableName)).Result()
    }
    keeper.AddColumn(ctx, appId, msg.TableName, field)
    return sdk.Result{}
}

func handleMsgDropColumn(ctx sdk.Context, keeper Keeper, msg MsgDropColumn) sdk.Result {
    appId, err := keeper.GetDatabaseId(ctx, msg.AppCode)
    if err != nil {
        return sdk.ErrUnknownRequest("Invalid app code").Result()
    }

    if !isAdmin(ctx, keeper, msg.AppCode, msg.Owner) {
        return sdk.ErrUnknownRequest("Not authorized").Result()
    }
    if !keeper.IsFieldPresent(ctx, appId, msg.TableName, msg.Field) {
        return sdk.ErrUnknownRequest(fmt.Sprintf("Field %s of table %s does not exist yet!", msg.Field, msg.TableName)).Result()
    }
    keeper.DropColumn(ctx, appId, msg.TableName, msg.Field)
    return sdk.Result{}
}

func handleMsgRenameColumn(ctx sdk.Context, keeper Keeper, msg MsgRenameColumn) sdk.Result {
    appId, err := keeper.GetDatabaseId(ctx, msg.AppCode)
    if err != nil {
        return sdk.ErrUnknownRequest("Invalid app code").Result()
    }

    if !isAdmin(ctx, keeper, msg.AppCode, msg.Owner) {
        return sdk.ErrUnknownRequest("Not authorized").Result()
    }
    if !keeper.IsFieldPresent(ctx, appId, msg.TableName, msg.OldField) {
        return sdk.ErrUnknownRequest(fmt.Sprintf("Field %s of table %s does not exist yet!", msg.OldField, msg.TableName)).Result()
    }

    newField := strings.ToLower(msg.NewField)
    if keeper.IsFieldPresent(ctx, appId, msg.TableName, newField) {
        return sdk.ErrUnknownRequest(fmt.Sprintf("Field %s of table %s exists already!", msg.NewField, msg.TableName)).Result()
    }
    keeper.RenameColumn(ctx, appId, msg.TableName, msg.OldField, newField)
    return sdk.Result{}
}

func handleMsgCreateIndex(ctx sdk.Context, keeper Keeper, msg MsgCreateIndex) sdk.Result {
    appId, err := keeper.GetDatabaseId(ctx, msg.AppCode)
    if err != nil {
        return sdk.ErrUnknownRequest("Invalid app code").Result()
    }

    if !isAdmin(ctx, keeper, msg.AppCode, msg.Owner) {
        return sdk.ErrUnknownRequest("Not authorized").Result()
    }
    if ! keeper.IsFieldPresent(ctx, appId, msg.TableName, msg.Field) {
        return sdk.ErrUnknownRequest(fmt.Sprintf("Field %s of table %s does not exist yet!", msg.Field, msg.TableName)).Result()
    }
    keeper.CreateIndex(ctx, appId, msg.Owner, msg.TableName, msg.Field)
    return sdk.Result{}
}

func handleMsgDropIndex(ctx sdk.Context, keeper Keeper, msg MsgDropIndex) sdk.Result {
    appId, err := keeper.GetDatabaseId(ctx, msg.AppCode)
    if err != nil {
        return sdk.ErrUnknownRequest("Invalid app code").Result()
    }

    if !isAdmin(ctx, keeper, msg.AppCode, msg.Owner) {
        return sdk.ErrUnknownRequest("Not authorized").Result()
    }
    if ! keeper.IsFieldPresent(ctx, appId, msg.TableName, msg.Field) {
        return sdk.ErrUnknownRequest(fmt.Sprintf("Field %s of table %s does not exist yet!", msg.Field, msg.TableName)).Result()
    }

    existingIndex, err := keeper.GetIndex(ctx, appId, msg.TableName)
    if err != nil {
        return sdk.ErrUnknownRequest(fmt.Sprintf("Table %s does not have any index yet!", msg.TableName)).Result()
    }
 
    if !utils.ItemExists(existingIndex, msg.Field) {
        return sdk.ErrUnknownRequest(fmt.Sprintf("Table %s does not have index on %s yet!", msg.TableName, msg.Field)).Result()
    }

    keeper.DropIndex(ctx, appId, msg.Owner, msg.TableName, msg.Field)
    return sdk.Result{}
}

func handleMsgModifyOption(ctx sdk.Context, keeper Keeper, msg MsgModifyOption) sdk.Result {
    appId, err := keeper.GetDatabaseId(ctx, msg.AppCode)
    if err != nil {
        return sdk.ErrUnknownRequest("Invalid app code").Result()
    }

    if !isAdmin(ctx, keeper, msg.AppCode, msg.Owner) {
        return sdk.ErrUnknownRequest("Not authorized").Result()
    }
    if !keeper.IsTablePresent(ctx, appId, msg.TableName) {
        return sdk.ErrUnknownRequest("Table name does not exist!").Result()
    }

    keeper.ModifyOption(ctx, appId, msg.Owner, msg.TableName, msg.Action, msg.Option)
    return sdk.Result{}
}

func handleMsgModifyColumnOption(ctx sdk.Context, keeper Keeper, msg MsgModifyColumnOption) sdk.Result {
    appId, err := keeper.GetDatabaseId(ctx, msg.AppCode)
    if err != nil {
        return sdk.ErrUnknownRequest("Invalid app code").Result()
    }

    if !isAdmin(ctx, keeper, msg.AppCode, msg.Owner) {
        return sdk.ErrUnknownRequest("Not authorized").Result()
    }
    if !keeper.IsTablePresent(ctx, appId, msg.TableName) {
        return sdk.ErrUnknownRequest("Table name does not exist!").Result()
    }

    keeper.ModifyColumnOption(ctx, appId, msg.Owner, msg.TableName, msg.FieldName, msg.Action, msg.Option)
    return sdk.Result{}
}

func handleMsgInsertRow(ctx sdk.Context, keeper Keeper, msg types.MsgInsertRow) sdk.Result {
    appId, err := keeper.GetDatabaseId(ctx, msg.AppCode)
    if err != nil {
        return sdk.ErrUnknownRequest("Invalid app code").Result()
    }

    if !keeper.IsTablePresent(ctx, appId, msg.TableName) {
        return sdk.ErrUnknownRequest(fmt.Sprintf("Table % does not exist!", msg.TableName)).Result()
    }
    
    var rowFields types.RowFields
    if err := json.Unmarshal(msg.Fields, &rowFields); err != nil {
        return sdk.ErrUnknownRequest("Failed to parse row fields!").Result()
    }

    _, err = keeper.Insert(ctx, appId, msg.TableName, rowFields, msg.Owner)
    if err != nil {
        return sdk.ErrUnknownRequest("Failed validation of inserting row").Result()
    }
    return sdk.Result{}
}

func handleMsgUpdateRow(ctx sdk.Context, keeper Keeper, msg types.MsgUpdateRow) sdk.Result {
    appId, err := keeper.GetDatabaseId(ctx, msg.AppCode)
    if err != nil {
        return sdk.ErrUnknownRequest("Invalid app code").Result()
    }

    if !keeper.IsTablePresent(ctx, appId, msg.TableName) {
        return sdk.ErrUnknownRequest(fmt.Sprintf("Table % does not exist!", msg.TableName)).Result()
    }

    options, _ := keeper.GetOption(ctx, appId, msg.TableName)
    if ! utils.ItemExists(options, string(types.TBLOPT_UPDATABLE)) {
        return sdk.ErrUnknownRequest(fmt.Sprintf("Table % is not updatable!", msg.TableName)).Result()
    }

    var rowFields types.RowFields
    if err := json.Unmarshal(msg.Fields, &rowFields); err != nil {
        return sdk.ErrUnknownRequest("Failed to parse row fields!").Result()
    }

    keeper.Update(ctx, appId, msg.TableName, msg.Id, rowFields, msg.Owner)
    return sdk.Result{}
}

func handleMsgDeleteRow(ctx sdk.Context, keeper Keeper, msg types.MsgDeleteRow) sdk.Result {
    appId, err := keeper.GetDatabaseId(ctx, msg.AppCode)
    if err != nil {
        return sdk.ErrUnknownRequest("Invalid app code").Result()
    }

    if !keeper.IsTablePresent(ctx, appId, msg.TableName) {
        return sdk.ErrUnknownRequest(fmt.Sprintf("Table % does not exist!", msg.TableName)).Result()
    }

    options, _ := keeper.GetOption(ctx, appId, msg.TableName)
    if ! utils.ItemExists(options, string(types.TBLOPT_DELETABLE)) {
        return sdk.ErrUnknownRequest(fmt.Sprintf("Table % is not updatable!", msg.TableName)).Result()
    }

    keeper.Delete(ctx, appId, msg.TableName, msg.Id, msg.Owner)
    return sdk.Result{}
}

func handleMsgAddAdminAccount(ctx sdk.Context, keeper Keeper, msg MsgAddAdminAccount) sdk.Result {
    appId, err := keeper.GetDatabaseId(ctx, msg.AppCode)
    if err != nil {
        return sdk.ErrUnknownRequest("Invalid app code").Result()
    }

    if !isAdmin(ctx, keeper, msg.AppCode, msg.Owner) {
        return sdk.ErrUnknownRequest("Not authorized").Result()
    }
    err = keeper.AddAdminAccount(ctx, appId, msg.AdminAddress)
    if err != nil {
        return sdk.ErrUnknownRequest(fmt.Sprintf("%v", err)).Result()
    }
    return sdk.Result{}
}

////////////////////
//                //
// helper methods //
//                //
////////////////////

func isAdmin(ctx sdk.Context, keeper Keeper, appCode string, address sdk.AccAddress) bool {
    adminAddresses := keeper.GetDatabaseAdmins(ctx, appCode)
    var is_admin = false
    for _, addr := range adminAddresses {
        if bytes.Compare(address, addr) == 0 {
            is_admin = true
            break
        }
    }
    return is_admin
}

