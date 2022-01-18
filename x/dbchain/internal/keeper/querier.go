package keeper

import (
    "encoding/json"
    "fmt"
    "github.com/dbchaincloud/cosmos-sdk/codec"
    sdk "github.com/dbchaincloud/cosmos-sdk/types"
    sdkerrors "github.com/dbchaincloud/cosmos-sdk/types/errors"
    "github.com/mr-tron/base58"
    abci "github.com/dbchaincloud/tendermint/abci/types"
    lua "github.com/yuin/gopher-lua"
    "github.com/yzhanginwa/dbchain/x/dbchain/internal/keeper/cache"
    "github.com/yzhanginwa/dbchain/x/dbchain/internal/types"
    "github.com/yzhanginwa/dbchain/x/dbchain/internal/utils"
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
    QueryOptionForInternal   = "table_option_for_internal"     // used for rest server to check if a table public
    QueryAssociation = "association"
    QueryCounterCache = "counter_cache"
    QueryColumnOption   = "column_option"
    QueryCounterInfo    = "counter_info"
    QueryColumnDataType = "column_data_type"
    QueryCanAddColumnOption = "can_add_column_option"
    QueryCanSetColumnDataType = "can_set_column_data_type"
    QueryCanInsertRow = "can_insert_row"
    QueryRow      = "find"
    QueryIdsBy    = "find_by"
    QueryAllIds   = "find_all"
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
    QueryByDynamicScript  = "dynamic_script"
    QueryTxSimpleResult    = "txSimpleResult"
    QueryAllAccounts       = "allAccounts"
    QueryDbchainTxNum      = "dbchainTxNum"
    QueryDbchainRecentTxNum  = "dbchainRecentTxNum"
    QueryApplicationUserFileVolumeLimit  = "application_user_file_volume_limit"
    QueryApplicationUserUsedFileVolume  = "application_user_used_file_volume"
    //add for bsb
    QueryAccountTxs = "account_txs"
    QueryAccountTxsByTime  = "account_txs_by_time"
    QueryTokenKeepers = "token_keepers"
    QueryLimitP2PTransferStatus = "limit_p2p_transfer_status"
    QueryUserPrivateKey = "get_user_private_key"
    QueryCurrentMinGasPrices = "current_min_gas_prices"

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
        case QueryOptionForInternal:
            return queryTableOptionForInternal(ctx, path[1:], req, keeper)
        case QueryAssociation:
            return queryAssociation(ctx, path[1:], req, keeper)
        case QueryCounterCache:
            return queryCounterCache(ctx, path[1:], req, keeper)
        case QueryColumnOption:
            return queryColumnOption(ctx, path[1:], req, keeper)
        case QueryCounterInfo:
            return queryCounterInfo(ctx, path[1:], req, keeper)
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
        case QueryByDynamicScript:
            return queryByDynamicScript(ctx, path[1:], req, keeper)
        case QueryTxSimpleResult:
            return queryTxSimpleResult(ctx, path[1:], req, keeper)
        case QueryAllAccounts:
            return queryAllAccount(ctx, path[1:], req, keeper)
        case QueryDbchainTxNum:
            return queryDbchainTxNum(ctx, path[1:], req, keeper)
        case QueryDbchainRecentTxNum:
            return queryDbchainRecentTxNum(ctx, path[1:], req, keeper)
        case QueryAccountTxs:
            return queryAccountTxs(ctx, path[1:], req, keeper)
        case QueryAccountTxsByTime:
            return queryAccountTxsByTime(ctx, path[1:], req, keeper)
        case QueryTokenKeepers:
            return queryTokenKeepers(ctx, path[1:], req, keeper)
        case QueryLimitP2PTransferStatus:
            return queryLimitP2PTransferStatus(ctx, path[1:], req, keeper)
        case QueryUserPrivateKey:
            return queryUserPrivateKey(ctx, path[1:], req, keeper)
        case QueryCurrentMinGasPrices:
            return queryCurrentMinGasPrices(ctx, path[1:], req, keeper)
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
    addr, err := utils.VerifyAccessCode(accessCode)
    if err != nil {
        return []byte{}, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "Access code is not valid!")
    }

    appId, err := keeper.GetDatabaseId(ctx, path[1])
    if err != nil {
        return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "Invalid app code")
    }

    if !keeper.isAdmin(ctx, appId, addr) {
        return []byte{}, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "Admin privilege is needed!")
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

func queryTableOptionForInternal(ctx sdk.Context, path []string, req abci.RequestQuery, keeper Keeper) ([]byte, error) {
    appId, err := keeper.GetDatabaseId(ctx, path[0])
    if err != nil {
        return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "Invalid app code")
    }

    tableName := path[1]
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

func queryAssociation(ctx sdk.Context, path []string, req abci.RequestQuery, keeper Keeper) ([]byte, error) {
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
    associations := keeper.GetTableAssociations(ctx, appId, tableName)

    if associations == nil {
        return []byte{}, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, fmt.Sprintf("Table %s does not have association",  tableName))
    }

    res, err := codec.MarshalJSONIndent(keeper.cdc, associations)
    if err != nil {
        panic("could not marshal result to JSON")
    }

    return res, nil
}

func queryCounterCache(ctx sdk.Context, path []string, req abci.RequestQuery, keeper Keeper) ([]byte, error) {
    accessCode:= path[0]
    addr, err := utils.VerifyAccessCode(accessCode)
    if err != nil {
        return []byte{}, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "Access code is not valid!")
    }

    appId, err := keeper.GetDatabaseId(ctx, path[1])
    if err != nil {
        return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "Invalid app code")
    }

    if !keeper.isAdmin(ctx, appId, addr) {
        return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "permission forbidden")
    }

    tableName := path[2]
    counterCache := keeper.GetCounterCache(ctx, appId, tableName)

    if counterCache == nil {
        return []byte{}, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, fmt.Sprintf("Table %s does not have association",  tableName))
    }

    res, err := codec.MarshalJSONIndent(keeper.cdc, counterCache)
    if err != nil {
        panic("could not marshal result to JSON")
    }

    return res, nil
}

func queryColumnOption(ctx sdk.Context, path []string, req abci.RequestQuery, keeper Keeper) ([]byte, error) {
    accessCode:= path[0]
    addr, err := utils.VerifyAccessCode(accessCode)
    if err != nil {
        return []byte{}, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "Access code is not valid!")
    }

    appId, err := keeper.GetDatabaseId(ctx, path[1])
    if err != nil {
        return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "Invalid app code")
    }

    if !keeper.isAdmin(ctx, appId, addr) {
        return []byte{}, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "Admin privilege is needed!")
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

func queryCounterInfo(ctx sdk.Context, path []string, req abci.RequestQuery, keeper Keeper) ([]byte, error) {
    accessCode:= path[0]
    addr, err := utils.VerifyAccessCode(accessCode)
    if err != nil {
        return []byte{}, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "Access code is not valid!")
    }

    appId, err := keeper.GetDatabaseId(ctx, path[1])
    if err != nil {
        return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "Invalid app code")
    }

    if !keeper.isAdmin(ctx, appId, addr) {
        return []byte{}, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "Admin privilege is needed!")
    }

    tableName := path[2]

    counterInfo, err := keeper.GetCounterInfo(ctx, appId, tableName)

    if err != nil {
        return []byte{}, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "find counter info err")
    }

    res, err := codec.MarshalJSONIndent(keeper.cdc, counterInfo)
    if err != nil {
        panic("could not marshal result to JSON")
    }

    return res, nil
}

func queryColumnDataType(ctx sdk.Context, path []string, req abci.RequestQuery, keeper Keeper) ([]byte, error) {
    accessCode:= path[0]
    addr, err := utils.VerifyAccessCode(accessCode)
    if err != nil {
        return []byte{}, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "Access code is not valid!")
    }

    appId, err := keeper.GetDatabaseId(ctx, path[1])
    if err != nil {
        return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "Invalid app code")
    }

    if !keeper.isAdmin(ctx, appId, addr) {
        return []byte{}, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "Admin privilege is needed!")
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
     openBase(L)
     registerTableType(L, ctx, appId, keeper, addr)
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
    result := []string{}

    appId, err := keeper.GetDatabaseId(ctx, appCode)
    if err != nil {
        return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "Invalid app code")
    }

    // handle groups
    groups := keeper.getGroups(ctx, appId)
    result = append(result, "Groups:")
    for _, group := range groups {
        result = append(result, fmt.Sprintf("\t%s", group))
        groupMembers := keeper.getGroupMembers(ctx, appId, group)
        for _, groupMember := range groupMembers {
            result = append(result, fmt.Sprintf("\t\t%s", groupMember.String()))
        }
    }

    // handle tables
    tables := keeper.GetTables(ctx, appId)
    result = append(result, "Tables:")
    for _, table := range tables {
        tableOptions, err := keeper.GetOption(ctx, appId, table)
        if err != nil {
            return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "Failed to get table options")
        }

        result = append(result, "")
        if len(tableOptions) > 0 {
            result = append(result, fmt.Sprintf("\t%s (%s)", table, strings.Join(tableOptions, ", ")))
        } else {
            result = append(result, fmt.Sprintf("\t%s", table))
        }

        tableObj, err := keeper.GetTable(ctx, appId, table)
        if err != nil {
            return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "Failed to get table columns")
        }

        if len(tableObj.Filter) > 0 {
            result = append(result, fmt.Sprintf("\t%s", "Filter:"))
            result = append(result, fmt.Sprintf("\t\t%s", tableObj.Filter))
        }

        if len(tableObj.Trigger) > 0 {
            result = append(result, fmt.Sprintf("\t%s", "Trigger:"))
            result = append(result, fmt.Sprintf("\t\t%s", tableObj.Trigger))
        }

        // handle fields
        for _, field := range tableObj.Fields {
            fieldOptions, err := keeper.GetColumnOption(ctx, appId, table, field)
            if err != nil {
                return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "Failed to get table columns")
            }

            if len(fieldOptions) > 0 {
                result = append(result, fmt.Sprintf("\t\t%s (%s)", field, strings.Join(fieldOptions, ", ")))
            } else {
                result = append(result, fmt.Sprintf("\t\t%s", field))
            }
        }
    }

    res, err := codec.MarshalJSONIndent(keeper.cdc, result)
    if err != nil {
        panic("could not marshal result to JSON")
    }

    return res, nil
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
    addr, err := utils.VerifyAccessCode(accessCode)
    if err != nil {
        return []byte{}, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "Access code is not valid!")
    }

    appId, err := keeper.GetDatabaseId(ctx, path[1])
    if err != nil {
        return []byte{}, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "Invalid app code")
    }

    if !keeper.isAdmin(ctx, appId, addr) {
        return []byte{}, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "Admin privilege is needed!")
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
    addr, err := utils.VerifyAccessCode(accessCode)
    if err != nil {
        return []byte{}, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "Access code is not valid!")
    }

    appId, err := keeper.GetDatabaseId(ctx, path[1])
    if err != nil {
        return []byte{}, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "Invalid app code")
    }

    if !keeper.isAdmin(ctx, appId, addr) {
        return []byte{}, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "Admin privilege is needed!")
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
        return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest,fmt.Sprintf("%v", err))
    }

    return res, nil
}

func queryByDynamicScript(ctx sdk.Context, path []string, req abci.RequestQuery, keeper Keeper) ([]byte, error) {
    accessCode:= path[0]
    addr, err := utils.VerifyAccessCode(accessCode)
    if err != nil {
        return []byte{}, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "Access code is not valid!")
    }

    appId, err := keeper.GetDatabaseId(ctx, path[1])
    encode := path[2]
    script, err := base58.Decode(encode)
    if err != nil {
        return []byte{}, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "Invalid query data")
    }

    res , err := keeper.DoDynamicScript(ctx, appId, string(script), addr)
    if err != nil {
        return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest,fmt.Sprintf("%v", err))
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

func queryAccountTxs(ctx sdk.Context, path []string, req abci.RequestQuery, keeper Keeper)([]byte, error) {
    accessCode := path[0]
    addr, err := utils.VerifyAccessCode(accessCode)
    if err != nil {
        return []byte{}, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "Access code is not valid!")
    }
    num := 0
    if len(path) > 1 {
        num , err = strconv.Atoi(path[1])
        if err != nil {
            num = 0
        }
    }
    txs := keeper.GetAddrTxs(ctx, addr, uint(num))
    res, err := json.Marshal(txs)
    if err != nil {
        panic("could not marshal result to JSON")
    }

    return res, nil
}

func queryAccountTxsByTime(ctx sdk.Context, path []string, req abci.RequestQuery, keeper Keeper)([]byte, error) {
    userAccountAddress := path[0]
    startDate := path[1]
    endDate := path[2]

    addr , err := sdk.AccAddressFromBech32(userAccountAddress)
    if err != nil {
        return []byte{}, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "userAccountAddress error")
    }

    txs := keeper.GetAddrTxsByTime(ctx, addr, startDate, endDate)
    res, err := json.Marshal(txs)
    if err != nil {
        panic("could not marshal result to JSON")
    }

    return res, nil
}

func queryTokenKeepers(ctx sdk.Context, path []string, req abci.RequestQuery, keeper Keeper)([]byte, error) {
    accessCode := path[0]
    addr, err := utils.VerifyAccessCode(accessCode)
    if err != nil {
        return []byte{}, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "Access code is not valid!")
    }
    admins := keeper.ShowTokenKeepers(ctx, addr)
    res, err := codec.MarshalJSONIndent(keeper.cdc, admins)
    if err != nil {
        panic("could not marshal result to JSON")
    }

    return res, nil
}
//queryLimitP2PTransferStatus
func queryLimitP2PTransferStatus(ctx sdk.Context, path []string, req abci.RequestQuery, keeper Keeper)([]byte, error) {
    accessCode := path[0]
    _, err := utils.VerifyAccessCode(accessCode)
    if err != nil {
        return []byte{}, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "Access code is not valid!")
    }
    limit := keeper.ShowCurrentLimitP2PTransferStatus(ctx)
    res, err := codec.MarshalJSONIndent(keeper.cdc, limit)
    if err != nil {
        panic("could not marshal result to JSON")
    }

    return res, nil
}

func queryUserPrivateKey(ctx sdk.Context, path []string, req abci.RequestQuery, keeper Keeper)([]byte, error) {
    addr := path[0]
    limit := keeper.GetUserPrivateInfo(ctx, addr)
    return limit, nil
}

func queryCurrentMinGasPrices(ctx sdk.Context, path []string, req abci.RequestQuery, keeper Keeper)([]byte, error) {
    minGasPrices := ctx.MinGasPrices()
    res, err := codec.MarshalJSONIndent(keeper.cdc, minGasPrices)
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

