package keeper

import (
    "errors"
    sdk "github.com/cosmos/cosmos-sdk/types"
    "github.com/yzhanginwa/dbchain/x/dbchain/internal/types"
)

func (k Keeper) AddFriend(ctx sdk.Context, owner sdk.AccAddress, ownerName string, friendAddr string, friendName string) error {
    ownerAddr:= owner.String()
    err := k.addFriend(ctx, ownerAddr, friendAddr, friendName)
    if err != nil {
        return err
    }

    k.addPendingFriend(ctx, ownerAddr, ownerName, friendAddr)
    return nil
}

func (k Keeper) DropFriend(ctx sdk.Context, owner sdk.AccAddress, friendAddr string) error {
    ownerAddr:= owner.String()
    return k.deleteFriend(ctx, ownerAddr, friendAddr)
}

func (k Keeper) RespondFriend(ctx sdk.Context, owner sdk.AccAddress, friendAddr string, action string) error {
    ownerAddr:= owner.String()
    if action == "delete" {
        k.deleteFriend(ctx, ownerAddr, friendAddr)
    } else if action == "reject" {
        k.deletePendingFriend(ctx, ownerAddr, friendAddr)
    } else if action == "accept" {
        pf, err := k.getPendingFriend(ctx, ownerAddr, friendAddr)
        if err != nil {
            return err
        }
        k.addFriend(ctx, ownerAddr, pf.Address, pf.Name)
        k.deletePendingFriend(ctx, ownerAddr, friendAddr)
    } else {
        return errors.New("Wrong action for responding friend")
    }

    return nil
}

func (k Keeper) GetFriends(ctx sdk.Context, owner sdk.AccAddress) []types.Friend {
    store := ctx.KVStore(k.storeKey)
    bech32 := owner.String()

    start, end := getFriendIteratorStartAndEndKey(bech32)
    iter := store.Iterator([]byte(start), []byte(end))
    var mold types.Friend
    var friends []types.Friend
    for ; iter.Valid(); iter.Next() {
        val := iter.Value()
        k.cdc.MustUnmarshalBinaryBare(val, &mold)
        friends = append(friends, mold)
    }
    return friends
}

func (k Keeper) GetPendingFriends(ctx sdk.Context, owner sdk.AccAddress) []types.Friend {
    store := ctx.KVStore(k.storeKey)
    bech32 := owner.String()

    start, end := getPendingFriendIteratorStartAndEndKey(bech32)
    iter := store.Iterator([]byte(start), []byte(end))
    var mold types.Friend
    var friends []types.Friend
    for ; iter.Valid(); iter.Next() {
        val := iter.Value()
        k.cdc.MustUnmarshalBinaryBare(val, &mold)
        friends = append(friends, mold)
    }
    return friends
}

//////////////////////
//                  //
// helper functions //
//                  //
//////////////////////

func (k Keeper) addFriend(ctx sdk.Context, ownerAddr string, friendAddr string, friendName string) error {
    store := ctx.KVStore(k.storeKey)
    key := getFriendKey(ownerAddr, friendAddr)

    bz := store.Get([]byte(key))
    if bz != nil {
        return errors.New("Friend existed already")
    }

    friend := types.NewFriend()
    friend.Address = friendAddr
    friend.Name    = friendName

    store.Set([]byte(key), k.cdc.MustMarshalBinaryBare(friend))
    return nil
}

func (k Keeper) deleteFriend(ctx sdk.Context, ownerAddr string, friendAddr string) error {
    store := ctx.KVStore(k.storeKey)
    key := getFriendKey(ownerAddr, friendAddr)
    if store.Has([]byte(key)) {
        store.Delete([]byte(key))
        return nil
    }
    return errors.New("Friend doesn't exist")
}

func (k Keeper) addPendingFriend(ctx sdk.Context, ownerAddr string, ownerName string, friendAddr string) error {
    store := ctx.KVStore(k.storeKey)
    key := getPendingFriendKey(friendAddr, ownerAddr)
    bz := store.Get([]byte(key))
    if bz != nil {
        return nil  // no need to return error
    }

    friend := types.NewFriend()
    friend.Address = ownerAddr
    friend.Name = ownerName

    store.Set([]byte(key), k.cdc.MustMarshalBinaryBare(friend))
    return nil
}

func (k Keeper) getPendingFriend(ctx sdk.Context, ownerAddr string, friendAddr string) (types.Friend, error) {
    store := ctx.KVStore(k.storeKey)
    key := getPendingFriendKey(ownerAddr, friendAddr)
    bz := store.Get([]byte(key))
    if bz == nil {
        return types.Friend{}, errors.New("Pending friend not found")
    }
    var friend types.Friend
    k.cdc.MustUnmarshalBinaryBare(bz, &friend)
    return friend, nil
}

func (k Keeper) deletePendingFriend(ctx sdk.Context, ownerAddr string, friendAddr string) error {
    store := ctx.KVStore(k.storeKey)
    key := getPendingFriendKey(ownerAddr, friendAddr)
    store.Delete([]byte(key))
    return nil
}
