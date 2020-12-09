package rest

import (
    "net/http"

    "github.com/dbchaincloud/cosmos-sdk/client/context"
    "github.com/yzhanginwa/dbchain/x/dbchain/internal/types"

    sdk "github.com/dbchaincloud/cosmos-sdk/types"
    "github.com/dbchaincloud/cosmos-sdk/types/rest"
    "github.com/dbchaincloud/cosmos-sdk/x/auth/client/utils"
)

type createTableReq struct  {
    BaseReq rest.BaseReq   `json:"base_req"`
    Owner   string         `json:"owner"`
    AppCode string         `json:"app_code"`
    Name    string         `json:"title"`
    Fields  []string       `json:"fields"`
}

func createTableHandler(cliCtx context.CLIContext) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        var req createTableReq

        if !rest.ReadRESTReq(w, r, cliCtx.Codec,  &req)  {
            rest.WriteErrorResponse(w, http.StatusBadRequest, "failed to parse request")
            return
        }

        baseReq := req.BaseReq.Sanitize()
        if !baseReq.ValidateBasic(w) {
            return
        }

        addr, err := sdk.AccAddressFromBech32(req.Owner)
        if err != nil {
            rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
            return
        }

        msg := types.NewMsgCreateTable(addr, req.AppCode, req.Name, req.Fields)
        err = msg.ValidateBasic()
        if err != nil {
            rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
            return
        }

        utils.WriteGenerateStdTxResponse(w, cliCtx, baseReq, []sdk.Msg{msg})
    }
}

