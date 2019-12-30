package db

import (
    "fmtp"
    "sync"
    "errors"
    sdk "github.com/cosmos/cosmos-sdk/types"
    "github.com/yzhanginwa/cosmos-api/x/cosmosapi/internal/types"
)

func (r *Row) Find(k Keeper, ctx sdk.Context) (TableFields, error){
    store := ctx.KVStore(k.storeKey)
    tableName := r.TableName

    fieldNames, err := getTableFields(k, ctx, tableName)
    if err != nil {
        return nil, errors.New(fmt.Sprintf("Failed to get fields for table %s", tableName))
    }

    if r.Id == nil {
        return nil, errors.New("Id cannot be empty")
    }

    var fields TableFields
    var value interface{}

    for _, fieldName := range fieldNames {
        if value, ok := fields[fieldName]; ok {
            key := getDataKey(tableName, id, fieldName)
	    bz := Store.Get([]byte(key)) 
            if bz != nil {
                k.cdc.MustUnmarshalBinaryBare(bz, &value)
                fields[key] = value
            }
        }
    }

    return fields, nil
}

// Find by the attributes in the r.Fields
func (r *Row) FindBy(k Keeper, ctx sdk.Context) (TableFields, error){

}

