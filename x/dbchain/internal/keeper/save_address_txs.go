package keeper

import (
	"encoding/json"
	sdk "github.com/dbchaincloud/cosmos-sdk/types"
	"time"
)

func (k Keeper) SaveAddrTxs ( ctx sdk.Context, addr sdk.AccAddress, value map[string]interface{}) error {
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

func (k Keeper) GetAddrTxs (ctx sdk.Context, addr sdk.AccAddress,  num ...uint ) []map[string]interface{} {
	store := DbChainStore(ctx, k.storeKey)

	start := getCurrentAddrTxsId(k, ctx, addr)
	end := uint(0)
	if len(num) != 0 && num[0] != 0{
		if start >= num[0] {
			end = start - num[0]
		}
	}

	result := make([]map[string]interface{}, 0)

	for i := start; i > end; i--{
		key := getAccountTxKey(addr.String(), i)
		bz, err := store.Get([]byte(key))
		if err != nil || bz == nil{
			continue
		}
		temp := make(map[string]interface{}, 0)
		err = json.Unmarshal(bz, &temp)
		if err != nil {
			continue
		}
		result = append(result, temp)
	}
	return result
}

func (k Keeper) GetAddrTxsByTime(ctx sdk.Context, addr sdk.AccAddress, startDate, endDate string) []map[string]interface{} {
	store := DbChainStore(ctx, k.storeKey)
	start := getCurrentAddrTxsId(k, ctx, addr)
	end := uint(0)
	result := make([]map[string]interface{}, 0)

	for i := start; i > end; i--{
		key := getAccountTxKey(addr.String(), i)
		bz, err := store.Get([]byte(key))
		if err != nil || bz == nil{
			continue
		}
		temp := make(map[string]interface{}, 0)
		err = json.Unmarshal(bz, &temp)
		if err != nil {
			continue
		}
		txTime :=temp["txTime"]
		if txTime == nil {
			continue
		}
		t := txTime.(string)
		status := checkTimeValid(endDate, startDate, t)
		if status == "high" {
			continue
		} else if status == "low" {
			break
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
	var nextId uint  = 0
	bz, _ := store.Get([]byte(key))
	if bz != nil {
		k.cdc.MustUnmarshalBinaryBare(bz, &nextId)
		return nextId - 1
	}

	return 0
}

func getNextSaveAddrTxsId(k Keeper, ctx sdk.Context, addr sdk.AccAddress) (uint, error) {
	store := DbChainStore(ctx, k.storeKey)
	mutex.Lock()
	defer mutex.Unlock()

	var nextIdKey = getNextAccountTxIdKey(addr.String())
	var nextId uint
	if bz, _ := store.Get([]byte(nextIdKey)); bz != nil {
		k.cdc.MustUnmarshalBinaryBare(bz, &nextId)
	} else {
		nextId = 1
	}

	err := store.Set([]byte(nextIdKey), k.cdc.MustMarshalBinaryBare(nextId + 1))
	if err != nil{
		return 0, err
	}

	return nextId, nil
}

// when txDate lower than startDate return low
// when txDate higher than endDate return  high
// when txDate lower than endDate and higher than startDate , return right
// endData is higher than startDate

func checkTimeValid(endTime , startTime , txDate string) string {
	end, _ := time.ParseInLocation("2006-01-02 15:04:05", endTime, time.Local)
	start, _ := time.ParseInLocation("2006-01-02 15:04:05", startTime, time.Local)
	txTime, _ := time.ParseInLocation("2006-01-02 15:04:05", txDate, time.Local)
	txSec := txTime.Unix()
	if txSec > end.Unix() {
		return "high"
	} else if txSec < start.Unix() {
		return "low"
	}
	return "right"
}