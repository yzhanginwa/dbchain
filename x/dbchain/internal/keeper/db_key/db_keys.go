package db_key

import (
    "fmt"
    "strings"
    "github.com/yzhanginwa/dbchain/x/dbchain/internal/utils"
)

const (
    KeyPrefixAppCode    = "appcode"
    KeyPrefixAppNextId  = "appnextid"
    KeyPrefixDb       = "db"
    KeyPrefixSysGroup = "sysgrp"
    KeyPrefixFriend         = "friend"
    KeyPrefixPendingFriend  = "pending_friend"

    KeyPrefixUser  = "usr"
    KeyPrefixMeta  = "mt"
    KeyPrefixData  = "dt"
    KeyPrefixIndex = "ix"
    KeyPrefixGroup = "grp"
    KeyPrefixGroupMemo = "grp_mm"
    KeyPrefixGroups = "grps"
)

//////////////////////////////////
//                              //
// application/database related //
//                              //
//////////////////////////////////

func GetDatabaseKey(appCode string) string {
    return fmt.Sprintf("%s:%s", KeyPrefixAppCode, appCode)
}

func GetDatabaseNextIdKey() string {
    return fmt.Sprintf("%s", KeyPrefixAppNextId)
}

func GetDatabaseUserKey(appId uint, user string) string {
    return fmt.Sprintf("%s:%d:%s:%s", KeyPrefixDb, appId, KeyPrefixUser, user)
}

func GetDatabaseIteratorStartAndEndKey() (string, string) {
    start := fmt.Sprintf("%s:", KeyPrefixAppCode)
    end   := fmt.Sprintf("%s;", KeyPrefixAppCode)
    return start, end
}

func GetAppCodeFromDatabaseKey(key string) string {
    arr := strings.Split(key, ":")
    return arr[1]
}

func GetDatabaseUserIteratorStartAndEndKey(appId uint) (string, string) {
    start := fmt.Sprintf("%s:%d:%s:", KeyPrefixDb, appId, KeyPrefixUser)
    end   := fmt.Sprintf("%s:%d:%s;", KeyPrefixDb, appId, KeyPrefixUser)
    return start, end
}

func GetUserFromDatabaseUserKey(key string) string {
    arr := strings.Split(key, ":")
    return arr[3]
}

///////////////////
//               //
// table related //
//               //
///////////////////

// to store name of all tables
func GetTablesKey(appId uint) string {
    return fmt.Sprintf("%s:%d:%s:tables", KeyPrefixDb, appId, KeyPrefixMeta)
}

// to store the id for next new record of a table
func GetNextIdKey(appId uint, tableName string) string {
    return fmt.Sprintf("%s:%d:%s:nextId:%s", KeyPrefixDb, appId, KeyPrefixMeta, tableName)
}

// to store the name of fields for a table
func GetTableKey(appId uint, tableName string) string {
    return fmt.Sprintf("%s:%d:%s:tn:%s", KeyPrefixDb, appId, KeyPrefixMeta, tableName)
}

// to store table fields which have index on
func GetMetaTableIndexKey(appId uint, tableName string) string {
    return fmt.Sprintf("%s:%d:%s:idx:%s", KeyPrefixDb, appId, KeyPrefixMeta, tableName)
}

// to store the options for a table
func GetTableOptionsKey(appId uint, tableName string) string {
    return fmt.Sprintf("%s:%d:%s:opt:%s", KeyPrefixDb, appId, KeyPrefixMeta, tableName)
}

func GetColumnOptionsKey(appId uint, tableName string, fieldName string) string {
    return fmt.Sprintf("%s:%d:%s:fldopt:%s:%s", KeyPrefixDb, appId, KeyPrefixMeta, tableName, fieldName)
}

//////////////////////
//                  //
// function related //
//                  //
//////////////////////

func GetFunctionKey(appId uint, functionName string) string {
    return fmt.Sprintf("%s:%d:%s:functionInfo:%s", KeyPrefixDb, appId, KeyPrefixMeta, functionName)
}

// to store name of all funcS
func GetFunctionsKey(appId uint) string {
    return fmt.Sprintf("%s:%d:%s:functions", KeyPrefixDb, appId, KeyPrefixMeta)
}

//////////////////
//              //
// data related //
//              //
//////////////////

// to store the index data (ids) of an index field
func GetIndexKey(appId uint, tableName string, fieldName string, value string) string {
    return fmt.Sprintf("%s:%d:%s:%s:%s:%s", KeyPrefixDb, appId, KeyPrefixIndex, tableName, fieldName, value)
}

// to get the start and end parameters of iterator which seeks index data for certain field
func GetIndexDataIteratorStartAndEndKey(appId uint, tableName string, fieldName string) (string, string) {
    start := fmt.Sprintf("%s:%d:%s:%s:%s:", KeyPrefixDb, appId, KeyPrefixIndex, tableName, fieldName)
    end   := fmt.Sprintf("%s:%d:%s:%s:%s;", KeyPrefixDb, appId, KeyPrefixIndex, tableName, fieldName)
    return start, end
}

// to store the value of a fields on a record of a table.
func GetDataKeyBytes(appId uint, tableName string, fieldName string, id uint) []byte {
    base := fmt.Sprintf("%s:%d:%s:%s:%s:", KeyPrefixDb, appId, KeyPrefixData, tableName, fieldName)
    idBytes := utils.IntToByteArray(int64(id))
    return append([]byte(base), idBytes...)
}

// to get the start and end parameters of iterator which seeks certain value of a field
func GetFieldDataIteratorStartAndEndKey(appId uint, tableName string, fieldName string) (string, string) {
    start := fmt.Sprintf("%s:%d:%s:%s:%s:", KeyPrefixDb, appId, KeyPrefixData, tableName, fieldName)
    end   := fmt.Sprintf("%s:%d:%s:%s:%s;", KeyPrefixDb, appId, KeyPrefixData, tableName, fieldName)
    return start, end
}

func GetIdFromDataKey(key []byte) uint {
    length := len(key)
    if length < 8 {
        panic("key length cannot less than 8")
    }
    id := utils.ByteArrayToInt(key[length-8:])  // 8 = 64 / 8
    return uint(id)
}

// func getFieldNameFromDataKey(key string) string {
//     arr := strings.Split(key, ":")
//     return arr[4]
// }

////////////////////
//                //
// friend related //
//                //
////////////////////

func GetFriendKey(owner string, friendAddr string) string {
    return fmt.Sprintf("%s:%s:%s", KeyPrefixFriend, owner, friendAddr)
}

func GetFriendIteratorStartAndEndKey(owner string) (string, string) {
    start := fmt.Sprintf("%s:%s:", KeyPrefixFriend, owner)
    end   := fmt.Sprintf("%s:%s;", KeyPrefixFriend, owner)
    return start, end
}

func GetPendingFriendKey(owner string, friendAddr string) string {
    return fmt.Sprintf("%s:%s:%s", KeyPrefixPendingFriend, owner, friendAddr)
}

func GetPendingFriendIteratorStartAndEndKey(owner string) (string, string) {
    start := fmt.Sprintf("%s:%s:", KeyPrefixPendingFriend, owner)
    end   := fmt.Sprintf("%s:%s;", KeyPrefixPendingFriend, owner)
    return start, end
}

///////////////////
//               //
// group related //
//               //
///////////////////

func GetGroupsKey(appId uint) string {
    return fmt.Sprintf("%s:%d:%s", KeyPrefixDb, appId, KeyPrefixGroups)
}

func GetGroupKey(appId uint, groupName string) string {
    return fmt.Sprintf("%s:%d:%s:%s", KeyPrefixDb, appId, KeyPrefixGroup, groupName)
}

func GetGroupMemoKey(appId uint, groupName string) string {
    return fmt.Sprintf("%s:%d:%s:%s", KeyPrefixDb, appId, KeyPrefixGroupMemo, groupName)
}

func GetAdminGroupKey(appId uint) string {
    return GetGroupKey(appId, "admin")
}

//////////////////
//              //
// system level //
//              //
//////////////////

func GetSysGroupKey(groupName string) string {
    return fmt.Sprintf("%s:%s", KeyPrefixSysGroup, groupName)
}

func GetSysAdminGroupKey() string {
    return GetSysGroupKey("admin")
}
