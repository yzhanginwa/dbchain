package querier_cache

import (
    "fmt"
    sdk "github.com/dbchaincloud/cosmos-sdk/types"
)

const (
    stringTrue  = "true"
    stringFalse = "false"
)

func GetIsTablePublic(appId uint, tableName string) (bool, error) {
    key := fmt.Sprintf("GetIsTablePublic:%d:%s", appId, tableName)
    result, err := theCache.Get([]byte(key))
    if err != nil {
        return false, err
    }
    if string(result) == stringTrue {
        return true, nil
    } else {
        return false, nil
    }
}

func SetIsTablePublic(appId uint, tableName string, result bool) (error) {
    key := fmt.Sprintf("GetIsTablePublic:%d:%s", appId, tableName)
    var resultBytes []byte
    if result {
        resultBytes = []byte(stringTrue)
    } else {
        resultBytes = []byte(stringFalse)
    }
    return theCache.Set([]byte(key), resultBytes, expiration)
}

func GetIdsBy(address sdk.AccAddress, appId uint, tableName, fieldName, value string) ([]byte, error) {
    result, err := GetIsTablePublic(appId, tableName)
    var key string
    if err == nil && result {
        key = getIdsByKey0(appId, tableName, fieldName, value)
    } else {
        key = getIdsByKey1(address, appId, tableName, fieldName, value)
    }
    return theCache.Get([]byte(key))
}

func SetIdsBy(address sdk.AccAddress, appId uint, tableName, fieldName, value string, toBeSaved[]byte) (error) {
    result, err := GetIsTablePublic(appId, tableName)
    var key string
    if err == nil && result {
        key = getIdsByKey0(appId, tableName, fieldName, value)
    } else {
        key = getIdsByKey1(address, appId, tableName, fieldName, value)
    }

    RegisterKeysOfTable(appId, tableName, key)
    return theCache.Set([]byte(key), []byte(toBeSaved), expiration)
}

//////////////////////
//                  //
// Helper functions //
//                  //
//////////////////////

func getIdsByKey0(appId uint, tableName, fieldName, value string) string {
    return fmt.Sprintf("GetIdsBy0:%d:%s:%s:%s", appId, tableName, fieldName, value)
}

func getIdsByKey1(address sdk.AccAddress, appId uint, tableName, fieldName, value string) string {
    return fmt.Sprintf("GetIdsBy1:%s:%d:%s:%s:%s", address.String(), appId, tableName, fieldName, value)
}
