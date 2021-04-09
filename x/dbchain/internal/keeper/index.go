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

    store := DbChainStore(ctx, k.storeKey)
    key := getMetaTableIndexKey(appId, tableName)
    var indexFields []string

    bz, err := store.Get([]byte(key))
    if err != nil{
        return err
    }
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
    store := DbChainStore(ctx, k.storeKey)

    indexFields, err := k.GetIndexFields(ctx, appId, tableName)
    if err != nil {
        return errors.New(fmt.Sprintf("Table %s does not have any index yet!", tableName))
    }

    if !utils.StringIncluded(indexFields, fieldName) {
        return errors.New(fmt.Sprintf("Table %s does not have index on %s yet!", tableName, fieldName))
    }

    indexFields = utils.RemoveStringFromSet(indexFields, fieldName)
    key := getMetaTableIndexKey(appId, tableName)
    if len(indexFields) < 1 {
        err := store.Delete([]byte(key))
        if err != nil{
            return err
        }
    } else {
        err := store.Set([]byte(key), k.cdc.MustMarshalBinaryBare(indexFields))
        if err != nil{
            return err
        }
    }

    // to delete index data for the existing records of the table
    if err := dropIndexData(k, ctx, appId, tableName, fieldName); err != nil {
        return err
    }
    return nil
}

func (k Keeper) GetIndexFields(ctx sdk.Context, appId uint, tableName string) ([]string, error) {
    store := DbChainStore(ctx, k.storeKey)
    key := getMetaTableIndexKey(appId, tableName)
    bz, err := store.Get([]byte(key))
    if err != nil{
        return nil, err
    }
    if bz == nil {
        return []string{}, nil
    }

    var index_fields []string
    k.cdc.MustUnmarshalBinaryBare(bz, &index_fields)
    return index_fields, nil
}

func (k Keeper) appendIndexForRow(ctx sdk.Context, appId uint, tableName string, id uint) (uint, error){
    store := DbChainStore(ctx, k.storeKey)
    indexFields, err := k.GetIndexFields(ctx, appId, tableName)
    if err != nil {
        return 0, errors.New(fmt.Sprintf("Failed to get index for table %s", tableName))
    }
    if id == 0 {
        return 0, errors.New(fmt.Sprintf("Id for table %s is invalid", tableName))
    }


    for _, indexField := range indexFields {
        var mold []string
        value, err := k.FindField(ctx, appId, tableName, id, indexField)
        if err != nil {
            return id, nil    // the value for this field is empty. we don't need to do anything. Because people would not search on an empty value.
        }
        key := getIndexKey(appId, tableName, indexField, value)
        bz, err := store.Get([]byte(key))
        if err != nil{
            return id, err
        }
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

    store := DbChainStore(ctx, k.storeKey)
    start, end := getFieldDataIteratorStartAndEndKey(appId, tableName, fieldName)
    iter := store.Iterator([]byte(start), []byte(end))
    for ; iter.Valid(); iter.Next() {
        if iter.Error() != nil{
            return iter.Error()
        }
        dataKey := iter.Key()
        id := getIdFromDataKey(dataKey)
        if isRowFrozen(store, appId, tableName, id) {
            continue
        }

        val := iter.Value()
        k.cdc.MustUnmarshalBinaryBare(val, &dataValue)

        indexKey := getIndexKey(appId, tableName, fieldName, dataValue)
        bz, err := store.Get([]byte(indexKey))
        if err != nil{
            return err
        }
        if bz != nil {
            k.cdc.MustUnmarshalBinaryBare(bz, &indexValue)
        } else {
            indexValue = []string{}
        }
        indexValue = append(indexValue, fmt.Sprint(id))
        err = store.Set([]byte(indexKey), k.cdc.MustMarshalBinaryBare(indexValue))
        if err != nil{
            return err
        }
    }
    return nil
}

func dropIndexData(k Keeper, ctx sdk.Context, appId uint, tableName, fieldName string) error {
    store := DbChainStore(ctx, k.storeKey)
    start, end := getIndexDataIteratorStartAndEndKey(appId, tableName, fieldName)

    iter := store.Iterator([]byte(start), []byte(end))
    for ; iter.Valid(); iter.Next() {
        if iter.Error() != nil{
            return iter.Error()
        }
        indexKey := iter.Key()
        err := store.Delete([]byte(indexKey))
        if err != nil{
            return err
        }
    }
    return nil
}
