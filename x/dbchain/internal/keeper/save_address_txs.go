package keeper

import (
	"encoding/json"
	"fmt"
	sdk "github.com/dbchaincloud/cosmos-sdk/types"
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
	var endYear,endMonth,endDay, startYear, startMonth, startDay int
	fmt.Sscanf(startDate, "%d-%d-%d", &startYear, &startMonth, &startDay)
	fmt.Sscanf(endDate, "%d-%d-%d", &endYear, &endMonth, &endDay)


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
		status := checkTimeValid(endYear, endMonth, endDay, startYear, startMonth, startDay, t)

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

// when txDate lower than startDate return low
// when txDate higher than endDate return  high
// when txDate lower than endDate and higher than startDate , return right
// endData is higher than startDate

func checkTimeValid(endYear, endMonth, endDay, startYear, startMonth, startDay int, txDate string)  string {
	txYear, txMonth, txDay := 0,0,0
	fmt.Sscanf(txDate, "%d-%d-%d", &txYear, &txMonth, & txDay)
	if txYear > endYear {
		return "high"
	} else if txYear < startYear {
		return "low"
	}

	if txMonth > endMonth {
		return "high"
	} else if txMonth < startMonth {
		return "low"
	}

	if txDay > endDay {
		return "high"
	} else if txDay < startDay {
		return "low"
	}
	return "right"

}