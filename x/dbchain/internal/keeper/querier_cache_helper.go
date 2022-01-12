package keeper

import (
    sdk "github.com/dbchaincloud/cosmos-sdk/types"
)

// NOTE!!! DO NOT cache the callback because the ctx and keeper may become obsolete or old
func GetQuerierCacheCallback1(ctx sdk.Context, keeper Keeper) func(uint, string) bool {
    return func(appId uint, tableName string) bool {
        return keeper.isTablePublic(ctx, appId, tableName)
    }
}
