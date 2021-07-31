package oracle

import (
    "fmt"
    "errors"
    sdk "github.com/dbchaincloud/cosmos-sdk/types"
    sdkerrors "github.com/dbchaincloud/cosmos-sdk/types/errors"
    dtypes "github.com/yzhanginwa/dbchain/x/dbchain/internal/types"
)

// MsgSend - high level transaction of the coin module
type MsgSend struct {
    FromAddress sdk.AccAddress `json:"from_address" yaml:"from_address"`
    ToAddress   sdk.AccAddress `json:"to_address" yaml:"to_address"`
    Amount      sdk.Coins      `json:"amount" yaml:"amount"`
}

func NewMsgSend(fromAddr, toAddr sdk.AccAddress, amount sdk.Coins) MsgSend {
    return MsgSend{FromAddress: fromAddr, ToAddress: toAddr, Amount: amount}
}

func (msg MsgSend) GetSignBytes() []byte {
    return sdk.MustSortJSON(aminoCdc.MustMarshalJSON(msg))
}

// Route Implements Msg.
func (msg MsgSend) Route() string { return dtypes.RouterKey }

// Type Implements Msg.
func (msg MsgSend) Type() string { return "send" }

// ValidateBasic Implements Msg.
func (msg MsgSend) ValidateBasic() error {
    if msg.FromAddress.Empty() {
        return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "missing sender address")
    }
    if msg.ToAddress.Empty() {
        return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "missing recipient address")
    }
    if !msg.Amount.IsValid() {
        return sdkerrors.Wrap(sdkerrors.ErrInvalidCoins, msg.Amount.String())
    }
    if !msg.Amount.IsAllPositive() {
        return sdkerrors.Wrap(sdkerrors.ErrInvalidCoins, msg.Amount.String())
    }
    return nil
}

// GetSigners Implements Msg.
func (msg MsgSend) GetSigners() []sdk.AccAddress {
    return []sdk.AccAddress{msg.FromAddress}
}

func GetSendTokenMsg(addr sdk.AccAddress) (UniversalMsg, error) {
    privKey, err := LoadPrivKey()
    if err != nil {
        fmt.Println("Failed to load oracle's private key!!!")
        return nil, errors.New("Failed to load oracle's private key!!!")
    }
    oracleAccAddr := sdk.AccAddress(privKey.PubKey().Address())

    oneCoin := sdk.NewCoin("dbctoken", sdk.NewInt(1))
    msg := NewMsgSend(oracleAccAddr, addr, []sdk.Coin{oneCoin})
    return msg, nil
}
