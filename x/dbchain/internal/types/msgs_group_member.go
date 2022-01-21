package types

import (
    sdk "github.com/cosmos/cosmos-sdk/types"
    sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

//////////////////////////
//                      //
// MsgModifyGroupMember //
//                      //
//////////////////////////

type MsgModifyGroupMember struct {
    AppCode string              `json:"app_code"`
    Group string                `json:"group"`
    Action string               `json:"action"`
    Member sdk.AccAddress       `json:"member"`
    Owner sdk.AccAddress        `json:"owner"`
}

func NewMsgModifyGroupMember(appCode string, group string, action string, member sdk.AccAddress, owner sdk.AccAddress) MsgModifyGroupMember {
    return MsgModifyGroupMember {
        AppCode: appCode,
        Group: group,
        Action: action,
        Member: member,
        Owner: owner,
    }
}

// Route should return the name of the module
func (msg MsgModifyGroupMember) Route() string { return RouterKey }

// Type should return the action
func (msg MsgModifyGroupMember) Type() string { return "modify_group_member" }

// ValidateBasic runs stateless checks on the message
func (msg MsgModifyGroupMember) ValidateBasic() error {
    if len(msg.AppCode) == 0 {
        return sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "App code cannot be empty")
    }
    if len(msg.Group) == 0 {
        return sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "Group name cannot be empty")
    }
    if msg.Action != "add" && msg.Action != "drop" {
        return sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "Wrong action")
    }
    if msg.Member.Empty() {
        return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, msg.Member.String())
    }
    if msg.Owner.Empty() {
        return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, msg.Owner.String())
    }
    return nil
}

// GetSignBytes encodes the message for signing
func (msg MsgModifyGroupMember) GetSignBytes() []byte {
    return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners defines whose signature is required
func (msg MsgModifyGroupMember) GetSigners() []sdk.AccAddress {
    return []sdk.AccAddress{msg.Owner}
}
