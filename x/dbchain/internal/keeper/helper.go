package keeper

import (
    "github.com/yzhanginwa/dbchain/x/dbchain/internal/utils"
    sdk "github.com/cosmos/cosmos-sdk/types"
    "strconv"
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

func (k Keeper) validateOwnField(ctx sdk.Context, appId uint, tableName, fieldName string, owner sdk.AccAddress) bool {
    foreignTableName, ok := utils.GetTableNameFromForeignKey(fieldName)
    if !ok {
        return false
    }

    store := ctx.KVStore(k.storeKey)
    ownerStr := owner.String()

    start, end := getFieldDataIteratorStartAndEndKey(appId, tableName, fieldName)
    iter := store.Iterator([]byte(start), []byte(end))
    var mold string

    for ; iter.Valid(); iter.Next() {
        val := iter.Value()
        k.cdc.MustUnmarshalBinaryBare(val, &mold)

        if !k.hasForeignRecord(ctx, appId, foreignTableName, mold, ownerStr) {
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

    return k.hasForeignRecord(ctx, appId, foreignTableName, value, owner.String())
}

func (k Keeper) hasForeignRecord(ctx sdk.Context, appId uint, tableName, idStr, ownerAddr string) bool {
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
