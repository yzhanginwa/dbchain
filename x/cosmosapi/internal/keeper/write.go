package keeper

import (
    "fmt"
    "strconv"
    "errors"
    sdk "github.com/cosmos/cosmos-sdk/types"
    "github.com/yzhanginwa/cosmos-api/x/cosmosapi/internal/other"
    "github.com/yzhanginwa/cosmos-api/x/cosmosapi/internal/types"
)


func (k Keeper) Insert(ctx sdk.Context, tableName string, fields types.RowFields, owner sdk.AccAddress) (uint, error){
    id, err := getNextId(k, ctx, tableName)
    if err != nil {
        return 0, errors.New(fmt.Sprintf("Failed to get id for table %s", tableName))
    }

    // to set the 2 special fields
    fields["id"] = strconv.Itoa(int(id))
    fields["created_by"] = owner.String()
    fields["created_at"] = other.GetCurrentBlockTime().String()

    k.Write(ctx, tableName, id, fields, owner)
    k.updateIndex(ctx, tableName, id, fields)
    return id, nil
}


//  // TODO: need to think over how and when to allow updating
//  func (k Keeper) Update(ctx sdk.Context, tableName string, id uint, fields types.RowFields) (uint, error){
//      // TODO: need to check the ownership of the record
//      k.Write(ctx, tableName, id, fields)
//      k.updateIndex(ctx, tableName, id, fields)
//      return id, nil
//  }


func (k Keeper) Write(ctx sdk.Context, tableName string, id uint, fields types.RowFields, owner sdk.AccAddress) (uint, error){
    store := ctx.KVStore(k.storeKey)

    fieldNames, err := k.getTableFields(ctx, tableName)
    if err != nil {
        return 0, errors.New(fmt.Sprintf("Failed to get fields for table %s", tableName))
    }

    if id == 0 {
        return 0, errors.New(fmt.Sprintf("Id for table %s is invalid", tableName))
    }

    for _, fieldName := range fieldNames {
        if value, ok := fields[fieldName]; ok {
            key := getDataKey(tableName, id, fieldName)
            store.Set([]byte(key), k.cdc.MustMarshalBinaryBare(value)) 
        }
    }

    return id, nil
}

func (k Keeper) Delete(ctx sdk.Context, tableName string, id uint) (uint, error){
    store := ctx.KVStore(k.storeKey)

    fieldNames, err := k.getTableFields(ctx, tableName)
    if err != nil {
        return 0, errors.New(fmt.Sprintf("Failed to get fields for table %s", tableName))
    }

    if id == 0 {
        return 0, errors.New("Id cannot be empty")
    }

    for _, fieldName := range fieldNames {
        key := getDataKey(tableName, id, fieldName)
    store.Delete([]byte(key)) 
    }

    return id, nil
}
