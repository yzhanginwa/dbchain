package main

import (
    "github.com/yzhanginwa/dbchain/x/dbchain/client/oracle"
    "os"
    "path"

    "github.com/cosmos/cosmos-sdk/client"
    "github.com/cosmos/cosmos-sdk/client/flags"
    sdk "github.com/cosmos/cosmos-sdk/types"
    "github.com/cosmos/cosmos-sdk/version"
    authrest "github.com/cosmos/cosmos-sdk/x/auth/client/rest"
    "github.com/spf13/cobra"
    "github.com/spf13/viper"
    "github.com/tendermint/tendermint/libs/cli"
    app "github.com/yzhanginwa/dbchain"
)

func main() {
    cobra.EnableCommandSorting = false

    cdc := app.MakeCodec()

    // Read in the configuration file for the sdk
    config := sdk.GetConfig()
    config.SetBech32PrefixForAccount(sdk.Bech32PrefixAccAddr, sdk.Bech32PrefixAccPub)
    config.SetBech32PrefixForValidator(sdk.Bech32PrefixValAddr, sdk.Bech32PrefixValPub)
    config.SetBech32PrefixForConsensusNode(sdk.Bech32PrefixConsAddr, sdk.Bech32PrefixConsPub)
    config.Seal()

    rootCmd := &cobra.Command{
        Use:   "dbchainoracle",
        Short: "dbchainoracle Client",
    }

    // Add --chain-id to persistent flags and mark it required
    rootCmd.PersistentFlags().String(flags.FlagChainID, "", "Chain ID of tendermint node")
    rootCmd.PersistentPreRunE = func(_ *cobra.Command, _ []string) error {
        return initConfig(rootCmd)
    }

    // Construct Root Command
    rootCmd.AddCommand(
        client.ConfigCmd(app.DefaultOracleHome),
        oracle.ServeCommand(cdc, registerRoutes),//use oracle server. its no limit to req.body
        version.Cmd,
        flags.NewCompletionCmd(rootCmd, true),
    )

    executor := cli.PrepareMainCmd(rootCmd, "NS", app.DefaultOracleHome)
    err := executor.Execute()
    if err != nil {
        panic(err)
    }
}

func registerRoutes(rs *oracle.RestServer) {
    client.RegisterRoutes(rs.CliCtx, rs.Mux)
    authrest.RegisterTxRoutes(rs.CliCtx, rs.Mux)
    app.ModuleBasics.RegisterRESTRoutes(rs.CliCtx, rs.Mux)
}

func initConfig(cmd *cobra.Command) error {
    home, err := cmd.PersistentFlags().GetString(cli.HomeFlag)
    if err != nil {
        return err
    }

    cfgFile := path.Join(home, "config", "config.toml")
    if _, err := os.Stat(cfgFile); err == nil {
        viper.SetConfigFile(cfgFile)

        if err := viper.ReadInConfig(); err != nil {
            return err
        }
    }
    if err := viper.BindPFlag(flags.FlagChainID, cmd.PersistentFlags().Lookup(flags.FlagChainID)); err != nil {
        return err
    }
    if err := viper.BindPFlag(cli.EncodingFlag, cmd.PersistentFlags().Lookup(cli.EncodingFlag)); err != nil {
        return err
    }
    return viper.BindPFlag(cli.OutputFlag, cmd.PersistentFlags().Lookup(cli.OutputFlag))
}
