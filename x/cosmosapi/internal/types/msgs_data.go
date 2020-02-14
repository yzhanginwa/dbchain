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

