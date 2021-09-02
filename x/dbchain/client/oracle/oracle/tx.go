package oracle

import (
    "encoding/json"
    "fmt"
    "github.com/cosmos/cosmos-sdk/client/context"
    "github.com/mr-tron/base58"
    "github.com/yzhanginwa/dbchain/x/dbchain/internal/types"
    "time"
    "encoding/hex"
    "github.com/spf13/viper"
    rpchttp "github.com/tendermint/tendermint/rpc/client/http"
    "github.com/tendermint/tendermint/crypto/secp256k1"
    sdk "github.com/cosmos/cosmos-sdk/types"
    "github.com/yzhanginwa/dbchain/x/dbchain/internal/utils"
)

const (
    BatchSize int = 100
    StdTxGasNum = 2000000
)

type UniversalMsg interface {
    GetSignBytes() []byte
}

var (
    messageChannel = make(chan []UniversalMsg, 5000)
    runnerIsRunning = false
)

func BuildTxsAndBroadcast(cliCtx context.CLIContext, msgs []UniversalMsg) {
    if !runnerIsRunning {
        runnerIsRunning = true
        go txRunner(cliCtx)
    }
    messageChannel <- msgs
}

func txRunner(cliCtx context.CLIContext) {
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
                key := getMsgKey(msg)

                // when using alipay to pay for the package, the client app would ask if the payment is finished.
                // Alipay would reply with outTradeNo and other payment info.
                // meanwhile the alipay notification service would send a notice to oracle to notify the success of a payment.
                // so oracle may generate 2 identical messages and put them into one transaction, which would cause transaction failure.

                if _, ok := hashFlag[key]; !ok {
                    queue = append(queue, msg)
                    hashFlag[key] = true
                }
            }
            if len(queue) >= BatchSize {
                err := executeTxs(cliCtx, queue, privKey, oracleAccAddr)
                if err != nil {
                    runnerIsRunning = false
                    return
                }
                queue = []UniversalMsg{}
                hashFlag = make(map[string]bool)
            }
        default:
            if len(queue) > 0 {
                err := executeTxs(cliCtx, queue, privKey, oracleAccAddr)
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

func executeTxs(cliCtx context.CLIContext, batch []UniversalMsg, privKey secp256k1.PrivKeySecp256k1, oracleAccAddr sdk.AccAddress) error {
    accNum, seq, err := GetAccountInfo(oracleAccAddr.String())
    if err != nil {
        fmt.Println("Failed to load oracle's account info!!!")
        return err
    }
    newBatch := make([]UniversalMsg, 0)
    for _, msg := range batch {
        if checkCanInsertRow(cliCtx, msg) {
            newBatch = append(newBatch, msg)
        }
    }

    txBytes, err := buildAndSignAndBuildTxBytes(newBatch, accNum, seq, privKey)
    if err != nil {
        return err
    }
    txHash := broadcastTxBytes(txBytes)
    waitUntilTxFinish(utils.MakeAccessCode(privKey), txHash)
    return nil
}

func buildTxAndBroadcast(cliCtx context.CLIContext, msg UniversalMsg) {
    msgs := []UniversalMsg{msg}
    BuildTxsAndBroadcast(cliCtx, msgs)
}

func buildAndSignAndBuildTxBytes(msgs []UniversalMsg, accNum uint64, seq uint64, privKey secp256k1.PrivKeySecp256k1) ([]byte, error) {
    size := len(msgs)
    stdFee := NewStdFee(uint64(StdTxGasNum * size), sdk.Coins{sdk.NewCoin("dbctoken", sdk.NewInt(int64(0)))})
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

func checkCanInsertRow(cliCtx context.CLIContext, msg UniversalMsg) bool {
    insertRow , ok := msg.(types.MsgInsertRow)
    if !ok {
        return true
    }
    rowFieldsJson := base58.Encode(insertRow.Fields)

    privKey, err := LoadPrivKey()
    if err != nil {
        fmt.Println("---------------------> failed msg is  : ", insertRow)
        return false
    }
    ac := utils.MakeAccessCode(privKey)
    res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/can_insert_row/%s/%s/%s/%s", "dbchain", ac, insertRow.AppCode, insertRow.TableName, rowFieldsJson), nil)
    if err != nil {
        fmt.Println("---------------------> msg is  : ", insertRow)
        return false
    }

    result := false
    err = json.Unmarshal(res, &result)
    if err != nil || !result {
        fmt.Println("---------------------> msg is  : ", insertRow)
        return false
    }
    ////////////////////////////////
    fmt.Println("---------------------> can_insert_row result : ", string(res))
    ////////////////////////////////
    return true
}

func getMsgKey(msg UniversalMsg) string {
    insertRow , ok := msg.(types.MsgInsertRow)
    if !ok {
        return hex.EncodeToString(msg.GetSignBytes())
    }
    if insertRow.TableName != "order_receipt" {
        return hex.EncodeToString(msg.GetSignBytes())
    }
    fields := make(map[string]string)
    err := json.Unmarshal(insertRow.Fields, &fields)
    if err != nil {
        return hex.EncodeToString(msg.GetSignBytes())
    }
    vendor_payment_no, ok := fields["vendor_payment_no"]
    if !ok {
        return hex.EncodeToString(msg.GetSignBytes())
    }
    return vendor_payment_no
}