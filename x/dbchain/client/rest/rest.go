package rest

import (
    "fmt"
    "github.com/cosmos/cosmos-sdk/client/context"
    "github.com/gorilla/mux"
)

// RegisterRoutes - Central function to define routes that get registered by the main application
func RegisterRoutes(cliCtx context.CLIContext, r *mux.Router, storeName string) {
    r.HandleFunc(fmt.Sprintf("/%s/check_chain_id/{%s}/{%s}", storeName, "accessToken", "chainId"), showCheckChainIdHandler(cliCtx, storeName)).Methods("GET")
    r.HandleFunc(fmt.Sprintf("/%s/is_sys_admin/{%s}", storeName, "accessToken"), showIsSysAdminHandler(cliCtx, storeName)).Methods("GET")
    r.HandleFunc(fmt.Sprintf("/%s/application/{%s}", storeName, "accessToken"), showApplicationsHandler(cliCtx, storeName)).Methods("GET")
    r.HandleFunc(fmt.Sprintf("/%s/application/{%s}/{%s}", storeName, "accessToken", "appCode"), showApplicationHandler(cliCtx, storeName)).Methods("GET")
    r.HandleFunc(fmt.Sprintf("/%s/admin_apps/{%s}", storeName, "accessToken"), showAdminAppsHandler(cliCtx, storeName)).Methods("GET")
    r.HandleFunc(fmt.Sprintf("/%s/is_app_user/{%s}/{%s}", storeName, "accessToken", "appCode"), showIsAppUserHandler(cliCtx, storeName)).Methods("GET")
    r.HandleFunc(fmt.Sprintf("/%s/app_users/{%s}/{%s}", storeName, "accessToken", "appCode"), showAppUsersHandler(cliCtx, storeName)).Methods("GET")
    r.HandleFunc(fmt.Sprintf("/%s/upload/{%s}", storeName, "accessToken"), uploadFileHandler(cliCtx)).Methods("POST")
    r.HandleFunc(fmt.Sprintf("/%s/tables", storeName), createTableHandler(cliCtx)).Methods("POST")
    r.HandleFunc(fmt.Sprintf("/%s/tables/{%s}/{%s}/{%s}", storeName, "accessToken", "appCode", "tableName"), showTableHandler(cliCtx, storeName)).Methods("GET")
    r.HandleFunc(fmt.Sprintf("/%s/tables/{%s}/{%s}", storeName, "accessToken", "appCode"), showTablesHandler(cliCtx, storeName)).Methods("GET")
    r.HandleFunc(fmt.Sprintf("/%s/table-options/{%s}/{%s}/{%s}", storeName, "accessToken", "appCode", "tableName"), showTableOptionsHandler(cliCtx, storeName)).Methods("GET")
    r.HandleFunc(fmt.Sprintf("/%s/functions/{%s}/{%s}", storeName, "accessToken", "appCode"), showFunctionsHandler(cliCtx, storeName)).Methods("GET")
    r.HandleFunc(fmt.Sprintf("/%s/functionInfo/{%s}/{%s}/{%s}", storeName, "accessToken", "appCode", "functionName"), showFunctionInfoHandler(cliCtx, storeName)).Methods("GET")
    r.HandleFunc(fmt.Sprintf("/%s/column-options/{%s}/{%s}/{%s}/{%s}", storeName, "accessToken", "appCode", "tableName", "fieldName"), showColumnOptionsHandler(cliCtx, storeName)).Methods("GET")
    r.HandleFunc(fmt.Sprintf("/%s/can_add_column_option/{%s}/{%s}/{%s}/{%s}/{%s}", storeName, "accessToken", "appCode", "tableName", "fieldName", "option"), showCanAddColumnOptionHandler(cliCtx, storeName)).Methods("GET")
    r.HandleFunc(fmt.Sprintf("/%s/can_insert_row/{%s}/{%s}/{%s}/{%s}", storeName, "accessToken", "appCode", "tableName", "rowFieldsJsonBase58"), showCanInsertRowHandler(cliCtx, storeName)).Methods("GET")
    r.HandleFunc(fmt.Sprintf("/%s/find/{%s}/{%s}/{%s}/{%s}", storeName, "accessToken", "appCode", "name", "id"), showRowHandler(cliCtx, storeName)).Methods("GET")
    r.HandleFunc(fmt.Sprintf("/%s/find_by/{%s}/{%s}/{%s}/{%s}/{%s}", storeName, "accessToken", "appCode", "name", "field", "value"), showIdsByHandler(cliCtx, storeName)).Methods("GET")
    r.HandleFunc(fmt.Sprintf("/%s/find_all/{%s}/{%s}/{%s}", storeName, "accessToken", "appCode", "name"), showAllIdsHandler(cliCtx, storeName)).Methods("GET")
    r.HandleFunc(fmt.Sprintf("/%s/friends/{%s}", storeName, "accessToken"), showFriends(cliCtx, storeName)).Methods("GET")
    r.HandleFunc(fmt.Sprintf("/%s/pending_friends/{%s}", storeName, "accessToken"), showPendingFriends(cliCtx, storeName)).Methods("GET")
    r.HandleFunc(fmt.Sprintf("/%s/groups/{%s}/{%s}", storeName, "accessToken", "appCode"), showGroups(cliCtx, storeName)).Methods("GET")
    r.HandleFunc(fmt.Sprintf("/%s/group/{%s}/{%s}/{%s}", storeName, "accessToken", "appCode", "groupName"), showGroupMembers(cliCtx, storeName)).Methods("GET")
    r.HandleFunc(fmt.Sprintf("/%s/group_memo/{%s}/{%s}/{%s}", storeName, "accessToken", "appCode", "groupName"), showGroupMemo(cliCtx, storeName)).Methods("GET")

    r.HandleFunc(fmt.Sprintf("/%s/querier/{%s}/{%s}/{%s}", storeName, "accessToken", "appCode", "querierBase58"), execQuerier(cliCtx, storeName)).Methods("GET")

    r.HandleFunc(fmt.Sprintf("/%s/oracle/send_verf_code/{%s}/{%s}", storeName, "accessToken", "mobile"), oracleSendVerfCode(cliCtx, storeName)).Methods("GEt")
    r.HandleFunc(fmt.Sprintf("/%s/oracle/verify_verf_code/{%s}/{%s}/{%s}", storeName, "accessToken", "mobile", "verificationCode"), oracleVerifyVerfCode(cliCtx, storeName)).Methods("GET")
    r.HandleFunc(fmt.Sprintf("/%s/oracle/verify_name_and_id_number/{%s}/{%s}/{%s}", storeName, "accessToken", "name", "id_number"), oracleVerifyNameAndIdNumber(cliCtx, storeName)).Methods("GET")
    r.HandleFunc(fmt.Sprintf("/%s/oracle/verify_corp_info/{%s}/{%s}/{%s}/{%s}", storeName, "accessToken", "corp_name", "reg_number", "credit_code"), oracleVerifyCorpInfo(cliCtx, storeName)).Methods("GET")
}
