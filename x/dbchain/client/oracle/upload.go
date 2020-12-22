package oracle

import (
    "fmt"
    "net/http"
    "github.com/cosmos/cosmos-sdk/client/context"
    "github.com/yzhanginwa/dbchain/x/dbchain/internal/utils"
    "github.com/cosmos/cosmos-sdk/types/rest"
    "github.com/gorilla/mux"
    shell "github.com/ipfs/go-ipfs-api"
)

func uploadFileHandler(cliCtx context.CLIContext) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        vars := mux.Vars(r)

        _, err := utils.VerifyAccessCode(vars["accessToken"])
        if err != nil {
            rest.WriteErrorResponse(w, http.StatusBadRequest, "failed to verify access token")
            return
        }

        file, _, err := r.FormFile("file")

        if err != nil {
            rest.WriteErrorResponse(w, http.StatusBadRequest, "failed to parse request")
            return
        }

	sh := shell.NewShell("localhost:5001")
	cid, err := sh.Add(file)
	if err != nil {
            rest.WriteErrorResponse(w, http.StatusBadRequest, fmt.Sprintf("%s", err))
            return
	}

        res := cid
        rest.PostProcessResponse(w, cliCtx, res)
    }
}

