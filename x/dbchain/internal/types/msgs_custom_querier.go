package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

////////////////////
//                //
// MsgAddfunction //
//                //
////////////////////

type MsgAddCustomQuerier struct {
	Owner sdk.AccAddress `json:"owner"`
	AppCode string       `json:"app_code"`
	QuerierName string  `json:"querier_name"`
	Description string     `json:"description"`
	Body string          `json:"body"`
}

// NewMsgCreatePoll is a constructor function for MsgCreatPoll
func NewMsgAddCustomQuerier(owner sdk.AccAddress, appCode, querierName, description, body string) MsgAddCustomQuerier {
	return MsgAddCustomQuerier {
		Owner: owner,
		AppCode: appCode,
		QuerierName: querierName,
		Description: description,
		Body: body,
	}
}

// Route should return the name of the module
func (msg MsgAddCustomQuerier) Route() string { return RouterKey }

// Type should return the action
func (msg MsgAddCustomQuerier) Type() string { return "add_function" }

// ValidateBasic runs stateless checks on the message
func (msg MsgAddCustomQuerier) ValidateBasic() error {
	if msg.Owner.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, msg.Owner.String())
	}
	if len(msg.AppCode) == 0 {
		return sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "App code cannot be empty")
	}
	if len(msg.Body) ==0 {
		return sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "Body cannot be empty")
	}

	return nil
}

// GetSignBytes encodes the message for signing
func (msg MsgAddCustomQuerier) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners defines whose signature is required
func (msg MsgAddCustomQuerier) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Owner}
}
