package types

import (
    sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var (
    ErrDatabaseNotExist = sdkerrors.Register(ModuleName, 1, "database does not exist")
)

