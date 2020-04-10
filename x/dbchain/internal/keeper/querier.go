package keeper

import (
    "fmt"
    "strconv"
    //"encoding/hex"

    "github.com/cosmos/cosmos-sdk/codec"
    sdk "github.com/cosmos/cosmos-sdk/types"
    sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
    abci "github.com/tendermint/tendermint/abci/types"
    "github.com/yzhanginwa/dbchain/x/dbchain/internal/utils"
)

// query endpoints supported by the dbchain service Querier
const (
    QueryApplication   = "application"
    QueryAppUsers      = "app_users"
    QueryIsAppUser     = "is_app_user"
    QueryAdminApps     = "admin_apps"
    QueryTables   = "tables"
    QueryIndex    = "index"
    QueryOption   = "option"
    QueryColumnOption   = "column_option"
    QueryRow      = "find"
    QueryIdsBy    = "find_by"
    QueryAllIds   = "find_all"
    QueryAdminGroup = "admin_group"
    QueryFriends  = "friends"
    QueryPendingFriends  = "pending_friends"
)


// NewQuerier is the module level router for state queries
func NewQuerier(keeper Keeper) sdk.Querier {
    return func(ctx sdk.Context, path []string, req abci.RequestQuery) (res []byte, err error) {
        switch path[0] {
        case QueryApplication:
            if len(path) > 2 {
                return queryApplication(ctx, path[1:], req, keeper)
            } else {
                return queryApplications(ctx, path[1:], req, keeper)
            }
        case QueryAppUsers:
            return queryAppUsers(ctx, path[1:], req, keeper)
        case QueryIsAppUser:
            return queryIsAppUser(ctx, path[1:], req, keeper)
        case QueryAdminApps:
            return queryAdminApps(ctx, path[1:], req, keeper)
        case QueryTables:
            if len(path) > 3 {
                return queryTable(ctx, path[1:], req, keeper)
            } else {
                return queryTables(ctx, path[1:], req, keeper)
            }
        case QueryIndex:
            return queryIndex(ctx, path[1:], req, keeper)
        case QueryOption:
            return queryOption(ctx, path[1:], req, keeper)
        case QueryColumnOption:
            return queryColumnOption(ctx, path[1:], req, keeper)
        case QueryRow:
            return queryRow(ctx, path[1:], req, keeper)
        case QueryIdsBy:
            return queryIdsBy(ctx, path[1:], req, keeper)
        case QueryAllIds:
            return queryAllIds(ctx, path[1:], req, keeper)
        case QueryAdminGroup:
            return queryAdminGroup(ctx, path[1:], req, keeper)
        case QueryFriends:
            return queryFriends(ctx, path[1:], req, keeper)
        case QueryPendingFriends:
            return queryPendingFriends(ctx, path[1:], req, keeper)
        default:
            return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "unknown dbchain query endpoint")
        }
    }
}

////////////////////////////////
//                            //
// query application/database //
//                            //
////////////////////////////////

// the the list of app code in the system
func queryApplications(ctx sdk.Context, path []string, req abci.RequestQuery, keeper Keeper) ([]byte, error) {
    accessCode:= path[0]
    _, err := utils.VerifyAccessCode(accessCode)
    if err != nil {
        return []byte{}, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "Access code is not valid!")
    }

    // we use the term database in the code
    applications := keeper.GetAllAppCode(ctx)

    res, err := codec.MarshalJSONIndent(keeper.cdc, applications)
    if err != nil {
        panic("could not marshal result to JSON")
    }

    return res, nil
}

func queryApplication(ctx sdk.Context, path []string, req abci.RequestQuery, keeper Keeper) ([]byte, error) {
    accessCode:= path[0]
    _, err := utils.VerifyAccessCode(accessCode)
    if err != nil {
        return []byte{}, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "Access code is not valid!")
    }

    appCode := path[1]
    database, err := keeper.getDatabase(ctx, appCode)

    if err != nil {
        return []byte{}, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, fmt.Sprintf("AppCode %s does not exist", appCode))
    }

    res, err := codec.MarshalJSONIndent(keeper.cdc, database)
    if err != nil {
        panic("could not marshal result to JSON")
    }

    return res, nil
}

func queryAppUsers(ctx sdk.Context, path []string, req abci.RequestQuery, keeper Keeper) ([]byte, error) {
    accessCode:= path[0]
    addr, err := utils.VerifyAccessCode(accessCode)
    if err != nil {
        return []byte{}, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "Access code is not valid!")
    }

    appCode := path[1]
    appId, err := keeper.GetDatabaseId(ctx, appCode)
    if err != nil {
        return []byte{}, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "Invalid app code")
    }

    users := keeper.GetDatabaseUsers(ctx, appId, addr)

    res, err := codec.MarshalJSONIndent(keeper.cdc, users)
    if err != nil {
        panic("could not marshal result to JSON")
    }

    return res, nil
}

func queryIsAppUser(ctx sdk.Context, path []string, req abci.RequestQuery, keeper Keeper) ([]byte, error) {
    accessCode:= path[0]
    addr, err := utils.VerifyAccessCode(accessCode)
    if err != nil {
        return []byte{}, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "Access code is not valid!")
    }

    appCode := path[1]
    appId, err := keeper.GetDatabaseId(ctx, appCode)
    if err != nil {
        return []byte{}, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "Invalid app code")
    }

    isAppUser := keeper.IsDatabaseUser(ctx, appId, addr)

    res, err := codec.MarshalJSONIndent(keeper.cdc, isAppUser)
    if err != nil {
        panic("could not marshal result to JSON")
    }

    return res, nil
}

func queryAdminApps(ctx sdk.Context, path []string, req abci.RequestQuery, keeper Keeper) ([]byte, error) {
    accessCode:= path[0]
    addr, err := utils.VerifyAccessCode(accessCode)
    if err != nil {
        return []byte{}, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "Access code is not valid!")
    }

    // we use the term database in the code
    adminApps := keeper.getAdminAppCode(ctx, addr)

    res, err := codec.MarshalJSONIndent(keeper.cdc, adminApps)
    if err != nil {
        panic("could not marshal result to JSON")
    }

    return res, nil
}

////////////////
//            //
// query meta //
//            //
////////////////

func queryTables(ctx sdk.Context, path []string, req abci.RequestQuery, keeper Keeper) ([]byte, error) {
    accessCode:= path[0]
    _, err := utils.VerifyAccessCode(accessCode)
    if err != nil {
        return []byte{}, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "Access code is not valid!")
    }

    appId, err := keeper.GetDatabaseId(ctx, path[1])
    if err != nil {
        return []byte{}, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "Invalid app code")
    }

    tables, err := keeper.GetTables(ctx, appId)

    if err != nil {
        return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "Can not get table names")
    }

    res, err := codec.MarshalJSONIndent(keeper.cdc, tables)
    if err != nil {
        panic("could not marshal result to JSON")
    }

    return res, nil
}

func queryTable(ctx sdk.Context, path []string, req abci.RequestQuery, keeper Keeper) ([]byte, error) {
    accessCode:= path[0]
    _, err := utils.VerifyAccessCode(accessCode)
    if err != nil {
        return []byte{}, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "Access code is not valid!")
    }

    appId, err := keeper.GetDatabaseId(ctx, path[1])
    if err != nil {
        return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "Invalid app code")
    }

    table, err := keeper.GetTable(ctx, appId, path[2])

    if err != nil {
        return []byte{}, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, fmt.Sprintf("Table %s does not exist",  table))
    }

    res, err := codec.MarshalJSONIndent(keeper.cdc, table)
    if err != nil {
        panic("could not marshal result to JSON")
    }

    return res, nil
}

func queryIndex(ctx sdk.Context, path []string, req abci.RequestQuery, keeper Keeper) ([]byte, error) {
    accessCode:= path[0]
    _, err := utils.VerifyAccessCode(accessCode)
    if err != nil {
        return []byte{}, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "Access code is not valid!")
    }

    appId, err := keeper.GetDatabaseId(ctx, path[1])
    if err != nil {
        return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "Invalid app code")
    }

    tableName := path[2]
    index, err := keeper.GetIndex(ctx, appId, tableName)

    if err != nil {
        return []byte{}, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, fmt.Sprintf("Table %s does not exist",  tableName))
    }

    res, err := codec.MarshalJSONIndent(keeper.cdc, index)
    if err != nil {
        panic("could not marshal result to JSON")
    }

    return res, nil
}

func queryOption(ctx sdk.Context, path []string, req abci.RequestQuery, keeper Keeper) ([]byte, error) {
    accessCode:= path[0]
    _, err := utils.VerifyAccessCode(accessCode)
    if err != nil {
        return []byte{}, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "Access code is not valid!")
    }

    appId, err := keeper.GetDatabaseId(ctx, path[1])
    if err != nil {
        return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "Invalid app code")
    }

    tableName := path[2]
    options, err := keeper.GetOption(ctx, appId, tableName)

    if err != nil {
        return []byte{}, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, fmt.Sprintf("Table %s does not exist",  tableName))
    }

    res, err := codec.MarshalJSONIndent(keeper.cdc, options)
    if err != nil {
        panic("could not marshal result to JSON")
    }

    return res, nil
}

func queryColumnOption(ctx sdk.Context, path []string, req abci.RequestQuery, keeper Keeper) ([]byte, error) {
    accessCode:= path[0]
    _, err := utils.VerifyAccessCode(accessCode)
    if err != nil {
        return []byte{}, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "Access code is not valid!")
    }

    appId, err := keeper.GetDatabaseId(ctx, path[1])
    if err != nil {
        return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "Invalid app code")
    }

    tableName := path[2]
    fieldName := path[3]

    options, err := keeper.GetColumnOption(ctx, appId, tableName, fieldName)

    if err != nil {
        return []byte{}, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, fmt.Sprintf("Field %s.%s does not exist",  tableName, fieldName))
    }

    res, err := codec.MarshalJSONIndent(keeper.cdc, options)
    if err != nil {
        panic("could not marshal result to JSON")
    }

    return res, nil
}

////////////////
//            //
// query data //
//            //
////////////////

func queryRow(ctx sdk.Context, path []string, req abci.RequestQuery, keeper Keeper) ([]byte, error) {
    accessCode:= path[0]
    addr, err := utils.VerifyAccessCode(accessCode)
    if err != nil {
        return []byte{}, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "Access code is not valid!")
    }

    appId, err := keeper.GetDatabaseId(ctx, path[1])
    if  err != nil {
        return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "Invalid app code")
    }

    tableName := path[2] 
    u32, err := strconv.ParseUint(path[3], 10, 32)
    fields, err := keeper.Find(ctx, appId, tableName, uint(u32), addr)

    if err != nil {
        return []byte{}, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, fmt.Sprintf("Table %s does not exist",  path[2]))
    }

    res, err := codec.MarshalJSONIndent(keeper.cdc, fields)
    if err != nil {
        panic("could not marshal result to JSON")
    }

    return res, nil
}

func queryIdsBy(ctx sdk.Context, path []string, req abci.RequestQuery, keeper Keeper) ([]byte, error) {
    accessCode:= path[0]
    addr, err := utils.VerifyAccessCode(accessCode)
    if err != nil {
        return []byte{}, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "Access code is not valid!")
    }

    appId, err := keeper.GetDatabaseId(ctx, path[1])
    if err != nil {
        return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "Invalid app code")
    }

    tableName := path[2]
    fieldName := path[3]
    value := path[4]
    ids := keeper.FindBy(ctx, appId, tableName, fieldName, value, addr)

    res, err := codec.MarshalJSONIndent(keeper.cdc, ids)
    if err != nil {
        panic("could not marshal result to JSON")
    }

    return res, nil
}

func queryAllIds(ctx sdk.Context, path []string, req abci.RequestQuery, keeper Keeper) ([]byte, error) {
    accessCode:= path[0]
    addr, err := utils.VerifyAccessCode(accessCode)
    if err != nil {
        return []byte{}, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "Access code is not valid!")
    }

    appId, err := keeper.GetDatabaseId(ctx, path[1])
    if err != nil {
        return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "Invalid app code")
    }

    var tableName = path[2]

    ids := keeper.FindAll(ctx, appId, tableName, addr)

    res, err := codec.MarshalJSONIndent(keeper.cdc, ids)
    if err != nil {
        panic("could not marshal result to JSON")
    }

    return res, nil
}

///////////////////
//               //
// query friends //
//               //
///////////////////

func queryFriends(ctx sdk.Context, path []string, req abci.RequestQuery, keeper Keeper) ([]byte, error) {
    accessCode:= path[0]
    addr, err := utils.VerifyAccessCode(accessCode)
    if err != nil {
        return []byte{}, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "Access code is not valid!")
    }

    friends := keeper.GetFriends(ctx, addr)
    res, err := codec.MarshalJSONIndent(keeper.cdc, friends)
    if err != nil {
        panic("could not marshal result to JSON")
    }

    return res, nil
}

func queryPendingFriends(ctx sdk.Context, path []string, req abci.RequestQuery, keeper Keeper) ([]byte, error) {
    accessCode:= path[0]
    addr, err := utils.VerifyAccessCode(accessCode)
    if err != nil {
        return []byte{}, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "Access code is not valid!")
    }

    friends := keeper.GetPendingFriends(ctx, addr)
    res, err := codec.MarshalJSONIndent(keeper.cdc, friends)
    if err != nil {
        panic("could not marshal result to JSON")
    }

    return res, nil
}

/////////////////
//             //
// query group //
//             //
/////////////////

func queryAdminGroup(ctx sdk.Context, path []string, req abci.RequestQuery, keeper Keeper) ([]byte, error) {
    accessCode:= path[0]
    _, err := utils.VerifyAccessCode(accessCode)
    if err != nil {
        return []byte{}, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "Access code is not valid!")
    }

    appId, err := keeper.GetDatabaseId(ctx, path[1])
    if err != nil {
        return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "Invalid app code")
    }

    adminAddresses := keeper.ShowAdminGroup(ctx, appId)

    res, err := codec.MarshalJSONIndent(keeper.cdc, adminAddresses)
    if err != nil {
        panic("could not marshal result to JSON")
    }

    return res, nil
}


//////////////////
//              //
// helper funcs //
//              //
//////////////////

