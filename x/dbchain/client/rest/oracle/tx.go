package oracle

import (
    "fmt"
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

    stdSignMsg := authtypes.StdSignMsg{
        ChainID:       "testnet",
        AccountNumber: accNum,
        Sequence:      seq,
        Memo:          "",
        Msgs:          msgs,
        Fee:           stdFee,
    }

    sig, err := privKey.Sign(stdSignMsg.Bytes())
    if err != nil {
        fmt.Println("Oracle: Failed to sign message!!!")
        return nil, err
    }

    stdSignature := authtypes.StdSignature {
        PubKey:    privKey.PubKey(),
        Signature: sig,
    }

    newStdTx := authtypes.NewStdTx(msgs, stdFee, []authtypes.StdSignature{stdSignature}, "")

    encoder := authtypes.DefaultTxEncoder(aminoCdc)
    txBytes, err := encoder(newStdTx)
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
