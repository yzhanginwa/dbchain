package types

import (
	//"encoding/base64"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

//////////////////////////
//                      //
// MsgSaveUserPrivateKey//
//                      //
//////////////////////////

// MsgCreateApplication defines a CreateApplication message
type MsgSaveUserPrivateKey struct {
	Owner   sdk.AccAddress `json:"owner"`
	User    string         `json:"user"`
	KeyInfo string         `json:"key_info"`
}
// NewMsgCreateApplication is a constructor function for MsgCreatTable
func NewMsgSaveUserPrivateKey(owner sdk.AccAddress, user,keyInfo string) MsgSaveUserPrivateKey {
	return MsgSaveUserPrivateKey {
		Owner: owner,
		User: user,
		KeyInfo: keyInfo,
	}
}

// Route should return the name of the module
func (msg MsgSaveUserPrivateKey) Route() string { return RouterKey }

// Type should return the action
func (msg MsgSaveUserPrivateKey) Type() string { return "save_user_private_key" }

// ValidateBasic runs stateless checks on the message
func (msg MsgSaveUserPrivateKey) ValidateBasic() error {
	if msg.Owner.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, msg.Owner.String())
	}
	if len(msg.KeyInfo) == 0 {
		return sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "keyInfo name cannot be empty")
	}
	return nil
}

// GetSignBytes encodes the message for signing
func (msg MsgSaveUserPrivateKey) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners defines whose signature is required
func (msg MsgSaveUserPrivateKey) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Owner}
}
