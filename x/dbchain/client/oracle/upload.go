package oracle

import (
    "encoding/json"
    "fmt"

    //"github.com/cosmos/cosmos-sdk/client/context"
    "github.com/cosmos/cosmos-sdk/client"

    sdk "github.com/cosmos/cosmos-sdk/types"
    "github.com/cosmos/cosmos-sdk/types/rest"
    "github.com/gorilla/mux"
    "github.com/yzhanginwa/dbchain/x/dbchain/internal/utils"
    "net/http"
    "strconv"
)

func uploadFileHandler(cliCtx client.Context) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        vars := mux.Vars(r)
        accessToken := vars["accessToken"]
        appcode := vars["appCode"]
        addr, _, err := utils.VerifyAccessCodeWithoutTimeChecking(accessToken)
        if err != nil {
            rest.WriteErrorResponse(w, http.StatusBadRequest, "failed to verify access token")
            return
        }

        if !checkAppUserFileVolumeLimit(cliCtx, accessToken, appcode, addr) {
            rest.WriteErrorResponse(w, http.StatusBadRequest, "your volume has been exhausted")
            return
        }

        file, _, err := r.FormFile("file")
        if err != nil {
            rest.WriteErrorResponse(w, http.StatusBadRequest, "failed to parse request")
            return
        }

        sh := utils.NewShellDbchain("localhost:5001")
        cid, err := sh.Add(file)
        if err != nil {
            rest.WriteErrorResponse(w, http.StatusBadRequest, fmt.Sprintf("%s", err))
            return
        }

        rest.PostProcessResponse(w, cliCtx, cid)
    }
}


func checkAppUserFileVolumeLimit(cliCtx client.Context, accessToken, appCode string, owner sdk.AccAddress) bool {

    //1、get App user file volume limit
    res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/application_user_file_volume_limit/%s/%s", "dbchain", accessToken, appCode), nil)
    if err != nil {
        return false
    }
    limitSize := ""
    err = json.Unmarshal(res, &limitSize)
    if err != nil || limitSize == "no limit"{
        return true //
    }
    //2、get user file volume of used
    res, _, err = cliCtx.QueryWithData(fmt.Sprintf("custom/%s/application_user_used_file_volume/%s/%s", "dbchain", accessToken, appCode), nil)
    if err != nil {
        return false
    }
    usedSize := ""
    err = json.Unmarshal(res, &usedSize)
    if err != nil {
        usedSize = "0"
    }
    //3、compare
    iLimitSize, _ := strconv.ParseInt(limitSize, 10, 64)
    if iLimitSize <= 0 {
        return true  //when size was set 0 or negative, it means no limit
    }
    iUsedSize, _  := strconv.ParseInt(usedSize, 10, 64)
    if iUsedSize < iLimitSize {
        return true
    }
    return false
}

