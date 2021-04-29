package rest

import (
    "encoding/json"
    "fmt"
    "github.com/cosmos/cosmos-sdk/client/context"
    shell "github.com/ipfs/go-ipfs-api"
    "github.com/yzhanginwa/dbchain/x/dbchain/internal/types"
    "net/http"
    "sync"

    "github.com/cosmos/cosmos-sdk/types/rest"

    "github.com/gorilla/mux"
)

func showCheckChainIdHandler(cliCtx context.CLIContext, storeName string) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        vars := mux.Vars(r)
        res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/check_chain_id/%s/%s", storeName, vars["accessToken"], vars["chainId"]), nil)
        if err != nil {
            rest.WriteErrorResponse(w, http.StatusNotFound, err.Error())
            return
        }
        rest.PostProcessResponse(w, cliCtx, res)
    }
}

func showIsSysAdminHandler(cliCtx context.CLIContext, storeName string) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        vars := mux.Vars(r)
        res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/is_sys_admin/%s", storeName, vars["accessToken"]), nil)
        if err != nil {
            rest.WriteErrorResponse(w, http.StatusNotFound, err.Error())
            return
        }
        rest.PostProcessResponse(w, cliCtx, res)
    }
}

func showApplicationsHandler(cliCtx context.CLIContext, storeName string) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        vars := mux.Vars(r)
        res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/application/%s", storeName, vars["accessToken"]), nil)
        if err != nil {
            rest.WriteErrorResponse(w, http.StatusNotFound, err.Error())
            return
        }
        rest.PostProcessResponse(w, cliCtx, res)
    }
}

func showApplicationHandler(cliCtx context.CLIContext, storeName string) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        vars := mux.Vars(r)
        res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/application/%s/%s", storeName, vars["accessToken"], vars["appCode"]), nil)
        if err != nil {
            rest.WriteErrorResponse(w, http.StatusNotFound, err.Error())
            return
        }
        rest.PostProcessResponse(w, cliCtx, res)
    }
}

func showApplicationUserFileVolumeLimit(cliCtx context.CLIContext, storeName string) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        vars := mux.Vars(r)
        res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/application_user_file_volume_limit/%s/%s", storeName, vars["accessToken"], vars["appCode"]), nil)
        if err != nil {
            rest.WriteErrorResponse(w, http.StatusNotFound, err.Error())
            return
        }
        rest.PostProcessResponse(w, cliCtx, res)
    }
}

func showApplicationUserUsedFileVolume(cliCtx context.CLIContext, storeName string) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        vars := mux.Vars(r)
        res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/application_user_used_file_volume/%s/%s", storeName, vars["accessToken"], vars["appCode"]), nil)
        if err != nil {
            rest.WriteErrorResponse(w, http.StatusNotFound, err.Error())
            return
        }
        rest.PostProcessResponse(w, cliCtx, res)
    }
}

func showAdminAppsHandler(cliCtx context.CLIContext, storeName string) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        vars := mux.Vars(r)
        res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/admin_apps/%s", storeName, vars["accessToken"]), nil)
        if err != nil {
            rest.WriteErrorResponse(w, http.StatusNotFound, err.Error())
            return
        }
        rest.PostProcessResponse(w, cliCtx, res)
    }
}

func showIsAppUserHandler(cliCtx context.CLIContext, storeName string) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        vars := mux.Vars(r)
        res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/is_app_user/%s/%s", storeName, vars["accessToken"], vars["appCode"]), nil)
        if err != nil {
            rest.WriteErrorResponse(w, http.StatusNotFound, err.Error())
            return
        }
        rest.PostProcessResponse(w, cliCtx, res)
    }
}

func showAppUsersHandler(cliCtx context.CLIContext, storeName string) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        vars := mux.Vars(r)
        res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/app_users/%s/%s", storeName, vars["accessToken"], vars["appCode"]), nil)
        if err != nil {
            rest.WriteErrorResponse(w, http.StatusNotFound, err.Error())
            return
        }
        rest.PostProcessResponse(w, cliCtx, res)
    }
}

func showTablesHandler(cliCtx context.CLIContext, storeName string) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        vars := mux.Vars(r)
        res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/tables/%s/%s", storeName, vars["accessToken"], vars["appCode"]), nil)
        if err != nil {
            rest.WriteErrorResponse(w, http.StatusNotFound, err.Error())
            return
        }
        rest.PostProcessResponse(w, cliCtx, res)
    }
}

func showTableHandler(cliCtx context.CLIContext, storeName string) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        vars := mux.Vars(r)
        res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/tables/%s/%s/%s", storeName, vars["accessToken"], vars["appCode"], vars["tableName"]), nil)
        if err != nil {
            rest.WriteErrorResponse(w, http.StatusNotFound, err.Error())
            return
        }
        rest.PostProcessResponse(w, cliCtx, res)
    }
}

func showFunctionsHandler(cliCtx context.CLIContext, storeName string) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        vars := mux.Vars(r)
        res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/functions/%s/%s", storeName, vars["accessToken"], vars["appCode"]), nil)
        if err != nil {
            rest.WriteErrorResponse(w, http.StatusNotFound, err.Error())
            return
        }
        rest.PostProcessResponse(w, cliCtx, res)

    }
}

func showFunctionInfoHandler(cliCtx context.CLIContext, storeName string) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        vars := mux.Vars(r)
        res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/functionInfo/%s/%s/%s", storeName, vars["accessToken"], vars["appCode"], vars["functionName"]), nil)
        if err != nil {
            rest.WriteErrorResponse(w, http.StatusNotFound, err.Error())
            return
        }
        rest.PostProcessResponse(w, cliCtx, res)
    }
}

func showCustomQueriersHandler(cliCtx context.CLIContext, storeName string) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        vars := mux.Vars(r)
        res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/customQueriers/%s/%s", storeName, vars["accessToken"], vars["appCode"]), nil)
        if err != nil {
            rest.WriteErrorResponse(w, http.StatusNotFound, err.Error())
            return
        }
        rest.PostProcessResponse(w, cliCtx, res)
    }
}

func showCustomQuerierInfoHandler(cliCtx context.CLIContext, storeName string) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        vars := mux.Vars(r)
        res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/customQuerierInfo/%s/%s/%s", storeName, vars["accessToken"], vars["appCode"], vars["querierName"]), nil)
        if err != nil {
            rest.WriteErrorResponse(w, http.StatusNotFound, err.Error())
            return
        }
        rest.PostProcessResponse(w, cliCtx, res)
    }
}

func showCallCustomQuerierHandler(cliCtx context.CLIContext, storeName string) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        vars := mux.Vars(r)
        res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/callCustomQuerier/%s/%s/%s/%s", storeName, vars["accessToken"], vars["appCode"], vars["querierName"], vars["params"]), nil)
        if err != nil {
            rest.WriteErrorResponse(w, http.StatusNotFound, err.Error())
            return
        }
        rest.PostProcessResponse(w, cliCtx, res)
    }
}

func showTableOptionsHandler(cliCtx context.CLIContext, storeName string) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        vars := mux.Vars(r)
        res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/option/%s/%s/%s", storeName, vars["accessToken"], vars["appCode"], vars["tableName"]), nil)
        if err != nil {
            rest.WriteErrorResponse(w, http.StatusNotFound, err.Error())
            return
        }
        rest.PostProcessResponse(w, cliCtx, res)
    }
}

func showColumnOptionsHandler(cliCtx context.CLIContext, storeName string) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        vars := mux.Vars(r)
        res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/column_option/%s/%s/%s/%s", storeName, vars["accessToken"], vars["appCode"], vars["tableName"], vars["fieldName"]), nil)
        if err != nil {
            rest.WriteErrorResponse(w, http.StatusNotFound, err.Error())
            return
        }
        rest.PostProcessResponse(w, cliCtx, res)
    }
}

func showColumnDataTypeHandler(cliCtx context.CLIContext, storeName string) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        vars := mux.Vars(r)
        res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/column_data_type/%s/%s/%s/%s", storeName, vars["accessToken"], vars["appCode"], vars["tableName"], vars["fieldName"]), nil)
        if err != nil {
            rest.WriteErrorResponse(w, http.StatusNotFound, err.Error())
            return
        }
        rest.PostProcessResponse(w, cliCtx, res)
    }
}

func showCanAddColumnOptionHandler(cliCtx context.CLIContext, storeName string) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        vars := mux.Vars(r)
        res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/can_add_column_option/%s/%s/%s/%s/%s", storeName, vars["accessToken"], vars["appCode"], vars["tableName"], vars["fieldName"], vars["option"]), nil)
        if err != nil {
            rest.WriteErrorResponse(w, http.StatusNotFound, err.Error())
            return
        }
        rest.PostProcessResponse(w, cliCtx, res)
    }
}

func showCanSetColumnDataTypeHandler(cliCtx context.CLIContext, storeName string) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        vars := mux.Vars(r)
        res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/can_set_column_data_type/%s/%s/%s/%s/%s", storeName, vars["accessToken"], vars["appCode"], vars["tableName"], vars["fieldName"], vars["dataType"]), nil)
        if err != nil {
            rest.WriteErrorResponse(w, http.StatusNotFound, err.Error())
            return
        }
        rest.PostProcessResponse(w, cliCtx, res)
    }
}

func showCanInsertRowHandler(cliCtx context.CLIContext, storeName string) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        vars := mux.Vars(r)
        rowFieldsJson := vars["rowFieldsJsonBase58"]
        res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/can_insert_row/%s/%s/%s/%s", storeName, vars["accessToken"], vars["appCode"], vars["tableName"], rowFieldsJson), nil)
        if err != nil {
            rest.WriteErrorResponse(w, http.StatusNotFound, err.Error())
            return
        }
        rest.PostProcessResponse(w, cliCtx, res)
    }
}

func showRowHandler(cliCtx context.CLIContext, storeName string) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        vars := mux.Vars(r)
        res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/find/%s/%s/%s/%s", storeName, vars["accessToken"], vars["appCode"], vars["name"], vars["id"]), nil)
        if err != nil {
            rest.WriteErrorResponse(w, http.StatusNotFound, err.Error())
            return
        }
        rest.PostProcessResponse(w, cliCtx, res)
    }
}

func showIdsByHandler(cliCtx context.CLIContext, storeName string) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        vars := mux.Vars(r)
        res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/find_by/%s/%s/%s/%s/%s", storeName, vars["accessToken"], vars["appCode"], vars["name"], vars["field"], vars["value"]), nil)
        if err != nil {
            rest.WriteErrorResponse(w, http.StatusNotFound, err.Error())
            return
        }
        rest.PostProcessResponse(w, cliCtx, res)
    }
}

func showAllIdsHandler(cliCtx context.CLIContext, storeName string) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        vars := mux.Vars(r)
        res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/find_all/%s/%s/%s", storeName, vars["accessToken"], vars["appCode"], vars["name"]), nil)
        if err != nil {
            rest.WriteErrorResponse(w, http.StatusNotFound, err.Error())
            return
        }
        rest.PostProcessResponse(w, cliCtx, res)
    }
}

func showFriends(cliCtx context.CLIContext, storeName string) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        vars := mux.Vars(r)
        res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/friends/%s", storeName, vars["accessToken"]), nil)
        if err != nil {
            rest.WriteErrorResponse(w, http.StatusNotFound, err.Error())
            return
        }
        rest.PostProcessResponse(w, cliCtx, res)
    }
}

func showPendingFriends(cliCtx context.CLIContext, storeName string) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        vars := mux.Vars(r)
        res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/pending_friends/%s", storeName, vars["accessToken"]), nil)
        if err != nil {
            rest.WriteErrorResponse(w, http.StatusNotFound, err.Error())
            return
        }
        rest.PostProcessResponse(w, cliCtx, res)
    }
}

func showGroups(cliCtx context.CLIContext, storeName string) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        vars := mux.Vars(r)
        res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/groups/%s/%s", storeName, vars["accessToken"], vars["appCode"]), nil)
        if err != nil {
            rest.WriteErrorResponse(w, http.StatusNotFound, err.Error())
            return
        }
        rest.PostProcessResponse(w, cliCtx, res)
    }
}

func showGroupMembers(cliCtx context.CLIContext, storeName string) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        vars := mux.Vars(r)
        res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/group/%s/%s/%s", storeName, vars["accessToken"], vars["appCode"], vars["groupName"]), nil)
        if err != nil {
            rest.WriteErrorResponse(w, http.StatusNotFound, err.Error())
            return
        }
        rest.PostProcessResponse(w, cliCtx, res)
    }
}

func showGroupMemo(cliCtx context.CLIContext, storeName string) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        vars := mux.Vars(r)
        res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/group_memo/%s/%s/%s", storeName, vars["accessToken"], vars["appCode"], vars["groupName"]), nil)
        if err != nil {
            rest.WriteErrorResponse(w, http.StatusNotFound, err.Error())
            return
        }
        rest.PostProcessResponse(w, cliCtx, res)
    }
}

func showIndex(cliCtx context.CLIContext, storeName string) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        vars := mux.Vars(r)
        res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/index/%s/%s/%s", storeName, vars["accessToken"], vars["appCode"], vars["tableName"]), nil)
        if err != nil {
            rest.WriteErrorResponse(w, http.StatusNotFound, err.Error())
            return
        }
        rest.PostProcessResponse(w, cliCtx, res)
    }
}

func execQuerier(cliCtx context.CLIContext, storeName string) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        vars := mux.Vars(r)
        res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/querier/%s/%s/%s", storeName, vars["accessToken"], vars["appCode"], vars["querierBase58"]), nil)
        if err != nil {
            rest.WriteErrorResponse(w, http.StatusNotFound, err.Error())
            return
        }
        rest.PostProcessResponse(w, cliCtx, res)
    }
}

func showTxSimpleResultHandler(cliCtx context.CLIContext, storeName string) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        vars := mux.Vars(r)
        res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/txSimpleResult/%s/%s", storeName, vars["accessToken"], vars["txHash"]), nil)
        if err != nil {
            rest.WriteErrorResponse(w, http.StatusNotFound, err.Error())
            return
        }
        rest.PostProcessResponse(w, cliCtx, res)
    }
}

func downloadFileHandler(cliCtx context.CLIContext, storeName string) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        vars := mux.Vars(r)
        accessToken := vars["accessToken"]
        appCode := vars["appCode"]
        tableName := vars["tableName"]
        id := vars["id"]
        fieldName := vars["fieldName"]
        var cid string
        //check type
        ch := make(chan string, 2)
        defer func() {
            close(ch)
        }()
        var wait sync.WaitGroup
        wait.Add(2)
        go func() {
            defer wait.Done()
            res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/column_data_type/%s/%s/%s/%s", storeName, accessToken, appCode, tableName, fieldName), nil)
            if err != nil {
                ch <- err.Error()
                return
            }
            var dataType string
            err = json.Unmarshal(res, &dataType)
            if err != nil {
                ch <- err.Error()
                return
            }
            if dataType != string(types.FLDTYP_FILE) {
                ch <- err.Error()
                return
            }
        }()

        //check value
        go func() {
            defer wait.Done()
            res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/find/%s/%s/%s/%s", storeName, accessToken, appCode, tableName, id), nil)
            if err != nil {
                ch <- err.Error()
                return
            }
            var  RowFields map[string]string
            err = json.Unmarshal(res, &RowFields)
            if err != nil {
                ch <- err.Error()
                return
            }
            cid = RowFields[fieldName]
        }()

        wait.Wait()
        if len(ch) > 0 {
            rest.PostProcessResponse(w, cliCtx, <-ch)
            return
        }
        //get file
        sh := shell.NewShell("localhost:5001")
        w.Header().Set("Content-Type", "application/octet-stream")
        reader,_ := sh.Cat(cid)
        buf := make([]byte, 4096)
        for {
            n , err := reader.Read(buf)
            if err != nil {
                if n < 4096 {
                    buf = buf[:n]
                }
                w.Write(buf)
                break
            }
            w.Write(buf)
        }
    }
}