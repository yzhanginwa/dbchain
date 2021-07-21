package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

///////////////////////////////
//                           //
// MsgChangeP2PTransferLimit //
//                           //
///////////////////////////////

// MsgChangeP2PTransferLimit defines a  limit of transfer message
type MsgModifyP2PTransferLimit struct {
	Owner sdk.AccAddress `json:"owner"`
	Limit bool       `json:"limit"`
}

// NewMsgChangeP2PTransferLimit is a constructor function for MsgChangeP2PTransferLimit
func NewMsgModifyP2PTransferLimit(owner sdk.AccAddress, limit bool) MsgModifyP2PTransferLimit {
	return MsgModifyP2PTransferLimit {
		Owner: owner,
		Limit: limit,
	}
}

// Route should return the name of the module
func (msg MsgModifyP2PTransferLimit) Route() string { return RouterKey }

// Type should return the action
func (msg MsgModifyP2PTransferLimit) Type() string { return "create_table" }

// ValidateBasic runs stateless checks on the message
func (msg MsgModifyP2PTransferLimit) ValidateBasic() error {
	if msg.Owner.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, msg.Owner.String())
	}
	return nil
}

// GetSignBytes encodes the message for signing
func (msg MsgModifyP2PTransferLimit) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners defines whose signature is required
func (msg MsgModifyP2PTransferLimit) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Owner}
}


///////////////////////////////
//                           //
// MsgModifySuperAdminMember //
//                           //
///////////////////////////////

// MsgChangeP2PTransferLimit is used to add or remove superAdmin
type MsgModifyChainSuperAdminMember struct {
	Owner sdk.AccAddress `json:"owner"`
	Action string       `json:"action"`  //add or remove
	Member sdk.AccAddress `json:"member"`
}

// NewMsgChangeP2PTransferLimit is a constructor function for NewMsgChangeP2PTransferLimit
func NewMsgChainModifySuperAdminMember(owner, member sdk.AccAddress, action string, ) MsgModifyChainSuperAdminMember {
	return MsgModifyChainSuperAdminMember {
		Owner: owner,
		Action: action,
		Member: member,
	}
}

// Route should return the name of the module
func (msg MsgModifyChainSuperAdminMember) Route() string { return RouterKey }

// Type should return the action
func (msg MsgModifyChainSuperAdminMember) Type() string { return "create_table" }

// ValidateBasic runs stateless checks on the message
func (msg MsgModifyChainSuperAdminMember) ValidateBasic() error {
	if msg.Owner.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, msg.Owner.String())
	}

	if msg.Member.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, msg.Member.String())
	}

	if msg.Action != "add" && msg.Action != "remove" {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "action can only be add or remove")
	}
	return nil
}

// GetSignBytes encodes the message for signing
func (msg MsgModifyChainSuperAdminMember) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners defines whose signature is required
func (msg MsgModifyChainSuperAdminMember) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Owner}
}