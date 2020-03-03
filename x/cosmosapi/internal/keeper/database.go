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

func (k Keeper) GetDatabaseAdmins(ctx sdk.Context, appCode string) []sdk.AccAddress {
    // TODO: we'll have a better way to maintain and retrive a group of database specific admins
    database, err := k.getDatabase(ctx, appCode)
    if err != nil {
        return []sdk.AccAddress{}
    } else {
        return []sdk.AccAddress{database.Owner}
    }
}

func (k Keeper) getDatabases(ctx sdk.Context) ([]string) {
    store := ctx.KVStore(k.storeKey)
    start, end := getDatabaseIteratorStartAndEndKey()
    iter := store.Iterator([]byte(start), []byte(end))
    var result []string

    for ; iter.Valid(); iter.Next() {
        key := iter.Key()
        keyString := string(key)
        appCode := getAppCodeFromDatabaseKey(keyString)
        result = append(result, appCode)
    }

    return result
}

func (k Keeper) GetDatabaseId(ctx sdk.Context, appCode string) (uint, error) {
    db, err := k.getDatabase(ctx, appCode)
    if err != nil {
        return 0, err
    } else {
        return db.AppId, nil
    }
}

func (k Keeper) getDatabase(ctx sdk.Context, appCode string) (types.Database, error) {
    store := ctx.KVStore(k.storeKey)
    key := getDatabaseKey(appCode)
    bz := store.Get([]byte(key))
    if bz == nil {
        return types.Database{}, errors.New(fmt.Sprintf("App code %s is invalid!", appCode))
    }
    var database types.Database
    k.cdc.MustUnmarshalBinaryBare(bz, &database)
    return database, nil
}

func (k Keeper) CreateDatabase(ctx sdk.Context, owner sdk.AccAddress, description string) error {
    store := ctx.KVStore(k.storeKey)
    newAppCode := generateNewAppCode(owner)
    key := getDatabaseKey(newAppCode)
    bz := store.Get([]byte(key))
    if bz != nil {
        return errors.New(fmt.Sprintf("Application code %s existed already!", newAppCode))
    }

    appId, _ := registerDatabaseId(k, ctx, newAppCode)
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

