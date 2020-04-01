package types

import (
    sdk "github.com/cosmos/cosmos-sdk/types"
    sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

////////////////////////
//                    //
// MsgAddAdminAccount //
//                    //
////////////////////////

type MsgAddAdminAccount struct {
    AppCode string              `json:"app_code"`
    AdminAddress sdk.AccAddress `json:"admin_address"`
    Owner sdk.AccAddress        `json:"owner"`
}

func NewMsgAddAdminAccount(appCode string, adminAddress sdk.AccAddress, owner sdk.AccAddress) MsgAddAdminAccount {
    return MsgAddAdminAccount {
        AppCode: appCode,
        AdminAddress: adminAddress,
        Owner: owner,
    }
}

// Route should return the name of the module
func (msg MsgAddAdminAccount) Route() string { return RouterKey }

// Type should return the action
func (msg MsgAddAdminAccount) Type() string { return "add_admin_account" }

// ValidateBasic runs stateless checks on the message
func (msg MsgAddAdminAccount) ValidateBasic() error {
    if len(msg.AppCode) == 0 {
        return sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "App code cannot be empty")
    }
    if msg.AdminAddress.Empty() {
        return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, msg.AdminAddress.String())
    }
    if msg.Owner.Empty() {
        return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, msg.Owner.String())
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

