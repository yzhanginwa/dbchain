package cache

import (
	"sync"
	"time"
)

//define a cache module to cache data in memory

type MemoryCache struct {
	data sync.Map	//key of map is used to save data, val of map is used to save timestamp
	ticker *time.Ticker
	expiration int64
}

func NewMemoryCache(ticker *time.Ticker, expiration int64) *MemoryCache {
	return &MemoryCache {
		ticker: ticker,
	}
}

func (m *MemoryCache) Set(key, val interface{}) {
	m.data.Store(key, val)
}

func (m *MemoryCache) Get(key interface{}) interface{}{
	data, _ := m.data.Load(key)
	return data
}

func (m *MemoryCache) Del(key interface{}) {
	m.data.Delete(key)
}

func (m *MemoryCache) Gc() {

	for {
		select {
		case <-m.ticker.C:
			nowTimeSecond := time.Now().Unix()
			m.data.Range(func(key, value interface{}) bool {
				saveTime, ok := value.(int64)
				if ok {
					if nowTimeSecond > saveTime + m.expiration {
						m.data.Delete(key)
					}
					return true
				}

				saveData , ok := value.(MakeNftInfo)
				if ok {
					if nowTimeSecond > saveData.TimeStamp + m.expiration {
						m.data.Delete(key)
					}
				}
				return true
			})
		}
	}
}

type MakeNftInfo struct {
	Data []byte
	TimeStamp int64
}