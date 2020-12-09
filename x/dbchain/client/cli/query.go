package cli

import (
    "encoding/json"
    "fmt"
    "github.com/mr-tron/base58"
    "github.com/dbchaincloud/tendermint/crypto/sm2"
    sdk "github.com/dbchaincloud/cosmos-sdk/types"
    "github.com/dbchaincloud/cosmos-sdk/client"
    "github.com/dbchaincloud/cosmos-sdk/client/context"
    "github.com/dbchaincloud/cosmos-sdk/client/flags"
    "github.com/dbchaincloud/cosmos-sdk/codec"

    "github.com/spf13/cobra"
    "github.com/dbchaincloud/tendermint/crypto/sm2"
    "github.com/yzhanginwa/dbchain/x/dbchain/client/oracle/oracle"
    "github.com/yzhanginwa/dbchain/x/dbchain/internal/types"
    "strings"
)

func GetQueryCmd(storeKey string, cdc *codec.Codec) *cobra.Command {
    dbchainQueryCmd := &cobra.Command{
        Use:                        types.ModuleName,
        Short:                      "Querying commands for the dbchain module",
        DisableFlagParsing:         true,
        SuggestionsMinimumDistance: 2,
        RunE:                       client.ValidateCmd,
    }
    dbchainQueryCmd.AddCommand(flags.GetCommands(
        GetCmdIsSysAdmin(storeKey, cdc),
        GetCmdApplication(storeKey, cdc),
        GetCmdAppUserFileVolumeLimit(storeKey,cdc),
        GetCmdAppUserUsedFileVolume(storeKey,cdc),
        GetCmdAppUsers(storeKey, cdc),
        GetCmdIsAppUser(storeKey, cdc),
        GetCmdTable(storeKey, cdc),
        GetCmdIndex(storeKey, cdc),
        GetCmdOption(storeKey, cdc),
        GetCmdAssociation(storeKey, cdc),
        GetCmdColumnOption(storeKey, cdc),
        GetCmdColumnDataType(storeKey, cdc),
        GetCmdCanAddColumnOption(storeKey, cdc),
        GetCmdCanSetColumnDataType(storeKey, cdc),
        GetCmdCanInsertRow(storeKey, cdc),
        GetCmdFindRow(storeKey, cdc),
        GetCmdFindIdsBy(storeKey, cdc),
        GetCmdFindAllIds(storeKey, cdc),
        GetCmdShowGroup(storeKey, cdc),
        GetCmdShowGroupMemo(storeKey, cdc),
        GetCmdShowFriends(storeKey, cdc),
        GetCmdShowPendingFriends(storeKey, cdc),
        GetCmdGetAccessCode(storeKey, cdc),
        GetCmdGetOracleInfo(storeKey, cdc),
        GetCmdExportDatabase(storeKey, cdc),
        GetCmdFunction(storeKey,cdc),
        GetCmdFunctionInfo(storeKey,cdc),
        GetCmdCustomQuerier(storeKey,cdc),
        GetCmdCustomQuerierInfo(storeKey,cdc),
        GetCmdCallCustomQuerier(storeKey,cdc),
        GetCmdTxSimpleResult(storeKey,cdc),
        GetCmdChainSuperAdmins(storeKey,cdc),
        GetCmdLimitP2PTransferStatus(storeKey,cdc),
    )...)
    return dbchainQueryCmd
}

func GetCmdIsSysAdmin(queryRoute string, cdc *codec.Codec) *cobra.Command {
    return &cobra.Command{
        Use: "is-sys-admin [access-code]",
        Short: "check whether user is system administrator",
        Args: cobra.ExactArgs(1),
        RunE: func(cmd *cobra.Command, args []string) error {
            cliCtx := context.NewCLIContext().WithCodec(cdc)

            accessCode := args[0]
            res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/is_sys_admin/%s", queryRoute, accessCode), nil)
            if err != nil {
                fmt.Printf("Failed to check whether you are a system administrator")
                return nil
            }

            var out types.QueryOfBoolean
            cdc.MustUnmarshalJSON(res, &out)
            return cliCtx.PrintOutput(out)
        },
    }
}

func GetCmdApplication(queryRoute string, cdc *codec.Codec) *cobra.Command {
    return &cobra.Command{
        Use: "application",
        Short: "query applications",
        Args: cobra.MinimumNArgs(1),
        RunE: func(cmd *cobra.Command, args []string) error {
            cliCtx := context.NewCLIContext().WithCodec(cdc)

            var path string
            // args[0] is the access token.
            if len(args) == 2 {
                path = fmt.Sprintf("custom/%s/application/%s/%s", queryRoute, args[0], args[1])
            } else {
                path = fmt.Sprintf("custom/%s/application/%s", queryRoute, args[0])
            }

            res, _, err := cliCtx.QueryWithData(path, nil)
            if err != nil {
                fmt.Print("could not get applications!")
                return nil
            }

            if len(args) == 2 {
                var out types.Database
                cdc.MustUnmarshalJSON(res, &out)
                return cliCtx.PrintOutput(out)
            } else {
                var out types.QueryTables
                cdc.MustUnmarshalJSON(res, &out)
                return cliCtx.PrintOutput(out)
            }
        },
    }
}

func GetCmdAppUserFileVolumeLimit(queryRoute string, cdc *codec.Codec) *cobra.Command {
    return &cobra.Command{
        Use: "app-user-file-volume-limit",
        Short: "show application user file volume limit",
        Args: cobra.ExactArgs(2),
        RunE: func(cmd *cobra.Command, args []string) error {
            cliCtx := context.NewCLIContext().WithCodec(cdc)

            accessCode := args[0]
            appCode   := args[1]
            res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/application_user_file_volume_limit/%s/%s", queryRoute, accessCode, appCode), nil)
            if err != nil {
                fmt.Printf("could not get users of application %s", appCode)
                return nil
            }

            var out string // QueryTables is a []string. It could be reused here
            cdc.MustUnmarshalJSON(res, &out)
            return cliCtx.PrintOutput(out)
        },
    }
}

func GetCmdAppUserUsedFileVolume(queryRoute string, cdc *codec.Codec) *cobra.Command {
    return &cobra.Command{
        Use: "app-user-used-file-volume",
        Short: "show application user file volume limit",
        Args: cobra.ExactArgs(2),
        RunE: func(cmd *cobra.Command, args []string) error {
            cliCtx := context.NewCLIContext().WithCodec(cdc)

            accessCode := args[0]
            appCode   := args[1]
            res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/application_user_used_file_volume/%s/%s", queryRoute, accessCode, appCode), nil)
            if err != nil {
                fmt.Println("could not get volume of user used")
                return nil
            }

            var out string // QueryTables is a []string. It could be reused here
            cdc.MustUnmarshalJSON(res, &out)
            return cliCtx.PrintOutput(out)
        },
    }
}

func GetCmdAppUsers(queryRoute string, cdc *codec.Codec) *cobra.Command {
    return &cobra.Command{
        Use: "app-users [accessCode] [appCode]",
        Short: "show app users",
        Args: cobra.ExactArgs(2),
        RunE: func(cmd *cobra.Command, args []string) error {
            cliCtx := context.NewCLIContext().WithCodec(cdc)

            accessCode := args[0]
            appCode   := args[1]
            res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/app_users/%s/%s", queryRoute, accessCode, appCode), nil)
            if err != nil {
                fmt.Printf("could not get users of application %s", appCode)
                return nil
            }

            var out types.QueryTables // QueryTables is a []string. It could be reused here
            cdc.MustUnmarshalJSON(res, &out)
            return cliCtx.PrintOutput(out)
        },
    }
}

func GetCmdIsAppUser(queryRoute string, cdc *codec.Codec) *cobra.Command {
    return &cobra.Command{
        Use: "is-app-user",
        Short: "check whether user is allowed for app",
        Args: cobra.ExactArgs(2),
        RunE: func(cmd *cobra.Command, args []string) error {
            cliCtx := context.NewCLIContext().WithCodec(cdc)

            accessCode := args[0]
            appCode   := args[1]
            res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/is_app_user/%s/%s", queryRoute, accessCode, appCode), nil)
            if err != nil {
                fmt.Printf("could not check user of application %s", appCode)
                return nil
            }

            var out types.QueryOfBoolean
            cdc.MustUnmarshalJSON(res, &out)
            return cliCtx.PrintOutput(out)
        },
    }
}

// GetCmdTables lists all table names
func GetCmdTable(queryRoute string, cdc *codec.Codec) *cobra.Command {
    return &cobra.Command{
        Use: "table", //args[0] and args[1] should be accessCode and appCode
        Short: "query tables",
        Args: cobra.MaximumNArgs(3),
        RunE: func(cmd *cobra.Command, args []string) error {
            cliCtx := context.NewCLIContext().WithCodec(cdc)
            var path string
            if len(args) == 3 {
                path = fmt.Sprintf("custom/%s/tables/%s/%s/%s", queryRoute, args[0], args[1], args[2])
            } else if len(args) == 2 {
                path = fmt.Sprintf("custom/%s/tables/%s/%s", queryRoute, args[0], args[1])
            } else {
                fmt.Printf("Need at least 2 parameters!")
                return nil
            }

            res, _, err := cliCtx.QueryWithData(path, nil)
            if err != nil {
                fmt.Printf("could not get table names")
                return nil
            }

            if len(args) == 3 {
                var out types.Table
                cdc.MustUnmarshalJSON(res, &out)
                return cliCtx.PrintOutput(out)
            } else {
                var out types.QueryTables
                cdc.MustUnmarshalJSON(res, &out)
                return cliCtx.PrintOutput(out)
            }
        },
    }
}

func GetCmdIndex(queryRoute string, cdc *codec.Codec) *cobra.Command {
    return &cobra.Command{
        Use: "index [accessCode] [appCode] [tableName]",
        Short: "show index",
        Args: cobra.ExactArgs(3),
        RunE: func(cmd *cobra.Command, args []string) error {
            cliCtx := context.NewCLIContext().WithCodec(cdc)

            accessCode := args[0]
            appCode   := args[1]
            tableName := args[2]
            res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/index/%s/%s/%s", queryRoute, accessCode, appCode, tableName), nil)
            if err != nil {
                fmt.Printf("could not index index of table %s", tableName)
                return nil
            }

            var out types.QueryTables // QueryTables is a []string. It could be reused here
            cdc.MustUnmarshalJSON(res, &out)
            return cliCtx.PrintOutput(out)
        },
    }
}

func GetCmdOption(queryRoute string, cdc *codec.Codec) *cobra.Command {
    return &cobra.Command{
        Use: "table-option [accessCode] [appCode] [tableName]",
        Short: "show table options",
        Args: cobra.ExactArgs(3),
        RunE: func(cmd *cobra.Command, args []string) error {
            cliCtx := context.NewCLIContext().WithCodec(cdc)

            accessCode := args[0]
            appCode    := args[1]
            tableName  := args[2]
            res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/option/%s/%s/%s", queryRoute, accessCode, appCode, tableName), nil)
            if err != nil {
                fmt.Printf("could not get options of table %s", tableName)
                return nil
            }

            var out types.QueryTables // QueryTables is a []string. It could be reused here
            cdc.MustUnmarshalJSON(res, &out)
            return cliCtx.PrintOutput(out)
        },
    }
}

func GetCmdAssociation(queryRoute string, cdc *codec.Codec) *cobra.Command {
    return &cobra.Command{
        Use: "table-association",
        Short: "show table options",
        Args: cobra.ExactArgs(3),
        RunE: func(cmd *cobra.Command, args []string) error {
            cliCtx := context.NewCLIContext().WithCodec(cdc)

            accessCode := args[0]
            appCode    := args[1]
            tableName  := args[2]
            res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/association/%s/%s/%s", queryRoute, accessCode, appCode, tableName), nil)
            if err != nil {
                fmt.Printf("could not get association of table %s", tableName)
                return nil
            }

            var out []types.Association // QueryTables is a []string. It could be reused here
            cdc.MustUnmarshalJSON(res, &out)
            return cliCtx.PrintOutput(out)
        },
    }
}

func GetCmdColumnOption(queryRoute string, cdc *codec.Codec) *cobra.Command {
    return &cobra.Command{
        Use: "column-option [accessCode] [appCode] [tableName] [fieldName]",
        Short: "show column options",
        Args: cobra.ExactArgs(4),
        RunE: func(cmd *cobra.Command, args []string) error {
            cliCtx := context.NewCLIContext().WithCodec(cdc)

            accessCode := args[0]
            appCode    := args[1]
            tableName  := args[2]
            fieldName  := args[3]
            res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/column_option/%s/%s/%s/%s", queryRoute, accessCode, appCode, tableName, fieldName), nil)
            if err != nil {
                fmt.Printf("could not get options of column %s of table %s", fieldName, tableName)
                return nil
            }

            var out types.QueryTables // QueryTables is a []string. It could be reused here
            cdc.MustUnmarshalJSON(res, &out)
            return cliCtx.PrintOutput(out)
        },
    }
}

func GetCmdColumnDataType(queryRoute string, cdc *codec.Codec) *cobra.Command {
    return &cobra.Command{
        Use: "column-data-type",
        Short: "show column data type",
        Args: cobra.ExactArgs(4),
        RunE: func(cmd *cobra.Command, args []string) error {
            cliCtx := context.NewCLIContext().WithCodec(cdc)

            accessCode := args[0]
            appCode    := args[1]
            tableName  := args[2]
            fieldName  := args[3]
            res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/column_data_type/%s/%s/%s/%s", queryRoute, accessCode, appCode, tableName, fieldName), nil)
            if err != nil {
                fmt.Printf("could not get options of column %s of table %s", fieldName, tableName)
                return nil
            }

            var out string
            cdc.MustUnmarshalJSON(res, &out)
            return cliCtx.PrintOutput(out)
        },
    }
}

func GetCmdCanAddColumnOption (queryRoute string, cdc *codec.Codec) *cobra.Command {
    return &cobra.Command{
        Use: "can-add-column-option [accessCode] [appCode] [tableName] [fieldName] [option]",
        Short: "test whether field option can be added",
        Args: cobra.ExactArgs(5),
        RunE: func(cmd *cobra.Command, args []string) error {
            cliCtx := context.NewCLIContext().WithCodec(cdc)

            accessCode := args[0]
            appCode    := args[1]
            tableName  := args[2]
            fieldName  := args[3]
            option     := args[4]

            res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/can_add_column_option/%s/%s/%s/%s/%s", queryRoute, accessCode, appCode, tableName, fieldName, option), nil)
            if err != nil {
                fmt.Printf("Failed to check whether field option can be added")
                return nil
            }

            var out bool
            cdc.MustUnmarshalJSON(res, &out)
            return cliCtx.PrintOutput(out)
        },
    }
}

func GetCmdCanSetColumnDataType (queryRoute string, cdc *codec.Codec) *cobra.Command {
    return &cobra.Command{
        Use: "can-set-column-data-type",
        Short: "test whether field data type can be set",
        Args: cobra.ExactArgs(5),
        RunE: func(cmd *cobra.Command, args []string) error {
            cliCtx := context.NewCLIContext().WithCodec(cdc)

            accessCode := args[0]
            appCode    := args[1]
            tableName  := args[2]
            fieldName  := args[3]
            dataType   := args[4]

            res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/can_set_column_data_type/%s/%s/%s/%s/%s", queryRoute, accessCode, appCode, tableName, fieldName, dataType), nil)
            if err != nil {
                fmt.Printf("Failed to check whether field option can be added")
                return nil
            }

            var out bool
            cdc.MustUnmarshalJSON(res, &out)
            return cliCtx.PrintOutput(out)
        },
    }
}

func GetCmdCanInsertRow(queryRoute string, cdc *codec.Codec) *cobra.Command {
    return &cobra.Command{
        Use: "can-insert-row [accessCode] [appCode] [tableName] [fields] [values]",
        Short: "test whether row can be inserted without violating any column option restricts",
        Args: cobra.ExactArgs(5),
        RunE: func(cmd *cobra.Command, args []string) error {
            cliCtx := context.NewCLIContext().WithCodec(cdc)

            accessCode := args[0]
            appCode    := args[1]
            tableName  := args[2]
            fields     := strings.Split(args[3], ",")
            values     := strings.Split(args[4], ",")

            rowFields := make(types.RowFields)
            for i, field := range fields {
                if i < len(values) {
                    rowFields[field] = values[i]
                }
            }

            rowFieldsJson, err := json.Marshal(rowFields)
            if err != nil {
                fmt.Printf("Failed to marshal to json for rowFields")
                return nil
            }

            res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/can_insert_row/%s/%s/%s/%s", queryRoute, accessCode, appCode, tableName, rowFieldsJson), nil)
            if err != nil {
                fmt.Printf("Failed to check whether row can be inserted")
                return nil
            }

            var out bool
            cdc.MustUnmarshalJSON(res, &out)
            return cliCtx.PrintOutput(out)
        },
    }
}

func GetCmdFindRow(queryRoute string, cdc *codec.Codec) *cobra.Command {
    return &cobra.Command{
        Use: "find [accessCode] [appCode] [tableName] [id]",
        Short: "find row",
        Args: cobra.ExactArgs(4),
        RunE: func(cmd *cobra.Command, args []string) error {
            cliCtx := context.NewCLIContext().WithCodec(cdc)

            accessCode := args[0]
            appCode    := args[1]
            tableName  := args[2]
            id         := args[3]
            res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/find/%s/%s/%s/%s", queryRoute, accessCode, appCode, tableName, id), nil)
            if err != nil {
                fmt.Printf("could not find row")
                return nil
            }

            var out types.QueryRowFields
            cdc.MustUnmarshalJSON(res, &out)
            return cliCtx.PrintOutput(out)
        },
    }
}

func GetCmdFindIdsBy(queryRoute string, cdc *codec.Codec) *cobra.Command {
    return &cobra.Command{
        Use: "find-by [accessCode] [appCode] [tableName] [fieldName] [value]",
        Short: "find by",
        Args: cobra.ExactArgs(5),
        RunE: func(cmd *cobra.Command, args []string) error {
            cliCtx := context.NewCLIContext().WithCodec(cdc)

            accessCode := args[0]
            appCode    := args[1]
            tableName  := args[2]
            fieldName  := args[3]
            value      := args[4]

            // args are accessCode, appCode, tableName, fieldName, and value respectively
            res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/find_by/%s/%s/%s/%s/%s", queryRoute, accessCode, appCode, tableName, fieldName, value), nil)
            if err != nil {
                fmt.Printf("could not find ids")
                return nil
            }

            var out types.QuerySliceOfString
            cdc.MustUnmarshalJSON(res, &out)
            return cliCtx.PrintOutput(out)
        },
    }
}

func GetCmdFindAllIds(queryRoute string, cdc *codec.Codec) *cobra.Command {
    return &cobra.Command{
        Use: "find-all [accessCode] [appCode] [tableName]",
        Short: "find all",
        Args: cobra.ExactArgs(3),
        RunE: func(cmd *cobra.Command, args []string) error {
            cliCtx := context.NewCLIContext().WithCodec(cdc)

            accessCode := args[0]
            appCode    := args[1]
            tableName  := args[2]

            res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/find_all/%s/%s/%s", queryRoute, accessCode, appCode, tableName), nil)
            if err != nil {
                fmt.Printf("could not find ids")
                return nil
            }

            var out types.QuerySliceOfString
            cdc.MustUnmarshalJSON(res, &out)
            return cliCtx.PrintOutput(out)
        },
    }
}

func GetCmdGetOracleInfo(queryRoute string, cdc *codec.Codec) *cobra.Command {
    return &cobra.Command{
        Use: "oracle-info",
        Short: "show oracle info",
        Args: cobra.ExactArgs(0),
        RunE: func(cmd *cobra.Command, args []string) error {
            cliCtx := context.NewCLIContext().WithCodec(cdc)

            privKey, err := oracle.LoadPrivKey()
            if err != nil {
                privKey := sm2.GenPrivKey()
                base58Str := base58.Encode(privKey[:])
                return cliCtx.PrintOutput(fmt.Sprintf("%s: %s", oracle.OracleEncryptedPrivKey, base58Str))
            }
            accAddr := sdk.AccAddress(privKey.PubKey().Address())
            return cliCtx.PrintOutput(fmt.Sprintf("Address: %s", accAddr.String()))
        },
    }
}

func GetCmdExportDatabase (queryRoute string, cdc *codec.Codec) *cobra.Command {
    return &cobra.Command{
        Use: "export-db [appCode]",
        Short: "export database schema",
        Args: cobra.ExactArgs(1),
        RunE: func(cmd *cobra.Command, args []string) error {
            cliCtx := context.NewCLIContext().WithCodec(cdc)

            appCode    := args[0]
            res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/export_database/%s", queryRoute, appCode), nil)
            if err != nil {
                fmt.Printf("Failed to export database")
                return nil
            }

            var out []string
            cdc.MustUnmarshalJSON(res, &out)
            for _, line := range out {
                fmt.Println(line)
            }
            return cliCtx.PrintOutput("")
        },
    }
}


func GetCmdFunction (queryRoute string, cdc *codec.Codec) *cobra.Command {
    return &cobra.Command{
        Use: "functions [accessCode] [appCode]",
        Short: "query functions",
        Args: cobra.ExactArgs(2),
        RunE: func(cmd *cobra.Command, args []string) error {
            cliCtx := context.NewCLIContext().WithCodec(cdc)

            accessCode := args[0]
            appCode    := args[1]
            res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/functions/%s/%s", queryRoute, accessCode, appCode), nil)
            if err != nil {
                fmt.Printf("could not get functions")
                return nil
            }

            var out []string
            cdc.MustUnmarshalJSON(res, &out)
            return cliCtx.PrintOutput(out)
        },
    }
}

func GetCmdFunctionInfo (queryRoute string, cdc *codec.Codec) *cobra.Command {
    return &cobra.Command{
        Use: "function-info [accessCode] [appCode] [functionName]",
        Short: "query function specific information",
        Args: cobra.ExactArgs(3),
        RunE: func(cmd *cobra.Command, args []string) error {
            cliCtx := context.NewCLIContext().WithCodec(cdc)

            accessCode   := args[0]
            appCode      := args[1]
            functionName := args[2]
            res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/functionInfo/%s/%s/%s", queryRoute, accessCode, appCode, functionName), nil)
            if err != nil {
                fmt.Printf("could not get functionInfo")
                return nil
            }

            var out types.Function
            cdc.MustUnmarshalJSON(res, &out)
            return cliCtx.PrintOutput(out)
        },
    }
}

func GetCmdCustomQuerier (queryRoute string, cdc *codec.Codec) *cobra.Command {
    return &cobra.Command{
        Use: "custom-queriers [accessCode] [appCode]",
        Short: "query queriers",
        Args: cobra.ExactArgs(2),
        RunE: func(cmd *cobra.Command, args []string) error {
            cliCtx := context.NewCLIContext().WithCodec(cdc)

            accessCode := args[0]
            appCode    := args[1]
            res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/customQueriers/%s/%s", queryRoute, accessCode, appCode), nil)
            if err != nil {
                fmt.Printf("could not get queriers")
                return nil
            }

            var out []string
            cdc.MustUnmarshalJSON(res, &out)
            return cliCtx.PrintOutput(out)
        },
    }
}

func GetCmdCustomQuerierInfo (queryRoute string, cdc *codec.Codec) *cobra.Command {
    return &cobra.Command{
        Use: "custom-querier-info [accessCode] [appCode] [functionName]",
        Short: "query querier specific information",
        Args: cobra.ExactArgs(3),
        RunE: func(cmd *cobra.Command, args []string) error {
            cliCtx := context.NewCLIContext().WithCodec(cdc)

            accessCode   := args[0]
            appCode      := args[1]
            functionName := args[2]
            res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/customQuerierInfo/%s/%s/%s", queryRoute, accessCode, appCode, functionName), nil)
            if err != nil {
                fmt.Printf("could not get querierInfo")
                return nil
            }

            var out types.Function
            cdc.MustUnmarshalJSON(res, &out)
            return cliCtx.PrintOutput(out)
        },
    }
}

func GetCmdTxSimpleResult (queryRoute string, cdc *codec.Codec) *cobra.Command {
    return &cobra.Command{
        Use: "tx-simple-result [accessCode] [txHash]",
        Short: "query whether the tx is successful",
        Args: cobra.ExactArgs(2),
        RunE: func(cmd *cobra.Command, args []string) error {
            cliCtx := context.NewCLIContext().WithCodec(cdc)

            accessCode   := args[0]
            txHash      := args[1]
            res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/txSimpleResult/%s/%s", queryRoute, accessCode, txHash), nil)

            if err != nil {
                fmt.Printf("could not get querierInfo")
                return nil
            }

            var out *types.TxStatus
            cdc.MustUnmarshalJSON(res, &out)
            return cliCtx.PrintOutput(out)
        },
    }
}

func GetCmdCallCustomQuerier (queryRoute string, cdc *codec.Codec) *cobra.Command {
    return &cobra.Command{
        Use: "call-custom-querier [accessCode] [appCode] [querierName] [params]",
        Short: "call custom querier to query data",
        Args: cobra.MinimumNArgs(4),
        RunE: func(cmd *cobra.Command, args []string) error {
            cliCtx := context.NewCLIContext().WithCodec(cdc)

            accessCode   := args[0]
            appCode      := args[1]
            querierName  := args[2]
            params       := args[3]

            res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/callCustomQuerier/%s/%s/%s/%s", queryRoute, accessCode, appCode, querierName, params), nil)
            if err != nil {
                fmt.Printf("could not get data")
                return nil
            }
            //TODO What kind of format is needed here
            var out = unmarshalCustomData(res, cdc)
            if out == nil {
                fmt.Println("invalid lua res data")
                return nil
            }
            return cliCtx.PrintOutput(out)
        },
    }
}

func GetCmdChainSuperAdmins (queryRoute string, cdc *codec.Codec) *cobra.Command {
    return &cobra.Command{
        Use: "chain-super-admins [accessCode]",
        Short: "query all admins. only admin can get data",
        Args: cobra.MinimumNArgs(1),
        RunE: func(cmd *cobra.Command, args []string) error {
            cliCtx := context.NewCLIContext().WithCodec(cdc)

            accessCode   := args[0]

            res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/chain_super_admins/%s", queryRoute, accessCode), nil)
            if err != nil {
                fmt.Printf("could not get data")
                return nil
            }
            //TODO What kind of format is needed here
            admins := make([]string, 0)
            cdc.MustUnmarshalJSON(res, &admins)
            return cliCtx.PrintOutput(admins)
        },
    }
}

func GetCmdLimitP2PTransferStatus (queryRoute string, cdc *codec.Codec) *cobra.Command {
    return &cobra.Command{
        Use: "limit-p2p-transfer-status [accessCode]",
        Short: "get current limit p2p transfer status",
        Args: cobra.MinimumNArgs(1),
        RunE: func(cmd *cobra.Command, args []string) error {
            cliCtx := context.NewCLIContext().WithCodec(cdc)

            accessCode   := args[0]

            res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/limit_p2p_transfer_status/%s", queryRoute, accessCode), nil)
            if err != nil {
                fmt.Printf("could not get data")
                return nil
            }
            //TODO What kind of format is needed here
            var limit bool
            cdc.MustUnmarshalJSON(res, &limit)
            return cliCtx.PrintOutput(limit)
        },
    }
}

func unmarshalCustomData(bz []byte, cdc *codec.Codec)interface{}{
    res1 := make([]uint, 0)
    err := cdc.UnmarshalJSON(bz, &res1)
    if err == nil{
        return res1
    }
    res2 := make(map[string]string)
    err = cdc.UnmarshalJSON(bz, &res2)
    if err == nil{
        return res2
    }
    res3 := make([]map[string]string, 0)
    err = cdc.UnmarshalJSON(bz, &res3)
    if err == nil{
        return res3
    }
    return nil
}
