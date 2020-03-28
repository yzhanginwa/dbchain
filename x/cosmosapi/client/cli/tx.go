package cli

import (
    "bufio"
    "fmt"
    "errors"
    "strings"
    "strconv"
    "encoding/json"
    "github.com/spf13/cobra"

    "github.com/cosmos/cosmos-sdk/client"
    "github.com/cosmos/cosmos-sdk/client/context"
    "github.com/cosmos/cosmos-sdk/client/flags"
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

    cosmosapiTxCmd.AddCommand(flags.PostCommands(
        GetCmdCreateApplication(cdc),
        GetCmdAddAppUser(cdc),
        GetCmdCreateTable(cdc),
        GetCmdDropTable(cdc),
        GetCmdAddColumn(cdc),
        GetCmdDropColumn(cdc),
        GetCmdRenameColumn(cdc),
        GetCmdCreateIndex(cdc),
        GetCmdDropIndex(cdc),
        GetCmdModifyOption(cdc),
        GetCmdModifyColumnOption(cdc),
        GetCmdInsertRow(cdc),
        GetCmdUpdateRow(cdc),
        GetCmdDeleteRow(cdc),
        GetCmdFreezeRow(cdc),
        GetCmdAddAdminAccount(cdc),
        GetCmdAddFriend(cdc),
        GetCmdRespondFriend(cdc),
    )...)

    return cosmosapiTxCmd
}

////////////////////
//                //
// schema related //
//                //
////////////////////

func GetCmdCreateApplication(cdc *codec.Codec) *cobra.Command {
    return &cobra.Command{
        Use:   "create-application",
        Short: "create a new application",
        Args:  cobra.ExactArgs(3),
        RunE: func(cmd *cobra.Command, args []string) error {
            cliCtx := context.NewCLIContext().WithCodec(cdc)
            inBuf := bufio.NewReader(cmd.InOrStdin())
            txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(cdc))

            name := args[0]
            description := args[1]
            var permissioned = true
            if args[2] == "no" || args[2] == "false" {
                permissioned = false
            }
            msg := types.NewMsgCreateApplication(cliCtx.GetFromAddress(), name, description, permissioned)
            err := msg.ValidateBasic()
            if err != nil {
                return err
            }

            return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
        },
    }
}

func GetCmdAddAppUser(cdc *codec.Codec) *cobra.Command {
    return &cobra.Command{
        Use:   "add-app-user [appCode] [address]",
        Short: "add application user",
        Args:  cobra.ExactArgs(2),
        RunE: func(cmd *cobra.Command, args []string) error {
            cliCtx := context.NewCLIContext().WithCodec(cdc)
            inBuf := bufio.NewReader(cmd.InOrStdin())
            txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(cdc))

            appCode := args[0]
            address := args[1]
            user, err := sdk.AccAddressFromBech32(address)
            if err != nil {
                return errors.New(fmt.Sprintf("Error %s", err))
            }

            msg := types.NewMsgAddDatabaseUser(cliCtx.GetFromAddress(), appCode, user)
            err = msg.ValidateBasic()
            if err != nil {
                return err
            }

            return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
        },
    }
}

func GetCmdCreateTable(cdc *codec.Codec) *cobra.Command {
    return &cobra.Command{
        Use:   "create-table [appCode] [name] [fields]",
        Short: "create a new table",
        Args:  cobra.ExactArgs(3),
        RunE: func(cmd *cobra.Command, args []string) error {
            cliCtx := context.NewCLIContext().WithCodec(cdc)
            inBuf := bufio.NewReader(cmd.InOrStdin())
            txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(cdc))

            appCode := args[0]
            name := args[1]
            fields := strings.Split(args[2], ",")
            msg := types.NewMsgCreateTable(cliCtx.GetFromAddress(), appCode, name, fields)
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
        Use:   "drop-table [appCode] [name]",
        Short: "drop a table",
        Args:  cobra.ExactArgs(2),
        RunE: func(cmd *cobra.Command, args []string) error {
            cliCtx := context.NewCLIContext().WithCodec(cdc)
            inBuf := bufio.NewReader(cmd.InOrStdin())
            txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(cdc))

            appCode := args[0]
            name := args[1]
            msg := types.NewMsgDropTable(cliCtx.GetFromAddress(), appCode, name)
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
        Use:   "add-column [appCode] [name] [field]",
        Short: "add a new column onto a table",
        Args:  cobra.ExactArgs(3),
        RunE: func(cmd *cobra.Command, args []string) error {
            cliCtx := context.NewCLIContext().WithCodec(cdc)
            inBuf := bufio.NewReader(cmd.InOrStdin())
            txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(cdc))

            appCode:= args[0]
            name   := args[1]
            field  := args[2]
            msg := types.NewMsgAddColumn(cliCtx.GetFromAddress(), appCode, name, field)
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
        Use:   "drop-column [appCode] [name] [field]",
        Short: "drop a column from a table",
        Args:  cobra.ExactArgs(3),
        RunE: func(cmd *cobra.Command, args []string) error {
            cliCtx := context.NewCLIContext().WithCodec(cdc)
            inBuf := bufio.NewReader(cmd.InOrStdin())
            txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(cdc))

            appCode := args[0]
            name    := args[1]
            field   := args[2]
            msg := types.NewMsgDropColumn(cliCtx.GetFromAddress(), appCode, name, field)
            err := msg.ValidateBasic()
            if err != nil {
                return err
            }

            return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
        },
    }
}

func GetCmdRenameColumn(cdc *codec.Codec) *cobra.Command {
    return &cobra.Command{
        Use:   "rename-column [appCode] [name] [old-field] [new-field",
        Short: "rename a column in a table",
        Args:  cobra.ExactArgs(4),
        RunE: func(cmd *cobra.Command, args []string) error {
            cliCtx := context.NewCLIContext().WithCodec(cdc)
            inBuf := bufio.NewReader(cmd.InOrStdin())
            txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(cdc))

            appCode  := args[0]
            name     := args[1]
            oldField := args[2]
            newField := args[3]
            msg := types.NewMsgRenameColumn(cliCtx.GetFromAddress(), appCode, name, oldField, newField)
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
        Use:   "create-index [appCode] [tableName] [field]",
        Short: "create a new index",
        Args:  cobra.ExactArgs(3),
        RunE: func(cmd *cobra.Command, args []string) error {
            cliCtx := context.NewCLIContext().WithCodec(cdc)
            inBuf := bufio.NewReader(cmd.InOrStdin())
            txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(cdc))

            appCode   := args[0]
            tableName := args[1]
            field     := args[2]
            msg := types.NewMsgCreateIndex(cliCtx.GetFromAddress(), appCode, tableName, field)
            err := msg.ValidateBasic()
            if err != nil {
                return err
            }

            return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
        },
    }
}

func GetCmdDropIndex(cdc *codec.Codec) *cobra.Command {
    return &cobra.Command{
        Use:   "drop-index [appCode] [tableName] [field]",
        Short: "drop an index",
        Args:  cobra.ExactArgs(3),
        RunE: func(cmd *cobra.Command, args []string) error {
            cliCtx := context.NewCLIContext().WithCodec(cdc)
            inBuf := bufio.NewReader(cmd.InOrStdin())
            txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(cdc))

            appCode   := args[0]
            tableName := args[1]
            field     := args[2]
            msg := types.NewMsgDropIndex(cliCtx.GetFromAddress(), appCode, tableName, field)
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
        Use:   "modify-table-option [appCode] [tableName] [action] [option]",
        Short: "modify table options",
        Args:  cobra.ExactArgs(4),
        RunE: func(cmd *cobra.Command, args []string) error {
            cliCtx := context.NewCLIContext().WithCodec(cdc)
            inBuf := bufio.NewReader(cmd.InOrStdin())
            txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(cdc))

            appCode   := args[0]
            tableName := args[1]
            action    := args[2]
            option    := args[3]

            msg := types.NewMsgModifyOption(cliCtx.GetFromAddress(), appCode, tableName, action, option)
            err := msg.ValidateBasic()
            if err != nil {
                return err
            }

            return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
        },
    }
}

func GetCmdModifyColumnOption(cdc *codec.Codec) *cobra.Command {
    return &cobra.Command{
        Use:   "modify-column-option [appCode] [tableName] [fieldName] [action] [option]",
        Short: "modify column options",
        Args:  cobra.ExactArgs(5),
        RunE: func(cmd *cobra.Command, args []string) error {
            cliCtx := context.NewCLIContext().WithCodec(cdc)
            inBuf := bufio.NewReader(cmd.InOrStdin())
            txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(cdc))

            appCode   := args[0]
            tableName := args[1]
            fieldName := args[2]
            action    := args[3]
            option    := args[4]

            msg := types.NewMsgModifyColumnOption(cliCtx.GetFromAddress(), appCode, tableName, fieldName, action, option)
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
        Use:   "insert-row [appCode] [tableName] [fields] [values]",
        Short: "create a new row",
        Args:  cobra.ExactArgs(4),
        RunE: func(cmd *cobra.Command, args []string) error {
            cliCtx := context.NewCLIContext().WithCodec(cdc)
            inBuf := bufio.NewReader(cmd.InOrStdin())
            txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(cdc))

            appCode := args[0]
            name    := args[1]
            fields  := strings.Split(args[2], ",")
            values  := strings.Split(args[3], ",")
            rowFields := make(types.RowFields)
            for i, field := range fields {
                if i < len(values) {
                    rowFields[field] = values[i]
                }
            }

            rowFieldsJson, err := json.Marshal(rowFields)
            if err != nil { return err } 

            msg := types.NewMsgInsertRow(cliCtx.GetFromAddress(), appCode, name, rowFieldsJson)
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
        Use:   "update-row [appCode] [tableName] [id] [fields] [values]",
        Short: "update a row",
        Args:  cobra.ExactArgs(5),
        RunE: func(cmd *cobra.Command, args []string) error {
            cliCtx := context.NewCLIContext().WithCodec(cdc)
            inBuf := bufio.NewReader(cmd.InOrStdin())
            txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(cdc))

            appCode := args[0]
            name    := args[1]
            id, err := strconv.ParseUint(args[2], 10, 0)
            if err != nil {
                return errors.New(fmt.Sprintf("Error %s", err))
            }
            fields := strings.Split(args[3], ",")
            values := strings.Split(args[4], ",")
            rowFields := make(types.RowFields)
            for i, field := range fields {
                if i < len(values) {
                    rowFields[field] = values[i]
                }
            }

            rowFieldsJson, err := json.Marshal(rowFields)
            if err != nil { return err }

            msg := types.NewMsgUpdateRow(cliCtx.GetFromAddress(), appCode, name, uint(id), rowFieldsJson)
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
        Use:   "delete-row [appCode] [tableName] [id]",
        Short: "delete a row",
        Args:  cobra.ExactArgs(3),
        RunE: func(cmd *cobra.Command, args []string) error {
            cliCtx := context.NewCLIContext().WithCodec(cdc)
            inBuf := bufio.NewReader(cmd.InOrStdin())
            txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(cdc))

            appCode := args[0]
            name    := args[1]
            id, err := strconv.ParseUint(args[2], 10, 0)
            if err != nil {
                return errors.New(fmt.Sprintf("Error %s", err))
            }

            msg := types.NewMsgDeleteRow(cliCtx.GetFromAddress(), appCode, name, uint(id))
            err = msg.ValidateBasic()
            if err != nil {
                return errors.New(fmt.Sprintf("Error %s", err))
            }

            return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
        },
    }
}

func GetCmdFreezeRow(cdc *codec.Codec) *cobra.Command {
    return &cobra.Command{
        Use:   "freeze-row [appCode] [tableName] [id]",
        Short: "freeze a row",
        Args:  cobra.ExactArgs(3),
        RunE: func(cmd *cobra.Command, args []string) error {
            cliCtx := context.NewCLIContext().WithCodec(cdc)
            inBuf := bufio.NewReader(cmd.InOrStdin())
            txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(cdc))

            appCode := args[0]
            tableName    := args[1]
            id, err := strconv.ParseUint(args[2], 10, 0)
            if err != nil {
                return errors.New(fmt.Sprintf("Error %s", err))
            }

            msg := types.NewMsgFreezeRow(cliCtx.GetFromAddress(), appCode, tableName, uint(id))
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
        Use:   "add-admin [appCode] [address]",
        Short: "add an account into admin group",
        Args:  cobra.ExactArgs(2),
        RunE: func(cmd *cobra.Command, args []string) error {
            cliCtx := context.NewCLIContext().WithCodec(cdc)
            inBuf := bufio.NewReader(cmd.InOrStdin())
            txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(cdc))

            appCode := args[0]
            adminAddress := args[1]
            addr, err := sdk.AccAddressFromBech32(adminAddress)

            if err != nil {
                return errors.New(fmt.Sprintf("Error %s", err))
            }

            msg := types.NewMsgAddAdminAccount(appCode, addr, cliCtx.GetFromAddress())
            err = msg.ValidateBasic()
            if err != nil {
                return errors.New(fmt.Sprintf("Error %s", err))
            }
            return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
        },
    }
}

////////////////
//            //
// add friend //
//            //
////////////////

func GetCmdAddFriend(cdc * codec.Codec) *cobra.Command {
    return &cobra.Command{
        Use:   "add-friend [my-name] [address] [name]",
        Short: "add a friend ",
        Args:  cobra.ExactArgs(3),
        RunE: func(cmd *cobra.Command, args []string) error {
            cliCtx := context.NewCLIContext().WithCodec(cdc)
            inBuf := bufio.NewReader(cmd.InOrStdin())
            txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(cdc))

            ownerName := args[0]
            address   := args[1]
            name      := args[2]
            msg := types.NewMsgAddFriend(cliCtx.GetFromAddress(), ownerName, address, name)
            err := msg.ValidateBasic()
            if err != nil {
                return errors.New(fmt.Sprintf("Error %s", err))
            }
            return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
        },
    }
}

func GetCmdRespondFriend(cdc * codec.Codec) *cobra.Command {
    return &cobra.Command{
        Use:   "respond-friend [address] [action]",
        Short: "Respond a friend. The action could be delete, accept, reject.",
        Args:  cobra.ExactArgs(2),
        RunE: func(cmd *cobra.Command, args []string) error {
            cliCtx := context.NewCLIContext().WithCodec(cdc)
            inBuf := bufio.NewReader(cmd.InOrStdin())
            txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(cdc))

            address := args[0]
            action  := args[1]
            msg := types.NewMsgRespondFriend(cliCtx.GetFromAddress(), address, action)
            err := msg.ValidateBasic()
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

