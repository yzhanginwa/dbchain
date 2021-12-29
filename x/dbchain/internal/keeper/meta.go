package keeper

import (
    "errors"
    "fmt"
    sdk "github.com/dbchaincloud/cosmos-sdk/types"
    "github.com/yzhanginwa/dbchain/x/dbchain/internal/keeper/cache"
    "github.com/yzhanginwa/dbchain/x/dbchain/internal/types"
    "github.com/yzhanginwa/dbchain/x/dbchain/internal/utils"
    "strconv"
    "strings"
)

/////////////////////////////
//                         //
// table related functions //
//                         //
/////////////////////////////

func (k Keeper) GetTables(ctx sdk.Context, appId uint) []string {
    store := DbChainStore(ctx, k.storeKey)
    tablesKey := getTablesKey(appId)
    bz, err := store.Get([]byte(tablesKey))
    if bz == nil || err != nil{
        return []string{}
    }
    var tableNames []string
    k.cdc.MustUnmarshalBinaryBare(bz, &tableNames)
    return tableNames
}


// Check if the table is present in the store or not
func (k Keeper) HasTable(ctx sdk.Context, appId uint, tableName string) bool {
    store := DbChainStore(ctx, k.storeKey)
    has ,err := store.Has([]byte(getTableKey(appId, tableName)))
    if err != nil{
        return false
    }
    return has
}

func (k Keeper) HasField(ctx sdk.Context, appId uint, tableName string, fieldName string) bool {
    table, err := k.GetTable(ctx, appId, tableName)
    if err != nil {
        return false
    }
    for _, f := range(table.Fields) {
        if f == fieldName {
            return true
        }
    }
    return false
}

// Create a new table
func (k Keeper) CreateTable(ctx sdk.Context, appId uint, owner sdk.AccAddress, tableName string, fieldNames []string) {
    store := DbChainStore(ctx, k.storeKey)
    table := types.NewTable()
    table.Owner = owner
    table.Name = tableName
    table.Fields = preProcessFields(fieldNames)
    // make Memos the same length as Fields
    fieldsLength := len(table.Fields)
    for len(table.Memos) < fieldsLength {
        table.Memos = append(table.Memos, "")
    }
    err := store.Set([]byte(getTableKey(appId, table.Name)), k.cdc.MustMarshalBinaryBare(table))
    if err != nil{
        return
    }

    var tables []string
    bz, err :=store.Get([]byte(getTablesKey(appId)))
    if err != nil{
        return
    }
    if bz == nil {
        tables = append(tables, table.Name)
    } else {
        k.cdc.MustUnmarshalBinaryBare(bz, &tables)
        tables = append(tables, table.Name)
    }
    store.Set([]byte(getTablesKey(appId)), k.cdc.MustMarshalBinaryBare(tables))
}


// Remove a table
func (k Keeper) DropTable(ctx sdk.Context, appId uint, owner sdk.AccAddress, tableName string) error {
    store := DbChainStore(ctx, k.storeKey)
    var tables []string
    bz, err :=store.Get([]byte(getTablesKey(appId)))
    if err != nil{
        return err
    }
    if bz != nil {
        k.cdc.MustUnmarshalBinaryBare(bz, &tables)
        for i, tbl := range tables {
            if tableName == tbl {
                //Drop column type and option first
                tableFields , _ := k.getTableFields(ctx, appId, tableName)
                //get counter cache field
                for _, field := range tableFields {
                    key := getColumnDataTypesKey(appId, tableName, field)
                    store.Delete([]byte(key))
                    key = getColumnOptionsKey(appId, tableName, field)
                    store.Delete([]byte(key))
                }

                //drop table
                tables = append(tables[:i], tables[i+1:]...)
                if len(tables) < 1 {
                    store.Delete([]byte(getTablesKey(appId)))
                } else {
                    store.Set([]byte(getTablesKey(appId)), k.cdc.MustMarshalBinaryBare(tables))
                }
                store.Delete([]byte(getTableKey(appId, tableName)))
                //delete rows of table
                start, end := getFieldDataIteratorStartAndEndKey(appId, tableName, "id")
                iter := store.Iterator([]byte(start), []byte(end))
                for ; iter.Valid(); iter.Next() {
                    if iter.Error() != nil{
                        return err
                    }
                    key := iter.Key()
                    id := getIdFromDataKey(key)
                    k.Delete(ctx, appId, tableName, id, owner)
                }
                //Drop index
                indexFields , err := k.GetIndexFields(ctx, appId, tableName)
                if err != nil {
                    break
                }
                for _, indexField := range indexFields {
                    k.DropIndex(ctx, appId, owner, tableName, indexField)
                }

                //Drop next_id key/value
                dropNextId(k, ctx, appId, tableName)
                //
                k.DeleteCounterCache(ctx, appId, tableName)

                break
            }
        }
    }
    cache.VoidTable(appId,tableName)
    return nil
}

func (k Keeper) DeleteCounterCache(ctx sdk.Context, appId uint, tableName string) bool {
    store := DbChainStore(ctx, k.storeKey)
    //1、as a main table
    counterCacheFields := k.GetCounterCacheFields(ctx, appId, tableName)
    if len(counterCacheFields) != 0 {
        key := getTableCounterCacheFieldKey(appId, tableName)
        err := store.Delete([]byte(key))
        if err != nil {
            return false
        }
    }
    for _, counterCacheField := range counterCacheFields {
        counterCaches := k.GetCounterCache(ctx, appId, counterCacheField.AssociationTable)
        exist := false
        for i, counterCache := range counterCaches {
            if counterCache.AssociationTable == tableName {
                exist = true
                counterCaches = append(counterCaches[: i], counterCaches[i+1 :]...)
            }
        }
        if !exist {
            continue
        }
        key := getTableCounterCacheInfoKey(appId, counterCacheField.AssociationTable)
        if len(counterCaches) == 0 {
            store.Delete([]byte(key))
        } else {
            bz := k.cdc.MustMarshalBinaryBare(counterCaches)
            err := store.Set([]byte(key), bz)
            if err != nil {
                return false
            }
        }
    }
    //2、as a satellite table
    counterCacheInfos := k.GetCounterCache(ctx, appId, tableName)
    if len(counterCacheInfos) != 0 {
        key := getTableCounterCacheInfoKey(appId, tableName)
        err := store.Delete([]byte(key))
        if err != nil {
            return false
        }
    }
    //
    for _ , counterCacheInfo := range counterCacheInfos {
        result, _ := k.DropColumn(ctx, appId, counterCacheInfo.AssociationTable, counterCacheInfo.CounterCacheField)
        if !result {
            return result
        }
        //delete association
        counterCacheFields := k.GetCounterCacheFields(ctx, appId, counterCacheInfo.AssociationTable)
        for i, counterCacheField := range counterCacheFields {
            if counterCacheField.AssociationTable == tableName {
                counterCacheFields = append(counterCacheFields[:i],counterCacheFields[i+1:]... )
                key := getTableCounterCacheFieldKey(appId, counterCacheInfo.AssociationTable)
                if len(counterCacheFields) != 0 {
                    store.Set([]byte(key), k.cdc.MustMarshalBinaryBare(counterCacheFields))
                } else {
                    store.Delete([]byte(key))
                }
                break
            }
        }
    }
    return true
}

func (k Keeper) DeleteCounterCacheField(ctx sdk.Context, appId uint, tableName,field string) bool {
    store := DbChainStore(ctx, k.storeKey)
    //1、as a main table
    counterCacheFields := k.GetCounterCacheFields(ctx, appId, tableName)
    var targetCounterCacheField  types.CounterCacheField
    for i, counterCacheField := range counterCacheFields {
        if counterCacheField.FieldName == field {
            targetCounterCacheField = counterCacheField
            counterCacheFields = append(counterCacheFields[:i], counterCacheFields[i+1:]...)
        }
    }
    if len(counterCacheFields) == 0 {
        key := getTableCounterCacheFieldKey(appId, tableName)
        err := store.Delete([]byte(key))
        if err != nil {
            return false
        }
    } else {
        bz := k.cdc.MustMarshalBinaryBare(counterCacheFields)
        key := getTableCounterCacheFieldKey(appId, tableName)
        err := store.Set([]byte(key),bz)
        if err != nil {
            return false
        }
    }

    //for _, counterCacheField := range counterCacheFields {
    counterCaches := k.GetCounterCache(ctx, appId, targetCounterCacheField.AssociationTable)
    exist := false
    for i, counterCache := range counterCaches {
        if counterCache.AssociationTable == tableName {
            exist = true
            counterCaches = append(counterCaches[: i], counterCaches[i+1 :]...)
        }
    }
    if exist {
        key := getTableCounterCacheInfoKey(appId, targetCounterCacheField.AssociationTable)
        if len(counterCaches) == 0 {
            store.Delete([]byte(key))
        } else {
            bz := k.cdc.MustMarshalBinaryBare(counterCaches)
            err := store.Set([]byte(key), bz)
            if err != nil {
                return false
            }
        }
    }

    //}
    //2、as a satellite table
    counterCacheInfos := k.GetCounterCache(ctx, appId, tableName)
    var targetCounterCacheInfo types.CounterCache
    for i, counterCacheInfo := range counterCacheInfos {
        if counterCacheInfo.ForeignKey == field {
            targetCounterCacheInfo = counterCacheInfo
            counterCacheInfos = append(counterCacheInfos[:i], counterCacheInfos[i+1:]...)
            break
        }
    }

    if len(counterCacheInfos) == 0 {
        key := getTableCounterCacheInfoKey(appId, tableName)
        err := store.Delete([]byte(key))
        if err != nil {
            return false
        }
    } else {
        bz := k.cdc.MustMarshalBinaryBare(counterCacheInfos)
        key := getTableCounterCacheInfoKey(appId, tableName)
        err := store.Set([]byte(key),bz)
        if err != nil {
            return false
        }
    }
    //delete association
    counterCacheFields = k.GetCounterCacheFields(ctx, appId, targetCounterCacheInfo.AssociationTable)
    for i, counterCacheField := range counterCacheFields {
        if counterCacheField.AssociationTable == tableName {
            counterCacheFields = append(counterCacheFields[:i],counterCacheFields[i+1:]... )
            key := getTableCounterCacheFieldKey(appId, targetCounterCacheInfo.AssociationTable)
            if len(counterCacheFields) != 0 {
                store.Set([]byte(key), k.cdc.MustMarshalBinaryBare(counterCacheFields))
            } else {
                store.Delete([]byte(key))
            }
            break
        }
    }
    return true
}

// Modify Table Association
func (k Keeper) ModifyTableAssociation(ctx sdk.Context, appId uint, tableName, option, associationMode, associationTable, foreignKey, method string) error {
    store := DbChainStore(ctx, k.storeKey)
    associations := make([]types.Association,0)
    key := getTableAssociationsKey(appId, tableName)

    bz , err := store.Get([]byte(key))
    if err != nil {
        return err
    }
    if bz != nil {
        k.cdc.MustUnmarshalBinaryBare(bz, &associations)
    }

    newAssociation := types.Association{
        AssociationMode: associationMode,
        AssociationTable: associationTable,
        ForeignKey: foreignKey,
        Method: method,
    }

    if option == "add" {
        for _, association := range associations {
            if association.Equal(newAssociation) {
                return errors.New("this association has been add")
            }
        }
        associations = append(associations, newAssociation)

    } else {
        hasAssociation := false
        for index, association := range associations {
            if association.Equal(newAssociation) {
                associations = append(associations[:index], associations[index+1:]...)
                hasAssociation = true
                break
            }
        }
        if !hasAssociation {
            return errors.New("this association does not exit")
        }
    }

    if len(associations) == 0 {
        store.Delete([]byte(key))
        return nil
    }

    bz = k.cdc.MustMarshalBinaryBare(associations)
    return store.Set([]byte(key), bz)
}

func (k Keeper)GetTableAssociations(ctx sdk.Context, appId uint, tableName string) []types.Association{
    store := DbChainStore(ctx, k.storeKey)
    associations := make([]types.Association,0)
    key := getTableAssociationsKey(appId, tableName)
    bz , err := store.Get([]byte(key))
    if err != nil {
        return nil
    }
    k.cdc.MustUnmarshalBinaryBare(bz, &associations)
    return associations
}

// Enable Table counter cache
func (k Keeper) AddCounterCache(ctx sdk.Context, appId uint, tableName, associationTable, foreignKey, counterCacheField, limit string) error {
    iLimit, err := strconv.Atoi(limit)
    if err != nil {
        return err
    }

    //save to associationTable
    newCounterCache := types.CounterCache{
        AssociationTable: tableName,
        ForeignKey: foreignKey,
        CounterCacheField: counterCacheField,
        Limit: limit,
    }

    store := DbChainStore(ctx, k.storeKey)
    counterCaches := make([]types.CounterCache, 0)
    key := getTableCounterCacheInfoKey(appId, associationTable)
    bz , err := store.Get([]byte(key))
    if err != nil {
        return err
    }

    if bz != nil {
        k.cdc.MustUnmarshalBinaryBare(bz, &counterCaches)
        for _, counterCache := range counterCaches {
            if counterCache.AssociationTable == tableName {
                return errors.New("this item has been added")
            }
        }
    }
    counterCaches = append(counterCaches, newCounterCache)
    bz = k.cdc.MustMarshalBinaryBare(counterCaches)
    err = store.Set([]byte(key), bz)
    if err != nil {
        return err
    }

    if !k.HasField(ctx, appId, associationTable, foreignKey) {
        return errors.New(fmt.Sprintf("Field %s of table %s not exists!", foreignKey, associationTable))
    }

    //add a new field for associationTable
    field := strings.ToLower(counterCacheField)
    if k.HasField(ctx, appId, tableName, field) {
        return errors.New(fmt.Sprintf("Field %s of table %s exists already!", counterCacheField, tableName))
    }
    _, err = k.AddColumn(ctx, appId, tableName, counterCacheField)
    if err != nil {
        return err
    }
    if !k.setCounterCacheColumnOption(ctx, appId, tableName, counterCacheField) {
        return errors.New("setCounterCacheColumnOption err")
    }


    //save counterCacheField
    existCounterCacheFields := make([]types.CounterCacheField, 0)
    associationTableCounterCacheFieldKey := getTableCounterCacheFieldKey(appId, tableName)
    bz, _ = store.Get([]byte(associationTableCounterCacheFieldKey))
    if bz != nil {
        k.cdc.MustUnmarshalBinaryBare(bz, &existCounterCacheFields)

        for _, existCounterCacheField := range existCounterCacheFields {
            if existCounterCacheField.FieldName == counterCacheField {
                return errors.New("this counterCacheField has been added")
            }
        }
    }

    newCounterCacheFields := types.CounterCacheField{
        AssociationTable: associationTable,
        FieldName: counterCacheField,
    }
    existCounterCacheFields = append(existCounterCacheFields, newCounterCacheFields)
    store.Set([]byte(associationTableCounterCacheFieldKey), k.cdc.MustMarshalBinaryBare(existCounterCacheFields))


    //update counter cache
    ids := k.findAllWithoutCheckPermission(ctx, appId, tableName)
    for _, id := range ids {
        sid := fmt.Sprintf("%d", id)
        counter := k.findByWithoutCheckPermission(ctx, appId, associationTable, foreignKey, []string{ sid })

        if iLimit > 0 && len(counter) > iLimit {
            return errors.New("counter number bigger than limit")
        }
        // set counter
        key := getDataKeyBytes(appId, tableName, counterCacheField, id)
        value := fmt.Sprintf("%d", len(counter))
        err := store.Set(key, k.cdc.MustMarshalBinaryBare(value))
        if err != nil {
            return err
        }

    }
    return nil
}

func (k Keeper)GetCounterCache(ctx sdk.Context, appId uint, tableName string) []types.CounterCache {
    store := DbChainStore(ctx, k.storeKey)
    counterCaches := make([]types.CounterCache, 0)
    key := getTableCounterCacheInfoKey(appId, tableName)
    bz , err := store.Get([]byte(key))
    if err != nil {
        return nil
    }
    if bz == nil{
        return nil
    }
    k.cdc.MustUnmarshalBinaryBare(bz, &counterCaches)
    return counterCaches
}

func (k Keeper)GetCounterCacheFields(ctx sdk.Context, appId uint, tableName string) []types.CounterCacheField {
    store := DbChainStore(ctx, k.storeKey)
    counterCacheFields := make([]types.CounterCacheField, 0)
    key := getTableCounterCacheFieldKey(appId, tableName)
    bz , err := store.Get([]byte(key))
    if err != nil {
        return nil
    }
    if bz == nil {
        return nil
    }
    k.cdc.MustUnmarshalBinaryBare(bz, &counterCacheFields)
    return counterCacheFields
}


func (k Keeper)GetTable(ctx sdk.Context, appId uint, tableName string) (types.Table, error){
    cTable, err := cache.GetTable(appId, tableName)
    if err == nil{
        return cTable, nil
    }
    table, err := k.RawGetTable(ctx, appId, tableName)
    if err != nil{
        return types.Table{}, err
    }
    cache.SetTable(appId, tableName, table)
    return table, nil
}

// Get a table 
func (k Keeper) RawGetTable(ctx sdk.Context, appId uint, tableName string) (types.Table, error) {
    store := DbChainStore(ctx, k.storeKey)
    bz, err := store.Get([]byte(getTableKey(appId, tableName)))
    if err != nil{
        return types.Table{}, err
    }
    if bz == nil {
        return types.Table{}, errors.New(fmt.Sprintf("table %s not found", tableName))
    }
    var table types.Table
    k.cdc.MustUnmarshalBinaryBare(bz, &table)
    return table, nil
}

// Add a field
func (k Keeper) AddColumn(ctx sdk.Context, appId uint, tableName string, fieldName string) (bool, error) {
    table, err := k.GetTable(ctx, appId, tableName)
    if err != nil {
        return false, err
    }

    for _, fld := range table.Fields {
        if fieldName == fld {
            return false, errors.New(fmt.Sprintf("field %s existed already", fieldName))
        }
    }

    table.Fields = append(table.Fields, fieldName)
    table.Memos = append(table.Memos, "")

    store := DbChainStore(ctx, k.storeKey)
    err = store.Set([]byte(getTableKey(appId, table.Name)), k.cdc.MustMarshalBinaryBare(table))
    if err != nil{
        return false, err
    }
    //void cache appTable
    cache.VoidTable(appId,table.Name)
    return true, nil
}

// Remove a field
func (k Keeper) DropColumn(ctx sdk.Context, appId uint, tableName string, fieldName string) (bool, error){
    table, err := k.GetTable(ctx, appId, tableName)
    if err != nil {
        return false, err
    }

    if isSystemField(fieldName) {
        return false, errors.New(fmt.Sprintf("cannot remove system fields"))
    }

    var foundField = false 
    for i, fld := range table.Fields {
        if fieldName == fld {
            foundField = true
            table.Fields = append(table.Fields[:i], table.Fields[i+1:]...)
            table.Memos  = append(table.Memos[:i],  table.Memos[i+1:]...)
            break
        }
    }
    if !foundField {
        return false, errors.New(fmt.Sprintf("field %s not existed", fieldName))
    }

    store := DbChainStore(ctx, k.storeKey)
    err = store.Set([]byte(getTableKey(appId, table.Name)), k.cdc.MustMarshalBinaryBare(table))
    if err != nil{
        return false, err
    }

    //void cache appTable
    cache.VoidTable(appId,table.Name)
    // Remove data of this dropped column
    removeDataOfColumn(k, ctx, appId, tableName, fieldName)
    // delete counterCache
    delStatus := k.DeleteCounterCacheField(ctx, appId, tableName,fieldName)
    if !delStatus {
        return false, errors.New("del counter cache field fail")
    }
    return true, nil
}

// Rename a field
func (k Keeper) RenameColumn(ctx sdk.Context, appId uint, tableName string, oldField string, newField string) (bool, error) {
    table, err := k.GetTable(ctx, appId, tableName)
    if err != nil {
        return false, err
    }

    if oldField == "" || newField == "" {
        return false, errors.New(fmt.Sprintf("cannot have empty field name"))
    }

    oldField = strings.ToLower(oldField)
    newField = strings.ToLower(newField)

    if oldField == "id" || newField == "id" {
        return false, errors.New(fmt.Sprintf("cannot rename field id"))
    }


    var foundField = false
    var index = 0
    for i, fld := range table.Fields {
        if oldField == fld {
            foundField = true
            index = i
            break
        }
        if newField == fld {
            return false, errors.New(fmt.Sprintf("cannot rename to field %s", newField))
        }
    }
    if !foundField {
        return false, errors.New(fmt.Sprintf("field %s not found", oldField))
    }

    table.Fields[index] = newField

    store := DbChainStore(ctx, k.storeKey)
    err = store.Set([]byte(getTableKey(appId, table.Name)), k.cdc.MustMarshalBinaryBare(table))
    if err != nil{
        return false, err
    }
    //void cache appTable
    cache.VoidTable(appId,table.Name)
    return true, nil
}

func (k Keeper) ModifyOption(ctx sdk.Context, appId uint, owner sdk.AccAddress, tableName string, action string, option string) {
    store := DbChainStore(ctx, k.storeKey)
    key := getTableOptionsKey(appId, tableName)
    var options []string
    var result []string

    bz, err := store.Get([]byte(key))
    if err != nil{
        return
    }
    if bz != nil {
        k.cdc.MustUnmarshalBinaryBare(bz, &options)
    }

    optionExisted := utils.ItemExists(options, option)
    if action == "add" {
        if optionExisted {
            return
        } else {
            result = append(options, option)
        }
    } else {
        if optionExisted {
            for _, opt := range options {
                if opt == option {
                    continue
                }
                result = append(result, opt)
            }
        } else {
            return
        }
    }
    if len(result) > 0 {
        store.Set([]byte(key), k.cdc.MustMarshalBinaryBare(result))
    } else {
        store.Delete([]byte(key))
    }
}

func (k Keeper) AddInsertFilter(ctx sdk.Context, appId uint, owner sdk.AccAddress, tableName string, filter string) bool {
    if err := checkLuaSyntax(filter); err != nil {
        return false
    } 

    table, err := k.GetTable(ctx, appId, tableName)
    if err != nil {
        return false
    }

    if len(table.Filter) > 0 {
        return false
    }

    table.Filter = filter

    store := DbChainStore(ctx, k.storeKey)
    err = store.Set([]byte(getTableKey(appId, table.Name)), k.cdc.MustMarshalBinaryBare(table))
    if err != nil{
        return false
    }
    //void cache appTable
    cache.VoidTable(appId,table.Name)
    return true
}

func (k Keeper) DropInsertFilter(ctx sdk.Context, appId uint, owner sdk.AccAddress, tableName string) bool {
    table, err := k.GetTable(ctx, appId, tableName)
    if err != nil {
        return false
    }

    table.Filter = ""
    store := DbChainStore(ctx, k.storeKey)
    err = store.Set([]byte(getTableKey(appId, table.Name)), k.cdc.MustMarshalBinaryBare(table))
    if err != nil{
        return false
    }
    //void cache appTable
    cache.VoidTable(appId,table.Name)
    return true
}

func (k Keeper) GetInsertFilter(ctx sdk.Context, appId uint, tableName string) string {
    table, err := k.GetTable(ctx, appId, tableName)
    if err != nil {
        return ""
    }

    return table.Filter
}

func (k Keeper) AddTrigger(ctx sdk.Context, appId uint, owner sdk.AccAddress, tableName string, trigger string) bool {
    if err := checkLuaSyntax(trigger); err != nil {
        return false
    }

    table, err := k.GetTable(ctx, appId, tableName)
    if err != nil {
        return false
    }

    if len(table.Trigger) > 0 {
        return false
    }

    table.Trigger = trigger

    store := DbChainStore(ctx, k.storeKey)
    err = store.Set([]byte(getTableKey(appId, table.Name)), k.cdc.MustMarshalBinaryBare(table))
    if err != nil{
        return false
    }
    //void cache appTable
    cache.VoidTable(appId,table.Name)
    return true
}

func (k Keeper) DropTrigger(ctx sdk.Context, appId uint, owner sdk.AccAddress, tableName string) bool {
    table, err := k.GetTable(ctx, appId, tableName)
    if err != nil {
        return false
    }

    table.Trigger = ""
    store := DbChainStore(ctx, k.storeKey)
    err = store.Set([]byte(getTableKey(appId, table.Name)), k.cdc.MustMarshalBinaryBare(table))
    if err != nil{
        return false
    }
    //void cache appTable
    cache.VoidTable(appId,table.Name)
    return true
}

func (k Keeper) SetTableMemo(ctx sdk.Context, appId uint, tableName, memo string, owner sdk.AccAddress) bool {
    table, err := k.GetTable(ctx, appId, tableName)
    if err != nil {
        return false
    }

    table.Memo = memo
    store := DbChainStore(ctx, k.storeKey)
    err = store.Set([]byte(getTableKey(appId, table.Name)), k.cdc.MustMarshalBinaryBare(table))
    if err != nil{
        return false
    }
    //void cache appTable
    cache.VoidTable(appId,table.Name)
    return true
}

func (k Keeper) GetOption(ctx sdk.Context, appId uint, tableName string) ([]string, error) {
    store := DbChainStore(ctx, k.storeKey)
    key := getTableOptionsKey(appId, tableName)
    bz, err := store.Get([]byte(key))
    if err != nil{
        return nil, err
    }
    if bz == nil {
        return []string{}, nil
    }
    var options []string
    k.cdc.MustUnmarshalBinaryBare(bz, &options)
    return options, nil
}

func (k Keeper) GetWritableByGroups(ctx sdk.Context, appId uint, tableName string) []string {
    options, _ := k.GetOption(ctx, appId, tableName)
    baseLen := len(types.TBLOPT_WRITABLE_BY)
    var result []string

    for _, option := range options {
        if len(option) > (baseLen + 2) {
            if option[:baseLen] == string(types.TBLOPT_WRITABLE_BY) {
                g := option[baseLen:]
                gl := len(g)
                if g[0] == '(' && g[gl-1] == ')' {
                    result = append(result, g[1:gl-1])
                }
            }
        }
    }
    return result
}

func (k Keeper) ModifyColumnOption(ctx sdk.Context, appId uint, owner sdk.AccAddress, tableName string, fieldName string, action string, option string) bool {
    if !validateColumnOption(option) {
        return false
    }

    store := DbChainStore(ctx, k.storeKey)
    key := getColumnOptionsKey(appId, tableName, fieldName)
    var options []string
    var result []string

    bz, err := store.Get([]byte(key))
    if err != nil{
        return false
    }
    if bz != nil {
        k.cdc.MustUnmarshalBinaryBare(bz, &options)
    }

    optionExisted := isColumnOptionIncluded(options, option)

    if action == "add" {
        if optionExisted {
            return false
        } else {
            switch types.FieldOption(option) {
            case types.FLDOPT_NOTNULL:
                if !k.validateNotNullField(ctx, appId, tableName, fieldName) {
                    return false
                }
            case types.FLDOPT_UNIQUE:
                if !isColumnValuesUnique(k, ctx, appId, tableName, fieldName) {
                    return false
                }
            case types.FLDOPT_OWN:
                if !k.validateOwnField(ctx, appId, tableName, fieldName) {
                    return false
                }
            case types.FLDOPT_COUNTER_CACHE:
                return false

            }
            result = append(options, option)
        }
    } else {
        if optionExisted {
            for _, opt := range options {
                if opt == option {
                    continue
                }
                result = append(result, opt)
            }
        } else {
            return false
        }
    }

    if len(result) > 0 {
        err = store.Set([]byte(key), k.cdc.MustMarshalBinaryBare(result))
    } else {
        err = store.Delete([]byte(key))
    }
    if err != nil{
        return false
    }
    return true
}

func (k Keeper) setCounterCacheColumnOption (ctx sdk.Context, appId uint,tableName string, fieldName string) bool {
    store := DbChainStore(ctx, k.storeKey)
    result := []string{ string(types.FLDOPT_COUNTER_CACHE) }
    key := getColumnOptionsKey(appId, tableName, fieldName)
    err := store.Set([]byte(key), k.cdc.MustMarshalBinaryBare(result))
    if err != nil {
        return false
    }
    return true
}

func (k Keeper) SetColumnDataType(ctx sdk.Context, appId uint, owner sdk.AccAddress, tableName string, fieldName string, dataType string) bool {
    if !validateColumnDataType(dataType) {
        return false
    }

    store := DbChainStore(ctx, k.storeKey)
    key := getColumnDataTypesKey(appId, tableName, fieldName)
    var currentDataType string

    bz, err := store.Get([]byte(key))
    if err != nil{
        return false
    }
    if bz != nil {
        k.cdc.MustUnmarshalBinaryBare(bz, &currentDataType)
    }

    if currentDataType == dataType {
        return false
    }

    switch types.FieldDataType(dataType) {
    case types.FLDTYP_INT:
        if !k.validateIntField(ctx, appId, tableName, fieldName) {
            return false
        }
    case types.FLDTYP_FILE:
        if !k.validateFileField(ctx, appId, tableName, fieldName) {
            return false
        }
    case types.FLDTYP_DECIMAL:
        if !k.validateDecimalField(ctx, appId, tableName, fieldName) {
            return false
        }
    case types.FLDTYP_ADDRESS:
        if !k.validateAddressField(ctx, appId, tableName, fieldName) {
            return false
        }
    case types.FLDTYP_TIME:
        if !k.validateTimeField(ctx, appId, tableName, fieldName) {
            return false
        }
    }

    err = store.Set([]byte(key), k.cdc.MustMarshalBinaryBare(dataType))
    if err != nil{
        return false
    }
    return true
}

func (k Keeper) SetColumnMemo(ctx sdk.Context, appId uint, owner sdk.AccAddress, tableName string, fieldName string, memo string) bool {
    table, err := k.GetTable(ctx, appId, tableName)
    if err != nil {
        return false
    }
    for i, f := range(table.Fields) {
        if f == fieldName {
            fieldsLength := len(table.Fields)
            for len(table.Memos) < fieldsLength {
                table.Memos = append(table.Memos, "")
            }
            table.Memos[i] = memo
            store := DbChainStore(ctx, k.storeKey)
            err := store.Set([]byte(getTableKey(appId, table.Name)), k.cdc.MustMarshalBinaryBare(table))
            if err != nil{
                return false
            }
            //void cache appTable
            cache.VoidTable(appId,table.Name)
            return true
        }
    }
    return false
}

func (k Keeper) GetColumnOption(ctx sdk.Context, appId uint, tableName string, fieldName string) ([]string, error) {
    store := DbChainStore(ctx, k.storeKey)
    key := getColumnOptionsKey(appId, tableName, fieldName)
    bz, err := store.Get([]byte(key))
    if err != nil{
        return nil, err
    }
    if bz == nil {
        return []string{}, nil
    }
    var options []string
    k.cdc.MustUnmarshalBinaryBare(bz, &options)
    return options, nil
}

func (k Keeper) GetCounterInfo(ctx sdk.Context, appId uint, tableName string) ([]map[string]string, error) {

    counterCacheFields := k.GetCounterCacheFields(ctx, appId, tableName)
    if len(counterCacheFields) == 0 {
        return nil, nil
    }
    result := make([]map[string]string, 0)
    for _, counterCacheField := range counterCacheFields {

        associationTable := counterCacheField.AssociationTable
        counterCacheInfos := k.GetCounterCache(ctx, appId, associationTable)
        for _, counterCacheInfo := range counterCacheInfos {
            if counterCacheInfo.AssociationTable == tableName {
                temp :=  map[string]string{
                    "counter_field" : counterCacheField.FieldName,
                    "association_table" : associationTable,
                    "foreign_key" : counterCacheInfo.ForeignKey,
                    "limit" : counterCacheInfo.Limit,
                }
                result = append(result, temp)
                break
            }
        }
    }

    return result, nil
}

func (k Keeper) GetColumnDataType(ctx sdk.Context, appId uint, tableName string, fieldName string) (string, error) {
    store := DbChainStore(ctx, k.storeKey)
    key := getColumnDataTypesKey(appId, tableName, fieldName)
    bz, err := store.Get([]byte(key))
    if err != nil{
        return "", err
    }
    if bz == nil {
        return "string", nil
    }
    var dataType string
    k.cdc.MustUnmarshalBinaryBare(bz, &dataType)
    return dataType, nil
}

func (k Keeper) GetCanAddColumnOption(ctx sdk.Context, appId uint, tableName, fieldName, option string) bool {
    switch types.FieldOption(option) {
    case types.FLDOPT_NOTNULL:
        if !k.validateNotNullField(ctx, appId, tableName, fieldName) {
            return false
        }
    case types.FLDOPT_UNIQUE:
        if !isColumnValuesUnique(k, ctx, appId, tableName, fieldName) {
            return false
        }
    case types.FLDOPT_OWN:
        if !k.validateOwnField(ctx, appId, tableName, fieldName) {
            return false
        }
    case types.FLDOPT_READABLE:
        return true
    }
    return true
}

func (k Keeper) GetCanSetColumnDataType(ctx sdk.Context, appId uint, tableName, fieldName, dataType string) bool {
    switch types.FieldDataType(dataType) {
        case types.FLDTYP_STRING:
            return true
        case types.FLDTYP_INT:
           if !k.validateIntField(ctx, appId, tableName, fieldName) {
               return false
           }
        case types.FLDTYP_FILE:
           if !k.validateFileField(ctx, appId, tableName, fieldName) {
               return false
           }
        case types.FLDTYP_DECIMAL:
            if !k.validateDecimalField(ctx, appId, tableName, fieldName) {
                return false
            }
        case types.FLDTYP_ADDRESS:
            if !k.validateAddressField(ctx, appId, tableName, fieldName) {
                return false
            }
        case types.FLDTYP_TIME:
            if !k.validateTimeField(ctx, appId, tableName, fieldName) {
                return false
            }
        default:
            return false
    }
    return true
}

////////////////////
//                //
// helper methods //
//                //
////////////////////

// to preprocess the new table field names
// to make sure the fields be lowercase
// to make sure field id be in place
func preProcessFields(fieldNames []string) []string {
    var result = []string{"id", "created_by", "created_at", "tx_hash"}
    m := map[string]bool{
        "id" : true,
        "created_by" : true,
        "created_at" : true,
        "tx_hash" : true,
    }

    for _, field := range fieldNames {
        newName := strings.ToLower(field)
        if newName == "" {
            continue
        }
        if m[newName] {
            continue
        } else {
           m[newName] = true
           result = append(result, newName)
        }
    }
    return result
}

func validateColumnOption(option string) bool {
    switch types.FieldOption(option) {
    case types.FLDOPT_NOTNULL:
        return true
    case types.FLDOPT_UNIQUE:
        return true
    case types.FLDOPT_OWN:
        return true
    case types.FLDOPT_READABLE:
        return true
    }

    if types.ValidateEnumColumnOption(option) {
        return true
    }

    return false
}

func validateColumnDataType(dataType string) bool {
    switch types.FieldDataType(dataType) {
    case types.FLDTYP_STRING, types.FLDTYP_INT, types.FLDTYP_FILE, types.FLDTYP_DECIMAL, types.FLDTYP_ADDRESS,
    types.FLDTYP_TIME:
        return true
    }

    return false
}

func isColumnOptionIncluded(options []string, option string) bool {
    for _, opt := range options {
        if opt == option {
            return true
        }
        if types.ValidateEnumColumnOption(opt) {
            if types.ValidateEnumColumnOption(option) {
                return true
            }
        }
    }
    return false
}


func isColumnValuesUnique(k Keeper, ctx sdk.Context, appId uint, tableName string, fieldName string) bool {
    store := DbChainStore(ctx, k.storeKey)
    flag := make(map[string]bool)

    start, end := getFieldDataIteratorStartAndEndKey(appId, tableName, fieldName)
    iter := store.Iterator([]byte(start), []byte(end))
    var mold string
    for ; iter.Valid(); iter.Next() {
        if iter.Error() != nil{
            return false
        }
        dataKey := iter.Key()
        id := getIdFromDataKey(dataKey)
        if isRowFrozen(store, appId, tableName, id) {
            continue
        }

        val := iter.Value()
        k.cdc.MustUnmarshalBinaryBare(val, &mold)
        if _, ok := flag[mold]; ok {
            return false
        }
        flag[mold] = true
    }
    return true
}

func removeDataOfColumn(k Keeper, ctx sdk.Context, appId uint, tableName, fieldName string) {
    store := DbChainStore(ctx, k.storeKey)

    start, end := getFieldDataIteratorStartAndEndKey(appId, tableName, fieldName)
    iter := store.Iterator([]byte(start), []byte(end))
    for ; iter.Valid(); iter.Next() {
        if iter.Error() != nil{
            return
        }
        key := iter.Key()
        err := store.Delete([]byte(key))
        if err != nil{
            return
        }
    }
}

func isRowFrozen(store *SafeStore, appId uint, tableName string, id uint) bool {
    key := getDataKeyBytes(appId, tableName, types.FLD_FROZEN_AT, id)
    bz,_ := store.Get(key)
    if bz != nil{
        return true
    }
    return false
}
