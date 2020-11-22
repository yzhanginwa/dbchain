package keeper

import (
    "fmt"
    "errors"
    sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) appendIndexForRow(ctx sdk.Context, appId uint, tableName string, id uint) (uint, error){
    store := ctx.KVStore(k.storeKey)
    indexFields, err := k.GetIndexFields(ctx, appId, tableName)
    if err != nil {
        return 0, errors.New(fmt.Sprintf("Failed to get index for table %s", tableName))
    }
    if id == 0 {
        return 0, errors.New(fmt.Sprintf("Id for table %s is invalid", tableName))
    }

    var mold []string
    for _, indexField := range indexFields {
        value, err := k.FindField(ctx, appId, tableName, id, indexField)
        if err != nil {
            return id, nil    // the value for this field is empty. we don't need to do anything. Because people would not search on an empty value.
        }
        key := getIndexKey(appId, tableName, indexField, value)
        bz := store.Get([]byte(key))
        if bz != nil {
            k.cdc.MustUnmarshalBinaryBare(bz, &mold)
        }
        mold = append(mold, fmt.Sprint(id))
        store.Set([]byte(key), k.cdc.MustMarshalBinaryBare(mold))
    }

    return id, nil
}

//////////////////////
//                  //
// helper functions //
//                  //
//////////////////////

