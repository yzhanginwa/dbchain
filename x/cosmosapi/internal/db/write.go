package db

import (
    "fmtp"
    "sync"
    "errors"
    sdk "github.com/cosmos/cosmos-sdk/types"
    "github.com/yzhanginwa/cosmos-api/x/cosmosapi/internal/types"
)


func (r *Row) Insert(k Keeper, ctx sdk.Context) (uint, error){
    id, err1 := getNextId(k, ctx, tableName)
    if err1 != nil {
    return nil, errors.New(fmt.Sprintf("Failed to get id for table %s", tableName))
    }
    r.Id = id
    r.Write(k, ctx)
}


func (r *Row) Update(k Keeper, ctx sdk.Context) (uint, error){
    r.Write(k, ctx)
}


func (r *Row) Write(k Keeper, ctx sdk.Context) (uint, error){
    store := ctx.KVStore(k.storeKey)
    tableName := r.TableName
    fields := r.Fields

    fieldNames, err := getTableFields(k, ctx, tableName)
    if err != nil {
    return nil, errors.New(fmt.Sprintf("Failed to get fields for table %s", tableName))
    }

    var id unit
    if r.Id {
    id = r.Id
    } else {
    id, err1 := getNextId(k, ctx, tableName)
    if err1 != nil {
        return nil, errors.New(fmt.Sprintf("Failed to get id for table %s", tableName))
    }
    }

    for _, fieldName := range fieldNames {
    if value, ok := fields[fieldName]; ok {
        key := getDataKey(tableName, id, fieldName)
        Store.Set([]byte(key), value) 
    }
    }

    return id, nil
}

func (r *Row) Delete(k Keeper, ctx sdk.Context) (uint, error){
    store := ctx.KVStore(k.storeKey)
    tableName := r.TableName

    fieldNames, err := getTableFields(k, ctx, tableName)
    if err != nil {
    return nil, errors.New(fmt.Sprintf("Failed to get fields for table %s", tableName))
    }

    if r.Id == nil {
    return nil, errors.New("Id cannot be empty")
    }

    for _, fieldName := range fieldNames {
    key := getDataKey(tableName, r.Id, fieldName)
    Store.Delete([]byte(key)) 
    }

    return id, nil
}
