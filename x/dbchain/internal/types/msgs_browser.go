package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

/////////////////////////
//                     //
//   MsgUpdateTotalTx  //
//                     //
////////////////////////

type MsgUpdateTotalTx struct {
	Owner sdk.AccAddress `json:"owner"`
	Data string  		 `json:"date"`
}

// NewMsgCreatePoll is a constructor function for MsgCreatPoll
func NewMsgUpdateTotalTx(owner sdk.AccAddress,data string) MsgUpdateTotalTx {
	return MsgUpdateTotalTx {
		Owner: owner,
		Data: data,
	}
}

// Route should return the name of the module
func (msg MsgUpdateTotalTx) Route() string { return RouterKey }

// Type should return the action
func (msg MsgUpdateTotalTx) Type() string { return "update_total_tx" }

// ValidateBasic runs stateless checks on the message
func (msg MsgUpdateTotalTx) ValidateBasic() error {
	if msg.Owner.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, msg.Owner.String())
	}
	return nil
}

// GetSignBytes encodes the message for signing
func (msg MsgUpdateTotalTx) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners defines whose signature is required
func (msg MsgUpdateTotalTx) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Owner}
}


/////////////////////////////
//                         //
//   MsgUpdateTxStatistic  //
//                         //
/////////////////////////////

type MsgUpdateTxStatistic struct {
	Owner sdk.AccAddress        `json:"owner"`
	Data  string				`json:"date"`
}

// NewMsgCreatePoll is a constructor function for MsgCreatPoll
func NewMsgUpdateTxStatistic(owner sdk.AccAddress, data string) MsgUpdateTxStatistic {

	return MsgUpdateTxStatistic {
		Owner: owner,
		Data: data,
	}
}

// Route should return the name of the module
func (msg MsgUpdateTxStatistic) Route() string { return RouterKey }

// Type should return the action
func (msg MsgUpdateTxStatistic) Type() string { return "update_total_tx" }

// ValidateBasic runs stateless checks on the message
func (msg MsgUpdateTxStatistic) ValidateBasic() error {
	if msg.Owner.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, msg.Owner.String())
	}
	return nil
}

// GetSignBytes encodes the message for signing
func (msg MsgUpdateTxStatistic) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners defines whose signature is required
func (msg MsgUpdateTxStatistic) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Owner}
}