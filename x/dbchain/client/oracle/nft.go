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
	nftAppCode = "KQ3TVRJC15"
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

	priceRegExp = `(^[1-9]\d*(\.\d{1,2})?$)|(^0(\.\d{1,2})?$)`
	invitationScore = 10
)

var priceRex *regexp.Regexp
var orderSet *nftOrderSet

func init() {
	// order cache
	orderSet = newNftOrderSet(time.Second * 300)
	go orderSet.GC()
	//session
	session.InitManager(
		session.SetCookieLifeTime(300),
	)
}

func nftUserRegister(cliCtx context.CLIContext, storeName string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data, err := readBodyData(r)
		if err != nil {
			generalResponse(w, map[string]string{ "error" : err.Error()})
		}
		/* user table
		tel code address invitation code
		input params
		tel password invitation_code
		*/
		tel := data["tel"]
		if verifyTelCache.Get(tel) == nil {
			generalResponse(w, map[string]string{"error" : "please verify tel first"})
			return
		}
		password := data["password"]//
		password, valid := passwordFormatCheck(password)
		if !valid {
			generalResponse(w, map[string]string{"error" : "format of password err"})
			return
		}
		invitationCode := data["invitation_code"]
		fieldsOfUser := map[string]string {
			"tel" : tel,
			"my_code" : genCode(4),
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
			generalResponse(w, map[string]string{"error" : "register fail"})
			return
		}
		if invitationCode != "" {
			updateScore(cliCtx, storeName, invitationCode)
		}
		generalResponse(w, map[string]string{"success" : "register success"})
		return
	}
}

func updateScore(cliCtx context.CLIContext, storeName string, invitationCode string) {

	ac := getOracleAc()
	//find userId
	userId, err := findByCoreIds(cliCtx, storeName, ac, nftAppCode, nftUserTable, "my_code", invitationCode)
	if err != nil || len(userId) == 0 {
		return
	}
	//find user score
	ac = getOracleAc()
	score, err := findByCore(cliCtx, storeName, ac, nftAppCode, nftScoreTable, "user_id", userId[0])
	if err != nil {
		return
	}
	owner := loadOracleAddr()
	msgs := make([]oracle.UniversalMsg, 0)
	lastToken := 0

	if len(score) != 0 {
		lastToken, _ = strconv.Atoi(score["token"])
		id, _ := strconv.Atoi(score["id"])
		msg := types.NewMsgFreezeRow(owner, nftAppCode, nftScoreTable, uint(id))
		if msg.ValidateBasic() != nil {
			return
		}
		msgs = append(msgs, msg)
	}
	currentToken := lastToken + invitationScore
	token := strconv.Itoa(currentToken)
	field , _ := json.Marshal(map[string]string{ "user_id" : userId[0], "token" : token})
	msg := types.NewMsgInsertRow(owner, nftAppCode, nftScoreTable, field)
	if msg.ValidateBasic() != nil {
		return
	}
	msgs = append(msgs, msg)
	oracle.BuildTxsAndBroadcast(cliCtx, msgs)
	return
}

func nftUserLogin(cliCtx context.CLIContext, storeName string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data, err := readBodyData(r)
		if err != nil {
			generalResponse(w, map[string]string{ "error" : err.Error()})
		}

		tel := data["tel"]
		password := data["password"]
		//query tel and password
		ac := getOracleAc()
		res, err := findByCore(cliCtx, storeName, ac, nftAppCode, "user", "tel", tel)
		if err != nil {
			generalResponse(w, map[string]string{"error" : err.Error()})
			return
		}

		if len(res) == 0 {
			generalResponse(w, map[string]string{"error" : "tel not exist"})
			return
		}

		id := res["id"]
		ac = getOracleAc()
		res, err = findByCore(cliCtx, storeName, ac, nftAppCode, "password", "user_id", id)
		if err != nil {
			generalResponse(w, map[string]string{"error" : err.Error()})
			return
		}

		if len(res) == 0 {
			generalResponse(w, map[string]string{"error" : "unknown error"})
			return
		}

		hspswd := res["password"]
		hs := sha256.Sum256([]byte(password))
		if hex.EncodeToString(hs[:]) == hspswd {
			if !saveSession(w, r , tel) {
				generalResponse(w, map[string]string{"error" : "save session err"})
				return
			}
			generalResponse(w, map[string]string{"success" : ""})
			return
		} else {
			generalResponse(w, map[string]string{"error" : "password err"})
			return
		}
	}
}

func nftMake(cliCtx context.CLIContext, storeName string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		//TODO 查询制作订单
		file, fileHeader, err := r.FormFile("file")
		if err != nil {
			generalResponse(w, map[string]string{"error" : err.Error()})
			return
		}

		cid, err := uploadFileToIpfs(file, fileHeader.Filename)
		if err != nil {
			generalResponse(w, map[string]string{"error" : err.Error()})
			return
		}

		name := r.FormValue("name")
		description := r.FormValue("description")
		number := r.FormValue("number")
		tel := r.FormValue("tel")
		if !verifySession(w, r, tel) {
			generalResponse(w, map[string]string{"error" : "please login first"})
			return
		}
		//find user_id
		ac := getOracleAc()
		res, err := findByCore(cliCtx, storeName, ac, nftAppCode, nftUserTable, "tel", tel)
		if err != nil {
			generalResponse(w, map[string]string{"error" : err.Error()})
			return
		}
		if len(res) == 0 {
			generalResponse(w, map[string]string{"error" : "tel not exist"})
			return
		}

		id := res["id"]
		//make callFunction data
		argument := make([]string, 0)
		fieldsOfDenom := map[string]string {
			"user_id" : id,
			"code" : codePre + genCode(16),
			"name" : name,
			"file" : cid,
			"description" : description,
			"number" : number,
		}
		strFieldsOfDenom, _ := json.Marshal(fieldsOfDenom)

		argument = append(argument, denomTable, string(strFieldsOfDenom), "nft")
		inumber, err  := strconv.ParseInt(number, 10, 64)
		if err != nil {
			generalResponse(w, map[string]string{"error" : "number err"})
			return
		}

		for i := 0; i < int(inumber); i++{

			fieldOfNFT := map[string]string {
				"denom_id" : "",
				"code" : codePre + genCode(16),
			}
			strFieldOfNFT, _ := json.Marshal(fieldOfNFT)
			argument = append(argument, string(strFieldOfNFT))
		}

		strArgument, _ := json.Marshal(argument)
		//make nft
		err = callFunction(cliCtx,nftAppCode, "makeNFT", string(strArgument))
		if err != nil {
			generalResponse(w, map[string]string{"error" : "make err"})
			return
		}
		generalResponse(w, map[string]string{"success" : "make NFT success"})
		return
	}
}

func nftPublish(cliCtx context.CLIContext, storeName string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		//denom tel
		data, err := readBodyData(r)
		if err != nil {
			generalResponse(w, map[string]string{ "error" : err.Error()})
			return
		}
		tel := data["tel"] //be used to check session
		if !verifySession(w, r, tel) {
			generalResponse(w, map[string]string{ "error" : "please login first"})
			return
		}
		denomId := data["denom_id"]
		price := data["price"]
		if !checkPriceValid(price) {
			generalResponse(w, map[string]string{"error" : "prices err"})
			return
		}
		//check if denomId exists
		if !checkIfDataExistInDatabaseByFindBy(cliCtx, storeName, nftAppCode, nftTable, "denom_id", denomId) {
			generalResponse(w, map[string]string{"error" : err.Error()})
			return
		}

		fields := map[string]string {
			"denom_id" : denomId,
			"price" : price,
		}
		owner := loadOracleAddr()
		//_, addr, _  := loadNftOraclePkAndAddr()
		strFields , _ := json.Marshal(fields)
		msg := types.NewMsgInsertRow(owner, nftAppCode, nftPublishTable, strFields)
		if msg.ValidateBasic() != nil {
			generalResponse(w, map[string]string{"error" : err.Error()})
			return
		}

		oracle.BuildTxsAndBroadcast(cliCtx, []oracle.UniversalMsg{msg})
		generalResponse(w, map[string]string{"success" : "publishing"})
		return
	}
}

func nftWithdraw(cliCtx context.CLIContext, storeName string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		//denom tel
		data, err := readBodyData(r)
		if err != nil {
			generalResponse(w, map[string]string{ "error" : err.Error()})
			return
		}
		tel := data["tel"] //be used to check session
		if !verifySession(w, r, tel) {
			generalResponse(w, map[string]string{ "error" : "please login first"})
			return
		}
		denomId := data["denom_id"]
		//check if denomId exists
		ac := getOracleAc()
		res, err := findByCore(cliCtx, storeName, ac, nftAppCode, nftPublishTable, "denom_id", denomId)
		if err != nil {
			generalResponse(w, map[string]string{"error" : err.Error()})
			return
		}

		if len(res) == 0 {
			generalResponse(w, map[string]string{"error" : "tel not exist"})
			return
		}
		strId := res["id"]
		id ,_ := strconv.Atoi(strId)
		addr := loadOracleAddr()
		msg := types.NewMsgFreezeRow(addr, nftAppCode, nftPublishTable, uint(id))
		if msg.ValidateBasic() != nil {
			generalResponse(w, map[string]string{"error" : "tel not exist"})
			return
		}
		oracle.BuildTxsAndBroadcast(cliCtx, []oracle.UniversalMsg{msg})
		generalResponse(w, map[string]string{"success" : ""})
		return
	}
}

//need txHash, so need special oracle
func nftBuy(cliCtx context.CLIContext, storeName string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		store, err := session.Start(stdCtx.Background(), w, r)
		ms , _ := store.Get("tel")
		fmt.Println(ms)

		data, err := readBodyData(r)
		if err != nil {
			generalResponse(w, map[string]string{ "error" : err.Error()})
			return
		}
		tel := data["tel"] //be used to check session
		if !verifySession(w, r, tel) {
			generalResponse(w, map[string]string{ "error" : "please login first"})
			return
		}
		denomId := data["denom_id"]
		redirectURL := data["redirect_url"]
		//1、下订单先检查库存
		ac := getOracleAc()
		nfts, err := findByAll(cliCtx, storeName, ac, nftAppCode, nftTable, "denom_id", denomId)
		if err != nil {
			generalResponse(w, map[string]string{ "error" : err.Error()})
			return
		}
		if len(nfts) == 0 {
			generalResponse(w, map[string]string{ "error" : "sell out"})
			return
		}
		//2、检查当前下单量
		if orderSet.Size() >= len(nfts) {
			generalResponse(w, map[string]string{ "error" : "all booked"})
			return
		}
		if !orderSet.Set(len(nfts)) {
			generalResponse(w, map[string]string{ "error" : "all booked"})
			return
		}
		//3、下单
		pk, addr , _ := loadSpecialPkForNtf()
		fields , _ := json.Marshal(map[string]string{
			"tel" : tel,
			"nft_id" : nfts[orderSet.Size()]["id"],
		})
		msg := types.NewMsgInsertRow(addr, nftAppCode, nftPublishOrder, fields)
		if msg.ValidateBasic() != nil {
			generalResponse(w, map[string]string{ "error" : "all booked"})
			return
		}

		_, status, errInfo := oracle.BuildAndSignBroadcastTx(cliCtx, []oracle.UniversalMsg{msg}, pk,  addr)
		if status != oracle.Success {
			generalResponse(w, map[string]string{ "error" : errInfo})
			return
		}

		//get Order Id
		ac = getOracleAc()
		order, err := findByCore(cliCtx, storeName, ac, nftAppCode, nftPublishOrder, "nft_id", denomId)
		if err != nil {
			generalResponse(w, map[string]string{ "error" : err.Error()})
			return
		}
		id := order["id"]
		// get Money
		ac = getOracleAc()
		publish, err := findByCore(cliCtx, storeName, ac, nftAppCode, nftPublishTable, "denom_id", denomId)
		if err != nil {
			generalResponse(w, map[string]string{ "error" : err.Error()})
			return
		}
		if len(publish) == 0 {
			generalResponse(w, map[string]string{ "error" : "sell out"})
			return
		}
		money := publish["price"]
		//4、 pay
		nftPay(w, redirectURL, money, nftAppCode + "_" + id)
		return
	}
}
//
func nftPay(w http.ResponseWriter, RedirectURL, Money, OutTradeNo string) {
	url, err := oraclePagePay(RedirectURL, Money, OutTradeNo)
	if err != nil {
		generalResponse(w, map[string]string{ "error" : err.Error()})
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write([]byte(url))

}

func nftSaveReceipt(cliCtx context.CLIContext, storeName string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		if _, err := aliClient.VerifySign(r.Form); err != nil {
			fmt.Println("aliClient.VerifySign  err : ", err.Error())
			w.Write([]byte("failed"))
			return
		}
		outTradeNo := strings.TrimSpace(r.Form.Get("out_trade_no"))
		totalAmount  := strings.TrimSpace(r.Form.Get("total_amount"))
		tradeNo := strings.TrimSpace(r.Form.Get("trade_no"))
		//get owner
		ac := getOracleAc()
		ss := strings.Split(outTradeNo, "_")
		id := ss[1]
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

		// nft transfer
		freezeIds, _ := json.Marshal([]string{order["nft_id"]})
		insertValue, _ := json.Marshal(map[string]string {
			"nft_id" : order["nft_id"],
			"owner" : user["address"],
		})

		denomId, sellOut := checkIfSellOut(cliCtx, storeName, order["nft_id"])
		argument := []string{ nftTable, string(freezeIds), nftCardBagTable, string(insertValue) }
		strArgument, _ := json.Marshal(argument)
		err = callFunction(cliCtx,nftAppCode, "nft_deliver", string(strArgument))
		if err != nil {
			fmt.Println("serious error ： ", tradeNo, " nft deliver fail" )
			return
		}
		generalResponse(w, map[string]string{"success" : "nft deliver success"})

		if sellOut {
			withdrawSoldOut(cliCtx, storeName, denomId)
		}
		return
	}
}

func nftTransfer(cliCtx context.CLIContext, storeName string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		data, err := readBodyData(r)
		if err != nil {
			generalResponse(w, map[string]string{ "error" : err.Error()})
			return
		}
		tel := data["tel"] //be used to check session
		if !verifySession(w, r, tel) {
			generalResponse(w, map[string]string{ "error" : "please login first"})
			return
		}
		nftId := data["nft_id"]
		toAddr := data["to_addr"]
		//1、check owner
		ac := getOracleAc()
		user, err := findByCore(cliCtx, storeName, ac, nftAppCode, nftUserTable, "tel", tel)
		if err != nil || len(user) == 0 {
			generalResponse(w, map[string]string{ "error" : "user don't exit"})
			return
		}
		addr := user["address"]
		//2、get nft info
		ac = getOracleAc()
		nft, err := findByCore(cliCtx, storeName, ac, nftAppCode, nftCardBagTable, "nft_id", nftId)
		if err != nil || len(nft) == 0 {
			generalResponse(w, map[string]string{ "error" : "nft don't exit"})
			return
		}
		if addr != nft["owner"] {
			generalResponse(w, map[string]string{ "error" : "permission forbidden"})
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
			generalResponse(w, map[string]string{"error" : "nft transfer fail"})
			return
		}
		generalResponse(w, map[string]string{"success" : "nft transfer success"})
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
		if !verifySession(w, r, tel) {
			generalResponse(w, map[string]string{"error" : "please login first"})
			return
		}

		userId, ok := CanEditPersonalInfo(cliCtx, storeName, tel)
		if ! ok {
			generalResponse(w, map[string]string{"error" : "edit err"})
			return
		}

		file, fileHeader, err := r.FormFile("file")
		if err != nil {
			generalResponse(w, map[string]string{"error" : err.Error()})
			return
		}

		cid, err := uploadFileToIpfs(file, fileHeader.Filename)
		if err != nil {
			generalResponse(w, map[string]string{"error" : err.Error()})
			return
		}

		fields,_ := json.Marshal(map[string]string{
			"user_id" : userId,
			"avatar" : cid,
			"nickname" : nickname,
			"description" : description,
		})

		owner := loadOracleAddr()
		msg := types.NewMsgInsertRow(owner, nftAppCode, nftUserInfoTable, fields)
		if msg.ValidateBasic() != nil {
			generalResponse(w, map[string]string{"error" : err.Error()})
			return
		}

		oracle.BuildTxsAndBroadcast(cliCtx, []oracle.UniversalMsg{msg})
		generalResponse(w, map[string]string{"success" : "edit success"})
		return
	}
}

//////////////////////////////////
//                              //
// help func                    //
//                              //
//////////////////////////////////
//register tel must verify code first, tel which was verified should be cached for register
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

	 	if hasNum && hasLower && hasUper {
			hs := sha256.Sum256([]byte(password))
			return hex.EncodeToString(hs[:]), true
		}

		if a <=  '9' && a >= '0' {
			hasNum = true
		} else if a <=  'z' && a >= 'a' {
			hasLower = true
		} else if a <= 'Z' && a >= 'A' {
			hasUper = true
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

func genCode(length int) string {
	code := make([]byte, length)
	rand.Read(code)
	//TODO check code exist or not
	return hex.EncodeToString(code)

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

func saveSession(w http.ResponseWriter, r *http.Request, val string) bool {
	store, err := session.Start(stdCtx.Background(), w, r)
	if err != nil {
		return false
	}
	store.Set("tel", val)
	err = store.Save()
	if err != nil {
		return false
	}
	return true
}

func verifySession(w http.ResponseWriter, r *http.Request, val string) bool {
	store, err := session.Start(stdCtx.Background(), w, r)
	if err != nil {
		return false
	}
	sessionTel, ok := store.Get("tel")
	if !ok {
		return false
	}
	if sessionTel.(string) != val {
		return false
	}
	return true
}