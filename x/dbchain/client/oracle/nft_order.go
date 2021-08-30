package oracle

import (
	"sync"
	"time"
)


////////////////////////
//                    //
//      order set     //
//                    //
////////////////////////

type nftOrderSet struct {
	ordersTimestamp map[string]map[int64]struct{}//sync.Map
	ordersDenom map[string]int
	t *time.Ticker
	lock sync.RWMutex
	invalidTime int64
}

func newNftOrderSet(duration time.Duration) *nftOrderSet {
	n := &nftOrderSet {
		t : time.NewTicker(duration),
		ordersTimestamp: make(map[string]map[int64]struct{}),
		ordersDenom : make(map[string]int),
		invalidTime : duration.Nanoseconds(),
	}
	return n
}

func (n *nftOrderSet) Set(denonId string, maxSize int) bool {
	n.lock.Lock()
	defer n.lock.Unlock()
	if maxSize <= 0 {
		return false
	}
	currentOrderSize, ok := n.ordersDenom[denonId]
	if !ok {
		temp := make(map[int64]struct{})
		nt := time.Now().UnixNano()
		//temp = append(temp, nt)
		temp[nt] = struct{}{}
		n.ordersTimestamp[denonId] = temp
		n.ordersDenom[denonId] = 1
		return true
	}

	if currentOrderSize >= maxSize {
		return false
	}
	currentOrderSize++
	n.ordersDenom[denonId] = currentOrderSize
	nt := time.Now().UnixNano()
	timeStamps := n.ordersTimestamp[denonId]
	timeStamps[nt] = struct{}{}
	return true
}

func (n * nftOrderSet) Delete(denomId string) {
	n.lock.Lock()
	defer n.lock.Unlock()
	delete(n.ordersDenom, denomId)
	delete(n.ordersTimestamp, denomId)
}

func (n * nftOrderSet)GC() {
	for {
		select {
		case _ = <-n.t.C :
			nowTime := time.Now().UnixNano()
			n.lock.Lock()

			for denomId, timeStamps := range n.ordersTimestamp {
				currentOrder := n.ordersDenom[denomId]
				for timeStamp, _ := range timeStamps {
						if nowTime - timeStamp > n.invalidTime {
							delete(timeStamps, timeStamp)
							currentOrder--
						}
				}
				if currentOrder <= 0 {
					delete(n.ordersDenom, denomId)
				} else {
					n.ordersDenom[denomId] = currentOrder
				}
			}
			n.lock.Unlock()
		}
	}
}

func (n *nftOrderSet) Size(denomId string) int {
	n.lock.RLock()
	defer n.lock.RUnlock()
	return n.ordersDenom[denomId]
}