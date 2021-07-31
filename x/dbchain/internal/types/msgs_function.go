package types

import (
    "encoding/base64"
    sdk "github.com/dbchaincloud/cosmos-sdk/types"
    sdkerrors "github.com/dbchaincloud/cosmos-sdk/types/errors"
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
    Description string   `json:"description"`
    Body string          `json:"body"`
}

// NewMsgCreatePoll is a constructor function for MsgCreatPoll
func NewMsgAddFunction(owner sdk.AccAddress, appCode, functionName, description, body string) MsgAddFunction {
    return MsgAddFunction {
        Owner: owner,
        AppCode: appCode,
        FunctionName: functionName,
        Description: base64.StdEncoding.EncodeToString([]byte(description)),
        Body: base64.StdEncoding.EncodeToString([]byte(body)),
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


/////////////////////
//                 //
// MsgDropFunction //
//                 //
/////////////////////

type MsgDropFunction struct {
    Owner sdk.AccAddress `json:"owner"`
    AppCode string       `json:"app_code"`
    FunctionName string  `json:"function_name"`
}

// NewMsgCreatePoll is a constructor function for MsgCreatPoll
func NewMsgDropFunction(owner sdk.AccAddress, appCode, functionName string) MsgDropFunction {
    return MsgDropFunction {
        Owner: owner,
        AppCode: appCode,
        FunctionName: functionName,
    }
}

// Route should return the name of the module
func (msg MsgDropFunction) Route() string { return RouterKey }

// Type should return the action
func (msg MsgDropFunction) Type() string { return "call_function" }

// ValidateBasic runs stateless checks on the message
func (msg MsgDropFunction) ValidateBasic() error {
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
func (msg MsgDropFunction) GetSignBytes() []byte {
    return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners defines whose signature is required
func (msg MsgDropFunction) GetSigners() []sdk.AccAddress {
    return []sdk.AccAddress{msg.Owner}
}