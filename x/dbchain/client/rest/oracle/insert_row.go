package oracle

import (
    "fmt"
    "encoding/json"
    sdk "github.com/cosmos/cosmos-sdk/types"
    "github.com/yzhanginwa/dbchain/x/dbchain/internal/types"
)

func InsertRow(appCode string, tableName string, rowFields types.RowFields) {
    rowFieldsJson, err := json.Marshal(rowFields)
    if err != nil {
        fmt.Println("Oracle: Failed to to json.Marshal!!!")
        return
    }

    privKey, err := LoadPrivKey()
    if err != nil {
        fmt.Println("Failed to load oracle's private key!!!")
        return
    }
    oracleAccAddr := sdk.AccAddress(privKey.PubKey().Address())

    msg := types.NewMsgInsertRow(oracleAccAddr, "0000000001", "authentication", rowFieldsJson)
    err = msg.ValidateBasic()
    if err != nil {
        fmt.Println("Oracle: Failed validate new message!!!")
        return
    }

    buildTxAndBroadcast(msg)
}
