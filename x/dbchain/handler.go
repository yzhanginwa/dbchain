package dbchain

import (
    "encoding/base64"
    "fmt"
    "strings"
    "bytes"
    "encoding/json"
    sdk "github.com/dbchaincloud/cosmos-sdk/types"
    sdkerrors "github.com/dbchaincloud/cosmos-sdk/types/errors"
    "github.com/yzhanginwa/dbchain/x/dbchain/internal/types"
    "github.com/yzhanginwa/dbchain/x/dbchain/internal/utils"
    "github.com/dbchaincloud/cosmos-sdk/version"
)

const (
    CommunityEdition = "dbChainCommunity"
)

var (
    AllowCreateApplication bool
)

// NewHandler returns a handler for "hain" type messages.
func NewHandler(keeper Keeper) sdk.Handler {
    return func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
        var result *sdk.Result
        var err error
        defer setTxStatus(ctx, keeper, &err)

        switch msg := msg.(type) {
        case MsgCreateApplication:
            result, err = handleMsgCreateApplication(ctx, keeper, msg)
        case MsgDropApplication:
            result, err = handleMsgDropApplication(ctx, keeper, msg)
        case MsgRecoverApplication:
            result, err = handleMsgRecoverApplication(ctx, keeper, msg)
        case MsgCreateSysDatabase:
            result, err = handleMsgCreateSysDatabase(ctx, keeper, msg)
        case MsgSetAppUserFileVolumeLimit:
            result, err = handleMsgSetAppUserFileVolumeLimit(ctx, keeper, msg)
        case MsgModifyDatabaseUser:
            result, err = handleMsgModifyDatabaseUser(ctx, keeper, msg)
        case MsgAddFunction:
            result, err = handleMsgAddFunction(ctx, keeper, msg)
        case MsgCallFunction:
            result, err = handleMsgCallFunction(ctx, keeper, msg)
        case MsgDropFunction:
            result, err = handleMsgDropFunction(ctx, keeper, msg)
        case MsgAddCustomQuerier:
            result, err = handleMsgAddCustomQuerier(ctx, keeper, msg)
        case MsgDropCustomQuerier:
            result, err = handleMsgDropCustomQuerier(ctx, keeper, msg)
        case MsgCreateTable:
            result, err = handleMsgCreateTable(ctx, keeper, msg)
        case MsgModifyTableAssociation:
            result, err = handleMsgModifyTableAssociation(ctx, keeper, msg)
        case MsgAddCounterCache:
            result, err = handleMsgAddCountCache(ctx, keeper, msg)
        case MsgDropTable:
            result, err = handleMsgDropTable(ctx, keeper, msg)
        case MsgAddColumn:
            result, err = handleMsgAddColumn(ctx, keeper, msg)
        case MsgDropColumn:
            result, err = handleMsgDropColumn(ctx, keeper, msg)
        case MsgRenameColumn:
            result, err = handleMsgRenameColumn(ctx, keeper, msg)
        case MsgCreateIndex:
            result, err = handleMsgCreateIndex(ctx, keeper, msg)
        case MsgDropIndex:
            result, err = handleMsgDropIndex(ctx, keeper, msg)
        case MsgModifyOption:
            result, err = handleMsgModifyOption(ctx, keeper, msg)
        case MsgAddInsertFilter:
            result, err = handleMsgAddInsertFilter(ctx, keeper, msg)
        case MsgDropInsertFilter:
            result, err = handleMsgDropInsertFilter(ctx, keeper, msg)
        case MsgAddTrigger:
            result, err = handleMsgAddTrigger(ctx, keeper, msg)
        case MsgDropTrigger:
            result, err = handleMsgDropTrigger(ctx, keeper, msg)
        case MsgSetTableMemo:
            result, err = handleMsgSetTableMemo(ctx, keeper, msg)
        case MsgModifyColumnOption:
            result, err = handleMsgModifyColumnOption(ctx, keeper, msg)
        case MsgSetColumnDataType:
            result, err = handleMsgSetColumnDataType(ctx, keeper, msg)
        case MsgSetColumnMemo:
            result, err = handleMsgSetColumnMemo(ctx, keeper, msg)
        case MsgInsertRow:
            result, err = handleMsgInsertRow(ctx, keeper, msg)
        case MsgUpdateRow:
            result, err = handleMsgUpdateRow(ctx, keeper, msg)
        case MsgDeleteRow:
            result, err = handleMsgDeleteRow(ctx, keeper, msg)
        case MsgFreezeRow:
            result, err = handleMsgFreezeRow(ctx, keeper, msg)
        case MsgModifyGroup:
            result, err = handleMsgModifyGroup(ctx, keeper, msg)
        case MsgSetGroupMemo:
            result, err = handleMsgSetGroupMemo(ctx, keeper, msg)
        case MsgModifyGroupMember:
            result, err = handleMsgModifyGroupMember(ctx, keeper, msg)
        case MsgAddFriend:
            result, err = handleMsgAddFriend(ctx, keeper, msg)
        case MsgDropFriend:
            result, err = handleMsgDropFriend(ctx, keeper, msg)
        case MsgRespondFriend:
            result, err = handleMsgRespondFriend(ctx, keeper, msg)
        case MsgSetSchemaStatus:
            result, err = handleMsgSetSchemaStatus(ctx, keeper, msg)
        case MsgSetDatabasePermission:
            result, err = handleMsgSetDatabasePermission(ctx, keeper, msg)
        case MsgUpdateTotalTx:
            result, err = handleMsgUpdateTotalTx(ctx, keeper, msg)
        case MsgUpdateTxStatistic:
            result, err = handleMsgUpdateTxStatistic(ctx, keeper, msg)
        case MsgModifyP2PTransferLimit:
            result, err = handleMsgModifyP2PTransferLimit(ctx, keeper, msg)
        case MsgModifyTokenKeeperMember:
            result, err = handleMsgModifyTokenKeeperMember(ctx, keeper, msg)
        case types.MsgSaveUserPrivateKey:
            result, err = handleMsgSaveUserPrivateKey(ctx, keeper, msg)
        default:
            errMsg := fmt.Sprintf("Unrecognized dbchain Msg type: %v", msg.Type())
            result, err = nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, errMsg)
        }
        return result, err
    }
}

// Handle a message to create application
func handleMsgCreateApplication(ctx sdk.Context, keeper Keeper, msg MsgCreateApplication) (*sdk.Result, error) {
    if !AllowCreateApplication {
        if !isSysAdmin(ctx, keeper, msg.Owner) {
            return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest,"Not authorized")
        }
    }

    if version.Name == CommunityEdition {
        var apps = keeper.GetAllAppCode(ctx)
        if len(apps) > 2 {
            return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "No more than 2 apps allowed")
        }
    }
    bz, err := base64.StdEncoding.DecodeString(msg.Description)
    if err != nil {
        return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest,fmt.Sprintf("%v", err))
    }
    msg.Description = string(bz)
    // We use the term database for internal use. To outside we use application to make users understand easily
    err = keeper.CreateDatabase(ctx, msg.Owner, msg.Name, msg.Description, msg.PermissionRequired, false)
    if err != nil {
        return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest,fmt.Sprintf("%v", err))
    }
    return &sdk.Result{}, nil
}

func handleMsgDropApplication(ctx sdk.Context, keeper Keeper, msg MsgDropApplication) (*sdk.Result, error) {
    appId, err := keeper.GetDatabaseId(ctx, msg.AppCode)
    if err != nil {
       return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "Invalid app code")
    }

    if !isAdmin(ctx, keeper, appId, msg.Owner) {
        return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "Not authorized")
    }

    err = keeper.DeleteApplication(ctx, msg.AppCode)
    if err != nil {
        return nil, err
    }
    return &sdk.Result{}, nil
}

func handleMsgRecoverApplication(ctx sdk.Context, keeper Keeper, msg MsgRecoverApplication) (*sdk.Result, error) {
    appId, err := keeper.GetDatabaseIdWithoutCheck(ctx, msg.AppCode)
    if err != nil {
        return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "Invalid app code")
    }

    if !isAdmin(ctx, keeper, appId, msg.Owner) {
        return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "Not authorized")
    }

    keeper.RecoverApplication(ctx, msg.AppCode)
    return &sdk.Result{}, nil
}

func handleMsgCreateSysDatabase(ctx sdk.Context, keeper Keeper, msg MsgCreateSysDatabase) (*sdk.Result, error) {
    // only sys admin can create the sys database
    if !isSysAdmin(ctx, keeper, msg.Owner) {
        return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest,"Not authorized")
    }

    err := keeper.CreateDatabase(ctx, msg.Owner, "sysdb", "database for the use of system only", false, true)
    if err != nil {
        return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest,fmt.Sprintf("%v", err))
    }
    return &sdk.Result{}, nil
}

func handleMsgSetAppUserFileVolumeLimit(ctx sdk.Context, keeper Keeper, msg MsgSetAppUserFileVolumeLimit) (*sdk.Result, error) {
    appId, err := keeper.GetDatabaseIdNotFrozen(ctx, msg.AppCode)
    if err != nil {
        return nil, err
    }
    // only sys admin can set file volume limit
    if !isAdmin(ctx, keeper, appId, msg.Owner) {
        return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest,"Not authorized")
    }

    err = keeper.SetAppUserFileVolumeLimit(ctx, appId, msg.Size)
    if err != nil {
        return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest,fmt.Sprintf("%v", err))
    }
    return &sdk.Result{}, nil
}

func handleMsgModifyDatabaseUser(ctx sdk.Context, keeper Keeper, msg MsgModifyDatabaseUser) (*sdk.Result, error) {
    // We use the term database for internal use. To outside we use application to make users understand easily
    err := keeper.ModifyDatabaseUser(ctx, msg.Owner, msg.AppCode, msg.Action, msg.User)
    if err != nil {
        return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest,fmt.Sprintf("%v", err))
    }
    return &sdk.Result{}, nil
}

func handleMsgAddFunction(ctx sdk.Context, keeper Keeper, msg MsgAddFunction) (*sdk.Result, error) {
    appId, err := keeper.GetDatabaseIdNotFrozen(ctx, msg.AppCode)
    if err != nil {
        return nil, err
    }

    if !isAdmin(ctx, keeper, appId, msg.Owner) {
        return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "Not authorized")
    }

    if version.Name == CommunityEdition {
        tables := keeper.GetTables(ctx, appId)
        if len(tables) > 29 {
            return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "No more than 30 tables allowed")
        }
    }
    bz , err := base64.StdEncoding.DecodeString(msg.Description)
    if err != nil {
        return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, err.Error())
    }
    msg.Description = string(bz)
    bz , err = base64.StdEncoding.DecodeString(msg.Body)
    if err != nil {
        return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, err.Error())
    }
    msg.Body = string(bz)

    //TODO Does it need to be checked that if the function has been added
    err = keeper.AddFunction(ctx, appId, msg.FunctionName, msg.Description, msg.Body, msg.Owner, 0)
    if err != nil{
        return nil, err
    }
    return &sdk.Result{}, nil
}

func handleMsgCallFunction(ctx sdk.Context, keeper Keeper, msg MsgCallFunction) (*sdk.Result, error) {
    appId, err := keeper.GetDatabaseId(ctx, msg.AppCode)
    if err != nil {
        return nil, err
    }

    err = keeper.CallFunction(ctx, appId, msg.Owner, msg.FunctionName, msg.Argument)
    if err != nil{
        return nil, err
    }
    return &sdk.Result{}, nil
}

func handleMsgDropFunction(ctx sdk.Context, keeper Keeper, msg MsgDropFunction) (*sdk.Result, error) {
    appId, err := keeper.GetDatabaseIdNotFrozen(ctx, msg.AppCode)
    if err != nil {
        return nil, err
    }

    if !isAdmin(ctx, keeper, appId, msg.Owner) {
        return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "Not authorized")
    }

    err = keeper.DropFunction(ctx, appId, msg.Owner, msg.FunctionName, 0)
    if err != nil{
        return nil, err
    }
    return &sdk.Result{}, nil
}

func handleMsgAddCustomQuerier(ctx sdk.Context, keeper Keeper, msg MsgAddCustomQuerier) (*sdk.Result, error) {
    appId, err := keeper.GetDatabaseIdNotFrozen(ctx, msg.AppCode)
    if err != nil {
        return nil, err
    }

    if !isAdmin(ctx, keeper, appId, msg.Owner) {
        return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "Not authorized")
    }

    if version.Name == CommunityEdition {
        tables := keeper.GetTables(ctx, appId)
        if len(tables) > 29 {
            return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "No more than 30 tables allowed")
        }
    }
    bz , err := base64.StdEncoding.DecodeString(msg.Description)
    if err != nil {
        return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, err.Error())
    }
    msg.Description = string(bz)
    bz , err = base64.StdEncoding.DecodeString(msg.Body)
    if err != nil {
        return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, err.Error())
    }
    msg.Body = string(bz)
    //TODO Does it need to be checked that if the function has been added
    err = keeper.AddFunction(ctx, appId, msg.QuerierName, msg.Description, msg.Body, msg.Owner, 1)
    if err != nil{
        return nil, err
    }
    return &sdk.Result{}, nil
}


func handleMsgDropCustomQuerier(ctx sdk.Context, keeper Keeper, msg MsgDropCustomQuerier) (*sdk.Result, error) {
    appId, err := keeper.GetDatabaseIdNotFrozen(ctx, msg.AppCode)
    if err != nil {
        return nil, err
    }

    if !isAdmin(ctx, keeper, appId, msg.Owner) {
        return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "Not authorized")
    }

    err = keeper.DropFunction(ctx, appId, msg.Owner, msg.QuerierName, 1)
    if err != nil{
        return nil, err
    }
    return &sdk.Result{}, nil
}

// Handle a message to create table 
func handleMsgCreateTable(ctx sdk.Context, keeper Keeper, msg MsgCreateTable) (*sdk.Result, error) {
    appId, err := keeper.GetDatabaseIdNotFrozen(ctx, msg.AppCode)
    if err != nil {
        return nil, err
    }

    if !isAdmin(ctx, keeper, appId, msg.Owner) {
        return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "Not authorized")
    }
 
    if version.Name == CommunityEdition {
        tables := keeper.GetTables(ctx, appId)
        if len(tables) > 29 {
            return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "No more than 30 tables allowed")
        }
    }

    if keeper.HasTable(ctx, appId, msg.TableName) {
        return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "Table name existed already!")
    }
    keeper.CreateTable(ctx, appId, msg.Owner, msg.TableName, msg.Fields)
    err = keeper.CreateIndex(ctx, appId, msg.Owner, msg.TableName, "created_by")
    if err != nil {
        return nil, err
    }
    return &sdk.Result{}, nil
}

// Handle a message to create table
func handleMsgModifyTableAssociation(ctx sdk.Context, keeper Keeper, msg MsgModifyTableAssociation) (*sdk.Result, error) {
    appId, err := keeper.GetDatabaseIdNotFrozen(ctx, msg.AppCode)
    if err != nil {
        return nil, err
    }

    if !isAdmin(ctx, keeper, appId, msg.Owner) {
        return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "Not authorized")
    }

    if !keeper.HasTable(ctx, appId, msg.TableName) {
        return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "Table does not existed")
    }
    err = keeper.ModifyTableAssociation(ctx, appId, msg.TableName, msg.Option, msg.AssociationMode, msg.AssociationTable, msg.ForeignKey, msg.Method)
    if err != nil {
        return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, err.Error())
    }
    return &sdk.Result{}, nil
}

func handleMsgAddCountCache(ctx sdk.Context, keeper Keeper, msg MsgAddCounterCache) (*sdk.Result, error) {
    appId, err := keeper.GetDatabaseIdNotFrozen(ctx, msg.AppCode)
    if err != nil {
        return nil, err
    }

    if !isAdmin(ctx, keeper, appId, msg.Owner) {
        return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "Not authorized")
    }

    if !keeper.HasTable(ctx, appId, msg.TableName) {
        return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "Table does not existed")
    }

    AssociationTable, err := keeper.GetTable(ctx, appId, msg.AssociationTable)
    if err != nil {
        return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, err.Error())
    }
    if !utils.StringIncluded(AssociationTable.Fields, msg.ForeignKey) {
        return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "ForeignKey does not existed")
    }


    err = keeper.AddCounterCache(ctx, appId, msg.TableName, msg.AssociationTable, msg.ForeignKey, msg.CounterCacheField, msg.Limit)
    if err != nil {
        return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, err.Error())
    }
    return &sdk.Result{}, nil
}

func handleMsgDropTable(ctx sdk.Context, keeper Keeper, msg MsgDropTable) (*sdk.Result, error) {
    appId, err := keeper.GetDatabaseIdNotFrozen(ctx, msg.AppCode)
    if err != nil {
        return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, err.Error())
    }

    if !isAdmin(ctx, keeper, appId, msg.Owner) {
        return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "Not authorized")
    }

    if !keeper.HasTable(ctx, appId, msg.TableName) {
        return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "Table name does not exist!")
    }
    keeper.DropTable(ctx, appId, msg.Owner, msg.TableName)
    return &sdk.Result{}, nil
}

func handleMsgAddColumn(ctx sdk.Context, keeper Keeper, msg MsgAddColumn) (*sdk.Result, error) {
    appId, err := keeper.GetDatabaseIdNotFrozen(ctx, msg.AppCode)
    if err != nil {
        return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, err.Error())
    }

    if !isAdmin(ctx, keeper, appId, msg.Owner) {
        return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "Not authorized")
    }

    field := strings.ToLower(msg.Field)
    if keeper.HasField(ctx, appId, msg.TableName, field) {
        return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, fmt.Sprintf("Field %s of table %s exists already!", msg.Field, msg.TableName))
    }
    _, err = keeper.AddColumn(ctx, appId, msg.TableName, field)
    if err != nil {
        return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest,fmt.Sprintf("%v", err))
    }
    return &sdk.Result{}, nil
}

func handleMsgDropColumn(ctx sdk.Context, keeper Keeper, msg MsgDropColumn) (*sdk.Result, error) {
    appId, err := keeper.GetDatabaseIdNotFrozen(ctx, msg.AppCode)
    if err != nil {
        return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, err.Error())
    }

    if !isAdmin(ctx, keeper, appId, msg.Owner) {
        return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "Not authorized")
    }
    if !keeper.HasField(ctx, appId, msg.TableName, msg.Field) {
        return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, fmt.Sprintf("Field %s of table %s does not exist yet!", msg.Field, msg.TableName))
    }

    _, err = keeper.DropColumn(ctx, appId, msg.TableName, msg.Field)
    if err != nil {
        return nil, err
    }
    return &sdk.Result{}, nil
}

func handleMsgRenameColumn(ctx sdk.Context, keeper Keeper, msg MsgRenameColumn) (*sdk.Result, error) {
    appId, err := keeper.GetDatabaseIdNotFrozen(ctx, msg.AppCode)
    if err != nil {
        return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, err.Error())
    }

    if !isAdmin(ctx, keeper, appId, msg.Owner) {
        return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "Not authorized")
    }
    if !keeper.HasField(ctx, appId, msg.TableName, msg.OldField) {
        return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, fmt.Sprintf("Field %s of table %s does not exist yet!", msg.OldField, msg.TableName))
    }

    newField := strings.ToLower(msg.NewField)
    if keeper.HasField(ctx, appId, msg.TableName, newField) {
        return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, fmt.Sprintf("Field %s of table %s exists already!", msg.NewField, msg.TableName))
    }
    _, err = keeper.RenameColumn(ctx, appId, msg.TableName, msg.OldField, newField)
    if err != nil {
        return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest,fmt.Sprintf("%v", err))
    }
    return &sdk.Result{}, nil
}

func handleMsgCreateIndex(ctx sdk.Context, keeper Keeper, msg MsgCreateIndex) (*sdk.Result, error) {
    appId, err := keeper.GetDatabaseIdNotFrozen(ctx, msg.AppCode)
    if err != nil {
        return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, err.Error())
    }

    if !isAdmin(ctx, keeper, appId, msg.Owner) {
        return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "Not authorized")
    }
    if !keeper.HasField(ctx, appId, msg.TableName, msg.Field) {
        return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, fmt.Sprintf("Field %s of table %s does not exist yet!", msg.Field, msg.TableName))
    }
    if err = keeper.CreateIndex(ctx, appId, msg.Owner, msg.TableName, msg.Field); err != nil {
        return nil, err
    }
    return &sdk.Result{}, nil
}

func handleMsgDropIndex(ctx sdk.Context, keeper Keeper, msg MsgDropIndex) (*sdk.Result, error) {
    appId, err := keeper.GetDatabaseIdNotFrozen(ctx, msg.AppCode)
    if err != nil {
        return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, err.Error())
    }

    if !isAdmin(ctx, keeper, appId, msg.Owner) {
        return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "Not authorized")
    }
    if !keeper.HasField(ctx, appId, msg.TableName, msg.Field) {
        return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, fmt.Sprintf("Field %s of table %s does not exist yet!", msg.Field, msg.TableName))
    }

    if err = keeper.DropIndex(ctx, appId, msg.Owner, msg.TableName, msg.Field); err != nil {
        return nil, err
    }
    return &sdk.Result{}, nil
}

func handleMsgModifyOption(ctx sdk.Context, keeper Keeper, msg MsgModifyOption) (*sdk.Result, error) {
    appId, err := keeper.GetDatabaseIdNotFrozen(ctx, msg.AppCode)
    if err != nil {
        return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, err.Error())
    }

    if !isAdmin(ctx, keeper, appId, msg.Owner) {
        return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "Not authorized")
    }
    if !keeper.HasTable(ctx, appId, msg.TableName) {
        return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "Table name does not exist!")
    }

    keeper.ModifyOption(ctx, appId, msg.Owner, msg.TableName, msg.Action, msg.Option)
    return &sdk.Result{}, nil
}

func handleMsgAddInsertFilter(ctx sdk.Context, keeper Keeper, msg MsgAddInsertFilter) (*sdk.Result, error) {
    appId, err := keeper.GetDatabaseIdNotFrozen(ctx, msg.AppCode)
    if err != nil {
        return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, err.Error())
    }

    if !isAdmin(ctx, keeper, appId, msg.Owner) {
        return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "Not authorized")
    }
    if !keeper.HasTable(ctx, appId, msg.TableName) {
        return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "Table name does not exist!")
    }
    bz ,err := base64.StdEncoding.DecodeString(msg.Filter)
    if err != nil {
        return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, err.Error())
    }
    msg.Filter = string(bz)

    if !keeper.AddInsertFilter(ctx, appId, msg.Owner, msg.TableName, msg.Filter) {
        return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "Invalid filter")
    }

    return &sdk.Result{}, nil
}

func handleMsgDropInsertFilter(ctx sdk.Context, keeper Keeper, msg MsgDropInsertFilter) (*sdk.Result, error) {
    appId, err := keeper.GetDatabaseIdNotFrozen(ctx, msg.AppCode)
    if err != nil {
        return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, err.Error())
    }

    if !isAdmin(ctx, keeper, appId, msg.Owner) {
        return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "Not authorized")
    }
    if !keeper.HasTable(ctx, appId, msg.TableName) {
        return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "Table name does not exist!")
    }

    if !keeper.DropInsertFilter(ctx, appId, msg.Owner, msg.TableName) {
        return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "Failed to drop insert filter")
    }

    return &sdk.Result{}, nil
}

func handleMsgAddTrigger(ctx sdk.Context, keeper Keeper, msg MsgAddTrigger) (*sdk.Result, error) {
    appId, err := keeper.GetDatabaseIdNotFrozen(ctx, msg.AppCode)
    if err != nil {
        return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, err.Error())
    }

    if !isAdmin(ctx, keeper, appId, msg.Owner) {
        return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "Not authorized")
    }
    if !keeper.HasTable(ctx, appId, msg.TableName) {
        return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "Table name does not exist!")
    }

    bz, err := base64.StdEncoding.DecodeString(msg.Trigger)
    if err != nil {
        return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, err.Error())
    }
    msg.Trigger = string(bz)
    if !keeper.AddTrigger(ctx, appId, msg.Owner, msg.TableName, msg.Trigger) {
        return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "Invalid trigger")
    }

    return &sdk.Result{}, nil
}

func handleMsgDropTrigger(ctx sdk.Context, keeper Keeper, msg MsgDropTrigger) (*sdk.Result, error) {
    appId, err := keeper.GetDatabaseIdNotFrozen(ctx, msg.AppCode)
    if err != nil {
        return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, err.Error())
    }

    if !isAdmin(ctx, keeper, appId, msg.Owner) {
        return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "Not authorized")
    }
    if !keeper.HasTable(ctx, appId, msg.TableName) {
        return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "Table name does not exist!")
    }

    if !keeper.DropTrigger(ctx, appId, msg.Owner, msg.TableName) {
        return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "Failed to drop trigger")
    }

    return &sdk.Result{}, nil
}

func handleMsgSetTableMemo(ctx sdk.Context, keeper Keeper, msg MsgSetTableMemo) (*sdk.Result, error) {
    appId, err := keeper.GetDatabaseIdNotFrozen(ctx, msg.AppCode)
    if err != nil {
        return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, err.Error())
    }

    if !isAdmin(ctx, keeper, appId, msg.Owner) {
        return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "Not authorized")
    }
    if !keeper.HasTable(ctx, appId, msg.TableName) {
        return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "Table name does not exist!")
    }

    bz ,err := base64.StdEncoding.DecodeString(msg.Memo)
    if err != nil {
        return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, err.Error())
    }
    msg.Memo = string(bz)
    if !keeper.SetTableMemo(ctx, appId, msg.TableName, msg.Memo, msg.Owner) {
        return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "Failed to drop trigger")
    }

    return &sdk.Result{}, nil
}

func handleMsgModifyColumnOption(ctx sdk.Context, keeper Keeper, msg MsgModifyColumnOption) (*sdk.Result, error) {
    appId, err := keeper.GetDatabaseIdNotFrozen(ctx, msg.AppCode)
    if err != nil {
        return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, err.Error())
    }

    if !isAdmin(ctx, keeper, appId, msg.Owner) {
        return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "Not authorized")
    }
    if !keeper.HasTable(ctx, appId, msg.TableName) {
        return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "Table name does not exist!")
    }

    if !keeper.ModifyColumnOption(ctx, appId, msg.Owner, msg.TableName, msg.FieldName, msg.Action, msg.Option) {
        return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "Invalid column option!")
    }
    return &sdk.Result{}, nil
}

func handleMsgSetColumnDataType(ctx sdk.Context, keeper Keeper, msg MsgSetColumnDataType) (*sdk.Result, error) {
    appId, err := keeper.GetDatabaseIdNotFrozen(ctx, msg.AppCode)
    if err != nil {
        return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, err.Error())
    }

    if !isAdmin(ctx, keeper, appId, msg.Owner) {
        return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "Not authorized")
    }
    if !keeper.HasTable(ctx, appId, msg.TableName) {
        return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "Table name does not exist!")
    }

    if !keeper.SetColumnDataType(ctx, appId, msg.Owner, msg.TableName, msg.FieldName, msg.DataType) {
        return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "Invalid column type!")
    }
    return &sdk.Result{}, nil
}

func handleMsgSetColumnMemo(ctx sdk.Context, keeper Keeper, msg MsgSetColumnMemo) (*sdk.Result, error) {
    appId, err := keeper.GetDatabaseIdNotFrozen(ctx, msg.AppCode)
    if err != nil {
        return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, err.Error())
    }

    if !isAdmin(ctx, keeper, appId, msg.Owner) {
        return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "Not authorized")
    }
    if !keeper.HasTable(ctx, appId, msg.TableName) {
        return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "Table name does not exist!")
    }
    bz ,err := base64.StdEncoding.DecodeString(msg.Memo)
    if err != nil {
        return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, err.Error())
    }
    msg.Memo = string(bz)

    if !keeper.SetColumnMemo(ctx, appId, msg.Owner, msg.TableName, msg.FieldName, msg.Memo) {
        return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "Failed to set column memo!")
    }
    return &sdk.Result{}, nil
}

func handleMsgInsertRow(ctx sdk.Context, keeper Keeper, msg types.MsgInsertRow) (*sdk.Result, error) {
    appId, err := keeper.GetDatabaseId(ctx, msg.AppCode)
    if err != nil {
        return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "Invalid app code")
    }

    if !keeper.HasTable(ctx, appId, msg.TableName) {
        return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, fmt.Sprintf("Table %s does not exist!", msg.TableName))
    }
    
    var rowFields types.RowFields
    if err := json.Unmarshal(msg.Fields, &rowFields); err != nil {
        return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest,"Failed to parse row fields!")
    }

    if version.Name == CommunityEdition {
        nextId, _ := keeper.PeekNextId(ctx, appId, msg.TableName)
        if nextId > 1000 {
            return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "No more than 1000 rows allowed")
        }
    }

    _, err = keeper.Insert(ctx, appId, msg.TableName, rowFields, msg.Owner)
    if err != nil {
        return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest,"Failed to insert row : " + err.Error())
    }
    return &sdk.Result{}, nil
}

func handleMsgUpdateRow(ctx sdk.Context, keeper Keeper, msg types.MsgUpdateRow) (*sdk.Result, error) {
    appId, err := keeper.GetDatabaseId(ctx, msg.AppCode)
    if err != nil {
        return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest,"Invalid app code")
    }

    if !keeper.HasTable(ctx, appId, msg.TableName) {
        return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest,fmt.Sprintf("Table %s does not exist!", msg.TableName))
    }

    options, _ := keeper.GetOption(ctx, appId, msg.TableName)
    if ! utils.ItemExists(options, string(types.TBLOPT_UPDATABLE)) {
        return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest,fmt.Sprintf("Table %s is not updatable!", msg.TableName))
    }

    var rowFields types.RowFields
    if err := json.Unmarshal(msg.Fields, &rowFields); err != nil {
        return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest,"Failed to parse row fields!")
    }

    keeper.Update(ctx, appId, msg.TableName, msg.Id, rowFields, msg.Owner)
    return &sdk.Result{}, nil
}

func handleMsgDeleteRow(ctx sdk.Context, keeper Keeper, msg types.MsgDeleteRow) (*sdk.Result, error) {
    appId, err := keeper.GetDatabaseId(ctx, msg.AppCode)
    if err != nil {
        return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest,"Invalid app code")
    }

    if !keeper.HasTable(ctx, appId, msg.TableName) {
        return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest,fmt.Sprintf("Table %s does not exist!", msg.TableName))
    }

    options, _ := keeper.GetOption(ctx, appId, msg.TableName)
    if ! utils.ItemExists(options, string(types.TBLOPT_DELETABLE)) {
        return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest,fmt.Sprintf("Table %s is not updatable!", msg.TableName))
    }

    _, err = keeper.Delete(ctx, appId, msg.TableName, msg.Id, msg.Owner)
    if err != nil {
        return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest,fmt.Sprintf("%v", err))
    }
    return &sdk.Result{}, nil
}

func handleMsgFreezeRow(ctx sdk.Context, keeper Keeper, msg types.MsgFreezeRow) (*sdk.Result, error) {
    appId, err := keeper.GetDatabaseId(ctx, msg.AppCode)
    if err != nil {
        return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest,"Invalid app code")
    }

    if !keeper.HasTable(ctx, appId, msg.TableName) {
        return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest,fmt.Sprintf("Table %s does not exist!", msg.TableName))
    }

    _, err = keeper.Freeze(ctx, appId, msg.TableName, msg.Id, msg.Owner)
    if err != nil {
        return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest,fmt.Sprintf("%v", err))
    }
    return &sdk.Result{}, nil
}

func handleMsgModifyGroup(ctx sdk.Context, keeper Keeper, msg MsgModifyGroup) (*sdk.Result, error) {
    appId, err := keeper.GetDatabaseId(ctx, msg.AppCode)
    if err != nil {
        return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest,"Invalid app code")
    }

    if msg.Group == "admin" {
        if !isSysAdmin(ctx, keeper, msg.Owner) {
            return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest,"Not authorized")
        }
    } else {
        if !isAdmin(ctx, keeper, appId, msg.Owner) {
            return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest,"Not authorized")
        }
    }

    err = keeper.ModifyGroup(ctx, appId, msg.Action, msg.Group)
    if err != nil {
        return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest,fmt.Sprintf("%v", err))
    }
    return &sdk.Result{}, nil
}

func handleMsgSetGroupMemo(ctx sdk.Context, keeper Keeper, msg MsgSetGroupMemo) (*sdk.Result, error) {
    appId, err := keeper.GetDatabaseId(ctx, msg.AppCode)
    if err != nil {
        return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest,"Invalid app code")
    }

    if msg.Group == "admin" {
        return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest,"No need to set memo for admin")
    }

    if !isAdmin(ctx, keeper, appId, msg.Owner) {
        return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest,"Not authorized")
    }

    bz ,err := base64.StdEncoding.DecodeString(msg.Memo)
    if err != nil {
        return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest,err.Error())
    }
    msg.Memo = string(bz)
    keeper.SetGroupMemo(ctx, appId, msg.Group, msg.Memo)
    return &sdk.Result{}, nil
}

func handleMsgModifyGroupMember(ctx sdk.Context, keeper Keeper, msg MsgModifyGroupMember) (*sdk.Result, error) {
    appId, err := keeper.GetDatabaseId(ctx, msg.AppCode)
    if err != nil {
        return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest,"Invalid app code")
    }

    if msg.Group == "admin" {
        if !(isSysAdmin(ctx, keeper, msg.Owner) || isAdmin(ctx, keeper, appId, msg.Owner)) {
            return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest,"Not authorized")
        }
    } else {
        if !isAdmin(ctx, keeper, appId, msg.Owner) {
            return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest,"Not authorized")
        }
    }

    err = keeper.ModifyGroupMember(ctx, appId, msg.Group, msg.Action, msg.Member)
    if err != nil {
        return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest,fmt.Sprintf("%v", err))
    }
    return &sdk.Result{}, nil
}

func handleMsgAddFriend(ctx sdk.Context, keeper Keeper, msg MsgAddFriend) (*sdk.Result, error) {
    err := keeper.AddFriend(ctx, msg.Owner, msg.OwnerName, msg.FriendAddr, msg.FriendName)
    if err != nil {
        return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest,fmt.Sprintf("%v", err))
    }
    return &sdk.Result{}, nil
}

func handleMsgDropFriend(ctx sdk.Context, keeper Keeper, msg MsgDropFriend) (*sdk.Result, error) {
    err := keeper.DropFriend(ctx, msg.Owner, msg.FriendAddr)
    if err != nil {
        return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest,fmt.Sprintf("%v", err))
    }
    return &sdk.Result{}, nil
}

func handleMsgRespondFriend(ctx sdk.Context, keeper Keeper, msg MsgRespondFriend) (*sdk.Result, error) {
    err := keeper.RespondFriend(ctx, msg.Owner, msg.FriendAddr, msg.Action)
    if err != nil {
        return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest,fmt.Sprintf("%v", err))
    }
    return &sdk.Result{}, nil
}

func handleMsgSetSchemaStatus(ctx sdk.Context, keeper Keeper, msg MsgSetSchemaStatus) (*sdk.Result, error) {
    appId, err := keeper.GetDatabaseId(ctx, msg.AppCode)
    if err != nil {
        return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest,"Invalid app code")
    }

    if !isSysAdmin(ctx, keeper, msg.Owner) && !isAdmin(ctx, keeper, appId, msg.Owner) {
        return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "Not authorized")
    }

    err = keeper.SetSchemaStatus(ctx, msg.Owner, msg.AppCode, msg.Status)
    if err != nil {
        return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest,fmt.Sprintf("%v", err))
    }
    return &sdk.Result{}, nil
}

func handleMsgSetDatabasePermission(ctx sdk.Context, keeper Keeper, msg MsgSetDatabasePermission) (*sdk.Result, error) {
    appId, err := keeper.GetDatabaseIdNotFrozen(ctx, msg.AppCode)
    if err != nil {
        return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest,err.Error())
    }

    if !isSysAdmin(ctx, keeper, msg.Owner) && !isAdmin(ctx, keeper, appId, msg.Owner) {
        return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "Not authorized")
    }

    err = keeper.SetDatabasePermission(ctx, msg.Owner, msg.AppCode, msg.PermissionRequired)
    if err != nil {
        return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest,fmt.Sprintf("%v", err))
    }
    return &sdk.Result{}, nil
}

////////////////////////
//                    //
// blockchain browser //
//                    //
////////////////////////

func handleMsgUpdateTotalTx(ctx sdk.Context, keeper Keeper, msg MsgUpdateTotalTx) (*sdk.Result, error) {

    appId, err := keeper.GetDatabaseId(ctx, "0000000001")
    if err != nil {
       return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest,"Invalid app code")
    }

    if !keeper.IsGroupMember(ctx, appId, "oracle", msg.Owner) {
        return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest,"Permission forbidden")
    }
    err = keeper.UpdateTotalTx(ctx, msg.Data)
    if err != nil {
        return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest,fmt.Sprintf("%v", err))
    }
    return &sdk.Result{}, nil
}

func handleMsgUpdateTxStatistic(ctx sdk.Context, keeper Keeper, msg MsgUpdateTxStatistic) (*sdk.Result, error) {
    appId, err := keeper.GetDatabaseId(ctx, "0000000001")
    if err != nil {
        return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest,"Invalid app code")
    }

    if !keeper.IsGroupMember(ctx, appId, "oracle", msg.Owner) {
        return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest,"Permission forbidden")
    }

    err = keeper.UpdateTxStatistic(ctx, msg.Data)
    if err != nil {
        return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest,fmt.Sprintf("%v", err))
    }
    return &sdk.Result{}, nil
}

func handleMsgModifyP2PTransferLimit(ctx sdk.Context, keeper Keeper, msg MsgModifyP2PTransferLimit) (*sdk.Result, error) {
    //only chain super admin can set limit
    err := keeper.SetP2PTransferLimit(ctx, msg.Owner, msg.Limit)
    if err != nil {
        return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest,fmt.Sprintf("%v", err))
    }
    return &sdk.Result{}, nil
}

func handleMsgModifyTokenKeeperMember(ctx sdk.Context, keeper Keeper, msg MsgModifyTokenKeeperMember) (*sdk.Result, error) {
    err := keeper.ModifyMemberOfTokenKeepers(ctx, msg.Owner, msg.Member, msg.Action)
    if err != nil {
        return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest,fmt.Sprintf("%v", err))
    }
    return &sdk.Result{}, nil
}

func handleMsgSaveUserPrivateKey(ctx sdk.Context, keeper Keeper, msg MsgSaveUserPrivateKey) (*sdk.Result, error) {
    err := keeper.SaveUserPrivateInfo(ctx, msg.Owner, msg.User, msg.KeyInfo)
    if err != nil {
        return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest,fmt.Sprintf("%v", err))
    }
    return &sdk.Result{}, nil
}

////////////////////
//                //
// helper methods //
//                //
////////////////////

func isSysAdmin(ctx sdk.Context, keeper Keeper, address sdk.AccAddress) bool {
    sysAdmins := keeper.GetSysAdmins(ctx)
    var is_sysAdmin = false
    for _, addr := range sysAdmins {
        if bytes.Compare(address, addr) == 0 {
            is_sysAdmin = true
            break
        }
    }
    return is_sysAdmin
}

func isAdmin(ctx sdk.Context, keeper Keeper, appId uint, address sdk.AccAddress) bool {
    adminAddresses := keeper.GetDatabaseAdmins(ctx, appId)
    var is_admin = false
    for _, addr := range adminAddresses {
        if bytes.Compare(address, addr) == 0 {
            is_admin = true
            break
        }
    }
    return is_admin
}
