package keeper

import (
    "errors"
    sdk "github.com/cosmos/cosmos-sdk/types"
    "github.com/yzhanginwa/cosmos-api/x/cosmosapi/internal/types"
)

func (k Keeper) AddFriend(ctx sdk.Context, owner sdk.AccAddress, friendAddr string, friendName string) error {
    store := ctx.KVStore(k.storeKey)
    key := getFriendKey(owner.String(), friendAddr)
    
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

//////////////////////
//                  //
// helper functions //
//                  //
//////////////////////

