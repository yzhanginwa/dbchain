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

    if field == "id" {
        return false, errors.New(fmt.Sprintf("cannot remove field id"))
    }

    var foundField = false 
    for i, fld := range table.Fields {
        if field == fld {
            foundField = true
            table.Fields = append(table.Fields[:i], table.Fields[i+1:]...)
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
    store := ctx.KVStore(k.storeKey)
    key := getTableInsertFilterKey(appId, tableName)
    var filters []string

    bz := store.Get([]byte(key))
    if bz != nil {
        k.cdc.MustUnmarshalBinaryBare(bz, &filters)
    }

    filters = append(filters, filter)
    store.Set([]byte(key), k.cdc.MustMarshalBinaryBare(filters))
    return true
}

func (k Keeper) DropInsertFilter(ctx sdk.Context, appId uint, owner sdk.AccAddress, tableName string, index int) bool {
    store := ctx.KVStore(k.storeKey)
    key := getTableInsertFilterKey(appId, tableName)
    var filters []string

    bz := store.Get([]byte(key))
    if bz != nil {
        k.cdc.MustUnmarshalBinaryBare(bz, &filters)
    }

    l := len(filters)
    if index < 0 || index > (l - 1) {
        return false
    }

    filters[index] = filters[l-1]
    filters[l-1] = ""
    filters = filters[:l-1]
    store.Set([]byte(key), k.cdc.MustMarshalBinaryBare(filters))
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

func (k Keeper) ModifyColumnOption(ctx sdk.Context, appId uint, owner sdk.AccAddress, tableName string, fieldName string, action string, option string) {
    store := ctx.KVStore(k.storeKey)
    key := getColumnOptionsKey(appId, tableName, fieldName)
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

func validateInsertFilterSyntax(k Keeper, ctx sdk.Context, appId uint, tableName string, filter string) bool {
    parser := ss.NewParser(strings.NewReader(filter),
        func(table, field string) bool {
            fieldNames, err := k.getTableFields(ctx, appId, table)
            if err != nil { return false }
            return utils.StringIncluded(fieldNames, field)
        },

        func(table, field string) (string, error) {
            //for now we get parent table solely from the field name
            if tn, ok := utils.GetTableNameFromForeignKey(field); ok {
                return tn, nil
            } else {
                return "", errors.New("Wrong reference field name")
            }
        },
    )
    err := parser.FilterCondition()
    if err != nil {
        return false
    }
    return true
}
