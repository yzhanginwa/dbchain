package keeper

import (
    "github.com/yzhanginwa/dbchain/x/dbchain/internal/utils"
)

func isSystemField(fieldName string) bool {
    systemFields := []string{"id", "created_by", "created_at"}
    return utils.ItemExists(systemFields, fieldName)
}

