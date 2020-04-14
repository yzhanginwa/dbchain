package keeper

import (
    "errors"
    "bytes"
    sdk "github.com/cosmos/cosmos-sdk/types"
    "github.com/yzhanginwa/dbchain/x/dbchain/internal/types"
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

func (k Keeper) GetSysAdmins(ctx sdk.Context) []string {
    store := ctx.KVStore(k.storeKey)

    var sysAdmins []string
    key := getSysAdminGroupKey()
    bz := store.Get([]byte(key))
    if bz != nil {
        k.cdc.MustUnmarshalBinaryBare(bz, &sysAdmins)
    }
    return sysAdmins
}

////////////////////
//                //
// Database level //
//                //
////////////////////

func (k Keeper) CreateGroup(ctx sdk.Context, appId uint, groupName string) error {
    store := ctx.KVStore(k.storeKey)
    key := getGroupsKey(appId)
    var groups []string
    bz := store.Get([]byte(key))
    if bz != nil {
        k.cdc.MustUnmarshalBinaryBare(bz, &groups)
        for _, grp := range groups {
            if groupName  == grp {
                return errors.New("Duplicate group name")
            }
        }
    }

    groups = append(groups, groupName) 
    store.Set([]byte(key), k.cdc.MustMarshalBinaryBare(groups))
    return nil
}

func (k Keeper) AddAdminAccount(ctx sdk.Context, appId uint, adminAddress sdk.AccAddress) error {
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
