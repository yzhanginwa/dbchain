package keeper

import (
    "fmt"
    "errors"
    sdk "github.com/cosmos/cosmos-sdk/types"
    "github.com/yzhanginwa/dbchain/x/dbchain/internal/utils"
)

// for now we only support indexes on single field
func (k Keeper) CreateIndex(ctx sdk.Context, appId uint, owner sdk.AccAddress, tableName string, fieldName string) error {
    // exclude the unwanted fields
    if fieldName == "id" {
        return errors.New("No index can be created on field id")
    }

    store := ctx.KVStore(k.storeKey)
    key := getMetaTableIndexKey(appId, tableName)
    var indexFields []string

    bz := store.Get([]byte(key))
    if bz != nil {
        k.cdc.MustUnmarshalBinaryBare(bz, &indexFields)
    }
    if utils.StringIncluded(indexFields, fieldName) {
        return errors.New(fmt.Sprintf("Fields %s is indexed already!", fieldName))
    }
    indexFields = append(indexFields, fieldName)
    store.Set([]byte(key), k.cdc.MustMarshalBinaryBare(indexFields))

    if err := createIndexData(k, ctx, appId, tableName, fieldName); err != nil {
        return err
    }

    return nil
}

func (k Keeper) DropIndex(ctx sdk.Context, appId uint, owner sdk.AccAddress, tableName string, fieldName string) error {
    store := ctx.KVStore(k.storeKey)

    indexFields, err := k.GetIndexFields(ctx, appId, tableName)
    if err != nil {
        return errors.New(fmt.Sprintf("Table %s does not have any index yet!", tableName))
    }

    if !utils.StringIncluded(indexFields, fieldName) {
        return errors.New(fmt.Sprintf("Table %s does not have index on %s yet!", tableName, fieldName))
    }

    utils.RemoveStringFromSet(indexFields, fieldName)
    key := getMetaTableIndexKey(appId, tableName)
    if len(indexFields) < 1 {
        store.Delete([]byte(key))
    } else {
        store.Set([]byte(key), k.cdc.MustMarshalBinaryBare(indexFields))
    }

    // to delete index data for the existing records of the table
    if err := dropIndexData(k, ctx, appId, tableName, fieldName); err != nil {
        return err
    }
    return nil
}

func (k Keeper) GetIndexFields(ctx sdk.Context, appId uint, tableName string) ([]string, error) {
    store := ctx.KVStore(k.storeKey)
    key := getMetaTableIndexKey(appId, tableName)
    bz := store.Get([]byte(key))
    if bz == nil {
        return []string{}, nil
    }

    var index_fields []string
    k.cdc.MustUnmarshalBinaryBare(bz, &index_fields)
    return index_fields, nil
}

func (k Keeper) appendIndexForRow(ctx sdk.Context, appId uint, tableName string, id uint) (uint, error){
    store := ctx.KVStore(k.storeKey)
    indexFields, err := k.GetIndexFields(ctx, appId, tableName)
    if err != nil {
        return 0, errors.New(fmt.Sprintf("Failed to get index for table %s", tableName))
    }
    if id == 0 {
        return 0, errors.New(fmt.Sprintf("Id for table %s is invalid", tableName))
    }

    var mold []string
    for _, indexField := range indexFields {
        value, err := k.FindField(ctx, appId, tableName, id, indexField)
        if err != nil {
            return id, nil    // the value for this field is empty. we don't need to do anything. Because people would not search on an empty value.
        }
        key := getIndexKey(appId, tableName, indexField, value)
        bz := store.Get([]byte(key))
        if bz != nil {
            k.cdc.MustUnmarshalBinaryBare(bz, &mold)
        }
        mold = append(mold, fmt.Sprint(id))
        store.Set([]byte(key), k.cdc.MustMarshalBinaryBare(mold))
    }

    return id, nil
}

//////////////////////
//                  //
// helper functions //
//                  //
//////////////////////

func createIndexData(k Keeper, ctx sdk.Context, appId uint, tableName, fieldName string) error {
    var dataValue string
    var indexValue []string

    store := ctx.KVStore(k.storeKey)
    start, end := getFieldDataIteratorStartAndEndKey(appId, tableName, fieldName)
    iter := store.Iterator([]byte(start), []byte(end))
    for ; iter.Valid(); iter.Next() {
        dataKey := iter.Key()
        id := getIdFromDataKey(dataKey)
        if isRowFrozen(store, appId, tableName, id) {
            continue
        }

        val := iter.Value()
        k.cdc.MustUnmarshalBinaryBare(val, &dataValue)

        indexKey := getIndexKey(appId, tableName, fieldName, dataValue)
        bz := store.Get([]byte(indexKey))
        if bz != nil {
            k.cdc.MustUnmarshalBinaryBare(bz, &indexValue)
        } else {
            indexValue = []string{}
        }
        indexValue = append(indexValue, fmt.Sprint(id))
        store.Set([]byte(indexKey), k.cdc.MustMarshalBinaryBare(indexValue))
    }
    return nil
}

func dropIndexData(k Keeper, ctx sdk.Context, appId uint, tableName, fieldName string) error {
    store := ctx.KVStore(k.storeKey)
    start, end := getFieldDataIteratorStartAndEndKey(appId, tableName, fieldName)
    iter := store.Iterator([]byte(start), []byte(end))
    for ; iter.Valid(); iter.Next() {
        indexKey := iter.Key()
        store.Delete([]byte(indexKey))
    }
    return nil
}
