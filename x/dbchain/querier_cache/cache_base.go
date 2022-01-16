package querier_cache

import (
    "fmt"
    "time"
    "math/rand"
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
    extraWait = 500                     // 500 milliseconds
)

var (
    theCache *freecache.Cache 
    theChannel = make(chan tableId, 50)    
    notificationBufferMap = make(map[uint]map[string]uint)
    myRand = rand.New(rand.NewSource(time.Now().UnixNano()))
)

func init() {
    theCache = freecache.NewCache(theCacheSize)
    debug.SetGCPercent(20)
    go handleTableExpiration()
}

// this function alleviates the "cache miss storm" problem
// by using random expiration time
func set(key, value []byte, expireSeconds int) (err error) {
    var realExpiration = 5
    if expireSeconds > 5 {
        realExpiration = expiration + myRand.Intn(int(expiration * 0.2))
    }
    return theCache.Set([]byte(key), value, realExpiration)
}

// this function elimiates the "cache miss storm" problem in a different way.
//
// when some highly requested concurrent queries miss the cache, they would
// make the same requests to query against the system before the first request
// returns and populates the cache. When this happens, the system may become
// very busy or unresponsive.
//
// This function lets the first request to get real data from system and pululate
// the cache. In the meantime it lets the other requests to wait until the
// cache is populated.
func get(key []byte) (value []byte, err error) {
    val1, err1 :=  theCache.Get(key)
    if err1 == nil {
        return val1,  err1
    }

    guardKey := []byte(fmt.Sprintf("gk_%s", key))
    for i := 0; i < 4; i++  {
        _, err2 :=  theCache.Get(guardKey)
        if err2!= nil {
            theCache.Set(guardKey, []byte("ANY"), 2)
            return val1, err1
        }
        time.Sleep(500 * time.Millisecond)
        val3, err3 := theCache.Get(key)
        if err3 == nil {
            return val3, err3
        }
    }

    return val1, err1
}

func del(key []byte) (affected bool) {
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
    if appId == 0 {
        theChannel <- tableId{0, ""}
        for k1 := range notificationBufferMap {
            v1 := notificationBufferMap[k1]
            for k2 := range v1 {
                theChannel <- tableId{k1, k2}         // k1: appId, k2: tableName 
                delete(v1, k2)
            }
            delete(notificationBufferMap, k1)
        }
    } else {
        if v, found := notificationBufferMap[appId]; found {
            v[tableName] = 1
        } else {
            notificationBufferMap[appId] = map[string]uint{tableName: 1}
        }
    }
}

func handleTableExpiration() {
    for {
        oneTableId := <-theChannel
        appId := oneTableId.AppId
        if appId == 0 {
            time.Sleep(extraWait * time.Millisecond)
            continue
        }
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
