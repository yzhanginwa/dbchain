package oracle

import (
    "fmt"
    "net/http"
    "io/ioutil"
    "github.com/cosmos/cosmos-sdk/x/auth/exported"
)

func getAccountInfo(address string) (uint64, uint64, error) {
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
