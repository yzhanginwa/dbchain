package types

import (
    sdk "github.com/cosmos/cosmos-sdk/types"
    sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
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
func (msg MsgCreateApplication) ValidateBasic() error {
    if msg.Owner.Empty() {
        return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, msg.Owner.String())
    }
    if len(msg.Name) == 0 {
        return sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "Application name cannot be empty")
    }
    if len(msg.Description) == 0 {
        return sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "Application description cannot be empty")
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
func (msg MsgAddDatabaseUser) ValidateBasic() error {
    if msg.Owner.Empty() {
        return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, msg.Owner.String())
    }
    if len(msg.AppCode) == 0 {
        return sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "Application Code cannot be empty")
    }
    if msg.User.Empty() {
        return sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "User cannot be empty")
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

//////////////////////////
//                      //
// MsgCreateSysDatabase //
//                      //
//////////////////////////

// MsgCreateSysDatabase defines a CreateApplication message
type MsgCreateSysDatabase struct {
    Owner sdk.AccAddress `json:"owner"`
}

// NewMsgCreateSysDatabase is a constructor function for MsgCreatTable
func NewMsgCreateSysDatabase(owner sdk.AccAddress) MsgCreateSysDatabase {
    return MsgCreateSysDatabase {
        Owner: owner,
    }
}

// Route should return the name of the module
func (msg MsgCreateSysDatabase) Route() string { return RouterKey }

// Type should return the action
func (msg MsgCreateSysDatabase) Type() string { return "create_sys_database" }

// ValidateBasic runs stateless checks on the message
func (msg MsgCreateSysDatabase) ValidateBasic() error {
    if msg.Owner.Empty() {
        return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, msg.Owner.String())
    }
    return nil
}

// GetSignBytes encodes the message for signing
func (msg MsgCreateSysDatabase) GetSignBytes() []byte {
    return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners defines whose signature is required
func (msg MsgCreateSysDatabase) GetSigners() []sdk.AccAddress {
    return []sdk.AccAddress{msg.Owner}
}

///////////////////////////
//                       //
// NewMsgSetSchemaStatus //
//                       //
///////////////////////////

// MsgSetSchemaStatus sets the status of SchemaFrozen of a database
type MsgSetSchemaStatus struct {
    Owner sdk.AccAddress `json:"owner"`
    AppCode string       `json:"app_code"`
    Status string        `json:"status"`
}

// NewMsgSetSchemaStatus is a constructor function for MsgCreatTable
func NewMsgSetSchemaStatus(owner sdk.AccAddress, appCode, status string) MsgSetSchemaStatus {
    return MsgSetSchemaStatus {
        Owner: owner,
        AppCode: appCode,
        Status: status,
    }
}

// Route should return the name of the module
func (msg MsgSetSchemaStatus) Route() string { return RouterKey }

// Type should return the action
func (msg MsgSetSchemaStatus) Type() string { return "set_schema_status" }

// ValidateBasic runs stateless checks on the message
func (msg MsgSetSchemaStatus) ValidateBasic() error {
    if msg.Owner.Empty() {
        return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, msg.Owner.String())
    }
    if len(msg.AppCode) == 0 {
        return sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "Application Code cannot be empty")
    }
    if len(msg.Status) == 0 {
        return sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "Status cannot be empty")
    }
    if msg.Status != "frozen" && msg.Status != "unfrozen" {
        return sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "Status has to be either frozen of unfrozen")
    }
    return nil
}

// GetSignBytes encodes the message for signing
func (msg MsgSetSchemaStatus) GetSignBytes() []byte {
    return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners defines whose signature is required
func (msg MsgSetSchemaStatus) GetSigners() []sdk.AccAddress {
    return []sdk.AccAddress{msg.Owner}
}
