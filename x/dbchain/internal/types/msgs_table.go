package types

import (
    "encoding/base64"
    "fmt"
    sdk "github.com/cosmos/cosmos-sdk/types"
    sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
    "github.com/yzhanginwa/dbchain/x/dbchain/internal/super_script"
    "strings"
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
        Filter: base64.StdEncoding.EncodeToString([]byte(filter)),
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
    //preProcess filter
    p := super_script.NewPreprocessor(strings.NewReader(msg.Filter))
    p.Process()
    if !p.Success {
        sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "syntax error")
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
        Trigger: base64.StdEncoding.EncodeToString([]byte(trigger)),
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
    //preProcess filter
    p := super_script.NewPreprocessor(strings.NewReader(msg.Trigger))
    p.Process()
    if !p.Success {
        return sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "syntax error")
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

/////////////////////
//                 //
// MsgSetTableMemo //
//                 //
/////////////////////

type MsgSetTableMemo struct {
    AppCode string           `json:"app_code"`
    TableName string         `json:"table_name"`
    Memo string              `json:"memo"`
    Owner sdk.AccAddress     `json:"owner"`
}

func NewMsgSetTableMemo(appCode, tableName, memo string, owner sdk.AccAddress) MsgSetTableMemo {
    return MsgSetTableMemo {
        AppCode: appCode,
        TableName: tableName,
        Memo: base64.StdEncoding.EncodeToString([]byte(memo)),
        Owner: owner,
    }
}

// Route should return the name of the module
func (msg MsgSetTableMemo) Route() string { return RouterKey }

// Type should return the action
func (msg MsgSetTableMemo) Type() string { return "set_table_memo" }

// ValidateBasic runs stateless checks on the message
func (msg MsgSetTableMemo) ValidateBasic() error {
    if len(msg.AppCode) == 0 {
        return sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "App code cannot be empty")
    }
    if len(msg.TableName) == 0 {
        return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "Table name cannot be empty")
    }
    if msg.Owner.Empty() {
        return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, msg.Owner.String())
    }
    return nil
}

// GetSignBytes encodes the message for signing
func (msg MsgSetTableMemo) GetSignBytes() []byte {
    return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners defines whose signature is required
func (msg MsgSetTableMemo) GetSigners() []sdk.AccAddress {
    return []sdk.AccAddress{msg.Owner}
}

///////////////////////////////
//                           //
// MsgModifyTableAssociation //
//                           //
///////////////////////////////

type MsgModifyTableAssociation struct {
    AppCode string           `json:"app_code"`
    TableName string         `json:"table_name"`
    AssociationMode string   `json:"association_mode"`
    AssociationTable string  `json:"association_table"`
    Method      string       `json:"method"`
    ForeignKey  string       `json:"foreign_key"`
    Option string            `json:"option"`

    Owner sdk.AccAddress     `json:"owner"`
}

func NewMsgModifyTableAssociation(appCode, tableName , associationMode, associationTable , method, foreignKey , option string, owner sdk.AccAddress, ) MsgModifyTableAssociation {
    return MsgModifyTableAssociation {
        AppCode: appCode,
        TableName: tableName,
        AssociationMode: associationMode,
        AssociationTable: associationTable,
        Method: method,
        ForeignKey: foreignKey,
        Option: option,
        Owner: owner,
    }
}

// Route should return the name of the module
func (msg MsgModifyTableAssociation) Route() string { return RouterKey }

// Type should return the action
func (msg MsgModifyTableAssociation) Type() string { return "set_table_association" }

// ValidateBasic runs stateless checks on the message
func (msg MsgModifyTableAssociation) ValidateBasic() error {
    if len(msg.AppCode) == 0 {
        return sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "App code cannot be empty")
    }
    if len(msg.TableName) == 0 {
        return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "Table name cannot be empty")
    }
    if msg.Option != "add" && msg.Option != "drop" {
        return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "Option only can be set add or drop")
    }

    if msg.AssociationMode != "has_one" && msg.AssociationMode != "has_many" && msg.AssociationMode != "belongs_to" {
        return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "")
    }

    if msg.AssociationTable == "" {
        return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "AssociationTable name cannot be empty")
    }

    if msg.Owner.Empty() {
        return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, msg.Owner.String())
    }
    return nil
}

// GetSignBytes encodes the message for signing
func (msg MsgModifyTableAssociation) GetSignBytes() []byte {
    return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners defines whose signature is required
func (msg MsgModifyTableAssociation) GetSigners() []sdk.AccAddress {
    return []sdk.AccAddress{msg.Owner}
}

///////////////////////////////
//                           //
// MsgEnableCountCache       //
//                           //
///////////////////////////////

type MsgAddCounterCache struct {
    AppCode string           `json:"app_code"`
    TableName string         `json:"table_name"`
    AssociationTable string  `json:"association_table"`
    ForeignKey  string       `json:"foreign_key"`
    CounterCacheField  string       `json:"counter_cache_field"`
    Limit       string       `json:"limit"`

    Owner sdk.AccAddress     `json:"owner"`
}

func NewMsgAddCounterCache(appCode, tableName, associationTable , foreignKey , counterCacheField, limit string, owner sdk.AccAddress, ) MsgAddCounterCache {
    return MsgAddCounterCache {
        AppCode: appCode,
        TableName: tableName,
        AssociationTable: associationTable,
        ForeignKey: foreignKey,
        CounterCacheField: counterCacheField,
        Limit: limit,
        Owner: owner,
    }
}

// Route should return the name of the module
func (msg MsgAddCounterCache) Route() string { return RouterKey }

// Type should return the action
func (msg MsgAddCounterCache) Type() string { return "enable_count_cache" }

// ValidateBasic runs stateless checks on the message
func (msg MsgAddCounterCache) ValidateBasic() error {
    if len(msg.AppCode) == 0 {
        return sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "App code cannot be empty")
    }
    if len(msg.TableName) == 0 {
        return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "Table name cannot be empty")
    }

    if msg.AssociationTable == "" {
        return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "AssociationTable name cannot be empty")
    }

    if msg.ForeignKey == "" {
        return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "ForeignKey cannot be empty")
    }

    if msg.CounterCacheField == "" {
        return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "CountCacheField name cannot be empty")
    }


    if msg.Owner.Empty() {
        return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, msg.Owner.String())
    }
    return nil
}

// GetSignBytes encodes the message for signing
func (msg MsgAddCounterCache) GetSignBytes() []byte {
    return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners defines whose signature is required
func (msg MsgAddCounterCache) GetSigners() []sdk.AccAddress {
    return []sdk.AccAddress{msg.Owner}
}