package oracle

import (
    "bytes"
    "encoding/hex"
    "encoding/json"
    "fmt"
    "github.com/dbchaincloud/cosmos-sdk/client/context"
    sdk "github.com/dbchaincloud/cosmos-sdk/types"

    "github.com/dbchaincloud/cosmos-sdk/x/auth/client/rest"
    std "github.com/dbchaincloud/cosmos-sdk/x/auth/types"
    "github.com/dbchaincloud/tendermint/crypto/sm2"
    "github.com/mr-tron/base58"
    "github.com/spf13/viper"
    "github.com/yzhanginwa/dbchain/x/dbchain/internal/types"
    "github.com/yzhanginwa/dbchain/x/dbchain/internal/utils"
    "io/ioutil"
    "net/http"
    "time"
)

const (
    BatchSize int = 10
    //BaseUrl = "http://192.168.0.19/relay/"
    BaseUrl = "http://192.168.0.19:3001/relay/"
)

const gasNum = 3000000
type UniversalMsg interface {
    GetSignBytes() []byte
    GetSigners() []sdk.AccAddress
}

var (
    messageChannel = make(chan []UniversalMsg, 1000)
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

func executeTxs(cliCtx context.CLIContext, batch []UniversalMsg, privKey sm2.PrivKeySm2, oracleAccAddr sdk.AccAddress) error {
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

    txBytes, err := buildAndSignAndBuildTxBytes(cliCtx, newBatch, accNum, seq, privKey)
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

func buildAndSignAndBuildTxBytes(cliCtx context.CLIContext, msgs []UniversalMsg, accNum uint64, seq uint64, privKey sm2.PrivKeySm2) ([]byte, error) {
    size := len(msgs)

    needFee , err := setStdFee(cliCtx, "dbchain", size)
    if err != nil {
        return nil, err
    }
    stdFee := NewStdFee(uint64(gasNum * size), needFee)
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

    txBytes, err := genPostTxBytes(cliCtx, stdFee, "",msgs, []StdSignature{stdSignature})
    if err != nil {
        fmt.Println("Oracle: Failed to marshal StdTx!!!")
        return nil, err
    }

    return txBytes, nil
}

func broadcastTxBytes(txBytes []byte) string {

    resp, err := http.Post(BaseUrl + "txs", "application/json", bytes.NewBuffer(txBytes))
    defer resp.Body.Close()
    if err != nil {
        fmt.Println(err)
        return ""
    }
    bz, err  := ioutil.ReadAll(resp.Body)
    if err != nil {
        fmt.Println(err)
        return ""
    }

    type result struct {
        Height string
        Txhash string
        Code int
    }
    temp := result{}
    json.Unmarshal(bz, &temp)
    return temp.Txhash

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

func getCurrentMinGasPrices(cliCtx context.CLIContext, storeName string) (sdk.DecCoins, error){
    ac := generateAccessToken()
    res, err := httpGetRequest(fmt.Sprintf("%s/dbchain/min_gas_prices/%s", BaseUrl, ac))
    if err != nil {
        return nil, err
    }
    type response struct {
        Height string
        Result sdk.DecCoins
    }
    temp := response{}
    err = json.Unmarshal(res, &temp)
    if err != nil {
        return nil, err
    }
    return temp.Result, nil
}

func httpGetRequest(url string) ([]byte, error) {
    resp, err := http.Get(url)
    defer resp.Body.Close()
    if err != nil {
        return nil, err
    }
    bz, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        return nil, err
    }
    return bz, nil
}

func setStdFee(cliCtx context.CLIContext, storeName string, size int)  (sdk.Coins, error) {
    minGasPrices , err := getCurrentMinGasPrices(cliCtx, "dbchain")
    if err != nil {
        return nil, err
    }
    if len(minGasPrices) == 0 {
        return sdk.Coins{sdk.NewCoin("dbctoken", sdk.NewInt(int64(0)))}, nil
    }

    requiredFees := make(sdk.Coins, len(minGasPrices))
    glDec := sdk.NewDec(int64(gasNum * size))
    for i, gp := range minGasPrices {
        fee := gp.Amount.Mul(glDec)
        requiredFees[i] = sdk.NewCoin(gp.Denom, fee.Ceil().RoundInt())
    }
    return requiredFees, nil
}

func genPostTxBytes(cliCtx context.CLIContext, Fee StdFee, Memo string, msgs []UniversalMsg, Signatures []StdSignature) ([]byte, error) {
    stdFee := std.StdFee(Fee)
    stdMsgs := make([]sdk.Msg  , 0)
    for _, m := range msgs {
        stdMsg := m.(sdk.Msg)
        stdMsgs = append(stdMsgs, stdMsg)
    }
    sdkStdSignature := make([]std.StdSignature, 0)
    for _, signature := range Signatures {
        sdkStdSignature = append(sdkStdSignature, std.StdSignature(signature))
    }
    newStdTx := std.NewStdTx(stdMsgs, stdFee, sdkStdSignature, Memo)
    txBytes, err := cliCtx.Codec.MarshalJSON(rest.BroadcastReq{Tx: newStdTx, Mode: "async"})
    return txBytes, err
}