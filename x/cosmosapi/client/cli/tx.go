package cli

import (
    "fmt"
    "errors"
    "strings"
    "strconv"
    "encoding/json"
    "github.com/spf13/cobra"

    "github.com/cosmos/cosmos-sdk/client"
    "github.com/cosmos/cosmos-sdk/client/context"
    "github.com/cosmos/cosmos-sdk/codec"
    sdk "github.com/cosmos/cosmos-sdk/types"
    "github.com/cosmos/cosmos-sdk/x/auth"
    "github.com/cosmos/cosmos-sdk/x/auth/client/utils"
    "github.com/yzhanginwa/cosmos-api/x/cosmosapi/internal/types"
)

func GetTxCmd(storeKey string, cdc *codec.Codec) *cobra.Command {
    cosmosapiTxCmd := &cobra.Command{
        Use:                        types.ModuleName,
        Short:                      "Cosmosapi transaction subcommands",
        DisableFlagParsing:         true,
        SuggestionsMinimumDistance: 2,
        RunE:                       client.ValidateCmd,
    }

    cosmosapiTxCmd.AddCommand(client.PostCommands(
        GetCmdCreateTable(cdc),
        GetCmdDropTable(cdc),
        GetCmdAddColumn(cdc),
        GetCmdDropColumn(cdc),
        GetCmdRenameField(cdc),
        GetCmdCreateIndex(cdc),
        GetCmdModifyOption(cdc),
        GetCmdModifyFieldOption(cdc),
        GetCmdInsertRow(cdc),
        GetCmdUpdateRow(cdc),
        GetCmdDeleteRow(cdc),
        GetCmdAddAdminAccount(cdc),
    )...)

    return cosmosapiTxCmd
}

////////////////////
//                //
// schema related //
//                //
////////////////////

// GetCmdCreatePoll is the CLI command for sending a CreatePoll transaction
func GetCmdCreateTable(cdc *codec.Codec) *cobra.Command {
    return &cobra.Command{
        Use:   "create-table [name] [fields]",
        Short: "create a new table",
        Args:  cobra.ExactArgs(2),
        RunE: func(cmd *cobra.Command, args []string) error {
            cliCtx := context.NewCLIContext().WithCodec(cdc)
            txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))

            name := args[0]
            fields := strings.Split(args[1], ",")
            msg := types.NewMsgCreateTable(cliCtx.GetFromAddress(), name, fields)
            err := msg.ValidateBasic()
            if err != nil {
                return err
            }

            return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
        },
    }
}

func GetCmdDropTable(cdc *codec.Codec) *cobra.Command {
    return &cobra.Command{
        Use:   "drop-table [name]",
        Short: "drop a table",
        Args:  cobra.ExactArgs(1),
        RunE: func(cmd *cobra.Command, args []string) error {
            cliCtx := context.NewCLIContext().WithCodec(cdc)
            txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))

            name := args[0]
            msg := types.NewMsgDropTable(cliCtx.GetFromAddress(), name)
            err := msg.ValidateBasic()
            if err != nil {
                return err
            }

            return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
        },
    }
}

func GetCmdAddColumn(cdc *codec.Codec) *cobra.Command {
    return &cobra.Command{
        Use:   "add-column [name] [field]",
        Short: "add a new column onto a table",
        Args:  cobra.ExactArgs(2),
        RunE: func(cmd *cobra.Command, args []string) error {
            cliCtx := context.NewCLIContext().WithCodec(cdc)
            txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))

            name := args[0]
            field := args[1]
            msg := types.NewMsgAddColumn(cliCtx.GetFromAddress(), name, field)
            err := msg.ValidateBasic()
            if err != nil {
                return err
            }

            return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
        },
    }
}

func GetCmdDropColumn(cdc *codec.Codec) *cobra.Command {
    return &cobra.Command{
        Use:   "drop-column [name] [field]",
        Short: "drop a column from a table",
        Args:  cobra.ExactArgs(2),
        RunE: func(cmd *cobra.Command, args []string) error {
            cliCtx := context.NewCLIContext().WithCodec(cdc)
            txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))

            name := args[0]
            field := args[1]
            msg := types.NewMsgDropColumn(cliCtx.GetFromAddress(), name, field)
            err := msg.ValidateBasic()
            if err != nil {
                return err
            }

            return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
        },
    }
}

func GetCmdRenameField(cdc *codec.Codec) *cobra.Command {
    return &cobra.Command{
        Use:   "rename-field [name] [old-field] [new-field",
        Short: "rename a field in a table",
        Args:  cobra.ExactArgs(3),
        RunE: func(cmd *cobra.Command, args []string) error {
            cliCtx := context.NewCLIContext().WithCodec(cdc)
            txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))

            name := args[0]
            oldField := args[1]
            newField := args[2]
            msg := types.NewMsgRenameField(cliCtx.GetFromAddress(), name, oldField, newField)
            err := msg.ValidateBasic()
            if err != nil {
                return err
            }

            return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
        },
    }
}

// GetCmdCreateIndex is the CLI command for sending a CreateIndex transaction
func GetCmdCreateIndex(cdc *codec.Codec) *cobra.Command {
    return &cobra.Command{
        Use:   "create-index [tableName] [field]",
        Short: "create a new index",
        Args:  cobra.ExactArgs(2),
        RunE: func(cmd *cobra.Command, args []string) error {
            cliCtx := context.NewCLIContext().WithCodec(cdc)
            txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))

            tableName := args[0]
            field := args[1]
            msg := types.NewMsgCreateIndex(cliCtx.GetFromAddress(), tableName, field)
            err := msg.ValidateBasic()
            if err != nil {
                return err
            }

            return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
        },
    }
}

func GetCmdModifyOption(cdc *codec.Codec) *cobra.Command {
    return &cobra.Command{
        Use:   "modify-table-option [tableName] [action] [option]",
        Short: "modify table options",
        Args:  cobra.ExactArgs(3),
        RunE: func(cmd *cobra.Command, args []string) error {
            cliCtx := context.NewCLIContext().WithCodec(cdc)
            txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))

            tableName := args[0]
            action := args[1]
            option := args[2]

            msg := types.NewMsgModifyOption(cliCtx.GetFromAddress(), tableName, action, option)
            err := msg.ValidateBasic()
            if err != nil {
                return err
            }

            return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
        },
    }
}

func GetCmdModifyFieldOption(cdc *codec.Codec) *cobra.Command {
    return &cobra.Command{
        Use:   "modify-field-option [tableName] [fieldName] [action] [option]",
        Short: "modify field options",
        Args:  cobra.ExactArgs(4),
        RunE: func(cmd *cobra.Command, args []string) error {
            cliCtx := context.NewCLIContext().WithCodec(cdc)
            txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))

            tableName := args[0]
            fieldName := args[1]
            action := args[2]
            option := args[3]

            msg := types.NewMsgModifyFieldOption(cliCtx.GetFromAddress(), tableName, fieldName, action, option)
            err := msg.ValidateBasic()
            if err != nil {
                return err
            }

            return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
        },
    }
}


///////////////////////////////
//                           //
// data manipulation related //
//                           //
///////////////////////////////

func GetCmdInsertRow(cdc *codec.Codec) *cobra.Command {
    return &cobra.Command{
        Use:   "insert-row [tableName] [fields] [values]",
        Short: "create a new row",
        Args:  cobra.ExactArgs(3),
        RunE: func(cmd *cobra.Command, args []string) error {
            cliCtx := context.NewCLIContext().WithCodec(cdc)
            txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))

            name := args[0]
            fields := strings.Split(args[1], ",")
            values := strings.Split(args[2], ",")
            rowFields := make(types.RowFields)
            for i, field := range fields {
                if i < len(values) {
                    rowFields[field] = values[i]
                }
            }

            rowFieldsJson, err := json.Marshal(rowFields)
            if err != nil { return err } 

            msg := types.NewMsgInsertRow(cliCtx.GetFromAddress(), name, rowFieldsJson)
            err = msg.ValidateBasic()
            if err != nil {
                return errors.New(fmt.Sprintf("Error %s", err)) 
            }

            return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
        },
    }
}

func GetCmdUpdateRow(cdc *codec.Codec) *cobra.Command {
    return &cobra.Command{
        Use:   "update-row [tableName] [id] [fields] [values]",
        Short: "update a row",
        Args:  cobra.ExactArgs(4),
        RunE: func(cmd *cobra.Command, args []string) error {
            cliCtx := context.NewCLIContext().WithCodec(cdc)
            txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))

            name := args[0]
            id, err := strconv.ParseUint(args[1], 10, 0)
            if err != nil {
                return errors.New(fmt.Sprintf("Error %s", err))
            }
            fields := strings.Split(args[2], ",")
            values := strings.Split(args[3], ",")
            rowFields := make(types.RowFields)
            for i, field := range fields {
                if i < len(values) {
                    rowFields[field] = values[i]
                }
            }

            rowFieldsJson, err := json.Marshal(rowFields)
            if err != nil { return err }

            msg := types.NewMsgUpdateRow(cliCtx.GetFromAddress(), name, uint(id), rowFieldsJson)
            err = msg.ValidateBasic()
            if err != nil {
                return errors.New(fmt.Sprintf("Error %s", err))
            }

            return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
        },
    }
}

func GetCmdDeleteRow(cdc *codec.Codec) *cobra.Command {
    return &cobra.Command{
        Use:   "delete-row [tableName] [id]",
        Short: "delete a row",
        Args:  cobra.ExactArgs(2),
        RunE: func(cmd *cobra.Command, args []string) error {
            cliCtx := context.NewCLIContext().WithCodec(cdc)
            txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))

            name := args[0]
            id, err := strconv.ParseUint(args[1], 10, 0)
            if err != nil {
                return errors.New(fmt.Sprintf("Error %s", err))
            }

            msg := types.NewMsgDeleteRow(cliCtx.GetFromAddress(), name, uint(id))
            err = msg.ValidateBasic()
            if err != nil {
                return errors.New(fmt.Sprintf("Error %s", err))
            }

            return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
        },
    }
}

/////////////////////////
//                     //
// admin group related //
//                     //
/////////////////////////

func GetCmdAddAdminAccount(cdc * codec.Codec) *cobra.Command {
    return &cobra.Command{
        Use:   "add-admin [address]",
        Short: "add an account into admin group",
        Args:  cobra.ExactArgs(1),
        RunE: func(cmd *cobra.Command, args []string) error {
            cliCtx := context.NewCLIContext().WithCodec(cdc)
            txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))

            addr, err := sdk.AccAddressFromBech32(args[0]) // args[0] is the new admin address

            if err != nil {
                return errors.New(fmt.Sprintf("Error %s", err))
            }

            msg := types.NewMsgAddAdminAccount(addr, cliCtx.GetFromAddress())
            err = msg.ValidateBasic()
            if err != nil {
                return errors.New(fmt.Sprintf("Error %s", err))
            }
            return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
        },
    }
}



//////////////////////
//                  //
// helper functions //
//                  //
//////////////////////

