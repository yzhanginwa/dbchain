package keeper

import (
    "fmt"
    "sync"
    "errors"
    sdk "github.com/cosmos/cosmos-sdk/types"
)

var mutex = &sync.Mutex{}
var NextIds = make(map[string]uint)
var nextAppId uint

func (k Keeper) PeekNextId(ctx sdk.Context, appId uint, tableName string) (uint, error) {
    var nextIdKey = getNextIdKey(appId, tableName)
    var nextId uint
    var found bool
    if nextId, found = NextIds[nextIdKey]; found {
        return nextId, nil
    } else {
        store := DbChainStore(ctx, k.storeKey)
        mutex.Lock()
        defer mutex.Unlock()

        bz, _ := store.Get([]byte(nextIdKey))
        if bz != nil {
            k.cdc.MustUnmarshalBinaryBare(bz, &nextId)
            NextIds[nextIdKey] = nextId
            return nextId, nil
        } else {
            return 0, fmt.Errorf("Failed to peek next id for table %s\n", tableName)
        }
    }
}

func getNextId(k Keeper, ctx sdk.Context, appId uint, tableName string) (uint, error) {
    store := DbChainStore(ctx, k.storeKey)
    mutex.Lock()
    defer mutex.Unlock()

    var nextIdKey = getNextIdKey(appId, tableName)
    var nextId uint
    var found bool
    if nextId, found = NextIds[nextIdKey]; found {
    } else if bz, err := store.Get([]byte(nextIdKey)); bz != nil {
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
    NextIds[nextIdKey] = nextId + 1

    return nextId, nil
}

func dropNextId(k Keeper, ctx sdk.Context, appId uint, tableName string) {
    store := DbChainStore(ctx, k.storeKey)
    mutex.Lock()
    defer mutex.Unlock()

    var nextIdKey = getNextIdKey(appId, tableName)
    delete(NextIds, nextIdKey)
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
