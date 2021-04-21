package oracle

import (
    "fmt"
    "time"
    "encoding/json"
    "net/http"
    "io/ioutil"
    "github.com/yzhanginwa/dbchain/x/dbchain/internal/types"
    "github.com/cosmos/cosmos-sdk/x/auth/exported"
)

func GetAccountInfo(address string) (uint64, uint64, error) {
    resp, err := http.Get(fmt.Sprintf("http://localhost:1317/auth/accounts/%s", address))
    if err != nil {
        fmt.Println("failed to get account info")
        return 0, 0, err
    }
    defer resp.Body.Close()
    body, err := ioutil.ReadAll(resp.Body)

    type MyAccount struct {
      Height string            `json:"height"`
      Result exported.Account  `json:"result"`
    }

    var account MyAccount
    if err := aminoCdc.UnmarshalJSON(body, &account); err != nil {
        fmt.Printf("failted to broadcast unmarshal account body\n")
        return 0, 0, err
    }

    seq := account.Result.GetSequence()
    accountNumber := account.Result.GetAccountNumber()

    return accountNumber, seq, nil
}

func waitUntilTxFinish(accessCode, txHash string) {
    for count := 10; count > 0; count-- {
        status, err := CheckTxStatus(accessCode, txHash)
        if err != nil || status != "processing" {
            break
        }
        time.Sleep(1 * time.Second)
    }
}

func CheckTxStatus(accessCode, txHash string) (string, error) {
    var txState types.TxStatus
    resp, err := http.Get(fmt.Sprintf("http://localhost:1317/dbchain/tx-simple-result/%s/%s", accessCode, txHash))
    if err != nil {
        return "", err
    } else {
        defer resp.Body.Close()
        body, err := ioutil.ReadAll(resp.Body)
        if err != nil {
            return "", err
        }
        if err = json.Unmarshal(body, &txState); err != nil {
            return "", err
        } else {
            return txState.State, nil
        }
    }
}
