package types

import (
    sdkerrors "github.com/dbchaincloud/cosmos-sdk/types/errors"
)

var (
    ErrDatabaseNotExist = sdkerrors.Register(ModuleName, 1, "database does not exist")
)

