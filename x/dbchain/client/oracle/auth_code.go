package oracle

import (
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/types/rest"
	"github.com/gorilla/mux"
	"github.com/yzhanginwa/dbchain/x/dbchain/client/oracle/authenticator"
	"github.com/yzhanginwa/dbchain/x/dbchain/client/oracle/oracle"
	"net/http"
)

const (
	 authCodeInfo = "auth_code_info"
)

func showUserShareKey(cliCtx context.CLIContext) http.HandlerFunc{
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		organization := vars["organization"]
		userName := vars["user_name"]

		fieldValue := map[string]string{
			"organization": organization,
			"user_name": userName,
		}

		storeName := "dbchain"
		appcode := "0000000001"
		tableName := authCodeInfo

		//1. query from database
		secrets, err := queryByWhere(cliCtx, storeName, appcode, tableName, fieldValue)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		if len(secrets) != 0 {
			rest.PostProcessResponse(w, cliCtx, secrets[0]["secret"])
			return
		}

		//2. if dont find generate
		ga := authenticator.NewGAuth()
		secret , err := ga.CreateSecret()
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		//3. save to database
		fieldValue["secret"] = secret
		oracleAccAddr := oracle.GetOracleAccAddr()
		SaveToOrderInfoTable(cliCtx, oracleAccAddr, fieldValue, authCodeInfo)
		//4. return
		rest.PostProcessResponse(w, cliCtx, secret)
		return
	}
}



func showVerifyAuthCode(cliCtx context.CLIContext) http.HandlerFunc{
	return func(w http.ResponseWriter, r *http.Request) {

		vars := mux.Vars(r)
		organization := vars["organization"]
		userName := vars["user_name"]
		verifyCode :=  vars["auth_code"]

		fieldValue := map[string]string {
			"organization": organization,
			"user_name": userName,
		}

		storeName := "dbchain"
		appcode := "0000000001"
		tableName := authCodeInfo

		secrets, err := queryByWhere(cliCtx, storeName, appcode, tableName, fieldValue)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		if len(secrets) == 0 {
			rest.WriteErrorResponse(w, http.StatusBadRequest, "user don't exist")
			return
		}

		secret := secrets[0]["secret"]
		ga := authenticator.NewGAuth()
		ret, err := ga.VerifyCode(secret, verifyCode, 1)
		if err != nil || ret != true {
			ret = false
		}

		rest.PostProcessResponse(w, cliCtx, ret)

	}
}