package cli

import (
    "bufio"
    "fmt"
    "encoding/json"
    "github.com/spf13/cobra"

    "io/ioutil"
    "os"

    "github.com/cosmos/cosmos-sdk/client/context"
    "github.com/cosmos/cosmos-sdk/codec"
    sdk "github.com/cosmos/cosmos-sdk/types"
    "github.com/cosmos/cosmos-sdk/x/auth"
    "github.com/cosmos/cosmos-sdk/x/auth/client/utils"
    "github.com/yzhanginwa/dbchain/x/dbchain/internal/types"
)

func GetCmdImportTableRows(cdc *codec.Codec) *cobra.Command {
    return &cobra.Command{
        Use:   "import-table-rows <app-code> <table-name> <filename>",
        Short: "import table rows from a json file",
        Args:  cobra.ExactArgs(3),
        RunE: func(cmd *cobra.Command, args []string) error {
            cliCtx := context.NewCLIContext().WithCodec(cdc)
            inBuf := bufio.NewReader(cmd.InOrStdin())
            txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(cdc))

            appCode := args[0]
            tableName := args[1]
            filename := args[2]

            tableRows, err := getTableRowsFromJsonFile(filename)
            if err != nil {
                return err
            }

            msgs, err := tableRowsToMsgs(cliCtx.GetFromAddress(), appCode, tableName, tableRows)
            if err != nil {
                return err
            }

            batchSize := 1             // the default gass can offord only onn insertion
            for len(msgs) > 0 {
                msgBatch := make([]sdk.Msg, 0)
                if len(msgs) > batchSize {
                    msgBatch = msgs[:batchSize]
                    msgs = msgs[batchSize:]
                } else {
                    msgBatch = msgs
                    msgs = make([]sdk.Msg, 0)
                }
                err =  utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, msgBatch)
                if err != nil {
                    return err
                }

                // // We need to wait for the above trasaction to be finished before handling the next one.
                // // Otherwise the next transaction would have the same sequnce number.
                err = waitUntilNextBlock()
                if err != nil {
                        return err
                }
            }
            return nil
        },
    }
}

func tableRowsToMsgs(ownerAddr sdk.AccAddress, appCode string, tableName string, tableRows *[]map[string]string) ([]sdk.Msg, error) {
    msgs := make([]sdk.Msg, 0)
    for _, row := range *tableRows {

        rowFields := make(types.RowFields)
        for f, value := range row {
            if f == "id" || f == "created_by" || f == "created_at" || f== "tx_hash" {
                continue
            }
            rowFields[f] = value
        }
        rowFieldsJson, err := json.Marshal(rowFields)
        if err != nil {
            fmt.Printf("Failed to marshal to json for rowFields")
            return msgs, err
        }

        msg := types.NewMsgInsertRow(ownerAddr, appCode, tableName, rowFieldsJson)
        msgs = append(msgs, msg)
    }   

    return msgs, nil
}

func getTableRowsFromJsonFile(filename string) (*[]map[string]string, error) {
    jsonFile, err := os.Open(filename)
    if err != nil {
        fmt.Println(err)
    }
    defer jsonFile.Close()

    byteValue, _ := ioutil.ReadAll(jsonFile)

    var tableRows []map[string]string
    json.Unmarshal(byteValue, &tableRows)

    return &tableRows, nil
}
