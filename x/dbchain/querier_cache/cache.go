package querier_cache

import (
    "fmt"
    "errors"
    "strings"
    "sort"
    sdk "github.com/dbchaincloud/cosmos-sdk/types"
)

const (
    stringTrue  = "true"
    stringFalse = "false"
)

func GetIsTablePublic(appId uint, tableName string) (bool, error) {
    key := getIsTablePublicKey(appId, tableName)
    result, err := get([]byte(key))
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
    key := getIsTablePublicKey(appId, tableName)
    var resultBytes []byte
    if result {
        resultBytes = []byte(stringTrue)
    } else {
        resultBytes = []byte(stringFalse)
    }
    return set([]byte(key), resultBytes, expiration)
}

func VoidIsTablePublic(appId uint, tableName string) {
    key := getIsTablePublicKey(appId, tableName)
    del([]byte(key))
}

func GetIdsBy(address sdk.AccAddress, appId uint, tableName, fieldName, value string, isTablePublic func(uint, string) bool) ([]byte, error) {
    result := isTablePublic(appId, tableName)
    var key string
    if result {
        key = getIdsByKey0(appId, tableName, fieldName, value)
    } else {
        key = getIdsByKey1(address, appId, tableName, fieldName, value)
    }
    return get([]byte(key))
}

func SetIdsBy(address sdk.AccAddress, appId uint, tableName, fieldName, value string, isTablePublic func(uint, string) bool, toBeSaved[]byte) (error) {
    result := isTablePublic(appId, tableName)
    var key string
    if result {
        key = getIdsByKey0(appId, tableName, fieldName, value)
    } else {
        key = getIdsByKey1(address, appId, tableName, fieldName, value)
    }

    RegisterKeysOfTable(appId, tableName, key)
    return set([]byte(key), []byte(toBeSaved), expiration)
}

func GetFind(address sdk.AccAddress, appId uint, tableName, rowId string, isTablePublic func(uint, string) bool) ([]byte, error) {
    result := isTablePublic(appId, tableName)
    var key string
    if result {
        key = getFindKey0(appId, tableName, rowId)
    } else {
        key = getFindKey1(address, appId, tableName, rowId)
    }
    return get([]byte(key))
}

func SetFind(address sdk.AccAddress, appId uint, tableName, rowId string, isTablePublic func(uint, string) bool, toBeSaved []byte) (error) {
    result := isTablePublic(appId, tableName)
    var key string
    if result {
        key = getFindKey0(appId, tableName, rowId)
    } else {
        key = getFindKey1(address, appId, tableName, rowId)
    }
    // No need to be invalidated when table inserted or row frozen
    return set([]byte(key), toBeSaved, expiration * 10)
}

func GetQuerier(address sdk.AccAddress, appId uint, querierObjs [](map[string]string), isTablePublic func(uint, string) bool) ([]byte, error) {
    tableName, err := findTableFromQuerierObjects(querierObjs)
    if err != nil {
        return []byte(""), nil
    }

    querierStr := querierObjectsToString(querierObjs)

    result := isTablePublic(appId, tableName)
    var key string
    if result {
        key = getQuerierKey0(appId, tableName, querierStr)
    } else {
        key = getQuerierKey1(address, appId, tableName, querierStr)
    }
    return get([]byte(key))
}

func SetQuerier(address sdk.AccAddress, appId uint, querierObjs [](map[string]string), isTablePublic func(uint, string) bool, toBeSaved []byte) (error) {
    tableName, err :=findTableFromQuerierObjects(querierObjs)
    if err != nil {
        return err
    }

    querierStr := querierObjectsToString(querierObjs)

    result := isTablePublic(appId, tableName)
    var key string
    if result {
        key = getQuerierKey0(appId, tableName, querierStr)
    } else {
        key = getQuerierKey1(address, appId, tableName, querierStr)
    }

    RegisterKeysOfTable(appId, tableName, key)
    return set([]byte(key), toBeSaved, expiration)
}

//////////////////////
//                  //
// Helper functions //
//                  //
//////////////////////

func getIsTablePublicKey(appId uint, tableName string) string {
    return fmt.Sprintf("GetIsTablePublic:%d:%s", appId, tableName)
}

func getIdsByKey0(appId uint, tableName, fieldName, value string) string {
    return fmt.Sprintf("GetIdsBy0:%d:%s:%s:%s", appId, tableName, fieldName, value)
}

func getIdsByKey1(address sdk.AccAddress, appId uint, tableName, fieldName, value string) string {
    return fmt.Sprintf("GetIdsBy1:%s:%d:%s:%s:%s", address.String(), appId, tableName, fieldName, value)
}

func getFindKey0(appId uint, tableName, rowId string) string {
    return fmt.Sprintf("GetFind0:%d:%s:%s", appId, tableName, rowId)
}

func getFindKey1(address sdk.AccAddress, appId uint, tableName, rowId string) string {
    return fmt.Sprintf("GetFind1:%s:%d:%s:%s", address.String(), appId, tableName, rowId)
}

func getQuerierKey0(appId uint, tableName, querierStr string) string {
    return fmt.Sprintf("GetQuerier0:%d:%s:%s", appId, tableName, querierStr)
}

func getQuerierKey1(address sdk.AccAddress, appId uint, tableName, querierStr string) string {
    return fmt.Sprintf("GetQuerier1:%s:%d:%s:%s", address.String(), appId, tableName, querierStr)
}

func findTableFromQuerierObjects(querierObjs [](map[string]string)) (string, error) {
    for _, item := range querierObjs {
        if item["method"] == "table" {
            return item["table"], nil
        }
    }
    return "", errors.New("no table found in the querier")
}

func querierObjectsToString(querierObjs [](map[string]string)) string {
    result := []string{}
    for _, item := range querierObjs {
        str := mapToString(item)
        result = append(result, str)
    }
    sort.Strings(result)
    return strings.Join(result, ":")
}

func mapToString(input map[string]string) string {
    result := []string{}
    for key, element := range input {
       result = append(result, key + ":" + element)
    }
    sort.Strings(result)
    return strings.Join(result, ":")
}
