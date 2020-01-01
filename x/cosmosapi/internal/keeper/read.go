package keeper

import (
    "fmt"
    "errors"
    sdk "github.com/cosmos/cosmos-sdk/types"
    "github.com/yzhanginwa/cosmos-api/x/cosmosapi/internal/types"
)

func (r *Row) Find(k Keeper, ctx sdk.Context) (types.RowFields, error){
    store := ctx.KVStore(k.storeKey)
    tableName := r.TableName

    fieldNames, err := k.getTableFields(ctx, tableName)
    if err != nil {
    return nil, errors.New(fmt.Sprintf("Failed to get fields for table %s", tableName))
    }

    if r.Id == 0 {
    return nil, errors.New("Id cannot be empty")
    }

    var fields types.RowFields

    for _, fieldName := range fieldNames {
    if value, ok := fields[fieldName]; ok {
        key := getDataKey(tableName, r.Id, fieldName)
        bz := store.Get([]byte(key)) 
        if bz != nil {
        k.cdc.MustUnmarshalBinaryBare(bz, &value)
        fields[key] = value
        }
    }
    }

    return fields, nil
}

// Find by the attributes in the r.Fields
//func (r *Row) FindBy(k Keeper, ctx sdk.Context) (types.RowFields, error){

//}

