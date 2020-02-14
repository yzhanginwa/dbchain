package types

import (
    sdk "github.com/cosmos/cosmos-sdk/types"
)

////////////////////////
//                    //
// MsgAddAdminAccount //
//                    //
////////////////////////

type MsgAddAdminAccount struct {
    AdminAddress sdk.AccAddress `json:"admin_address"`
    Owner sdk.AccAddress        `json:"owner"`
}

func NewMsgAddAdminAccount(adminAddress sdk.AccAddress, owner sdk.AccAddress) MsgAddAdminAccount {
    return MsgAddAdminAccount {
        AdminAddress: adminAddress,
        Owner: owner,
    }
}

// Route should return the name of the module
func (msg MsgAddAdminAccount) Route() string { return RouterKey }

// Type should return the action
func (msg MsgAddAdminAccount) Type() string { return "add_admin_account" }

// ValidateBasic runs stateless checks on the message
func (msg MsgAddAdminAccount) ValidateBasic() sdk.Error {
    if msg.AdminAddress.Empty() {
        return sdk.ErrInvalidAddress(msg.AdminAddress.String())
    }
    if msg.Owner.Empty() {
        return sdk.ErrInvalidAddress(msg.Owner.String())
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

