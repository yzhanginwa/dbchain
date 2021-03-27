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
table name : buyerorder
fields:
	buyer          //address of buyer
	out_trade_no   //商家订单号

table name : orderinfo
fields:
	out_trade_no	//商家订单号
	trade_no		//支付宝交易号
	total_amount	//交易的订单金额
	trade_status	//交易状态
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
 	IsTest = true
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

func oracleCallAliPagePay(cliCtx context.CLIContext, storeName string) http.HandlerFunc {
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
		OutTradeNo := r.Form.Get("out_trade_no")
		appcode := r.Form.Get("appcode")
		sellableid := r.Form.Get("sellableid")
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

func getOrderMoney(cliCtx context.CLIContext, storeName, appcode, sellableid, buyer string) (map[string]string,error) {

	sellableResult, err := getSellableInfo(cliCtx, storeName, appcode, sellableid)
	if err != nil {
		return nil, err
	}
	originId := sellableResult["originid"]
	if originId == "" {
		return sellableResult, nil
	}
	//check order validity
	//query 0000000001 order_receipt to conform originId exist
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
	if len(preOrder) != 1 {
		return nil, errors.New("submit order fialded")
	}
	if preOrder["sellable_id"] != sellableResult["originid"] {
		return nil, errors.New("submit order fialded")
	}
	return sellableResult, nil
}

func getSellableInfo(cliCtx context.CLIContext, storeName, appcode, sellableid string ) (map[string]string, error){
	tableName := "sellable"
	privKey, err := oracle.LoadPrivKey()
	if err != nil {
		return nil, err
	}
	ac := utils.MakeAccessCode(privKey)
	res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/find/%s/%s/%s/%s", storeName, ac, appcode, tableName, sellableid), nil)
	if err != nil {
		return nil, err
	}

	sellableResult := make(map[string]string)
	err = json.Unmarshal(res, &sellableResult)
	if err != nil {
		return nil, err
	}
	return sellableResult, nil
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
	privKey, err := oracle.LoadPrivKey()
	if err != nil {
		return nil, err
	}
	ac := utils.MakeAccessCode(privKey)
	res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/find/%s/%s/%s/%s", storeName, ac, appcode, "order", id), nil)
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

		fieldValue := map[string]string{
			"orderid": outTradeNo,
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
		SaveToOrderInfoTable(oracleAccAddr, res, OrderReceipt)
		bz , _ := json.Marshal(res)
		rest.PostProcessResponse(w, cliCtx, bz)
	}
}


func oracleSavePayStatus(cliCtx context.CLIContext, storeName string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		r.ParseForm()
		outTradeNo := r.Form.Get("out_trade_no")
		total_amount  := r.Form.Get("total_amount")
		trade_no := r.Form.Get("trade_no")
		res := newOrderReceiptData(cliCtx, storeName, outTradeNo, total_amount, trade_no)
		oracleAccAddr := oracle.GetOracleAccAddr()
		SaveToOrderInfoTable(oracleAccAddr, res, OrderReceipt)
	}
}

func newOrderReceiptData(cliCtx context.CLIContext, storeName , out_trade_no, total_amount, trade_no string)map[string]string{
	outTradeNo := out_trade_no
	info := strings.Split(outTradeNo,"-")

	res := make(map[string]string)
	res["appcode"] = info[0]
	res["orderid"] = info[1]
	//TODO GetOwner
	orderInfo , err := getOrderInfo(cliCtx, storeName, info[0], info[1])
	if err != nil {
		res["owner"] = ""
	} else {
		res["owner"] = orderInfo["created_by"]
	}

	res["amount"]  = total_amount
	res["expiration_date"] = calcExpirationDate(cliCtx, storeName, info[0], res["owner"],orderInfo["sellable_id"])
	res["vendor"]  = "alipay"
	res["vendor_payment_no"] = trade_no
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
	if lastOrderInfo["sellable_id"] == "1" && lastOrderInfo["sellable_id"] != sellableid {
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

func SaveToOrderInfoTable(oracleAddr sdk.AccAddress,  row map[string]string, tableName string) error{
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
	oracle.BuildTxsAndBroadcast([]oracle.UniversalMsg{msg})
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
