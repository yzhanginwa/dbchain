package cache 

import (
    "errors"
    dbk "github.com/yzhanginwa/dbchain/x/dbchain/internal/keeper/db_key"
    "github.com/yzhanginwa/dbchain/x/dbchain/internal/types"
    "sync"
    "time"
)

////////////////////
//                //
// database cache //
//                //
////////////////////

var (
    database = make(map[string]types.Database)
    appIdToCode = make(map[uint]string)
    appTables = make(map[string]types.Table)
    TxStatusCache sync.Map
)


const (
    TxStateInvalidTime     = 600
    TxInvalidCheckRunTime  = 20 * time.Millisecond
    TxStateFail            = "fail"
    TxStateSuccess         = "success"
    TxStatePending         = "pending"
    TxStateProcessing      = "processing"
)

func GetDatabase(appCode string) (types.Database, bool) {
    value, ok := database[appCode]
    return value, ok
}

func VoidDatabase(appCode string) {
    delete(database, appCode)
}

func SetDatabase(appCode string, db types.Database) {
    database[appCode] = db
    appIdToCode[db.AppId] = appCode
}

func GetAppCodeById(appId uint) (string, bool) {
    value, ok := appIdToCode[appId]
    return value, ok
}

//operate appFields cash data
func GetTable(appId uint, tableName string)(types.Table, error){
    key := dbk.GetTableKey(appId, tableName)
    table, ok := appTables[key]
    if !ok {
        return types.Table{}, errors.New("appTables has no this table")
    }
    return table,nil
}

func SetTable(appId uint, tableName string, table types.Table){
    key := dbk.GetTableKey(appId, tableName)
    appTables[key] = table
}

func VoidTable(appId uint, tableName string){
    key := dbk.GetTableKey(appId, tableName)
    delete(appTables, key)
}
