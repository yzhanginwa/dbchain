package db

import (
    "fmtp"
    "sync"
    sdk "github.com/cosmos/cosmos-sdk/types"
    "github.com/yzhanginwa/cosmos-api/x/cosmosapi/internal/types"
)

func getNextIdKey(tableName) {
    return fmt.Sprintf("%s:nextId:%s", types.KeyPrefixMeta, tableName)
}

func getTableKey(tableName) {
    return fmt.Sprintf("%s:tn:%s", types.KeyPrefixMeta, tableName)
}

func getDataKey(tableName string, id uint, fieldName string) {
    return fmt.Sprintf("%s:%s:%d:%s", types.KeyPrefixData, tableName, id, fieldName)
}

