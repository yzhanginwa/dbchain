package oracle

import (
    "fmt"
    "encoding/json"
    sdk "github.com/cosmos/cosmos-sdk/types"
    "github.com/yzhanginwa/dbchain/x/dbchain/internal/types"
)

func InsertRows(appCode string, tableName string, rowFieldss []types.RowFields) {
    privKey, err := LoadPrivKey()
    if err != nil {
        fmt.Println("Failed to load oracle's private key!!!")
        return
    }

    oracleAccAddr := sdk.AccAddress(privKey.PubKey().Address())

    msgs := []UniversalMsg{}
    for _, rowFields := range rowFieldss {
        rowFieldsJson, err := json.Marshal(rowFields)
        if err != nil {
            fmt.Println("Oracle: Failed to to json.Marshal!!!")
            return
        }


        msg := types.NewMsgInsertRow(oracleAccAddr, "0000000001", "authentication", rowFieldsJson)
        err = msg.ValidateBasic()
        if err != nil {
            fmt.Println("Oracle: Failed validate new message!!!")
            return
        }

        msgs = append(msgs, msg)
    }
    BuildTxsAndBroadcast(msgs)
}

func InsertRow(appCode string, tableName string, rowFields types.RowFields) {
    rowFieldss := []types.RowFields{rowFields}
    InsertRows(appCode, tableName, rowFieldss)
}
