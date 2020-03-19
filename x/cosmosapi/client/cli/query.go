package cli

import (
    "fmt"
    "github.com/cosmos/cosmos-sdk/client"
    "github.com/cosmos/cosmos-sdk/client/context"
    "github.com/cosmos/cosmos-sdk/codec"
    "github.com/yzhanginwa/cosmos-api/x/cosmosapi/internal/types"
    "github.com/spf13/cobra"
)

func GetQueryCmd(storeKey string, cdc *codec.Codec) *cobra.Command {
    cosmosapiQueryCmd := &cobra.Command{
        Use:                        types.ModuleName,
        Short:                      "Querying commands for the cosmosapi module",
        DisableFlagParsing:         true,
        SuggestionsMinimumDistance: 2,
        RunE:                       client.ValidateCmd,
    }
    cosmosapiQueryCmd.AddCommand(client.GetCommands(
        GetCmdApplication(storeKey, cdc),
        GetCmdTable(storeKey, cdc),
        GetCmdIndex(storeKey, cdc),
        GetCmdOption(storeKey, cdc),
        GetCmdColumnOption(storeKey, cdc),
        GetCmdFindRow(storeKey, cdc),
        GetCmdFindIdsBy(storeKey, cdc),
        GetCmdFindAllIds(storeKey, cdc),
        GetCmdShowAdminGroup(storeKey, cdc),
        GetCmdGetAccessCode(storeKey, cdc),
    )...)
    return cosmosapiQueryCmd
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
            res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/find/%s/%s/%s/%", queryRoute, accessCode, appCode, tableName, id), nil)
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
