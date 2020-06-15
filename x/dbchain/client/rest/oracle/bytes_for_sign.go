package oracle

import (
    "encoding/json"
    "github.com/tendermint/tendermint/crypto"
    sdk "github.com/cosmos/cosmos-sdk/types"
)

type StdSignDoc struct {
    AccountNumber uint64            `json:"account_number" yaml:"account_number"`
    ChainID       string            `json:"chain_id" yaml:"chain_id"`
    Fee           json.RawMessage   `json:"fee" yaml:"fee"`
    Memo          string            `json:"memo" yaml:"memo"`
    Msgs          []json.RawMessage `json:"msgs" yaml:"msgs"`
    Sequence      uint64            `json:"sequence" yaml:"sequence"`
}

type StdFee struct {
    Amount sdk.Coins `json:"amount" yaml:"amount"`
    Gas    uint64    `json:"gas" yaml:"gas"`
}

type StdSignature struct {
        crypto.PubKey `json:"pub_key" yaml:"pub_key"` // optional
        Signature     []byte                          `json:"signature" yaml:"signature"`
}

type StdTx struct {
    Msgs       []sdk.Msg      `json:"msg" yaml:"msg"`
    Fee        StdFee         `json:"fee" yaml:"fee"`
    Signatures []StdSignature `json:"signatures" yaml:"signatures"`
    Memo       string         `json:"memo" yaml:"memo"`
}

func init () {
    //aminoCdc.RegisterInterface((*sdk.Msg)(nil), nil)
    //aminoCdc.RegisterInterface((*sdk.Tx)(nil), nil)
    //aminoCdc.RegisterConcrete(types.MsgInsertRow{}, "dbchain/InsertRow", nil)
    //cryptoamino.RegisterAmino(aminoCdc)
    //authtypes.RegisterCodec(aminoCdc)
}

func StdSignBytes(chainID string, accnum uint64, sequence uint64, fee StdFee, msgs []sdk.Msg, memo string) []byte {
    msgsBytes := make([]json.RawMessage, 0, len(msgs))
    for _, msg := range msgs {
        msgsBytes = append(msgsBytes, json.RawMessage(GetSignBytes(msg)))
    }
    bz, err := aminoCdc.MarshalJSON(StdSignDoc{
        AccountNumber: accnum,
        ChainID:       chainID,
        Fee:           json.RawMessage(fee.Bytes()),
        Memo:          memo,
        Msgs:          msgsBytes,
        Sequence:      sequence,
    })
    if err != nil {
        panic(err)
    }
    return sdk.MustSortJSON(bz)
}

func GetSignBytes(msg interface{}) []byte {
    //func MustMarshalJSON(x interface{}) []byte
    return sdk.MustSortJSON(aminoCdc.MustMarshalJSON(msg))
}

// NewStdFee returns a new instance of StdFee
func NewStdFee(gas uint64, amount sdk.Coins) StdFee {
        return StdFee{
                Amount: amount,
                Gas:    gas,
        }
}

// Bytes for signing later
func (fee StdFee) Bytes() []byte {
        // normalize. XXX
        // this is a sign of something ugly
        // (in the lcd_test, client side its null,
        // server side its [])
        if len(fee.Amount) == 0 {
                fee.Amount = sdk.NewCoins()
        }
        bz, err := aminoCdc.MarshalJSON(fee) // TODO
        if err != nil {
                panic(err)
        }
        return bz
}

func NewStdTx(msgs []sdk.Msg, fee StdFee, sigs []StdSignature, memo string) StdTx {
    return StdTx{
        Msgs:       msgs,
        Fee:        fee,
        Signatures: sigs,
        Memo:       memo,
    }
}
