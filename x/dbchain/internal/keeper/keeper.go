package keeper

import (
    "fmt"
    "github.com/cosmos/cosmos-sdk/codec"
    sdk "github.com/cosmos/cosmos-sdk/types"
    "github.com/cosmos/cosmos-sdk/x/bank"
    "github.com/tendermint/tendermint/libs/log"
    "github.com/yzhanginwa/dbchain/x/dbchain/internal/types"
)

// Keeper maintains the link to storage and exposes getter/setter methods for the various parts of the state machine
type Keeper struct {
    CoinKeeper bank.Keeper

    storeKey sdk.StoreKey // Unexposed key to access store from sdk.Context

    cdc *codec.Codec // The wire codec for binary encoding/decoding.
}


// NewKeeper creates new instances of the dbchain Keeper
func NewKeeper(coinKeeper bank.Keeper, storeKey sdk.StoreKey, cdc *codec.Codec) Keeper {
    return Keeper{
        CoinKeeper: coinKeeper,
        storeKey:   storeKey,
        cdc:        cdc,
    }
}

// Logger returns a module-specific logger.
func (k Keeper) Logger(ctx sdk.Context) log.Logger {
    return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

//////////////////////
//                  //
// helper functions //
//                  //
//////////////////////

