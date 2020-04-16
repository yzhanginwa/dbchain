package types

import (
    sdk "github.com/cosmos/cosmos-sdk/types"
    sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

////////////////////
//                //
// MsgModifyGroup //
//                //
////////////////////

type MsgModifyGroup struct {
    AppCode string           `json:"app_code"`
    Action string            `json:"action"`
    Group string             `json:"group"`
    Owner sdk.AccAddress     `json:"owner"`
}

func NewMsgModifyGroup(appCode string, action string, group string, owner sdk.AccAddress) MsgModifyGroup {
    return MsgModifyGroup {
        AppCode: appCode,
        Action: action,
        Group: group,
        Owner: owner,
    }
}

// Route should return the name of the module
func (msg MsgModifyGroup) Route() string { return RouterKey }

// Type should return the action
func (msg MsgModifyGroup) Type() string { return "modify_group" }

// ValidateBasic runs stateless checks on the message
func (msg MsgModifyGroup) ValidateBasic() error {
    if len(msg.AppCode) == 0 {
        return sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "App code cannot be empty")
    }
    if msg.Action != "add" && msg.Action != "drop" {
        return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "Wrong action")
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
func (msg MsgModifyGroup) GetSignBytes() []byte {
    return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners defines whose signature is required
func (msg MsgModifyGroup) GetSigners() []sdk.AccAddress {
    return []sdk.AccAddress{msg.Owner}
}
