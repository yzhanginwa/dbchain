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

func (k Keeper) GetSysAdmins(ctx sdk.Context) []sdk.AccAddress {
    store := ctx.KVStore(k.storeKey)

    var sysAdmins []sdk.AccAddress
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

func (k Keeper) ModifyGroup(ctx sdk.Context, appId uint, action string, groupName string) error {
    store := ctx.KVStore(k.storeKey)
    key := getGroupsKey(appId)
    var groups []string
    bz := store.Get([]byte(key))
    if bz != nil {
        k.cdc.MustUnmarshalBinaryBare(bz, &groups)
    }

    var position = -1
    for i, grp := range groups {
        if groupName  == grp {
            position = i
            break
        }
    }

    if action == "add" {
        if position > -1 {
            return errors.New("Duplicate group name")
        } else {
            groups = append(groups, groupName)
        }
    } else {
        if groupName == "admin" {
            return errors.New("You should not drop Admin group")
        }
        if position > -1 {
            if len(k.ShowGroup(ctx, appId, groupName)) > 0 {
                return errors.New("Group is not empty")
            }

            l := len(groups)
            groups[position] = groups[l-1]
            groups[l-1] = ""
            groups = groups[:l-1]
        } else {
            return errors.New("Group not exist")
        }
    }

    if len(groups) < 1 {
        store.Delete([]byte(key))
    } else {
        store.Set([]byte(key), k.cdc.MustMarshalBinaryBare(groups))
    }
    return nil
}

func (k Keeper) ModifyGroupMember(ctx sdk.Context, appId uint, group string, action string, member sdk.AccAddress) error {
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

    var position = -1
    for i, addr := range members {
        if bytes.Compare(member, addr) == 0 {
            position = i
            break
        }
    }

    if action == "add" {
        if position > -1 {
            return errors.New("Duplicate member")
        } else {
            members = append(members, member)
        }
    } else {
        if position > -1 {
            l := len(members)
            members[position] = members[l-1]
            members[l-1] = sdk.AccAddress{}
            members = members[:l-1]
        } else {
            return errors.New("Member not exist")
        }
    }

    if len(members) < 1 {
        store.Delete([]byte(key))
    } else {
        store.Set([]byte(key), k.cdc.MustMarshalBinaryBare(members))
    }
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
