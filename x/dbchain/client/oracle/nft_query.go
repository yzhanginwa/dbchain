package oracle

import (
	"encoding/json"
	"fmt"
	"github.com/dbchaincloud/cosmos-sdk/client/context"
	"github.com/dbchaincloud/cosmos-sdk/types/rest"
	"github.com/gorilla/mux"
	"github.com/mr-tron/base58"
	oerr "github.com/yzhanginwa/dbchain/x/dbchain/client/oracle/error"
	"github.com/yzhanginwa/dbchain/x/dbchain/client/oracle/oracle"
	"github.com/yzhanginwa/dbchain/x/dbchain/internal/types"
	"io/ioutil"
	"net/http"
	"strconv"
)

const (
	BaseUrl = oracle.BaseUrl + "dbchain/"
	ErrCode = "err_code"
	ErrInfo = "err_info"
	SuccessInfo = "success_info"
)

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
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write(res)
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
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write(res)
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
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write(res)
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
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write(res)
	}
}

func nftFindPopularAuthor(cliCtx context.CLIContext, storeName string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		numbers := vars["numbers"]
		inumber, err := strconv.Atoi(numbers)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusNotFound, "numbers err")
			generalResponse(w, map[string]string{
				ErrInfo : oerr.ErrDescription[oerr.ParamsErrCode],
				ErrCode : oerr.ParamsErrCode,
			})
			return
		}
		_ = numbers
		queryString := `[{"method":"table","table":"denom"},{"method":"select","fields":"user_id"}]`
		baseQueryString := base58.Encode([]byte(queryString))

		ac := getOracleAc()
		requestUrl := fmt.Sprintf("%s/querier/%s/%s/%s", BaseUrl,ac, nftAppCode,baseQueryString)
		res, err := httpGetRequest(requestUrl)
		if err != nil {
			generalResponse(w, map[string]string{
				ErrInfo : err.Error(),
				ErrCode : oerr.UndefinedErrCode,
			})
			return
		}

		type response struct {
			Height string
			Result []map[string]string
		}
		temp := response{}
		err = json.Unmarshal(res, &temp)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusNotFound, err.Error())
			return
		}

		Statistics := make(map[string]int, 0)
		for _, v := range temp.Result {
			userId := v["user_id"]
			num := Statistics[userId]
			num++
			Statistics[userId] = num
		}

		popularAuthor := make(map[int]string)
		count := 0
		min := 0
		for user, num := range Statistics {
			count++
			if min == 0 {
				min = num
			}

			if count <= inumber {
				popularAuthor[num] = user
			} else {
				if min < num {
					delete(popularAuthor, min)
					popularAuthor[num] = user
					min = num
				}
			}
		}
		//query popular
		type userInfo struct {
			Height string
			Result []map[string]string
		}
		tempUserInfo := userInfo{}
		popularAuthorsInfo := make([]map[string]string, 0)
		for _, userid := range popularAuthor {
			ac := getOracleAc()
			userInfo, err := findByCore(cliCtx, storeName, ac, nftAppCode, nftUserTable, "id", userid)
			if err != nil || userInfo == nil {
				continue
			}

			queryString := `[{"method":"table","table":"user_info"},{"method":"select","fields":"user_id,avatar,nickname"},{"method" : "where", "field" : "user_id", "operator" : "=", "value" : "` + userid + `"}]`
			baseQueryString := base58.Encode([]byte(queryString))

			requestUrl := fmt.Sprintf("%s/querier/%s/%s/%s", BaseUrl,ac, nftAppCode,baseQueryString)
			res, err := httpGetRequest(requestUrl)
			if err != nil {
				rest.WriteErrorResponse(w, http.StatusNotFound, err.Error())
				return
			}
			err = json.Unmarshal(res, &tempUserInfo)
			if err != nil {
				continue
			}
			if len(tempUserInfo.Result) > 0 {
				tempUserInfo.Result[0]["tel"] = userInfo["tel"]
				popularAuthorsInfo = append(popularAuthorsInfo, tempUserInfo.Result[0])
			} else {
				popularAuthorsInfo = append(popularAuthorsInfo, map[string]string{
					"user_id" : userid,
					"avatar" : "",
					"nickname" : "",
					"tel" : userInfo["tel"],
				})
			}
		}

		bz, _ := json.Marshal(popularAuthorsInfo)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write(bz)
	}
}

func nftFindLastestNft(cliCtx context.CLIContext, storeName string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		numbers := vars["numbers"]
		inumber, err := strconv.Atoi(numbers)
		if err != nil {
			generalResponse(w, map[string]string{
				ErrInfo : oerr.ErrDescription[oerr.ParamsErrCode],
				ErrCode : oerr.ParamsErrCode})
			return
		}
		queryString := `[{"method":"table","table":"nft_publish"},{"method":"select","fields":"id"}]`
		ids := queryByQuerier(queryString)
		if len(ids) == 0 {
			generalResponse(w, map[string]string{
				ErrInfo : oerr.ErrDescription[oerr.UndefinedErrCode],
				ErrCode : oerr.UndefinedErrCode,
			})
			return
		}
		var validIds []map[string]string
		if len(ids) > inumber {
			staid := len(ids) - inumber
			validIds = ids[staid : ]
		} else {
			validIds = ids
		}
		nfts := make([]map[string]string, 0)
		for _, validId := range validIds {
			id  := validId["id"]
			ac := getOracleAc()
			queryString := fmt.Sprintf("%s/find/%s/%s/%s/%s", BaseUrl, ac, nftAppCode, nftPublishTable, id)
			publishInfo, err := findRow(cliCtx, queryString)
			if err != nil {
				continue
			}
			denomId := publishInfo["denom_id"]
			queryString = fmt.Sprintf("%s/find/%s/%s/%s/%s", BaseUrl, ac, nftAppCode, denomTable, denomId)
			ntfInfo, err := findRow(cliCtx, queryString)
			if err != nil {
				continue
			}
			userInfo , err := findByCore(cliCtx, storeName, ac, nftAppCode, nftUserInfoTable, "user_id", ntfInfo["user_id"])
			if err != nil {
				continue
			}
			ntfInfo["avatar"] = userInfo["avatar"]
			ntfInfo["nickname"] = userInfo["nickname"]
			nfts = append(nfts, ntfInfo)
		}
		bz, _ := json.Marshal(nfts)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write(bz)
	}
}

func nftFindNftDetails(cliCtx context.CLIContext, storeName string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		denomId := vars["denom_id"]
		ac := getOracleAc()
		queryString := fmt.Sprintf("%s/find/%s/%s/%s/%s", BaseUrl, ac, nftAppCode, denomTable, denomId)
		ntfInfo, err := findRow(cliCtx, queryString)
		if err != nil {
			generalResponse(w, map[string]string{
				ErrInfo : "find denom err",
				ErrCode : oerr.UndefinedErrCode,
			})
			return
		}

		userId := ntfInfo["user_id"]
		userInfo , err := findByCore(cliCtx, storeName, ac, nftAppCode, nftUserInfoTable, "user_id", userId)
		if err != nil || userInfo == nil {
			ntfInfo["avatar"] =  ""
			ntfInfo["nickname"] = ""
		} else {
			ntfInfo["avatar"] =  userInfo["avatar"]
			ntfInfo["nickname"] = userInfo["nickname"]
		}

		publishInfo , err := findByCore(cliCtx, storeName, ac, nftAppCode, nftPublishTable, "denom_id", denomId)
		if err != nil {
			generalResponse(w, map[string]string{
				ErrInfo : "find price err",
				ErrCode : oerr.UndefinedErrCode,
			})
			return
		}
		ntfInfo["price"] = publishInfo["price"]

		nfts , err := findByCoreIds(cliCtx, storeName, ac, nftAppCode, nftTable, "denom_id", denomId)
		if err != nil {
			generalResponse(w, map[string]string{
				ErrInfo : "find nft err",
				ErrCode : oerr.UndefinedErrCode,
			})
			return
		}
		ntfInfo["remaining"] = strconv.Itoa(len(nfts))
		bz, _ := json.Marshal(ntfInfo)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write(bz)
	}
}


func nftUserInfo(cliCtx context.CLIContext, storeName string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		tel := vars["tel"]
		_, ok := verifySession(w, r, tel)
		if !ok {
			generalResponse(w, map[string]string{
				ErrInfo : oerr.ErrDescription[oerr.UnLoginErrCode],
				ErrCode : oerr.UnLoginErrCode})
			return
		}
		ac := getOracleAc()
		res, err := findByCore(cliCtx, storeName, ac, nftAppCode, nftUserTable, "tel", tel)
		if err != nil || res == nil{
			generalResponse(w, map[string]string{
				ErrInfo : "find user info err",
				ErrCode : oerr.UndefinedErrCode,
			})
			return
		}
		result := map[string]string {
			"tel" : res["tel"],
			"address" : res["address"],
		}

		userid := res["id"]
		res, err = findByCore(cliCtx, storeName, ac, nftAppCode, nftUserInfoTable, "user_id", userid)
		if err != nil {
			generalResponse(w, map[string]string{
				ErrInfo : "find user info err",
				ErrCode : oerr.UndefinedErrCode,
			})
			return
		}
		if len(res) != 0 {
			result["avatar"] =       res["avatar"]
			result["nickname"] =     res["nickname"]
			result["description"] =  res["description"]
		}
		bz, _ := json.Marshal(result)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write(bz)
	}
}


func nftsOfUserMake(cliCtx context.CLIContext, storeName string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		tel := vars["tel"]
		_, ok := verifySession(w, r, tel)
		if !ok {
			generalResponse(w, map[string]string{
				ErrInfo : oerr.ErrDescription[oerr.UnLoginErrCode],
				ErrCode : oerr.UnLoginErrCode},
			)
			return
		}
		ac := getOracleAc()
		res, err := findByCore(cliCtx, storeName, ac, nftAppCode, nftUserTable, "tel", tel)
		if err != nil || res == nil {
			generalResponse(w, map[string]string{
				ErrInfo : "find user info err",
				ErrCode : oerr.UndefinedErrCode,
			})
			return
		}
		userid := res["id"]
		denoms, err := findByAll(cliCtx, storeName, ac, nftAppCode, denomTable, "user_id", userid)
		if err != nil {
			generalResponse(w, map[string]string{
				ErrInfo : "find nft err",
				ErrCode : oerr.UndefinedErrCode,
			})
			return
		}
		for _, denom := range denoms {
			denomId := denom["id"]
			ac := getOracleAc()
			publishId, err := findByCoreIds(cliCtx, storeName, ac, nftAppCode, nftPublishTable, "denom_id", denomId)
			if err != nil {
				continue
			}
			if len(publishId) != 0 {
				denom["publish"] = "true"
			} else {
				denom["publish"] = "false"
			}
		}
		bz, _ := json.Marshal(denoms)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write(bz)
	}
}


func nftsOfUserBuy(cliCtx context.CLIContext, storeName string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		tel := vars["tel"]
		_, ok := verifySession(w, r, tel)
		if !ok {
			generalResponse(w, map[string]string{
				ErrInfo : oerr.ErrDescription[oerr.UnLoginErrCode],
				ErrCode : oerr.UnLoginErrCode},
			)
			return
		}
		ac := getOracleAc()
		res, err := findByCore(cliCtx, storeName, ac, nftAppCode, nftUserTable, "tel", tel)
		if err != nil || res == nil {
			generalResponse(w, map[string]string{
				ErrInfo : "find user info err",
				ErrCode : oerr.UndefinedErrCode,
			})
			return
		}
		addr := res["address"]
		boughts, err := findByAll(cliCtx, storeName, ac, nftAppCode, nftCardBagTable, "owner", addr)
		if err != nil {
			generalResponse(w, map[string]string{
				ErrInfo : "find nft err",
				ErrCode : oerr.UndefinedErrCode,
			})
			return
		}
		temp := make(map[string]struct{})
		result := make([]map[string]string, 0)
		for _, bought := range boughts {
			ac := getOracleAc()
			nftId := bought["nft_id"]

			if _, ok := temp[nftId]; ok {
				continue
			} else {
				temp[nftId] = struct{}{}
			}
			id := bought["id"]
			lastNftTransferInfo, err := findByCore(cliCtx, storeName, ac, nftAppCode, nftCardBagTable, "nft_id", nftId)
			if err != nil || lastNftTransferInfo == nil {
				continue
			}
			if id == lastNftTransferInfo["id"] {
				nftINfo, err := findNftInfoByNftId(cliCtx, storeName, nftId)
				if err != nil {
					continue
				}
				result = append(result, nftINfo)
			}
		}
		bz, _ := json.Marshal(result)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write(bz)
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
	if len(out) == 0 {
		return nil, nil
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


func queryByQuerier(queryString string) []map[string]string {
	baseQueryString := base58.Encode([]byte(queryString))
	ac := getOracleAc()
	requestUrl := fmt.Sprintf("%s/querier/%s/%s/%s", BaseUrl,ac, nftAppCode,baseQueryString)
	res, err := httpGetRequest(requestUrl)
	if err != nil {
		return nil
	}

	//query popular
	type userInfo struct {
		Height string
		Result []map[string]string
	}
	tempUserInfo := userInfo{}
	err = json.Unmarshal(res, &tempUserInfo)
	if err != nil {
		return nil
	}
	return tempUserInfo.Result
}

func findNftInfoByNftId(cliCtx context.CLIContext, storeName, nftId string) (map[string]string, error) {
	ac := getOracleAc()

	queryString := fmt.Sprintf("%s/find/%s/%s/%s/%s", BaseUrl, ac, nftAppCode, nftTable, nftId)
	nftInfo, err := findRow(cliCtx, queryString)
	if err !=nil || nftInfo == nil {
		return nftInfo, err
	}

	denomInfo, err := findByCore(cliCtx, storeName, ac, nftAppCode, denomTable, "id", nftInfo["denom_id"])
	if err != nil || denomInfo == nil {
		return nftInfo, err
	}
	nftInfo["name"] = denomInfo["name"]
	nftInfo["file"] = denomInfo["file"]
	nftInfo["description"] = denomInfo["description"]
	return nftInfo, nil
}