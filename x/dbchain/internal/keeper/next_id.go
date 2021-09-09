package keeper

import (
    "errors"
    "fmt"
    sdk "github.com/cosmos/cosmos-sdk/types"
    "sync"
)

var mutex = &sync.Mutex{}
var nextAppId uint

func (k Keeper) PeekNextId(ctx sdk.Context, appId uint, tableName string) (uint, error) {
    var nextIdKey = getNextIdKey(appId, tableName)
    var nextId uint
    store := DbChainStore(ctx, k.storeKey)
    if bz, _ := store.Get([]byte(nextIdKey)); bz != nil {
        k.cdc.MustUnmarshalBinaryBare(bz, &nextId)
        return nextId, nil
    } else {
        return 0, nil
    }
}

func getNextId(k Keeper, ctx sdk.Context, appId uint, tableName string) (uint, error) {
    store := DbChainStore(ctx, k.storeKey)
    mutex.Lock()
    defer mutex.Unlock()

    var nextIdKey = getNextIdKey(appId, tableName)
    var nextId uint
    if bz, err := store.Get([]byte(nextIdKey)); bz != nil {
        k.cdc.MustUnmarshalBinaryBare(bz, &nextId)
    } else if bz, err = store.Get([]byte(getTableKey(appId, tableName))); bz != nil {
        nextId = 1
    } else if err != nil{
        return 0, err
    }else {
        return 0, errors.New(fmt.Sprintf("Invalid table name %s", tableName))
    }

    err := store.Set([]byte(nextIdKey), k.cdc.MustMarshalBinaryBare(nextId + 1))
    if err != nil{
        return 0, err
    }

    return nextId, nil
}

func dropNextId(k Keeper, ctx sdk.Context, appId uint, tableName string) {
    store := DbChainStore(ctx, k.storeKey)
    mutex.Lock()
    defer mutex.Unlock()

    var nextIdKey = getNextIdKey(appId, tableName)
    store.Delete([]byte(nextIdKey))
}

func registerDatabaseId(k Keeper, ctx sdk.Context, appCode string) (uint, error) {
    store := DbChainStore(ctx, k.storeKey)
    mutex.Lock()
    defer mutex.Unlock()

    var nextAppIdKey = getDatabaseNextIdKey()
    if nextAppId < 1 {
        bz, err := store.Get([]byte(nextAppIdKey))
        if err != nil{
            return 0, err
        }
        if bz == nil {
            nextAppId = 1
        } else {
            k.cdc.MustUnmarshalBinaryBare(bz, &nextAppId)
        }
    }

    err := store.Set([]byte(nextAppIdKey), k.cdc.MustMarshalBinaryBare(nextAppId + 1))
    if err != nil{
        return 0, err
    }
    currentAppId := nextAppId
    nextAppId += 1
    return currentAppId, nil
}
