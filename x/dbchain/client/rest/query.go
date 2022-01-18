package rest

import (
    "encoding/hex"
    "encoding/json"
    "fmt"
    "github.com/cosmos/cosmos-sdk/client/context"
    "github.com/cosmos/cosmos-sdk/client/flags"
    keys2 "github.com/cosmos/cosmos-sdk/client/keys"
    "github.com/cosmos/cosmos-sdk/crypto/keys"
    sdk "github.com/cosmos/cosmos-sdk/types"
    "github.com/cosmos/cosmos-sdk/types/rest"
    "github.com/cosmos/cosmos-sdk/x/auth/client/utils"
    "github.com/cosmos/go-bip39"
    shell "github.com/ipfs/go-ipfs-api"
    "github.com/mr-tron/base58"
    "github.com/spf13/viper"
    "github.com/tendermint/tendermint/crypto"
    tmamino "github.com/tendermint/tendermint/crypto/encoding/amino"
    "github.com/yzhanginwa/dbchain/x/dbchain/client/oracle/oracle"
    "github.com/yzhanginwa/dbchain/x/dbchain/internal/types"
    "io"
    "io/ioutil"
    "net/http"
    "strconv"
    "strings"
    "sync"

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

func showTablesDetailHandler(cliCtx context.CLIContext, storeName string) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        vars := mux.Vars(r)
        res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/tables_detail/%s/%s", storeName, vars["accessToken"], vars["appCode"]), nil)
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

func showTableDetailHandler(cliCtx context.CLIContext, storeName string) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        vars := mux.Vars(r)
        res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/table_detail/%s/%s/%s", storeName, vars["accessToken"], vars["appCode"], vars["tableName"]), nil)
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
        r.ParseForm()
        originParams := r.Form["params"]
        params := strings.Join(originParams, "/")
        params = base58.Encode([]byte(params))
        res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/callCustomQuerier/%s/%s/%s/%s", storeName, vars["accessToken"], vars["appCode"], vars["querierName"], params), nil)
        if err != nil {
            rest.WriteErrorResponse(w, http.StatusNotFound, err.Error())
            return
        }
        rest.PostProcessResponse(w, cliCtx, res)
    }
}

func showCallDynamicScriptHandler(cliCtx context.CLIContext, storeName string) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        vars := mux.Vars(r)
        res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/dynamic_script/%s/%s/%s", storeName, vars["accessToken"], vars["appCode"], vars["script"]), nil)
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

func showTableAssociationsHandler(cliCtx context.CLIContext, storeName string) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        vars := mux.Vars(r)
        res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/association/%s/%s/%s", storeName, vars["accessToken"], vars["appCode"], vars["tableName"]), nil)
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
            //n maybe less than 4096 but err is nil, so it needs be checked everytime
            if n < 4096 {
                w.Write(buf[:n])
            } else {
                w.Write(buf)
            }
        }
    }
}

//add for bsn
func showAccountTxs(cliCtx context.CLIContext, storeName string) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        vars := mux.Vars(r)
        r.ParseForm()
        number := r.Form["number"]
        queryString := fmt.Sprintf("custom/%s/account_txs/%s", storeName, vars["accessToken"])
        if number != nil && len(number) != 0 {
            queryString += "/" + number[0]
        }
        res, _, err := cliCtx.QueryWithData(queryString, nil)
        if err != nil {
            rest.WriteErrorResponse(w, http.StatusNotFound, err.Error())
            return
        }
        rest.PostProcessResponse(w, cliCtx, res)
    }
}

func showChainSuperAdmins(cliCtx context.CLIContext, storeName string) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        vars := mux.Vars(r)
        res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/chain_super_admins/%s", storeName, vars["accessToken"]), nil)
        if err != nil {
            rest.WriteErrorResponse(w, http.StatusNotFound, err.Error())
            return
        }
        rest.PostProcessResponse(w, cliCtx, res)
    }
}

func showLimitP2PTransferStatus(cliCtx context.CLIContext, storeName string) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        vars := mux.Vars(r)
        res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/limit_p2p_transfer_status/%s", storeName, vars["accessToken"]), nil)
        if err != nil {
            rest.WriteErrorResponse(w, http.StatusNotFound, err.Error())
            return
        }
        rest.PostProcessResponse(w, cliCtx, res)
    }
}

func showAllTxs(cliCtx context.CLIContext, storeName string) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {

        vars := mux.Vars(r)
        startHeight := vars["start_height"]
        endHeight := vars["end_height"]

        start, err := strconv.Atoi(startHeight)
        if err != nil {
            generalResponse(w,map[string]string { "error" : "invalid parameter" })
            return
        }

        end, err := strconv.Atoi(endHeight)
        if err != nil {
            generalResponse(w,map[string]string { "error" : "invalid parameter" })
            return
        }

        if end - start <= 0 {
            generalResponse(w,map[string]string { "error" : "invalid parameter" })
            return
        }

        if start - end > 17280 {
            generalResponse(w,map[string]string { "error" : "number of query blocks can not more than 17280" })
            return
        }

        node, err := cliCtx.GetNode()
        if err != nil {
            generalResponse(w,map[string]string { "error" : "GetNode err : " + err.Error()})
            return
        }
        result := make([]sdk.TxResponse,0)
        for i := start ; i <= end; i++ {
            height := int64(i)
            block, err := node.Block(&height)
            if err != nil {
                rest.WriteErrorResponse(w, http.StatusBadRequest, "get block err : " + err.Error())
                return
            }
            Txs := block.Block.Txs
            for _,tx := range Txs {
                txha := hex.EncodeToString(tx.Hash())
                out, err := utils.QueryTx(cliCtx,txha)
                if err != nil {
                    continue
                }
                result = append(result, out)
            }

        }
        rest.PostProcessResponseBare(w, cliCtx, result)
    }
}


func applyAccountInfo() http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        kb, err := keys.NewKeyring(sdk.KeyringServiceName(), keys.BackendOS, viper.GetString(flags.FlagHome), nil)
        if err != nil {
           generalResponse(w, map[string]string {"error " : err.Error()})
           return
        }
        //genName
        nameBytes := make([]byte, 24)
        keyName := ""
        for i := 0; i < 10; i++ {
           io.ReadFull(crypto.CReader(), nameBytes)
           keyName = hex.EncodeToString(nameBytes)
           info, err := kb.Get(keyName)
           if err != nil {
               continue
           }
           if info != nil && i == 9 {
               generalResponse(w, map[string]string {"error " : "generate key pairs err"})
               return
           } else if info == nil {
               break
           }
        }

        info, secret, err := CreateMnemonic(kb, keyName, keys.English, keys2.DefaultKeyPass, keys.Secp256k1)
        if err != nil {
           generalResponse(w, map[string]string {"error " : "generate key pairs err"})
           return
        }

        pk, err := kb.ExportPrivateKeyObject(keyName, keys2.DefaultKeyPass)
        if err != nil {
           generalResponse(w, map[string]string {"error " : "generate key pairs err"})
           return
        }

        add := sdk.AccAddress(info.GetPubKey().Address())

        data := map[string]string {
            "publicKey" : hex.EncodeToString(pk.PubKey().Bytes()),
            "privateKey" : hex.EncodeToString(pk.Bytes()),
            "address" : add.String(),
            "mnemonic" : secret,
        }
        generalResponse(w, data)
        return

    }
}

func applyAccountInfoByPublicKey() http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        result, err  := ioutil.ReadAll(r.Body)
        if err != nil {
            generalResponse(w, map[string]string{"error" : err.Error()})
            return
        }
        postData := make(map[string]string)
        err = json.Unmarshal(result, &postData)
        if err != nil {
            generalResponse(w, map[string]string{"error" : err.Error()})
            return
        }
        publicKey := postData["publicKey"]
        pubBytes , err := hex.DecodeString(publicKey)
        if err != nil {
            generalResponse(w, map[string]string{"error" : err.Error()})
            return
        }
        pubKey, err  := tmamino.PubKeyFromBytes(pubBytes)
        if err != nil {
            generalResponse(w, map[string]string{"error" : "Public key format should be hexadecimal string"})
            return
        }

        add := sdk.AccAddress(pubKey.Address())
        data := map[string]string {
            "publicKey" : publicKey,
            "address" : add.String(),
        }
        generalResponse(w, data)
        return

    }
}

func rechargeTx(cliCtx context.CLIContext, storeName string) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        data, err := readBodyData(r)
        if err != nil {
            generalResponse(w, map[string]string{ "error" : err.Error()})
        }
        bsnAddress := data["bsnAddress"]
        userAccountAddress := data["userAccountAddress"]
        rechargeGas := data["rechargeGas"]
        tx, status, errInfo := sendFromBsnAddressToUserAddress(cliCtx, bsnAddress, userAccountAddress, rechargeGas)
        generalResponse(w, map[string]interface{}{
            "txHash" : tx,
            "state" : status,
            "remarks" : errInfo,
        })
        return
    }
}

func getAccountTxByTime(cliCtx context.CLIContext, storeName string) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        data, err := readBodyData(r)
        if err != nil {
            generalResponse(w, map[string]string{ "error" : err.Error()})
        }

        userAccountAddress := data["userAccountAddress"]
        startDate := data["startDate"]
        endDate := data["endDate"]
        if userAccountAddress == "" || startDate == "" || endDate == "" {
            generalResponse(w, map[string]string{"error" : "expect 3 parameters : userAccountAddress, startDate, endDate"})
            return
        }

        year, month, day := 0,0,0
        nStartDate , _ := fmt.Sscanf(startDate,"%d-%d-%d", &year, &month, &day)
        nEndDate , _ := fmt.Sscanf(endDate,"%d-%d-%d", &year, &month, &day)
        if nStartDate != 3 || nEndDate != 3 {
            generalResponse(w, map[string]string{"error" : "time format error, it should be  yyyy-mm-dd"})
            return
        }



        res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/account_txs_by_time/%s/%s/%s", storeName, userAccountAddress, startDate, endDate), nil)
        if err != nil {
            rest.WriteErrorResponse(w, http.StatusNotFound, err.Error())
            return
        }
        rest.PostProcessResponse(w, cliCtx, res)
    }
}
///////////////////
//               //
//   help func   //
//               //
///////////////////

func generalResponse(w http.ResponseWriter, data interface{}) {
    bz,_ := json.Marshal(data)
    w.Header().Set("Content-Type", "application/json")
    _, _ = w.Write(bz)
}

func CreateMnemonic(
    kb keys.Keybase,name string, language keys.Language, passwd string, algo keys.SigningAlgo,
) (keys.Info,  string,  error) {

    entropy, err := bip39.NewEntropy(128)
    if err != nil {
        return  nil, "", err
    }

    mnemonic, err := bip39.NewMnemonic(entropy)
    if err != nil {
        return nil, "", err
    }

    info, err := kb.CreateAccount( name, mnemonic, keys.DefaultBIP39Passphrase, passwd, sdk.GetConfig().GetFullFundraiserPath(), algo)
    if err != nil {
        return nil, "", err
    }
    return info, mnemonic, nil
}

func readBodyData(r *http.Request) (map[string]string, error) {
    result, err  := ioutil.ReadAll(r.Body)
    if err != nil {
        return nil , err
    }

    postData := make(map[string]string)
    err = json.Unmarshal(result, &postData)
    if err != nil {
        return nil, err
    }
    return postData, nil
}

func sendFromBsnAddressToUserAddress(cliCtx context.CLIContext, bsnAddress, userAccountAddress, rechargeGas string) (string, int, string){
    from, err  := sdk.AccAddressFromBech32(bsnAddress)
    if err != nil {
        return "", oracle.Failed, err.Error()
    }
    to, err := sdk.AccAddressFromBech32(userAccountAddress)
    if err != nil {
        return "", oracle.Failed, err.Error()
    }
    coins, err := sdk.ParseCoins(rechargeGas)
    if err != nil {
        return "", oracle.Failed, err.Error()
    }

    kb, err := keys.NewKeyring(sdk.KeyringServiceName(), keys.BackendOS, viper.GetString(flags.FlagHome), nil)
    if err != nil {
        return "", oracle.Failed, err.Error()
    }

    info , err := kb.GetByAddress(from)
    if err != nil {
        return "", oracle.Failed, err.Error()
    }
    fmt.Println(info.GetName())

    pk , err := kb.ExportPrivateKeyObject(info.GetName(), keys2.DefaultKeyPass)

    msg := oracle.NewMsgSend(from, to, coins)
    txHash, status, errInfo := oracle.BuildAndSignBroadcastTx(cliCtx, []oracle.UniversalMsg{msg}, pk, from)
    fmt.Println(txHash, status, errInfo)
    return txHash, status, errInfo
}
