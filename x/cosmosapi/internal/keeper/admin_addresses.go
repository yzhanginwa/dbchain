package keeper

import (
    "errors"
    "bytes"
    sdk "github.com/cosmos/cosmos-sdk/types"
    "github.com/yzhanginwa/cosmos-api/x/cosmosapi/internal/types"
)

//////////////////
//              //
// system level //
//              //
//////////////////

func (k Keeper) CreateGenesisAdminGroup(ctx sdk.Context, genesisState types.GenesisState) {
    store := ctx.KVStore(k.storeKey)

    key := getSysAdminGroupKey()
    adminAddresses := genesisState.AdminAddresses
    store.Set([]byte(key), k.cdc.MustMarshalBinaryBare(adminAddresses))
}

////////////////////
//                //
// Database level //
//                //
////////////////////

func (k Keeper) AddAdminAccount(ctx sdk.Context, appId uint, adminAddress sdk.AccAddress, owner sdk.AccAddress) error {
    store := ctx.KVStore(k.storeKey)
    key := getAdminGroupKey(appId)

    var adminAddresses []sdk.AccAddress
    bz := store.Get([]byte(key))
    if bz != nil {
        k.cdc.MustUnmarshalBinaryBare(bz, &adminAddresses)
    }

    for _, addr := range adminAddresses {
        if bytes.Compare(adminAddress, addr) == 0 {
            return errors.New("Duplicate admin address found")
        }
    }

    adminAddresses = append(adminAddresses, adminAddress)
    store.Set([]byte(key), k.cdc.MustMarshalBinaryBare(adminAddresses))
    return nil
}

func (k Keeper) ShowAdminGroup(ctx sdk.Context, appId uint) []sdk.AccAddress {
    store := ctx.KVStore(k.storeKey)
    key := getAdminGroupKey(appId)

    bz := store.Get([]byte(key))
    if bz == nil {
        return []sdk.AccAddress{}
    }
    var adminAddresses []sdk.AccAddress
    k.cdc.MustUnmarshalBinaryBare(bz, &adminAddresses)
    return adminAddresses
}
