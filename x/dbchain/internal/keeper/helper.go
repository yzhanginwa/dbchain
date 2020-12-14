package keeper

import (
    "strconv"
    "github.com/yzhanginwa/dbchain/x/dbchain/internal/utils"
    sdk "github.com/cosmos/cosmos-sdk/types"
    "github.com/ipfs/go-cid"
)

func isSystemField(fieldName string) bool {
    systemFields := []string{"id", "created_by", "created_at"}
    return utils.ItemExists(systemFields, fieldName)
}

func (k Keeper) isAdmin(ctx sdk.Context, appId uint, addr sdk.AccAddress) bool {
    admins := k.getGroupMembers(ctx, appId, "admin")
    if utils.AddressIncluded(admins, addr) {
        return true
    }
    return false
}

func (k Keeper) isAuditor(ctx sdk.Context, appId uint, addr sdk.AccAddress) bool {
    auditors := k.getGroupMembers(ctx, appId, "audit")
    if utils.AddressIncluded(auditors, addr) {
        return true
    }
    return false
}

func (k Keeper) validateNotNullField(ctx sdk.Context, appId uint, tableName, fieldName string) bool {
    store := DbChainStore(ctx, k.storeKey)

    start, end := getFieldDataIteratorStartAndEndKey(appId, tableName, fieldName)
    iter := store.Iterator([]byte(start), []byte(end))
    var mold string
    var lastId uint = 0
    for ; iter.Valid(); iter.Next() {
        if iter.Error() != nil{
            return false
        }
        key := iter.Key()
        id := getIdFromDataKey(key)
        val := iter.Value()
        k.cdc.MustUnmarshalBinaryBare(val, &mold)
        if mold == "" || (id - lastId > 1) {
            return false
        }
        lastId = id
    }
    return true
}

func (k Keeper) validateIntField(ctx sdk.Context, appId uint, tableName, fieldName string) bool {
    store := DbChainStore(ctx, k.storeKey)

    start, end := getFieldDataIteratorStartAndEndKey(appId, tableName, fieldName)
    iter := store.Iterator([]byte(start), []byte(end))
    var mold string

    for ; iter.Valid(); iter.Next() {
        if iter.Error() != nil{
            return false
        }
        val := iter.Value()
        k.cdc.MustUnmarshalBinaryBare(val, &mold)
        if _, err := strconv.Atoi(mold); err != nil {
            return false
        }
    }
    return true
}

func (k Keeper) validateFileField(ctx sdk.Context, appId uint, tableName, fieldName string) bool {
    store := DbChainStore(ctx, k.storeKey)

    start, end := getFieldDataIteratorStartAndEndKey(appId, tableName, fieldName)
    iter := store.Iterator([]byte(start), []byte(end))
    var mold string

    for ; iter.Valid(); iter.Next() {
        if iter.Error() != nil{
            return false
        }
        val := iter.Value()
        k.cdc.MustUnmarshalBinaryBare(val, &mold)
        if _, err := cid.Decode(mold); err != nil {
            return false
        }
    }
    return true
}

func (k Keeper) validateOwnField(ctx sdk.Context, appId uint, tableName, fieldName string) bool {
    foreignTableName, ok := utils.GetTableNameFromForeignKey(fieldName)
    if !ok {
        return false
    }

    store := DbChainStore(ctx, k.storeKey)

    start, end := getFieldDataIteratorStartAndEndKey(appId, tableName, fieldName)
    iter := store.Iterator([]byte(start), []byte(end))
    var mold string

    for ; iter.Valid(); iter.Next() {
        if iter.Error() != nil{
            return false
        }
        key := iter.Key()
        id := getIdFromDataKey(key)
        owner, err := k.FindField(ctx, appId, tableName, id, "created_by")
        if err != nil {
            return false
        }

        val := iter.Value()
        k.cdc.MustUnmarshalBinaryBare(val, &mold)
        if !k.hasForeignRecordOfOwn(ctx, appId, foreignTableName, mold, owner) {
            return false
        }
    }
    return true
}

func (k Keeper) validateOwnFieldOnValue(ctx sdk.Context, appId uint, fieldName, value string, owner sdk.AccAddress) bool {
    foreignTableName, ok := utils.GetTableNameFromForeignKey(fieldName)
    if !ok {
        return false
    }

    return k.hasForeignRecordOfOwn(ctx, appId, foreignTableName, value, owner.String())
}

func (k Keeper) hasForeignRecordOfOwn(ctx sdk.Context, appId uint, tableName, idStr, ownerAddr string) bool {
    u64, err := strconv.ParseUint(idStr, 10, 64)
    if err != nil {
        return false
    }
    foreignOwner, err := k.FindField(ctx, appId, tableName, uint(u64), "created_by")
    if err == nil {
        if foreignOwner != ownerAddr {
            return false
        }
    } else {
        return false
    }

    return true
}
