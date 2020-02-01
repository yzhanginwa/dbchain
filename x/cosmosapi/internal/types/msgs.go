package types

import (
    sdk "github.com/cosmos/cosmos-sdk/types"
)

// RouterKey is the module name router key
const RouterKey = ModuleName // this was defined in your key.go file

////////////////////
//                //
// MsgCreateTable //
//                //
////////////////////

// MsgCreateTable defines a CreateTable message
type MsgCreateTable struct {
    Owner sdk.AccAddress `json:"owner"`
    TableName string     `json:"table_name"`
    Fields []string      `json:"fields"`
}

// NewMsgCreateTable is a constructor function for MsgCreatTable
func NewMsgCreateTable(owner sdk.AccAddress, tableName string, fields []string) MsgCreateTable {
    return MsgCreateTable {
        Owner: owner,
        TableName: tableName,
        Fields: fields,
    }
}

// Route should return the name of the module
func (msg MsgCreateTable) Route() string { return RouterKey }

// Type should return the action
func (msg MsgCreateTable) Type() string { return "create_table" }

// ValidateBasic runs stateless checks on the message
func (msg MsgCreateTable) ValidateBasic() sdk.Error {
    if msg.Owner.Empty() {
        return sdk.ErrInvalidAddress(msg.Owner.String())
    }
    if len(msg.TableName) == 0 {
        return sdk.ErrUnknownRequest("Table name cannot be empty")
    }
    if len(msg.Fields) ==0 {
        return sdk.ErrUnknownRequest("Fields cannot be empty")
    }
    return nil
}

// GetSignBytes encodes the message for signing
func (msg MsgCreateTable) GetSignBytes() []byte {
    return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners defines whose signature is required
func (msg MsgCreateTable) GetSigners() []sdk.AccAddress {
    return []sdk.AccAddress{msg.Owner}
}

/////////////////
//             //
// MsgAddField //
//             //
/////////////////

type MsgAddField struct {
    Owner sdk.AccAddress `json:"owner"`
    TableName string     `json:"table_name"`
    Field string         `json:"field"`
}

// NewMsgCreatePoll is a constructor function for MsgCreatPoll
func NewMsgAddField(owner sdk.AccAddress, tableName string, field string) MsgAddField {
    return MsgAddField {
        Owner: owner,
        TableName: tableName,
        Field: field,
    }
}

// Route should return the name of the module
func (msg MsgAddField) Route() string { return RouterKey }

// Type should return the action
func (msg MsgAddField) Type() string { return "add_field" }

// ValidateBasic runs stateless checks on the message
func (msg MsgAddField) ValidateBasic() sdk.Error {
    if msg.Owner.Empty() {
        return sdk.ErrInvalidAddress(msg.Owner.String())
    }
    if len(msg.TableName) == 0 {
        return sdk.ErrUnknownRequest("Table name cannot be empty")
    }
    if len(msg.Field) ==0 {
        return sdk.ErrUnknownRequest("Field cannot be empty")
    }
    return nil
}

// GetSignBytes encodes the message for signing
func (msg MsgAddField) GetSignBytes() []byte {
    return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners defines whose signature is required
func (msg MsgAddField) GetSigners() []sdk.AccAddress {
    return []sdk.AccAddress{msg.Owner}
}

////////////////////
//                //
// MsgRemoveField //
//                //
////////////////////

type MsgRemoveField struct {
    Owner sdk.AccAddress `json:"owner"`
    TableName string     `json:"table_name"`
    Field string         `json:"field"`
}

func NewMsgRemoveField(owner sdk.AccAddress, tableName string, field string) MsgRemoveField {
    return MsgRemoveField {
        Owner: owner,
        TableName: tableName,
        Field: field,
    }
}

// Route should return the name of the module
func (msg MsgRemoveField) Route() string { return RouterKey }

// Type should return the action
func (msg MsgRemoveField) Type() string { return "remove_field" }

// ValidateBasic runs stateless checks on the message
func (msg MsgRemoveField) ValidateBasic() sdk.Error {
    if msg.Owner.Empty() {
        return sdk.ErrInvalidAddress(msg.Owner.String())
    }
    if len(msg.TableName) == 0 {
        return sdk.ErrUnknownRequest("Table name cannot be empty")
    }
    if len(msg.Field) ==0 {
        return sdk.ErrUnknownRequest("Field cannot be empty")
    }
    return nil
}

// GetSignBytes encodes the message for signing
func (msg MsgRemoveField) GetSignBytes() []byte {
    return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners defines whose signature is required
func (msg MsgRemoveField) GetSigners() []sdk.AccAddress {
    return []sdk.AccAddress{msg.Owner}
}

////////////////////
//                //
// MsgRenameField //
//                //
////////////////////

type MsgRenameField struct {
    Owner sdk.AccAddress `json:"owner"`
    TableName string     `json:"table_name"`
    OldField string      `json:"old_field"`
    NewField string      `json:"new_field"`
}

func NewMsgRenameField(owner sdk.AccAddress, tableName string, oldField string, newField string) MsgRenameField {
    return MsgRenameField {
        Owner: owner,
        TableName: tableName,
        OldField: oldField,
        NewField: newField,
    }
}

// Route should return the name of the module
func (msg MsgRenameField) Route() string { return RouterKey }

// Type should return the action
func (msg MsgRenameField) Type() string { return "rename_field" }

// ValidateBasic runs stateless checks on the message
func (msg MsgRenameField) ValidateBasic() sdk.Error {
    if msg.Owner.Empty() {
        return sdk.ErrInvalidAddress(msg.Owner.String())
    }
    if len(msg.TableName) == 0 {
        return sdk.ErrUnknownRequest("Table name cannot be empty")
    }
    if len(msg.OldField) ==0 {
        return sdk.ErrUnknownRequest("Old field cannot be empty")
    }
    if len(msg.NewField) ==0 {
        return sdk.ErrUnknownRequest("New field cannot be empty")
    }
    return nil
}

// GetSignBytes encodes the message for signing
func (msg MsgRenameField) GetSignBytes() []byte {
    return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners defines whose signature is required
func (msg MsgRenameField) GetSigners() []sdk.AccAddress {
    return []sdk.AccAddress{msg.Owner}
}

////////////////////
//                //
// MsgCreateIndex //
//                //
////////////////////

type MsgCreateIndex struct {
    Owner sdk.AccAddress `json:"owner"`
    TableName string     `json:"table_name"`
    Field string         `json:"field"`
}

// NewMsgCreatePoll is a constructor function for MsgCreatPoll
func NewMsgCreateIndex(owner sdk.AccAddress, tableName string, field string) MsgCreateIndex {
    return MsgCreateIndex {
        Owner: owner,
        TableName: tableName,
        Field: field,
    }
}

// Route should return the name of the module
func (msg MsgCreateIndex) Route() string { return RouterKey }

// Type should return the action
func (msg MsgCreateIndex) Type() string { return "create_index" }

// ValidateBasic runs stateless checks on the message
func (msg MsgCreateIndex) ValidateBasic() sdk.Error {
    if msg.Owner.Empty() {
        return sdk.ErrInvalidAddress(msg.Owner.String())
    }
    if len(msg.TableName) == 0 {
        return sdk.ErrUnknownRequest("Table name cannot be empty")
    }
    if len(msg.Field) ==0 {
        return sdk.ErrUnknownRequest("Field cannot be empty")
    }
    return nil
}

// GetSignBytes encodes the message for signing
func (msg MsgCreateIndex) GetSignBytes() []byte {
    return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners defines whose signature is required
func (msg MsgCreateIndex) GetSigners() []sdk.AccAddress {
    return []sdk.AccAddress{msg.Owner}
}

//////////////////
//              //
// MsgInsertRow //
//              //
//////////////////

// MsgCreatePoll defines a CreatePoll message
type MsgInsertRow struct {
    Owner sdk.AccAddress `json:"owner"`
    TableName string     `json:"table_name"`
    Fields RowFieldsJson `json:"fields"`
}

// NewMsgCreatePoll is a constructor function for MsgCreatPoll
func NewMsgInsertRow(owner sdk.AccAddress, tableName string, fieldsJson RowFieldsJson) MsgInsertRow {
    return MsgInsertRow{
        Owner: owner,
        TableName: tableName,
        Fields: fieldsJson,
    }
}

// Route should return the name of the module
func (msg MsgInsertRow) Route() string { return RouterKey }

// Type should return the action
func (msg MsgInsertRow) Type() string { return "insert_row" }

// ValidateBasic runs stateless checks on the message
func (msg MsgInsertRow) ValidateBasic() sdk.Error {
    if msg.Owner.Empty() {
        return sdk.ErrInvalidAddress(msg.Owner.String())
    }
    if len(msg.TableName) == 0 {
        return sdk.ErrUnknownRequest("Table name cannot be empty")
    }
    if len(msg.Fields) ==0 {
        return sdk.ErrUnknownRequest("Fields cannot be empty")
    }
    return nil
}

// GetSignBytes encodes the message for signing
func (msg MsgInsertRow) GetSignBytes() []byte {
    return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners defines whose signature is required
func (msg MsgInsertRow) GetSigners() []sdk.AccAddress {
    return []sdk.AccAddress{msg.Owner}
}

////////////////////////
//                    //
// MsgAddAdminAccount //
//                    //
////////////////////////

type MsgAddAdminAccount struct {
    AdminAddress sdk.AccAddress `json:"admin_address"`
    Owner sdk.AccAddress        `json:"owner"`
}

func NewMsgAddAdminAccount(adminAddress sdk.AccAddress, owner sdk.AccAddress) MsgAddAdminAccount {
    return MsgAddAdminAccount {
        AdminAddress: adminAddress,
        Owner: owner,
    }
}

// Route should return the name of the module
func (msg MsgAddAdminAccount) Route() string { return RouterKey }

// Type should return the action
func (msg MsgAddAdminAccount) Type() string { return "add_admin_account" }

// ValidateBasic runs stateless checks on the message
func (msg MsgAddAdminAccount) ValidateBasic() sdk.Error {
    if msg.AdminAddress.Empty() {
        return sdk.ErrInvalidAddress(msg.AdminAddress.String())
    }
    if msg.Owner.Empty() {
        return sdk.ErrInvalidAddress(msg.Owner.String())
    }
    return nil
}

// GetSignBytes encodes the message for signing
func (msg MsgAddAdminAccount) GetSignBytes() []byte {
    return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners defines whose signature is required
func (msg MsgAddAdminAccount) GetSigners() []sdk.AccAddress {
    return []sdk.AccAddress{msg.Owner}
}

