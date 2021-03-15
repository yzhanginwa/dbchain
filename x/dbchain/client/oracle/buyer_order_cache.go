package oracle

import "sync"

//
var buyerOrder sync.Map

func SetBuyerOrder(buyer ,order string){
	buyerOrder.Store(buyer,order)
}

func GetOrder(buyer string) string {
	res, ok := buyerOrder.Load(buyer)
	if !ok {
		return ""
	}
	return res.(string)
}

func DelOrder(buyer string) {
	buyerOrder.Delete(buyer)
	return
}