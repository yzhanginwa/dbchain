package keeper

import (
    sdk "github.com/cosmos/cosmos-sdk/types"
    "github.com/yzhanginwa/cosmos-api/x/cosmosapi/internal/types"
)


func (k Keeper) CreateGenesisAdminGroup(ctx sdk.Context, genesisState types.GenesisState) {
    store := ctx.KVStore(k.storeKey)

    key := getAdminGroupKey()
    adminAddresses := genesisState.AdminAddresses
    store.Set([]byte(key), k.cdc.MustMarshalBinaryBare(adminAddresses))
}

