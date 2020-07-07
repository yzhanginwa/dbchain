package keeper

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
    KeyPrefixGroups = "grps"
)

//////////////////////////////////
//                              //
// application/database related //
//                              //
//////////////////////////////////

func getDatabaseKey(appCode string) string {
    return fmt.Sprintf("%s:%s", KeyPrefixAppCode, appCode)
}

func getDatabaseNextIdKey() string {
    return fmt.Sprintf("%s", KeyPrefixAppNextId)
}

func getDatabaseUserKey(appId uint, user string) string {
    return fmt.Sprintf("%s:%d:%s:%s", KeyPrefixDb, appId, KeyPrefixUser, user)
}

func getDatabaseIteratorStartAndEndKey() (string, string) {
    start := fmt.Sprintf("%s:", KeyPrefixAppCode)
    end   := fmt.Sprintf("%s;", KeyPrefixAppCode)
    return start, end
}

func getAppCodeFromDatabaseKey(key string) string {
    arr := strings.Split(key, ":")
    return arr[1]
}

func getDatabaseUserIteratorStartAndEndKey(appId uint) (string, string) {
    start := fmt.Sprintf("%s:%d:%s:", KeyPrefixDb, appId, KeyPrefixUser)
    end   := fmt.Sprintf("%s:%d:%s;", KeyPrefixDb, appId, KeyPrefixUser)
    return start, end
}

func getUserFromDatabaseUserKey(key string) string {
    arr := strings.Split(key, ":")
    return arr[3]
}

///////////////////
//               //
// table related //
//               //
///////////////////

// to store name of all tables
func getTablesKey(appId uint) string {
    return fmt.Sprintf("%s:%d:%s:tables", KeyPrefixDb, appId, KeyPrefixMeta)
}

// to store the id for next new record of a table
func getNextIdKey(appId uint, tableName string) string {
    return fmt.Sprintf("%s:%d:%s:nextId:%s", KeyPrefixDb, appId, KeyPrefixMeta, tableName)
}

// to store the name of fields for a table
func getTableKey(appId uint, tableName string) string {
    return fmt.Sprintf("%s:%d:%s:tn:%s", KeyPrefixDb, appId, KeyPrefixMeta, tableName)
}

// to store table fields which have index on
func getMetaTableIndexKey(appId uint, tableName string) string {
    return fmt.Sprintf("%s:%d:%s:idx:%s", KeyPrefixDb, appId, KeyPrefixMeta, tableName)
}

// to store the options for a table
func getTableOptionsKey(appId uint, tableName string) string {
    return fmt.Sprintf("%s:%d:%s:opt:%s", KeyPrefixDb, appId, KeyPrefixMeta, tableName)
}

// to store the insert filters for a table
func getTableInsertFilterKey(appId uint, tableName string) string {
    return fmt.Sprintf("%s:%d:%s:insfltr:%s", KeyPrefixDb, appId, KeyPrefixMeta, tableName)
}

func getColumnOptionsKey(appId uint, tableName string, fieldName string) string {
    return fmt.Sprintf("%s:%d:%s:fldopt:%s:%s", KeyPrefixDb, appId, KeyPrefixMeta, tableName, fieldName)
}


//////////////////
//              //
// data related //
//              //
//////////////////

// to store the id of a indexed field
func getIndexKey(appId uint, tableName string, fieldName string, value string) string {
    return fmt.Sprintf("%s:%d:%s:%s:%s:%s", KeyPrefixDb, appId, KeyPrefixIndex, tableName, fieldName, value)
}

// to store the value of a fields on a record of a table.
func getDataKeyBytes(appId uint, tableName string, fieldName string, id uint) []byte {
    base := fmt.Sprintf("%s:%d:%s:%s:%s:", KeyPrefixDb, appId, KeyPrefixData, tableName, fieldName)
    idBytes := utils.IntToByteArray(int64(id))
    return append([]byte(base), idBytes...)
}

// to get the start and end parameters of iterator which seeks certain value of a field
func getFieldDataIteratorStartAndEndKey(appId uint, tableName string, fieldName string) (string, string) {
    start := fmt.Sprintf("%s:%d:%s:%s:%s:", KeyPrefixDb, appId, KeyPrefixData, tableName, fieldName)
    end   := fmt.Sprintf("%s:%d:%s:%s:%s;", KeyPrefixDb, appId, KeyPrefixData, tableName, fieldName)
    return start, end
}

func getIdFromDataKey(key []byte) uint {
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

func getFriendKey(owner string, friendAddr string) string {
    return fmt.Sprintf("%s:%s:%s", KeyPrefixFriend, owner, friendAddr)
}

func getFriendIteratorStartAndEndKey(owner string) (string, string) {
    start := fmt.Sprintf("%s:%s:", KeyPrefixFriend, owner)
    end   := fmt.Sprintf("%s:%s;", KeyPrefixFriend, owner)
    return start, end
}

func getPendingFriendKey(owner string, friendAddr string) string {
    return fmt.Sprintf("%s:%s:%s", KeyPrefixPendingFriend, owner, friendAddr)
}

func getPendingFriendIteratorStartAndEndKey(owner string) (string, string) {
    start := fmt.Sprintf("%s:%s:", KeyPrefixPendingFriend, owner)
    end   := fmt.Sprintf("%s:%s;", KeyPrefixPendingFriend, owner)
    return start, end
}

///////////////////
//               //
// group related //
//               //
///////////////////

func getGroupsKey(appId uint) string {
    return fmt.Sprintf("%s:%d:%s", KeyPrefixDb, appId, KeyPrefixGroups)
}

func getGroupKey(appId uint, groupName string) string {
    return fmt.Sprintf("%s:%d:%s:%s", KeyPrefixDb, appId, KeyPrefixGroup, groupName)
}

func getAdminGroupKey(appId uint) string {
    return getGroupKey(appId, "admin")
}

//////////////////
//              //
// system level //
//              //
//////////////////

func getSysGroupKey(groupName string) string {
    return fmt.Sprintf("%s:%s", KeyPrefixSysGroup, groupName)
}

func getSysAdminGroupKey() string {
    return getSysGroupKey("admin")
}
