package types

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// x/bank module sentinel errors
var (
	ErrNoInputs            = sdkerrors.Register(ModuleName, 101, "no inputs to send transaction")
	ErrNoOutputs           = sdkerrors.Register(ModuleName, 102, "no outputs to send transaction")
	ErrInputOutputMismatch = sdkerrors.Register(ModuleName, 103, "sum inputs != sum outputs")
	ErrSendDisabled        = sdkerrors.Register(ModuleName, 104, "send transactions are disabled")
)
