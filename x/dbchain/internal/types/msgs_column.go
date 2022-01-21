package types

import (
    "encoding/base64"
    sdk "github.com/cosmos/cosmos-sdk/types"
    sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

//////////////////
//              //
// MsgAddColumn //
//              //
//////////////////

type MsgAddColumn struct {
    Owner sdk.AccAddress `json:"owner"`
    AppCode string       `json:"app_code"`
    TableName string     `json:"table_name"`
    Field string         `json:"field"`
}

// NewMsgCreatePoll is a constructor function for MsgCreatPoll
func NewMsgAddColumn(owner sdk.AccAddress, appCode string, tableName string, field string) MsgAddColumn {
    return MsgAddColumn {
        Owner: owner,
        AppCode: appCode,
        TableName: tableName,
        Field: field,
    }
}

// Route should return the name of the module
func (msg MsgAddColumn) Route() string { return RouterKey }

// Type should return the action
func (msg MsgAddColumn) Type() string { return "add_column" }

// ValidateBasic runs stateless checks on the message
func (msg MsgAddColumn) ValidateBasic() error {
    if msg.Owner.Empty() {
        return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, msg.Owner.String())
    }
    if len(msg.AppCode) == 0 {
        return sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "App code cannot be empty")
    }
    if len(msg.TableName) == 0 {
        return sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "Table name cannot be empty")
    }
    if len(msg.Field) ==0 {
        return sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "Field cannot be empty")
    }
    if !validateMetaName(msg.Field) {
        return sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "Field name is invalid")
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

///////////////////
//               //
// MsgDropColumn //
//               //
///////////////////

type MsgDropColumn struct {
    Owner sdk.AccAddress `json:"owner"`
    AppCode string       `json:"app_code"`
    TableName string     `json:"table_name"`
    Field string         `json:"field"`
}

func NewMsgDropColumn(owner sdk.AccAddress, appCode string, tableName string, field string) MsgDropColumn {
    return MsgDropColumn {
        Owner: owner,
        AppCode: appCode,
        TableName: tableName,
        Field: field,
    }
}

// Route should return the name of the module
func (msg MsgDropColumn) Route() string { return RouterKey }

// Type should return the action
func (msg MsgDropColumn) Type() string { return "drop_column" }

// ValidateBasic runs stateless checks on the message
func (msg MsgDropColumn) ValidateBasic() error {
    if msg.Owner.Empty() {
        return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, msg.Owner.String())
    }
    if len(msg.AppCode) == 0 {
        return sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "App code cannot be empty")
    }
    if len(msg.TableName) == 0 {
        return sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "Table name cannot be empty")
    }
    if len(msg.Field) ==0 {
        return sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "Field cannot be empty")
    }

    return nil
}

// GetSignBytes encodes the message for signing
func (msg MsgDropColumn) GetSignBytes() []byte {
    return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners defines whose signature is required
func (msg MsgDropColumn) GetSigners() []sdk.AccAddress {
    return []sdk.AccAddress{msg.Owner}
}

/////////////////////
//                 //
// MsgRenameColumn //
//                 //
/////////////////////

type MsgRenameColumn struct {
    Owner sdk.AccAddress `json:"owner"`
    AppCode string       `json:"app_code"`
    TableName string     `json:"table_name"`
    OldField string      `json:"old_field"`
    NewField string      `json:"new_field"`
}

func NewMsgRenameColumn(owner sdk.AccAddress, appCode string, tableName string, oldField string, newField string) MsgRenameColumn {
    return MsgRenameColumn {
        Owner: owner,
        AppCode: appCode,
        TableName: tableName,
        OldField: oldField,
        NewField: newField,
    }
}

// Route should return the name of the module
func (msg MsgRenameColumn) Route() string { return RouterKey }

// Type should return the action
func (msg MsgRenameColumn) Type() string { return "rename_column" }

// ValidateBasic runs stateless checks on the message
func (msg MsgRenameColumn) ValidateBasic() error {
    if msg.Owner.Empty() {
        return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, msg.Owner.String())
    }
    if len(msg.AppCode) == 0 {
        return sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "App code cannot be empty")
    }
    if len(msg.TableName) == 0 {
        return sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "Table name cannot be empty")
    }
    if len(msg.OldField) ==0 {
        return sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "Old field cannot be empty")
    }
    if len(msg.NewField) ==0 {
        return sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "New field cannot be empty")
    }
    if !validateMetaName(msg.NewField) {
        return sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "New field name is invalid")
    }


    return nil
}

// GetSignBytes encodes the message for signing
func (msg MsgRenameColumn) GetSignBytes() []byte {
    return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners defines whose signature is required
func (msg MsgRenameColumn) GetSigners() []sdk.AccAddress {
    return []sdk.AccAddress{msg.Owner}
}

////////////////////
//                //
// MsgCreateIndex //
//                //
////////////////////

type MsgCreateIndex struct {
    Owner sdk.AccAddress `json:"owner"`
    AppCode string       `json:"app_code"`
    TableName string     `json:"table_name"`
    Field string         `json:"field"`
}

// NewMsgCreatePoll is a constructor function for MsgCreatPoll
func NewMsgCreateIndex(owner sdk.AccAddress, appCode string, tableName string, field string) MsgCreateIndex {
    return MsgCreateIndex {
        Owner: owner,
        AppCode: appCode,
        TableName: tableName,
        Field: field,
    }
}

// Route should return the name of the module
func (msg MsgCreateIndex) Route() string { return RouterKey }

// Type should return the action
func (msg MsgCreateIndex) Type() string { return "create_index" }

// ValidateBasic runs stateless checks on the message
func (msg MsgCreateIndex) ValidateBasic() error {
    if msg.Owner.Empty() {
        return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, msg.Owner.String())
    }
    if len(msg.AppCode) == 0 {
        return sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "App code cannot be empty")
    }
    if len(msg.TableName) == 0 {
        return sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "Table name cannot be empty")
    }
    if len(msg.Field) ==0 {
        return sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "Field cannot be empty")
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
// MsgDropIndex //
//              //
//////////////////

type MsgDropIndex struct {
    Owner sdk.AccAddress `json:"owner"`
    AppCode string       `json:"app_code"`
    TableName string     `json:"table_name"`
    Field string         `json:"field"`
}

// NewMsgCreatePoll is a constructor function for MsgCreatPoll
func NewMsgDropIndex(owner sdk.AccAddress, appCode string, tableName string, field string) MsgDropIndex {
    return MsgDropIndex {
        Owner: owner,
        AppCode: appCode,
        TableName: tableName,
        Field: field,
    }
}

// Route should return the name of the module
func (msg MsgDropIndex) Route() string { return RouterKey }

// Type should return the action
func (msg MsgDropIndex) Type() string { return "drop_index" }

// ValidateBasic runs stateless checks on the message
func (msg MsgDropIndex) ValidateBasic() error {
    if msg.Owner.Empty() {
        return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, msg.Owner.String())
    }
    if len(msg.AppCode) == 0 {
        return sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "App code cannot be empty")
    }
    if len(msg.TableName) == 0 {
        return sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "Table name cannot be empty")
    }
    if len(msg.Field) ==0 {
        return sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "Field cannot be empty")
    }

    return nil
}

// GetSignBytes encodes the message for signing
func (msg MsgDropIndex) GetSignBytes() []byte {
    return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners defines whose signature is required
func (msg MsgDropIndex) GetSigners() []sdk.AccAddress {
    return []sdk.AccAddress{msg.Owner}
}

///////////////////////////
//                       //
// MsgModifyColumnOption //
//                       //
///////////////////////////

type MsgModifyColumnOption struct {
    Owner sdk.AccAddress `json:"owner"`
    AppCode string       `json:"app_code"`
    TableName string     `json:"table_name"`
    FieldName string     `json:"field_name"`
    Action string        `json:"action"`
    Option string        `json:"option"`
}

// NewMsgCreatePoll is a constructor function for MsgCreatPoll
func NewMsgModifyColumnOption(owner sdk.AccAddress, appCode string, tableName string, fieldName string, action string, option string) MsgModifyColumnOption {
    return MsgModifyColumnOption {
        Owner: owner,
        AppCode: appCode,
        TableName: tableName,
        FieldName: fieldName,
        Action: action,
        Option: option,
    }
}

// Route should return the name of the module
func (msg MsgModifyColumnOption) Route() string { return RouterKey }

// Type should return the action
func (msg MsgModifyColumnOption) Type() string { return "modify_column_option" }

// ValidateBasic runs stateless checks on the message
func (msg MsgModifyColumnOption) ValidateBasic() error {
    if msg.Owner.Empty() {
        return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, msg.Owner.String())
    }
    if len(msg.AppCode) == 0 {
        return sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "App code cannot be empty")
    }
    if len(msg.TableName) == 0 {
        return sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "Table name cannot be empty")
    }
    if len(msg.FieldName) == 0 {
        return sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "Field name cannot be empty")
    }
    if len(msg.Action) ==0 {
        return sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "Action cannot be empty")
    }

    if !(msg.Action == "add" || msg.Action == "drop") {
        return sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "Action has to be either add or drop")
    }

    if len(msg.Option) ==0 {
        return sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "Option cannot be empty")
    }
    return nil
}

// GetSignBytes encodes the message for signing
func (msg MsgModifyColumnOption) GetSignBytes() []byte {
    return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners defines whose signature is required
func (msg MsgModifyColumnOption) GetSigners() []sdk.AccAddress {
    return []sdk.AccAddress{msg.Owner}
}

///////////////////////////
//                       //
// MsgSetColumnDataType  //
//                       //
///////////////////////////

type MsgSetColumnDataType struct {
    Owner sdk.AccAddress `json:"owner"`
    AppCode string       `json:"app_code"`
    TableName string     `json:"table_name"`
    FieldName string     `json:"field_name"`
    DataType string      `json:"data_type"`
}

func NewMsgSetColumnDataType(owner sdk.AccAddress, appCode string, tableName string, fieldName string, dataType string) MsgSetColumnDataType {
    return MsgSetColumnDataType {
        Owner: owner,
        AppCode: appCode,
        TableName: tableName,
        FieldName: fieldName,
        DataType: dataType,
    }
}

// Route should return the name of the module
func (msg MsgSetColumnDataType) Route() string { return RouterKey }

// Type should return the action
func (msg MsgSetColumnDataType) Type() string { return "set_column_data_type" }

// ValidateBasic runs stateless checks on the message
func (msg MsgSetColumnDataType) ValidateBasic() error {
    if msg.Owner.Empty() {
        return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, msg.Owner.String())
    }
    if len(msg.AppCode) == 0 {
        return sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "App code cannot be empty")
    }
    if len(msg.TableName) == 0 {
        return sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "Table name cannot be empty")
    }
    if len(msg.FieldName) == 0 {
        return sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "Field name cannot be empty")
    }
    if len(msg.DataType) ==0 {
        return sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "DataType cannot be empty")
    }
    return nil
}

// GetSignBytes encodes the message for signing
func (msg MsgSetColumnDataType) GetSignBytes() []byte {
    return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners defines whose signature is required
func (msg MsgSetColumnDataType) GetSigners() []sdk.AccAddress {
    return []sdk.AccAddress{msg.Owner}
}

//////////////////////
//                  //
// MsgSetColumnMemo //
//                  //
//////////////////////


type MsgSetColumnMemo struct {
    AppCode string           `json:"app_code"`
    TableName string         `json:"table_name"`
    FieldName string         `json:"field_name"`
    Memo string              `json:"memo"`
    Owner sdk.AccAddress     `json:"owner"`
}

func NewMsgSetColumnMemo(appCode, tableName, fieldName, memo string, owner sdk.AccAddress) MsgSetColumnMemo {
    return MsgSetColumnMemo {
        AppCode: appCode,
        TableName: tableName,
        FieldName: fieldName,
        Memo: base64.StdEncoding.EncodeToString([]byte(memo)),
        Owner: owner,
    }
}

// Route should return the name of the module
func (msg MsgSetColumnMemo) Route() string { return RouterKey }

// Type should return the action
func (msg MsgSetColumnMemo) Type() string { return "set_column_memo" }

// ValidateBasic runs stateless checks on the message
func (msg MsgSetColumnMemo) ValidateBasic() error {
    if len(msg.AppCode) == 0 {
        return sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "App code cannot be empty")
    }
    if len(msg.TableName) == 0 {
        return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "Table name cannot be empty")
    }
    if len(msg.FieldName) == 0 {
        return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "Column name cannot be empty")
    }
    if msg.Owner.Empty() {
        return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, msg.Owner.String())
    }
    return nil
}

// GetSignBytes encodes the message for signing
func (msg MsgSetColumnMemo) GetSignBytes() []byte {
    return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners defines whose signature is required
func (msg MsgSetColumnMemo) GetSigners() []sdk.AccAddress {
    return []sdk.AccAddress{msg.Owner}
}
