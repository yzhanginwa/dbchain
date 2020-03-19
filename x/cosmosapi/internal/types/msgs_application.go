package types

import (
    sdk "github.com/cosmos/cosmos-sdk/types"
)

//////////////////////////
//                      //
// MsgCreateApplication //
//                      //
//////////////////////////

// MsgCreateApplication defines a CreateApplication message
type MsgCreateApplication struct {
    Owner sdk.AccAddress `json:"owner"`
    Name string          `json:"name"`
    Description string   `json:"description"`
    Permissioned bool    `json:"permissioned"`
}

// NewMsgCreateApplication is a constructor function for MsgCreatTable
func NewMsgCreateApplication(owner sdk.AccAddress, name string, description string, permissioned bool) MsgCreateApplication {
    return MsgCreateApplication {
        Owner: owner,
        Name: name,
        Description: description,
        Permissioned: permissioned,
    }
}

// Route should return the name of the module
func (msg MsgCreateApplication) Route() string { return RouterKey }

// Type should return the action
func (msg MsgCreateApplication) Type() string { return "create_application" }

// ValidateBasic runs stateless checks on the message
func (msg MsgCreateApplication) ValidateBasic() sdk.Error {
    if msg.Owner.Empty() {
        return sdk.ErrInvalidAddress(msg.Owner.String())
    }
    if len(msg.Name) == 0 {
        return sdk.ErrUnknownRequest("Application name cannot be empty")
    }
    if len(msg.Description) == 0 {
        return sdk.ErrUnknownRequest("Application description cannot be empty")
    }
    return nil
}

// GetSignBytes encodes the message for signing
func (msg MsgCreateApplication) GetSignBytes() []byte {
    return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners defines whose signature is required
func (msg MsgCreateApplication) GetSigners() []sdk.AccAddress {
    return []sdk.AccAddress{msg.Owner}
}

////////////////////////
//                    //
// MsgAddDatabaseUser //
//                    //
////////////////////////

type MsgAddDatabaseUser struct {
    Owner sdk.AccAddress `json:"owner"`
    AppCode string       `json:"app_code"`
    User sdk.AccAddress  `json:"description"`
}

// NewMsgAddDatabaseUser is a constructor function for MsgCreatTable
func NewMsgAddDatabaseUser(owner sdk.AccAddress, appcode string, user sdk.AccAddress) MsgAddDatabaseUser {
    return MsgAddDatabaseUser {
        Owner: owner,
        AppCode: appcode,
        User: user,
    }
}

// Route should return the name of the module
func (msg MsgAddDatabaseUser) Route() string { return RouterKey }

// Type should return the action
func (msg MsgAddDatabaseUser) Type() string { return "add_database_user" }

// ValidateBasic runs stateless checks on the message
func (msg MsgAddDatabaseUser) ValidateBasic() sdk.Error {
    if msg.Owner.Empty() {
        return sdk.ErrInvalidAddress(msg.Owner.String())
    }
    if len(msg.AppCode) == 0 {
        return sdk.ErrUnknownRequest("Application Code cannot be empty")
    }
    if msg.User.Empty() {
        return sdk.ErrUnknownRequest("User cannot be empty")
    }
    return nil
}

// GetSignBytes encodes the message for signing
func (msg MsgAddDatabaseUser) GetSignBytes() []byte {
    return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners defines whose signature is required
func (msg MsgAddDatabaseUser) GetSigners() []sdk.AccAddress {
    return []sdk.AccAddress{msg.Owner}
}
