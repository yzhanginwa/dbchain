package keeper

import (
    "fmt"
)

const (
    KeyPrefixMeta  = "mt"
    KeyPrefixData  = "dt"
    KeyPrefixIndex = "ix"
)

func getNextIdKey(tableName string) string {
    return fmt.Sprintf("%s:nextId:%s", KeyPrefixMeta, tableName)
}

func getTableKey(tableName string) string {
    return fmt.Sprintf("%s:tn:%s", KeyPrefixMeta, tableName)
}

func getDataKey(tableName string, id uint, fieldName string) string {
    return fmt.Sprintf("%s:%s:%d:%s", KeyPrefixData, tableName, id, fieldName)
}

