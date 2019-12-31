package keeper

import (
    "os"
    //"fmt"
    "errors"
    "github.com/cosmos/cosmos-sdk/codec"
    sdk "github.com/cosmos/cosmos-sdk/types"
    "github.com/cosmos/cosmos-sdk/x/bank"
    "github.com/tendermint/tendermint/libs/log"
    "github.com/yzhanginwa/cosmos-api/x/cosmosapi/internal/types"
)

var (
    logger = defaultLogger()
)

func defaultLogger() log.Logger {
    return log.NewTMLogger(log.NewSyncWriter(os.Stdout)).With("ethan1", "ethan2")
}

// Keeper maintains the link to storage and exposes getter/setter methods for the various parts of the state machine
type Keeper struct {
    CoinKeeper bank.Keeper

    storeKey sdk.StoreKey // Unexposed key to access store from sdk.Context

    cdc *codec.Codec // The wire codec for binary encoding/decoding.
}


// NewKeeper creates new instances of the cosmosapi Keeper
func NewKeeper(coinKeeper bank.Keeper, storeKey sdk.StoreKey, cdc *codec.Codec) Keeper {
    return Keeper{
        CoinKeeper: coinKeeper,
        storeKey:   storeKey,
        cdc:        cdc,
    }
}

// Check if the poll id is present in the store or not
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
}


// Gets a poll for an id
func (k Keeper) GetTable(ctx sdk.Context, name string) (types.Table, error) {
    store := ctx.KVStore(k.storeKey)
    bz := store.Get([]byte(getTableKey(name)))
    if bz == nil {
        return types.Table{}, errors.New("not found table")
    }
    var table types.Table
    k.cdc.MustUnmarshalBinaryBare(bz, &table)
    return table, nil
}

//////////////////////
//                  //
// helper functions //
//                  //
//////////////////////

