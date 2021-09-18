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
	"sort"
	"strconv"
	"time"
)

const (
	ErrCode = "err_code"
	ErrInfo = "err_info"
	Result = "result"
)
var BaseUrl = oracle.BaseUrl + "dbchain/"


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

		nums := make([]int, 0)
		numAndUser := make(map[int]string)
		for user, num := range Statistics {
			if _, ok := numAndUser[num]; !ok {
				numAndUser[num] = user
				nums = append(nums, num)
			}
		}
		sort.Ints(nums)
		popularAuthor := make([]string, 0)
		if len(nums) > inumber {
			nums = nums[len(nums) - inumber : ]
		}
		for i := len(nums) - 1; i >= 0; i-- {
			num := nums[i]
			popularAuthor = append(popularAuthor, numAndUser[num])
		}

		//query popular
		type userInfoStruct struct {
			Height string
			Result []map[string]string
		}

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
			tempUserInfo := userInfoStruct{}
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

		successDataResponse(w, popularAuthorsInfo)
		return
	}
}

func nftFindLastestNft(cliCtx context.CLIContext, storeName string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		page := vars["page"]
		numbers := vars["numbers"]
		inumber, err := strconv.Atoi(numbers)
		if err != nil || inumber <= 0 {
			generalResponse(w, map[string]string{
				ErrInfo : oerr.ErrDescription[oerr.ParamsErrCode],
				ErrCode : oerr.ParamsErrCode})
			return
		}
		ipage, err := strconv.Atoi(page)
		if err != nil || ipage <= 0 {
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
		if len(ids) >= ipage * inumber {
			staid := len(ids) - ipage * inumber
			end := len(ids) - inumber * (ipage - 1)
			validIds = ids[staid : end]
		} else if (len(ids) < ipage * inumber) && (len(ids) >= inumber * (ipage - 1)){
			end := len(ids) - inumber * (ipage - 1)
			validIds = ids[ : end]
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
			ntfInfo["price"] = publishInfo["price"]
			nfts = append(nfts, ntfInfo)
		}

		successDataResponse(w, nfts)
		return
	}
}

func nftFindNftDetails(cliCtx context.CLIContext, storeName string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		denomId := vars["denom_id"]
		ac := getOracleAc()
		queryString := fmt.Sprintf("%s/find/%s/%s/%s/%s", BaseUrl, ac, nftAppCode, denomTable, denomId)
		ntfInfo, err := findRow(cliCtx, queryString)
		if err != nil || ntfInfo == nil {
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

		if len(publishInfo) != 0 {
			ntfInfo["publish"] = "true"
		} else {
			ntfInfo["publish"] = "false"
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
		ntfInfo["collected"] = "false"
		//
		queryUserId, _ := verifySession(w, r)
		if queryUserId != "" {
			ac = getOracleAc()
			queryString := `[{"method":"table","table":"nft_collection"},{"method":"select","fields":"user_id,denom_id"},{"method":"where", "field" : "user_id", "operator" : "=", "value" : "` + queryUserId + `"},{"method":"where", "field" : "denom_id", "operator" : "=", "value" : "` + denomId + `"}]`
			collect := queryByQuerier(queryString)
			if len(collect) != 0 {
				ntfInfo["collected"] = "true"
			}
		}
		//
		collectNum , err := findByCoreIds(cliCtx, storeName, ac, nftAppCode, nftCollection, "denom_id", denomId)
		ntfInfo["collected_number"] = strconv.Itoa(len(collectNum))

		successDataResponse(w, ntfInfo)
		return
	}
}


func nftUserInfo(cliCtx context.CLIContext, storeName string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userId, ok := verifySession(w, r)
		if !ok {
			generalResponse(w, map[string]string{
				ErrInfo : oerr.ErrDescription[oerr.UnLoginErrCode],
				ErrCode : oerr.UnLoginErrCode})
			return
		}
		ac := getOracleAc()
		queryString := fmt.Sprintf("%s/find/%s/%s/%s/%s", BaseUrl, ac, nftAppCode, nftUserTable, userId)
		res, err := findRow(cliCtx, queryString)
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
			"my_code" : res["my_code"],
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
		} else {
			result["avatar"] =       ""
			result["nickname"] =     ""
			result["description"] =  ""
		}

		result["authentication"] = "false"
		if userAuthenticationStatus(userId) {
			result["authentication"] = "true"
		}

		res, err = findByCore(cliCtx, storeName, ac, nftAppCode, nftScoreTable, "user_id", userid)
		if err != nil {
			generalResponse(w, map[string]string{
				ErrInfo : "find user info err",
				ErrCode : oerr.UndefinedErrCode,
			})
			return
		}
		if len(res) != 0 {
			result["token"] = res["token"]
		} else {
			result["token"] = "0"
		}

		successDataResponse(w, result)
		return
	}
}

func nftOfUserBasicInfo(cliCtx context.CLIContext, storeName string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		userId := vars["user_id"]
		ac := getOracleAc()
		queryString := fmt.Sprintf("%s/find/%s/%s/%s/%s", BaseUrl, ac, nftAppCode, nftUserTable, userId)
		res, err := findRow(cliCtx, queryString)
		if err != nil || res == nil{
			generalResponse(w, map[string]string{
				ErrInfo : "find user info err",
				ErrCode : oerr.UndefinedErrCode,
			})
			return
		}
		result := map[string]string {
			"address" : res["address"],
		}

		res, err = findByCore(cliCtx, storeName, ac, nftAppCode, nftUserInfoTable, "user_id", userId)
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

		successDataResponse(w, result)
		return
	}
}

func nftUserIncome(cliCtx context.CLIContext, storeName string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		userId := vars["user_id"]
		ac := getOracleAc()
		queryString := fmt.Sprintf("%s/find/%s/%s/%s/%s", BaseUrl, ac, nftAppCode, nftUserTable, userId)
		userInfo, err := findRow(cliCtx, queryString)
		if err != nil || userInfo == nil{
			generalResponse(w, map[string]string{
				ErrInfo : "find user info err",
				ErrCode : oerr.UndefinedErrCode,
			})
			return
		}
		seller := userInfo["address"]

		receipts, err := findByAll(cliCtx, storeName, ac, nftAppCode, nftOrderReceipt, "seller", seller)
		if err != nil {
			generalResponse(w, map[string]string{
				ErrInfo : "find user info err",
				ErrCode : oerr.UndefinedErrCode,
			})
			return
		}

		money := 0.0
		for _, receipt := range receipts {
			amount := receipt["amount"]
			fAmount, err := strconv.ParseFloat(amount, 64)
			if err != nil {
				continue
			}
			money += fAmount
		}
		result := map[string]string{
			"money" : fmt.Sprintf("%f", money),
		}
		successDataResponse(w, result)
		return
	}
}

func nftUserIncomeByTime(cliCtx context.CLIContext, storeName string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		userId := vars["user_id"]
		startTimeStr := vars["start_time"]
		startTime, err := time.Parse("2006-01-02", startTimeStr)
		if err != nil {
			generalResponse(w, map[string]string{
				ErrInfo : oerr.ErrDescription[oerr.ParamsErrCode],
				ErrCode : oerr.ParamsErrCode,
			})
			return
		}
		endTimeStr := vars["end_time"]
		endTime, err := time.Parse("2006-01-02", endTimeStr)
		if err != nil {
			generalResponse(w, map[string]string{
				ErrInfo : oerr.ErrDescription[oerr.ParamsErrCode],
				ErrCode : oerr.ParamsErrCode,
			})
			return
		}
		startStamp := startTime.UnixNano() / 1000000
		endStamp := endTime.UnixNano() / 1000000


		ac := getOracleAc()
		queryString := fmt.Sprintf("%s/find/%s/%s/%s/%s", BaseUrl, ac, nftAppCode, nftUserTable, userId)
		userInfo, err := findRow(cliCtx, queryString)
		if err != nil || userInfo == nil{
			generalResponse(w, map[string]string{
				ErrInfo : "find user info err",
				ErrCode : oerr.UndefinedErrCode,
			})
			return
		}
		seller := userInfo["address"]

		receipts, err := findByAll(cliCtx, storeName, ac, nftAppCode, nftOrderReceipt, "seller", seller)
		if err != nil {
			generalResponse(w, map[string]string{
				ErrInfo : "find user info err",
				ErrCode : oerr.UndefinedErrCode,
			})
			return
		}

		money := 0.0
		for _, receipt := range receipts {
			createdTime := receipt["created_at"]
			createStamp , _ := strconv.ParseInt(createdTime, 10, 64)
			if createStamp < startStamp || createStamp > endStamp {
				continue
			}

			amount := receipt["amount"]
			fAmount, err := strconv.ParseFloat(amount, 64)
			if err != nil {
				continue
			}
			money += fAmount
		}
		result := map[string]string{
			"money" : fmt.Sprintf("%f", money),
		}
		successDataResponse(w, result)
		return
	}
}

func nftUserCollectInfo(cliCtx context.CLIContext, storeName string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		numOrDetail := vars["number_or_detail"]
		userId := vars["user_id"]

		queryString := `[{"method":"table","table":"nft_collection"},{"method":"select","fields":"denom_id"},{"method":"where", "field" : "user_id", "operator" : "=", "value" : "` + userId + `"}]`
		collects := queryByQuerier(queryString)
		if numOrDetail == "number" {
			number := strconv.Itoa(len(collects))
			result := map[string]string {
				"number" : number,
			}
			successDataResponse(w, result)
			return
		}

		result := make([]map[string]string, 0)
		for _, collect := range collects {
			denomId := collect["denom_id"]
			ac := getOracleAc()
			queryString := fmt.Sprintf("%s/find/%s/%s/%s/%s", BaseUrl, ac, nftAppCode, denomTable, denomId)
			denomInfo, err := findRow(cliCtx, queryString)
			if err != nil {
				continue
			}
			if len(denomInfo) != 0 {
				AuthorInfo := findAuthorInfoFromDenomId(cliCtx, storeName, denomId)
				denomInfo["avatar"] =  AuthorInfo["avatar"]
				denomInfo["nickname"] = AuthorInfo["nickname"]
				result = append(result, denomInfo)
			}
		}
		successDataResponse(w, result)
		return

	}
}

func nftUserAllTokenRecord(cliCtx context.CLIContext, storeName string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userId, ok := verifySession(w, r)
		if !ok {
			generalResponse(w, map[string]string{
				ErrInfo : oerr.ErrDescription[oerr.UnLoginErrCode],
				ErrCode : oerr.UnLoginErrCode})
			return
		}
		queryString := `[{"method":"table","table":"score"},{"method":"select","fields":"token,action,memo,increment,created_at"},{"method" : "where", "field" : "user_id", "operator" : "=", "value" : "` + userId + `"}]`
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

		type tokenRecord struct {
			Height string
			Result []map[string]string
		}
		temp := tokenRecord{}
		json.Unmarshal(res, &temp)

		successDataResponse(w, temp.Result)
		return
	}
}

func nftUserInvitationRecord(cliCtx context.CLIContext, storeName string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		userId, ok := verifySession(w, r)
		if !ok {
			generalResponse(w, map[string]string{
				ErrInfo : oerr.ErrDescription[oerr.UnLoginErrCode],
				ErrCode : oerr.UnLoginErrCode})
			return
		}
		ac := getOracleAc()
		queryString := fmt.Sprintf("%s/find/%s/%s/%s/%s", BaseUrl, ac, nftAppCode, nftUserTable, userId)
		userInfo, err := findRow(cliCtx, queryString)
		if err != nil || userInfo == nil {
			generalResponse(w, map[string]string{
				ErrInfo : "find denom err",
				ErrCode : oerr.UndefinedErrCode,
			})
			return
		}
		myCode := userInfo["my_code"]
		queryString = `[{"method":"table","table":"user"},{"method":"select","fields":"id,tel,created_at"},{"method" : "where", "field" : "invitation_code", "operator" : "=", "value" : "` + myCode + `"}]`
		inviteInfos := queryByQuerier(queryString)
		for _ ,inviteInfo := range inviteInfos {
			inviteInfo["action"] = "+" + strconv.Itoa(invitationScore)
			tel := inviteInfo["tel"]
			tel = tel[:3] + "****" + tel[7:]
			inviteInfo["tel"] = tel
			queryString = `[{"method":"table","table":"user_info"},{"method":"select","fields":"nickname"},{"method" : "where", "field" : "user_id", "operator" : "=", "value" : "` + inviteInfo["id"] + `"}]`
			inviteeInfo := queryByQuerier(queryString)
			if len(inviteeInfo) > 0 {
				inviteInfo["nickname"] = inviteeInfo[0]["nickname"]
			} else {
				inviteInfo["nickname"] = ""
			}
			delete(inviteInfo, "id")
		}

		successDataResponse(w, inviteInfos)
		return
	}}

func nftUserNftOrderNumber(cliCtx context.CLIContext, storeName string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userId, ok := verifySession(w, r)
		if !ok {
			generalResponse(w, map[string]string{
				ErrInfo : oerr.ErrDescription[oerr.UnLoginErrCode],
				ErrCode : oerr.UnLoginErrCode,
			})
			return
		}

		ac := getOracleAc()
		queryString := fmt.Sprintf("%s/find/%s/%s/%s/%s", BaseUrl, ac, nftAppCode, nftUserTable, userId)
		userInfo, err :=  findRow(cliCtx, queryString)
		if err != nil || len(userInfo) == 0 {
			generalResponse(w, map[string]string{
				ErrInfo : "find user info err",
				ErrCode : oerr.UndefinedErrCode,
			})
			return
		}

		tel := userInfo["tel"]
		orderIds, err := findByCoreIds(cliCtx, storeName, ac, nftAppCode, nftPublishOrder, "tel", tel)
		if err != nil {
			generalResponse(w, map[string]string{
				ErrInfo : "find user info err",
				ErrCode : oerr.UndefinedErrCode,
			})
			return
		}
		numOfPayOrder := 0
		for _, orderId := range orderIds {
			ac := getOracleAc()
			val := nftAppCode + "_buy_" + orderId
			payed, _ := findByCoreIds(cliCtx, storeName, ac, nftAppCode, nftOrderReceipt, "orderid", val)
			if len(payed) != 0 {
				numOfPayOrder++
			}
		}
		result := map[string]string{
			"nft_order_numbers" : strconv.Itoa(numOfPayOrder),
		}

		successDataResponse(w, result)
		return
	}
}

func nftUserNftNumber(cliCtx context.CLIContext, storeName string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		userId := vars["user_id"]

		ac := getOracleAc()
		denomIds, err := findByCoreIds(cliCtx, storeName, ac, nftAppCode, denomTable, "user_id", userId)
		if err != nil {
			generalResponse(w, map[string]string{
				ErrInfo : "find user info err",
				ErrCode : oerr.UndefinedErrCode,
			})
			return
		}
		result := map[string]string{
			"nft_numbers" : strconv.Itoa(len(denomIds)),
		}

		successDataResponse(w, result)
		return
	}
}

func nftsOfUserMake(cliCtx context.CLIContext, storeName string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		userid := vars["user_id"]
		publishStatus := vars["publish_status"]
		if publishStatus != "all" && publishStatus != "published" {
			generalResponse(w, map[string]string{
				ErrInfo : oerr.ErrDescription[oerr.ParamsErrCode],
				ErrCode : oerr.ParamsErrCode},
			)
			return
		}
		if publishStatus == "all" {
			_, ok := verifySession(w, r)
			if !ok {
				generalResponse(w, map[string]string{
					ErrInfo : oerr.ErrDescription[oerr.UnLoginErrCode],
					ErrCode : oerr.UnLoginErrCode},
				)
				return
			}
		}
		ac := getOracleAc()
		denoms, err := findByAll(cliCtx, storeName, ac, nftAppCode, denomTable, "user_id", userid)
		if err != nil {
			generalResponse(w, map[string]string{
				ErrInfo : "find nft err",
				ErrCode : oerr.UndefinedErrCode,
			})
			return
		}
		result := make([]map[string]string, 0)
		for _, denom := range denoms {
			denomId := denom["id"]
			ac := getOracleAc()
			publish, err := findByCore(cliCtx, storeName, ac, nftAppCode, nftPublishTable, "denom_id", denomId)
			if err != nil {
				continue
			}

			if publishStatus == "published" {
				if len(publish) == 0 {
					continue
				}
			} else {
				if len(publish) != 0 {
					denom["publish"] = "true"
				} else {
					denom["publish"] = "false"
				}
			}
			denom["price"] = publish["price"]
			result = append(result, denom)

		}
		successDataResponse(w, result)
		return
	}
}


func nftsOfUserBuy(cliCtx context.CLIContext, storeName string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userId, ok := verifySession(w, r)
		if !ok {
			generalResponse(w, map[string]string{
				ErrInfo : oerr.ErrDescription[oerr.UnLoginErrCode],
				ErrCode : oerr.UnLoginErrCode},
			)
			return
		}
		ac := getOracleAc()
		queryString := fmt.Sprintf("%s/find/%s/%s/%s/%s", BaseUrl, ac, nftAppCode, nftUserTable, userId)
		res, err := findRow(cliCtx, queryString)
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
				userInfo := findAuthorInfoByNftId(cliCtx, storeName, nftId)
				nftINfo["avatar"] = userInfo["avatar"]
				nftINfo["nickname"] = userInfo["nickname"]
				result = append(result, nftINfo)
			}
		}
		successDataResponse(w, result)
		return
	}
}


//////////////////////////
//                      //
//      help func       //
//                      //
//////////////////////////

func userAuthenticationStatus(userId string) bool {
	queryString := `[{"method":"table","table":"nft_real_name_authentication"},{"method":"select","fields":"user_id"},{"method":"where", "field" : "user_id", "operator" : "=", "value" : "` + userId + `"}]`
	userAuthenticationInfo := queryByQuerier(queryString)
	if len(userAuthenticationInfo) > 0 {
		return true
	}
	return false
}

func userProductionPermission(userId string) bool {
	queryString := `[{"method":"table","table":"nft_production_permission"},{"method":"select","fields":"user_id"},{"method":"where", "field" : "user_id", "operator" : "=", "value" : "` + userId + `"}]`
	permission := queryByQuerier(queryString)
	if len(permission) > 0 {
		return true
	}
	return false
}

func findAuthorInfoByNftId(cliCtx context.CLIContext, storeName, nftId string) map[string]string {
	ac := getOracleAc()
	queryString := fmt.Sprintf("%s/find/%s/%s/%s/%s", BaseUrl, ac, nftAppCode, nftTable, nftId)
	nftInfo, err := findRow(cliCtx, queryString)
	if err != nil {
		return map[string]string{}
	}
	userInfo := findAuthorInfoFromDenomId(cliCtx, storeName, nftInfo["denom_id"])
	return userInfo
}

func findAuthorInfoFromDenomId(cliCtx context.CLIContext, storeName, denomId string) map[string]string{
	ac := getOracleAc()
	queryString := fmt.Sprintf("%s/find/%s/%s/%s/%s", BaseUrl, ac, nftAppCode, denomTable, denomId)
	denomInfo, err := findRow(cliCtx, queryString)
	if err != nil {
		return map[string]string{}
	}
	ac = getOracleAc()
	userInfo, err :=findByCore(cliCtx, storeName, ac, nftAppCode, nftUserInfoTable, "user_id", denomInfo["user_id"])
	if err != nil {
		return map[string]string{}
	}
	return userInfo
}

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