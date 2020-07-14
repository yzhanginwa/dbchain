package keeper

import (
    "errors"
    sdk "github.com/cosmos/cosmos-sdk/types"
    "github.com/yzhanginwa/dbchain/x/dbchain/internal/utils"
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
