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
    "github.com/yzhanginwa/dbchain/x/dbchain/internal/keeper/import_export"
)

const (
    batchSize = 4
)

func GetCmdImportDatabase(cdc *codec.Codec) *cobra.Command {
    return &cobra.Command{
        Use:   "import-database <app-code> <filename>",
        Short: "import database tables from file",
        Args:  cobra.ExactArgs(2),
        RunE: func(cmd *cobra.Command, args []string) error {
            cliCtx := context.NewCLIContext().WithCodec(cdc)
            inBuf := bufio.NewReader(cmd.InOrStdin())
            txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(cdc))

            appCode := args[0]
            filename := args[1]

            database, err := getDatabaseFromJsonFile(filename)
            if err != nil {
                return err
            }

            msgs, err := databaseToMsgs(cliCtx.GetFromAddress(), appCode, database)
            if err != nil {
                return err
            }

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

func databaseToMsgs(ownerAddr sdk.AccAddress, appCode string, database *import_export.Database) ([]sdk.Msg, error) {
    msgs := make([]sdk.Msg, 0)
    for _, table := range database.Tables {
        err := importTableMsg(ownerAddr, appCode, table, &msgs)
        if err != nil {
            return msgs, err
        }
    }

    for _, fn := range database.CustomFns {
        msgs = append(msgs, types.NewMsgAddFunction(ownerAddr, appCode, fn.Name, fn.Description, fn.Body))
    }

    for _, fn := range database.CustomQueriers {
        msgs = append(msgs, types.NewMsgAddCustomQuerier(ownerAddr, appCode, fn.Name, fn.Description, fn.Body))
    }

    return msgs, nil
}

func importTableMsg(ownerAddr sdk.AccAddress, appCode string, table import_export.Table, msgs *[]sdk.Msg) error {

    var fieldNames []string
    for _, field := range table.Fields {
        fieldNames = append(fieldNames, field.Name)
    }
    msg := types.NewMsgCreateTable(ownerAddr, appCode, table.Name, fieldNames)
    if err := msg.ValidateBasic(); err != nil {
        return err
    }
    *msgs = append(*msgs, msg)

    for _, tableOption := range table.Options {
        *msgs = append(*msgs, types.NewMsgModifyOption(ownerAddr, appCode, table.Name, "add", tableOption))
    }

    if len(table.Filter) > 0 {
        *msgs = append(*msgs, types.NewMsgAddInsertFilter(ownerAddr, appCode, table.Name, table.Filter))
    }

    if len(table.Trigger) > 0 {
        *msgs = append(*msgs, types.NewMsgAddTrigger(ownerAddr, appCode, table.Name, table.Trigger))
    }

    if len(table.Memo) > 0 {
        *msgs = append(*msgs, types.NewMsgSetTableMemo(appCode, table.Name, table.Memo, ownerAddr))
    }

    for _, field := range table.Fields {
        if field.FieldType != "string" {
            *msgs = append(*msgs, types.NewMsgSetColumnDataType(ownerAddr, appCode, table.Name, field.Name, field.FieldType))
        }

        for _, attr := range field.PropertyArr {
            *msgs = append(*msgs, types.NewMsgModifyColumnOption(ownerAddr, appCode, table.Name, field.Name, "add", attr))
        }

        if field.IsIndex {
            *msgs = append(*msgs, types.NewMsgCreateIndex(ownerAddr, appCode, table.Name, field.Name))
        }

        if len(field.Memo) > 0 {
            *msgs = append(*msgs, types.NewMsgSetColumnMemo(appCode, table.Name, field.Name, field.Memo, ownerAddr))
        }
    }

    return nil
}

func getDatabaseFromJsonFile(filename string) (*import_export.Database, error) {
    jsonFile, err := os.Open(filename)
    if err != nil {
        fmt.Println(err)
    }
    defer jsonFile.Close()

    byteValue, _ := ioutil.ReadAll(jsonFile)

    var database import_export.Database
    json.Unmarshal(byteValue, &database)

    return &database, nil
}

