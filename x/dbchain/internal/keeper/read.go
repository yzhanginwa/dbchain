package keeper

import (
    "fmt"
    "errors"
    "strconv"
    sdk "github.com/cosmos/cosmos-sdk/types"
    "github.com/yzhanginwa/dbchain/x/dbchain/internal/types"
    "github.com/yzhanginwa/dbchain/x/dbchain/internal/utils"
)


func (k Keeper) DoFind(ctx sdk.Context, appId uint, tableName string, id uint) (types.RowFields, error){
    store := ctx.KVStore(k.storeKey)

    fieldNames, err := k.getTableFields(ctx, appId, tableName)
    if err != nil {
        return nil, errors.New(fmt.Sprintf("Failed to get fields for table %s", tableName))
    }

    if id == 0 {
        return nil, errors.New("Id cannot be 0")
    }

    var fields = make(types.RowFields)
    var value string

    for _, fieldName := range fieldNames {
        key := getDataKey(appId, tableName, fieldName, id)
        bz := store.Get([]byte(key)) 
        if bz != nil {
            k.cdc.MustUnmarshalBinaryBare(bz, &value)
            fields[fieldName] = value
        }
    }

    return fields, nil
}

func (k Keeper) FindField(ctx sdk.Context, appId uint, tableName string, id uint, fieldName string) (string, error){
    store := ctx.KVStore(k.storeKey)

    if !k.HasField(ctx, appId, tableName, fieldName) {
        return "", errors.New("Field not existed")
    }

    if id == 0 {
        return "", errors.New("Id cannot be 0")
    }

    key := getDataKey(appId, tableName, fieldName, id)
    bz := store.Get([]byte(key))
    if bz != nil {
        var value string
        k.cdc.MustUnmarshalBinaryBare(bz, &value)
        return value, nil
    }
    return "", errors.New("Field not found")
}

func (k Keeper) Find(ctx sdk.Context, appId uint, tableName string, id uint, owner sdk.AccAddress) (types.RowFields, error){
    var ids []uint
    ids = append(ids, id)

    // if public table, return all ids
    if !k.isTablePublic(ctx, appId, tableName) {
        ids = k.filterOwnIds(ctx, appId, tableName, ids, owner)
        if len(ids) < 1 {
            return nil, errors.New(fmt.Sprintf("Failed to get fields for id %d", id))
        }
    }

    return k.DoFind(ctx, appId, tableName, id)
}

// Find by an attribute in the r.Fields
func (k Keeper) FindBy(ctx sdk.Context, appId uint, tableName string, field string,  value string, owner sdk.AccAddress) []uint {
    store := ctx.KVStore(k.storeKey)

    var hasIndex bool
    indexFields, err := k.GetIndex(ctx, appId, tableName)
    if err == nil {
        for _, item := range(indexFields) {
            if item == field {
                hasIndex = true
                break
            }
        }
    }

    var result []uint
    if hasIndex {
        key := getIndexKey(appId, tableName, field, value)
        bz := store.Get([]byte(key))
        if bz == nil {
            result = []uint{}
        } else {
            k.cdc.MustUnmarshalBinaryBare(bz, &result)
        }
    } else {
        // partial table scanning
        start, end := getFieldDataIteratorStartAndEndKey(appId, tableName, field)
        iter := store.Iterator([]byte(start), []byte(end))
        var mold string
        for ; iter.Valid(); iter.Next() {
            key := iter.Key()
            keyString := string(key)
            fn := getFieldNameFromDataKey(keyString)
            if fn == field {
                val := iter.Value()
                k.cdc.MustUnmarshalBinaryBare(val, &mold)
                if mold == value {
                    id := getIdFromDataKey(keyString)
                    u64, _ := strconv.ParseUint(id, 10, 64)
                    result = append(result, uint(u64))
                }
            }
        }
    }

    // if public table, return all ids
    if k.isTablePublic(ctx, appId, tableName) {
        return result
    } else {
        return k.filterOwnIds(ctx, appId, tableName, result, owner)
    }
}

func (k Keeper) FindAll(ctx sdk.Context, appId uint, tableName string, owner sdk.AccAddress) []uint {
    store := ctx.KVStore(k.storeKey)
    var result []uint

    // full table scanning
    start, end := getFieldDataIteratorStartAndEndKey(appId, tableName, "id")
    iter := store.Iterator([]byte(start), []byte(end))
    var currentId = "0"
    for ; iter.Valid(); iter.Next() {
        key := iter.Key()
        keyString := string(key)
        id := getIdFromDataKey(keyString)
        if id != currentId {
            currentId = id
            u64, _ := strconv.ParseUint(id, 10, 64)
            result = append(result, uint(u64))
        }
    }

    // if public table, return all ids
    if k.isTablePublic(ctx, appId, tableName) {
        return result
    } else {
        return k.filterOwnIds(ctx, appId, tableName, result, owner)
    }
}

//////////////////
//              //
// helper funcs //
//              //
//////////////////

func (k Keeper) isTablePublic(ctx sdk.Context, appId uint, tableName string) bool {
    tableOptions, _ := k.GetOption(ctx, appId, tableName)
    return utils.ItemExists(tableOptions, string(types.TBLOPT_PUBLIC))
}

func (k Keeper) filterOwnIds(ctx sdk.Context, appId uint,  tableName string, ids []uint, owner sdk.AccAddress) []uint {
    store := ctx.KVStore(k.storeKey)
    var ownerString string = owner.String()

    var result = []uint{}
    var mold string
    for _, id := range ids {
        key := getDataKey(appId, tableName, "created_by", uint(id))
        bz := store.Get([]byte(key))
        if bz != nil {
            k.cdc.MustUnmarshalBinaryBare(bz, &mold)
            if mold == ownerString {
                result = append(result, uint(id))
            }
        }
    }
    return result
}
