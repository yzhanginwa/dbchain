package types

import (
    "fmt"
    sdk "github.com/cosmos/cosmos-sdk/types"
    sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

////////////////////
//                //
// MsgCreateTable //
//                //
////////////////////

// MsgCreateTable defines a CreateTable message
type MsgCreateTable struct {
    Owner sdk.AccAddress `json:"owner"`
    AppCode string       `json:"app_code"`
    TableName string     `json:"table_name"`
    Fields []string      `json:"fields"`
}

// NewMsgCreateTable is a constructor function for MsgCreatTable
func NewMsgCreateTable(owner sdk.AccAddress, appCode string, tableName string, fields []string) MsgCreateTable {
    return MsgCreateTable {
        Owner: owner,
        AppCode: appCode,
        TableName: tableName,
        Fields: fields,
    }
}

// Route should return the name of the module
func (msg MsgCreateTable) Route() string { return RouterKey }

// Type should return the action
func (msg MsgCreateTable) Type() string { return "create_table" }

// ValidateBasic runs stateless checks on the message
func (msg MsgCreateTable) ValidateBasic() error {
    if msg.Owner.Empty() {
        return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, msg.Owner.String())
    }
    if len(msg.AppCode) == 0 {
        return sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "App code cannot be empty")
    }
    if len(msg.TableName) == 0 {
        return sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "Table name cannot be empty")
    }
    if !validateMetaName(msg.TableName) {
        return sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "Table name is invalid")
    }

    if len(msg.Fields) ==0 {
        return sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "Fields cannot be empty")
    }

    for _, fld := range msg.Fields {
        if !validateMetaName(fld) {
            return sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, fmt.Sprintf("Field name %s is invalid", fld))
        }
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

//////////////////
//              //
// MsgDropTable //
//              //
//////////////////

// MsgDropTable defines a DropTable message
type MsgDropTable struct {
    Owner sdk.AccAddress `json:"owner"`
    AppCode string       `json:"app_code"`
    TableName string     `json:"table_name"`
}

// NewMsgDropTable is a constructor function for MsgDropTable
func NewMsgDropTable(owner sdk.AccAddress, appCode string, tableName string) MsgDropTable {
    return MsgDropTable {
        Owner: owner,
        AppCode: appCode,
        TableName: tableName,
    }
}

// Route should return the name of the module
func (msg MsgDropTable) Route() string { return RouterKey }

// Type should return the action
func (msg MsgDropTable) Type() string { return "drop_table" }

// ValidateBasic runs stateless checks on the message
func (msg MsgDropTable) ValidateBasic() error {
    if msg.Owner.Empty() {
        return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, msg.Owner.String())
    }
    if len(msg.AppCode) == 0 {
        return sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "App code cannot be empty")
    }
    if len(msg.TableName) == 0 {
        return sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "Table name cannot be empty")
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

/////////////////////
//                 //
// MsgModifyOption //
//                 //
/////////////////////

type MsgModifyOption struct {
    Owner sdk.AccAddress `json:"owner"`
    AppCode string       `json:"app_code"`
    TableName string     `json:"table_name"`
    Action string        `json:"action"`
    Option string        `json:"option"`
}

// NewMsgCreatePoll is a constructor function for MsgCreatPoll
func NewMsgModifyOption(owner sdk.AccAddress, appCode string, tableName string, action string, option string) MsgModifyOption {
    return MsgModifyOption {
        Owner: owner,
        AppCode: appCode,
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
func (msg MsgModifyOption) ValidateBasic() error {
    if msg.Owner.Empty() {
        return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, msg.Owner.String())
    }
    if len(msg.AppCode) == 0 {
        return sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "App code cannot be empty")
    }
    if len(msg.TableName) == 0 {
        return sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "Table name cannot be empty")
    }
    if len(msg.Action) == 0 {
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
func (msg MsgModifyOption) GetSignBytes() []byte {
    return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners defines whose signature is required
func (msg MsgModifyOption) GetSigners() []sdk.AccAddress {
    return []sdk.AccAddress{msg.Owner}
}

////////////////////////
//                    //
// MsgAddInsertFilter //
//                    //
////////////////////////

type MsgAddInsertFilter struct {
    Owner sdk.AccAddress `json:"owner"`
    AppCode string       `json:"app_code"`
    TableName string     `json:"table_name"`
    Filter string        `json:"filter"`
}

// NewMsgCreatePoll is a constructor function for MsgCreatPoll
func NewMsgAddInsertFilter(owner sdk.AccAddress, appCode string, tableName string, filter string) MsgAddInsertFilter {
    return MsgAddInsertFilter {
        Owner: owner,
        AppCode: appCode,
        TableName: tableName,
        Filter: filter,
    }
}

// Route should return the name of the module
func (msg MsgAddInsertFilter) Route() string { return RouterKey }

// Type should return the action
func (msg MsgAddInsertFilter) Type() string { return "add_insert_filter" }

// ValidateBasic runs stateless checks on the message
func (msg MsgAddInsertFilter) ValidateBasic() error {
    if msg.Owner.Empty() {
        return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, msg.Owner.String())
    }
    if len(msg.AppCode) == 0 {
        return sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "App code cannot be empty")
    }
    if len(msg.TableName) == 0 {
        return sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "Table name cannot be empty")
    }
    if len(msg.Filter) ==0 {
        return sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "Filter cannot be empty")
    }
    return nil
}

// GetSignBytes encodes the message for signing
func (msg MsgAddInsertFilter) GetSignBytes() []byte {
    return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners defines whose signature is required
func (msg MsgAddInsertFilter) GetSigners() []sdk.AccAddress {
    return []sdk.AccAddress{msg.Owner}
}

///////////////////
//               //
// MsgAddTrigger //
//               //
///////////////////

type MsgAddTrigger struct {
    Owner sdk.AccAddress `json:"owner"`
    AppCode string       `json:"app_code"`
    TableName string     `json:"table_name"`
    Trigger string       `json:"trigger"`
}

// NewMsgCreatePoll is a constructor function for MsgCreatPoll
func NewMsgAddTrigger(owner sdk.AccAddress, appCode string, tableName string, trigger string) MsgAddTrigger {
    return MsgAddTrigger {
        Owner: owner,
        AppCode: appCode,
        TableName: tableName,
        Trigger: trigger,
    }
}

// Route should return the name of the module
func (msg MsgAddTrigger) Route() string { return RouterKey }

// Type should return the action
func (msg MsgAddTrigger) Type() string { return "add_trigger" }

// ValidateBasic runs stateless checks on the message
func (msg MsgAddTrigger) ValidateBasic() error {
    if msg.Owner.Empty() {
        return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, msg.Owner.String())
    }
    if len(msg.AppCode) == 0 {
        return sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "App code cannot be empty")
    }
    if len(msg.TableName) == 0 {
        return sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "Table name cannot be empty")
    }
    if len(msg.Trigger) ==0 {
        return sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "Trigger cannot be empty")
    }
    return nil
}

// GetSignBytes encodes the message for signing
func (msg MsgAddTrigger) GetSignBytes() []byte {
    return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners defines whose signature is required
func (msg MsgAddTrigger) GetSigners() []sdk.AccAddress {
    return []sdk.AccAddress{msg.Owner}
}

///////////////////
//               //
// MsgDropTrigger //
//               //
///////////////////

type MsgDropTrigger struct {
    Owner sdk.AccAddress `json:"owner"`
    AppCode string       `json:"app_code"`
    TableName string     `json:"table_name"`
}

// NewMsgCreatePoll is a constructor function for MsgCreatPoll
func NewMsgDropTrigger(owner sdk.AccAddress, appCode string, tableName string) MsgDropTrigger {
    return MsgDropTrigger {
        Owner: owner,
        AppCode: appCode,
        TableName: tableName,
    }
}

// Route should return the name of the module
func (msg MsgDropTrigger) Route() string { return RouterKey }

// Type should return the action
func (msg MsgDropTrigger) Type() string { return "drop_trigger" }

// ValidateBasic runs stateless checks on the message
func (msg MsgDropTrigger) ValidateBasic() error {
    if msg.Owner.Empty() {
        return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, msg.Owner.String())
    }
    if len(msg.AppCode) == 0 {
        return sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "App code cannot be empty")
    }
    if len(msg.TableName) == 0 {
        return sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "Table name cannot be empty")
    }
    return nil
}

// GetSignBytes encodes the message for signing
func (msg MsgDropTrigger) GetSignBytes() []byte {
    return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners defines whose signature is required
func (msg MsgDropTrigger) GetSigners() []sdk.AccAddress {
    return []sdk.AccAddress{msg.Owner}
}

/////////////////////////
//                     //
// MsgDropInsertFilter //
//                     //
/////////////////////////

type MsgDropInsertFilter struct {
    Owner sdk.AccAddress `json:"owner"`
    AppCode string       `json:"app_code"`
    TableName string     `json:"table_name"`
}

// NewMsgCreatePoll is a constructor function for MsgCreatPoll
func NewMsgDropInsertFilter(owner sdk.AccAddress, appCode string, tableName string) MsgDropInsertFilter {
    return MsgDropInsertFilter {
        Owner: owner,
        AppCode: appCode,
        TableName: tableName,
    }
}

// Route should return the name of the module
func (msg MsgDropInsertFilter) Route() string { return RouterKey }

// Type should return the action
func (msg MsgDropInsertFilter) Type() string { return "drop_insert_filter" }

// ValidateBasic runs stateless checks on the message
func (msg MsgDropInsertFilter) ValidateBasic() error {
    if msg.Owner.Empty() {
        return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, msg.Owner.String())
    }
    if len(msg.AppCode) == 0 {
        return sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "App code cannot be empty")
    }
    if len(msg.TableName) == 0 {
        return sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "Table name cannot be empty")
    }
    return nil
}

// GetSignBytes encodes the message for signing
func (msg MsgDropInsertFilter) GetSignBytes() []byte {
    return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners defines whose signature is required
func (msg MsgDropInsertFilter) GetSigners() []sdk.AccAddress {
    return []sdk.AccAddress{msg.Owner}
}
