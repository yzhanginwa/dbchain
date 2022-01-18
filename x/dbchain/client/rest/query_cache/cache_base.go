package query_cache

import (
    "fmt"
    "time"
    "os"
    "syscall"
    "bufio"
    "strings"
    "math/rand"
    "encoding/binary"
    "runtime/debug"
    "github.com/coocood/freecache"
)

type tableId struct {
    AppCode string 
    TableName string
}

const (
    theCacheSize = 500 * 1024 * 1024    // 500 MB
    expiration = 300                    // 300 seconds
    extraWait = 500                     // 500 milliseconds
    fifoName = "/tmp/dbchain-fifo-1.pipe"
)

var (
    theCache *freecache.Cache 
    notificationBufferMap = make(map[string]map[string]uint)
    myRand = rand.New(rand.NewSource(time.Now().UnixNano()))
    fifoFile *os.File
    fifoReader *bufio.Reader
)


func init() {
    fmt.Println("create fifo in this executable")
    createFifo()
    // somehow the dbchaind includes this package and
    // the init() got invoked when it started.
    // we only want this goroutine to run in dbchaincli rest-server,
    // so we launch it in function "cacheInstance" to avoid the above issue
    //go handleTableExpiration()
}

func cacheInstance() *freecache.Cache {
    if theCache == nil {
        theCache = freecache.NewCache(theCacheSize)
        debug.SetGCPercent(20)
        emptyFifo()
        go handleTableExpiration()
    }
    return theCache
}

func createFifo() {
    if _, err := os.Stat(fifoName); err != nil {
        err := syscall.Mkfifo(fifoName, 0600)
        if err != nil {
            panic(fmt.Sprintf("Failed to create named pipe %s\n", fifoName))
        }
    }
}

func emptyFifo() {
    file, err := os.OpenFile(fifoName, syscall.O_RDONLY, 0)
    if err != nil {
        panic(fmt.Sprintf("Failed to open named pipie %s\n", fifoName))
    }

    buf := make([]byte, 10)
    for {
        n, err := file.Read(buf)
        if err != nil {
            break;
        }
        if n < 10 {
            break;
        }
    }
    defer file.Close()
}

func readFifo() (string, string) {
    if fifoFile == nil {
        f, err := os.OpenFile(fifoName, syscall.O_RDONLY, 0)
        if err != nil {
            panic(fmt.Sprintf("Failed to open named pipie %s\n", fifoName))
        }
        fifoFile = f
        fifoReader = bufio.NewReader(fifoFile)
    }

    line, err := fifoReader.ReadString('\n')
    if err != nil {
        // usually this happens when dbchaind is down
        time.Sleep(time.Second)
    }

    line = strings.TrimSuffix(line, "\n")
    s := strings.Split(line, ",")
    if len(s) == 2 {
        return s[0], s[1]
    }
    return "", ""
}


// this function alleviates the "cache miss storm" problem
// by using random expiration time
func set(key, value []byte, expireSeconds int) (err error) {
    theCache := cacheInstance()
    var realExpiration = 5
    if expireSeconds > 5 {
        realExpiration = expiration + myRand.Intn(int(expiration * 0.2))
    }
    return theCache.Set([]byte(key), value, realExpiration)
}

// this function elimiates some forms of the "cache miss storm" problem.
//
// when some highly requested concurrent queries miss the cache, they would
// make the same requests to query against the system before the first request
// returns and populates the cache. When this happens, the system may become
// very busy or unresponsive.
//
// This function lets the first request to get real data from system and polulate
// the cache. In the meantime it lets the other requests to wait until the
// cache is populated or timeout.
func get(key []byte) (value []byte, err error) {
    theCache := cacheInstance()
    val1, err1 :=  theCache.Get(key)
    if err1 == nil {
        return val1,  err1
    }

    guardKey := []byte(fmt.Sprintf("gk_%s", key))
    for i := 0; i < 8; i++  {
        _, err2 :=  theCache.Get(guardKey)
        if err2!= nil {
            theCache.Set(guardKey, []byte("ANY"), 4)
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
    theCache := cacheInstance()
    return theCache.Del(key)
}

func RegisterKeysOfTable(appCode, tableName string, dataKey string) (error) {
    theCache := cacheInstance()
    counterKey := getTableKeyCounterKey(appCode, tableName)
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

    tableKey := getTableKey(appCode, tableName, value)
    return theCache.Set([]byte(tableKey), []byte(dataKey), expiration)
}

func handleTableExpiration() {
    theCache := cacheInstance()
    for {
        appCode, tableName := readFifo()
        if appCode == "" {
            continue
        }
        if appCode == "_" {
            time.Sleep(extraWait * time.Millisecond)
            continue
        }
        counterKey := getTableKeyCounterKey(appCode, tableName)
        val, err := theCache.Get([]byte(counterKey))
        if err != nil {
            continue
        }
        val_uint32 := binary.LittleEndian.Uint32(val)
        theCache.Del([]byte(counterKey))

        for i := uint32(1); i <= val_uint32; i++ {
            tableKey := getTableKey(appCode, tableName, uint32(i))
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

func getTableKeyCounterKey(appCode, tableName string) string {
    return fmt.Sprintf("TableKeyCounter:%s:%s", appCode, tableName)
}

func getTableKey(appCode, tableName string, index uint32) string {
    return fmt.Sprintf("TableKey:%s:%s:%d", appCode, tableName, index)
}
