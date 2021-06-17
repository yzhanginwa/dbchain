package oracle

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/cosmos/cosmos-sdk/client/context"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/rest"
	"github.com/gorilla/mux"
	"github.com/mr-tron/base58"
	"github.com/smartwalle/alipay/v3"
	"github.com/spf13/viper"
	"github.com/yzhanginwa/dbchain/x/dbchain/client/oracle/oracle"
	"github.com/yzhanginwa/dbchain/x/dbchain/internal/types"
	"github.com/yzhanginwa/dbchain/x/dbchain/internal/utils"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

/**************************************************************
table name : order:
	sellable_id          //商品id

table name : sellable
fields:
	orginid	    //
	name		//
	type	    //
	term_days	//
	price       //
	memo		//
	volume_limit    //
	records_limit   //
***************************************************************/


var (
	aliClient *alipay.Client
	notifyUrl string
)
const (
	kAppId = "2021002129602543"
	OrderSub   = "YTBox"
	AlipayNotifyURL  = "alipay-notify-url"
 	IsProduction = true
 	OrderReceipt = "order_receipt"
 	IsTest = false
 	AliPay = "alipay"
 	WechatPay = "wechatpay"
 	DbcToken = "dbctoken"
 	ApplePay = "applepay"
)

func init(){
	var err error
	if aliClient, err = alipay.New(kAppId, cmsDbchainCloudPriKey, IsProduction); err != nil {
		fmt.Println("init aliclient err : ", err)
		os.Exit(-1)
	}
	// 使用支付宝证书
	if err = aliClient.LoadAppPublicCert(appCertPublicKey); err != nil {
		fmt.Println("load app public cert from file fail : ", err)
		os.Exit(-1)
	}
	if err = aliClient.LoadAliPayRootCert(alipayRootCert); err != nil {
		log.Println("load alipay root cert from file fail : ", err)
		os.Exit(-1)
	}
	if err = aliClient.LoadAliPayPublicCert(alipayCertPublicKeyRSA2); err != nil {
		log.Println("load aliPay public cert from file fail : ", err)
		os.Exit(-1)
	}
	//need create table buyerorder and orderinfo
}

func oracleRecipientAddress(cliCtx context.CLIContext, storeName string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		accessCode      		:= vars["accessToken"]
		_, err := utils.VerifyAccessCode(accessCode)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusNotFound, err.Error())
			return
		}
		privKey, err := oracle.LoadPrivKey()
		oracleAddress := sdk.AccAddress(privKey.PubKey().Address())
		res := map[string]string{
			"recipient" : oracleAddress.String(),
		}
		bz , err := json.Marshal(res)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusNotFound, "marshal recipient address failed")
			return
		}
		rest.PostProcessResponse(w, cliCtx, bz)
	}
}
func oracleCallDbcPay(cliCtx context.CLIContext, storeName string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		accessCode      := vars["accessToken"]
		payType         := vars["payType"]
		r.ParseForm()

		buyer, err := utils.VerifyAccessCode(accessCode)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusNotFound, err.Error())
			return
		}
		ReturnURL := r.Form.Get("returnURL")
		OutTradeNo := strings.TrimSpace(r.Form.Get("out_trade_no"))
		appcode := r.Form.Get("appcode")
		sellableid := r.Form.Get("sellableid")
		paymentId  := r.Form.Get("paymentid")
		vendor     := r.Form.Get("vendor")
		tableName  := r.Form.Get("tablename")
		if payType == "app" && (vendor == DbcToken || vendor == ApplePay){
			if tableName == ""{
				tableName = "payment"
			}
			receiptData := r.Form.Get("receipt_data")
			internalPurchase(cliCtx, storeName, OutTradeNo, tableName, receiptData, paymentId, vendor, buyer, w)
			return
		}
		//query from user sellable to get money
		res, err := getOrderMoney(cliCtx, storeName, appcode, sellableid, buyer.String())
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusNotFound, err.Error())
			return
		}

		Money := res["price"]
		//test
		if IsTest {
			switch sellableid {
			case "1":
				Money = "0.01"
			case "2":
				Money = "0.02"
			case "3":
				Money = "0.03"
			case "4":
				Money = "0.04"
			case "5":
				Money = "0.05"
			}
		}
		//end
		url := ""
		if payType == "page" {
			url , err = oraclePagePay(ReturnURL, Money, OutTradeNo)
		} else if payType == "app" {
			url , err = oracleAppPay(Money, OutTradeNo)
		} else {
			rest.WriteErrorResponse(w, http.StatusNotFound, "pay type wrong")
			return
		}
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusNotFound, "Failed to submit order : " + err.Error())
			return
		}

		result := make(map[string]string)
		result["url"] = url
		bz , err := json.Marshal(result)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusNotFound, "marshal rowFields failed")
			return
		}
		rest.PostProcessResponse(w, cliCtx, bz)
		return
	}
}

func oracleApplepay(cliCtx context.CLIContext, storeName string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		accessCode      := vars["accessToken"]
		r.ParseForm()

		buyer, err := utils.VerifyAccessCode(accessCode)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusNotFound, err.Error())
			return
		}

		OutTradeNo := strings.TrimSpace(r.Form.Get("out_trade_no"))
		receiptData := r.Form.Get("receipt_data")

		internalPurchase(cliCtx, storeName, OutTradeNo, "", receiptData, "", ApplePay, buyer, w)
		return

	}
}

func internalPurchase(cliCtx context.CLIContext, storeName, OutTradeNo, tableName, receiptData, paymentId, vendor string , buyer sdk.AccAddress, w http.ResponseWriter) {
	if !checkAppleUnverifyReceipt {
		go checkAppleReceiptRunner(cliCtx)
		checkAppleUnverifyReceipt = true
	}
	appcodeAndOrderId := strings.Split(OutTradeNo,"-")
	if len(appcodeAndOrderId) != 2 {
		rest.WriteErrorResponse(w, http.StatusNotFound, "outTradeNo error")
		return
	}
	appcode := appcodeAndOrderId[0]
	if vendor == DbcToken {
		bz , err := callDbcTokenPay(cliCtx, storeName, appcode, tableName, paymentId)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		rest.PostProcessResponse(w, cliCtx, bz)
	} else if vendor ==  ApplePay {
		appleTransactionId, applePayType, err := verifyApplePay(cliCtx, storeName, OutTradeNo, buyer.String(), receiptData)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			if err.Error() == "failed to access Apple server" {
				recheck := unVerifyAppleReceiptData{
					storeName,
					OutTradeNo,
					buyer.String(),
					receiptData,
				}
				unVerifyAppleReceiptBuf<-recheck

			}
			return
		}
		bz , err := callDbcApplePay(cliCtx, storeName, OutTradeNo, appleTransactionId, applePayType)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		rest.PostProcessResponse(w, cliCtx, bz)
	}
	return
}

func callDbcTokenPay(cliCtx context.CLIContext, storeName, appcode, tableName ,paymentId string) ([]byte, error){
	//check
	//comfirm table has payment option
	hasPaymentOption := checkTableOption(cliCtx, storeName, appcode, tableName, types.TBLOPT_PAYMENT)
	if !hasPaymentOption {
		return nil, errors.New("your table does not has payment option")
	}

	ac := getOracleAc()
	queryString := fmt.Sprintf("custom/%s/find/%s/%s/%s/%s", storeName, ac, appcode, tableName, paymentId)
	paymentInfo ,err := oracleQueryUserTable(cliCtx, queryString)
	if err != nil {
		return nil, err
	}
	//TODO need a special address
	privKey, err := oracle.LoadPrivKey()
	recipientAddress := sdk.AccAddress(privKey.PubKey().Address())
	if paymentInfo["recipient"]  != recipientAddress.String() {
		return  nil, errors.New("invalid recipient address")
	}
	orderId := paymentInfo["orderid"]
	orderInfo, err  := getOrderInfo(cliCtx, storeName, appcode, orderId)
	if err != nil {
		return nil, err
	}
	sellableId := orderInfo["sellable_id"]
	sellableInfo, err := getSellableInfo(cliCtx, storeName, appcode, sellableId)
	if err != nil {
		return  nil, err
	}
	//TODO need tokenNum field
	if sellableInfo["token_price"] != paymentInfo["amount"] {
		return nil, errors.New("amount err")
	}
	//save to order_receipt
	amount := paymentInfo["amount"]
	owner := orderInfo["created_by"]
	expiration_date := calcExpirationDate(cliCtx, storeName, appcode, owner ,orderInfo["sellable_id"])

	res := newOrderReceiptDataCore(appcode, orderId, owner, amount, expiration_date, DbcToken, paymentId)
	oracleAccAddr := oracle.GetOracleAccAddr()
	SaveToOrderInfoTable(cliCtx, oracleAccAddr, res, OrderReceipt)
	bz , _ := json.Marshal(res)
	return bz, nil
}

func checkTableOption(cliCtx context.CLIContext, storeName, appcode, tableName string, TableOption types.TableOption ) bool {
	ac := getOracleAc()
	//comfirm table has payment option
	options, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/option/%s/%s/%s", storeName, ac, appcode, tableName), nil)
	if err != nil {
		return  false
	}
	var out types.QueryTables // QueryTables is a []string. It could be reused here
	json.Unmarshal(options, &out)

	hasPaymentOption := false
	for _, option := range out {
		if option == string(TableOption) {
			hasPaymentOption = true
			break
		}
	}
	return hasPaymentOption
}
func getOrderMoney(cliCtx context.CLIContext, storeName, appcode, sellableid, buyer string) (map[string]string,error) {

	sellableResult, err := getSellableInfo(cliCtx, storeName, appcode, sellableid)
	if err != nil {
		return nil, err
	}
	originCode := sellableResult["origin_code"]
	if originCode == "" {
		return sellableResult, nil
	}
	//check order validity
	//query 0000000001 order_receipt to conform origin_code exist
	fieldValue := map[string]string{
		"appcode": appcode,
		"owner":   buyer,
	}
	orderReceipt, err := getOrderReceiptInfo(cliCtx, storeName, fieldValue)
	if err != nil {
		return nil, err
	}
	if len(orderReceipt) < 1 {
		return nil, errors.New("submit order fialded")
	}
	orderId := orderReceipt[len(orderReceipt)-1]["orderid"]
	//query from user order table
	preOrder, err := getOrderInfo(cliCtx, storeName, appcode, orderId)
	if err != nil {
		return nil, err
	}
	if len(preOrder) == 0 {
		return nil, errors.New("submit order fialded")
	}
	if preOrder["sellable_id"] != sellableResult["origin_code"] {
		return nil, errors.New("submit order fialded")
	}
	return sellableResult, nil
}

func getSellableInfo(cliCtx context.CLIContext, storeName, appcode, sellableid string ) (map[string]string, error){
	tableName := "sellable"
	fieldName := "code"
	value := sellableid
	ac := getOracleAc()
	res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/find_by/%s/%s/%s/%s/%s", storeName, ac, appcode, tableName, fieldName, value), nil)
	if err != nil {
		fmt.Printf("could not find ids")
		return nil, err
	}

	var out types.QuerySliceOfString
	json.Unmarshal(res, &out)
	if len(out) < 1 {
		return nil, errors.New("could not find sellable id")
	}
	id := out[len(out) - 1]
	queryString := fmt.Sprintf("custom/%s/find/%s/%s/%s/%s", storeName, ac, appcode, tableName, id)
	return oracleQueryUserTable(cliCtx, queryString)
}

func getOrderReceiptInfo(cliCtx context.CLIContext, storeName string, fieldValue map[string]string) ([]map[string]string, error){

	querierObjs := []map[string]string{}
	var ent map[string]string
	ent = map[string]string{
		"method": "table",
		"table":  "order_receipt",
	}
	querierObjs = append(querierObjs, ent)

	for field, val := range fieldValue {
		ent = map[string]string{
			"method":   "where",
			"field":    field,
			"value":    val,
			"operator": "==",
		}
		querierObjs = append(querierObjs, ent)
	}
	return  querierQuery(cliCtx, storeName, "0000000001", querierObjs)
}

func querierQuery(cliCtx context.CLIContext, storeName , appCode string, querierObjs []map[string]string, ) ([]map[string]string, error){
	bz , err := json.Marshal(querierObjs)
	if err != nil {
		return nil, err
	}
	querierBase58 := base58.Encode(bz)
	privKey, err := oracle.LoadPrivKey()
	if err != nil {
		return nil, err
	}
	ac := utils.MakeAccessCode(privKey)
	orderReceiptInfo, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/querier/%s/%s/%s", storeName, ac, appCode, querierBase58), nil)
	if err != nil {
		return nil, err
	}

	orderReceipt := make([]map[string]string, 0)
	if string(orderReceiptInfo) == "[]"{
		return orderReceipt,nil
	}
	err = json.Unmarshal(orderReceiptInfo, &orderReceipt)
	if err != nil {
		return nil, err
	}
	return orderReceipt, nil
}

func getOrderInfo(cliCtx context.CLIContext, storeName, appcode, id string) (map[string]string, error) {
	ac := getOracleAc()
	queryString := fmt.Sprintf("custom/%s/find/%s/%s/%s/%s", storeName, ac, appcode, "order", id)
	return oracleQueryUserTable(cliCtx, queryString)
}

func getOracleAc() string {
	privKey, err := oracle.LoadPrivKey()
	if err != nil {
		return ""
	}
	ac := utils.MakeAccessCode(privKey)
	return ac
}

func oracleQueryUserTable(cliCtx context.CLIContext, query string) (map[string]string, error) {
	res, _, err := cliCtx.QueryWithData(query, nil)
	if err != nil {
		return nil, err
	}
	orderInfo := make(map[string]string)
	err = json.Unmarshal(res, &orderInfo)
	if err != nil {
		return nil, err
	}
	return orderInfo, err
}

func oraclePagePay(ReturnURL, Money , OutTradeNo string) (string, error) {
	var p = alipay.TradePagePay{}
	p.ReturnURL = ReturnURL
	p.OutTradeNo = OutTradeNo
	//TODO total amount need to calc
	p.TotalAmount = Money
	p.ProductCode = "FAST_INSTANT_TRADE_PAY"
	p.Subject = OrderSub
	p.NotifyURL = loadAlipayNotifyUrl()
	url, err := aliClient.TradePagePay(p)
	if err != nil {
		return "", err
	}
	return url.String(), nil
}

func oracleAppPay(Money , OutTradeNo string) (string, error) {
	var p = alipay.TradeAppPay{}
	p.NotifyURL = loadAlipayNotifyUrl()
	p.OutTradeNo = OutTradeNo
	p.TotalAmount = Money
	p.ProductCode = "FAST_INSTANT_TRADE_PAY"
	p.Subject = OrderSub
	result, err := aliClient.TradeAppPay(p)
	if err != nil {
		return "", err
	}
	return result, nil
}


func oracleQueryPayStatus(cliCtx context.CLIContext, storeName string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		outTradeNo := r.Form.Get("out_trade_no")
		outTradeNo = strings.TrimSpace(outTradeNo)
		//out_trade_no format appcode-id
		info := strings.Split(outTradeNo,"-")
		if len(info)  != 2 {
			rest.WriteErrorResponse(w, http.StatusBadRequest, errors.New("out_trade_no format err").Error())
			return
		}
		appcode := info[0]
		orderid := info[1]

		fieldValue := map[string]string{
			"orderid": orderid,
			"appcode": appcode,
		}
		orderReceipt, err := getOrderReceiptInfo(cliCtx, storeName, fieldValue)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		if len(orderReceipt) != 0 {
			bz, _ := json.Marshal(orderReceipt[0])
			rest.PostProcessResponse(w, cliCtx, bz)
			return
		}
		//query from ali
		aliOrderStatus, err := OracleQueryAliOrder(outTradeNo)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		if aliOrderStatus["trade_status"] != string(alipay.TradeStatusSuccess) {
			rest.WriteErrorResponse(w, http.StatusBadRequest, aliOrderStatus["trade_status"])
			return
		}
		// save to OrderReceipt table
		total_amount  := aliOrderStatus["total_amount"]
		trade_no := aliOrderStatus["trade_no"]
		res := newOrderReceiptData(cliCtx, storeName, outTradeNo, total_amount, trade_no)

		oracleAccAddr := oracle.GetOracleAccAddr()
		SaveToOrderInfoTable(cliCtx, oracleAccAddr, res, OrderReceipt)
		bz , _ := json.Marshal(res)
		rest.PostProcessResponse(w, cliCtx, bz)
	}
}


func oracleSavePayStatus(cliCtx context.CLIContext, storeName string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {


		r.ParseForm()
		if _, err := aliClient.VerifySign(r.Form); err != nil {
			fmt.Println("aliClient.VerifySign  err : ", err.Error())
			w.Write([]byte("failed"))
			return
		}
		outTradeNo := strings.TrimSpace(r.Form.Get("out_trade_no"))
		total_amount  := strings.TrimSpace(r.Form.Get("total_amount"))
		trade_no := strings.TrimSpace(r.Form.Get("trade_no"))
		fmt.Println("===> outTradeNo: ", outTradeNo, "===> trade_no: ", trade_no, "===>")
		res := newOrderReceiptData(cliCtx, storeName, outTradeNo, total_amount, trade_no)
		oracleAccAddr := oracle.GetOracleAccAddr()
		SaveToOrderInfoTable(cliCtx, oracleAccAddr, res, OrderReceipt)
		w.Write([]byte("success"))
	}
}

func newOrderReceiptData(cliCtx context.CLIContext, storeName , out_trade_no, total_amount, trade_no string)map[string]string{
	outTradeNo := out_trade_no
	info := strings.Split(outTradeNo,"-")
	//TODO GetOwner
	orderInfo , err := getOrderInfo(cliCtx, storeName, info[0], info[1])
	owner := ""
	if err != nil {
		owner = ""
	} else {
		owner = orderInfo["created_by"]
	}
	expiration_date := calcExpirationDate(cliCtx, storeName, info[0], owner ,orderInfo["sellable_id"])
	res := newOrderReceiptDataCore(info[0],info[1], owner, total_amount,expiration_date, AliPay, trade_no)
	return res
}

func newOrderReceiptDataCore(appcode, orderid, owner, amount, expiration_date, vendor,  vendor_payment_no string) map[string]string{
	res := make(map[string]string)
	res["appcode"] = appcode
	res["orderid"] = orderid
	res["owner"] = owner
	res["amount"]  = amount
	res["expiration_date"] = expiration_date
	res["vendor"]  = vendor
	res["vendor_payment_no"] = vendor_payment_no
	return res

}

func calcExpirationDate(cliCtx context.CLIContext, storeName string, appcode ,owner ,sellableid string) string{
	//当前购买套餐
	sellableInfo ,err := getSellableInfo(cliCtx, storeName, appcode, sellableid)
	if err != nil || sellableInfo["term_days"] == ""{
		return ""
	}
	fieldValues := map[string]string {
		"appcode" : appcode,
		"owner"   : owner,
	}
	//购买过的套餐
	Purchaseds, err := getOrderReceiptInfo(cliCtx, storeName, fieldValues)
	if err != nil {
		//if something unknown happend , set time from now
		//TODO
		return ""
	}
	if len(Purchaseds) == 0 {
		termDays := sellableInfo["term_days"]
		addDays, _ := strconv.Atoi(termDays)
		t := time.Now()
		t = t.Add(time.Hour * 24 * time.Duration(addDays))
		return  fmt.Sprintf("%d", t.UnixNano()/1000000)
	}

	lastOrderInfo, err := getOrderInfo(cliCtx, storeName, appcode, Purchaseds[len(Purchaseds) -1 ]["orderid"])
	if err != nil {
		return ""
	}
	termDays := sellableInfo["term_days"]
	addDays, _ := strconv.Atoi(termDays)
	//升级套餐
	if lastOrderInfo["sellable_id"] == sellableInfo["origin_code"] && lastOrderInfo["sellable_id"] != sellableid {
		t := time.Now()
		t = t.Add(time.Hour * 24 * time.Duration(addDays))
		return  fmt.Sprintf("%d", t.UnixNano()/1000000)
	}
	//续费套餐
	expirationDate := Purchaseds[len(Purchaseds) -1 ]["expiration_date"]
	timeSatmep, err := strconv.ParseInt(expirationDate, 10,64)
	if err != nil {
		return ""
	}
	t := time.Unix(timeSatmep/1000,0)
	t = t.Add(time.Hour * 24 * time.Duration(addDays))
	return fmt.Sprintf("%d", t.UnixNano()/1000000)
}

func OracleQueryAliOrder(outTradeNo string) (map[string]string, error){
	p := alipay.TradeQuery{
		OutTradeNo : outTradeNo,
	}
	query, err := aliClient.TradeQuery(p)
	if err != nil {
		return nil, err
	}
	result := make(map[string]string)
	result["out_trade_no"] = query.Content.OutTradeNo
	result["trade_no"] = query.Content.TradeNo
	result["total_amount"] = query.Content.TotalAmount
	result["trade_status"] = string(query.Content.TradeStatus)
	return result, nil
}

func SaveToOrderInfoTable(cliCtx context.CLIContext, oracleAddr sdk.AccAddress,  row map[string]string, tableName string) error{
	//write to buyerorder table
	rowFields := make(types.RowFields)
	for k, v := range row {
		rowFields[k] = v
	}

	js,err := json.Marshal(rowFields)
	if err != nil {
		return err
	}
	msg := types.NewMsgInsertRow(oracleAddr, "0000000001", tableName, js)
	err = msg.ValidateBasic()
	if err != nil {
		return err
	}
	oracle.BuildTxsAndBroadcast(cliCtx, []oracle.UniversalMsg{msg})
	return nil
}

func loadAliPrivateKey(path string) string {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return ""
	}
	return string(b)
}

func loadAlipayNotifyUrl() string{
	if notifyUrl != ""{
		return notifyUrl
	}
	notifyUrl = viper.GetString(AlipayNotifyURL)
	notifyUrl += "dbchain/oracle/dbcpay_notify"
	return notifyUrl
}
