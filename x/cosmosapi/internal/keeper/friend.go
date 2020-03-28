package keeper

import (
    "errors"
    sdk "github.com/cosmos/cosmos-sdk/types"
    "github.com/yzhanginwa/cosmos-api/x/cosmosapi/internal/types"
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
