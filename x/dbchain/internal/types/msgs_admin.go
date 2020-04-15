package types

import (
    sdk "github.com/cosmos/cosmos-sdk/types"
    sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

////////////////////////
//                    //
// MsgAddGroupMember //
//                    //
////////////////////////

type MsgAddGroupMember struct {
    AppCode string              `json:"app_code"`
    Group string                `json:"group"`
    Member sdk.AccAddress       `json:"member"`
    Owner sdk.AccAddress        `json:"owner"`
}

func NewMsgAddGroupMember(appCode string, group string, member sdk.AccAddress, owner sdk.AccAddress) MsgAddGroupMember {
    return MsgAddGroupMember {
        AppCode: appCode,
        Group: group,
        Member: member,
        Owner: owner,
    }
}

// Route should return the name of the module
func (msg MsgAddGroupMember) Route() string { return RouterKey }

// Type should return the action
func (msg MsgAddGroupMember) Type() string { return "add_group_member" }

// ValidateBasic runs stateless checks on the message
func (msg MsgAddGroupMember) ValidateBasic() error {
    if len(msg.AppCode) == 0 {
        return sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "App code cannot be empty")
    }
    if len(msg.Group) == 0 {
        return sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "Group name cannot be empty")
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
func (msg MsgAddGroupMember) GetSignBytes() []byte {
    return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners defines whose signature is required
func (msg MsgAddGroupMember) GetSigners() []sdk.AccAddress {
    return []sdk.AccAddress{msg.Owner}
}

