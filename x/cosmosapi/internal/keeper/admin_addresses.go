package keeper

import (
    "errors"
    "bytes"
    sdk "github.com/cosmos/cosmos-sdk/types"
    "github.com/yzhanginwa/cosmos-api/x/cosmosapi/internal/types"
)


func (k Keeper) CreateGenesisAdminGroup(ctx sdk.Context, genesisState types.GenesisState) {
    store := ctx.KVStore(k.storeKey)

    key := getAdminGroupKey()
    adminAddresses := genesisState.AdminAddresses
    store.Set([]byte(key), k.cdc.MustMarshalBinaryBare(adminAddresses))
}

func (k Keeper) AddAdminAccount(ctx sdk.Context, adminAddress sdk.AccAddress, owner sdk.AccAddress) (bool, error){
    store := ctx.KVStore(k.storeKey)
    key := getAdminGroupKey()

    bz := store.Get([]byte(key))
    if bz == nil {
        return false, errors.New("No admin group found")
    }
    var adminAddresses []sdk.AccAddress
    k.cdc.MustUnmarshalBinaryBare(bz, &adminAddresses)

    var owner_was_admin = false
    for _, addr := range adminAddresses {
        if bytes.Compare(adminAddress, addr) == 0 {
            return false, errors.New("Duplicate admin address found")
        }
        if bytes.Compare(owner, addr) == 0 {
            owner_was_admin = true
        }
    }

    if !owner_was_admin {
        return false, errors.New("Not authorized signer")
    }

    adminAddresses = append(adminAddresses, adminAddress)
    store.Set([]byte(key), k.cdc.MustMarshalBinaryBare(adminAddresses))
    return true, nil
}
