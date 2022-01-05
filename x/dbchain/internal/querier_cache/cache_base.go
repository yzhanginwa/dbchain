package querier_cache

import (
    "runtime/debug"
    "github.com/coocood/freecache"
)

const (
    theCacheSize = 500 * 1024 * 1024    // 500 MB
    expiration = 300                    // 300 seconds
)

var (
    theCache *freecache.Cache 
)

func init() {
    theCache = freecache.NewCache(theCacheSize)
    debug.SetGCPercent(20)
}

func Set(key, value []byte, expireSeconds int) (err error) {
    return theCache.Set(key, value, expiration)
}

func Get(key []byte) (value []byte, err error) {
    return theCache.Get(key)
}

func Del(key []byte) (affected bool) {
    return theCache.Del(key)
}

