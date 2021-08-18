package oracle

import (
	"math"
	"sync"
	"time"
)


////////////////////////
//                    //
//      order set     //
//                    //
////////////////////////

type nftOrderSet struct {
	orders map[int]int64//sync.Map
	t *time.Ticker
	lock sync.RWMutex
}

func newNftOrderSet(duration time.Duration) *nftOrderSet {
	n := &nftOrderSet {
		t : time.NewTicker(duration),
		orders: make(map[int]int64),
	}
	return n
}

func (n *nftOrderSet) Set(maxSize int) bool {
	n.lock.Lock()
	defer n.lock.Unlock()
	if len(n.orders) >= maxSize {
		return false
	}
	index := len(n.orders) % math.MaxInt32
	n.orders[index] = time.Now().Unix()
	return true
}

func (n * nftOrderSet) Delete(orderId int) {
	n.lock.Lock()
	defer n.lock.Unlock()
	delete(n.orders, orderId)
}

func (n * nftOrderSet)GC() {
	for {
		select {
		case _ = <-n.t.C :
			nowTime := time.Now().Unix()
			n.lock.Lock()
			for orderId, orderTime := range n.orders {
				if nowTime - orderTime > 300 {
					delete(n.orders, orderId)
				}
			}
			n.lock.Unlock()
		}
	}
}

func (n *nftOrderSet) Size() int {
	n.lock.RLock()
	defer n.lock.RUnlock()
	return len(n.orders)
}