package oracle

import (
	"bytes"
	stdCtx "context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/dbchaincloud/cosmos-sdk/client/context"
	sdk "github.com/dbchaincloud/cosmos-sdk/types"
	"github.com/dbchaincloud/tendermint/crypto"
	"github.com/dbchaincloud/tendermint/crypto/sm2"
	"github.com/go-session/session"
	"github.com/smartwalle/alipay/v3"
	"github.com/spf13/viper"
	"github.com/yzhanginwa/dbchain/x/dbchain/client/oracle/cache"
	oerr "github.com/yzhanginwa/dbchain/x/dbchain/client/oracle/error"
	"github.com/yzhanginwa/dbchain/x/dbchain/client/oracle/oracle"
	"github.com/yzhanginwa/dbchain/x/dbchain/internal/types"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"
)

const (
	//nftAppCode = "9SXMMWWR8A"
	nftAppCodeKey = "app-code"
	codePre = "dbc"
	//tables
	nftUserTable = "user"
	nftPasswordTable = "password"
	nftUserInfoTable = "user_info"
	nftScoreTable = "score"
	nftTable = "nft"
	denomTable = "denom"
	nftPublishTable = "nft_publish"
	nftMarketTable = "nft_market"
	nftCardBagTable = "nft_card_bag"
	nftPublishOrder = "nft_publish_order"
	nftOrderReceipt = "nft_order_receipt"
	nftMakeOrder = "nft_make_order"
	nftCollection = "nft_collection"
	nftRealNameAuthentication = "nft_real_name_authentication"
	//only user in this table can make nft
	nftProductionPermission = "nft_production_permission"

	priceRegExp = `(^[1-9]\d*(\.\d{1,2})?$)|(^0(\.\d{1,2})?$)`
	invitationScore = 1000
	nftMakePricePerNft = 0.01
)

var priceRex *regexp.Regexp
var orderSet *nftOrderSet
var makeOrderCache *cache.MemoryCache
//var nftAppCode = "3CASSKY7NQ"
var NFTAppCode = ""

func init() {
	// order cache
	orderSet = newNftOrderSet(time.Second * 300)
	go orderSet.GC()
	makeOrderCache = cache.NewMemoryCache(time.NewTicker(30 * time.Minute), 1800)
	go makeOrderCache.Gc()

	//session
	session.InitManager(
		session.SetCookieLifeTime(86400),
		session.SetEnableSIDInHTTPHeader(true),
	)
}

func nftUserRegister(cliCtx context.CLIContext, storeName string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data, err := readBodyData(r)
		if err != nil {
			generalResponse(w, map[string]string{
				ErrInfo : oerr.ErrDescription[oerr.UndefinedErrCode],
				ErrCode : oerr.UndefinedErrCode})
			return

		}
		verificationCode := data["verification_code"]
		_ = verificationCode
		/* user table
		tel code address invitation code
		input params
		tel password invitation_code
		*/
		//TODO check if register
		tel := data["tel"]
		ac := getOracleAc()
		nftAppCode := LoadNFTAppCode()
		res , _ := findByCore(cliCtx, storeName, ac, nftAppCode, nftUserTable, "tel", tel)
		if len(res) != 0 {
			generalResponse(w, map[string]string{
				ErrInfo : oerr.ErrDescription[oerr.RegisteredErrCode],
				ErrCode : oerr.RegisteredErrCode})
			return
		}
		if !VerifyVerfCode(tel, tel, verificationCode) {
			generalResponse(w, map[string]string{
				ErrInfo : oerr.ErrDescription[oerr.TelNoVerifyCode],
				ErrCode : oerr.TelNoVerifyCode})
			return
		}

		password := data["password"]//
		password, valid := passwordFormatCheck(password)
		if !valid {
			generalResponse(w, map[string]string{
				ErrInfo : oerr.ErrDescription[oerr.FormatErrCode],
				ErrCode : oerr.FormatErrCode,
			})
			return
		}
		myCode := genCode(4, cliCtx, storeName, nftUserTable, "my_code")
		if myCode == "" {
			generalResponse(w, map[string]string{
				ErrInfo : "gen invitation code fail",
				ErrCode : oerr.UndefinedErrCode,
			})
			return
		}
		invitationCode := data["invitation_code"]
		fieldsOfUser := map[string]string {
			"tel" : tel,
			"my_code" : myCode,
			"address" : genUserAddress(),
			"invitation_code" : invitationCode,
		}
		//只管提交交易，前端负责查询
		fieldsOfPassword := map[string]string{
			"user_id" : "",
			"password" : password,
		}
		strFieldsOfUser, _ := json.Marshal(fieldsOfUser)
		strFieldsOfPassword, _ := json.Marshal(fieldsOfPassword)
		argument := []string{nftUserTable,string(strFieldsOfUser), nftPasswordTable ,string(strFieldsOfPassword)}
		strArgument, _ := json.Marshal(argument)

		err = callFunction(cliCtx,nftAppCode, "register", string(strArgument))
		if err != nil {
			generalResponse(w, map[string]string{
				ErrInfo : "register fail",
				ErrCode : oerr.UndefinedErr,
			})
			return
		}
		if invitationCode != "" {
			updateScoreOfInvite(cliCtx, storeName, invitationCode, "+", invitationScore, "Invite users")
		}
		generalResponse(w, map[string]string{
			ErrInfo : oerr.ErrDescription[oerr.SuccessCode],
			ErrCode : oerr.SuccessCode,
		})
		return
	}
}

func updateScoreOfInvite(cliCtx context.CLIContext, storeName string, invitationCode, action string, increment int, memo string) error {
	ac := getOracleAc()
	nftAppCode := LoadNFTAppCode()
	//find userId
	userId, err := findByCoreIds(cliCtx, storeName, ac, nftAppCode, nftUserTable, "my_code", invitationCode)
	if err != nil || len(userId) == 0 {
		return err
	}
	return updateScoreCore(cliCtx, storeName, userId[0], action, increment, memo)
}

func updateScoreCore(cliCtx context.CLIContext, storeName string, userId, action string, increment int, memo string) error {
	//find user score
	ac := getOracleAc()
	nftAppCode := LoadNFTAppCode()
	score, err := findByCore(cliCtx, storeName, ac, nftAppCode, nftScoreTable, "user_id", userId)
	if err != nil {
		return err
	}
	owner := loadOracleAddr()
	msgs := make([]oracle.UniversalMsg, 0)
	lastToken := 0

	if len(score) != 0 {
		lastToken, _ = strconv.Atoi(score["token"])
	}
	currentToken := 0
	if action == "+" {
		currentToken = lastToken + increment
	} else {
		if lastToken < increment {
			return errors.New("not enough token")
		}
		currentToken = lastToken - increment
	}

	token := strconv.Itoa(currentToken)
	field , _ := json.Marshal(map[string]string{
		"user_id" : userId,
		"token" : token,
		"action" : action + strconv.Itoa(increment),
		"memo" : memo})
	msg := types.NewMsgInsertRow(owner, nftAppCode, nftScoreTable, field)
	if msg.ValidateBasic() != nil {
		return errors.New("internal err")
	}
	msgs = append(msgs, msg)
	oracle.BuildTxsAndBroadcast(cliCtx, msgs)
	return nil
}

func nftUserLogin(cliCtx context.CLIContext, storeName string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data, err := readBodyData(r)
		if err != nil {
			generalResponse(w, map[string]string{
				ErrInfo : oerr.ErrDescription[oerr.ParamsErrCode],
				ErrCode : oerr.ParamsErrCode,
			})
		}

		tel := data["tel"]
		password := data["password"]
		//query tel and password
		ac := getOracleAc()
		nftAppCode := LoadNFTAppCode()
		res, err := findByCore(cliCtx, storeName, ac, nftAppCode, "user", "tel", tel)
		if err != nil {
			generalResponse(w, map[string]string{
				ErrInfo : err.Error(),
				ErrCode : oerr.UndefinedErrCode,
			})
			return
		}

		if len(res) == 0 {
			generalResponse(w, map[string]string{
				ErrInfo : oerr.ErrDescription[oerr.UnregisterErrCode],
				ErrCode : oerr.UnregisterErrCode,
			})
			return
		}

		userId := res["id"]
		ac = getOracleAc()
		res, err = findByCore(cliCtx, storeName, ac, nftAppCode, "password", "user_id", userId)
		if err != nil {
			generalResponse(w, map[string]string{
				ErrInfo : err.Error(),
				ErrCode : oerr.UndefinedErrCode,
			})
			return
		}

		if len(res) == 0 {
			generalResponse(w, map[string]string{
				ErrInfo : oerr.ErrDescription[oerr.UndefinedErrCode],
				ErrCode : oerr.UndefinedErrCode,
			})
			return
		}

		hspswd := res["password"]
		hs := sha256.Sum256([]byte(password))
		if hex.EncodeToString(hs[:]) == hspswd {
			if !saveSession(w, r , tel, userId) {
				generalResponse(w, map[string]string{
					ErrInfo : "save session err",
					ErrCode : oerr.UndefinedErrCode,
				})
				return
			}
			result := map[string]string{ "user_id" : userId}
			successDataResponse(w, result)
			return
		} else {
			generalResponse(w, map[string]string{
				ErrInfo : oerr.ErrDescription[oerr.PasswordErrCode],
				ErrCode : oerr.PasswordErrCode,
			})
			return
		}
	}
}

func nftUserLogout(cliCtx context.CLIContext, storeName string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		if deleteSession(w, r) {
			generalResponse(w, map[string]string{
				ErrInfo : oerr.ErrDescription[oerr.SuccessCode],
				ErrCode : oerr.SuccessCode,
			})
			return
		}
		generalResponse(w, map[string]string{
			ErrInfo : oerr.ErrDescription[oerr.ServerErrCode],
			ErrCode : oerr.ServerErrCode,
		})
		return
	}
}

func nftUserResetPassword(cliCtx context.CLIContext, storeName string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data, err := readBodyData(r)
		if err != nil {
			generalResponse(w, map[string]string{
				ErrInfo : oerr.ErrDescription[oerr.ParamsErrCode],
				ErrCode : oerr.ParamsErrCode,
			})
		}

		tel := data["tel"]
		verificationCode := data["verification_code"]
		password := data["password"]
		password, valid := passwordFormatCheck(password)
		if !valid {
			generalResponse(w, map[string]string{
				ErrInfo : oerr.ErrDescription[oerr.FormatErrCode],
				ErrCode : oerr.FormatErrCode,
			})
			return
		}
		//query tel and password
		ac := getOracleAc()
		nftAppCode := LoadNFTAppCode()
		res, err := findByCore(cliCtx, storeName, ac, nftAppCode, "user", "tel", tel)
		if err != nil {
			generalResponse(w, map[string]string{
				ErrInfo : err.Error(),
				ErrCode : oerr.UndefinedErrCode,
			})
			return
		}

		if len(res) == 0 {
			generalResponse(w, map[string]string{
				ErrInfo : oerr.ErrDescription[oerr.UnregisterErrCode],
				ErrCode : oerr.UnregisterErrCode,
			})
			return
		}

		if !VerifyVerfCode(tel, tel, verificationCode) {
			generalResponse(w, map[string]string{
				ErrInfo : oerr.ErrDescription[oerr.TelNoVerifyCode],
				ErrCode : oerr.TelNoVerifyCode})
			return
		}

		userId := res["id"]
		ac = getOracleAc()
		res, err = findByCore(cliCtx, storeName, ac, nftAppCode, "password", "user_id", userId)
		if err != nil {
			generalResponse(w, map[string]string{
				ErrInfo : err.Error(),
				ErrCode : oerr.UndefinedErrCode,
			})
			return
		}

		owner := loadOracleAddr()
		msgs := make([]oracle.UniversalMsg, 0)
		if len(res) != 0 {
			pswdId, _ := strconv.Atoi(res["id"])
			msgFreeze := types.NewMsgFreezeRow(owner, nftAppCode, nftPasswordTable, uint(pswdId))
			msgs = append(msgs, msgFreeze)
		}

		//
		fieldsOfPassword := map[string]string{
			"user_id" : userId,
			"password" : password,
		}
		bz , _ := json.Marshal(fieldsOfPassword)
		msgInsert := types.NewMsgInsertRow(owner, nftAppCode, nftPasswordTable, bz)
		msgs = append(msgs, msgInsert)

		oracle.BuildTxsAndBroadcast(cliCtx, msgs)
		generalResponse(w, map[string]string{
			ErrInfo : oerr.ErrDescription[oerr.SuccessCode],
			ErrCode : oerr.SuccessCode,
		})
		return
	}
}

func nftMakeBefore(cliCtx context.CLIContext, storeName string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		file, fileHeader, err := r.FormFile("file")
		if err != nil {
			generalResponse(w, map[string]string{
				ErrInfo : err.Error(),
				ErrCode : oerr.UndefinedErrCode,
			})
			return
		}

		cid, err := uploadFileToIpfs(file, fileHeader.Filename)
		if err != nil {
			generalResponse(w, map[string]string{
				ErrInfo : err.Error(),
				ErrCode : oerr.UndefinedErrCode,
			})
			return
		}

		name := r.FormValue("name")
		description := r.FormValue("description")
		redirectURL := r.FormValue("redirect_url")
		number := r.FormValue("number")
		tel := r.FormValue("tel")
		payType :=  r.FormValue("pay_type")
		_, ok := verifySession(w, r)
		if !ok {
			generalResponse(w, map[string]string{
				ErrInfo : oerr.ErrDescription[oerr.UnLoginErrCode],
				ErrCode : oerr.UnLoginErrCode,
			})
			return
		}
		//check if user exist
		ac := getOracleAc()
		nftAppCode := LoadNFTAppCode()
		res, err := findByCore(cliCtx, storeName, ac, nftAppCode, nftUserTable, "tel", tel)
		if err != nil {
			generalResponse(w, map[string]string{
				ErrInfo : err.Error(),
				ErrCode : oerr.UndefinedErrCode,
			})
			return
		}
		if len(res) == 0 {
			generalResponse(w, map[string]string{
				ErrInfo : "tel not exist",
				ErrCode : oerr.UndefinedErrCode,
			})
			return
		}
		//下订单
		inumber, err := strconv.Atoi(number)
		if err != nil {
			generalResponse(w, map[string]string{
				ErrInfo : oerr.ErrDescription[oerr.ParamsErrCode],
				ErrCode : oerr.ParamsErrCode,
			})
			return
		}
		if inumber > 100 || inumber <= 0 {
			generalResponse(w, map[string]string{
				ErrInfo : "the quantity is between 1 and 100",
				ErrCode : oerr.UndefinedErr,
			})
			return
		}
		price := fmt.Sprintf("%f", float32(inumber) * nftMakePricePerNft)

		pk, addr , _ := loadSpecialPkForNtf()
		fields , _ := json.Marshal(map[string]string{
			"tel" : tel,
			"price" : price,
		})

		msg := types.NewMsgInsertRow(addr, nftAppCode, nftMakeOrder, fields)
		if msg.ValidateBasic() != nil {
			generalResponse(w, map[string]string{
				ErrInfo : "all booked",
				ErrCode : oerr.UndefinedErrCode,
			})
			return
		}

		_, status, errInfo := oracle.BuildAndSignBroadcastTx(cliCtx, []oracle.UniversalMsg{msg}, pk,  addr)
		if status != oracle.Success {
			generalResponse(w, map[string]string{
				ErrInfo : errInfo,
				ErrCode : oerr.UndefinedErrCode,
			})
			return
		}

		//get Order Id
		ac = getOracleAc()
		order, err := findByCore(cliCtx, storeName, ac, nftAppCode, nftMakeOrder, "tel", tel)
		if err != nil {
			generalResponse(w, map[string]string{
				ErrInfo : err.Error(),
				ErrCode : oerr.UndefinedErrCode,
			})
			return
		}
		id := order["id"]
		nftPay(w, payType, redirectURL, price, nftAppCode + "_make_" + id)
		//cache data
		data, _ := json.Marshal(map[string]string {
			"name" : name,
			"description" : description,
			"number" : number,
			"tel" : tel,
			"cid" : cid,
		})
		makeOrderCache.Set(nftAppCode + "_make_" + id, cache.MakeNftInfo{Data: data, TimeStamp: time.Now().Unix()})
	}
}

func nftMakeCore(cliCtx context.CLIContext, storeName string, outTradeNo string)  {
	    cacheData := makeOrderCache.Get(outTradeNo)
	    if cacheData == nil {
			return
		}
	    nftInfo := cacheData.(cache.MakeNftInfo)
	    info := make(map[string]string)
	    json.Unmarshal(nftInfo.Data, &info)
		name := info["name"]
		description := info["description"]
		number := info["number"]
		tel := info["tel"]
		cid := info["cid"]

		ac := getOracleAc()
		nftAppCode := LoadNFTAppCode()
		res, _ := findByCore(cliCtx, storeName, ac, nftAppCode, nftUserTable, "tel", tel)
		id := res["id"]
		denomCode := genCode(16, cliCtx, storeName, denomTable, "code")
		if denomCode == "" {
			fmt.Println("error : gen denom token err")
			return
		}
		//make callFunction data
		argument := make([]string, 0)
		fieldsOfDenom := map[string]string {
			"user_id" : id,
			"code" : codePre + denomCode,
			"name" : name,
			"file" : cid,
			"description" : description,
			"number" : number,
		}
		strFieldsOfDenom, _ := json.Marshal(fieldsOfDenom)

		argument = append(argument, denomTable, string(strFieldsOfDenom), "nft")
		inumber, _  := strconv.ParseInt(number, 10, 64)

		for i := 0; i < int(inumber); i++{

			nftCode := genCode(16, cliCtx, storeName, nftTable, "code")
			if nftCode == "" {
				fmt.Println("error : gen nft token err")
				return
			}
			fieldOfNFT := map[string]string {
				"denom_id" : "",
				"code" : codePre + nftCode,
			}
			strFieldOfNFT, _ := json.Marshal(fieldOfNFT)
			argument = append(argument, string(strFieldOfNFT))
		}

		strArgument, _ := json.Marshal(argument)
		//make nft
		err := callFunction(cliCtx,nftAppCode, "makeNFT", string(strArgument))
		if err != nil {
			fmt.Println("error ： make err")
			return
		}
		return
}

func nftMakeOld(cliCtx context.CLIContext, storeName string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		//TODO 查询制作订单
		file, fileHeader, err := r.FormFile("file")
		if err != nil {
			generalResponse(w, map[string]string{
				ErrInfo : err.Error(),
				ErrCode : oerr.UndefinedErrCode,
			})
			return
		}

		cid, err := uploadFileToIpfs(file, fileHeader.Filename)
		if err != nil {
			generalResponse(w, map[string]string{
				ErrInfo : err.Error(),
				ErrCode : oerr.UndefinedErrCode,
			})
			return
		}

		name := r.FormValue("name")
		description := r.FormValue("description")
		number := r.FormValue("number")
		tel := r.FormValue("tel")
		_, ok := verifySession(w, r)
		if !ok {
			generalResponse(w, map[string]string{
				ErrInfo : oerr.ErrDescription[oerr.UnLoginErrCode],
				ErrCode : oerr.UnLoginErrCode,
			})
			return
		}
		//find user_id
		ac := getOracleAc()
		nftAppCode := LoadNFTAppCode()
		res, err := findByCore(cliCtx, storeName, ac, nftAppCode, nftUserTable, "tel", tel)
		if err != nil {
			generalResponse(w, map[string]string{
				ErrInfo : err.Error(),
				ErrCode : oerr.UndefinedErrCode,
			})
			return
		}
		if len(res) == 0 {
			generalResponse(w, map[string]string{
				ErrInfo : "tel not exist",
				ErrCode : oerr.UndefinedErrCode,
			})
			return
		}

		id := res["id"]
		//make callFunction data
		argument := make([]string, 0)
		denomCode := genCode(16, cliCtx, storeName, denomTable, "code")
		if denomCode == "" {
			generalResponse(w, map[string]string{
				ErrInfo : "gen denom token err",
				ErrCode : oerr.UndefinedErrCode,
			})
			return
		}
		fieldsOfDenom := map[string]string {
			"user_id" : id,
			"code" : codePre + denomCode,
			"name" : name,
			"file" : cid,
			"description" : description,
			"number" : number,
		}
		strFieldsOfDenom, _ := json.Marshal(fieldsOfDenom)

		argument = append(argument, denomTable, string(strFieldsOfDenom), "nft")
		inumber, err  := strconv.ParseInt(number, 10, 64)
		if err != nil {
			generalResponse(w, map[string]string{
				ErrInfo : oerr.ErrDescription[oerr.ParamsErrCode],
				ErrCode : oerr.ParamsErrCode,
			})
			return
		}
		if inumber > 100 || inumber <= 0 {
			generalResponse(w, map[string]string{
				ErrInfo : "the quantity is between 1 and 100",
				ErrCode : oerr.UndefinedErr,
			})
			return
		}

		for i := 0; i < int(inumber); i++{

			nftCode := genCode(16, cliCtx, storeName, nftTable, "code")
			if nftCode == "" {
				generalResponse(w, map[string]string{
					ErrInfo : "gen nft token err",
					ErrCode : oerr.UndefinedErrCode,
				})
				return
			}

			fieldOfNFT := map[string]string {
				"denom_id" : "",
				"code" : codePre + nftCode,
			}
			strFieldOfNFT, _ := json.Marshal(fieldOfNFT)
			argument = append(argument, string(strFieldOfNFT))
		}

		strArgument, _ := json.Marshal(argument)
		//make nft
		err = callFunction(cliCtx,nftAppCode, "makeNFT", string(strArgument))
		if err != nil {
			generalResponse(w, map[string]string{
				ErrInfo : "make err",
				ErrCode : oerr.UndefinedErrCode,
			})
			return
		}
		generalResponse(w, map[string]string{
			ErrInfo : oerr.ErrDescription[oerr.SuccessCode],
			ErrCode : oerr.SuccessCode,
		})
		return
	}
}

func nftPublish(cliCtx context.CLIContext, storeName string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		//denom tel
		data, err := readBodyData(r)
		if err != nil {
			generalResponse(w, map[string]string{
				ErrInfo : err.Error(),
				ErrCode : oerr.UndefinedErrCode,
			})
			return
		}
		userId, ok := verifySession(w, r)
		if !ok {
			generalResponse(w, map[string]string{
				ErrInfo : oerr.ErrDescription[oerr.UnLoginErrCode],
				ErrCode : oerr.UnLoginErrCode,
			})
			return
		}

		if !userAuthenticationStatus(userId) {
			generalResponse(w, map[string]string{
				ErrInfo : oerr.ErrDescription[oerr.UnauthorizedErrCode],
				ErrCode : oerr.UnauthorizedErrCode,
			})
			return
		}

		denomId := data["denom_id"]
		price := data["price"]
		if !checkPriceValid(price) {
			generalResponse(w, map[string]string{
				ErrInfo : oerr.ErrDescription[oerr.ParamsErrCode],
				ErrCode : oerr.ParamsErrCode,
			})
			return
		}
		//check if denomId exists
		nftAppCode := LoadNFTAppCode()
		if !checkIfDataExistInDatabaseByFindBy(cliCtx, storeName, nftAppCode, nftTable, "denom_id", denomId) {
			generalResponse(w, map[string]string{
				ErrInfo : "nft dont exist",
				ErrCode : oerr.UndefinedErrCode,
			})
			return
		}

		if checkIfDataExistInDatabaseByFindBy(cliCtx, storeName, nftAppCode, nftPublishTable, "denom_id", denomId) {
			generalResponse(w, map[string]string{
				ErrInfo : "nfts have been published",
				ErrCode : oerr.UndefinedErrCode,
			})
			return
		}

		fields := map[string]string {
			"denom_id" : denomId,
			"price" : price,
		}
		owner := loadOracleAddr()
		strFields , _ := json.Marshal(fields)
		msg := types.NewMsgInsertRow(owner, nftAppCode, nftPublishTable, strFields)
		if msg.ValidateBasic() != nil {
			generalResponse(w, map[string]string{
				ErrInfo : oerr.ErrDescription[oerr.UndefinedErrCode],
				ErrCode : oerr.UndefinedErrCode,
			})
			return
		}

		oracle.BuildTxsAndBroadcast(cliCtx, []oracle.UniversalMsg{msg})
		generalResponse(w, map[string]string{
			ErrInfo : oerr.ErrDescription[oerr.SuccessCode],
			ErrCode : oerr.SuccessCode,
		})
		return
	}
}

func nftWithdraw(cliCtx context.CLIContext, storeName string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		//denom tel
		data, err := readBodyData(r)
		if err != nil {
			generalResponse(w, map[string]string{
				ErrInfo : err.Error(),
				ErrCode : oerr.UndefinedErrCode,
			})
			return
		}
		_, ok := verifySession(w, r)
		if !ok {
			generalResponse(w, map[string]string{
				ErrInfo : oerr.ErrDescription[oerr.UnLoginErrCode],
				ErrCode : oerr.UnLoginErrCode,
			})
			return
		}
		denomId := data["denom_id"]
		//check if denomId exists
		ac := getOracleAc()
		nftAppCode := LoadNFTAppCode()
		res, err := findByCore(cliCtx, storeName, ac, nftAppCode, nftPublishTable, "denom_id", denomId)
		if err != nil {
			generalResponse(w, map[string]string{
				ErrInfo : err.Error(),
				ErrCode : oerr.UndefinedErrCode,
			})
			return
		}

		if len(res) == 0 {
			generalResponse(w, map[string]string{
				ErrInfo : "tel not exist",
				ErrCode : oerr.UndefinedErrCode,
			})
			return
		}
		strId := res["id"]
		id ,_ := strconv.Atoi(strId)
		addr := loadOracleAddr()
		msg := types.NewMsgFreezeRow(addr, nftAppCode, nftPublishTable, uint(id))
		if msg.ValidateBasic() != nil {
			generalResponse(w, map[string]string{
				ErrInfo : oerr.ErrDescription[oerr.UndefinedErrCode],
				ErrCode : oerr.UndefinedErrCode,
			})
			return
		}
		oracle.BuildTxsAndBroadcast(cliCtx, []oracle.UniversalMsg{msg})
		generalResponse(w, map[string]string{
			ErrInfo : oerr.ErrDescription[oerr.SuccessCode],
			ErrCode : oerr.SuccessCode,
		})
		return
	}
}

//need txHash, so need special oracle
func nftBuy(cliCtx context.CLIContext, storeName string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		nftAppCode := LoadNFTAppCode()
		store, err := session.Start(stdCtx.Background(), w, r)
		ms , _ := store.Get("tel")
		fmt.Println(ms)

		data, err := readBodyData(r)
		if err != nil {
			generalResponse(w, map[string]string{
				ErrInfo : err.Error(),
				ErrCode : oerr.UndefinedErrCode,
			})
			return
		}
		tel := data["tel"] //be used to check session
		payType := data["pay_type"]
		userId, ok := verifySession(w, r)
		if !ok {
			generalResponse(w, map[string]string{
				ErrInfo : oerr.ErrDescription[oerr.UnLoginErrCode],
				ErrCode : oerr.UnLoginErrCode,
			})
			return
		}
		if !userAuthenticationStatus(userId) {
			generalResponse(w, map[string]string{
				ErrInfo : oerr.ErrDescription[oerr.UnauthorizedErrCode],
				ErrCode : oerr.UnauthorizedErrCode,
			})
			return
		}

		denomId := data["denom_id"]
		redirectURL := data["redirect_url"]
		//1、下订单先检查库存
		ac := getOracleAc()
		nfts, err := findByAll(cliCtx, storeName, ac, nftAppCode, nftTable, "denom_id", denomId)
		if err != nil {
			generalResponse(w, map[string]string{
				ErrInfo : err.Error(),
				ErrCode : oerr.UndefinedErrCode,
			})

			return
		}
		if len(nfts) == 0 {
			generalResponse(w, map[string]string{
				ErrInfo : oerr.ErrDescription[oerr.SoldOutErrCode],
				ErrCode : oerr.SoldOutErrCode,
			})
			return
		}
		//2、检查当前下单量
		if orderSet.Size(denomId) >= len(nfts) {
			generalResponse(w, map[string]string{
				ErrInfo : "all booked",
				ErrCode : oerr.AllBookedErrCode,
			})
			return
		}
		if !orderSet.Set(denomId, len(nfts)) {
			generalResponse(w, map[string]string{
				ErrInfo : "all booked",
				ErrCode : oerr.AllBookedErrCode,
			})
			return
		}
		//3、下单
		pk, addr , _ := loadSpecialPkForNtf()
		orderNftId :=nfts[orderSet.Size(denomId) - 1]["id"]
		fields , _ := json.Marshal(map[string]string{
			"tel" : tel,
			"nft_id" : orderNftId,
		})
		msg := types.NewMsgInsertRow(addr, nftAppCode, nftPublishOrder, fields)
		if msg.ValidateBasic() != nil {
			generalResponse(w, map[string]string{
				ErrInfo : oerr.ErrDescription[oerr.UndefinedErrCode],
				ErrCode : oerr.UndefinedErrCode,
			})
			return
		}

		_, status, errInfo := oracle.BuildAndSignBroadcastTx(cliCtx, []oracle.UniversalMsg{msg}, pk,  addr)
		if status != oracle.Success {
			generalResponse(w, map[string]string{
				ErrInfo : errInfo,
				ErrCode : oerr.UndefinedErrCode,
			})
			return
		}

		//get Order Id
		ac = getOracleAc()
		order, err := findByCore(cliCtx, storeName, ac, nftAppCode, nftPublishOrder, "nft_id", orderNftId)
		if err != nil {
			generalResponse(w, map[string]string{
				ErrInfo : err.Error(),
				ErrCode : oerr.UndefinedErrCode,
			})
			return
		}
		id := order["id"]
		// get Money
		ac = getOracleAc()
		publish, err := findByCore(cliCtx, storeName, ac, nftAppCode, nftPublishTable, "denom_id", denomId)
		if err != nil {
			generalResponse(w, map[string]string{
				ErrInfo : err.Error(),
				ErrCode : oerr.UndefinedErrCode,
			})
			return
		}
		if len(publish) == 0 {
			generalResponse(w, map[string]string{
				ErrInfo : oerr.ErrDescription[oerr.SoldOutErrCode],
				ErrCode : oerr.SoldOutErrCode,
			})
			return
		}
		money := publish["price"]
		//4、 pay
		nftPay(w, payType,redirectURL, money, nftAppCode + "_buy_" + id)
		return
	}
}

func nftBuyCore( cliCtx context.CLIContext, storeName string, nftId, addr, outTradeNo ,money string) {
	// nft transfer
	nftAppCode := LoadNFTAppCode()
	freezeIds, _ := json.Marshal([]string{nftId})
	insertValue, _ := json.Marshal(map[string]string {
		"nft_id" : nftId,
		"owner" : addr,
	})

	denomId, sellOut := checkIfSellOut(cliCtx, storeName, nftId)
	argument := []string{ nftTable, string(freezeIds), nftCardBagTable, string(insertValue) }
	strArgument, _ := json.Marshal(argument)
	err := callFunction(cliCtx,nftAppCode, "nft_deliver", string(strArgument))
	if err != nil {
		fmt.Println("serious error ： ", outTradeNo, " nft deliver fail" )
		return
	}
	//update token for seller and buyer
	sellerInfo := findAuthorInfoByNftId(cliCtx, storeName, nftId)
	ac := getOracleAc()
	buyerInfo, err := findByCore(cliCtx, storeName, ac, nftAppCode, nftUserTable, "address", addr)
	if err != nil {
		fmt.Println("find buyer info err" )
		return
	}
	imoney, _ := strconv.ParseFloat(money, 32)
	if int(imoney) != 0 {
		updateScoreCore(cliCtx, storeName, sellerInfo["user_id"], "+", int(imoney), "sell")
		updateScoreCore(cliCtx, storeName, buyerInfo["id"], "+", int(imoney), "buy")
	}

	if sellOut {
		withdrawSoldOut(cliCtx, storeName, denomId)
	}
	return
}
//
func nftPay(w http.ResponseWriter, PayType ,RedirectURL, Money, OutTradeNo string) {
	var url string
	var err error
	if PayType == "app" {
		url, err = oracleAppPay(Money, OutTradeNo)
	} else if PayType == "web"{
		url, err = oraclePagePay(RedirectURL, Money, OutTradeNo)
	} else {
		url, err = oracleAppPay(Money, OutTradeNo)
	}

	if err != nil {
		generalResponse(w, map[string]string{
			ErrInfo : err.Error(),
			ErrCode : oerr.UndefinedErr,
		})
		return
	}

	result := map[string]string {
		"url" : url,
	}
	successDataResponse(w, result)
	return
}

func nftSaveReceipt(cliCtx context.CLIContext, storeName string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		outTradeNo := r.Form.Get("out_trade_no")
		totalAmount, tradeNo := verifyAlipay(outTradeNo)
		if totalAmount == "" {
			if _, err := aliClient.VerifySign(r.Form); err != nil {
				fmt.Println("aliClient.VerifySign  err : ", err.Error())
				w.Write([]byte("failed"))
				return
			}
			totalAmount  = strings.TrimSpace(r.Form.Get("total_amount"))
			tradeNo = strings.TrimSpace(r.Form.Get("trade_no"))
		}

		//get owner
		ac := getOracleAc()
		nftAppCode := LoadNFTAppCode()
		ss := strings.Split(outTradeNo, "_")
		id := ss[2]
		queryString := fmt.Sprintf("custom/%s/find/%s/%s/%s/%s", storeName, ac, nftAppCode, nftPublishOrder, id)
		order, err := oracleQueryUserTable(cliCtx, queryString)
		if err != nil {
			fmt.Println("serious error ： ", tradeNo, " get order fail, but payed" )
			return
		}

		ac = getOracleAc()
		user, err := findByCore(cliCtx, storeName, ac, nftAppCode, nftUserTable, "tel", order["tel"])
		if err != nil || len(user) == 0 {
			fmt.Println("serious error ： ", tradeNo, " get user addr, but payed" )
		}
		//
		fields , _ := json.Marshal(map[string]string {
			"appcode" : nftAppCode,
			"orderid" : outTradeNo,
			"owner" : user["address"],
			"amount"  : totalAmount,
			"vendor"  : "alipay",
			"vendor_payment_no" : tradeNo,
		})

		err = insertRowWithTx(cliCtx, nftAppCode, nftOrderReceipt, fields)
		if err != nil {
			fmt.Println("serious error ： ", tradeNo, " save receipt fail, but payed" )
			return
		}
		if ss[1] == "make" {
			nftMakeCore(cliCtx, storeName, outTradeNo)
		} else {
			nftBuyCore(cliCtx,storeName, order["nft_id"], user["address"], tradeNo, totalAmount)
		}
		return
	}
}

func nftSaveReceiptInitiative(cliCtx context.CLIContext, storeName string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		outTradeNo := r.Form.Get("out_trade_no")
		totalAmount, tradeNo := verifyAlipay(outTradeNo)
		if totalAmount == "" {
			if _, err := aliClient.VerifySign(r.Form); err != nil {
				fmt.Println("aliClient.VerifySign  err : ", err.Error())
				w.Write([]byte("failed"))
				return
			}
			totalAmount  = strings.TrimSpace(r.Form.Get("total_amount"))
			tradeNo = strings.TrimSpace(r.Form.Get("trade_no"))
		}

		//get owner
		ac := getOracleAc()
		nftAppCode := LoadNFTAppCode()
		ss := strings.Split(outTradeNo, "_")
		id := ss[2]
		queryString := fmt.Sprintf("custom/%s/find/%s/%s/%s/%s", storeName, ac, nftAppCode, nftPublishOrder, id)
		order, err := oracleQueryUserTable(cliCtx, queryString)
		if err != nil {
			fmt.Println("serious error ： ", tradeNo, " get order fail, but payed" )
			generalResponse(w, map[string]string{
				ErrInfo : "serious error ： " + tradeNo + " get order fail, but payed",
				ErrCode : oerr.UndefinedErrCode,
			})
			return
		}

		ac = getOracleAc()
		user, err := findByCore(cliCtx, storeName, ac, nftAppCode, nftUserTable, "tel", order["tel"])
		if err != nil || len(user) == 0 {
			fmt.Println("serious error ： ", tradeNo, " get user addr, but payed" )
			generalResponse(w, map[string]string{
				ErrInfo : "serious error ： " + tradeNo + " get user addr, but payed",
				ErrCode : oerr.UndefinedErrCode,
			})
			return
		}
		//
		seller := ""
		if ss[1] == "buy" {
			seller = getNftMakerFromOutTradeNum(cliCtx, storeName, outTradeNo)
		}
		fields , _ := json.Marshal(map[string]string {
			"appcode" : nftAppCode,
			"orderid" : outTradeNo,
			"owner" : user["address"],
			"amount"  : totalAmount,
			"vendor"  : "alipay",
			"vendor_payment_no" : tradeNo,
			"seller" : seller,
		})
		err = insertRowWithTx(cliCtx, nftAppCode, nftOrderReceipt, fields)
		if err != nil {
			fmt.Println("serious error ： ", tradeNo, " save receipt fail, but payed" )
			generalResponse(w, map[string]string{
				ErrInfo : "serious error ： " + tradeNo + " save receipt fail, but payed",
				ErrCode : oerr.UndefinedErrCode,
			})
			return
		}
		if ss[1] == "make" {
			nftMakeCore(cliCtx, storeName, outTradeNo)
		} else {
			nftBuyCore(cliCtx,storeName, order["nft_id"], user["address"], tradeNo, totalAmount)
		}
		generalResponse(w, map[string]string{
			ErrInfo : oerr.ErrDescription[oerr.SuccessCode],
			ErrCode : oerr.SuccessCode,
		})
		return
	}
}

func nftTransfer(cliCtx context.CLIContext, storeName string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		data, err := readBodyData(r)
		if err != nil {
			generalResponse(w, map[string]string{
				ErrInfo : err.Error(),
				ErrCode : oerr.UndefinedErrCode,
			})
			return
		}
		tel := data["tel"] //be used to check session
		_, ok := verifySession(w, r)
		if !ok {
			generalResponse(w, map[string]string{
				ErrInfo : oerr.ErrDescription[oerr.UnLoginErrCode],
				ErrCode : oerr.UnLoginErrCode,
			})
			return
		}
		nftId := data["nft_id"]
		toAddr := data["to_addr"]
		//1、check owner
		ac := getOracleAc()
		nftAppCode := LoadNFTAppCode()
		user, err := findByCore(cliCtx, storeName, ac, nftAppCode, nftUserTable, "tel", tel)
		if err != nil || len(user) == 0 {
			generalResponse(w, map[string]string{
				ErrInfo : "user don't exit",
				ErrCode : oerr.UndefinedErrCode,
			})
			return
		}
		addr := user["address"]
		//2、get nft info
		ac = getOracleAc()
		nft, err := findByCore(cliCtx, storeName, ac, nftAppCode, nftCardBagTable, "nft_id", nftId)
		if err != nil || len(nft) == 0 {
			generalResponse(w, map[string]string{
				ErrInfo : "nft don't exit",
				ErrCode : oerr.UndefinedErrCode,
			})
			return
		}
		if addr != nft["owner"] {
			generalResponse(w, map[string]string{
				ErrInfo : "permission forbidden",
				ErrCode : oerr.UndefinedErrCode,
			})
			return
		}

		// nft transfer
		freezeIds, _ := json.Marshal([]string{nft["id"]})
		insertValue, _ := json.Marshal(map[string]string {
			"nft_id" : nft["nft_id"],
			"owner" : toAddr,
			"last_id" : nft["id"],
		})

		argument := []string{ nftTable, string(freezeIds), nftCardBagTable, string(insertValue) }
		strArgument, _ := json.Marshal(argument)
		err = callFunction(cliCtx,nftAppCode, "nft_deliver", string(strArgument))
		if err != nil {
			generalResponse(w, map[string]string{
				ErrInfo : "nft transfer fail",
				ErrCode : oerr.UndefinedErrCode,
			})
			return
		}
		generalResponse(w, map[string]string{
			ErrInfo : oerr.ErrDescription[oerr.SuccessCode],
			ErrCode : oerr.SuccessCode,
		})
		return
	}
}

func nftEditPersonalInformation(cliCtx context.CLIContext, storeName string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		//If information has been edit , return
		//
		nickname := r.FormValue("nickname")
		description := r.FormValue("description")
		tel := r.FormValue("tel")
		_, ok := verifySession(w, r)
		if !ok {
			generalResponse(w, map[string]string{
				ErrInfo : oerr.ErrDescription[oerr.UnLoginErrCode],
				ErrCode : oerr.UnLoginErrCode,
			})
			return
		}

		userId, ok := CanEditPersonalInfo(cliCtx, storeName, tel)
		if ! ok {
			generalResponse(w, map[string]string{
				ErrInfo : "edit err",
				ErrCode : oerr.UndefinedErrCode,
			})
			return
		}

		file, fileHeader, err := r.FormFile("file")
		if err != nil {
			generalResponse(w, map[string]string{
				ErrInfo : err.Error(),
				ErrCode : oerr.UndefinedErrCode,
			})
			return
		}

		cid, err := uploadFileToIpfs(file, fileHeader.Filename)
		if err != nil {
			generalResponse(w, map[string]string{
				ErrInfo : err.Error(),
				ErrCode : oerr.UndefinedErrCode,
			})
			return
		}

		fields,_ := json.Marshal(map[string]string{
			"user_id" : userId,
			"avatar" : cid,
			"nickname" : nickname,
			"description" : description,
		})

		owner := loadOracleAddr()
		nftAppCode := LoadNFTAppCode()
		msg := types.NewMsgInsertRow(owner, nftAppCode, nftUserInfoTable, fields)
		if msg.ValidateBasic() != nil {
			generalResponse(w, map[string]string{
				ErrInfo : err.Error(),
				ErrCode : oerr.UndefinedErrCode,
			})
			return
		}

		oracle.BuildTxsAndBroadcast(cliCtx, []oracle.UniversalMsg{msg})
		generalResponse(w, map[string]string{
			ErrInfo : oerr.ErrDescription[oerr.SuccessCode],
			ErrCode : oerr.SuccessCode,
		})
		return
	}
}

func nftCollect(cliCtx context.CLIContext, storeName string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		data, err := readBodyData(r)
		if err != nil {
			generalResponse(w, map[string]string{
				ErrInfo : err.Error(),
				ErrCode : oerr.UndefinedErrCode,
			})
			return
		}
		denomId := data["denom_id"]
		userId, ok := verifySession(w, r)
		if !ok {
			generalResponse(w, map[string]string{
				ErrInfo : oerr.ErrDescription[oerr.UnLoginErrCode],
				ErrCode : oerr.UnLoginErrCode,
			})
			return
		}
		queryString := `[{"method":"table","table":"nft_collection"},{"method":"select","fields":"user_id,denom_id"},{"method":"where", "field" : "user_id", "operator" : "=", "value" : "` + userId + `"},{"method":"where", "field" : "denom_id", "operator" : "=", "value" : "` + denomId + `"}]`
		collect := queryByQuerier(queryString)
		if len(collect) != 0 {
			generalResponse(w, map[string]string{
				ErrInfo : "have collected",
				ErrCode : oerr.UndefinedErr,
			})
			return
		}

		fields,_ := json.Marshal(map[string]string{
			"user_id" : userId,
			"denom_id" : denomId,
		})

		owner := loadOracleAddr()
		nftAppCode := LoadNFTAppCode()
		msg := types.NewMsgInsertRow(owner, nftAppCode, nftCollection, fields)
		if msg.ValidateBasic() != nil {
			generalResponse(w, map[string]string{
				ErrInfo : err.Error(),
				ErrCode : oerr.UndefinedErrCode,
			})
			return
		}

		oracle.BuildTxsAndBroadcast(cliCtx, []oracle.UniversalMsg{msg})
		generalResponse(w, map[string]string{
			ErrInfo : oerr.ErrDescription[oerr.SuccessCode],
			ErrCode : oerr.SuccessCode,
		})
		return
	}
}

func nftCancelCollect(cliCtx context.CLIContext, storeName string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data, err := readBodyData(r)
		if err != nil {
			generalResponse(w, map[string]string{
				ErrInfo : err.Error(),
				ErrCode : oerr.UndefinedErrCode,
			})
			return
		}

		nftAppCode := LoadNFTAppCode()
		denomId := data["denom_id"]
		userId, ok := verifySession(w, r)
		if !ok {
			generalResponse(w, map[string]string{
				ErrInfo : oerr.ErrDescription[oerr.UnLoginErrCode],
				ErrCode : oerr.UnLoginErrCode,
			})
			return
		}
		queryString := `[{"method":"table","table":"nft_collection"},{"method":"select","fields":"id"},{"method":"where", "field" : "user_id", "operator" : "=", "value" : "` + userId + `"},{"method":"where", "field" : "denom_id", "operator" : "=", "value" : "` + denomId + `"}]`
		collect := queryByQuerier(queryString)
		if len(collect) == 0 {
			generalResponse(w, map[string]string{
				ErrInfo : "not collected",
				ErrCode : oerr.UndefinedErr,
			})
			return
		}

		owner := loadOracleAddr()
		id , _ := strconv.Atoi(collect[0]["id"])
		msg := types.NewMsgFreezeRow(owner, nftAppCode, nftCollection, uint(id))
		if msg.ValidateBasic() != nil {
			generalResponse(w, map[string]string{
				ErrInfo : err.Error(),
				ErrCode : oerr.UndefinedErrCode,
			})
			return
		}

		oracle.BuildTxsAndBroadcast(cliCtx, []oracle.UniversalMsg{msg})
		generalResponse(w, map[string]string{
			ErrInfo : oerr.ErrDescription[oerr.SuccessCode],
			ErrCode : oerr.SuccessCode,
		})
		return
	}
}

//////////////////////////////////
//                              //
// help func                    //
//                              //
//////////////////////////////////
//register tel must verify code first, tel which was verified should be cached for register
func getNftMakerFromOutTradeNum(cliCtx context.CLIContext, storeName, outTradeNumber string) string {
	nftAppCode := LoadNFTAppCode()
	tradeInfo := strings.Split(outTradeNumber, "_")
	orderId := tradeInfo[2]
	ac := getOracleAc()
	BaseUrl := LoadNFTBaseUrl()
	queryString := fmt.Sprintf("%s/find/%s/%s/%s/%s", BaseUrl, ac, nftAppCode, nftPublishOrder, orderId)
	orderInfo, err := findRow(cliCtx, queryString)
	if err != nil {
		return ""
	}
	ac = getOracleAc()
	queryString = fmt.Sprintf("%s/find/%s/%s/%s/%s", BaseUrl, ac, nftAppCode, nftTable, orderInfo["nft_id"])
	nftInfo, err := findRow(cliCtx, queryString)
	if err != nil {
		return ""
	}
	ac = getOracleAc()
	queryString = fmt.Sprintf("%s/find/%s/%s/%s/%s", BaseUrl, ac, nftAppCode, denomTable, nftInfo["denom_id"])
	denomInfo, err := findRow(cliCtx, queryString)
	if err != nil {
		return ""
	}
	ac = getOracleAc()
	queryString = fmt.Sprintf("%s/find/%s/%s/%s/%s", BaseUrl, ac, nftAppCode, nftUserTable, denomInfo["user_id"])
	userInfo, err := findRow(cliCtx, queryString)
	if err != nil {
		return ""
	}
	return userInfo["address"]
}

func verifyTel() bool {
	//TODO
	return true
}

func passwordFormatCheck(password string) (string,bool) {
	 hasNum, hasLower, hasUper := false, false, false
	 if len(password) < 8 {
		 return "", false
	 }

	 for _, a := range password {

		if a <=  '9' && a >= '0' {
			hasNum = true
		} else if a <=  'z' && a >= 'a' {
			hasLower = true
		} else if a <= 'Z' && a >= 'A' {
			hasUper = true
		}

		 if hasNum && hasLower && hasUper {
			 hs := sha256.Sum256([]byte(password))
			 return hex.EncodeToString(hs[:]), true
		 }
	 }
	return "", false
}

func insertUserTable(fields map[string]string) {
	//TODO ner app
	//tableName := "user"
	//insertRow(nftAppCode, tableName, fields)
}

func insertRowWithTx (cliCtx context.CLIContext, appcode, tableName string, fields []byte) error {
	pk, addr , _ := loadSpecialPkForNtf()
	msg := types.NewMsgInsertRow(addr, appcode, tableName, fields)
	if err := msg.ValidateBasic(); err  != nil {
		return err
	}

	_, status, errInfo := oracle.BuildAndSignBroadcastTx(cliCtx, []oracle.UniversalMsg{msg}, pk,  addr)
	if status != oracle.Success {
		return errors.New(errInfo)
	}
	return nil
}

func freezeRow(appcode, tableName string , id int) {

}

func genUserAddress() string {
	pk := sm2.GenPrivKey()
	addr := sdk.AccAddress(pk.PubKey().Address())
	return addr.String()
}

func genCode(length int, cliCtx context.CLIContext, storeName ,tableName , field string) string {
	nftAppCode := LoadNFTAppCode()
	for i := 0; i < 10; i++{
		code := make([]byte, length)
		rand.Read(code)
		strCode := hex.EncodeToString(code)
		ac := getOracleAc()
		ids, _ := findByCoreIds(cliCtx, storeName, ac, nftAppCode, tableName, field, strCode)
		if len(ids) ==  0 {
			return strCode
		}
	}
	return ""
}

func callFunction(cliCtx context.CLIContext, appCode, functionName, argument string) error {
	addr := loadOracleAddr()
	if addr == nil {
		return errors.New("loadOracleArr err")
	}
	msg := types.NewMsgCallFunction(addr, appCode, functionName, argument)
	oracle.BuildTxsAndBroadcast(cliCtx,  []oracle.UniversalMsg{msg})
	return nil
}

func loadOracleAddr() sdk.AccAddress {
	pk, err := loadOraclePrivateKey()
	if err != nil {
		return nil
	}
	addr := sdk.AccAddress(pk.PubKey().Address())
	return addr
}

func uploadFileToIpfs(file io.Reader, fileName string) (string, error) {
	ac := getOracleAc()
	nftAppCode := LoadNFTAppCode()
	BaseUrl := LoadNFTBaseUrl()
	url := fmt.Sprintf("%s/upload/%s/%s", BaseUrl, ac, nftAppCode)
	return postFile(file, fileName, url)
}


func postFile(file io.Reader, fileName string, targetUrl string) (string,error) {
	bodyBuf := &bytes.Buffer{}
	bodyWriter := multipart.NewWriter(bodyBuf)

	//关键的一步操作
	fileWriter, err := bodyWriter.CreateFormFile("file", fileName)
	if err != nil {
		fmt.Println("error writing to buffer")
		return "", err
	}
	//iocopy
	_, err = io.Copy(fileWriter, file)
	if err != nil {
		return "", err
	}

	contentType := bodyWriter.FormDataContentType()
	bodyWriter.Close()

	resp, err := http.Post(targetUrl, contentType, bodyBuf)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	result := make(map[string]string, 0)
	json.Unmarshal(respBody, &result)
	cid := result["result"]
	if cid == "" {
		return "", errors.New("upload file fail")
	} else {
		return cid, nil
	}
}


func genNFTCode() string {
	code := make([]byte, 16)
	rand.Read(code)
	//TODO check code exist or not
	return hex.EncodeToString(code)
}

func checkUniqueCode(code string) bool {
	//TODO
	return true
}

func checkIfDataExistInDatabaseByFindBy(cliCtx context.CLIContext, storeName , nftAppCode, tableName, field, value string) bool {
	ac := getOracleAc()
	res, err := findByCore(cliCtx, storeName, ac, nftAppCode, tableName, field, value)
	if err != nil {
		return false
	}

	if len(res) == 0 {
		return false
	}
	return true
}

func checkPriceValid(price string) bool {
	//最多两位小数
	if priceRex == nil {
		priceRex , _ = regexp.Compile(priceRegExp)
	}
	return priceRex.MatchString(price)
}

func loadNftOraclePkAndAddr() (crypto.PrivKey, sdk.AccAddress, error) {
	//TODO,
	return nil, nil, nil
}

func checkIfSellOut(cliCtx context.CLIContext, storeName , nftId string) (string, bool){
	// the last nft was sold, it should be withdraw
	ac := getOracleAc()
	nftAppCode := LoadNFTAppCode()
	queryString := fmt.Sprintf("custom/%s/find/%s/%s/%s/%s", storeName, ac, nftAppCode, nftTable, nftId)
	nft, err := oracleQueryUserTable(cliCtx, queryString)
	if err != nil || nft == nil {
		fmt.Println("serious error ： get nft fail, but payed" )
		return "", false
	}
	denomId := nft["denom_id"]
	ids, err := findByCoreIds(cliCtx, storeName, ac, nftAppCode, nftTable, "denom_id", denomId)
	if err != nil {
		fmt.Println("find nfts err")
	}
	if len(ids) <= 1 {
		return denomId, true
	}
	return "", false
}

func withdrawSoldOut(cliCtx context.CLIContext, storeName, denomId string) {
	ac := getOracleAc()
	nftAppCode := LoadNFTAppCode()
	res, err := findByCore(cliCtx, storeName, ac, nftAppCode, nftPublishTable, "denom_id", denomId)
	if err != nil || len(res) == 0{
		return
	}
	strId := res["id"]
	id ,_ := strconv.Atoi(strId)
	addr := loadOracleAddr()
	msg := types.NewMsgFreezeRow(addr, nftAppCode, nftPublishTable, uint(id))
	if msg.ValidateBasic() != nil {
		return
	}
	oracle.BuildTxsAndBroadcast(cliCtx, []oracle.UniversalMsg{msg})
	return
}

func saveSession(w http.ResponseWriter, r *http.Request, tel, userId string) bool {
	store, err := session.Start(stdCtx.Background(), w, r)
	if err != nil {
		return false
	}
	store.Set("tel", tel)
	store.Set("userId", userId)
	err = store.Save()
	if err != nil {
		return false
	}
	return true
}

func verifySession(w http.ResponseWriter, r *http.Request) (string, bool) {
	store, err := session.Start(stdCtx.Background(), w, r)
	if err != nil {
		return "", false
	}
	_, ok := store.Get("tel")
	if !ok {
		return "", false
	}

	userId, _ := store.Get("userId")
	return userId.(string), true
}

func deleteSession(w http.ResponseWriter, r *http.Request)  bool {
	err := session.Destroy(stdCtx.Background(), w, r)
	if err != nil {
		return false
	}
	return  true
}

func verifyAlipay(outTradeNo string) (string, string){
	aliOrderStatus, err := OracleQueryAliOrder(outTradeNo)
	if err != nil {
		return "", ""
	}
	if aliOrderStatus["trade_status"] != string(alipay.TradeStatusSuccess) {
		return "", ""
	}
	// save to OrderReceipt table
	totalAmount  := aliOrderStatus["total_amount"]
	tradeNo := aliOrderStatus["trade_no"]
	return totalAmount, tradeNo
}

func LoadNFTAppCode() string {
	if NFTAppCode != "" {
		return NFTAppCode
	}
	NFTAppCode = viper.GetString(nftAppCodeKey)
	if NFTAppCode == "" {
		fmt.Println("serious err : NFTAppCode is null")
	}
	return NFTAppCode
}