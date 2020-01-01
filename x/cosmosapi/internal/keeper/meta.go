package keeper

import (
    "fmt"
    "errors"
    sdk "github.com/cosmos/cosmos-sdk/types"
    "github.com/yzhanginwa/cosmos-api/x/cosmosapi/internal/types"
)

func (k Keeper) getTables(ctx sdk.Context) ([]string, error) {
    store := ctx.KVStore(k.storeKey)
    tablesKey := getTablesKey()
    bz := store.Get([]byte(tablesKey))
    if bz == nil {
        return nil, errors.New("No tables found")
    }
    var tableNames []string
    k.cdc.MustUnmarshalBinaryBare(bz, &tableNames)
    return tableNames, nil
}


// Check if the table is present in the store or not
func (k Keeper) IsTablePresent(ctx sdk.Context, name string) bool {
    store := ctx.KVStore(k.storeKey)
    return store.Has([]byte(getTableKey(name)))
}


// Create a new table
func (k Keeper) CreateTable(ctx sdk.Context, owner sdk.AccAddress, name string, fields []string) {
    store := ctx.KVStore(k.storeKey)
    table := types.NewTable()
    table.Owner = owner
    table.Name = name
    table.Fields = fields 
    store.Set([]byte(getTableKey(table.Name)), k.cdc.MustMarshalBinaryBare(table))

    var tables []string
    bz :=store.Get([]byte(getTablesKey()))
    if bz == nil {
        tables = append(tables, table.Name)
    } else {
        k.cdc.MustUnmarshalBinaryBare(bz, &tables)
        tables = append(tables, table.Name)
    }
    store.Set([]byte(getTablesKey()), k.cdc.MustMarshalBinaryBare(tables))
}


// Get a table 
func (k Keeper) GetTable(ctx sdk.Context, name string) (types.Table, error) {
    store := ctx.KVStore(k.storeKey)
    bz := store.Get([]byte(getTableKey(name)))
    if bz == nil {
        return types.Table{}, errors.New(fmt.Sprintf("table %s not found", name))
    }
    var table types.Table
    k.cdc.MustUnmarshalBinaryBare(bz, &table)
    return table, nil
}

////////////////////
//                //
// helper methods //
//                //
////////////////////

