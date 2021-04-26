package oracle

import (
    "fmt"
    "time"
    "encoding/hex"
    "github.com/spf13/viper"
    rpchttp "github.com/tendermint/tendermint/rpc/client/http"
    "github.com/tendermint/tendermint/crypto/secp256k1"
    sdk "github.com/cosmos/cosmos-sdk/types"
    "github.com/yzhanginwa/dbchain/x/dbchain/internal/utils"
)

const (
    BatchSize int = 10
)

type UniversalMsg interface {
    GetSignBytes() []byte
}

var (
    messageChannel = make(chan []UniversalMsg, 1000)
    runnerIsRunning = false
)

func BuildTxsAndBroadcast(msgs []UniversalMsg) {
    if !runnerIsRunning {
        runnerIsRunning = true
        go txRunner()
    }
    messageChannel <- msgs
}

func txRunner() {
    privKey, err := LoadPrivKey()
    if err != nil {
        fmt.Println("Failed to load oracle's private key!!!")
        runnerIsRunning = false
        return
    }
    oracleAccAddr := sdk.AccAddress(privKey.PubKey().Address())
    queue := []UniversalMsg{}
    hashFlag := make(map[string]bool)

    for {
        select {
        case msgs := <- messageChannel:
            for _, msg := range msgs {
                tmpHash := hex.EncodeToString(msg.GetSignBytes())

                // when using alipay to pay for the package, the client app would ask if the payment is finished.
                // Alipay would reply with outTradeNo and other payment info.
                // meanwhile the alipay notification service would send a notice to oracle to notify the success of a payment.
                // so oracle may generate 2 identical messages and put them into one transaction, which would cause transaction failure.

                if _, ok := hashFlag[tmpHash]; !ok {
                    queue = append(queue, msg)
                    hashFlag[tmpHash] = true
                }
            }
            if len(queue) >= BatchSize {
                err := executeTxs(queue, privKey, oracleAccAddr)
                if err != nil {
                    runnerIsRunning = false
                    return
                }
                queue = []UniversalMsg{}
                hashFlag = make(map[string]bool)
            }
        default:
            if len(queue) > 0 {
                err := executeTxs(queue, privKey, oracleAccAddr)
                if err != nil {
                    runnerIsRunning = false
                    return
                }
                queue = []UniversalMsg{}
                hashFlag = make(map[string]bool)
            } else {
                time.Sleep(2 * time.Second)
            }
        }
    }
}

func executeTxs(batch []UniversalMsg, privKey secp256k1.PrivKeySecp256k1, oracleAccAddr sdk.AccAddress) error {
    accNum, seq, err := GetAccountInfo(oracleAccAddr.String())
    if err != nil {
        fmt.Println("Failed to load oracle's account info!!!")
        return err
    }

    txBytes, err := buildAndSignAndBuildTxBytes(batch, accNum, seq, privKey)
    if err != nil {
        return err
    }
    txHash := broadcastTxBytes(txBytes)
    waitUntilTxFinish(utils.MakeAccessCode(privKey), txHash)
    return nil
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

func broadcastTxBytes(txBytes []byte) string {
    rpc, err := rpchttp.New("http://localhost:26657", "/websocket")
    if err != nil {
        fmt.Printf("failted to get client: %v\n", err)
        return ""
    }

    resp, err := rpc.BroadcastTxAsync(txBytes)
    if err != nil {
        fmt.Printf("failted to broadcast transaction: %v\n", err)
        return ""
    } else {
        return hex.EncodeToString(resp.Hash)
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
