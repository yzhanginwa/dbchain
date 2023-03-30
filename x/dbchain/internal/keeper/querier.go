package keeper

import (
    "encoding/json"
    "fmt"
    "github.com/cosmos/cosmos-sdk/codec"
    sdk "github.com/cosmos/cosmos-sdk/types"
    sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
    "github.com/mr-tron/base58"
    abci "github.com/tendermint/tendermint/abci/types"
    lua "github.com/yuin/gopher-lua"
    "github.com/yzhanginwa/dbchain/x/dbchain/internal/keeper/cache"
    "github.com/yzhanginwa/dbchain/x/dbchain/internal/types"
    "github.com/yzhanginwa/dbchain/x/dbchain/internal/utils"
    "github.com/yzhanginwa/dbchain/x/dbchain/internal/keeper/import_export"
    "strconv"
    "strings"
    "time"
)

// query endpoints supported by the dbchain service Querier
const (
    QueryCheckChainId  = "check_chain_id"
    QueryIsSysAdmin    = "is_sys_admin"
    QueryApplication   = "application"
    QueryApplicationBrowser   = "application_browser"
    QueryAppUsers      = "app_users"
    QueryIsAppUser     = "is_app_user"
    QueryAdminApps     = "admin_apps"
    QueryTables   = "tables"
    QueryIndex    = "index"
    QueryOption   = "option"
    QueryColumnOption   = "column_option"
    QueryColumnDataType = "column_data_type"
    QueryCanAddColumnOption = "can_add_column_option"
    QueryCanSetColumnDataType = "can_set_column_data_type"
    QueryCanInsertRow = "can_insert_row"
    QueryRow      = "find"
    QueryIdsBy    = "find_by"
    QueryAllIds   = "find_all"
    QueryMaxId    = "max_id"
    QueryGroups   = "groups"
    QueryGroup    = "group"
    QueryGroupMemo = "group_memo"
    QueryFriends  = "friends"
    QueryPendingFriends  = "pending_friends"
    QueryQuerier  = "querier"
    QueryExportDB = "export_database"
    QueryFunctions = "functions"
    QueryFunctionInfo = "functionInfo"
    QueryCustomQueriers  = "customQueriers"
    QueryCustomQuerierInfo = "customQuerierInfo"
    QueryCallCustomQuerier = "callCustomQuerier"
    QueryTxSimpleResult    = "txSimpleResult"
    QueryAllAccounts       = "allAccounts"
    QueryDbchainTxNum      = "dbchainTxNum"
    QueryDbchainRecentTxNum  = "dbchainRecentTxNum"
    QueryApplicationUserFileVolumeLimit  = "application_user_file_volume_limit"
    QueryApplicationUserUsedFileVolume  = "application_user_used_file_volume"
)


// NewQuerier is the module level router for state queries
func NewQuerier(keeper Keeper) sdk.Querier {
    return func(ctx sdk.Context, path []string, req abci.RequestQuery) (res []byte, err error) {
        switch path[0] {
        case QueryCheckChainId:
            return queryCheckChainId(ctx, path[1:], req, keeper)
        case QueryIsSysAdmin:
            return queryIsSysAdmin(ctx, path[1:], req, keeper)
        case QueryApplication:
            if len(path) > 2 {
                return queryApplication(ctx, path[1:], req, keeper)
            } else {
                return queryApplications(ctx, path[1:], req, keeper)
            }
        case QueryApplicationUserFileVolumeLimit :
            return queryApplicationUserFileVolumeLimit(ctx, path[1:], req, keeper)
        case QueryApplicationUserUsedFileVolume :
            return queryApplicationUserUsedFileVolume(ctx, path[1:], req, keeper)
        case QueryApplicationBrowser:
            return queryApplicationsBrowser(ctx, path[1:], req, keeper)
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
        case QueryColumnDataType:
            return queryColumnDataType(ctx, path[1:], req, keeper)
        case QueryCanAddColumnOption:
            return queryCanAddColumnOption(ctx, path[1:], req, keeper)
        case QueryCanSetColumnDataType:
            return queryCanAddColumnDataType(ctx, path[1:], req, keeper)
        case QueryCanInsertRow:
            return queryCanInsertRow(ctx, path[1:], req, keeper)
        case QueryRow:
            return queryRow(ctx, path[1:], req, keeper)
        case QueryIdsBy:
            return queryIdsBy(ctx, path[1:], req, keeper)
        case QueryAllIds:
            return queryAllIds(ctx, path[1:], req, keeper)
        case QueryMaxId:
            return queryMaxId(ctx, path[1:], req, keeper)
        case QueryGroups:
            return queryGroups(ctx, path[1:], req, keeper)
        case QueryGroup:
            return queryGroup(ctx, path[1:], req, keeper)
        case QueryGroupMemo:
            return queryGroupMemo(ctx, path[1:], req, keeper)
        case QueryFriends:
            return queryFriends(ctx, path[1:], req, keeper)
        case QueryPendingFriends:
            return queryPendingFriends(ctx, path[1:], req, keeper)
        case QueryQuerier:
            return queryQuerier(ctx, path[1:], req, keeper)
        case QueryExportDB:
            return queryExportDatabase(ctx, path[1:], req, keeper)
        case QueryFunctions:
            return queryFunctions(ctx, path[1:], req, keeper)
        case QueryFunctionInfo:
            return queryFunctionsInfo(ctx, path[1:], req, keeper)
        case QueryCustomQueriers:
            return queryCustomQueriers(ctx, path[1:], req, keeper)
        case QueryCustomQuerierInfo:
            return queryCustomQuerierInfo(ctx, path[1:], req, keeper)
        case QueryCallCustomQuerier:
            return queryCallCustomQuerier(ctx, path[1:], req, keeper)
        case QueryTxSimpleResult:
            return queryTxSimpleResult(ctx, path[1:], req, keeper)
        case QueryAllAccounts:
            return queryAllAccount(ctx, path[1:], req, keeper)
        case QueryDbchainTxNum:
            return queryDbchainTxNum(ctx, path[1:], req, keeper)
        case QueryDbchainRecentTxNum:
            return queryDbchainRecentTxNum(ctx, path[1:], req, keeper)
        default:
            return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "unknown dbchain query endpoint")
        }
    }
}

/////////////////////////////////////////////////
//                                             //
// query whether the given chain Id is correct //
//                                             //
/////////////////////////////////////////////////

func queryCheckChainId(ctx sdk.Context, path []string, req abci.RequestQuery, keeper Keeper) ([]byte, error) {
    if len(path) != 2 {
        return []byte{}, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "Number of query parameters is wrong!")
    }
    accessCode:= path[0]
    _, err := utils.VerifyAccessCode(accessCode)
    if err != nil {
        return []byte{}, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "Access code is not valid!")
    }

    testChainId := path[1]
    chainId := ctx.ChainID()
    result := (testChainId == chainId)

    res, err := codec.MarshalJSONIndent(keeper.cdc, result)
    if err != nil {
        panic("could not marshal result to JSON")
    }

    return res, nil
}

////////////////////////////////////
//                                //
// query whether user is sysadmin //
//                                //
////////////////////////////////////

func queryIsSysAdmin(ctx sdk.Context, path []string, req abci.RequestQuery, keeper Keeper) ([]byte, error) {
    accessCode:= path[0]
    addr, err := utils.VerifyAccessCode(accessCode)
    if err != nil {
        return []byte{}, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "Access code is not valid!")
    }

    isSysAdmin:= keeper.IsSysAdmin(ctx, addr)

    res, err := codec.MarshalJSONIndent(keeper.cdc, isSysAdmin)
    if err != nil {
        panic("could not marshal result to JSON")
    }

    return res, nil
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
    //if database expiration delete and return null
    if database.Deleted == true && database.Expiration <= time.Now().Unix(){
        keeper.PurgeApplication(ctx, appCode)
        return []byte{}, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, fmt.Sprintf("AppCode %s does not exist", appCode))
    }

    res, err := codec.MarshalJSONIndent(keeper.cdc, database)
    if err != nil {
        panic("could not marshal result to JSON")
    }

    return res, nil
}

func queryApplicationUserFileVolumeLimit(ctx sdk.Context, path []string, req abci.RequestQuery, keeper Keeper) ([]byte, error) {

    accessCode:= path[0]
    _, _, err := utils.VerifyAccessCodeWithoutTimeChecking(accessCode)
    if err != nil {
        return []byte{}, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "Access code is not valid!")
    }

    appCode := path[1]
    appId, err := keeper.GetDatabaseId(ctx, appCode)
    if err != nil {
        return []byte{}, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "Invalid app code")
    }

    size := keeper.GetApplicationUserFileVolumeLimit(ctx, appId)
    res, err := codec.MarshalJSONIndent(keeper.cdc, size)
    if err != nil {
        panic("could not marshal result to JSON")
    }

    return res, nil
}

func queryApplicationUserUsedFileVolume(ctx sdk.Context, path []string, req abci.RequestQuery, keeper Keeper) ([]byte, error) {

    accessCode:= path[0]
    addr, _, err := utils.VerifyAccessCodeWithoutTimeChecking(accessCode)
    if err != nil {
        return []byte{}, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "Access code is not valid!")
    }

    appCode := path[1]
    appId, err := keeper.GetDatabaseId(ctx, appCode)
    if err != nil {
        return []byte{}, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "Invalid app code")
    }

    size := keeper.GetApplicationUserUsedFileVolume(ctx, appId, addr)
    res, err := codec.MarshalJSONIndent(keeper.cdc, size)
    if err != nil {
        panic("could not marshal result to JSON")
    }

    return res, nil
}

func queryApplicationsBrowser(ctx sdk.Context, path []string, req abci.RequestQuery, keeper Keeper) ([]byte, error) {

    applications := keeper.GetAllAppCode(ctx)
    res, err := codec.MarshalJSONIndent(keeper.cdc, len(applications))
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

    if !keeper.isAdmin(ctx, appId, addr) {
        return []byte{}, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "Admin privilege is needed!")
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

    tables := keeper.GetTables(ctx, appId)
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
    index, err := keeper.GetIndexFields(ctx, appId, tableName)

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

func queryColumnDataType(ctx sdk.Context, path []string, req abci.RequestQuery, keeper Keeper) ([]byte, error) {
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

    dataType, err := keeper.GetColumnDataType(ctx, appId, tableName, fieldName)

    if err != nil {
        return []byte{}, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, fmt.Sprintf("Field %s.%s does not exist",  tableName, fieldName))
    }

    res, err := codec.MarshalJSONIndent(keeper.cdc, dataType)
    if err != nil {
        panic("could not marshal result to JSON")
    }

    return res, nil
}

func queryCanAddColumnOption(ctx sdk.Context, path []string, req abci.RequestQuery, keeper Keeper) ([]byte, error) {
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
    option    := path[4]

    result := keeper.GetCanAddColumnOption(ctx, appId, tableName, fieldName, option)

    res, err := codec.MarshalJSONIndent(keeper.cdc, result)
    if err != nil {
        panic("could not marshal result to JSON")
    }

    return res, nil
}

func queryCanAddColumnDataType(ctx sdk.Context, path []string, req abci.RequestQuery, keeper Keeper) ([]byte, error) {
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
    dataType  := path[4]

    var result bool
    result = keeper.GetCanSetColumnDataType(ctx, appId, tableName, fieldName, dataType)

    res, err := codec.MarshalJSONIndent(keeper.cdc, result)
    if err != nil {
        panic("could not marshal result to JSON")
    }

    return res, nil
}

 func queryCanInsertRow(ctx sdk.Context, path []string, req abci.RequestQuery, keeper Keeper) ([]byte, error) {
    accessCode:= path[0]
    addr, err := utils.VerifyAccessCode(accessCode)
    if err != nil {
        return []byte{}, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "Access code is not valid!")
    }

    appId, err := keeper.GetDatabaseId(ctx, path[1])
    if err != nil {
        return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "Invalid app code")
    }

    tableName     := path[2]
    encodeFields  := path[3]

     rowFieldsJson, err := base58.Decode(encodeFields)
     if err != nil {
         return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest,"Failed to parse row fields!")
     }

    var rowFields types.RowFields
    if err := json.Unmarshal([]byte(rowFieldsJson), &rowFields); err != nil {
        return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest,"Failed to parse row fields!")
    }

    result := true
     L := lua.NewState(lua.Options{
         SkipOpenLibs : true,
     })
     defer L.Close()
     L.SetGlobal("IsRegisterData",lua.LBool(false))
    _, err = keeper.PreInsertCheck(ctx, appId, tableName, rowFields, addr, L)
    if err != nil {
        result = false
    }

    res, err := codec.MarshalJSONIndent(keeper.cdc, result)
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
    ids := keeper.FindBy(ctx, appId, tableName, fieldName, []string{value}, addr)

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

func queryMaxId(ctx sdk.Context, path []string, req abci.RequestQuery, keeper Keeper) ([]byte, error) {
    accessCode:= path[0]
    _, err := utils.VerifyAccessCode(accessCode)
    if err != nil {
        return []byte{}, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "Access code is not valid!")
    }

    appId, err := keeper.GetDatabaseId(ctx, path[1])
    if err != nil {
        return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "Invalid app code")
    }

    var tableName = path[2]
    nextId, _ := keeper.PeekNextId(ctx, appId, tableName)

    var maxId int
    if nextId > 0 {
        maxId = int(nextId) - 1
    } else {
        maxId = 0
    }

    res, err := codec.MarshalJSONIndent(keeper.cdc, maxId)
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

func queryGroups(ctx sdk.Context, path []string, req abci.RequestQuery, keeper Keeper) ([]byte, error) {
    accessCode:= path[0]
    _, err := utils.VerifyAccessCode(accessCode)
    if err != nil {
        return []byte{}, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "Access code is not valid!")
    }

    appId, err := keeper.GetDatabaseId(ctx, path[1])
    if err != nil {
        return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "Invalid app code")
    }

    groups := keeper.getGroups(ctx, appId)

    res, err := codec.MarshalJSONIndent(keeper.cdc, groups)
    if err != nil {
        panic("could not marshal result to JSON")
    }

    return res, nil
}

func queryGroup(ctx sdk.Context, path []string, req abci.RequestQuery, keeper Keeper) ([]byte, error) {
    accessCode:= path[0]
    _, err := utils.VerifyAccessCode(accessCode)
    if err != nil {
        return []byte{}, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "Access code is not valid!")
    }

    appId, err := keeper.GetDatabaseId(ctx, path[1])
    if err != nil {
        return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "Invalid app code")
    }

    groupName := path[2]
    addresses := keeper.getGroupMembers(ctx, appId, groupName)

    res, err := codec.MarshalJSONIndent(keeper.cdc, addresses)
    if err != nil {
        panic("could not marshal result to JSON")
    }

    return res, nil
}

func queryGroupMemo(ctx sdk.Context, path []string, req abci.RequestQuery, keeper Keeper) ([]byte, error) {
    accessCode:= path[0]
    _, err := utils.VerifyAccessCode(accessCode)
    if err != nil {
        return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "Access code is not valid!")
    }

    appId, err := keeper.GetDatabaseId(ctx, path[1])
    if err != nil {
        return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "Invalid app code")
    }

    groupName := path[2]
    memo := keeper.getGroupMembersMemo(ctx, appId, groupName)

    res, err := codec.MarshalJSONIndent(keeper.cdc, memo)
    if err != nil {
        panic("could not marshal result to JSON")
    }

    return res, nil
}

/////////////////////
//                 //
// export database //
//                 //
/////////////////////

func queryExportDatabase (ctx sdk.Context, path []string, req abci.RequestQuery, keeper Keeper) ([]byte, error) {
    appCode := path[0]
    appId, err := keeper.GetDatabaseId(ctx, appCode)
    if err != nil {
        return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "Invalid app code")
    }
    
    ieDatabase := import_export.Database{}
    ieDatabase.Appcode = appCode
    database, err := keeper.getDatabaseById(ctx, appId)
    if err != nil {
        return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "Failed to get database by id")
    }
    ieDatabase.Name = database.Name
    ieDatabase.Memo = database.Description

    ieTables := []import_export.Table{}

    tableNames := keeper.GetTables(ctx, appId)
    for _, tableName := range tableNames {
        ieTable := import_export.Table{}
        ieTable.Name = tableName
     
        tableOptions, err := keeper.GetOption(ctx, appId, tableName)
        if err != nil {
            return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "Failed to get table options")
        }

        ieTable.Options = tableOptions

        tableObj, err := keeper.GetTable(ctx, appId, tableName)
        if err != nil {
            return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "Failed to get table columns")
        }

        ieTable.Filter = tableObj.Filter
        ieTable.Trigger = tableObj.Trigger
        ieTable.Memo = tableObj.Memo
        ieFields := []import_export.Field{}

        // handle fields
        for index, fieldName := range tableObj.Fields {
            if fieldName == "id" || fieldName == "created_by" || fieldName == "created_at" || fieldName == "tx_hash" {
                continue
            }
            ieField := import_export.Field{}
            ieField.Name = fieldName

            fieldOptions, err := keeper.GetColumnOption(ctx, appId, tableName, fieldName)
            if err != nil {
                return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "Failed to get table column options")
            }
            ieField.PropertyArr = fieldOptions

            fieldDataType, err := keeper.GetColumnDataType(ctx, appId, tableName, fieldName)
            if err != nil {
                return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "Failed to get table column data type")
            }
            ieField.FieldType = fieldDataType

            if len(tableObj.Memos) > index {
                ieField.Memo = tableObj.Memos[index]
            }

            ieFields = append(ieFields, ieField)
        }
        ieTable.Fields = ieFields
        ieTables = append(ieTables, ieTable)
    }
    ieDatabase.Tables = ieTables

    ieFuncs := _generateImportExportFuncs(ctx, keeper, appId, 0)
    ieDatabase.CustomFns = ieFuncs

    ieQueriers := _generateImportExportFuncs(ctx, keeper, appId, 1)
    ieDatabase.CustomQueriers = ieQueriers

    res, err := codec.MarshalJSONIndent(keeper.cdc, ieDatabase)
    if err != nil {
        panic("could not marshal result to JSON")
    }

    return res, nil
}

func _generateImportExportFuncs(ctx sdk.Context, keeper Keeper, appId uint, funcOrQuerier int) []import_export.CustomFn {
    funcNames := keeper.GetFunctions(ctx, appId, funcOrQuerier)
    ieFuncs := []import_export.CustomFn{}
    for _, funcName := range funcNames {
        functionInfo := keeper.GetFunctionInfo(ctx, appId, funcName, funcOrQuerier)
        ieFunc := import_export.CustomFn{}
        ieFunc.Name = functionInfo.Name
        ieFunc.Owner = functionInfo.Owner.String()
        ieFunc.Description = functionInfo.Description
        ieFunc.Body = functionInfo.Body
        ieFuncs = append(ieFuncs, ieFunc)
    }
    return ieFuncs
}

/////////////////////
//                 //
// query functions //
//                 //
/////////////////////

func queryFunctions(ctx sdk.Context, path []string, req abci.RequestQuery, keeper Keeper) ([]byte, error) {
    accessCode:= path[0]
    _, err := utils.VerifyAccessCode(accessCode)
    if err != nil {
        return []byte{}, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "Access code is not valid!")
    }

    appId, err := keeper.GetDatabaseId(ctx, path[1])
    if err != nil {
        return []byte{}, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "Invalid app code")
    }

    functions := keeper.GetFunctions(ctx, appId, FuncHandleType)
    res, err := codec.MarshalJSONIndent(keeper.cdc, functions)
    if err != nil {
        panic("could not marshal result to JSON")
    }

    return res, nil
}


func queryFunctionsInfo(ctx sdk.Context, path []string, req abci.RequestQuery, keeper Keeper) ([]byte, error) {
    accessCode:= path[0]
    _, err := utils.VerifyAccessCode(accessCode)
    if err != nil {
        return []byte{}, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "Access code is not valid!")
    }

    appId, err := keeper.GetDatabaseId(ctx, path[1])
    if err != nil {
        return []byte{}, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "Invalid app code")
    }

    functionInfo := keeper.GetFunctionInfo(ctx, appId, path[2], 0)
    res, err := codec.MarshalJSONIndent(keeper.cdc, functionInfo)
    if err != nil {
        panic("could not marshal result to JSON")
    }

    return res, nil
}

func queryTxSimpleResult(ctx sdk.Context, path []string, req abci.RequestQuery, keeper Keeper) ([]byte, error) {
    accessCode:= path[0]
    _, err := utils.VerifyAccessCode(accessCode)
    if err != nil {
        return []byte{}, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "Access code is not valid!")
    }

    var txState *types.TxStatus
    nowTime := time.Now().Unix()
    txHash := path[1]
    txHash = strings.ToLower(txHash)
    txStateIt,ok := cache.TxStatusCache.Load(txHash)
    if !ok {
        errStr := "Tx : " + path[1] + " is unhandled" + ". Please check again later !"
        txState = types.NewTxStatus(cache.TxStatePending, 0, errStr, nowTime)
    } else {
        txState = txStateIt.(*types.TxStatus)
        //The information has expired and needs to be deleted
        if nowTime - txState.GetTimeStamp() > cache.TxStateInvalidTime {
            cache.TxStatusCache.Delete(txHash)
        }
    }

    res, err := codec.MarshalJSONIndent(keeper.cdc, txState)
    if err != nil {
        panic("could not marshal result to JSON")
    }

    return res, nil
}

/////////////////////
//                 //
// query queriers //
//                 //
/////////////////////

func queryCustomQueriers(ctx sdk.Context, path []string, req abci.RequestQuery, keeper Keeper) ([]byte, error) {
    accessCode:= path[0]
    _, err := utils.VerifyAccessCode(accessCode)
    if err != nil {
        return []byte{}, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "Access code is not valid!")
    }

    appId, err := keeper.GetDatabaseId(ctx, path[1])
    if err != nil {
        return []byte{}, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "Invalid app code")
    }

    functions := keeper.GetFunctions(ctx, appId, QueryHandleType)
    res, err := codec.MarshalJSONIndent(keeper.cdc, functions)
    if err != nil {
        panic("could not marshal result to JSON")
    }

    return res, nil
}


func queryCustomQuerierInfo(ctx sdk.Context, path []string, req abci.RequestQuery, keeper Keeper) ([]byte, error) {
    accessCode:= path[0]
    _, err := utils.VerifyAccessCode(accessCode)
    if err != nil {
        return []byte{}, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "Access code is not valid!")
    }

    appId, err := keeper.GetDatabaseId(ctx, path[1])
    if err != nil {
        return []byte{}, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "Invalid app code")
    }

    querierInfo := keeper.GetFunctionInfo(ctx, appId, path[2], 1)
    res, err := codec.MarshalJSONIndent(keeper.cdc, querierInfo)
    if err != nil {
        panic("could not marshal result to JSON")
    }

    return res, nil
}

func queryCallCustomQuerier(ctx sdk.Context, path []string, req abci.RequestQuery, keeper Keeper) ([]byte, error) {
    accessCode:= path[0]
    addr, err := utils.VerifyAccessCode(accessCode)
    if err != nil {
        return []byte{}, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "Access code is not valid!")
    }

    appId, err := keeper.GetDatabaseId(ctx, path[1])
    if err != nil {
        return []byte{}, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "Invalid app code")
    }

    querierInfo := keeper.GetFunctionInfo(ctx, appId, path[2], 1)
    res , err := keeper.DoCustomQuerier(ctx, appId, querierInfo, path[3], addr)
    //res, err := codec.MarshalJSONIndent(keeper.cdc, functionInfo)
    if err != nil {
        panic("could not marshal result to JSON")
    }

    return res, nil
}

func queryAllAccount(ctx sdk.Context, path []string, req abci.RequestQuery, keeper Keeper)([]byte, error){
   accounts := keeper.GetAllAccounts(ctx)
   accountNum := len(accounts)
   res, err := json.Marshal(accountNum)
   if err != nil {
       panic("could not marshal result to JSON")
   }
   return res, nil
}

func queryDbchainTxNum(ctx sdk.Context, path []string, req abci.RequestQuery, keeper Keeper)([]byte, error){
    res, err := keeper.GetDbchainTxNum(ctx)
    if err != nil {
        panic("could not marshal result to JSON")
    }
    return res, nil
}

func queryDbchainRecentTxNum(ctx sdk.Context, path []string, req abci.RequestQuery, keeper Keeper)([]byte, error){
    res, err := keeper.GetDbchainRecentTxNum(ctx)
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

