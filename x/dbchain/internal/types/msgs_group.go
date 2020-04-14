package types

import (
    sdk "github.com/cosmos/cosmos-sdk/types"
    sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

/////////////////
//             //
// MsgAddGroup //
//             //
/////////////////

type MsgCreateGroup struct {
    AppCode string           `json:"app_code"`
    Group string             `json:"group"`
    Owner sdk.AccAddress     `json:"owner"`
}

func NewMsgCreateGroup(appCode string, group string, owner sdk.AccAddress) MsgCreateGroup {
    return MsgCreateGroup {
        AppCode: appCode,
        Group: group,
        Owner: owner,
    }
}

// Route should return the name of the module
func (msg MsgCreateGroup) Route() string { return RouterKey }

// Type should return the action
func (msg MsgCreateGroup) Type() string { return "create_group" }

// ValidateBasic runs stateless checks on the message
func (msg MsgCreateGroup) ValidateBasic() error {
    if len(msg.AppCode) == 0 {
        return sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "App code cannot be empty")
    }
    if len(msg.Group) == 0 {
        return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "Group name cannot be empty")
    }
    if msg.Owner.Empty() {
        return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, msg.Owner.String())
    }
    return nil
}

// GetSignBytes encodes the message for signing
func (msg MsgCreateGroup) GetSignBytes() []byte {
    return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners defines whose signature is required
func (msg MsgCreateGroup) GetSigners() []sdk.AccAddress {
    return []sdk.AccAddress{msg.Owner}
}
