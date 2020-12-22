package oracle

import (
    "fmt"
    "errors"
    "net/http"
    "io/ioutil"
    "encoding/json"
    "github.com/mr-tron/base58"
    "github.com/yzhanginwa/dbchain/x/dbchain/internal/utils"
)

type SuperQuerierClient struct {
    AppCode string
    Commands [](map[string]string)
}

func NewSuperQuerierClient(appCode string) SuperQuerierClient {
    return SuperQuerierClient{
        AppCode: appCode,
        Commands: [](map[string]string){},
    }
}

func (client *SuperQuerierClient) Table(tableName string) *SuperQuerierClient {
    command := map[string]string{
        "method": "table",
        "table": tableName,
    }
    client.Commands = append(client.Commands, command)
    return client
}

func (client *SuperQuerierClient) Equal(fieldName, value string) *SuperQuerierClient {
    command := map[string]string{
        "method": "equal",
        "field": fieldName,
        "value": value,
    }
    client.Commands = append(client.Commands, command)
    return client
}

func (client *SuperQuerierClient) Last() *SuperQuerierClient {
    command := map[string]string{
        "method": "last",
    }
    client.Commands = append(client.Commands, command)
    return client
}

func (client *SuperQuerierClient) Execute() ([]interface{}, error) {
    accessToken := generateAccessToken()
    appCode := client.AppCode
    commandsJson, err := json.Marshal(client.Commands)
    if err != nil {
        return nil, errors.New("Failed to marshall commands")
    }
    querierBase58 := base58.Encode(commandsJson)

    resp, err := http.Get(fmt.Sprintf("http://localhost:1317/dbchain/querier/%s/%s/%s", accessToken, appCode, querierBase58))
    if err != nil {
        fmt.Println("failed to get account info")
        return nil, err
    }
    defer resp.Body.Close()

    body, _ := ioutil.ReadAll(resp.Body)
    var decodedBody interface{}
    err = json.Unmarshal(body, &decodedBody)
    if err != nil {
        return nil, err
    }
    result := UnWrap(decodedBody, "result")
    return result.([]interface{}), nil
}

//////////////////////
//                  //
// helper functions // 
//                  //
//////////////////////

func generateAccessToken() string {
    privKey, err := LoadPrivKey()
    if err != nil {
        panic("failed to load oracle's private key!!!")
    }
    accessToken := utils.MakeAccessCode(privKey)
    return accessToken
}

