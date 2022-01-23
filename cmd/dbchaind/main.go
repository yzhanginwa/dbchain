package main

import (
    "encoding/json"
    "io"

    "github.com/cosmos/cosmos-sdk/server"
    "github.com/cosmos/cosmos-sdk/x/staking"

    "github.com/spf13/cobra"
    "github.com/tendermint/tendermint/libs/cli"
    "github.com/tendermint/tendermint/libs/log"

    sdk "github.com/cosmos/cosmos-sdk/types"
    "github.com/cosmos/cosmos-sdk/x/auth"
    genutilcli "github.com/cosmos/cosmos-sdk/x/genutil/client/cli"
    app "github.com/yzhanginwa/dbchain"

    //bankmodule "github.com/yzhanginwa/dbchain/x/bank"
    bankmodule "github.com/cosmos/cosmos-sdk/x/bank"

    dbcmodule "github.com/yzhanginwa/dbchain/x/dbchain"
    dbchaincli "github.com/yzhanginwa/dbchain/x/dbchain/client/cli"

    abci "github.com/tendermint/tendermint/abci/types"
    tmtypes "github.com/tendermint/tendermint/types"
    dbm "github.com/tendermint/tm-db"
)

func main() {
    go statusReport()
    go dbcmodule.TxCacheInvalid()
    cobra.EnableCommandSorting = false

    cdc := app.MakeCodec()

    config := sdk.GetConfig()
    config.SetBech32PrefixForAccount(app.Bech32PrefixAccAddr, app.Bech32PrefixAccPub)
    config.SetBech32PrefixForValidator(app.Bech32PrefixValAddr, app.Bech32PrefixValPub)
    config.SetBech32PrefixForConsensusNode(app.Bech32PrefixConsAddr, app.Bech32PrefixConsPub)
    config.Seal()

    ctx := server.NewDefaultContext()

    rootCmd := &cobra.Command{
        Use:               "dbchaind",
        Short:             "dbchain App Daemon (server)",
        PersistentPreRunE: server.PersistentPreRunEFn(ctx),
    }

    rootCmd.PersistentFlags().BoolVar(&dbcmodule.AllowCreateApplication,
                                      "allow-create-application",
                                      false,
                                      "allow non-admin users to create application")

    rootCmd.PersistentFlags().Int64Var(&bankmodule.ExistentialDeposit,
                                        "existential-deposit",
                                        0,
                                        "deposits with the least per address")

    // CLI commands to initialize the chain
    rootCmd.AddCommand(
        genutilcli.InitCmd(ctx, cdc, app.ModuleBasics, app.DefaultNodeHome),
        genutilcli.CollectGenTxsCmd(ctx, cdc, auth.GenesisAccountIterator{}, app.DefaultNodeHome),
        genutilcli.GenTxCmd(
            ctx, cdc, app.ModuleBasics, staking.AppModuleBasic{},
            auth.GenesisAccountIterator{}, app.DefaultNodeHome, app.DefaultCLIHome,
        ),
        genutilcli.ValidateGenesisCmd(ctx, cdc, app.ModuleBasics),
        // AddGenesisAccountCmd allows users to add accounts to the genesis file
        AddGenesisAccountCmd(ctx, cdc, app.DefaultNodeHome, app.DefaultCLIHome),
        dbchaincli.AddGenesisAdminAccountCmd(ctx, cdc, app.DefaultNodeHome, app.DefaultCLIHome),
    )

    server.AddCommands(ctx, cdc, rootCmd, newApp, exportAppStateAndTMValidators)

    // prepare and add flags
    executor := cli.PrepareBaseCmd(rootCmd, "NS", app.DefaultNodeHome)
    err := executor.Execute()
    if err != nil {
        panic(err)
    }
}

func newApp(logger log.Logger, db dbm.DB, traceStore io.Writer) abci.Application {
    return app.NewDbChainApp(logger, db)
}

func exportAppStateAndTMValidators(
    logger log.Logger, db dbm.DB, traceStore io.Writer, height int64, forZeroHeight bool, jailWhiteList []string,
) (json.RawMessage, []tmtypes.GenesisValidator, error) {

    if height != -1 {
        nsApp := app.NewDbChainApp(logger, db)
        err := nsApp.LoadHeight(height)
        if err != nil {
            return nil, nil, err
        }
        return nsApp.ExportAppStateAndValidators(forZeroHeight, jailWhiteList)
    }

    nsApp := app.NewDbChainApp(logger, db)

    return nsApp.ExportAppStateAndValidators(forZeroHeight, jailWhiteList)
}
