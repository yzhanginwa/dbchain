package keeper

import (
    "fmt"
    "regexp"
    "strconv"
    "errors"
    sdk "github.com/cosmos/cosmos-sdk/types"
    "github.com/yzhanginwa/dbchain/x/dbchain/internal/types"
    "github.com/yzhanginwa/dbchain/x/dbchain/internal/utils"
    "strings"
)


func (k Keeper) DoFind(ctx sdk.Context, appId uint, tableName string, id uint) (types.RowFields, error){
    store := DbChainStore(ctx, k.storeKey)

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
        key := getDataKeyBytes(appId, tableName, fieldName, id)
        bz, err := store.Get(key)
        if err != nil{
            return nil,err
        }
        if bz != nil {
            k.cdc.MustUnmarshalBinaryBare(bz, &value)
            fields[fieldName] = value
        }
    }

    return fields, nil
}

func (k Keeper) FindField(ctx sdk.Context, appId uint, tableName string, id uint, fieldName string) (string, error){
    store := DbChainStore(ctx, k.storeKey)

    if !k.HasField(ctx, appId, tableName, fieldName) {
        return "", errors.New("Field not existed")
    }

    if id == 0 {
        return "", errors.New("Id cannot be 0")
    }

    key := getDataKeyBytes(appId, tableName, fieldName, id)
    bz, err := store.Get(key)
    if err != nil{
        return "", err
    }
    if bz != nil {
        var value string
        k.cdc.MustUnmarshalBinaryBare(bz, &value)
        return value, nil
    }
    return "", errors.New("Field not found")
}

func (k Keeper) Find(ctx sdk.Context, appId uint, tableName string, id uint, user sdk.AccAddress) (types.RowFields, error){
    if !k.isOwnId(ctx,appId, tableName, id, user) {
        return nil, errors.New(fmt.Sprintf("Failed to get fields for id %d", id))
    }

    return k.DoFind(ctx, appId, tableName, id)
}

// Find by an attribute in the r.Fields
func (k Keeper) FindBy(ctx sdk.Context, appId uint, tableName string, field string,  values []string, user sdk.AccAddress) []uint {
    store := DbChainStore(ctx, k.storeKey)

    hasIndex := false
    indexFields, err := k.GetIndexFields(ctx, appId, tableName)
    if err == nil {
        for _, item := range(indexFields) {
            if item == field {
                hasIndex = true
                break
            }
        }
    }

    results := []uint{}
    if hasIndex {
        for i := 0; i < len(values); i++ {
            value := values[i]
            key := getIndexKey(appId, tableName, field, value)
            bz, err := store.Get([]byte(key))
            if err != nil{
                return nil
            }
            var result []string
            if bz != nil {
                k.cdc.MustUnmarshalBinaryBare(bz, &result)
                for _, sId := range result {
                   id , err := strconv.ParseUint(sId, 10, 32)
                   if err != nil {
                       continue
                   }
                   results = append(results, uint(id))
                }
            }
        }
    } else {
        // partial table scanning
        start, end := getFieldDataIteratorStartAndEndKey(appId, tableName, field)
        iter := store.Iterator([]byte(start), []byte(end))
        var mold string
        for ; iter.Valid(); iter.Next() {
            if iter.Error() != nil{
                return nil
            }
            key := iter.Key()
            val := iter.Value()
            k.cdc.MustUnmarshalBinaryBare(val, &mold)
            if utils.StringIncluded(values, mold) {
                id := getIdFromDataKey(key)
                if isRowFrozen(store, appId, tableName, id) {
                    continue;
                }
                results = append(results, id)
            }
        }
    }

    // if public table or auditor user, return all ids
    if k.isTablePublic(ctx, appId, tableName) || k.isAuditor(ctx, appId, user) {
        return results
    } else {
        return k.filterOwnIds(ctx, appId, tableName, results, user)
    }
}

func (k Keeper) Where(ctx sdk.Context, appId uint, tableName string, field string, operator string, value string, reg *regexp.Regexp, user sdk.AccAddress) []uint {
    //TODO: consider if the field has index and how to make use of it
    store := DbChainStore(ctx, k.storeKey)
    isInteger := k.isTypeOfInteger(ctx, appId, tableName, field)
    results := []uint{}
    if field == "id" && (operator ==  "==" || operator ==  "=") {
        id , err := strconv.ParseUint(value, 10, 32)
        if err != nil || !k.isOwnId(ctx, appId, tableName, uint(id), user) {
            return results
        }
        results = append(results, uint(id))
        return results
    }


    if k.isIndexField(ctx, appId, tableName, field) {
        if operator ==  "==" || operator ==  "="  {
            var mold []string
            key := getIndexKey(appId, tableName, field, value)
            bz, err := store.Get([]byte(key))
            if err != nil{
                return results
            }
            if bz != nil {
                k.cdc.MustUnmarshalBinaryBare(bz, &mold)
            }
            for _, sId := range mold {
                id , err := strconv.ParseUint(sId, 10, 32)
                if err != nil {
                    continue
                }
                results = append(results, uint(id))
            }
            return results
        } else {
            start, end := getIndexDataIteratorStartAndEndKey(appId, tableName, field)
            iter := store.Iterator([]byte(start), []byte(end))
            for ; iter.Valid(); iter.Next() {
                if iter.Error() != nil {
                    continue
                }
                key := iter.Key()
                val := iter.Value()
                sliceKey := strings.Split(string(key),":")
                matching := fieldValueCompare(isInteger, operator, sliceKey[len(sliceKey)-1], value, reg)
                if matching {
                    var mold []string
                    k.cdc.MustUnmarshalBinaryBare(val, &mold)
                    for _, sId := range mold {
                        id , err := strconv.ParseUint(sId, 10, 32)
                        if err != nil {
                            continue
                        }
                        results = append(results, uint(id))
                    }
                }
            }
            return results
        }
    }

    start, end := getFieldDataIteratorStartAndEndKey(appId, tableName, field)
    iter := store.Iterator([]byte(start), []byte(end))
    var mold string
    for ; iter.Valid(); iter.Next() {
        if iter.Error() != nil{
            return nil
        }
        key := iter.Key()
        val := iter.Value()
        k.cdc.MustUnmarshalBinaryBare(val, &mold)

        matching := fieldValueCompare(isInteger, operator, mold, value, reg)
        if matching {
            id := getIdFromDataKey(key)
            if isRowFrozen(store, appId, tableName, id) {
                continue;
            }
            results = append(results, id)
        }
    }

    if k.isTablePublic(ctx, appId, tableName) || k.isAuditor(ctx, appId, user) {
        return results
    } else {
        return k.filterOwnIds(ctx, appId, tableName, results, user)
    }
}

func (k Keeper) FindAll(ctx sdk.Context, appId uint, tableName string, user sdk.AccAddress) []uint {
    store := DbChainStore(ctx, k.storeKey)
    var result []uint

    // full table scanning
    start, end := getFieldDataIteratorStartAndEndKey(appId, tableName, "id")
    iter := store.Iterator([]byte(start), []byte(end))
    for ; iter.Valid(); iter.Next() {
        if iter.Error() != nil{
            return nil
        }
        key := iter.Key()
        id := getIdFromDataKey(key)
        if isRowFrozen(store, appId, tableName, id) {
            continue;
        }
        result = append(result, id)
    }

    if k.isTablePublic(ctx, appId, tableName) || k.isAuditor(ctx, appId, user) {
        return result
    } else {
        return k.filterOwnIds(ctx, appId, tableName, result, user)
    }
}

//////////////////
//              //
// helper funcs //
//              //
//////////////////
func (k Keeper) isIndexField(ctx sdk.Context, appId uint, tableName, field string) bool {
    indexFields, err := k.GetIndexFields(ctx, appId, tableName)
    if err != nil {
        return false
    }
    for _,indexField := range indexFields {
        if field == indexField {
            return true
        }
    }
    return false
}

func (k Keeper) isTablePublic(ctx sdk.Context, appId uint, tableName string) bool {
    tableOptions, _ := k.GetOption(ctx, appId, tableName)
    return utils.ItemExists(tableOptions, string(types.TBLOPT_PUBLIC))
}

func (k Keeper) isOwnId(ctx sdk.Context, appId uint, tableName string, id uint, user sdk.AccAddress) bool{
    var ids []uint
    ids = append(ids, id)

    // if public table, return all ids
    if !k.isTablePublic(ctx, appId, tableName) && !k.isAuditor(ctx, appId, user) {
        ids = k.filterOwnIds(ctx, appId, tableName, ids, user)
        if len(ids) < 1 {
            return false
        }
    }
    return true
}

func (k Keeper) filterOwnIds(ctx sdk.Context, appId uint,  tableName string, ids []uint, user sdk.AccAddress) []uint {
    store := DbChainStore(ctx, k.storeKey)
    var userString string = user.String()

    var result = []uint{}
    var mold string
    for _, id := range ids {
        key := getDataKeyBytes(appId, tableName, "created_by", uint(id))
        bz, err := store.Get(key)
        if err != nil{
            return nil
        }
        if bz != nil {
            k.cdc.MustUnmarshalBinaryBare(bz, &mold)
            if mold == userString {
                result = append(result, uint(id))
            }
        }
    }
    return result
}

func (k Keeper) isTypeOfInteger(ctx sdk.Context, appId uint, tableName, fieldName string) bool {
    fieldDataType, _ := k.GetColumnDataType(ctx, appId, tableName, fieldName)
    return utils.StringIncluded(fieldDataType, string(types.FLDTYP_INT))
}

func fieldValueCompare(isInteger bool, operator, left, right string, reg *regexp.Regexp) bool {
    matching := false

    if isInteger {
        l, err := strconv.Atoi(left)
        if err != nil {
            return false
        }
        r, err := strconv.Atoi(right)
        if err != nil {
            return false
        }

        switch operator {
        case "=", "==" :
            if l == r{
                matching = true
            }
        case ">" :
            if l > r {
                matching = true
            }
        case ">=" :
            if l >= r {
                matching = true
            }
        case "<" :
            if l < r {
                matching = true
            }
        case "<=" :
            if l <= r {
                matching = true
            }
        case "<>" :
            if l != r {
                matching = true
            }
        }

    } else {
        switch operator {
        case "=", "==" :
            if left == right {
                matching = true
            }
        case ">" :
            if left > right {
                matching = true
            }
        case ">=" :
            if left >= right {
                matching = true
            }
        case "<" :
            if left < right {
                matching = true
            }
        case "<=" :
            if left <= right {
                matching = true
            }
        case "<>" :
            if left != right {
                matching = true
            }
        case "like":
            if reg.MatchString(left) {
                matching = true
            }

        }
    }
    return matching
}
