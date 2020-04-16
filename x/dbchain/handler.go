package dbchain

import (
    "fmt"
    "strings"
    "bytes"
    "encoding/json"
    sdk "github.com/cosmos/cosmos-sdk/types"
    sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
    "github.com/yzhanginwa/dbchain/x/dbchain/internal/types"
    "github.com/yzhanginwa/dbchain/x/dbchain/internal/utils"
    "github.com/cosmos/cosmos-sdk/version"
)

const (
    CommunityEdition = "dbChainCommunity"
)

var (
    AllowCreateApplication bool
)

// NewHandler returns a handler for "nameservice" type messages.
func NewHandler(keeper Keeper) sdk.Handler {
    return func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
        switch msg := msg.(type) {
        case MsgCreateApplication:
            return handleMsgCreateApplication(ctx, keeper, msg)
        case MsgAddDatabaseUser:
            return handleMsgAddDatabaseUser(ctx, keeper, msg)
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
        case MsgFreezeRow:
            return handleMsgFreezeRow(ctx, keeper, msg)
        case MsgModifyGroup:
            return handleMsgModifyGroup(ctx, keeper, msg)
        case MsgModifyGroupMember:
            return handleMsgModifyGroupMember(ctx, keeper, msg)
        case MsgAddFriend:
            return handleMsgAddFriend(ctx, keeper, msg)
        case MsgRespondFriend:
            return handleMsgRespondFriend(ctx, keeper, msg)
        default:
            errMsg := fmt.Sprintf("Unrecognized dbchain Msg type: %v", msg.Type())
            return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, errMsg)
        }
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
        if len(apps) > 9 {
            return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "No more than 10 apps allowed")
        }
    }
    // We use the term database for internal use. To outside we use application to make users understand easily
    keeper.CreateDatabase(ctx, msg.Owner, msg.Name, msg.Description, msg.Permissioned)
    return &sdk.Result{}, nil
}

func handleMsgAddDatabaseUser(ctx sdk.Context, keeper Keeper, msg MsgAddDatabaseUser) (*sdk.Result, error) {
    // We use the term database for internal use. To outside we use application to make users understand easily
    keeper.AddDatabaseUser(ctx, msg.Owner, msg.AppCode, msg.User)
    return &sdk.Result{}, nil
}

// Handle a message to create table 
func handleMsgCreateTable(ctx sdk.Context, keeper Keeper, msg MsgCreateTable) (*sdk.Result, error) {
    appId, err := keeper.GetDatabaseId(ctx, msg.AppCode)
    if err != nil {
        return nil, err
    }

    if !isAdmin(ctx, keeper, msg.AppCode, msg.Owner) {
        return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "Not authorized")
    }
 
    if version.Name == CommunityEdition {
        tables, err := keeper.GetTables(ctx, appId)
        if err != nil {
            return nil, err
        }
        if len(tables) > 29 {
            return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "No more than 30 tables allowed")
        }
    }

    if keeper.HasTable(ctx, appId, msg.TableName) {
        return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "Table name existed already!")
    }
    keeper.CreateTable(ctx, appId, msg.Owner, msg.TableName, msg.Fields)
    return &sdk.Result{}, nil
}

func handleMsgDropTable(ctx sdk.Context, keeper Keeper, msg MsgDropTable) (*sdk.Result, error) {
    appId, err := keeper.GetDatabaseId(ctx, msg.AppCode)
    if err != nil {
        return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "Invalid app code")
    }

    if !isAdmin(ctx, keeper, msg.AppCode, msg.Owner) {
        return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "Not authorized")
    }

    if !keeper.HasTable(ctx, appId, msg.TableName) {
        return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "Table name does not exist!")
    }
    keeper.DropTable(ctx, appId, msg.Owner, msg.TableName)
    return &sdk.Result{}, nil
}

func handleMsgAddColumn(ctx sdk.Context, keeper Keeper, msg MsgAddColumn) (*sdk.Result, error) {
    appId, err := keeper.GetDatabaseId(ctx, msg.AppCode)
    if err != nil {
        return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "Invalid app code")
    }

    if !isAdmin(ctx, keeper, msg.AppCode, msg.Owner) {
        return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "Not authorized")
    }

    field := strings.ToLower(msg.Field)
    if keeper.HasField(ctx, appId, msg.TableName, field) {
        return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, fmt.Sprintf("Field %s of table %s exists already!", msg.Field, msg.TableName))
    }
    keeper.AddColumn(ctx, appId, msg.TableName, field)
    return &sdk.Result{}, nil
}

func handleMsgDropColumn(ctx sdk.Context, keeper Keeper, msg MsgDropColumn) (*sdk.Result, error) {
    appId, err := keeper.GetDatabaseId(ctx, msg.AppCode)
    if err != nil {
        return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "Invalid app code")
    }

    if !isAdmin(ctx, keeper, msg.AppCode, msg.Owner) {
        return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "Not authorized")
    }
    if !keeper.HasField(ctx, appId, msg.TableName, msg.Field) {
        return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, fmt.Sprintf("Field %s of table %s does not exist yet!", msg.Field, msg.TableName))
    }
    keeper.DropColumn(ctx, appId, msg.TableName, msg.Field)
    return &sdk.Result{}, nil
}

func handleMsgRenameColumn(ctx sdk.Context, keeper Keeper, msg MsgRenameColumn) (*sdk.Result, error) {
    appId, err := keeper.GetDatabaseId(ctx, msg.AppCode)
    if err != nil {
        return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "Invalid app code")
    }

    if !isAdmin(ctx, keeper, msg.AppCode, msg.Owner) {
        return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "Not authorized")
    }
    if !keeper.HasField(ctx, appId, msg.TableName, msg.OldField) {
        return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, fmt.Sprintf("Field %s of table %s does not exist yet!", msg.OldField, msg.TableName))
    }

    newField := strings.ToLower(msg.NewField)
    if keeper.HasField(ctx, appId, msg.TableName, newField) {
        return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, fmt.Sprintf("Field %s of table %s exists already!", msg.NewField, msg.TableName))
    }
    keeper.RenameColumn(ctx, appId, msg.TableName, msg.OldField, newField)
    return &sdk.Result{}, nil
}

func handleMsgCreateIndex(ctx sdk.Context, keeper Keeper, msg MsgCreateIndex) (*sdk.Result, error) {
    appId, err := keeper.GetDatabaseId(ctx, msg.AppCode)
    if err != nil {
        return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "Invalid app code")
    }

    if !isAdmin(ctx, keeper, msg.AppCode, msg.Owner) {
        return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "Not authorized")
    }
    if !keeper.HasField(ctx, appId, msg.TableName, msg.Field) {
        return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, fmt.Sprintf("Field %s of table %s does not exist yet!", msg.Field, msg.TableName))
    }
    keeper.CreateIndex(ctx, appId, msg.Owner, msg.TableName, msg.Field)
    return &sdk.Result{}, nil
}

func handleMsgDropIndex(ctx sdk.Context, keeper Keeper, msg MsgDropIndex) (*sdk.Result, error) {
    appId, err := keeper.GetDatabaseId(ctx, msg.AppCode)
    if err != nil {
        return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "Invalid app code")
    }

    if !isAdmin(ctx, keeper, msg.AppCode, msg.Owner) {
        return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "Not authorized")
    }
    if !keeper.HasField(ctx, appId, msg.TableName, msg.Field) {
        return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, fmt.Sprintf("Field %s of table %s does not exist yet!", msg.Field, msg.TableName))
    }

    existingIndex, err := keeper.GetIndex(ctx, appId, msg.TableName)
    if err != nil {
        return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, fmt.Sprintf("Table %s does not have any index yet!", msg.TableName))
    }
 
    if !utils.ItemExists(existingIndex, msg.Field) {
        return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, fmt.Sprintf("Table %s does not have index on %s yet!", msg.TableName, msg.Field))
    }

    keeper.DropIndex(ctx, appId, msg.Owner, msg.TableName, msg.Field)
    return &sdk.Result{}, nil
}

func handleMsgModifyOption(ctx sdk.Context, keeper Keeper, msg MsgModifyOption) (*sdk.Result, error) {
    appId, err := keeper.GetDatabaseId(ctx, msg.AppCode)
    if err != nil {
        return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "Invalid app code")
    }

    if !isAdmin(ctx, keeper, msg.AppCode, msg.Owner) {
        return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "Not authorized")
    }
    if !keeper.HasTable(ctx, appId, msg.TableName) {
        return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "Table name does not exist!")
    }

    keeper.ModifyOption(ctx, appId, msg.Owner, msg.TableName, msg.Action, msg.Option)
    return &sdk.Result{}, nil
}

func handleMsgModifyColumnOption(ctx sdk.Context, keeper Keeper, msg MsgModifyColumnOption) (*sdk.Result, error) {
    appId, err := keeper.GetDatabaseId(ctx, msg.AppCode)
    if err != nil {
        return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "Invalid app code")
    }

    if !isAdmin(ctx, keeper, msg.AppCode, msg.Owner) {
        return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "Not authorized")
    }
    if !keeper.HasTable(ctx, appId, msg.TableName) {
        return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "Table name does not exist!")
    }

    keeper.ModifyColumnOption(ctx, appId, msg.Owner, msg.TableName, msg.FieldName, msg.Action, msg.Option)
    return &sdk.Result{}, nil
}

func handleMsgInsertRow(ctx sdk.Context, keeper Keeper, msg types.MsgInsertRow) (*sdk.Result, error) {
    appId, err := keeper.GetDatabaseId(ctx, msg.AppCode)
    if err != nil {
        return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "Invalid app code")
    }

    if !keeper.HasTable(ctx, appId, msg.TableName) {
        return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, fmt.Sprintf("Table % does not exist!", msg.TableName))
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
        return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest,"Failed validation of inserting row")
    }
    return &sdk.Result{}, nil
}

func handleMsgUpdateRow(ctx sdk.Context, keeper Keeper, msg types.MsgUpdateRow) (*sdk.Result, error) {
    appId, err := keeper.GetDatabaseId(ctx, msg.AppCode)
    if err != nil {
        return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest,"Invalid app code")
    }

    if !keeper.HasTable(ctx, appId, msg.TableName) {
        return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest,fmt.Sprintf("Table % does not exist!", msg.TableName))
    }

    options, _ := keeper.GetOption(ctx, appId, msg.TableName)
    if ! utils.ItemExists(options, string(types.TBLOPT_UPDATABLE)) {
        return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest,fmt.Sprintf("Table % is not updatable!", msg.TableName))
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
        return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest,fmt.Sprintf("Table % does not exist!", msg.TableName))
    }

    options, _ := keeper.GetOption(ctx, appId, msg.TableName)
    if ! utils.ItemExists(options, string(types.TBLOPT_DELETABLE)) {
        return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest,fmt.Sprintf("Table % is not updatable!", msg.TableName))
    }

    keeper.Delete(ctx, appId, msg.TableName, msg.Id, msg.Owner)
    return &sdk.Result{}, nil
}

func handleMsgFreezeRow(ctx sdk.Context, keeper Keeper, msg types.MsgFreezeRow) (*sdk.Result, error) {
    appId, err := keeper.GetDatabaseId(ctx, msg.AppCode)
    if err != nil {
        return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest,"Invalid app code")
    }

    if !keeper.HasTable(ctx, appId, msg.TableName) {
        return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest,fmt.Sprintf("Table % does not exist!", msg.TableName))
    }

    keeper.Freeze(ctx, appId, msg.TableName, msg.Id, msg.Owner)
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
        if !isAdmin(ctx, keeper, msg.AppCode, msg.Owner) {
            return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest,"Not authorized")
        }
    }

    err = keeper.ModifyGroup(ctx, appId, msg.Action, msg.Group)
    if err != nil {
        return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest,fmt.Sprintf("%v", err))
    }
    return &sdk.Result{}, nil
}

func handleMsgModifyGroupMember(ctx sdk.Context, keeper Keeper, msg MsgModifyGroupMember) (*sdk.Result, error) {
    appId, err := keeper.GetDatabaseId(ctx, msg.AppCode)
    if err != nil {
        return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest,"Invalid app code")
    }

    if msg.Group == "admin" {
        if !isSysAdmin(ctx, keeper, msg.Owner) {
            return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest,"Not authorized")
        }
    } else {
        if !isAdmin(ctx, keeper, msg.AppCode, msg.Owner) {
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

func handleMsgRespondFriend(ctx sdk.Context, keeper Keeper, msg MsgRespondFriend) (*sdk.Result, error) {
    err := keeper.RespondFriend(ctx, msg.Owner, msg.FriendAddr, msg.Action)
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

