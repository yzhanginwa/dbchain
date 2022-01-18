package query_cache

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

func getIsTablePublic(appCode, tableName string) (bool, error) {
    key := getIsTablePublicKey(appCode, tableName)
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

func setIsTablePublic(appCode, tableName string, result bool) (error) {
    key := getIsTablePublicKey(appCode, tableName)
    var resultBytes []byte
    if result {
        resultBytes = []byte(stringTrue)
    } else {
        resultBytes = []byte(stringFalse)
    }
    return set([]byte(key), resultBytes, expiration)
}

func VoidIsTablePublic(appCode, tableName string) {
    key := getIsTablePublicKey(appCode, tableName)
    del([]byte(key))
}

func GetIdsBy(address sdk.AccAddress, appCode, tableName, fieldName, value string) ([]byte, error) {
    result := isTablePublic(appCode, tableName)
    var key string
    if result {
        key = getIdsByKey0(appCode, tableName, fieldName, value)
    } else {
        key = getIdsByKey1(address, appCode, tableName, fieldName, value)
    }
    return get([]byte(key))
}

func SetIdsBy(address sdk.AccAddress, appCode, tableName, fieldName, value string, toBeSaved[]byte) (error) {
    result := isTablePublic(appCode, tableName)
    var key string
    if result {
        key = getIdsByKey0(appCode, tableName, fieldName, value)
    } else {
        key = getIdsByKey1(address, appCode, tableName, fieldName, value)
    }

    RegisterKeysOfTable(appCode, tableName, key)
    return set([]byte(key), []byte(toBeSaved), expiration)
}

func GetFind(address sdk.AccAddress, appCode, tableName, rowId string) ([]byte, error) {
    result := isTablePublic(appCode, tableName)
    var key string
    if result {
        key = getFindKey0(appCode, tableName, rowId)
    } else {
        key = getFindKey1(address, appCode, tableName, rowId)
    }
    return get([]byte(key))
}

func SetFind(address sdk.AccAddress, appCode, tableName, rowId string, toBeSaved []byte) (error) {
    result := isTablePublic(appCode, tableName)
    var key string
    if result {
        key = getFindKey0(appCode, tableName, rowId)
    } else {
        key = getFindKey1(address, appCode, tableName, rowId)
    }
    // No need to be invalidated when table inserted or row frozen
    return set([]byte(key), toBeSaved, expiration * 10)
}

func GetQuerier(address sdk.AccAddress, appCode string, querierObjs [](map[string]string)) ([]byte, error) {
    tableName, err := findTableFromQuerierObjects(querierObjs)
    if err != nil {
        return []byte(""), nil
    }

    querierStr := querierObjectsToString(querierObjs)

    result := isTablePublic(appCode, tableName)
    var key string
    if result {
        key = getQuerierKey0(appCode, tableName, querierStr)
    } else {
        key = getQuerierKey1(address, appCode, tableName, querierStr)
    }
    return get([]byte(key))
}

func SetQuerier(address sdk.AccAddress, appCode string, querierObjs [](map[string]string), toBeSaved []byte) (error) {
    tableName, err :=findTableFromQuerierObjects(querierObjs)
    if err != nil {
        return err
    }

    querierStr := querierObjectsToString(querierObjs)

    result := isTablePublic(appCode, tableName)
    var key string
    if result {
        key = getQuerierKey0(appCode, tableName, querierStr)
    } else {
        key = getQuerierKey1(address, appCode, tableName, querierStr)
    }

    RegisterKeysOfTable(appCode, tableName, key)
    return set([]byte(key), toBeSaved, expiration)
}

//////////////////////
//                  //
// Helper functions //
//                  //
//////////////////////

func isTablePublic(appCode, tableName string) bool {
    return true
}

func getIsTablePublicKey(appCode, tableName string) string {
    return fmt.Sprintf("GetIsTablePublic:%s:%s", appCode, tableName)
}

func getIdsByKey0(appCode, tableName, fieldName, value string) string {
    return fmt.Sprintf("GetIdsBy0:%s:%s:%s:%s", appCode, tableName, fieldName, value)
}

func getIdsByKey1(address sdk.AccAddress, appCode, tableName, fieldName, value string) string {
    return fmt.Sprintf("GetIdsBy1:%s:%s:%s:%s:%s", address.String(), appCode, tableName, fieldName, value)
}

func getFindKey0(appCode, tableName, rowId string) string {
    return fmt.Sprintf("GetFind0:%s:%s:%s", appCode, tableName, rowId)
}

func getFindKey1(address sdk.AccAddress, appCode, tableName, rowId string) string {
    return fmt.Sprintf("GetFind1:%s:%s:%s:%s", address.String(), appCode, tableName, rowId)
}

func getQuerierKey0(appCode, tableName, querierStr string) string {
    return fmt.Sprintf("GetQuerier0:%s:%s:%s", appCode, tableName, querierStr)
}

func getQuerierKey1(address sdk.AccAddress, appCode, tableName, querierStr string) string {
    return fmt.Sprintf("GetQuerier1:%s:%s:%s:%s", address.String(), appCode, tableName, querierStr)
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
