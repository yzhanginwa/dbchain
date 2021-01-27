package keeper

import (
    "errors"
    "fmt"
    "strconv"
    sdk "github.com/cosmos/cosmos-sdk/types"
    "github.com/yzhanginwa/dbchain/x/dbchain/internal/utils"
    "github.com/yzhanginwa/dbchain/x/dbchain/internal/other"
)

// return the function which checks whether a table contains a field
func getScriptValidationCallbackOne(k Keeper, ctx sdk.Context, appId uint, tableName string) func(string, string) bool {
    return func(table, field string) bool {
        if table == "" {
            table = tableName
        }
        fieldNames, err := k.getTableFields(ctx, appId, table)
        if err != nil { return false }
        return utils.StringIncluded(fieldNames, field)
    }
}

// return the function which returns the parent table name of a reference field in certain table
func getScriptValidationCallbackTwo(k Keeper, ctx sdk.Context, appId uint, tableName string) func(string, string) (string, error) {
    return func(table, field string) (string, error) {
        //for now we get parent table solely from the field name
        if tn, ok := utils.GetTableNameFromForeignKey(field); ok {
            return tn, nil
        } else {
            return "", errors.New("Wrong reference field name")
        }
    }
}

func getGetFieldValueCallback(k Keeper, ctx sdk.Context, appId uint, owner sdk.AccAddress) func(string, uint, string) string {
    return func(tableName string, id uint, fieldName string) string {
        result, _ := k.FindField(ctx, appId, tableName, id, fieldName)
        return result
    }
}

func getGetTableValueCallback(k Keeper, ctx sdk.Context, appId uint, owner sdk.AccAddress) func([](map[string]string)) [](map[string]string) {
    return func(querierObjs [](map[string]string)) [](map[string]string) {
        qo := map[string]string{
            "method": "select",
            "fields": "id",
        }
        newQuerierObjs := append(querierObjs, qo)
        result, err := querierSuperHandler(ctx, k, appId, newQuerierObjs, owner)
        if err != nil {
            return [](map[string]string){}
        }
        return result
    }
}

func getInsertCallback(k Keeper, ctx sdk.Context, appId uint, owner sdk.AccAddress) func(string, map[string]string) {
    return func(tableName string, value map[string]string) {
        id, err := getNextId(k, ctx, appId, tableName)
        if err != nil {
            return
        }

        value["id"] = strconv.Itoa(int(id))
        value["created_by"] = owner.String()
        value["created_at"] = fmt.Sprintf("%d", other.GetCurrentBlockTime().Unix())

        k.Write(ctx, appId, tableName, id, value, owner)
    }
}
