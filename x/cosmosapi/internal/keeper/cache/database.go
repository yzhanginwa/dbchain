package cache 

import (
    "github.com/yzhanginwa/cosmos-api/x/dbchain/internal/types"
)

////////////////////
//                //
// database cache //
//                //
////////////////////

var (
    database = make(map[string]types.Database)
    appIdToCode = make(map[uint]string)
)

func GetDatabase(appCode string) (types.Database, bool) {
    value, ok := database[appCode]
    return value, ok
}

func SetDatabase(appCode string, db types.Database) {
    database[appCode] = db
    appIdToCode[db.AppId] = appCode
}

func GetAppCodeById(appId uint) (string, bool) {
    value, ok := appIdToCode[appId]
    return value, ok
}

