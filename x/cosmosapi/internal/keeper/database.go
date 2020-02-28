package keeper

import (
    "crypto/sha256"
    "encoding/base64"
    "fmt"
    "errors"
    sdk "github.com/cosmos/cosmos-sdk/types"
    "github.com/yzhanginwa/cosmos-api/x/cosmosapi/internal/types"
    "github.com/yzhanginwa/cosmos-api/x/cosmosapi/internal/other"
)

func (k Keeper) CreateDatabase(ctx sdk.Context, owner sdk.AccAddress, description string) error {
    store := ctx.KVStore(k.storeKey)
    newAppCode := generateNewAppCode(owner)
    key := getDatabaseKey(newAppCode)
    bz := store.Get([]byte(key))
    if bz != nil {
        return errors.New(fmt.Sprintf("Application code %s existed already!", newAppCode))
    }

    appId, _ := getDatabaseId(k, ctx, newAppCode)
    db := types.NewDatabase()
    db.Owner = owner
    db.Description = description
    db.AppCode = newAppCode
    db.AppId = appId
    store.Set([]byte(key), k.cdc.MustMarshalBinaryBare(db))

    return nil 
}

////////////////////
//                //
// helper methods //
//                //
////////////////////

func generateNewAppCode(owner sdk.AccAddress) string {
    blockTime := other.GetCurrentBlockTime().String()
    hashedBytes := sha256.Sum256([]byte(blockTime + owner.String()))
    hashStr := base64.StdEncoding.EncodeToString(hashedBytes[:])
    return hashStr[:10]
}

