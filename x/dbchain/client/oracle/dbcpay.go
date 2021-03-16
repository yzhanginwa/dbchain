package oracle

import (
	"encoding/json"
	"fmt"
	"github.com/cosmos/cosmos-sdk/client/context"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/rest"
	"github.com/gorilla/mux"
	"github.com/smartwalle/alipay/v3"
	"github.com/smartwalle/xid"
	"github.com/yzhanginwa/dbchain/x/dbchain/client/oracle/oracle"
	"github.com/yzhanginwa/dbchain/x/dbchain/internal/types"
	"github.com/yzhanginwa/dbchain/x/dbchain/internal/utils"
	"io/ioutil"
	"log"
	"net/http"
	"os"
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


var aliClient *alipay.Client
const (
	kAppId = "2021000117616474"
	BuyerOrder = "buyerorder"
	OrderInfo  = "orderinfo"
	OrderSub   = "YTBox"
	NotifyURL  = "https://controlpanel.dbchain.cloud/relay/dbchain/oracle/dbcpay_notify"
	//it should be true when switch to production environment
 	IsProduction = false
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

		addr, err := utils.VerifyAccessCode(accessCode)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusNotFound, err.Error())
			return
		}
		ReturnURL := r.Form.Get("returnURL")
		Money := r.Form.Get("money")
		OutTradeNo := fmt.Sprintf("%d", xid.Next())


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
		//write to buyerorder table
		rowFields := make(types.RowFields)
		rowFields["buyer"] = addr.String()
		rowFields["out_trade_no"] = OutTradeNo
		js, err := json.Marshal(rowFields)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusNotFound, "marshal rowFields failed")
			return
		}
		oracleAccAddr := oracle.GetOracleAccAddr()
		msg := types.NewMsgInsertRow(oracleAccAddr, "0000000001", BuyerOrder, js)
		err = msg.ValidateBasic()
		if err != nil {
			rest.PostProcessResponse(w, cliCtx, "new save order msg err")
			return
		}
		oracle.BuildTxsAndBroadcast([]oracle.UniversalMsg{msg})
		//set to cache
		SetBuyerOrder(addr.String(), OutTradeNo)


		result := make(map[string]string)
		result["url"] = url
		result["out_trade_no"] = OutTradeNo
		bz , err := json.Marshal(result)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusNotFound, "marshal rowFields failed")
			return
		}
		rest.PostProcessResponse(w, cliCtx, bz)
		return
	}
}

func oraclePagePay(ReturnURL, Money , OutTradeNo string) (string, error) {
	var p = alipay.TradePagePay{}
	p.ReturnURL = ReturnURL
	p.OutTradeNo = OutTradeNo
	//TODO total amount need to calc
	p.TotalAmount = Money
	p.ProductCode = "FAST_INSTANT_TRADE_PAY"
	p.Subject = OrderSub
	p.NotifyURL = NotifyURL
	url, err := aliClient.TradePagePay(p)
	if err != nil {
		return "", err
	}
	return url.String(), nil
}

func oracleAppPay(Money , OutTradeNo string) (string, error) {
	var p = alipay.TradeAppPay{}
	p.NotifyURL = NotifyURL
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

func oracleQuerySubmitOrderStatus(cliCtx context.CLIContext, storeName string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		accessToken := vars["accessToken"]
		r.ParseForm()
		outTradeNo := r.Form.Get("out_trade_no")

		addr, err := utils.GetAddrFromAccessCode(accessToken)
		if err != nil {
			rest.WriteErrorResponse(w,  http.StatusNotFound, err.Error())
			return
		}
		order := GetOrder(addr.String())
		if order == "" {
			rest.WriteErrorResponse(w,  http.StatusNotFound, "Not found")
			return
		}
		oracleAdd := oracle.GetOracleAccAddr()
		res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/submit_order_status/%s/%s/%s", storeName, vars["accessToken"], outTradeNo, oracleAdd), nil)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusNotFound, err.Error())
			return
		}
		if string(res) == `"Success"`{
			DelOrder(addr.String())
		}
		rest.PostProcessResponse(w, cliCtx, res)
		return
	}
}

func oracleQueryPayStatus(cliCtx context.CLIContext, storeName string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		r.ParseForm()
		outTradeNo := r.Form.Get("out_trade_no")
		oracleAdd := oracle.GetOracleAccAddr()
		res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/order_status/%s/%s/%s", storeName, vars["accessToken"], outTradeNo, oracleAdd), nil)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusNotFound, err.Error())
			return
		}
		res = checkSaveToOrderInfo(res, oracleAdd)
		rest.PostProcessResponse(w, cliCtx, res)
	}
}

func checkSaveToOrderInfo(src []byte, oracleAddr sdk.AccAddress)[]byte {
	temp := make(map[string]string)
	err := json.Unmarshal(src, &temp)
	if err != nil {
		return src
	}
	if temp["SaveToTable"] != "true" {
		return src
	}
	delete(temp,"SaveToTable")
	SaveToOrderInfoTable(oracleAddr, temp, OrderInfo)
	bz , _ := json.Marshal(temp)
	return bz
}
func oracleSavePayStatus(cliCtx context.CLIContext, storeName string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		r.ParseForm()
		outTradeNo := r.Form.Get("out_trade_no")
		res,err := OracleQueryAliOrder(outTradeNo)
		if err != nil {
			return
		}

		// query table first
		priv, err := oracle.LoadPrivKey()
		if err != nil {
			return
		}
		ac := utils.MakeAccessCode(priv)
		oracleAdd := oracle.GetOracleAccAddr()
		bz, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/order_status/%s/%s/%s", storeName, ac, outTradeNo, oracleAdd), nil)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusNotFound, err.Error())
			return
		}
		if bz == nil || len(bz) == 0 {
			oracleAccAddr := sdk.AccAddress(priv.PubKey().Address())
			SaveToOrderInfoTable(oracleAccAddr, res, OrderInfo)
		}
	}
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

