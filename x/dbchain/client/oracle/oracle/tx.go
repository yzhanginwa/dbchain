package oracle

import (
    "fmt"
    "github.com/spf13/viper"
    rpchttp "github.com/tendermint/tendermint/rpc/client/http"
    "github.com/tendermint/tendermint/crypto/secp256k1"
    sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
    BatchSize int = 10
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

    batchMsgs := makeBatches(msgs, BatchSize)
    for _, batch := range batchMsgs {
        txBytes, err := buildAndSignAndBuildTxBytes(batch, accNum, seq, privKey)
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

func buildAndSignAndBuildTxBytes(msgs []UniversalMsg, accNum uint64, seq uint64, privKey secp256k1.PrivKeySecp256k1) ([]byte, error) {
    size := len(msgs)
    stdFee := NewStdFee(uint64(200000 * size), sdk.Coins{sdk.NewCoin("dbctoken", sdk.NewInt(int64(0)))})
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
    rpc, err := rpchttp.New("http://localhost:26657", "/websocket")
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

func makeBatches(msgs []UniversalMsg, batchSize int) [][]UniversalMsg {
    result := [][]UniversalMsg{}
    if len(msgs) == 0 || batchSize < 1 {
        return result
    }

    for len(msgs) >= batchSize  {
        result = append(result, msgs[:batchSize])
        msgs = msgs[batchSize:]
    }

    if len(msgs) > 0 {
        result = append(result, msgs[:])
    }

    return result
}
