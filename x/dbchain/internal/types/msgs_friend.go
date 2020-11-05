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
    OwnerName string        `json:"owner_name"`
    FriendAddr string       `json:"friend_addr"`
    FriendName string       `json:"friend_name"`
}

// NewMsgAddFriend is a constructor function for MsgCreatTable
func NewMsgAddFriend(owner sdk.AccAddress, ownerName string, friendAddr string, friendName string) MsgAddFriend {
    return MsgAddFriend {
        Owner: owner,
        OwnerName: ownerName,
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
    if len(msg.OwnerName) == 0 {
        return sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "Owner name cannot be empty")
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

///////////////////
//               //
// MsgDropFriend //
//               //
///////////////////

type MsgDropFriend struct {
    Owner sdk.AccAddress    `json:"owner"`
    FriendAddr string       `json:"friend_addr"`
}

// NewMsgDropFriend is a constructor function for MsgCreatTable
func NewMsgDropFriend(owner sdk.AccAddress, friendAddr string) MsgDropFriend {
    return MsgDropFriend {
        Owner: owner,
        FriendAddr: friendAddr,
    }
}

// Route should return the name of the module
func (msg MsgDropFriend) Route() string { return RouterKey }

// Type should return the action
func (msg MsgDropFriend) Type() string { return "drop_friend" }

// ValidateBasic runs stateless checks on the message
func (msg MsgDropFriend) ValidateBasic() error {
    if msg.Owner.Empty() {
        return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, msg.Owner.String())
    }
    if len(msg.FriendAddr) == 0 {
        return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "Friend address cannot be empty")
    }

    return nil
}

// GetSignBytes encodes the message for signing
func (msg MsgDropFriend) GetSignBytes() []byte {
    return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners defines whose signature is required
func (msg MsgDropFriend) GetSigners() []sdk.AccAddress {
    return []sdk.AccAddress{msg.Owner}
}

//////////////////////
//                  //
// MsgRespondFriend //
//                  //
//////////////////////

// MsgRespondFriend defines a CreateTable message
type MsgRespondFriend struct {
    Owner sdk.AccAddress    `json:"owner"`
    FriendAddr string       `json:"friend_addr"`
    Action string           `json:"action"`
}

// NewMsgRespondFriend is a constructor function for MsgCreatTable
func NewMsgRespondFriend(owner sdk.AccAddress, friendAddr string, action string) MsgRespondFriend {
    return MsgRespondFriend {
        Owner: owner,
        FriendAddr: friendAddr,
        Action: action,
    }
}

// Route should return the name of the module
func (msg MsgRespondFriend) Route() string { return RouterKey }

// Type should return the action
func (msg MsgRespondFriend) Type() string { return "respond_friend" }

// ValidateBasic runs stateless checks on the message
func (msg MsgRespondFriend) ValidateBasic() error {
    if msg.Owner.Empty() {
        return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, msg.Owner.String())
    }
    if len(msg.FriendAddr) == 0 {
        return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "Friend address cannot be empty")
    }
    if len(msg.Action) == 0 {
        return sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "Friend name cannot be empty")
    }

    return nil
}

// GetSignBytes encodes the message for signing
func (msg MsgRespondFriend) GetSignBytes() []byte {
    return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners defines whose signature is required
func (msg MsgRespondFriend) GetSigners() []sdk.AccAddress {
    return []sdk.AccAddress{msg.Owner}
}
