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

// MsgCreatePoll defines a CreatePoll message
type MsgCreateTable struct {
        Owner sdk.AccAddress `json:"owner"`
	Name string          `json:"name"`
	Fields []string      `json:"fields"`
}

// NewMsgCreatePoll is a constructor function for MsgCreatPoll
func NewMsgCreateTable(owner sdk.AccAddress, name string, fields []string) MsgCreateTable {
	return MsgCreateTable {
                Owner: owner,
		Name: name,
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
	if len(msg.Name) == 0 {
		return sdk.ErrUnknownRequest("Name cannot be empty")
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

