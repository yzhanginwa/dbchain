package keeper

import (
    "github.com/yzhanginwa/dbchain/x/dbchain/internal/utils"
    sdk "github.com/cosmos/cosmos-sdk/types"
)

func isSystemField(fieldName string) bool {
    systemFields := []string{"id", "created_by", "created_at"}
    return utils.ItemExists(systemFields, fieldName)
}

func (k Keeper) isAdmin(ctx sdk.Context, appId uint, addr sdk.AccAddress) bool {
    admins := k.getGroupMembers(ctx, appId, "admin")
    if utils.AddressIncluded(admins, addr) {
        return true
    }
    return false
}

