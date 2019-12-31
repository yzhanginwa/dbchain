package db

import (
    "errors"
)

type TableFields map[string]interface{}

type Row struct {
    TableName: string   `json:"table_name"`
    Id: uint            `json:"id"`
    Fields: TableFields `json:"columns"`
}

func NewRow(tableName string, id uint, fields TableFields) {
    return Row {
    TableName: tableName,
    Id: id,
    Fields: fields,
    }
}

var mutex = &sync.Mutex{}
NextIds := make(map[string]unit)

func getNextId(k Keeper, ctx sdk.Context, tableName string) (uint, errors) {
    store := ctx.KVStore(k.storeKey)
    mutex.Lock()
    defer mutex.Unlock()

    var nextIdKey = getNextIdKey(tableName)
    var nextId uint

    if nextId = NextIds[tableName] {
    } else if bz := store.Get([]byte(nextIdKey)) {
    k.cdc.MustUnmarshalBinaryBare(bz, &nextId)
    } else if store.get([]byte(getTableKey(tableName))) {
    nextId = 1
    } else {
    return -1, errors.New(fmt.Sprintf("Invalid table name %s", tableName))
    }

    store.Set([]byte(nextIdKey), nextId + 1)
    NextIds[tableName] = nextId + 1

    return nextId, nil
}

func getTableFields(k Keeper, ctx sdk.Context, tableName string) []string {
    store := ctx.KVStore(k.storeKey)
    tableKey := getTableKey(tableName)
    bz := store.Get([]byte(tableKey)) {
    if bz == nil {
    return nil, errors.New(fmt.SprintF("Table %s does not exist", tableName))
    }
    var fieldNames []string
    k.cdc.MustUnmarshalBinaryBare(bz, &fieldNames)
    return fieldNames
}
