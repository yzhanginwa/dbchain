package oracle

import (
    "fmt"
    "github.com/spf13/viper"
    amino "github.com/tendermint/go-amino"
    cryptoamino "github.com/tendermint/tendermint/crypto/encoding/amino"
    rpcclient "github.com/tendermint/tendermint/rpc/client"
    "github.com/tendermint/tendermint/crypto/secp256k1"
    sdk "github.com/cosmos/cosmos-sdk/types"
    authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
    "github.com/yzhanginwa/dbchain/x/dbchain/internal/types"
)

var (
    aminoCdc = amino.NewCodec()
)

func init () {
    aminoCdc.RegisterInterface((*sdk.Msg)(nil), nil)
    aminoCdc.RegisterInterface((*sdk.Tx)(nil), nil)
    aminoCdc.RegisterConcrete(types.MsgInsertRow{}, "dbchain/InsertRow", nil)
    cryptoamino.RegisterAmino(aminoCdc)
    authtypes.RegisterCodec(aminoCdc)
}

func buildTxAndBroadcast(msg sdk.Msg) {
    privKey, err := LoadPrivKey()
    if err != nil {
        fmt.Println("Failed to load oracle's private key!!!")
        return
    }
    oracleAccAddr := sdk.AccAddress(privKey.PubKey().Address())
    accNum, seq, err := getAccountInfo(oracleAccAddr.String())
    if err != nil {
        fmt.Println("Failed to load oracle's account info!!!")
        return
    }

    txBytes, err := buildAndSignAndBuildTxBytes(msg, accNum, seq, privKey)
    if err != nil {
        return
    }
    broadcastTxBytes(txBytes)
}

func buildAndSignAndBuildTxBytes(msg sdk.Msg, accNum uint64, seq uint64, privKey secp256k1.PrivKeySecp256k1) ([]byte, error) {
    msgs := []sdk.Msg{msg}
    stdFee := authtypes.NewStdFee(200000, sdk.Coins{sdk.NewCoin("dbctoken", sdk.NewInt(1))})
    chainId := viper.GetString("chain-id")
    stdSignMsgBytes := authtypes.StdSignBytes(chainId, accNum, seq, stdFee, msgs, "")

    //type StdSignDoc struct {
    //        AccountNumber uint64            `json:"account_number" yaml:"account_number"`
    //        ChainID       string            `json:"chain_id" yaml:"chain_id"`
    //        Fee           json.RawMessage   `json:"fee" yaml:"fee"`
    //        Memo          string            `json:"memo" yaml:"memo"`
    //        Msgs          []json.RawMessage `json:"msgs" yaml:"msgs"`
    //        Sequence      uint64            `json:"sequence" yaml:"sequence"`
    //}

    //func StdSignBytes(chainID string, accnum uint64, sequence uint64, fee StdFee, msgs []sdk.Msg, memo string) []byte {
    //        msgsBytes := make([]json.RawMessage, 0, len(msgs))
    //        for _, msg := range msgs {
    //                msgsBytes = append(msgsBytes, json.RawMessage(msg.GetSignBytes()))
    //        }
    //        bz, err := ModuleCdc.MarshalJSON(StdSignDoc{
    //                AccountNumber: accnum,
    //                ChainID:       chainID,
    //                Fee:           json.RawMessage(fee.Bytes()),
    //                Memo:          memo,
    //                Msgs:          msgsBytes,
    //                Sequence:      sequence,
    //        })
    //        if err != nil {
    //                panic(err)
    //        }
    //        return sdk.MustSortJSON(bz)
    //}


    //func (msg MsgInsertRow) GetSignBytes() []byte {
    //    return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
    //}

    //func MustMarshalJSON(x interface{}) []byte


    //type StdFee struct {
    //        Amount sdk.Coins `json:"amount" yaml:"amount"`
    //        Gas    uint64    `json:"gas" yaml:"gas"`
    //}
    //
    //// NewStdFee returns a new instance of StdFee
    //func NewStdFee(gas uint64, amount sdk.Coins) StdFee {
    //        return StdFee{
    //                Amount: amount,
    //                Gas:    gas,
    //        }
    //}
    //
    //// Bytes for signing later
    //func (fee StdFee) Bytes() []byte {
    //        // normalize. XXX
    //        // this is a sign of something ugly
    //        // (in the lcd_test, client side its null,
    //        // server side its [])
    //        if len(fee.Amount) == 0 {
    //                fee.Amount = sdk.NewCoins()
    //        }
    //        bz, err := ModuleCdc.MarshalJSON(fee) // TODO
    //        if err != nil {
    //                panic(err)
    //        }
    //        return bz
    //}


    sig, err := privKey.Sign(stdSignMsgBytes)

    if err != nil {
        fmt.Println("Oracle: Failed to sign message!!!")
        return nil, err
    }

    stdSignature := authtypes.StdSignature {
        PubKey:    privKey.PubKey(),
        Signature: sig,
    }

    newStdTx := authtypes.NewStdTx(msgs, stdFee, []authtypes.StdSignature{stdSignature}, "")


    //type StdTx struct {
    //        Msgs       []sdk.Msg      `json:"msg" yaml:"msg"`
    //        Fee        StdFee         `json:"fee" yaml:"fee"`
    //        Signatures []StdSignature `json:"signatures" yaml:"signatures"`
    //        Memo       string         `json:"memo" yaml:"memo"`
    //}
    //
    //func NewStdTx(msgs []sdk.Msg, fee StdFee, sigs []StdSignature, memo string) StdTx {
    //        return StdTx{
    //                Msgs:       msgs,
    //                Fee:        fee,
    //                Signatures: sigs,
    //                Memo:       memo,
    //        }
    //}


    txBytes, err := aminoCdc.MarshalBinaryLengthPrefixed(newStdTx)
    if err != nil {
        fmt.Println("Oracle: Failed to marshal StdTx!!!")
        return nil, err
    }

    return txBytes, nil

    //cliCtx.BroadcastTxAsync(txBytes)
}

func broadcastTxBytes(txBytes []byte) {
    rpc, err := rpcclient.NewHTTP("http://localhost:26657", "/websocket")
    if err != nil {
        fmt.Printf("failted to get client: %v\n", err)
        return
    }

    _, err = rpc.BroadcastTxAsync(txBytes)
    if err != nil {
        fmt.Printf("failted to broadcast transaction: %v\n", err)
        return
    }
}
