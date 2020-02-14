package types

import (
    sdk "github.com/cosmos/cosmos-sdk/types"
)

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

////////////////////
//                //
// MsgDropTable //
//                //
////////////////////

// MsgDropTable defines a DropTable message
type MsgDropTable struct {
    Owner sdk.AccAddress `json:"owner"`
    TableName string     `json:"table_name"`
}

// NewMsgDropTable is a constructor function for MsgDropTable
func NewMsgDropTable(owner sdk.AccAddress, tableName string) MsgDropTable {
    return MsgDropTable {
        Owner: owner,
        TableName: tableName,
    }
}

// Route should return the name of the module
func (msg MsgDropTable) Route() string { return RouterKey }

// Type should return the action
func (msg MsgDropTable) Type() string { return "remove_table" }

// ValidateBasic runs stateless checks on the message
func (msg MsgDropTable) ValidateBasic() sdk.Error {
    if msg.Owner.Empty() {
        return sdk.ErrInvalidAddress(msg.Owner.String())
    }
    if len(msg.TableName) == 0 {
        return sdk.ErrUnknownRequest("Table name cannot be empty")
    }
    return nil
}

// GetSignBytes encodes the message for signing
func (msg MsgDropTable) GetSignBytes() []byte {
    return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners defines whose signature is required
func (msg MsgDropTable) GetSigners() []sdk.AccAddress {
    return []sdk.AccAddress{msg.Owner}
}

//////////////////
//              //
// MsgAddColumn //
//              //
//////////////////

type MsgAddColumn struct {
    Owner sdk.AccAddress `json:"owner"`
    TableName string     `json:"table_name"`
    Field string         `json:"field"`
}

// NewMsgCreatePoll is a constructor function for MsgCreatPoll
func NewMsgAddColumn(owner sdk.AccAddress, tableName string, field string) MsgAddColumn {
    return MsgAddColumn {
        Owner: owner,
        TableName: tableName,
        Field: field,
    }
}

// Route should return the name of the module
func (msg MsgAddColumn) Route() string { return RouterKey }

// Type should return the action
func (msg MsgAddColumn) Type() string { return "add_field" }

// ValidateBasic runs stateless checks on the message
func (msg MsgAddColumn) ValidateBasic() sdk.Error {
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
func (msg MsgAddColumn) GetSignBytes() []byte {
    return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners defines whose signature is required
func (msg MsgAddColumn) GetSigners() []sdk.AccAddress {
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

/////////////////////
//                 //
// MsgModifyOption //
//                 //
/////////////////////

type MsgModifyOption struct {
    Owner sdk.AccAddress `json:"owner"`
    TableName string     `json:"table_name"`
    Action string        `json:"action"`
    Option string        `json:"option"`
}

// NewMsgCreatePoll is a constructor function for MsgCreatPoll
func NewMsgModifyOption(owner sdk.AccAddress, tableName string, action string, option string) MsgModifyOption {
    return MsgModifyOption {
        Owner: owner,
        TableName: tableName,
        Action: action,
        Option: option,
    }
}

// Route should return the name of the module
func (msg MsgModifyOption) Route() string { return RouterKey }

// Type should return the action
func (msg MsgModifyOption) Type() string { return "modify_option" }

// ValidateBasic runs stateless checks on the message
func (msg MsgModifyOption) ValidateBasic() sdk.Error {
    if msg.Owner.Empty() {
        return sdk.ErrInvalidAddress(msg.Owner.String())
    }
    if len(msg.TableName) == 0 {
        return sdk.ErrUnknownRequest("Table name cannot be empty")
    }
    if len(msg.Action) ==0 {
        return sdk.ErrUnknownRequest("Action cannot be empty")
    }

    if !(msg.Action == "add" || msg.Action == "remove") {
        return sdk.ErrUnknownRequest("Action has to be either add or remove")
    }

    if len(msg.Option) ==0 {
        return sdk.ErrUnknownRequest("Option cannot be empty")
    }
    return nil
}

// GetSignBytes encodes the message for signing
func (msg MsgModifyOption) GetSignBytes() []byte {
    return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners defines whose signature is required
func (msg MsgModifyOption) GetSigners() []sdk.AccAddress {
    return []sdk.AccAddress{msg.Owner}
}


//////////////////////////
//                      //
// MsgModifyFieldOption //
//                      //
//////////////////////////

type MsgModifyFieldOption struct {
    Owner sdk.AccAddress `json:"owner"`
    TableName string     `json:"table_name"`
    FieldName string     `json:"field_name"`
    Action string        `json:"action"`
    Option string        `json:"option"`
}

// NewMsgCreatePoll is a constructor function for MsgCreatPoll
func NewMsgModifyFieldOption(owner sdk.AccAddress, tableName string, fieldName string, action string, option string) MsgModifyFieldOption {
    return MsgModifyFieldOption {
        Owner: owner,
        TableName: tableName,
        FieldName: fieldName,
        Action: action,
        Option: option,
    }
}

// Route should return the name of the module
func (msg MsgModifyFieldOption) Route() string { return RouterKey }

// Type should return the action
func (msg MsgModifyFieldOption) Type() string { return "modify_field_option" }

// ValidateBasic runs stateless checks on the message
func (msg MsgModifyFieldOption) ValidateBasic() sdk.Error {
    if msg.Owner.Empty() {
        return sdk.ErrInvalidAddress(msg.Owner.String())
    }
    if len(msg.TableName) == 0 {
        return sdk.ErrUnknownRequest("Table name cannot be empty")
    }
    if len(msg.FieldName) == 0 {
        return sdk.ErrUnknownRequest("Field name cannot be empty")
    }
    if len(msg.Action) ==0 {
        return sdk.ErrUnknownRequest("Action cannot be empty")
    }

    if !(msg.Action == "add" || msg.Action == "remove") {
        return sdk.ErrUnknownRequest("Action has to be either add or remove")
    }

    if len(msg.Option) ==0 {
        return sdk.ErrUnknownRequest("Option cannot be empty")
    }
    return nil
}

// GetSignBytes encodes the message for signing
func (msg MsgModifyFieldOption) GetSignBytes() []byte {
    return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners defines whose signature is required
func (msg MsgModifyFieldOption) GetSigners() []sdk.AccAddress {
    return []sdk.AccAddress{msg.Owner}
}

