package types

import (
    sdk "github.com/cosmos/cosmos-sdk/types"
    sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

////////////////////
//                //
// MsgAddfunction //
//                //
////////////////////

type MsgAddFunction struct {
    Owner sdk.AccAddress `json:"owner"`
    AppCode string       `json:"app_code"`
    FunctionName string  `json:"function_name"`
    Parameter string     `json:"parameter"`
    Body string          `json:"body"`
}

// NewMsgCreatePoll is a constructor function for MsgCreatPoll
func NewMsgAddFunction(owner sdk.AccAddress, appCode, functionName, parameter, body string) MsgAddFunction {
    return MsgAddFunction {
        Owner: owner,
        AppCode: appCode,
        FunctionName: functionName,
        Parameter: parameter,
        Body: body,
    }
}

// Route should return the name of the module
func (msg MsgAddFunction) Route() string { return RouterKey }

// Type should return the action
func (msg MsgAddFunction) Type() string { return "add_function" }

// ValidateBasic runs stateless checks on the message
func (msg MsgAddFunction) ValidateBasic() error {
    if msg.Owner.Empty() {
        return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, msg.Owner.String())
    }
    if len(msg.AppCode) == 0 {
        return sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "App code cannot be empty")
    }
    if len(msg.FunctionName) == 0 {
        return sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "Function name cannot be empty")
    }
    if len(msg.Body) ==0 {
        return sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "Body cannot be empty")
    }
    return nil
}

// GetSignBytes encodes the message for signing
func (msg MsgAddFunction) GetSignBytes() []byte {
    return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners defines whose signature is required
func (msg MsgAddFunction) GetSigners() []sdk.AccAddress {
    return []sdk.AccAddress{msg.Owner}
}

/////////////////////
//                 //
// MsgCallFunction //
//                 //
/////////////////////

type MsgCallFunction struct {
    Owner sdk.AccAddress `json:"owner"`
    AppCode string       `json:"app_code"`
    FunctionName string  `json:"function_name"`
    Argument string      `json:"argument"`
}

// NewMsgCreatePoll is a constructor function for MsgCreatPoll
func NewMsgCallFunction(owner sdk.AccAddress, appCode, functionName, argument string) MsgCallFunction {
    return MsgCallFunction {
        Owner: owner,
        AppCode: appCode,
        FunctionName: functionName,
        Argument: argument,
    }
}

// Route should return the name of the module
func (msg MsgCallFunction) Route() string { return RouterKey }

// Type should return the action
func (msg MsgCallFunction) Type() string { return "call_function" }

// ValidateBasic runs stateless checks on the message
func (msg MsgCallFunction) ValidateBasic() error {
    if msg.Owner.Empty() {
        return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, msg.Owner.String())
    }
    if len(msg.AppCode) == 0 {
        return sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "App code cannot be empty")
    }
    if len(msg.FunctionName) == 0 {
        return sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "Function name cannot be empty")
    }
    return nil
}

// GetSignBytes encodes the message for signing
func (msg MsgCallFunction) GetSignBytes() []byte {
    return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners defines whose signature is required
func (msg MsgCallFunction) GetSigners() []sdk.AccAddress {
    return []sdk.AccAddress{msg.Owner}
}
