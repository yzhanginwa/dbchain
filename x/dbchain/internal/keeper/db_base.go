package keeper

import (
    "fmt"
    "errors"
    sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) getTableFields(ctx sdk.Context, appId uint, tableName string) ([]string, error) {
    table, err := k.GetTable(ctx, appId, tableName)
    if err != nil {
        return nil, errors.New(fmt.Sprintf("Failed to access table %s", tableName))
    }
    fieldNames := table.Fields
    return fieldNames, nil
}

