package keeper

import (
    "fmt"
    "errors"
    "bytes"
    sdk "github.com/cosmos/cosmos-sdk/types"
    "github.com/yzhanginwa/dbchain/x/dbchain/internal/types"
    "github.com/yzhanginwa/dbchain/x/dbchain/internal/utils"
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

func (k Keeper) AddGroupMember(ctx sdk.Context, appId uint, group string, member sdk.AccAddress) error {
    groups := k.ShowGroups(ctx, appId)
    if !utils.ItemExists(groups, group) {
        return errors.New(fmt.Sprintf("Group %s does not exist", group))
    }

    store := ctx.KVStore(k.storeKey)
    key := getGroupKey(appId, group)

    var members []sdk.AccAddress
    bz := store.Get([]byte(key))
    if bz != nil {
        k.cdc.MustUnmarshalBinaryBare(bz, &members)
    }

    for _, addr := range members {
        if bytes.Compare(member, addr) == 0 {
            return errors.New("Duplicate admin address found")
        }
    }

    members  = append(members, member)
    store.Set([]byte(key), k.cdc.MustMarshalBinaryBare(members))
    return nil
}

func (k Keeper) ShowGroup(ctx sdk.Context, appId uint, groupName string) []sdk.AccAddress {
    store := ctx.KVStore(k.storeKey)
    key := getGroupKey(appId, groupName)

    bz := store.Get([]byte(key))
    if bz == nil {
        return []sdk.AccAddress{}
    }
    var addresses []sdk.AccAddress
    k.cdc.MustUnmarshalBinaryBare(bz, &addresses)
    return addresses
}

func (k Keeper) ShowGroups(ctx sdk.Context, appId uint) []string {
    store := ctx.KVStore(k.storeKey)
    key := getGroupsKey(appId)

    bz := store.Get([]byte(key))
    if bz == nil {
        return []string{}
    }
    var groups []string
    k.cdc.MustUnmarshalBinaryBare(bz, &groups)
    return groups
}
