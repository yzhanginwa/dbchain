package cli

import (
    "fmt"
    "github.com/mr-tron/base58"
    "github.com/tendermint/tendermint/crypto/secp256k1"
    sdk "github.com/cosmos/cosmos-sdk/types"
    "github.com/cosmos/cosmos-sdk/client"
    "github.com/cosmos/cosmos-sdk/client/context"
    "github.com/cosmos/cosmos-sdk/client/flags"
    "github.com/cosmos/cosmos-sdk/codec"
    "github.com/yzhanginwa/dbchain/x/dbchain/client/rest/oracle"
    "github.com/yzhanginwa/dbchain/x/dbchain/internal/types"
    "github.com/spf13/cobra"
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
        GetCmdAppUsers(storeKey, cdc),
        GetCmdIsAppUser(storeKey, cdc),
        GetCmdTable(storeKey, cdc),
        GetCmdIndex(storeKey, cdc),
        GetCmdOption(storeKey, cdc),
        GetCmdColumnOption(storeKey, cdc),
        GetCmdCanAddColumnOption(storeKey, cdc),
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
    )...)
    return dbchainQueryCmd
}

func GetCmdIsSysAdmin(queryRoute string, cdc *codec.Codec) *cobra.Command {
    return &cobra.Command{
        Use: "is-sys-admin",
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

func GetCmdAppUsers(queryRoute string, cdc *codec.Codec) *cobra.Command {
    return &cobra.Command{
        Use: "app-users",
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
        Use: "table",
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
        Use: "index",
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
        Use: "table-option",
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

func GetCmdColumnOption(queryRoute string, cdc *codec.Codec) *cobra.Command {
    return &cobra.Command{
        Use: "column-option",
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

func GetCmdCanAddColumnOption (queryRoute string, cdc *codec.Codec) *cobra.Command {
    return &cobra.Command{
        Use: "can-add-column-option",
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

func GetCmdFindRow(queryRoute string, cdc *codec.Codec) *cobra.Command {
    return &cobra.Command{
        Use: "find",
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
        Use: "find-by",
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
            res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/find_by/%s/%s/%s", queryRoute, accessCode, appCode, tableName, fieldName, value), nil)
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
        Use: "find-all",
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
                privKey := secp256k1.GenPrivKey()
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
        Use: "export-db",
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
