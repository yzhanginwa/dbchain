package cli

import (
    "bufio"
    "encoding/json"
    "errors"
    "fmt"
    "github.com/dbchaincloud/tendermint/crypto/algo"
    "github.com/dbchaincloud/tendermint/crypto/secp256k1"
    "github.com/dbchaincloud/tendermint/crypto/sm2"
    "github.com/mr-tron/base58"
    "github.com/spf13/cobra"
    "github.com/spf13/viper"
    "github.com/yzhanginwa/dbchain/x/bank"
    "os"
    "path"
    "strconv"
    "strings"

    "github.com/dbchaincloud/cosmos-sdk/client"
    "github.com/dbchaincloud/cosmos-sdk/client/context"
    "github.com/dbchaincloud/cosmos-sdk/client/flags"
    "github.com/dbchaincloud/cosmos-sdk/codec"
    sdk "github.com/dbchaincloud/cosmos-sdk/types"
    "github.com/dbchaincloud/cosmos-sdk/x/auth"
    "github.com/dbchaincloud/cosmos-sdk/x/auth/client/utils"
    account "github.com/dbchaincloud/cosmos-sdk/x/auth/types"
    "github.com/yzhanginwa/dbchain/x/dbchain/internal/types"
)

func GetTxCmd(storeKey string, cdc *codec.Codec) *cobra.Command {
    dbchainTxCmd := &cobra.Command{
        Use:                        types.ModuleName,
        Short:                      "DbChain transaction subcommands",
        DisableFlagParsing:         true,
        SuggestionsMinimumDistance: 2,
        RunE:                       client.ValidateCmd,
    }

    dbchainTxCmd.AddCommand(flags.PostCommands(
        GetCmdCreateApplication(cdc),
        GetCmdDropApplication(cdc),
        GetCmdRecoverApplication(cdc),
        GetCmdCreateSysDatabase(cdc),
        GetCmdSetAppUserFileVolumeLimit(cdc),
        GetCmdModifyAppUser(cdc),
        GetCmdSetAppPermission(cdc),
        GetCmdAddFunction(cdc),
        GetCmdCallFunction(cdc),
        GetCmdDropFunction(cdc),
        GetCmdAddCustomQuerier(cdc),
        GetCmdDropCustomQuerier(cdc),
        GetCmdCreateTable(cdc),
        GetCmdModifyTableAssociation(cdc),
        GetCmdAddCounterCache(cdc),
        GetCmdDropTable(cdc),
        GetCmdAddColumn(cdc),
        GetCmdDropColumn(cdc),
        GetCmdRenameColumn(cdc),
        GetCmdCreateIndex(cdc),
        GetCmdDropIndex(cdc),
        GetCmdModifyOption(cdc),
        GetCmdAddInsertFilter(cdc),
        GetCmdDropInsertFilter(cdc),
        GetCmdAddTrigger(cdc),
        GetCmdDropTrigger(cdc),
        GetCmdSetTableMemo(cdc),
        GetCmdModifyColumnOption(cdc),
        GetCmdSetColumnDataType(cdc),
        GetCmdSetColumnMemo(cdc),
        GetCmdInsertRow(cdc),
        GetCmdUpdateRow(cdc),
        GetCmdDeleteRow(cdc),
        GetCmdFreezeRow(cdc),
        GetCmdModifyGroup(cdc),
        GetCmdSetGroupMemo(cdc),
        GetCmdModifyGroupMember(cdc),
        GetCmdAddFriend(cdc),
        GetCmdDropFriend(cdc),
        GetCmdRespondFriend(cdc),
        GetCmdFreezeSchema(cdc),
        GetCmdUnfreezeSchema(cdc),
        GetCmdResetTxTotalTxs(cdc),
        GetCmdModifyTokenKeepers(cdc),
        GetCmdModifyP2PTransferLimit(cdc),
    )...)

    return dbchainTxCmd
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
            var permissionRequired = true
            if args[2] == "no" || args[2] == "false" {
                permissionRequired = false
            }
            msg := types.NewMsgCreateApplication(cliCtx.GetFromAddress(), name, description, permissionRequired)
            err := msg.ValidateBasic()
            if err != nil {
                return err
            }

            return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
        },
    }
}

func GetCmdDropApplication(cdc *codec.Codec) *cobra.Command {
    return &cobra.Command{
        Use:   "drop-application",
        Short: "drop a application",
        Args:  cobra.ExactArgs(1),
        RunE: func(cmd *cobra.Command, args []string) error {
            cliCtx := context.NewCLIContext().WithCodec(cdc)
            inBuf := bufio.NewReader(cmd.InOrStdin())
            txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(cdc))

            appCode := args[0]
            msg := types.NewMsgDropApplication(cliCtx.GetFromAddress(), appCode)
            err := msg.ValidateBasic()
            if err != nil {
                return err
            }

            return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
        },
    }
}

func GetCmdRecoverApplication(cdc *codec.Codec) *cobra.Command {
    return &cobra.Command{
        Use:   "recover-application",
        Short: "recover a application",
        Args:  cobra.ExactArgs(1),
        RunE: func(cmd *cobra.Command, args []string) error {
            cliCtx := context.NewCLIContext().WithCodec(cdc)
            inBuf := bufio.NewReader(cmd.InOrStdin())
            txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(cdc))

            appCode := args[0]
            msg := types.NewMsgRecoverApplication(cliCtx.GetFromAddress(), appCode)
            err := msg.ValidateBasic()
            if err != nil {
                return err
            }

            return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
        },
    }
}

func GetCmdCreateSysDatabase(cdc *codec.Codec) *cobra.Command {
    return &cobra.Command{
        Use:   "create-sys-database",
        Short: "create a system database",
        Args:  cobra.ExactArgs(0),
        RunE: func(cmd *cobra.Command, args []string) error {
            cliCtx := context.NewCLIContext().WithCodec(cdc)
            inBuf := bufio.NewReader(cmd.InOrStdin())
            txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(cdc))

            msgs, err := createSysDatabaseMsg(cliCtx, cliCtx.GetFromAddress())
            if err != nil {
                return err
            }
            return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, msgs)
        },
    }
}

func GetCmdSetAppUserFileVolumeLimit(cdc *codec.Codec) *cobra.Command {
    return &cobra.Command{
        Use:   "set-app-user-file-limit",
        Short: "set application user file volume limit. Uint of size is byte. when size was set 0 or negative, it means no limit",
        Args:  cobra.ExactArgs(2),
        RunE: func(cmd *cobra.Command, args []string) error {
            cliCtx := context.NewCLIContext().WithCodec(cdc)
            inBuf := bufio.NewReader(cmd.InOrStdin())
            txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(cdc))

            appCode := args[0]
            size    := args[1]
            msg := types.NewMsgSetAppUserFileVolumeLimit(cliCtx.GetFromAddress(), appCode, size)
            err := msg.ValidateBasic()
            if err != nil {
                return errors.New(fmt.Sprintf("Error %s", err))
            }
            return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
        },
    }
}

func createSysDatabaseMsg(cliCtx context.CLIContext, adminAddr sdk.AccAddress)([]sdk.Msg, error) {
    msgs := make([]sdk.Msg, 0)

    oracleAddr ,err := checkOracleInfo()
    if err != nil {
        return nil, err
    }

    if !checkAddrBalance(cliCtx, oracleAddr) {
        sendMsg := sendOneCoins(adminAddr, oracleAddr)
        msgs = append(msgs, sendMsg)
    }

    var msg sdk.Msg
    msg = types.NewMsgCreateSysDatabase(adminAddr)
    err = msg.ValidateBasic()
    if err != nil {
       return nil, err
    }

    msgs = append(msgs, msg)

    tables := map[string][]string{
        "authentication" : []string{"address", "type", "value"},
        "order_receipt"  : []string{"appcode", "owner", "orderid", "amount", "expiration_date", "vendor", "vendor_payment_no"},
    }
    for tableName, fileds := range tables {
        msg := types.NewMsgCreateTable(adminAddr,"0000000001",tableName, fileds)
        if err := msg.ValidateBasic(); err != nil {
            return nil, err
        }
        msgs = append(msgs, msg)
    }
    msgs = append(msgs, types.NewMsgModifyColumnOption(adminAddr, "0000000001", "order_receipt", "vendor_payment_no", "add", string(types.FLDOPT_UNIQUE)))
    msgs = append(msgs, types.NewMsgModifyGroup("0000000001", "add", "oracle", adminAddr))
    msgs = append(msgs, types.NewMsgModifyGroupMember("0000000001", "oracle", "add", oracleAddr, adminAddr))
    msgs = append(msgs, types.NewMsgSetGroupMemo("0000000001","oracle","oracleOfThisChain",adminAddr))
    //add table option
    msgs = append(msgs, types.NewMsgModifyOption(adminAddr, "0000000001", "authentication", "add", "writable-by(oracle)" ))
    msgs = append(msgs, types.NewMsgModifyOption(adminAddr, "0000000001", "order_receipt", "add", "writable-by(oracle)" ))
    return msgs, nil

}

func checkOracleInfo() (sdk.AccAddress, error){
    DefaultOracleHome := os.ExpandEnv(types.OracleHome)
    tempViper := viper.New()

    cfgFile := path.Join(DefaultOracleHome, "config", "config.toml")
    if _, err := os.Stat(cfgFile); err != nil {
        return nil, err
    }
    tempViper.SetConfigFile(cfgFile)
    if err := tempViper.ReadInConfig(); err != nil {
        return nil, err
    }

    oraclePrivKey := tempViper.GetString("oracle-encrypted-key")
    if oraclePrivKey == "" {
        fmt.Println(`please execute cmd : "dbchainoracle  query oracle-info " `)
        fmt.Println(`then set it in the "`, cfgFile, `" like the following format :`)
        fmt.Println(`oracle-encrypted-key = "you private key"`)
        return nil, errors.New("oracle-encrypted-key not set")
    }

    pkBytes, err:= base58.Decode(oraclePrivKey)
    if err != nil {
        return nil, err
    }
    var addr sdk.AccAddress
    switch algo.Algo {
    case algo.SM2:
        var privKey sm2.PrivKeySm2
        copy(privKey[:], pkBytes)
        addr = sdk.AccAddress(privKey.PubKeySm2().Address())
    default:
        var privKey secp256k1.PrivKeySecp256k1
        copy(privKey[:], pkBytes)
        addr = sdk.AccAddress(privKey.PubKey().Address())
    }
    return addr, nil
}

func checkAddrBalance(cliCtx context.CLIContext, addr sdk.AccAddress) bool {
    accGetter := account.NewAccountRetriever(cliCtx)
    acc , err := accGetter.GetAccount(addr)
    if err != nil {
        return false
    }
    coins := acc.GetCoins()
    if coins.Empty() {
        return false
    }
    return true
}

func sendOneCoins(from , to  sdk.AccAddress) sdk.Msg {
    oneCoin := sdk.NewCoin("dbctoken", sdk.NewInt(1))
    msg := bank.NewMsgSend(from, to, []sdk.Coin{oneCoin})
    return msg
}

func GetCmdSetAppPermission(cdc * codec.Codec) *cobra.Command {
    return &cobra.Command{
        Use:   "set-app-permission [database] [permission_required]",
        Short: "Set the permission_required status of database",
        Args:  cobra.ExactArgs(2),
        RunE: func(cmd *cobra.Command, args []string) error {
            cliCtx := context.NewCLIContext().WithCodec(cdc)
            inBuf := bufio.NewReader(cmd.InOrStdin())
            txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(cdc))

            appCode            := args[0]
            permissionRequired := args[1]
            msg := types.NewMsgSetDatabasePermission(cliCtx.GetFromAddress(), appCode, permissionRequired)
            err := msg.ValidateBasic()
            if err != nil {
                return errors.New(fmt.Sprintf("Error %s", err))
            }
            return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
        },
    }
}

func GetCmdModifyAppUser(cdc *codec.Codec) *cobra.Command {
    return &cobra.Command{
        Use:   "modify-app-user [appCode] [action] [address]",
        Short: "modify application user",
        Args:  cobra.ExactArgs(3),
        RunE: func(cmd *cobra.Command, args []string) error {
            cliCtx := context.NewCLIContext().WithCodec(cdc)
            inBuf := bufio.NewReader(cmd.InOrStdin())
            txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(cdc))

            appCode := args[0]
            action  := args[1]   // action has to be either 'add' or 'drop'
            address := args[2]
            user, err := sdk.AccAddressFromBech32(address)
            if err != nil {
                return errors.New(fmt.Sprintf("Error %s", err))
            }

            msg := types.NewMsgModifyDatabaseUser(cliCtx.GetFromAddress(), appCode, action, user)
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

func GetCmdModifyTableAssociation(cdc *codec.Codec) *cobra.Command {
    return &cobra.Command{
        Use:   "modify-table-association [appCode] [tableName] [associationMode] [associationTable] [method] [foreignKey] [option]",
        Short: "add or drop table association a new table",
        Args:  cobra.ExactArgs(7),
        RunE: func(cmd *cobra.Command, args []string) error {
            cliCtx := context.NewCLIContext().WithCodec(cdc)
            inBuf := bufio.NewReader(cmd.InOrStdin())
            txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(cdc))

            appCode := args[0]
            option := args[1]
            tableName := args[2]
            associationMode := args[3]
            associationTable := args[4]
            method := args[5]
            foreignKey := args[6]


            msg := types.NewMsgModifyTableAssociation(appCode,tableName,associationMode,associationTable,method,foreignKey,option,cliCtx.GetFromAddress())
            err := msg.ValidateBasic()
            if err != nil {
                return err
            }

            return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
        },
    }
}

func GetCmdAddCounterCache(cdc *codec.Codec) *cobra.Command {
    return &cobra.Command{
        Use:   "add-counter-cache [appCode] [tableName] [associationTable] [foreignKey] [counterCacheField] [limit]",
        Short: "add a counter cache field for a table, foreignKey is a field of this table, counterCacheField is a new field which will be add to associationTable. when limit is 0 or " +
            "negative , it means its no limit",
        Args:  cobra.ExactArgs(6),
        RunE: func(cmd *cobra.Command, args []string) error {
            cliCtx := context.NewCLIContext().WithCodec(cdc)
            inBuf := bufio.NewReader(cmd.InOrStdin())
            txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(cdc))

            appCode := args[0]
            tableName := args[1]
            associationTable := args[2]
            foreignKey := args[3]
            counterCacheField := args[4]
            limit := args[5]


            msg := types.NewMsgEnableCounterCache(appCode, tableName, associationTable, foreignKey, counterCacheField, limit, cliCtx.GetFromAddress())
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

////////////////////
//                //
// Set table memo //
//                //
////////////////////

func GetCmdSetTableMemo(cdc * codec.Codec) *cobra.Command {
    return &cobra.Command{
        Use:   "set-table-memo [appCode] [table] [memo]",
        Short: "set table memo",
        Args:  cobra.ExactArgs(3),
        RunE: func(cmd *cobra.Command, args []string) error {
            cliCtx := context.NewCLIContext().WithCodec(cdc)
            inBuf := bufio.NewReader(cmd.InOrStdin())
            txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(cdc))

            appCode   := args[0]
            tableName := args[1]
            memo      := args[2]

            msg := types.NewMsgSetTableMemo(appCode, tableName, memo, cliCtx.GetFromAddress())
            err := msg.ValidateBasic()
            if err != nil {
                return errors.New(fmt.Sprintf("Error %s", err))
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

func GetCmdSetColumnDataType(cdc *codec.Codec) *cobra.Command {
    return &cobra.Command{
        Use:   "modify-column-data-type [appCode] [tableName] [fieldName] [type]",
        Short: "modify column data type, support int , file ,decimal",
        Args:  cobra.ExactArgs(4),
        RunE: func(cmd *cobra.Command, args []string) error {
            cliCtx := context.NewCLIContext().WithCodec(cdc)
            inBuf := bufio.NewReader(cmd.InOrStdin())
            txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(cdc))

            appCode   := args[0]
            tableName := args[1]
            fieldName := args[2]
            dataType  := args[3]

            msg := types.NewMsgSetColumnDataType(cliCtx.GetFromAddress(), appCode, tableName, fieldName, dataType)
            err := msg.ValidateBasic()
            if err != nil {
                return err
            }

            return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
        },
    }
}

/////////////////////
//                 //
// Set column memo //
//                 //
/////////////////////

func GetCmdSetColumnMemo(cdc * codec.Codec) *cobra.Command {
    return &cobra.Command{
        Use:   "set-column-memo [appCode] [table] [field] [memo]",
        Short: "set column memo",
        Args:  cobra.ExactArgs(4),
        RunE: func(cmd *cobra.Command, args []string) error {
            cliCtx := context.NewCLIContext().WithCodec(cdc)
            inBuf := bufio.NewReader(cmd.InOrStdin())
            txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(cdc))

            appCode   := args[0]
            tableName := args[1]
            fieldName := args[2]
            memo      := args[3]

            msg := types.NewMsgSetColumnMemo(appCode, tableName, fieldName, memo, cliCtx.GetFromAddress())
            err := msg.ValidateBasic()
            if err != nil {
                return errors.New(fmt.Sprintf("Error %s", err))
            }
            return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
        },
    }
}

///////////////////////////////
//                           //
// validation for new record //
//                           //
///////////////////////////////

func GetCmdAddInsertFilter(cdc *codec.Codec) *cobra.Command {
    return &cobra.Command{
        Use:   "add-insert-filter [appCode] [tableName] [filter-text]",
        Short: "add an insert filter",
        Args:  cobra.ExactArgs(3),
        RunE: func(cmd *cobra.Command, args []string) error {
            cliCtx := context.NewCLIContext().WithCodec(cdc)
            inBuf := bufio.NewReader(cmd.InOrStdin())
            txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(cdc))

            appCode   := args[0]
            tableName := args[1]
            filter    := args[2]

            msg := types.NewMsgAddInsertFilter(cliCtx.GetFromAddress(), appCode, tableName, filter)
            err := msg.ValidateBasic()
            if err != nil {
                return errors.New(fmt.Sprintf("Error %s", err))
            }

            return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
        },
    }
}

func GetCmdDropInsertFilter(cdc *codec.Codec) *cobra.Command {
    return &cobra.Command{
        Use:   "drop-insert-filter [appCode] [tableName]",
        Short: "drop an insert filter",
        Args:  cobra.ExactArgs(2),
        RunE: func(cmd *cobra.Command, args []string) error {
            cliCtx := context.NewCLIContext().WithCodec(cdc)
            inBuf := bufio.NewReader(cmd.InOrStdin())
            txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(cdc))

            appCode   := args[0]
            tableName := args[1]

            msg := types.NewMsgDropInsertFilter(cliCtx.GetFromAddress(), appCode, tableName)
            err := msg.ValidateBasic()
            if err != nil {
                return errors.New(fmt.Sprintf("Error %s", err))
            }

            return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
        },
    }
}

////////////////////////////
//                        //
// trigger for new record //
//                        //
////////////////////////////

func GetCmdAddTrigger(cdc *codec.Codec) *cobra.Command {
    return &cobra.Command{
        Use:   "add-trigger [appCode] [tableName] [trigger-text]",
        Short: "add a trigger",
        Args:  cobra.ExactArgs(3),
        RunE: func(cmd *cobra.Command, args []string) error {
            cliCtx := context.NewCLIContext().WithCodec(cdc)
            inBuf := bufio.NewReader(cmd.InOrStdin())
            txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(cdc))

            appCode   := args[0]
            tableName := args[1]
            trigger   := args[2]

            msg := types.NewMsgAddTrigger(cliCtx.GetFromAddress(), appCode, tableName, trigger)
            err := msg.ValidateBasic()
            if err != nil {
                return errors.New(fmt.Sprintf("Error %s", err))
            }

            return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
        },
    }
}

func GetCmdDropTrigger(cdc *codec.Codec) *cobra.Command {
    return &cobra.Command{
        Use:   "drop-trigger [appCode] [tableName]",
        Short: "drop a trigger",
        Args:  cobra.ExactArgs(2),
        RunE: func(cmd *cobra.Command, args []string) error {
            cliCtx := context.NewCLIContext().WithCodec(cdc)
            inBuf := bufio.NewReader(cmd.InOrStdin())
            txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(cdc))

            appCode   := args[0]
            tableName := args[1]

            msg := types.NewMsgDropTrigger(cliCtx.GetFromAddress(), appCode, tableName)
            err := msg.ValidateBasic()
            if err != nil {
                return errors.New(fmt.Sprintf("Error %s", err))
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
// modify group member //
//                     //
/////////////////////////

func GetCmdModifyGroupMember(cdc * codec.Codec) *cobra.Command {
    return &cobra.Command{
        Use:   "modify-group-member [appCode] [group] [action] [address]",
        Short: "add/drop account into/from a group",
        Args:  cobra.ExactArgs(4),
        RunE: func(cmd *cobra.Command, args []string) error {
            cliCtx := context.NewCLIContext().WithCodec(cdc)
            inBuf := bufio.NewReader(cmd.InOrStdin())
            txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(cdc))

            appCode       := args[0]
            groupName     := args[1]
            action        := args[2]
            memberAddress := args[3]
            addr, err := sdk.AccAddressFromBech32(memberAddress)

            if err != nil {
                return errors.New(fmt.Sprintf("Error %s", err))
            }

            msg := types.NewMsgModifyGroupMember(appCode, groupName, action, addr, cliCtx.GetFromAddress())
            err = msg.ValidateBasic()
            if err != nil {
                return errors.New(fmt.Sprintf("Error %s", err))
            }
            return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
        },
    }
}

//////////////////
//              //
// Modify group //
//              //
//////////////////

func GetCmdModifyGroup(cdc * codec.Codec) *cobra.Command {
    return &cobra.Command{
        Use:   "modify-group [appCode] [action] [group]",
        Short: "add/drop group for a database",
        Args:  cobra.ExactArgs(3),
        RunE: func(cmd *cobra.Command, args []string) error {
            cliCtx := context.NewCLIContext().WithCodec(cdc)
            inBuf := bufio.NewReader(cmd.InOrStdin())
            txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(cdc))

            appCode   := args[0]
            action    := args[1]
            groupName := args[2]

            msg := types.NewMsgModifyGroup(appCode, action, groupName, cliCtx.GetFromAddress())
            err := msg.ValidateBasic()
            if err != nil {
                return errors.New(fmt.Sprintf("Error %s", err))
            }
            return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
        },
    }
}

////////////////////
//                //
// Set group memo //
//                //
////////////////////

func GetCmdSetGroupMemo(cdc * codec.Codec) *cobra.Command {
    return &cobra.Command{
        Use:   "set-group-memo [appCode] [group] [memo]",
        Short: "set group memo",
        Args:  cobra.ExactArgs(3),
        RunE: func(cmd *cobra.Command, args []string) error {
            cliCtx := context.NewCLIContext().WithCodec(cdc)
            inBuf := bufio.NewReader(cmd.InOrStdin())
            txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(cdc))

            appCode   := args[0]
            groupName := args[1]
            memo      := args[2]

            msg := types.NewMsgSetGroupMemo(appCode, groupName, memo, cliCtx.GetFromAddress())
            err := msg.ValidateBasic()
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

/////////////////
//             //
// drop friend //
//             //
/////////////////

func GetCmdDropFriend(cdc * codec.Codec) *cobra.Command {
    return &cobra.Command{
        Use:   "drop-friend [address]",
        Short: "drop a friend ",
        Args:  cobra.ExactArgs(1),
        RunE: func(cmd *cobra.Command, args []string) error {
            cliCtx := context.NewCLIContext().WithCodec(cdc)
            inBuf := bufio.NewReader(cmd.InOrStdin())
            txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(cdc))

            address   := args[0]
            msg := types.NewMsgDropFriend(cliCtx.GetFromAddress(), address)
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

////////////////////////////
//                        //
// Freeze/Unfreeze schema //
//                        //
////////////////////////////

func GetCmdFreezeSchema(cdc * codec.Codec) *cobra.Command {
    return &cobra.Command{
        Use:   "freeze-schema [database]",
        Short: "Freeze the schma of a database",
        Args:  cobra.ExactArgs(1),
        RunE: func(cmd *cobra.Command, args []string) error {
            cliCtx := context.NewCLIContext().WithCodec(cdc)
            inBuf := bufio.NewReader(cmd.InOrStdin())
            txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(cdc))

            appCode := args[0]
            msg := types.NewMsgSetSchemaStatus(cliCtx.GetFromAddress(), appCode, "frozen")
            err := msg.ValidateBasic()
            if err != nil {
                return errors.New(fmt.Sprintf("Error %s", err))
            }
            return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
        },
    }
}

func GetCmdUnfreezeSchema(cdc * codec.Codec) *cobra.Command {
    return &cobra.Command{
        Use:   "unfreeze-schema [database]",
        Short: "Unfreeze the schma of a database",
        Args:  cobra.ExactArgs(1),
        RunE: func(cmd *cobra.Command, args []string) error {
            cliCtx := context.NewCLIContext().WithCodec(cdc)
            inBuf := bufio.NewReader(cmd.InOrStdin())
            txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(cdc))

            appCode := args[0]
            msg := types.NewMsgSetSchemaStatus(cliCtx.GetFromAddress(), appCode, "unfrozen" )
            err := msg.ValidateBasic()
            if err != nil {
                return errors.New(fmt.Sprintf("Error %s", err))
            }
            return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
        },
    }
}

func GetCmdResetTxTotalTxs(cdc * codec.Codec) *cobra.Command {
    return &cobra.Command{
        Use:   "reset-total-txs",
        Short: "reset total txs statistics which is used for block browser",
        Args:  cobra.ExactArgs(0),
        RunE: func(cmd *cobra.Command, args []string) error {
            cliCtx := context.NewCLIContext().WithCodec(cdc)
            inBuf := bufio.NewReader(cmd.InOrStdin())
            txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(cdc))

            data := map[string]int64 {
                "txNum" : 0,
                "date" : 0,
            }
            bz , _ := json.Marshal(data)
            msg := types.NewMsgUpdateTotalTx(cliCtx.GetFromAddress(), string(bz))
            err := msg.ValidateBasic()
            if err != nil {
                return errors.New(fmt.Sprintf("Error %s", err))
            }
            return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
        },
    }
}

// GetCmdModifyChainSuperAdmins(cdc),
//        GetCmdChangeP2PTransferLimit(cdc),

func GetCmdModifyTokenKeepers(cdc * codec.Codec) *cobra.Command {
    return &cobra.Command{
        Use:   "modify-token-keeper [action] [address]",
        Short: "add or remove token keeper, need two agrs : address, action",
        Args:  cobra.ExactArgs(2),
        RunE: func(cmd *cobra.Command, args []string) error {
            cliCtx := context.NewCLIContext().WithCodec(cdc)
            inBuf := bufio.NewReader(cmd.InOrStdin())
            txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(cdc))

            action := args[0]
            address, err  := sdk.AccAddressFromBech32(args[1])
            if err != nil {
                return errors.New("invalid address")
            }
            msg := types.NewMsgChainModifyTokenKeeperMember(cliCtx.GetFromAddress(), address, action)
            err = msg.ValidateBasic()
            if err != nil {
                return errors.New(fmt.Sprintf("Error %s", err))
            }
            return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
        },
    }
}

func GetCmdModifyP2PTransferLimit(cdc * codec.Codec) *cobra.Command {
    return &cobra.Command{
        Use:   "modify-p2p-transfer-limit",
        Short: "limit p2p transfer or not, set true of false, only token keeper can submit this tx successfully",
        Args:  cobra.ExactArgs(1),
        RunE: func(cmd *cobra.Command, args []string) error {
            cliCtx := context.NewCLIContext().WithCodec(cdc)
            inBuf := bufio.NewReader(cmd.InOrStdin())
            txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(cdc))

            limitStr := args[0]
            var limit = false
            if limitStr == "true" {
                limit = true
            }
            msg := types.NewMsgModifyP2PTransferLimit(cliCtx.GetFromAddress(), limit)
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

