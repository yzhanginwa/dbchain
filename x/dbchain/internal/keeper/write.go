package keeper

import (
    "fmt"
    "strconv"
    "errors"
    sdk "github.com/cosmos/cosmos-sdk/types"
    "github.com/yzhanginwa/dbchain/x/dbchain/internal/other"
    "github.com/yzhanginwa/dbchain/x/dbchain/internal/types"
    "github.com/yzhanginwa/dbchain/x/dbchain/internal/utils"
    shell "github.com/ipfs/go-ipfs-api"
)


func (k Keeper) Insert(ctx sdk.Context, appId uint, tableName string, fields types.RowFields, owner sdk.AccAddress) (uint, error){
    if !k.IsDatabaseUser(ctx, appId, owner) {
        return 0, errors.New(fmt.Sprintf("Do not have user permission on database %d", appId))
    }
    if(!k.haveWritePermission(ctx, appId, tableName, owner)) {
        return 0, errors.New(fmt.Sprintf("Do not have permission inserting table %s", tableName))
    }

    if(!k.validateInsertion(ctx, appId, tableName, fields, owner)) {
        return 0, errors.New(fmt.Sprintf("Failed validation when inserting table %s", tableName))
    }

    // as far the first go routine to be used
    go k.tryToPinFile(ctx, appId, tableName, fields, owner)

    id, err := getNextId(k, ctx, appId, tableName)
    if err != nil {
        return 0, errors.New(fmt.Sprintf("Failed to get id for table %s", tableName))
    }

    // to set the 2 special fields
    fields["id"] = strconv.Itoa(int(id))
    fields["created_by"] = owner.String()
    fields["created_at"] = other.GetCurrentBlockTime().String()

    k.Write(ctx, appId, tableName, id, fields, owner)
    k.updateIndex(ctx, appId, tableName, id, fields)
    return id, nil
}


// TODO: need to think over how and when to allow updating
func (k Keeper) Update(ctx sdk.Context, appId uint, tableName string, id uint, fields types.RowFields, owner sdk.AccAddress) (uint, error){
    // TODO: need to check the ownership of the record
    k.Write(ctx, appId, tableName, id, fields, owner)
    k.updateIndex(ctx, appId, tableName, id, fields)
    return id, nil
}


func (k Keeper) Write(ctx sdk.Context, appId uint, tableName string, id uint, fields types.RowFields, owner sdk.AccAddress) (uint, error){
    store := ctx.KVStore(k.storeKey)

    fieldNames, err := k.getTableFields(ctx, appId, tableName)
    if err != nil {
        return 0, errors.New(fmt.Sprintf("Failed to get fields for table %s", tableName))
    }

    if id == 0 {
        return 0, errors.New(fmt.Sprintf("Id for table %s is invalid", tableName))
    }

    for _, fieldName := range fieldNames {
        if value, ok := fields[fieldName]; ok {
            key := getDataKey(appId, tableName, id, fieldName)
            store.Set([]byte(key), k.cdc.MustMarshalBinaryBare(value)) 
        }
    }

    return id, nil
}

func (k Keeper) Delete(ctx sdk.Context, appId uint, tableName string, id uint, owner sdk.AccAddress) (uint, error){
    store := ctx.KVStore(k.storeKey)

    fieldNames, err := k.getTableFields(ctx, appId, tableName)
    if err != nil {
        return 0, errors.New(fmt.Sprintf("Failed to get fields for table %s", tableName))
    }

    if id == 0 {
        return 0, errors.New("Id cannot be empty")
    }

    for _, fieldName := range fieldNames {
        key := getDataKey(appId, tableName, id, fieldName)
        store.Delete([]byte(key)) 
    }

    // TODO: to remove the related indexes
    return id, nil
}

func (k Keeper) Freeze(ctx sdk.Context, appId uint, tableName string, id uint, owner sdk.AccAddress) (uint, error){
    store := ctx.KVStore(k.storeKey)

    if id == 0 {
        return 0, errors.New("Id cannot be empty")
    }

    keyAt := getDataKey(appId, tableName, id, types.FLD_FROZEN_AT)
    bz := store.Get([]byte(keyAt))
    if bz != nil {
        return id, nil
    }
    store.Set([]byte(keyAt), k.cdc.MustMarshalBinaryBare(other.GetCurrentBlockTime().String()))

    keyBy := getDataKey(appId, tableName, id, types.FLD_FROZEN_BY)
    store.Set([]byte(keyBy), k.cdc.MustMarshalBinaryBare(owner.String()))

    // TODO: to remove the related indexes
    return id, nil
}

//////////////////
//              //
// helper funcs //
//              //
//////////////////

func isSystemField(fieldName string) bool {
    systemFields := []string{"id", "created_by", "created_at"}
    return utils.ItemExists(systemFields, fieldName)
}

func (k Keeper) haveWritePermission(ctx sdk.Context, appId uint, tableName string, owner sdk.AccAddress) bool {
    options, _ := k.GetOption(ctx, appId, tableName)
    if utils.ItemExists(options, string(types.TBLOPT_ADMIN_ONLY)) {
        admins := k.ShowAdminGroup(ctx, appId)
        if utils.AddressIncluded(admins, owner) {
            return true
        }
        return false
    }
    return true
}

// for now, we check the filed non-null option
func (k Keeper) validateInsertion(ctx sdk.Context, appId uint, tableName string, fields types.RowFields, owner sdk.AccAddress) bool {
    fieldNames, err := k.getTableFields(ctx, appId, tableName)
    if err != nil {
        return(false)
    }

    for _, fieldName := range fieldNames {
        if(isSystemField(fieldName)) {
            continue
        }
        fieldOptions, _ := k.GetColumnOption(ctx, appId, tableName, fieldName)
        // TODO: use a constant for the possible options
        if(utils.ItemExists(fieldOptions, string(types.FLDOPT_NOTNULL))) {
            if value, ok := fields[fieldName]; ok {
                if(len(value) < 1) {
                    return(false)
                }
            } else {
              return(false)
            }
        }

        if(utils.ItemExists(fieldOptions, string(types.FLDOPT_UNIQUE))) {
            if value, ok := fields[fieldName]; ok {
                if(len(value)>0) {
                    ids := k.FindBy(ctx, appId, tableName, fieldName, value, owner)
                    if len(ids) > 0 {
                        return(false)
                    }
                }
            }
        }

        if(utils.ItemExists(fieldOptions, string(types.FLDOPT_OWN))) {
            if tn, ok := utils.GetTableNameFromForeignKey(fieldName); ok {
                foreignId := fields[fieldName]
                u64, err := strconv.ParseUint(foreignId, 10, 64)
                if err != nil {
                    return false
                }
                foreignOwner, err := k.FindField(ctx, appId, tn, uint(u64), "created_by")
                if err == nil {
                    if foreignOwner != owner.String() {
                        return false
                    }
                } else {
                    return false
                }
            } else {
                return false
            }
        }
    }
    return(true)
}

func (k Keeper) tryToPinFile(ctx sdk.Context, appId uint, tableName string, fields types.RowFields, owner sdk.AccAddress) bool {
    fieldNames, err := k.getTableFields(ctx, appId, tableName)
    if err != nil {
        return(false)
    }

    for _, fieldName := range fieldNames {
        if(isSystemField(fieldName)) {
            continue
        }
        fieldOptions, _ := k.GetColumnOption(ctx, appId, tableName, fieldName)
        if(utils.ItemExists(fieldOptions, string(types.FLDOPT_FILE))) {
            if value, ok := fields[fieldName]; ok {
                sh := shell.NewShell("localhost:5001")
                err =sh.Pin(value)
                if err != nil {
                    logger := k.Logger(ctx)
                    logger.Error(fmt.Sprintf("Failed to pin ipfs cid %s", value))
                    return false
                }
            }
        }
    }
    return true
}
