
///////////////////////////////////////////////////////////////////////////////////
//                                                                               //
//  this file is used to verify apple pay and save data to receipt_order table   //
//	verification process :                                                       //
//  1. get receipt-data and orderInfo from client                                //
//  2. send encoding receipt-data to apple server                                //
//  3. get ProductId and TransactionId from data of apple response               //
//  4. verify ProductId and orderInfo to make they are the same                  //
//  5. build tx to save transaction info to receipt_order table                  //
//                                                                               //
///////////////////////////////////////////////////////////////////////////////////

package oracle

import (
	"bytes"
	"encoding/json"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/yzhanginwa/dbchain/x/dbchain/client/oracle/oracle"
	"net/http"
	"strconv"
	"strings"
)

type ReceiptFormat1 struct {
	OriginalPurchaseDatePst string `json:"original_purchase_date_pst"`
	PurchaseDateMs          string `json:"purchase_date_ms"`
	UniqueIdentifier        string `json:"unique_identifier"`
	OriginalTransactionId   string `json:"original_transaction_id"`
	Bvrs                    string `json:"bvrs"`
	TransactionId           string `json:"transaction_id"`
	Quantity                string `json:"quantity"`
	UniqueVendorIdentifier  string `json:"unique_vendor_identifier"`
	ItemId                  string `json:"item_id"`
	ProductId               string `json:"product_id"`
	PurchaseDate            string `json:"purchase_date"`
	OriginalPurchaseDate    string `json:"original_purchase_date"`
	PurchaseDatePst         string `json:"purchase_date_pst"`
	Bid                     string `json:"bid"`
	OriginalPurchaseDateMs  string `json:"original_purchase_date_ms"`
}

type ReceiptFormat2 struct {
	ReceiptType                string        `json:"receipt_type"`
	AdamId                     int           `json:"adam_id"`
	AppItemId                  int           `json:"app_item_id"`
	BundleId                   string        `json:"bundle_id"`
	ApplicationVersion         string        `json:"application_version"`
	DownloadId                 int           `json:"download_id"`
	VersionExternalIdentifier  int           `json:"version_external_identifier"`
	ReceiptCreationDate        string        `json:"receipt_creation_date"`
	ReceiptCreationDateMs      string        `json:"receipt_creation_date_ms"`
	ReceiptCreationDatePst     string        `json:"receipt_creation_date_pst"`
	RequestDate                string        `json:"request_date"`
	RequestDateMs              string        `json:"request_date_ms"`
	RequestDatePst             string        `json:"request_date_pst"`
	OriginalPurchaseDate       string        `json:"original_purchase_date"`
	OriginalPurchaseDateMs     string        `json:"original_purchase_date_ms"`
	OriginalPurchaseDatePst    string        `json:"original_purchase_date_pst"`
	OriginalApplicationVersion string        `json:"original_application_version"`
	InApp                      []InAppFormat `json:"in_app"`
}

type InAppFormat struct {
	Quantity                string `json:"quantity"`
	ProductId               string `json:"product_id"`
	TransactionId           string `json:"transaction_id"`
	OriginalTransactionId   string `json:"original_transaction_id"`
	PurchaseDate            string `json:"purchase_date"`
	PurchaseDateMs          string `json:"purchase_date_ms"`
	PurchaseDatePst         string `json:"purchase_date_pst"`
	OriginalPurchaseDate    string `json:"original_purchase_date"`
	OriginalPurchaseDateMs  string `json:"original_purchase_date_ms"`
	OriginalPurchaseDatePst string `json:"original_purchase_date_pst"`
	IsTrialPeriod           string `json:"is_trial_period"`
}

type AppleResponse1 struct {
	Receipt ReceiptFormat1 `json:"receipt"`
	Status  int            `json:"status"`
}

type AppleResponse2 struct {
	Receipt     ReceiptFormat2 `json:"receipt"`
	Environment string         `json:"environment"`
	Status      int            `json:"status"`
}

const (
	SUCCESS = 0
)

var sandUrl = "https://sandbox.itunes.apple.com/verifyReceipt"
var productUrl = "https://buy.itunes.apple.com/verifyReceipt"

func verifyApplePay(cliCtx context.CLIContext, storeName, outTradeNo, buyer string, receiptData string) (string,bool) {
	//get transaction info from apple server
	cli := http.Client{}
	ReceiptData := map[string]string{"receipt-data": receiptData}
	bz, _ := json.Marshal(ReceiptData)
	contentType := "application/json; charset=utf-8"
	reader := bytes.NewReader(bz)
	postUrl := ""
	if IsTest {
		postUrl = sandUrl
	} else {
		postUrl = productUrl
	}
	resp, err := cli.Post(postUrl, contentType, reader)
	if err != nil {
		return err.Error(), false
	}
	buffer := make([]byte, 0)
	buf := make([]byte, 4096)
	for {
		//n maybe less than 4096 when err is nil, so we should check it everytime
		n, err := resp.Body.Read(buf)
		if err != nil {
			if n < 4096 {
				buffer = append(buffer, buf[:n]...)
			} else {
				buffer = append(buffer, buf...)
			}
			break
		}
		if n < 4096 {
			buffer = append(buffer, buf[:n]...)
		} else {
			buffer = append(buffer, buf...)
		}

	}

	//decode data witch received from apple server
	var ResponseFormat1 AppleResponse1
	var ResponseFormat2 AppleResponse2
	var appleTransactionId string
	var productId string

	json.Unmarshal(buffer, &ResponseFormat2)
	if len(ResponseFormat2.Receipt.InApp) > 0 {
		if ResponseFormat2.Status != SUCCESS {
			return "ResponseFormat2.Status != SUCCESS", false
		}
		lastInApp := getLastestPurchase(ResponseFormat2.Receipt.InApp)
		productId = lastInApp.ProductId
		appleTransactionId = lastInApp.TransactionId
	} else {
		json.Unmarshal(buffer, &ResponseFormat1)
		if ResponseFormat1.Status != SUCCESS {
			return "ResponseFormat1.Status != SUCCESS", false
		}
		productId = ResponseFormat1.Receipt.ProductId
		appleTransactionId = ResponseFormat1.Receipt.TransactionId
	}

	//to make sure product of payed is same with order table
	appCodeAndOrderId := strings.Split(outTradeNo,"-")
	if len(appCodeAndOrderId) != 2 {
		return "appCodeAndOrderId len err", false
	}
	orderInfo , err := getOrderInfo(cliCtx, storeName, appCodeAndOrderId[0], appCodeAndOrderId[1])
	if err != nil {
		return err.Error(), false
	}
	if orderInfo["sellable_id"] != productId || orderInfo["created_by"] != buyer {
		return "orderInfo[\"sellable_id\"] != productId || orderInfo[\"created_by\"] != buyer", false
	}
	//to make sure there is no this orderid in order_receipt
	fieldValue := map[string]string{
		"orderid": appCodeAndOrderId[1],
		"appcode": appCodeAndOrderId[0],
	}
	orderReceipt, err := getOrderReceiptInfo(cliCtx, storeName, fieldValue)
	if len(orderReceipt) != 0 {
		return "len(orderReceipt) != 0", false
	}
	return appleTransactionId, true
}

func callDbcApplePay(cliCtx context.CLIContext, storeName, outTradeNo,  appleTransactionId string) ([]byte, error){
	appCodeAndOrderId := strings.Split(outTradeNo,"-")
	appcode := appCodeAndOrderId[0]
	orderId := appCodeAndOrderId[1]
	orderInfo, err  := getOrderInfo(cliCtx, storeName, appcode, orderId)
	if err != nil {
		return nil, err
	}
	sellableId := orderInfo["sellable_id"]
	sellableInfo, err := getSellableInfo(cliCtx, storeName, appcode, sellableId)
	if err != nil {
		return  nil, err
	}
	//save to order_receipt
	amount := sellableInfo["price"]
	owner := orderInfo["created_by"]
	expiration_date := calcExpirationDate(cliCtx, storeName, appcode, owner ,orderInfo["sellable_id"])

	res := newOrderReceiptDataCore(appcode, orderId, owner, amount, expiration_date, ApplePay, appleTransactionId)
	oracleAccAddr := oracle.GetOracleAccAddr()
	SaveToOrderInfoTable(cliCtx, oracleAccAddr, res, OrderReceipt)
	bz , _ := json.Marshal(res)
	return bz, nil
}

func getLastestPurchase(InAppDatas []InAppFormat) InAppFormat {
	index := -1
	var timeStamp int64 = -1
	for i, InApp := range InAppDatas {
		temp := InApp.PurchaseDateMs
		tempTimeStamp , err := strconv.ParseInt(temp, 10, 64)
		if err != nil {
			continue
		}
		if tempTimeStamp > timeStamp {
			timeStamp = tempTimeStamp
			index = i
		}
	}
	if index == -1 {
		return InAppFormat{}
	}
	return InAppDatas[index]
}