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
        return 0, nil
    }
}

func getNextId(k Keeper, ctx sdk.Context, appId uint, tableName string) (uint, error) {
    store := ctx.KVStore(k.storeKey)
    mutex.Lock()
    defer mutex.Unlock()

    var nextIdKey = getNextIdKey(appId, tableName)
    var nextId uint
    var found bool
    if nextId, found = NextIds[nextIdKey]; found {
    } else if bz := store.Get([]byte(nextIdKey)); bz != nil {
        k.cdc.MustUnmarshalBinaryBare(bz, &nextId)
    } else if bz = store.Get([]byte(getTableKey(appId, tableName))); bz != nil {
        nextId = 1
    } else {
        return 0, errors.New(fmt.Sprintf("Invalid table name %s", tableName))
    }

    store.Set([]byte(nextIdKey), k.cdc.MustMarshalBinaryBare(nextId + 1))
    NextIds[nextIdKey] = nextId + 1

    return nextId, nil
}

func registerDatabaseId(k Keeper, ctx sdk.Context, appCode string) (uint, error) {
    store := ctx.KVStore(k.storeKey)
    mutex.Lock()
    defer mutex.Unlock()

    var nextAppIdKey = getDatabaseNextIdKey()
    if nextAppId < 1 {
        bz := store.Get([]byte(nextAppIdKey))
        if bz == nil {
            nextAppId = 1
        } else {
            k.cdc.MustUnmarshalBinaryBare(bz, &nextAppId)
        }
    }

    store.Set([]byte(nextAppIdKey), k.cdc.MustMarshalBinaryBare(nextAppId + 1))
    currentAppId := nextAppId
    nextAppId += 1
    return currentAppId, nil
}
