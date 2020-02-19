package keeper

import (
    "fmt"
    "errors"
    "strings"
    sdk "github.com/cosmos/cosmos-sdk/types"
    "github.com/yzhanginwa/cosmos-api/x/cosmosapi/internal/types"
    "github.com/yzhanginwa/cosmos-api/x/cosmosapi/internal/utils"
)

/////////////////////////////
//                         //
// table related functions //
//                         //
/////////////////////////////

func (k Keeper) getTables(ctx sdk.Context) ([]string, error) {
    store := ctx.KVStore(k.storeKey)
    tablesKey := getTablesKey()
    bz := store.Get([]byte(tablesKey))
    if bz == nil {
        return nil, errors.New("No tables found")
    }
    var tableNames []string
    k.cdc.MustUnmarshalBinaryBare(bz, &tableNames)
    return tableNames, nil
}


// Check if the table is present in the store or not
func (k Keeper) IsTablePresent(ctx sdk.Context, name string) bool {
    store := ctx.KVStore(k.storeKey)
    return store.Has([]byte(getTableKey(name)))
}

func (k Keeper) IsFieldPresent(ctx sdk.Context, tableName string, field string) bool {
    table, err := k.GetTable(ctx, tableName)
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
func (k Keeper) CreateTable(ctx sdk.Context, owner sdk.AccAddress, name string, fields []string) {
    store := ctx.KVStore(k.storeKey)
    table := types.NewTable()
    table.Owner = owner
    table.Name = name
    table.Fields = preProcessFields(fields)
    store.Set([]byte(getTableKey(table.Name)), k.cdc.MustMarshalBinaryBare(table))

    var tables []string
    bz :=store.Get([]byte(getTablesKey()))
    if bz == nil {
        tables = append(tables, table.Name)
    } else {
        k.cdc.MustUnmarshalBinaryBare(bz, &tables)
        tables = append(tables, table.Name)
    }
    store.Set([]byte(getTablesKey()), k.cdc.MustMarshalBinaryBare(tables))
}

// Remove a table
func (k Keeper) DropTable(ctx sdk.Context, owner sdk.AccAddress, name string) {
    store := ctx.KVStore(k.storeKey)
    var tables []string
    bz :=store.Get([]byte(getTablesKey()))
    if bz != nil {
        k.cdc.MustUnmarshalBinaryBare(bz, &tables)
        for i, tbl := range tables {
            if name == tbl {
                tables = append(tables[:i], tables[i+1:]...)
                if len(tables) < 1 {
                    store.Delete([]byte(getTablesKey()))
                } else {
                    store.Set([]byte(getTablesKey()), k.cdc.MustMarshalBinaryBare(tables))
                }
                store.Delete([]byte(getTableKey(name)))
                break
            }
        }
    }
}

// Get a table 
func (k Keeper) GetTable(ctx sdk.Context, name string) (types.Table, error) {
    store := ctx.KVStore(k.storeKey)
    bz := store.Get([]byte(getTableKey(name)))
    if bz == nil {
        return types.Table{}, errors.New(fmt.Sprintf("table %s not found", name))
    }
    var table types.Table
    k.cdc.MustUnmarshalBinaryBare(bz, &table)
    return table, nil
}

// Add a field
func (k Keeper) AddColumn(ctx sdk.Context, name string, field string) (bool, error) {
    table, err := k.GetTable(ctx, name)
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
    store.Set([]byte(getTableKey(table.Name)), k.cdc.MustMarshalBinaryBare(table))
    return true, nil
}

// Remove a field
func (k Keeper) DropColumn(ctx sdk.Context, name string, field string) (bool, error){
    table, err := k.GetTable(ctx, name)
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
    store.Set([]byte(getTableKey(table.Name)), k.cdc.MustMarshalBinaryBare(table))
    return true, nil
}

// Rename a field
func (k Keeper) RenameColumn(ctx sdk.Context, name string, oldField string, newField string) (bool, error) {
    table, err := k.GetTable(ctx, name)
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
    store.Set([]byte(getTableKey(table.Name)), k.cdc.MustMarshalBinaryBare(table))
    return true, nil
}

func (k Keeper) ModifyOption(ctx sdk.Context, owner sdk.AccAddress, tableName string, action string, option string) {
    store := ctx.KVStore(k.storeKey)
    key := getTableOptionsKey(tableName)
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

func (k Keeper) GetOption(ctx sdk.Context, tableName string) ([]string, error) {
    store := ctx.KVStore(k.storeKey)
    key := getTableOptionsKey(tableName)
    bz := store.Get([]byte(key))
    if bz == nil {
        return []string{}, nil
    }
    var options []string
    k.cdc.MustUnmarshalBinaryBare(bz, &options)
    return options, nil
}

func (k Keeper) ModifyColumnOption(ctx sdk.Context, owner sdk.AccAddress, tableName string, fieldName string, action string, option string) {
    store := ctx.KVStore(k.storeKey)
    key := getFieldOptionsKey(tableName, fieldName)
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

func (k Keeper) GetFieldOption(ctx sdk.Context, tableName string, fieldName string) ([]string, error) {
    store := ctx.KVStore(k.storeKey)
    key := getFieldOptionsKey(tableName, fieldName)
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
func (k Keeper) CreateIndex(ctx sdk.Context, owner sdk.AccAddress, tableName string, field string) {
    store := ctx.KVStore(k.storeKey)
    key := getMetaTableIndexKey(tableName)
    var index_fields []string

    bz := store.Get([]byte(key))
    if bz != nil {
        k.cdc.MustUnmarshalBinaryBare(bz, &index_fields)
    }
    index_fields = append(index_fields, field)
    store.Set([]byte(key), k.cdc.MustMarshalBinaryBare(index_fields))
    // TODO: to create index data for the existing records of the table
}

func (k Keeper) DropIndex(ctx sdk.Context, owner sdk.AccAddress, tableName string, field string) {
    store := ctx.KVStore(k.storeKey)
    key := getMetaTableIndexKey(tableName)
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

func (k Keeper) GetIndex(ctx sdk.Context, tableName string) ([]string, error) {
    store := ctx.KVStore(k.storeKey)
    key := getMetaTableIndexKey(tableName)
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
