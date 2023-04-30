package cli

import (
    "fmt"
    "errors"
    "time"
    "net/http"
    "encoding/json"

)

func waitUntilNextBlock() error {
    currentHeight, err := getCurrentBlockNumber()
    if err != nil {
        return errors.New("Failed to get block height")
    }

    for i := 0; i < 5; i++ {
        time.Sleep(2 * time.Second)
        newHeight, err := getCurrentBlockNumber()
        if err != nil {
            return errors.New("Failed to get block height")
        }
        if newHeight != currentHeight {
            return nil
        }
    }
    return errors.New("Failed to get block height")
}

type ResponseJson struct {
    JsonRpc string           `json:"jsonrpc"`
    Id int                   `json:"id"`
    Result ResultJson        `json:"result"`
}
type ResultJson struct {
    BlockId interface{}      `json:"block_id"`
    Block BlockJson          `json:"block"`
}
type BlockJson struct {
    Header HeaderJson        `json:"header"`
    X map[string]interface{} `json:"-"`
}
type HeaderJson struct {
    Height string            `json:"height"`
    X map[string]interface{} `json:"-"`
}

/*
{
  "jsonrpc": "2.0",
  "id": -1,
  "result": {
    "block_id": {
      "hash": "89E6E7B4D6326F2D2304E37538E9B86961AE95B28E6B5715B14B61483E304DAE",
      ......
    },
    "block": {
      "header": {
        "chain_id": "testnet",
        "height": "2645",
        "time": "2023-04-30T07:21:27.843630866Z",
        ......
        ......
        ......
      },
    }
  }
}
*/

func getCurrentBlockNumber() (string, error) {
    url := "http://localhost:26657/block"
    resp, err := http.Get(url)
    if err != nil {
        fmt.Println("Error:", err)
        return "", err
    }

    defer resp.Body.Close()

    var data ResponseJson
    err = json.NewDecoder(resp.Body).Decode(&data)
    if err != nil {
        fmt.Println("Error:", err)
        return "", err
    }
    return data.Result.Block.Header.Height, nil
}
