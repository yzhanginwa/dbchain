package querier_cache

import (
    "fmt"
    "time"
    "encoding/binary"
    "runtime/debug"
    "github.com/coocood/freecache"
)

type tableId struct {
    AppId uint
    TableName string
}

const (
    theCacheSize = 500 * 1024 * 1024    // 500 MB
    expiration = 300                    // 300 seconds
)

var (
    theCache *freecache.Cache 
    theChannel = make(chan tableId, 50)    
)

func init() {
    theCache = freecache.NewCache(theCacheSize)
    debug.SetGCPercent(20)
    go handleTableExpiration()
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

func RegisterKeysOfTable(appId uint, tableName string, dataKey string) (error) {
    counterKey := getTableKeycounterKey(appId, tableName)
    val, err := theCache.Get([]byte(counterKey))
    var value uint32
    if err != nil {
        value = 0
    } else {
        value = binary.LittleEndian.Uint32(val)
    }
    
    value += 1
    bytes := make([]byte, 4)
    binary.LittleEndian.PutUint32(bytes, value)
    err = theCache.Set([]byte(counterKey), bytes, expiration)
    if err != nil {
        panic(fmt.Sprintf("Failed to set cache %s = %d", counterKey, value))
    }

    tableKey := getTableKey(appId, tableName, value)
    return theCache.Set([]byte(tableKey), []byte(dataKey), expiration)
}

func NotifyTableExpiration(appId uint, tableName string) {
    theChannel <- tableId{appId, tableName} 
}

func handleTableExpiration() {
    for {
        oneTableId := <-theChannel
        time.Sleep(2 * time.Second)     // just in case the message is in a bigger transaction
        appId := oneTableId.AppId
        tableName := oneTableId.TableName
        counterKey := getTableKeycounterKey(appId, tableName)
        val, err := theCache.Get([]byte(counterKey))
        if err != nil {
            continue
        }
        val_uint32 := binary.LittleEndian.Uint32(val)
        theCache.Del([]byte(counterKey))

        for i := uint32(1); i <= val_uint32; i++ {
            tableKey := getTableKey(appId, tableName, uint32(i))
            val1, err1 := theCache.Get([]byte(tableKey))
            if err1 == nil {
                theCache.Del(val1)
            }
            theCache.Del([]byte(tableKey))
        }
    }
}
    

//////////////////////
//                  //
// Helper functions //
//                  //
//////////////////////

func getTableKeycounterKey(appId uint, tableName string) string {
    return fmt.Sprintf("TableKeyCounter:%d:%s", appId, tableName)
}

func getTableKey(appId uint, tableName string, index uint32) string {
    return fmt.Sprintf("TableKey:%d:%s:%d", appId, tableName, index)
}
