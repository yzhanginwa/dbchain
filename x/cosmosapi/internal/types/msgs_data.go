package types

import (
    sdk "github.com/cosmos/cosmos-sdk/types"
)


//////////////////
//              //
// MsgInsertRow //
//              //
//////////////////

// MsgCreatePoll defines a CreatePoll message
type MsgInsertRow struct {
    Owner sdk.AccAddress `json:"owner"`
    AppCode string       `json:"app_code"`
    TableName string     `json:"table_name"`
    Fields RowFieldsJson `json:"fields"`
}

// NewMsgCreatePoll is a constructor function for MsgCreatPoll
func NewMsgInsertRow(owner sdk.AccAddress, appCode string, tableName string, fieldsJson RowFieldsJson) MsgInsertRow {
    return MsgInsertRow{
        Owner: owner,
        AppCode: appCode,
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
    if len(msg.AppCode) == 0 {
        return sdk.ErrUnknownRequest("App code cannot be empty")
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

//////////////////
//              //
// MsgUpdateRow //
//              //
//////////////////

type MsgUpdateRow struct {
    Owner sdk.AccAddress `json:"owner"`
    AppCode string       `json:"app_code"`
    TableName string     `json:"table_name"`
    Id uint              `json:"id"`
    Fields RowFieldsJson `json:"fields"`
}

func NewMsgUpdateRow(owner sdk.AccAddress, appCode string, tableName string, id uint, fieldsJson RowFieldsJson) MsgUpdateRow {
    return MsgUpdateRow{
        Owner: owner,
        AppCode: appCode,
        TableName: tableName,
        Id: id,
        Fields: fieldsJson,
    }
}

// Route should return the name of the module
func (msg MsgUpdateRow) Route() string { return RouterKey }

// Type should return the action
func (msg MsgUpdateRow) Type() string { return "update_row" }

// ValidateBasic runs stateless checks on the message
func (msg MsgUpdateRow) ValidateBasic() sdk.Error {
    if msg.Owner.Empty() {
        return sdk.ErrInvalidAddress(msg.Owner.String())
    }
    if len(msg.AppCode) == 0 {
        return sdk.ErrUnknownRequest("App code cannot be empty")
    }
    if len(msg.TableName) == 0 {
        return sdk.ErrUnknownRequest("Table name cannot be empty")
    }
    if msg.Id ==0 {
        return sdk.ErrUnknownRequest("Id cannot be zero")
    }
    if len(msg.Fields) ==0 {
        return sdk.ErrUnknownRequest("Fields cannot be empty")
    }
    return nil
}

// GetSignBytes encodes the message for signing
func (msg MsgUpdateRow) GetSignBytes() []byte {
    return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners defines whose signature is required
func (msg MsgUpdateRow) GetSigners() []sdk.AccAddress {
    return []sdk.AccAddress{msg.Owner}
}

//////////////////
//              //
// MsgDeleteRow //
//              //
//////////////////

type MsgDeleteRow struct {
    Owner sdk.AccAddress `json:"owner"`
    AppCode string       `json:"app_code"`
    TableName string     `json:"table_name"`
    Id uint              `json:"id"`
}

func NewMsgDeleteRow(owner sdk.AccAddress, appCode string, tableName string, id uint) MsgDeleteRow {
    return MsgDeleteRow{
        Owner: owner,
        AppCode: appCode,
        TableName: tableName,
        Id: id,
    }
}

// Route should return the name of the module
func (msg MsgDeleteRow) Route() string { return RouterKey }

// Type should return the action
func (msg MsgDeleteRow) Type() string { return "delete_row" }

// ValidateBasic runs stateless checks on the message
func (msg MsgDeleteRow) ValidateBasic() sdk.Error {
    if msg.Owner.Empty() {
        return sdk.ErrInvalidAddress(msg.Owner.String())
    }
    if len(msg.AppCode) == 0 {
        return sdk.ErrUnknownRequest("App code cannot be empty")
    }
    if len(msg.TableName) == 0 {
        return sdk.ErrUnknownRequest("Table name cannot be empty")
    }
    if msg.Id ==0 {
        return sdk.ErrUnknownRequest("Id cannot be zero")
    }
    return nil
}

// GetSignBytes encodes the message for signing
func (msg MsgDeleteRow) GetSignBytes() []byte {
    return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners defines whose signature is required
func (msg MsgDeleteRow) GetSigners() []sdk.AccAddress {
    return []sdk.AccAddress{msg.Owner}
}

//////////////////
//              //
// MsgFreezeRow //
//              //
//////////////////

type MsgFreezeRow struct {
    Owner sdk.AccAddress `json:"owner"`
    AppCode string       `json:"app_code"`
    TableName string     `json:"table_name"`
    Id uint              `json:"id"`
}

func NewMsgFreezeRow(owner sdk.AccAddress, appCode string, tableName string, id uint) MsgFreezeRow {
    return MsgFreezeRow{
        Owner: owner,
        AppCode: appCode,
        TableName: tableName,
        Id: id,
    }
}

// Route should return the name of the module
func (msg MsgFreezeRow) Route() string { return RouterKey }

// Type should return the action
func (msg MsgFreezeRow) Type() string { return "freeze_row" }

// ValidateBasic runs stateless checks on the message
func (msg MsgFreezeRow) ValidateBasic() sdk.Error {
    if msg.Owner.Empty() {
        return sdk.ErrInvalidAddress(msg.Owner.String())
    }
    if len(msg.AppCode) == 0 {
        return sdk.ErrUnknownRequest("App code cannot be empty")
    }
    if len(msg.TableName) == 0 {
        return sdk.ErrUnknownRequest("Table name cannot be empty")
    }
    if msg.Id ==0 {
        return sdk.ErrUnknownRequest("Id cannot be zero")
    }
    return nil
}

// GetSignBytes encodes the message for signing
func (msg MsgFreezeRow) GetSignBytes() []byte {
    return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners defines whose signature is required
func (msg MsgFreezeRow) GetSigners() []sdk.AccAddress {
    return []sdk.AccAddress{msg.Owner}
}

