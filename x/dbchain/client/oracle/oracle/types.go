package oracle

import (
    "encoding/json"
    "github.com/dbchaincloud/tendermint/crypto"
    sdk "github.com/dbchaincloud/cosmos-sdk/types"
)

////////////
//        // 
// StdFee //
//        // 
////////////

type StdFee struct {
    Amount sdk.Coins `json:"amount" yaml:"amount"`
    Gas    uint64    `json:"gas" yaml:"gas"`
}

func NewStdFee(gas uint64, amount sdk.Coins) StdFee {
    return StdFee{
        Amount: amount,
        Gas:    gas,
    }
}

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

//////////////////
//              //
// StdSignature //
//              //
//////////////////

type StdSignature struct {
    crypto.PubKey `json:"pub_key" yaml:"pub_key"` // optional
    Signature     []byte                          `json:"signature" yaml:"signature"`
}

///////////
//       //
// StdTx //
//       //
///////////

type StdTx struct {
    Msgs       []UniversalMsg `json:"msg" yaml:"msg"`
    Fee        StdFee         `json:"fee" yaml:"fee"`
    Signatures []StdSignature `json:"signatures" yaml:"signatures"`
    Memo       string         `json:"memo" yaml:"memo"`
}

func NewStdTx(msgs []UniversalMsg, fee StdFee, sigs []StdSignature, memo string) StdTx {
    return StdTx{
        Msgs:       msgs,
        Fee:        fee,
        Signatures: sigs,
        Memo:       memo,
    }
}

////////////////
//            //
// StdSignDoc //
//            //
////////////////

type StdSignDoc struct {
    AccountNumber uint64            `json:"account_number" yaml:"account_number"`
    ChainID       string            `json:"chain_id" yaml:"chain_id"`
    Fee           json.RawMessage   `json:"fee" yaml:"fee"`
    Memo          string            `json:"memo" yaml:"memo"`
    Msgs          []json.RawMessage `json:"msgs" yaml:"msgs"`
    Sequence      uint64            `json:"sequence" yaml:"sequence"`
}
