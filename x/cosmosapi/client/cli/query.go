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
    )...)
    return cosmosapiQueryCmd
}

func GetCmdApplication(queryRoute string, cdc *codec.Codec) *cobra.Command {
    return &cobra.Command{
        Use: "application",
        Short: "query applications",
        Args: cobra.MaximumNArgs(1),
        RunE: func(cmd *cobra.Command, args []string) error {
            cliCtx := context.NewCLIContext().WithCodec(cdc)

            var path string
            if len(args) == 1 {
                path = fmt.Sprintf("custom/%s/application/%s", queryRoute, args[0])
            } else {
                path = fmt.Sprintf("custom/%s/application", queryRoute)
            }

            res, _, err := cliCtx.QueryWithData(path, nil)
            if err != nil {
                fmt.Print("could not get applications!")
                return nil
            }

            if len(args) == 1 {
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
        Args: cobra.MaximumNArgs(1),
        RunE: func(cmd *cobra.Command, args []string) error {
            cliCtx := context.NewCLIContext().WithCodec(cdc)
            var path string
            if len(args) == 1 {
                path = fmt.Sprintf("custom/%s/tables/%s", queryRoute, args[0])
            } else {
                path = fmt.Sprintf("custom/%s/tables", queryRoute)
            }

            res, _, err := cliCtx.QueryWithData(path, nil)
            if err != nil {
                fmt.Printf("could not get table names")
                return nil
            }

            if len(args) == 1 {
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
        Args: cobra.ExactArgs(1),
        RunE: func(cmd *cobra.Command, args []string) error {
            cliCtx := context.NewCLIContext().WithCodec(cdc)

            res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/index/%s", queryRoute, args[0]), nil)
            if err != nil {
                fmt.Printf("could not index index of table %s", args[0])
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
        Args: cobra.ExactArgs(1),
        RunE: func(cmd *cobra.Command, args []string) error {
            cliCtx := context.NewCLIContext().WithCodec(cdc)

            // args[0] is table name
            res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/option/%s", queryRoute, args[0]), nil)
            if err != nil {
                fmt.Printf("could not get options of table %s", args[0])
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
        Args: cobra.ExactArgs(2),
        RunE: func(cmd *cobra.Command, args []string) error {
            cliCtx := context.NewCLIContext().WithCodec(cdc)

            // args[0] is table name
            // args[1] is field name
            res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/column_option/%s/%s", queryRoute, args[0], args[1]), nil)
            if err != nil {
                fmt.Printf("could not get options of column %s of table %s", args[1], args[0])
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
        Args: cobra.ExactArgs(2),
        RunE: func(cmd *cobra.Command, args []string) error {
            cliCtx := context.NewCLIContext().WithCodec(cdc)

            res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/find/%s/%s", queryRoute, args[0], args[1]), nil)
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
        Args: cobra.ExactArgs(3),
        RunE: func(cmd *cobra.Command, args []string) error {
            cliCtx := context.NewCLIContext().WithCodec(cdc)

            // args are tableName, fieldName, and value respectively
            res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/find_by/%s/%s/%s", queryRoute, args[0], args[1], args[2]), nil)
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
        Args: cobra.ExactArgs(1),
        RunE: func(cmd *cobra.Command, args []string) error {
            cliCtx := context.NewCLIContext().WithCodec(cdc)

            // args are tableName
            res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/find_all/%s", queryRoute, args[0]), nil)
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
