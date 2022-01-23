package oracle

import (
    "fmt"
    "encoding/json"

    //"github.com/cosmos/cosmos-sdk/client/context"
    "github.com/cosmos/cosmos-sdk/client"

    sdk "github.com/cosmos/cosmos-sdk/types"
    "github.com/yzhanginwa/dbchain/x/dbchain/internal/types"
)

func GetInsertRowMsgs(appCode string, tableName string, rowFieldss []types.RowFields) []UniversalMsg {
    privKey, err := LoadPrivKey()
    if err != nil {
        fmt.Println("Failed to load oracle's private key!!!")
        return []UniversalMsg{}
    }

    oracleAccAddr := sdk.AccAddress(privKey.PubKey().Address())

    msgs := []UniversalMsg{}
    for _, rowFields := range rowFieldss {
        rowFieldsJson, err := json.Marshal(rowFields)
        if err != nil {
            fmt.Println("Oracle: Failed to to json.Marshal!!!")
            return []UniversalMsg{}
        }


        msg := types.NewMsgInsertRow(oracleAccAddr, appCode, tableName, rowFieldsJson)
        err = msg.ValidateBasic()
        if err != nil {
            fmt.Println("Oracle: Failed validate new message!!!")
            return []UniversalMsg{}
        }

        msgs = append(msgs, msg)
    }
    return msgs
}

func InsertRows(cliCtx client.Context, appCode string, tableName string, rowFieldss []types.RowFields) {
    msgs := GetInsertRowMsgs(appCode, tableName, rowFieldss)
    if len(msgs) > 0 {
        BuildTxsAndBroadcast(cliCtx, msgs)
    }
}

func InsertRow(cliCtx client.Context, appCode string, tableName string, rowFields types.RowFields) {
    rowFieldss := []types.RowFields{rowFields}
    InsertRows(cliCtx, appCode, tableName, rowFieldss)
}

func GetOracleAccAddr() sdk.AccAddress{
    privKey, err := LoadPrivKey()
    if err != nil {
        return nil
    }

    oracleAccAddr := sdk.AccAddress(privKey.PubKey().Address())
    return oracleAccAddr
}
