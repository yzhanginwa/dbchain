package keeper

import (
    "crypto/sha256"
    "github.com/mr-tron/base58"
    "fmt"
    "strings"
    "errors"
    "bytes"
    sdk "github.com/cosmos/cosmos-sdk/types"
    "github.com/yzhanginwa/dbchain/x/dbchain/internal/types"
    "github.com/yzhanginwa/dbchain/x/dbchain/internal/other"
    "github.com/yzhanginwa/dbchain/x/dbchain/internal/keeper/cache"
)

func (k Keeper) GetDatabaseAdmins(ctx sdk.Context, appId uint) []sdk.AccAddress {
    return k.getGroupMembers(ctx, appId, "admin")
}

func (k Keeper) GetAllAppCode(ctx sdk.Context) ([]string) {
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

func (k Keeper) getAdminAppCode(ctx sdk.Context, address sdk.AccAddress) ([]string) {
    all := k.GetAllAppCode(ctx)
    var result []string

    for _, appCode := range all {
        appId, err := k.GetDatabaseId(ctx, appCode)
        if err != nil {
            return []string{}
        }
        adminAddresses := k.getGroupMembers(ctx, appId, "admin")
        for _, addr := range adminAddresses {
            if bytes.Compare(address, addr) == 0 {
                result = append(result, appCode)
                break
            }
        }
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


func (k Keeper) getDatabaseById(ctx sdk.Context, appId uint) (types.Database, error) {
    if appCode, ok := cache.GetAppCodeById(appId); ok {
        return k.getDatabase(ctx, appCode)
    } else {
        return types.Database{}, errors.New(fmt.Sprintf("AppID %d is invalid!", appId))
    }
}

func (k Keeper) getDatabase(ctx sdk.Context, appCode string) (types.Database, error) {
    if cached_db, ok := cache.GetDatabase(appCode); ok {
        return cached_db, nil
    }
    db, err := k.getDatabaseRaw(ctx, appCode)
    if err == nil {
        cache.SetDatabase(appCode, db)
    }
    return db, err
}

func (k Keeper) getDatabaseRaw(ctx sdk.Context, appCode string) (types.Database, error) {
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

func (k Keeper) CreateDatabase(ctx sdk.Context, owner sdk.AccAddress, name string, description string, permissioned bool, system bool) error {
    store := ctx.KVStore(k.storeKey)
    newAppCode := generateNewAppCode(owner)
    if system {
        newAppCode = "0000000001"
    }

    key := getDatabaseKey(newAppCode)
    bz := store.Get([]byte(key))
    if bz != nil {
        return errors.New(fmt.Sprintf("Application code %s existed already!", newAppCode))
    }

    appId, _ := registerDatabaseId(k, ctx, newAppCode)
    db := types.NewDatabase()
    db.Owner = owner
    db.Name = name
    db.Description = description
    db.Permissioned = permissioned
    db.AppCode = newAppCode
    db.AppId = appId
    store.Set([]byte(key), k.cdc.MustMarshalBinaryBare(db))

    // Add owner into the admin group of the database
    k.ModifyGroup(ctx, appId, "add", "admin")
    k.ModifyGroupMember(ctx, appId, "admin", "add", owner)

    // Add owner as one of database users if this application requires permission for users
    if permissioned {
        if err := k.AddDatabaseUser(ctx, owner, newAppCode, owner); err != nil {
            return errors.New("Failed to add owner as database user!")
        }
    }
    return nil 
}

func (k Keeper) AddDatabaseUser(ctx sdk.Context, owner sdk.AccAddress, appCode string, user sdk.AccAddress) error {
    store := ctx.KVStore(k.storeKey)

    appId, err := k.GetDatabaseId(ctx, appCode)
    if err != nil {
        return errors.New(fmt.Sprintf("Application code %s does not exist!", appCode))
    }

    key := getDatabaseUserKey(appId, user.String())
    store.Set([]byte(key), k.cdc.MustMarshalBinaryBare(user))
    return nil
}

func (k Keeper) DatabaseUserExists(ctx sdk.Context, appId uint, user sdk.AccAddress) bool {
    store := ctx.KVStore(k.storeKey)

    key := getDatabaseUserKey(appId, user.String())
    bz := store.Get([]byte(key))
    if bz == nil {
        return false
    }
    return true
}

func (k Keeper) GetDatabaseUsers(ctx sdk.Context, appId uint, owner sdk.AccAddress) []string {
    store := ctx.KVStore(k.storeKey)
    start, end := getDatabaseUserIteratorStartAndEndKey(appId)
    iter := store.Iterator([]byte(start), []byte(end))
    var result = []string{}

    for ; iter.Valid(); iter.Next() {
        key := iter.Key()
        keyString := string(key)
        user := getUserFromDatabaseUserKey(keyString)
        result = append(result, user)
    }

    return result
}

func (k Keeper) IsDatabaseUser(ctx sdk.Context, appId uint, owner sdk.AccAddress) bool {
    database, err := k.getDatabaseById(ctx, appId)
    if err != nil {
        return false
    }
    if database.Permissioned {
        if k.DatabaseUserExists(ctx, appId, owner) {
            return true
        } else {
            return false
        }
    } else {
        return true
    }
}

func (k Keeper) SetSchemaStatus(ctx sdk.Context, owner sdk.AccAddress, appCode, status string) error {
    frozen_status := true
    if status == "unfrozen" {    // status must be either "frozen" or "unfrozen"
        frozen_status = false
    }

    store := ctx.KVStore(k.storeKey)
    key := getDatabaseKey(appCode)
    bz := store.Get([]byte(key))
    if bz == nil {
        return errors.New(fmt.Sprintf("App code %s is invalid!", appCode))
    }
    var database types.Database
    k.cdc.MustUnmarshalBinaryBare(bz, &database)

    if database.SchemaFrozen == frozen_status {
        return errors.New("No need to do anything!")
    }

    database.SchemaFrozen = frozen_status
    store.Set([]byte(key), k.cdc.MustMarshalBinaryBare(database))
    cache.VoidDatabase(appCode)
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
    hashStr := base58.Encode(hashedBytes[:])
    code:= hashStr[:10]
    code = strings.ToUpper(code)
    return code
}

