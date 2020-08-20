package oracle

import (
    "fmt"
    "github.com/spf13/viper"
    rpcclient "github.com/tendermint/tendermint/rpc/client"
    "github.com/tendermint/tendermint/crypto/secp256k1"
    sdk "github.com/cosmos/cosmos-sdk/types"
)

type UniversalMsg interface{}

func BuildTxsAndBroadcast(msgs []UniversalMsg) {
    privKey, err := LoadPrivKey()
    if err != nil {
        fmt.Println("Failed to load oracle's private key!!!")
        return
    }
    oracleAccAddr := sdk.AccAddress(privKey.PubKey().Address())
    accNum, seq, err := GetAccountInfo(oracleAccAddr.String())
    if err != nil {
        fmt.Println("Failed to load oracle's account info!!!")
        return
    }

    for _, msg := range msgs {
        txBytes, err := buildAndSignAndBuildTxBytes(msg, accNum, seq, privKey)
        if err != nil {
            return
        }
        broadcastTxBytes(txBytes)
        seq += 1
    }
}

func buildTxAndBroadcast(msg UniversalMsg) {
    msgs := []UniversalMsg{msg}
    BuildTxsAndBroadcast(msgs)
}

func buildAndSignAndBuildTxBytes(msg UniversalMsg, accNum uint64, seq uint64, privKey secp256k1.PrivKeySecp256k1) ([]byte, error) {
    msgs := []UniversalMsg{msg}
    stdFee := NewStdFee(200000, sdk.Coins{sdk.NewCoin("dbctoken", sdk.NewInt(1))})
    chainId := viper.GetString("chain-id")
    stdSignMsgBytes := StdSignBytes(chainId, accNum, seq, stdFee, msgs, "")

    sig, err := privKey.Sign(stdSignMsgBytes)

    if err != nil {
        fmt.Println("Oracle: Failed to sign message!!!")
        return nil, err
    }

    stdSignature := StdSignature {
        PubKey:    privKey.PubKey(),
        Signature: sig,
    }

    newStdTx := NewStdTx(msgs, stdFee, []StdSignature{stdSignature}, "")
    txBytes, err := aminoCdc.MarshalBinaryLengthPrefixed(newStdTx)
    if err != nil {
        fmt.Println("Oracle: Failed to marshal StdTx!!!")
        return nil, err
    }

    return txBytes, nil
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
