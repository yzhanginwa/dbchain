package keeper

import (
    "fmt"
    "errors"
    "strings"
    sdk "github.com/cosmos/cosmos-sdk/types"
    "github.com/yzhanginwa/dbchain/x/dbchain/internal/types"
    "github.com/yzhanginwa/dbchain/x/dbchain/internal/utils"
    ss "github.com/yzhanginwa/dbchain/x/dbchain/internal/super_script"
)

/////////////////////////////
//                         //
// table related functions //
//                         //
/////////////////////////////

func (k Keeper) GetTables(ctx sdk.Context, appId uint) []string {
    store := ctx.KVStore(k.storeKey)
    tablesKey := getTablesKey(appId)
    bz := store.Get([]byte(tablesKey))
    if bz == nil {
        return []string{}
    }
    var tableNames []string
    k.cdc.MustUnmarshalBinaryBare(bz, &tableNames)
    return tableNames
}


// Check if the table is present in the store or not
func (k Keeper) HasTable(ctx sdk.Context, appId uint, name string) bool {
    store := ctx.KVStore(k.storeKey)
    return store.Has([]byte(getTableKey(appId, name)))
}

func (k Keeper) HasField(ctx sdk.Context, appId uint, tableName string, field string) bool {
    table, err := k.GetTable(ctx, appId, tableName)
    if err != nil {
        return false
    }
    for _, f := range(table.Fields) {
        if f == field {
            return true
        }
    }
    return false
}

// Create a new table
func (k Keeper) CreateTable(ctx sdk.Context, appId uint, owner sdk.AccAddress, name string, fields []string) {
    store := ctx.KVStore(k.storeKey)
    table := types.NewTable()
    table.Owner = owner
    table.Name = name
    table.Fields = preProcessFields(fields)
    // make Memos the same length as Fields
    fieldsLength := len(table.Fields)
    for len(table.Memos) < fieldsLength {
        table.Memos = append(table.Memos, "")
    }
    store.Set([]byte(getTableKey(appId, table.Name)), k.cdc.MustMarshalBinaryBare(table))

    var tables []string
    bz :=store.Get([]byte(getTablesKey(appId)))
    if bz == nil {
        tables = append(tables, table.Name)
    } else {
        k.cdc.MustUnmarshalBinaryBare(bz, &tables)
        tables = append(tables, table.Name)
    }
    store.Set([]byte(getTablesKey(appId)), k.cdc.MustMarshalBinaryBare(tables))
}

// Remove a table
func (k Keeper) DropTable(ctx sdk.Context, appId uint, owner sdk.AccAddress, name string) {
    store := ctx.KVStore(k.storeKey)
    var tables []string
    bz :=store.Get([]byte(getTablesKey(appId)))
    if bz != nil {
        k.cdc.MustUnmarshalBinaryBare(bz, &tables)
        for i, tbl := range tables {
            if name == tbl {
                tables = append(tables[:i], tables[i+1:]...)
                if len(tables) < 1 {
                    store.Delete([]byte(getTablesKey(appId)))
                } else {
                    store.Set([]byte(getTablesKey(appId)), k.cdc.MustMarshalBinaryBare(tables))
                }
                store.Delete([]byte(getTableKey(appId, name)))
                break
            }
        }
    }
}

// Get a table 
func (k Keeper) GetTable(ctx sdk.Context, appId uint, name string) (types.Table, error) {
    store := ctx.KVStore(k.storeKey)
    bz := store.Get([]byte(getTableKey(appId, name)))
    if bz == nil {
        return types.Table{}, errors.New(fmt.Sprintf("table %s not found", name))
    }
    var table types.Table
    k.cdc.MustUnmarshalBinaryBare(bz, &table)
    return table, nil
}

// Add a field
func (k Keeper) AddColumn(ctx sdk.Context, appId uint, name string, field string) (bool, error) {
    table, err := k.GetTable(ctx, appId, name)
    if err != nil {
        return false, err
    }

    for _, fld := range table.Fields {
        if field == fld {
            return false, errors.New(fmt.Sprintf("field %s existed already", field))
        }
    }

    table.Fields = append(table.Fields, field)
    table.Memos = append(table.Memos, "")

    store := ctx.KVStore(k.storeKey)
    store.Set([]byte(getTableKey(appId, table.Name)), k.cdc.MustMarshalBinaryBare(table))
    return true, nil
}

// Remove a field
func (k Keeper) DropColumn(ctx sdk.Context, appId uint, name string, field string) (bool, error){
    table, err := k.GetTable(ctx, appId, name)
    if err != nil {
        return false, err
    }

    if isSystemField(field) {
        return false, errors.New(fmt.Sprintf("cannot remove field id"))
    }

    var foundField = false 
    for i, fld := range table.Fields {
        if field == fld {
            foundField = true
            table.Fields = append(table.Fields[:i], table.Fields[i+1:]...)
            table.Memos  = append(table.Memos[:i],  table.Memos[i+1:]...)
            break
        }
    }
    if !foundField {
        return false, errors.New(fmt.Sprintf("field %s not existed", field))
    }

    store := ctx.KVStore(k.storeKey)
    store.Set([]byte(getTableKey(appId, table.Name)), k.cdc.MustMarshalBinaryBare(table))
    return true, nil
}

// Rename a field
func (k Keeper) RenameColumn(ctx sdk.Context, appId uint, name string, oldField string, newField string) (bool, error) {
    table, err := k.GetTable(ctx, appId, name)
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

    store := ctx.KVStore(k.storeKey)
    store.Set([]byte(getTableKey(appId, table.Name)), k.cdc.MustMarshalBinaryBare(table))
    return true, nil
}

func (k Keeper) ModifyOption(ctx sdk.Context, appId uint, owner sdk.AccAddress, tableName string, action string, option string) {
    store := ctx.KVStore(k.storeKey)
    key := getTableOptionsKey(appId, tableName)
    var options []string
    var result []string

    bz := store.Get([]byte(key))
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
    if !validateInsertFilterSyntax(k, ctx, appId, tableName, filter) {
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

    store := ctx.KVStore(k.storeKey)
    store.Set([]byte(getTableKey(appId, table.Name)), k.cdc.MustMarshalBinaryBare(table))
    return true
}

func (k Keeper) DropInsertFilter(ctx sdk.Context, appId uint, owner sdk.AccAddress, tableName string) bool {
    table, err := k.GetTable(ctx, appId, tableName)
    if err != nil {
        return false
    }

    table.Filter = ""
    store := ctx.KVStore(k.storeKey)
    store.Set([]byte(getTableKey(appId, table.Name)), k.cdc.MustMarshalBinaryBare(table))
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
    if !validateTriggerSyntax(k, ctx, appId, tableName, trigger) {
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

    store := ctx.KVStore(k.storeKey)
    store.Set([]byte(getTableKey(appId, table.Name)), k.cdc.MustMarshalBinaryBare(table))
    return true
}

func (k Keeper) DropTrigger(ctx sdk.Context, appId uint, owner sdk.AccAddress, tableName string) bool {
    table, err := k.GetTable(ctx, appId, tableName)
    if err != nil {
        return false
    }

    table.Trigger = ""
    store := ctx.KVStore(k.storeKey)
    store.Set([]byte(getTableKey(appId, table.Name)), k.cdc.MustMarshalBinaryBare(table))
    return true
}

func (k Keeper) GetOption(ctx sdk.Context, appId uint, tableName string) ([]string, error) {
    store := ctx.KVStore(k.storeKey)
    key := getTableOptionsKey(appId, tableName)
    bz := store.Get([]byte(key))
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

    store := ctx.KVStore(k.storeKey)
    key := getColumnOptionsKey(appId, tableName, fieldName)
    var options []string
    var result []string

    bz := store.Get([]byte(key))
    if bz != nil {
        k.cdc.MustUnmarshalBinaryBare(bz, &options)
    }

    optionExisted := isColumnOptionIncluded(options, option)

    if action == "add" {
        if optionExisted {
            return false
        } else {
            switch types.FieldOption(option) {
            // unique field should not contain duplicate values before becoming unique
            case types.FLDOPT_UNIQUE:
                if !isColumnValuesUnique(k, ctx, appId, tableName, fieldName) {
                    return false
                }
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
        store.Set([]byte(key), k.cdc.MustMarshalBinaryBare(result))
    } else {
        store.Delete([]byte(key))
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
            store := ctx.KVStore(k.storeKey)
            store.Set([]byte(getTableKey(appId, table.Name)), k.cdc.MustMarshalBinaryBare(table))
            return true
        }
    }
    return false
}

func (k Keeper) GetColumnOption(ctx sdk.Context, appId uint, tableName string, fieldName string) ([]string, error) {
    store := ctx.KVStore(k.storeKey)
    key := getColumnOptionsKey(appId, tableName, fieldName)
    bz := store.Get([]byte(key))
    if bz == nil {
        return []string{}, nil
    }
    var options []string
    k.cdc.MustUnmarshalBinaryBare(bz, &options)
    return options, nil
}

/////////////////////////////
//                         //
// index related functions //
//                         //
/////////////////////////////

// for now we only support indexes on single field
func (k Keeper) CreateIndex(ctx sdk.Context, appId uint, owner sdk.AccAddress, tableName string, field string) {
    store := ctx.KVStore(k.storeKey)
    key := getMetaTableIndexKey(appId, tableName)
    var index_fields []string

    bz := store.Get([]byte(key))
    if bz != nil {
        k.cdc.MustUnmarshalBinaryBare(bz, &index_fields)
    }
    index_fields = append(index_fields, field)
    store.Set([]byte(key), k.cdc.MustMarshalBinaryBare(index_fields))
    // TODO: to create index data for the existing records of the table
}

func (k Keeper) DropIndex(ctx sdk.Context, appId uint, owner sdk.AccAddress, tableName string, field string) {
    store := ctx.KVStore(k.storeKey)
    key := getMetaTableIndexKey(appId, tableName)
    var indexFields []string

    bz := store.Get([]byte(key))
    if bz != nil {
        k.cdc.MustUnmarshalBinaryBare(bz, &indexFields)
        for i, fld := range indexFields {
            if field == fld {
                indexFields = append(indexFields[:i], indexFields[i+1:]...)
                if len(indexFields) < 1 {
                    store.Delete([]byte(key))
                } else {
                    store.Set([]byte(key), k.cdc.MustMarshalBinaryBare(indexFields))
                }
                break
            }
        }
    }

    // TODO: to delete index data for the existing records of the table
}

func (k Keeper) GetIndex(ctx sdk.Context, appId uint, tableName string) ([]string, error) {
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

////////////////////
//                //
// helper methods //
//                //
////////////////////

// to preprocess the new table field names
// to make sure the fields be lowercase
// to make sure field id be in place
func preProcessFields(fields []string) []string {
    var result = []string{"id", "created_by", "created_at"}
    m := make(map[string]bool)
    m["id"] = true

    for _, field := range fields {
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
    case types.FLDOPT_INT:
        return true
    case types.FLDOPT_NOTNULL:
        return true
    case types.FLDOPT_UNIQUE:
        return true
    case types.FLDOPT_OWN:
        return true
    case types.FLDOPT_FILE:
        return true
    }

    if types.ValidateEnumColumnOption(option) {
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

func validateInsertFilterSyntax(k Keeper, ctx sdk.Context, appId uint, tableName string, filter string) bool {
    fn1 := getScriptValidationCallbackOne(k, ctx, appId, tableName)
    fn2 := getScriptValidationCallbackTwo(k, ctx, appId, tableName)

    parser := ss.NewParser(strings.NewReader(filter), fn1, fn2)
    err := parser.ParseFilter()
    if err != nil {
        return false
    }
    return true
}

func validateTriggerSyntax(k Keeper, ctx sdk.Context, appId uint, tableName string, trigger string) bool {
    fn1 := getScriptValidationCallbackOne(k, ctx, appId, tableName)
    fn2 := getScriptValidationCallbackTwo(k, ctx, appId, tableName)

    parser := ss.NewParser(strings.NewReader(trigger), fn1, fn2)
    err := parser.ParseTrigger()
    if err != nil {
        return false
    }
    return true
}

func isColumnValuesUnique(k Keeper, ctx sdk.Context, appId uint, tableName string, fieldName string) bool {
    store := ctx.KVStore(k.storeKey)
    flag := make(map[string]bool)

    start, end := getFieldDataIteratorStartAndEndKey(appId, tableName, fieldName)
    iter := store.Iterator([]byte(start), []byte(end))
    var mold string
    for ; iter.Valid(); iter.Next() {
        val := iter.Value()
        k.cdc.MustUnmarshalBinaryBare(val, &mold)
        if _, ok := flag[mold]; ok {
            return false
        }
        flag[mold] = true
    }
    return true
}
