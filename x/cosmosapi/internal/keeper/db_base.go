package keeper

import (
    "fmt"
    "errors"
    sdk "github.com/cosmos/cosmos-sdk/types"
)


///////////
///////////

func getTableFields(k Keeper, ctx sdk.Context, tableName string) ([]string, error) {
    store := ctx.KVStore(k.storeKey)
    tableKey := getTableKey(tableName)
    bz := store.Get([]byte(tableKey))
    if bz == nil {
        return nil, errors.New(fmt.Sprintf("Table %s does not exist", tableName))
    }
    var fieldNames []string
    k.cdc.MustUnmarshalBinaryBare(bz, &fieldNames)
    return fieldNames, nil
}

