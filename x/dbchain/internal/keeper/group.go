package keeper

import (
    "bytes"
    "errors"
    "fmt"
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
    store := DbChainStore(ctx, k.storeKey)

    key := getSysAdminGroupKey()
    adminAddresses := genesisState.AdminAddresses
    store.Set([]byte(key), k.cdc.MustMarshalBinaryBare(adminAddresses))
}

func (k Keeper) GetSysAdmins(ctx sdk.Context) []sdk.AccAddress {
    store := DbChainStore(ctx, k.storeKey)

    var sysAdmins []sdk.AccAddress
    key := getSysAdminGroupKey()
    bz, err := store.Get([]byte(key))
    if err != nil{
        return sysAdmins
    }
    if bz != nil {
        k.cdc.MustUnmarshalBinaryBare(bz, &sysAdmins)
    }
    return sysAdmins
}

func (k Keeper) IsSysAdmin(ctx sdk.Context, addr sdk.AccAddress) bool {
    sysAdmins := k.GetSysAdmins(ctx)
    for _, sysAdmin := range sysAdmins {
        if addr.Equals(sysAdmin) {
            return true
        }
    }
    return false
}

func (k Keeper) IsGroupMember(ctx sdk.Context, appId uint, groupName string, addr sdk.AccAddress) bool {
    groupMember := k.getGroupMembers(ctx, appId, groupName)
    for _, member := range groupMember {
        if member.String() == addr.String() {
            return true
        }
    }
    return false
}
////////////////////
//                //
// Database level //
//                //
////////////////////

func (k Keeper) ModifyGroup(ctx sdk.Context, appId uint, action string, groupName string) error {
    store := DbChainStore(ctx, k.storeKey)
    key := getGroupsKey(appId)
    var groups []string
    bz, err := store.Get([]byte(key))
    if err != nil{
        return err
    }
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
        if groupName == "admin" || groupName == "auditor" {
            return errors.New("You are not supposed to delete system group")
        }
        if position > -1 {
            if len(k.getGroupMembers(ctx, appId, groupName)) > 0 {
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

func (k Keeper) SetGroupMemo(ctx sdk.Context, appId uint, groupName string, memo string) {
    store := DbChainStore(ctx, k.storeKey)
    key := getGroupMemoKey(appId, groupName)
    store.Set([]byte(key), k.cdc.MustMarshalBinaryBare(memo))
}

func (k Keeper) ModifyGroupMember(ctx sdk.Context, appId uint, group string, action string, member sdk.AccAddress) error {
    groups := k.getGroups(ctx, appId)
    if !utils.ItemExists(groups, group) {
        return errors.New(fmt.Sprintf("Group %s does not exist", group))
    }

    store := DbChainStore(ctx, k.storeKey)
    key := getGroupKey(appId, group)

    var members []sdk.AccAddress
    bz, err := store.Get([]byte(key))
    if err != nil{
        return err
    }
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

func (k Keeper) getGroupMembers(ctx sdk.Context, appId uint, groupName string) []sdk.AccAddress {
    store := DbChainStore(ctx, k.storeKey)
    key := getGroupKey(appId, groupName)

    bz, err := store.Get([]byte(key))
    if bz == nil || err != nil{
        return []sdk.AccAddress{}
    }
    var addresses []sdk.AccAddress
    k.cdc.MustUnmarshalBinaryBare(bz, &addresses)
    return addresses
}

func (k Keeper) getGroupMembersMemo(ctx sdk.Context, appId uint, groupName string) string {
    store := DbChainStore(ctx, k.storeKey)
    key := getGroupMemoKey(appId, groupName)
    bz, err := store.Get([]byte(key))

    if bz == nil || err != nil{
        return ""
    }

    var memo string
    k.cdc.MustUnmarshalBinaryBare(bz, &memo)
    return memo
}

func (k Keeper) getGroups(ctx sdk.Context, appId uint) []string {
    store := DbChainStore(ctx, k.storeKey)
    key := getGroupsKey(appId)

    bz, err := store.Get([]byte(key))
    if bz == nil || err != nil{
        return []string{}
    }
    var groups []string
    k.cdc.MustUnmarshalBinaryBare(bz, &groups)
    return groups
}

func (k Keeper) getGroupsDetail(ctx sdk.Context, appId uint) []map[string]interface{} {
    var groupsDetail []map[string]interface{}
    groups := k.getGroups(ctx, appId)
    for _,groupName := range groups {
        groupDetail := make(map[string]interface{})
        members := k.getGroupMembers(ctx, appId, groupName)
        memo := k.getGroupMembersMemo(ctx, appId, groupName)
        groupDetail["group_name"] = groupName
        groupDetail["group_members"] = members
        groupDetail["group_memo"] = memo
        groupsDetail = append(groupsDetail, groupDetail)
    }
    return  groupsDetail
}