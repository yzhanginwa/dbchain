package keeper

import (
    "errors"
    sdk "github.com/cosmos/cosmos-sdk/types"
    "github.com/yzhanginwa/cosmos-api/x/cosmosapi/internal/types"
)

func (k Keeper) AddFriend(ctx sdk.Context, owner sdk.AccAddress, ownerName string, friendAddr string, friendName string) error {
    store := ctx.KVStore(k.storeKey)
    ownerStr := owner.String()
    key := getFriendKey(ownerStr, friendAddr)

    bz := store.Get([]byte(key))
    if bz != nil {
        return errors.New("Friend existed already")
    }

    friend := types.NewFriend()
    friend.Address = friendAddr
    friend.Name    = friendName

    store.Set([]byte(key), k.cdc.MustMarshalBinaryBare(friend))

    // try to add self as pending friend of the friend

    key = getPendingFriendKey(friendAddr, ownerStr)
    bz = store.Get([]byte(key))
    if bz != nil {
        return nil  // no need to return error
    }

    friend =types.NewFriend()
    friend.Address = ownerStr
    friend.Name = ownerName

    store.Set([]byte(key), k.cdc.MustMarshalBinaryBare(friend)) 
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

//////////////////////
//                  //
// helper functions //
//                  //
//////////////////////

