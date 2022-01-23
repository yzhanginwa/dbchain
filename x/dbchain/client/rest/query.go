package rest

import (
    "context"
    "encoding/hex"
    "encoding/json"
    "fmt"

    //"github.com/cosmos/cosmos-sdk/client/context"
    "github.com/cosmos/cosmos-sdk/client"

    sdk "github.com/cosmos/cosmos-sdk/types"
    "github.com/cosmos/cosmos-sdk/types/rest"

    //"github.com/cosmos/cosmos-sdk/x/auth/client/utils"
    authtx "github.com/cosmos/cosmos-sdk/x/auth/tx"

    "github.com/mr-tron/base58"
    "github.com/yzhanginwa/dbchain/x/dbchain/internal/types"
    mutils "github.com/yzhanginwa/dbchain/x/dbchain/internal/utils"
    "net/http"
    "strconv"
    "strings"
    "sync"

    "github.com/gorilla/mux"
)

func showCheckChainIdHandler(cliCtx client.Context, storeName string) http.HandlerFunc {
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

func showIsSysAdminHandler(cliCtx client.Context, storeName string) http.HandlerFunc {
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

func showApplicationsHandler(cliCtx client.Context, storeName string) http.HandlerFunc {
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

func showApplicationHandler(cliCtx client.Context, storeName string) http.HandlerFunc {
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

func showApplicationUserFileVolumeLimit(cliCtx client.Context, storeName string) http.HandlerFunc {
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

func showApplicationUserUsedFileVolume(cliCtx client.Context, storeName string) http.HandlerFunc {
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

func showAdminAppsHandler(cliCtx client.Context, storeName string) http.HandlerFunc {
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

func showIsAppUserHandler(cliCtx client.Context, storeName string) http.HandlerFunc {
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

func showAppUsersHandler(cliCtx client.Context, storeName string) http.HandlerFunc {
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

func showTablesHandler(cliCtx client.Context, storeName string) http.HandlerFunc {
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

func showTableHandler(cliCtx client.Context, storeName string) http.HandlerFunc {
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

func showFunctionsHandler(cliCtx client.Context, storeName string) http.HandlerFunc {
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

func showFunctionInfoHandler(cliCtx client.Context, storeName string) http.HandlerFunc {
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

func showCustomQueriersHandler(cliCtx client.Context, storeName string) http.HandlerFunc {
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

func showCustomQuerierInfoHandler(cliCtx client.Context, storeName string) http.HandlerFunc {
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

func showCallCustomQuerierHandler(cliCtx client.Context, storeName string) http.HandlerFunc {
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

func showCallDynamicScriptHandler(cliCtx client.Context, storeName string) http.HandlerFunc {
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

func showTableOptionsHandler(cliCtx client.Context, storeName string) http.HandlerFunc {
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

func showTableAssociationsHandler(cliCtx client.Context, storeName string) http.HandlerFunc {
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

func showColumnOptionsHandler(cliCtx client.Context, storeName string) http.HandlerFunc {
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

func showCounterCacheHandler(cliCtx client.Context, storeName string) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        vars := mux.Vars(r)
        res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/counter_info/%s/%s/%s", storeName, vars["accessToken"], vars["appCode"], vars["tableName"]), nil)
        if err != nil {
            rest.WriteErrorResponse(w, http.StatusNotFound, err.Error())
            return
        }
        rest.PostProcessResponse(w, cliCtx, res)
    }
}


func showColumnDataTypeHandler(cliCtx client.Context, storeName string) http.HandlerFunc {
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

func showCanAddColumnOptionHandler(cliCtx client.Context, storeName string) http.HandlerFunc {
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

func showCanSetColumnDataTypeHandler(cliCtx client.Context, storeName string) http.HandlerFunc {
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

func showCanInsertRowHandler(cliCtx client.Context, storeName string) http.HandlerFunc {
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

func showRowHandler(cliCtx client.Context, storeName string) http.HandlerFunc {
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

func showIdsByHandler(cliCtx client.Context, storeName string) http.HandlerFunc {
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

func showAllIdsHandler(cliCtx client.Context, storeName string) http.HandlerFunc {
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

func showFriends(cliCtx client.Context, storeName string) http.HandlerFunc {
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

func showPendingFriends(cliCtx client.Context, storeName string) http.HandlerFunc {
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

func showGroups(cliCtx client.Context, storeName string) http.HandlerFunc {
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

func showGroupMembers(cliCtx client.Context, storeName string) http.HandlerFunc {
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

func showGroupMemo(cliCtx client.Context, storeName string) http.HandlerFunc {
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

func showIndex(cliCtx client.Context, storeName string) http.HandlerFunc {
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

func execQuerier(cliCtx client.Context, storeName string) http.HandlerFunc {
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

func showTxSimpleResultHandler(cliCtx client.Context, storeName string) http.HandlerFunc {
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

func downloadFileHandler(cliCtx client.Context, storeName string) http.HandlerFunc {
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
        sh := mutils.NewShellDbchain("localhost:5001")
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
func showAccountTxs(cliCtx client.Context, storeName string) http.HandlerFunc {
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

func showTokenKeepers(cliCtx client.Context, storeName string) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        vars := mux.Vars(r)
        res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/token_keepers/%s", storeName, vars["accessToken"]), nil)
        if err != nil {
            rest.WriteErrorResponse(w, http.StatusNotFound, err.Error())
            return
        }
        rest.PostProcessResponse(w, cliCtx, res)
    }
}

func showLimitP2PTransferStatus(cliCtx client.Context, storeName string) http.HandlerFunc {
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

func showAllTxs(cliCtx client.Context, storeName string) http.HandlerFunc {
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
        result := make([]*sdk.TxResponse,0)
        for i := start ; i <= end; i++ {
            height := int64(i)
            block, err := node.Block(context.Background(), &height)
            if err != nil {
                rest.WriteErrorResponse(w, http.StatusBadRequest, "get block err : " + err.Error())
                return
            }
            Txs := block.Block.Txs
            for _,tx := range Txs {
                txha := hex.EncodeToString(tx.Hash())
                out, err := authtx.QueryTx(cliCtx,txha)
                if err != nil {
                    continue
                }
                result = append(result, out)
            }

        }
        rest.PostProcessResponseBare(w, cliCtx, result)
    }
}

func showCurrentMinGasPrices(cliCtx client.Context, storeName string) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        vars := mux.Vars(r)
        ac := vars["accessToken"]
        _, err := mutils.VerifyAccessCode(ac)
        if err != nil {
            rest.WriteErrorResponse(w, http.StatusNotFound, err.Error())
            return
        }

        res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/current_min_gas_prices/%s", storeName, vars["accessToken"]), nil)
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
