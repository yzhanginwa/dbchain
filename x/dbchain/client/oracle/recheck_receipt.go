package oracle

import (
	"github.com/dbchaincloud/cosmos-sdk/client/context"
	"sync"
	"time"
)


type unVerifyAppleReceiptData struct {
	storeName string
	outTradeNo string
	buyer string
	receiptData string
}
type receiptPool struct {
	receipts []unVerifyAppleReceiptData
	lock sync.Mutex
}

func newReceiptPool() *receiptPool{
	return &receiptPool {
		receipts: make([]unVerifyAppleReceiptData, 0),
		lock: sync.Mutex{},
	}
}


func (r *receiptPool)ReadN(N int) []unVerifyAppleReceiptData {
	r.lock.Lock()
	defer r.lock.Unlock()
	var res []unVerifyAppleReceiptData
	length := len(r.receipts)
	if length < 1 {
		return nil
	} else if length <= N {
		res = r.receipts
		r.receipts = make([]unVerifyAppleReceiptData, 0)
	} else {
		res = r.receipts[:N]
		r.receipts = r.receipts[N:]
	}

	return res
}


func (r *receiptPool)Write(receipt unVerifyAppleReceiptData)  {
	r.lock.Lock()
	defer r.lock.Unlock()
	r.receipts = append(r.receipts, receipt)
}

var checkAppleUnverifyReceipt = false
var unVerifyAppleReceiptBuf = make(chan unVerifyAppleReceiptData, 100)
var recheckAppleReceiptBuf = make(chan unVerifyAppleReceiptData, 100)


func checkAppleReceiptRunner(cliCtx context.CLIContext) {
	receiptPoolCache := newReceiptPool()
	go func() {
		for {
			select {
			case receipt := <- unVerifyAppleReceiptBuf:
				appleTransactionId, applePayType, err := verifyApplePay(cliCtx, receipt.storeName, receipt.outTradeNo, receipt.buyer, receipt.receiptData)
				if err != nil {
					if err.Error() ==  "failed to access Apple server" {
						receiptPoolCache.Write(receipt)
					}
					break
				}
				callDbcApplePay(cliCtx, receipt.storeName, receipt.outTradeNo, appleTransactionId, applePayType)
			}
		}
	}()

	go func() {
		for {
			select {
			case receipt := <- recheckAppleReceiptBuf:
				appleTransactionId, applePayType, err := verifyApplePay(cliCtx, receipt.storeName, receipt.outTradeNo, receipt.buyer, receipt.receiptData)
				if err != nil  {
					if err.Error() ==  "failed to access Apple server"{
						receiptPoolCache.Write(receipt)
					}
					break
				}
				callDbcApplePay(cliCtx, receipt.storeName, receipt.outTradeNo, appleTransactionId, applePayType)
			default:
				needRecheckReceipts :=  receiptPoolCache.ReadN(100)
				if needRecheckReceipts == nil {
					time.Sleep(time.Minute)
				} else {
					for _, receipt := range needRecheckReceipts {
						recheckAppleReceiptBuf <- receipt
					}
				}
			}
		}
	}()

}