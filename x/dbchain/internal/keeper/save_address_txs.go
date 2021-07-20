package keeper

import (
	"encoding/json"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) SaveAddrTxs ( ctx sdk.Context, addr sdk.AccAddress, value map[string]string) error {
	store := DbChainStore(ctx, k.storeKey)
	id, err := getNextSaveAddrTxsId(k, ctx, addr)
	if err != nil {
		return err
	}
	key := getAccountTxKey(addr.String(), id)
	bz , err := json.Marshal(value)
	if err != nil {
		return err
	}
	err = store.Set([]byte(key), bz)
	return err
}

func (k Keeper) GetAddrTxs (ctx sdk.Context, addr sdk.AccAddress,  num ...uint ) []map[string]string {
	store := DbChainStore(ctx, k.storeKey)

	start := getCurrentAddrTxsId(k, ctx, addr)
	end := uint(0)
	if len(num) != 0 {
		if start >= num[0] {
			end = start - end
		}
	}

	result := make([]map[string]string, 0)

	for i := start; i > end; i--{
		key := getAccountTxKey(addr.String(), i)
		bz, err := store.Get([]byte(key))
		if err != nil || bz == nil{
			continue
		}
		temp := make(map[string]string, 0)
		err = json.Unmarshal(bz, &temp)
		if err != nil {
			continue
		}
		result = append(result, temp)
	}
	return result
}

////////////////////////////
//                        //
//       help func        //
//                        //
////////////////////////////

func getCurrentAddrTxsId(k Keeper, ctx sdk.Context, addr sdk.AccAddress) uint {
	store := DbChainStore(ctx, k.storeKey)
	mutex.Lock()
	defer mutex.Unlock()
	var key = getNextAccountTxIdKey(addr.String())
	nextId ,ok := NextIds[key]
	if  ok {
		return nextId -1
	} else {
		bz, _ := store.Get([]byte(key))
		if bz != nil {
			k.cdc.MustUnmarshalBinaryBare(bz, &nextId)
			return nextId - 1
		}
	}
	return 0
}

func getNextSaveAddrTxsId(k Keeper, ctx sdk.Context, addr sdk.AccAddress) (uint, error) {
	store := DbChainStore(ctx, k.storeKey)
	mutex.Lock()
	defer mutex.Unlock()

	var nextIdKey = getNextAccountTxIdKey(addr.String())
	var nextId uint
	var found bool
	if nextId, found = NextIds[nextIdKey]; found {

	} else if bz, _ := store.Get([]byte(nextIdKey)); bz != nil {
		k.cdc.MustUnmarshalBinaryBare(bz, &nextId)
	} else {
		nextId = 1
	}

	err := store.Set([]byte(nextIdKey), k.cdc.MustMarshalBinaryBare(nextId + 1))
	if err != nil{
		return 0, err
	}
	NextIds[nextIdKey] = nextId + 1

	return nextId, nil
}