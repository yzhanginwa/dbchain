package keeper

import (
    "fmt"
    "errors"
    sdk "github.com/cosmos/cosmos-sdk/types"
    "github.com/yzhanginwa/cosmos-api/x/cosmosapi/internal/types"
)

func (k Keeper) Find(ctx sdk.Context, tableName string, id uint) (types.RowFields, error){
    store := ctx.KVStore(k.storeKey)

    fieldNames, err := k.getTableFields(ctx, tableName)
    if err != nil {
        return nil, errors.New(fmt.Sprintf("Failed to get fields for table %s", tableName))
    }

    if id == 0 {
        return nil, errors.New("Id cannot be 0")
    }

    var fields = make(types.RowFields)
    var value string

    for _, fieldName := range fieldNames {
        key := getDataKey(tableName, id, fieldName)
        bz := store.Get([]byte(key)) 
        if bz != nil {
            k.cdc.MustUnmarshalBinaryBare(bz, &value)
            fields[fieldName] = value
        }
    }

    return fields, nil
}

// Find by the attributes in the r.Fields
//func (r *Row) FindBy(k Keeper, ctx sdk.Context) (types.RowFields, error){

//}

