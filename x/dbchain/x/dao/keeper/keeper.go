package keeper

import (
    "github.com/dbchaincloud/cosmos-sdk/codec"
    sdk "github.com/dbchaincloud/cosmos-sdk/types"
    "github.com/yzhanginwa/dbchain/x/dao/internal/types"
)

type Keeper struct {
    storeKey sdk.StoreKey
    cdc      *codec.Codec
}

func NewKeeper(storeKey sdk.StoreKey, cdc *codec.Codec) Keeper {
    return Keeper{
        storeKey: storeKey,
        cdc:      cdc,
    }
}
