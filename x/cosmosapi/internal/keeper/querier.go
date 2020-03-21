package keeper

import (
    "fmt"
    "strconv"
    //"encoding/hex"

    "github.com/cosmos/cosmos-sdk/codec"
    sdk "github.com/cosmos/cosmos-sdk/types"
    abci "github.com/tendermint/tendermint/abci/types"
    "github.com/yzhanginwa/cosmos-api/x/cosmosapi/internal/utils"
)

// query endpoints supported by the cosmosapi service Querier
const (
    QueryApplication   = "application"
    QueryAppUsers      = "app_users"
    QueryAdminApps     = "admin_apps"
    QueryTables   = "tables"
    QueryIndex    = "index"
    QueryOption   = "option"
    QueryColumnOption   = "column_option"
    QueryRow      = "find"
    QueryIdsBy    = "find_by"
    QueryAllIds   = "find_all"
    QueryAdminGroup = "admin_group"
)


// NewQuerier is the module level router for state queries
func NewQuerier(keeper Keeper) sdk.Querier {
    return func(ctx sdk.Context, path []string, req abci.RequestQuery) (res []byte, err sdk.Error) {
        switch path[0] {
        case QueryApplication:
            if len(path) > 2 {
                return queryApplication(ctx, path[1:], req, keeper)
            } else {
                return queryApplications(ctx, path[1:], req, keeper)
            }
        case QueryAppUsers:
            return queryAppUsers(ctx, path[1:], req, keeper)
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
        default:
            return nil, sdk.ErrUnknownRequest("unknown cosmosapi query endpoint")
        }
    }
}

////////////////////////////////
//                            //
// query application/database //
//                            //
////////////////////////////////

// the the list of app code in the system
func queryApplications(ctx sdk.Context, path []string, req abci.RequestQuery, keeper Keeper) ([]byte, sdk.Error) {
    accessCode:= path[0]
    _, err := utils.VerifyAccessCode(accessCode)
    if err != nil {
        return []byte{}, sdk.ErrUnknownRequest("Access code is not valid!")
    }

    // we use the term database in the code
    applications := keeper.getAllAppCode(ctx)

    res, err := codec.MarshalJSONIndent(keeper.cdc, applications)
    if err != nil {
        panic("could not marshal result to JSON")
    }

    return res, nil
}

func queryApplication(ctx sdk.Context, path []string, req abci.RequestQuery, keeper Keeper) ([]byte, sdk.Error) {
    accessCode:= path[0]
    _, err := utils.VerifyAccessCode(accessCode)
    if err != nil {
        return []byte{}, sdk.ErrUnknownRequest("Access code is not valid!")
    }

    appCode := path[1]
    database, err := keeper.getDatabase(ctx, appCode)

    if err != nil {
        return []byte{}, sdk.ErrUnknownRequest(fmt.Sprintf("AppCode %s does not exist", appCode))
    }

    res, err := codec.MarshalJSONIndent(keeper.cdc, database)
    if err != nil {
        panic("could not marshal result to JSON")
    }

    return res, nil
}

func queryAppUsers(ctx sdk.Context, path []string, req abci.RequestQuery, keeper Keeper) ([]byte, sdk.Error) {
    accessCode:= path[0]
    addr, err := utils.VerifyAccessCode(accessCode)
    if err != nil {
        return []byte{}, sdk.ErrUnknownRequest("Access code is not valid!")
    }

    appCode := path[1]
    appId, err := keeper.GetDatabaseId(ctx, appCode)
    if err != nil {
        return []byte{}, sdk.ErrUnknownRequest("Invalid app code")
    }

    users := keeper.GetDatabaseUsers(ctx, appId, addr)

    res, err := codec.MarshalJSONIndent(keeper.cdc, users)
    if err != nil {
        panic("could not marshal result to JSON")
    }

    return res, nil
}

func queryAdminApps(ctx sdk.Context, path []string, req abci.RequestQuery, keeper Keeper) ([]byte, sdk.Error) {
    accessCode:= path[0]
    addr, err := utils.VerifyAccessCode(accessCode)
    if err != nil {
        return []byte{}, sdk.ErrUnknownRequest("Access code is not valid!")
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

func queryTables(ctx sdk.Context, path []string, req abci.RequestQuery, keeper Keeper) ([]byte, sdk.Error) {
    accessCode:= path[0]
    _, err := utils.VerifyAccessCode(accessCode)
    if err != nil {
        return []byte{}, sdk.ErrUnknownRequest("Access code is not valid!")
    }

    appId, err := keeper.GetDatabaseId(ctx, path[1])
    if err != nil {
        return []byte{}, sdk.ErrUnknownRequest("Invalid app code")
    }

    tables, err := keeper.getTables(ctx, appId)

    if err != nil {
        return nil, sdk.ErrUnknownRequest("Can not get table names")
    }

    res, err := codec.MarshalJSONIndent(keeper.cdc, tables)
    if err != nil {
        panic("could not marshal result to JSON")
    }

    return res, nil
}

func queryTable(ctx sdk.Context, path []string, req abci.RequestQuery, keeper Keeper) ([]byte, sdk.Error) {
    accessCode:= path[0]
    _, err := utils.VerifyAccessCode(accessCode)
    if err != nil {
        return []byte{}, sdk.ErrUnknownRequest("Access code is not valid!")
    }

    appId, err := keeper.GetDatabaseId(ctx, path[1])
    if err != nil {
        return nil, sdk.ErrUnknownRequest("Invalid app code")
    }

    table, err := keeper.GetTable(ctx, appId, path[2])

    if err != nil {
        return []byte{}, sdk.ErrUnknownRequest(fmt.Sprintf("Table %s does not exist",  table))
    }

    res, err := codec.MarshalJSONIndent(keeper.cdc, table)
    if err != nil {
        panic("could not marshal result to JSON")
    }

    return res, nil
}

func queryIndex(ctx sdk.Context, path []string, req abci.RequestQuery, keeper Keeper) ([]byte, sdk.Error) {
    accessCode:= path[0]
    _, err := utils.VerifyAccessCode(accessCode)
    if err != nil {
        return []byte{}, sdk.ErrUnknownRequest("Access code is not valid!")
    }

    appId, err := keeper.GetDatabaseId(ctx, path[1])
    if err != nil {
        return nil, sdk.ErrUnknownRequest("Invalid app code")
    }

    tableName := path[2]
    index, err := keeper.GetIndex(ctx, appId, tableName)

    if err != nil {
        return []byte{}, sdk.ErrUnknownRequest(fmt.Sprintf("Table %s does not exist",  tableName))
    }

    res, err := codec.MarshalJSONIndent(keeper.cdc, index)
    if err != nil {
        panic("could not marshal result to JSON")
    }

    return res, nil
}

func queryOption(ctx sdk.Context, path []string, req abci.RequestQuery, keeper Keeper) ([]byte, sdk.Error) {
    accessCode:= path[0]
    _, err := utils.VerifyAccessCode(accessCode)
    if err != nil {
        return []byte{}, sdk.ErrUnknownRequest("Access code is not valid!")
    }

    appId, err := keeper.GetDatabaseId(ctx, path[1])
    if err != nil {
        return nil, sdk.ErrUnknownRequest("Invalid app code")
    }

    tableName := path[2]
    options, err := keeper.GetOption(ctx, appId, tableName)

    if err != nil {
        return []byte{}, sdk.ErrUnknownRequest(fmt.Sprintf("Table %s does not exist",  tableName))
    }

    res, err := codec.MarshalJSONIndent(keeper.cdc, options)
    if err != nil {
        panic("could not marshal result to JSON")
    }

    return res, nil
}

func queryColumnOption(ctx sdk.Context, path []string, req abci.RequestQuery, keeper Keeper) ([]byte, sdk.Error) {
    accessCode:= path[0]
    _, err := utils.VerifyAccessCode(accessCode)
    if err != nil {
        return []byte{}, sdk.ErrUnknownRequest("Access code is not valid!")
    }

    appId, err := keeper.GetDatabaseId(ctx, path[1])
    if err != nil {
        return nil, sdk.ErrUnknownRequest("Invalid app code")
    }

    tableName := path[2]
    fieldName := path[3]

    options, err := keeper.GetColumnOption(ctx, appId, tableName, fieldName)

    if err != nil {
        return []byte{}, sdk.ErrUnknownRequest(fmt.Sprintf("Field %s.%s does not exist",  tableName, fieldName))
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

func queryRow(ctx sdk.Context, path []string, req abci.RequestQuery, keeper Keeper) ([]byte, sdk.Error) {
    accessCode:= path[0]
    addr, err := utils.VerifyAccessCode(accessCode)
    if err != nil {
        return []byte{}, sdk.ErrUnknownRequest("Access code is not valid!")
    }

    appId, err := keeper.GetDatabaseId(ctx, path[1])
    if  err != nil {
        return nil, sdk.ErrUnknownRequest("Invalid app code")
    }

    tableName := path[2] 
    u32, err := strconv.ParseUint(path[3], 10, 32)
    fields, err := keeper.Find(ctx, appId, tableName, uint(u32), addr)

    if err != nil {
        return []byte{}, sdk.ErrUnknownRequest(fmt.Sprintf("Table %s does not exist",  path[2]))
    }

    res, err := codec.MarshalJSONIndent(keeper.cdc, fields)
    if err != nil {
        panic("could not marshal result to JSON")
    }

    return res, nil
}

func queryIdsBy(ctx sdk.Context, path []string, req abci.RequestQuery, keeper Keeper) ([]byte, sdk.Error) {
    accessCode:= path[0]
    addr, err := utils.VerifyAccessCode(accessCode)
    if err != nil {
        return []byte{}, sdk.ErrUnknownRequest("Access code is not valid!")
    }

    appId, err := keeper.GetDatabaseId(ctx, path[1])
    if err != nil {
        return nil, sdk.ErrUnknownRequest("Invalid app code")
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

func queryAllIds(ctx sdk.Context, path []string, req abci.RequestQuery, keeper Keeper) ([]byte, sdk.Error) {
    accessCode:= path[0]
    addr, err := utils.VerifyAccessCode(accessCode)
    if err != nil {
        return []byte{}, sdk.ErrUnknownRequest("Access code is not valid!")
    }

    appId, err := keeper.GetDatabaseId(ctx, path[1])
    if err != nil {
        return nil, sdk.ErrUnknownRequest("Invalid app code")
    }

    var tableName = path[2]

    ids := keeper.FindAll(ctx, appId, tableName, addr)

    res, err := codec.MarshalJSONIndent(keeper.cdc, ids)
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

func queryAdminGroup(ctx sdk.Context, path []string, req abci.RequestQuery, keeper Keeper) ([]byte, sdk.Error) {
    appId, err := keeper.GetDatabaseId(ctx, path[0])
    if err != nil {
        return nil, sdk.ErrUnknownRequest("Invalid app code")
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

