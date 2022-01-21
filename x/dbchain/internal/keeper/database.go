package keeper

import (
    "bytes"
    "crypto/sha256"
    "errors"
    "fmt"
    sdk "github.com/cosmos/cosmos-sdk/types"
    "github.com/mr-tron/base58"
    "github.com/yzhanginwa/dbchain/x/dbchain/internal/keeper/cache"
    "github.com/yzhanginwa/dbchain/x/dbchain/internal/other"
    "github.com/yzhanginwa/dbchain/x/dbchain/internal/types"
    "github.com/yzhanginwa/dbchain/x/dbchain/internal/utils"
    "strconv"
    "strings"
    "time"
)

func (k Keeper) GetDatabaseAdmins(ctx sdk.Context, appId uint) []sdk.AccAddress {
    return k.getGroupMembers(ctx, appId, "admin")
}

func (k Keeper) GetAllAppCode(ctx sdk.Context) ([]string) {
    store := DbChainStore(ctx, k.storeKey)
    start, end := getDatabaseIteratorStartAndEndKey()
    iter := store.Iterator([]byte(start), []byte(end))
    var result []string

    for ; iter.Valid(); iter.Next() {
        if iter.Error() != nil{
            return nil
        }
        key := iter.Key()
        keyString := string(key)
        appCode := getAppCodeFromDatabaseKey(keyString)
        //
        _, err := k.getDatabase(ctx, appCode)
        if err != nil {
            continue
        }
        result = append(result, appCode)
    }

    return result
}

func (k Keeper) getAdminAppCode(ctx sdk.Context, address sdk.AccAddress) ([]string) {
    all := k.GetAllAppCode(ctx)
    var result []string

    for _, appCode := range all {
        appId, err := k.GetDatabaseIdWithoutCheck(ctx, appCode)
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
    } else if db.Deleted == true {
        return 0, errors.New("database has been deleted")
    } else {
        return db.AppId, nil
    }
}

func (k Keeper) GetDatabaseIdNotFrozen(ctx sdk.Context, appCode string) (uint, error){
    db, err := k.getDatabase(ctx, appCode)
    if err != nil {
        return 0, err
    }  else if db.Deleted == true || db.SchemaFrozen == true {
        return 0, errors.New("database has been deleted or frozen")
    } else {
        return db.AppId, nil
    }
}

func (k Keeper) GetDatabaseIdDataNotFrozen(ctx sdk.Context, appCode string) (uint, error){
    db, err := k.getDatabase(ctx, appCode)
    if err != nil {
        return 0, err
    }  else if db.Deleted == true || db.DataFrozen == true {
        return 0, errors.New("database has been deleted or data of database has been frozen")
    } else {
        return db.AppId, nil
    }
}

func (k Keeper) GetDatabaseIdSchemaDataNotFrozen(ctx sdk.Context, appCode string) (uint, error){
    db, err := k.getDatabase(ctx, appCode)
    if err != nil {
        return 0, err
    }  else if db.Deleted == true || db.DataFrozen == true || db.SchemaFrozen == true {
        return 0, errors.New("database has been deleted or frozen")
    } else {
        return db.AppId, nil
    }
}

func (k Keeper) GetDatabaseIdWithoutCheck(ctx sdk.Context, appCode string) (uint, error){
    db, err := k.getDatabase(ctx, appCode)
    if err != nil {
        return 0, err
    }  else {
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
    store := DbChainStore(ctx, k.storeKey)
    key := getDatabaseKey(appCode)
    bz, err := store.Get([]byte(key))
    if err != nil{
        return types.Database{},err
    }
    if bz == nil {
        return types.Database{}, errors.New(fmt.Sprintf("App code %s is invalid!", appCode))
    }
    var database types.Database
    k.cdc.MustUnmarshalBinaryBare(bz, &database)
    return database, nil
}

func (k Keeper) CreateDatabase(ctx sdk.Context, owner sdk.AccAddress, name string, description string, permissionRequired bool, system bool) error {
    store := DbChainStore(ctx, k.storeKey)
    newAppCode := generateNewAppCode(owner)
    if system {
        newAppCode = "0000000001"
    }

    key := getDatabaseKey(newAppCode)
    bz, err := store.Get([]byte(key))
    if err != nil{
        return err
    }
    if bz != nil {
        return errors.New(fmt.Sprintf("Application code %s existed already!", newAppCode))
    }

    appId, _ := registerDatabaseId(k, ctx, newAppCode)
    db := types.NewDatabase()
    db.Owner = owner
    db.Name = name
    db.Description = description
    db.PermissionRequired = permissionRequired
    db.Deleted = false
    db.Expiration = 0
    db.AppCode = newAppCode
    db.AppId = appId
    err = store.Set([]byte(key), k.cdc.MustMarshalBinaryBare(db))
    if err != nil{
        return err
    }

    // Add owner into the admin group of the database
    k.ModifyGroup(ctx, appId, "add", "admin")
    k.ModifyGroupMember(ctx, appId, "admin", "add", owner)

    // create auditor group without adding member for now
    k.ModifyGroup(ctx, appId, "add", "auditor")

    // Add owner as one of database users if this application requires permission for users
    if permissionRequired {
        if err := k.ModifyDatabaseUser(ctx, owner, newAppCode, "add", owner); err != nil {
            return errors.New("Failed to add owner as database user!")
        }
    }
    return nil 
}

func (k Keeper) DeleteApplication(ctx sdk.Context, appcode string) error{
    store := DbChainStore(ctx, k.storeKey)
    appKey := getDatabaseKey(appcode)
    appId, err := k.GetDatabaseId(ctx, appcode)
    if err != nil {
        return err
    }
    tables := k.GetTables(ctx, appId)
    bz, err :=store.Get([]byte(appKey))
    if err != nil || bz == nil {
        return err
    }
    var db types.Database
    k.cdc.MustUnmarshalBinaryBare(bz, &db)
    cache.VoidDatabase(appcode)

    if db.Deleted != false {
        return errors.New(fmt.Sprintf("database %s has been deleted", appcode))
    }
    db.Deleted = true
    //Due 30 days later
    t := time.Now().Add(time.Second * 86400 * 30)
    db.Expiration = t.Unix()
    store.Set([]byte(appKey), k.cdc.MustMarshalBinaryBare(db))
    for _, table := range tables {
        cache.VoidTable(appId, table)
    }
    return nil
}

func (k Keeper) PurgeApplication(ctx sdk.Context, appcode string) {
    store := DbChainStore(ctx, k.storeKey)
    appKey := getDatabaseKey(appcode)
    appId, err := k.GetDatabaseIdWithoutCheck(ctx, appcode)
    if err != nil {
        return
    }
    tables := k.GetTables(ctx, appId)
    bz, err :=store.Get([]byte(appKey))
    if err != nil || bz == nil {
        return
    }
    var db types.Database
    k.cdc.MustUnmarshalBinaryBare(bz, &db)

    store.Delete([]byte(appKey))
    for _, table := range tables {
        k.DropTable(ctx, appId, db.Owner, table)
    }
    //drop functions
    functions := k.GetFunctions(ctx, appId, 0)
    for _, name := range functions {
        k.DropFunction(ctx, appId, db.Owner, name, 0)
    }
    //drop querier
    queriers := k.GetFunctions(ctx, appId, 1)
    for _, name := range queriers {
        k.DropFunction(ctx, appId, db.Owner, name, 1)
    }
}

func (k Keeper) RecoverApplication(ctx sdk.Context, appcode string) {
    store := DbChainStore(ctx, k.storeKey)
    appKey := getDatabaseKey(appcode)

    bz, err := store.Get([]byte(appKey))
    if err != nil || bz == nil {
        return
    }
    var db types.Database
    k.cdc.MustUnmarshalBinaryBare(bz, &db)
    db.Deleted = false
    db.Expiration = 0
    store.Set([]byte(appKey), k.cdc.MustMarshalBinaryBare(db))
    cache.VoidDatabase(appcode)
}

func (k Keeper) ModifyDatabaseUser(ctx sdk.Context, owner sdk.AccAddress, appCode, action string, user sdk.AccAddress) error {
    store := DbChainStore(ctx, k.storeKey)

    appId, err := k.GetDatabaseId(ctx, appCode)
    if err != nil {
        return errors.New(fmt.Sprintf("Application code %s does not exist!", appCode))
    }


    if k.DatabaseUserExists(ctx, appId, user) {
        if action == "add" {
            return errors.New("Database user existed already!")
        } else {
            key := getDatabaseUserKey(appId, user.String())
            err := store.Delete([]byte(key))
            if err != nil{
                return err
            }
        }
    } else {
        if action == "add" {
            key := getDatabaseUserKey(appId, user.String())
            err := store.Set([]byte(key), k.cdc.MustMarshalBinaryBare(user))
            if err != nil{
                return err
            }
        } else {
            return errors.New("Database user does not exist!")
        }
    }
    return nil
}

func (k Keeper) DatabaseUserExists(ctx sdk.Context, appId uint, user sdk.AccAddress) bool {
    store := DbChainStore(ctx, k.storeKey)

    key := getDatabaseUserKey(appId, user.String())
    bz, err := store.Get([]byte(key))
    if bz == nil || err != nil{
        return false
    }
    return true
}

func (k Keeper) SetAppUserFileVolumeLimit(ctx sdk.Context, appId uint, size  string) error {
    store := DbChainStore(ctx, k.storeKey)
    key := getDatabaseUserFileVolumeLimitKey(appId)
    err := store.Set([]byte(key), k.cdc.MustMarshalBinaryBare(size))
    return err
}

func (k Keeper) UpdateAppUserUsedFileVolume(ctx sdk.Context, appId uint, user string, size  string) error {
    store := DbChainStore(ctx, k.storeKey)
    key := GetDatabaseUserUsedFileVolumeLimitKey(appId, user)
    newSize := ""
    bz, err := store.Get([]byte(key))
    if err != nil ||  bz == nil{
        newSize = size
    } else {
        originSize := ""
        k.cdc.MustUnmarshalBinaryBare(bz, &originSize)
        iOriginSize, _ :=  strconv.ParseInt(originSize, 10, 64)
        iSize, err := strconv.ParseInt(size, 10, 64)
        if err != nil {
            return errors.New("file size err")
        }
        newSize = fmt.Sprintf("%d", iSize + iOriginSize)
    }
    err = store.Set([]byte(key), k.cdc.MustMarshalBinaryBare(newSize))
    return err
}

func (k Keeper) RestoreVolume(ctx sdk.Context, appId uint, cids []string, user string) error {
    store := DbChainStore(ctx, k.storeKey)
    key := GetDatabaseUserUsedFileVolumeLimitKey(appId, user)
    bz, err := store.Get([]byte(key))
    if err != nil || bz == nil{
        return errors.New("get user used size err")
    }
    usedSize := ""
    k.cdc.MustUnmarshalBinaryBare(bz, &usedSize)
    iUsedSize, _ :=  strconv.ParseUint(usedSize, 10, 64)

    sh := utils.NewShellDbchain("localhost:5001")
    for _, cid := range cids {
        size,err := getUploadFileSize(sh, cid)
        if err != nil {
            continue
        }
        if iUsedSize > size {
            iUsedSize -= size
        }
    }
    newSize := fmt.Sprintf("%d", iUsedSize)
    err = store.Set([]byte(key), k.cdc.MustMarshalBinaryBare(newSize))
    return err
}

func (k Keeper) findFileTypeField(ctx sdk.Context, appId uint, tableName string) []string {
    fileFields := make([]string, 0)
    fields, err := k.getTableFields(ctx, appId, tableName)
    if err != nil {
        return fileFields
    }
    for _, field := range fields {
        dataType, err := k.GetColumnDataType(ctx, appId, tableName, field)
        if err != nil {
            continue
        }
        if dataType == string(types.FLDTYP_FILE) {
            fileFields = append(fileFields, dataType)
        }
    }
    return fileFields
}

func (k Keeper) GetApplicationUserFileVolumeLimit(ctx sdk.Context, appId uint) string {
    store := DbChainStore(ctx, k.storeKey)
    key := getDatabaseUserFileVolumeLimitKey(appId)
    bz, err := store.Get([]byte(key))
    if err != nil || bz == nil {
        return "no limit"
    }
    size := ""
    k.cdc.MustUnmarshalBinaryBare(bz, &size)
    return size
}

func (k Keeper) GetApplicationUserUsedFileVolume(ctx sdk.Context, appId uint, user sdk.AccAddress) string {
    store := DbChainStore(ctx, k.storeKey)
    key := GetDatabaseUserUsedFileVolumeLimitKey(appId, user.String())
    bz, err := store.Get([]byte(key))
    if err != nil || bz == nil {
        return "0"
    }
    size := ""
    k.cdc.MustUnmarshalBinaryBare(bz, &size)
    return size
}

func (k Keeper) GetDatabaseUsers(ctx sdk.Context, appId uint, owner sdk.AccAddress) []string {
    store := DbChainStore(ctx, k.storeKey)
    start, end := getDatabaseUserIteratorStartAndEndKey(appId)
    iter := store.Iterator([]byte(start), []byte(end))
    var result = []string{}

    for ; iter.Valid(); iter.Next() {
        if iter.Error() != nil{
            return nil
        }
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
    if database.PermissionRequired {
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

    store := DbChainStore(ctx, k.storeKey)
    key := getDatabaseKey(appCode)
    bz, err := store.Get([]byte(key))
    if err != nil{
        return err
    }
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

func (k Keeper) SetDatabaseDataStatus(ctx sdk.Context, owner sdk.AccAddress, appCode, status string) error {
    frozen_status := true
    if status == "unfrozen" {    // status must be either "frozen" or "unfrozen"
        frozen_status = false
    }

    store := DbChainStore(ctx, k.storeKey)
    key := getDatabaseKey(appCode)
    bz, err := store.Get([]byte(key))
    if err != nil{
        return err
    }
    if bz == nil {
        return errors.New(fmt.Sprintf("App code %s is invalid!", appCode))
    }
    var database types.Database
    k.cdc.MustUnmarshalBinaryBare(bz, &database)

    if database.DataFrozen == frozen_status {
        return errors.New("No need to do anything!")
    }

    database.DataFrozen = frozen_status
    store.Set([]byte(key), k.cdc.MustMarshalBinaryBare(database))
    cache.VoidDatabase(appCode)
    return nil
}

func (k Keeper) SetDatabasePermission(ctx sdk.Context, owner sdk.AccAddress, appCode, permissionRequired string) error {
    permissionStatus := true
    if permissionRequired == "unrequired" {    // permission must be either "required" or "unrequired"
        permissionStatus = false
    }

    store := DbChainStore(ctx, k.storeKey)
    key := getDatabaseKey(appCode)
    bz, err := store.Get([]byte(key))
    if err != nil{
        return err
    }
    if bz == nil {
        return errors.New(fmt.Sprintf("App code %s is invalid!", appCode))
    }
    var database types.Database
    k.cdc.MustUnmarshalBinaryBare(bz, &database)

    if database.PermissionRequired == permissionStatus {
        return errors.New("No need to do anything!")
    }

    database.PermissionRequired = permissionStatus
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

