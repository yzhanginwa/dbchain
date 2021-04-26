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
    Owner sdk.AccAddress    `json:"owner"`
    Name string             `json:"name"`
    Description string      `json:"description"`
    PermissionRequired bool `json:"permission_required"`
}

// NewMsgCreateApplication is a constructor function for MsgCreatTable
func NewMsgCreateApplication(owner sdk.AccAddress, name string, description string, permissionRequired bool) MsgCreateApplication {
    return MsgCreateApplication {
        Owner: owner,
        Name: name,
        Description: description,
        PermissionRequired: permissionRequired,
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

//////////////////////////
//                      //
//  MsgDropApplication  //
//                      //
//////////////////////////

type MsgDropApplication struct {
    Owner sdk.AccAddress `json:"owner"`
    AppCode string       `json:"app_code"`
}

func NewMsgDropApplication(owner sdk.AccAddress, appcode string) MsgDropApplication {
    return MsgDropApplication {
        Owner: owner,
        AppCode: appcode,
    }
}

// Route should return the name of the module
func (msg MsgDropApplication) Route() string { return RouterKey }

// Type should return the action
func (msg MsgDropApplication) Type() string { return "drop_database_user" }

// ValidateBasic runs stateless checks on the message
func (msg MsgDropApplication) ValidateBasic() error {
    if msg.Owner.Empty() {
        return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, msg.Owner.String())
    }
    if len(msg.AppCode) == 0 {
        return sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "Application Code cannot be empty")
    }
    return nil
}

// GetSignBytes encodes the message for signing
func (msg MsgDropApplication) GetSignBytes() []byte {
    return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners defines whose signature is required
func (msg MsgDropApplication) GetSigners() []sdk.AccAddress {
    return []sdk.AccAddress{msg.Owner}
}

//////////////////////////
//                      //
//  MsgDropApplication  //
//                      //
//////////////////////////

type MsgRecoverApplication struct {
    Owner sdk.AccAddress `json:"owner"`
    AppCode string       `json:"app_code"`
}

func NewMsgRecoverApplication(owner sdk.AccAddress, appcode string) MsgRecoverApplication {
    return MsgRecoverApplication {
        Owner: owner,
        AppCode: appcode,
    }
}

// Route should return the name of the module
func (msg MsgRecoverApplication) Route() string { return RouterKey }

// Type should return the action
func (msg MsgRecoverApplication) Type() string { return "drop_database_user" }

// ValidateBasic runs stateless checks on the message
func (msg MsgRecoverApplication) ValidateBasic() error {
    if msg.Owner.Empty() {
        return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, msg.Owner.String())
    }
    if len(msg.AppCode) == 0 {
        return sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "Application Code cannot be empty")
    }
    return nil
}

// GetSignBytes encodes the message for signing
func (msg MsgRecoverApplication) GetSignBytes() []byte {
    return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners defines whose signature is required
func (msg MsgRecoverApplication) GetSigners() []sdk.AccAddress {
    return []sdk.AccAddress{msg.Owner}
}

///////////////////////////
//                       //
// MsgModifyDatabaseUser //
//                       //
///////////////////////////

type MsgModifyDatabaseUser struct {
    Owner sdk.AccAddress `json:"owner"`
    AppCode string       `json:"app_code"`
    Action string        `json:"action"`
    User sdk.AccAddress  `json:"user"`
}

func NewMsgModifyDatabaseUser(owner sdk.AccAddress, appcode, action string, user sdk.AccAddress) MsgModifyDatabaseUser {
    return MsgModifyDatabaseUser {
        Owner: owner,
        AppCode: appcode,
        Action: action,
        User: user,
    }
}

// Route should return the name of the module
func (msg MsgModifyDatabaseUser) Route() string { return RouterKey }

// Type should return the action
func (msg MsgModifyDatabaseUser) Type() string { return "modify_database_user" }

// ValidateBasic runs stateless checks on the message
func (msg MsgModifyDatabaseUser) ValidateBasic() error {
    if msg.Owner.Empty() {
        return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, msg.Owner.String())
    }
    if len(msg.AppCode) == 0 {
        return sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "Application Code cannot be empty")
    }
    if msg.Action != "add" && msg.Action != "drop" {
        return sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "Action has to be either add or drop")
    }
    if msg.User.Empty() {
        return sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "User cannot be empty")
    }
    return nil
}

// GetSignBytes encodes the message for signing
func (msg MsgModifyDatabaseUser) GetSignBytes() []byte {
    return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners defines whose signature is required
func (msg MsgModifyDatabaseUser) GetSigners() []sdk.AccAddress {
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

/////////////////////////////////
//                             //
// NewMsgSetDatabasePermission //
//                             //
/////////////////////////////////

// MsgSetDatabasePermission sets the PermissionRequired of a database
type MsgSetDatabasePermission struct {
    Owner sdk.AccAddress       `json:"owner"`
    AppCode string             `json:"app_code"`
    PermissionRequired string  `json:"permission_required"`
}

// NewMsgSetDatabasePermission is a constructor function for MsgCreatTable
func NewMsgSetDatabasePermission(owner sdk.AccAddress, appCode, permissionRequired string) MsgSetDatabasePermission {
    return MsgSetDatabasePermission {
        Owner: owner,
        AppCode: appCode,
        PermissionRequired: permissionRequired,
    }
}

// Route should return the name of the module
func (msg MsgSetDatabasePermission) Route() string { return RouterKey }

// Type should return the action
func (msg MsgSetDatabasePermission) Type() string { return "set_database_permission" }

// ValidateBasic runs stateless checks on the message
func (msg MsgSetDatabasePermission) ValidateBasic() error {
    if msg.Owner.Empty() {
        return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, msg.Owner.String())
    }
    if len(msg.AppCode) == 0 {
        return sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "Application Code cannot be empty")
    }
    if len(msg.PermissionRequired) == 0 {
        return sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "Permission cannot be empty")
    }
    if msg.PermissionRequired != "required" && msg.PermissionRequired != "unrequired" {
        return sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "Permission has to be either required or unrequired")
    }
    return nil
}

// GetSignBytes encodes the message for signing
func (msg MsgSetDatabasePermission) GetSignBytes() []byte {
    return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners defines whose signature is required
func (msg MsgSetDatabasePermission) GetSigners() []sdk.AccAddress {
    return []sdk.AccAddress{msg.Owner}
}

/////////////////////////////////
//                             //
//    MsgSetFileVolumeLimit    //
//                             //
/////////////////////////////////

type MsgSetAppUserFileVolumeLimit struct {
    Owner sdk.AccAddress       `json:"owner"`
    AppCode string             `json:"app_code"`
    Size    string             `json:"size"`
}

// NewMsgSetDatabasePermission is a constructor function for MsgCreatTable
func NewMsgSetAppUserFileVolumeLimit(owner sdk.AccAddress, appCode, size string) MsgSetAppUserFileVolumeLimit {
    return MsgSetAppUserFileVolumeLimit {
        Owner: owner,
        AppCode: appCode,
        Size: size,
    }
}

// Route should return the name of the module
func (msg MsgSetAppUserFileVolumeLimit) Route() string { return RouterKey }

// Type should return the action
func (msg MsgSetAppUserFileVolumeLimit) Type() string { return "set_app_user_file_volume_limit" }

// ValidateBasic runs stateless checks on the message
func (msg MsgSetAppUserFileVolumeLimit) ValidateBasic() error {
    if msg.Owner.Empty() {
        return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, msg.Owner.String())
    }
    if len(msg.AppCode) == 0 {
        return sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "Application Code cannot be empty")
    }
    if len(msg.Size) == 0 {
        return sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "Size cannot be empty")
    }

    return nil
}

// GetSignBytes encodes the message for signing
func (msg MsgSetAppUserFileVolumeLimit) GetSignBytes() []byte {
    return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners defines whose signature is required
func (msg MsgSetAppUserFileVolumeLimit) GetSigners() []sdk.AccAddress {
    return []sdk.AccAddress{msg.Owner}
}