package oracle

import (
    "encoding/json"
    "fmt"
    //"github.com/cosmos/cosmos-sdk/x/auth/exported"
    authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
    "io/ioutil"
    "net/http"
    "time"
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
      Height string              `json:"height"`
      Result authtypes.AccountI  `json:"result"`
    }

    var account MyAccount
    if err := json.Unmarshal(body, &account); err != nil {
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
        if err != nil || (status == "success" || status == "fail") {
            break
        }
        time.Sleep(1 * time.Second)
    }
}

func CheckTxStatus(accessCode, txHash string) (string, error) {
    type  TxStatus struct {
        State string    `json:"state"`
        Index string    `json:"index"`
        Err  string     `json:"err"`
    }
    type response struct {
        Height string
        Result TxStatus
    }

    var txState response
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
            return txState.Result.State, nil
        }
    }
}
