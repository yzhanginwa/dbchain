package types

import (
    sdk "github.com/cosmos/cosmos-sdk/types"
    sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

//////////////////
//              //
// MsgAddFriend //
//              //
//////////////////

// MsgAddFriend defines a CreateTable message
type MsgAddFriend struct {
    Owner sdk.AccAddress    `json:"owner"`
    FriendAddr string       `json:"friend_addr"`
    FriendName string       `json:"friend_name"`
}

// NewMsgAddFriend is a constructor function for MsgCreatTable
func NewMsgAddFriend(owner sdk.AccAddress, friendAddr string, friendName string) MsgAddFriend {
    return MsgAddFriend {
        Owner: owner,
        FriendAddr: friendAddr,
        FriendName: friendName,
    }
}

// Route should return the name of the module
func (msg MsgAddFriend) Route() string { return RouterKey }

// Type should return the action
func (msg MsgAddFriend) Type() string { return "add_friend" }

// ValidateBasic runs stateless checks on the message
func (msg MsgAddFriend) ValidateBasic() error {
    if msg.Owner.Empty() {
        return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, msg.Owner.String())
    }
    if len(msg.FriendAddr) == 0 {
        return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "Friend address cannot be empty")
    }
    if len(msg.FriendName) == 0 {
        return sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "Friend name cannot be empty")
    }

    return nil
}

// GetSignBytes encodes the message for signing
func (msg MsgAddFriend) GetSignBytes() []byte {
    return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners defines whose signature is required
func (msg MsgAddFriend) GetSigners() []sdk.AccAddress {
    return []sdk.AccAddress{msg.Owner}
}
