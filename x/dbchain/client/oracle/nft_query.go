package oracle

import (
	"encoding/json"
	"fmt"
	"github.com/dbchaincloud/cosmos-sdk/client/context"
	"github.com/dbchaincloud/cosmos-sdk/types/rest"
	"github.com/gorilla/mux"
	"github.com/yzhanginwa/dbchain/x/dbchain/client/oracle/oracle"
	"github.com/yzhanginwa/dbchain/x/dbchain/internal/types"
	"io/ioutil"
	"net/http"
)

const BaseUrl = oracle.BaseUrl + "dbchain/"

func CanEditPersonalInfo(cliCtx context.CLIContext, storeName string, tel string) (string,bool) {
	ac := getOracleAc()
	userId, err := findByCoreIds(cliCtx, storeName, ac, nftAppCode, nftUserTable, "tel", tel)
	if err != nil || len(userId) == 0{
		return "",false
	}
	editId, err := findByCoreIds(cliCtx, storeName, ac, nftAppCode, nftUserInfoTable, "user_id", userId[0])
	if err != nil {
		return "", false
	}

	if len(editId) == 0 {
		return userId[0], true
	}
	return "", false
}

func nftFindById(cliCtx context.CLIContext, storeName string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		ac := getOracleAc()
		requestUrl := fmt.Sprintf("%s/find/%s/%s/%s/%s", BaseUrl,ac, nftAppCode,vars["name"], vars["id"])
		res, err := httpGetRequest(requestUrl)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusNotFound, err.Error())
			return
		}
		rest.PostProcessResponse(w, cliCtx, res)
	}
}

func nftFindByField(cliCtx context.CLIContext, storeName string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		ac := getOracleAc()
		requestUrl := fmt.Sprintf("%s/find_by/%s/%s/%s/%s/%s", BaseUrl,ac, nftAppCode,vars["name"], vars["field"], vars["value"])
		res, err := httpGetRequest(requestUrl)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusNotFound, err.Error())
			return
		}
		rest.PostProcessResponse(w, cliCtx, res)
	}
}

func nftFindAll(cliCtx context.CLIContext, storeName string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		ac := getOracleAc()
		requestUrl := fmt.Sprintf("%s/find_all/%s/%s/%s", BaseUrl,ac, nftAppCode,vars["name"])
		res, err := httpGetRequest(requestUrl)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusNotFound, err.Error())
			return
		}
		rest.PostProcessResponse(w, cliCtx, res)
	}
}

func nftFindByQuerier(cliCtx context.CLIContext, storeName string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		ac := getOracleAc()
		requestUrl := fmt.Sprintf("%s/querier/%s/%s/%s", BaseUrl,ac, nftAppCode,vars["querierBase58"])
		res, err := httpGetRequest(requestUrl)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusNotFound, err.Error())
			return
		}
		rest.PostProcessResponse(w, cliCtx, res)
	}
}

//////////////////////////
//                      //
//      help func       //
//                      //
//////////////////////////

func httpGetRequest(url string) ([]byte, error) {
	resp, err := http.Get(url)
	defer resp.Body.Close()
	if err != nil {
		return nil, err
	}
	bz, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return bz, nil
}


func findByAll(cliCtx context.CLIContext, storeName, ac, appcode, tableName, fieldName, value string ) ([]map[string]string, error){
	out, err := findByCoreIds(cliCtx, storeName, ac, appcode, tableName, fieldName, value)
	if err != nil {
		return nil, err
	}
	result := make([]map[string]string, 0)
	for _, id := range out {
		queryString := fmt.Sprintf("%s/find/%s/%s/%s/%s", BaseUrl, ac, appcode, tableName, id)
		res, err := findRow(cliCtx, queryString)
		if err != nil {
			continue
		}
		result = append(result, res)
	}
	return result, nil
}


func findByCore(cliCtx context.CLIContext, storeName, ac, appcode, tableName, fieldName, value string ) (map[string]string, error){
	out, err := findByCoreIds(cliCtx, storeName, ac, appcode, tableName, fieldName, value)
	if err != nil {
		return nil, err
	}
	id := out[len(out) - 1]
	queryString := fmt.Sprintf("%s/find/%s/%s/%s/%s", BaseUrl, ac, appcode, tableName, id)
	return findRow(cliCtx, queryString)
}

func findByCoreIds(cliCtx context.CLIContext, storeName, ac, appcode, tableName, fieldName, value string ) (types.QuerySliceOfString, error) {
	requestUrl := fmt.Sprintf("%s/find_by/%s/%s/%s/%s/%s", BaseUrl, ac, appcode, tableName, fieldName, value)
	res, err := httpGetRequest(requestUrl)
	if err != nil {
		fmt.Printf("could not find ids")
		return nil, err
	}
	type response struct {
		Height string
		Result types.QuerySliceOfString
	}
	temp := response{}
	json.Unmarshal(res, &temp)
	return temp.Result, nil
}

func findRow(cliCtx context.CLIContext, query string) (map[string]string, error) {
	res, err := httpGetRequest(query)
	if err != nil {
		return nil, err
	}
	type response struct {
		Height string
		Result map[string]string
	}

	temp := response{}
	err = json.Unmarshal(res, &temp)
	if err != nil {
		return nil, err
	}
	return temp.Result, err
}