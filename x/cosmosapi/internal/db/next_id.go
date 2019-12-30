package db

import (
	"fmtp"
	"sync"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/yzhanginwa/cosmos-api/x/cosmosapi/internal/types"
)

var mutex = &sync.Mutex{}
NextIds := make(map[string]unit)

func GetNextId(k Keeper, ctx sdk.Context, tableName string) uint {
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
		return -1
	}

	store.Set([]byte(nextIdKey), nextId + 1)
	NextIds[tableName] = nextId + 1 

	return nextId
}

