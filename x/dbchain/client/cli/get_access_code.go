package cli

import (
    "errors"
    "strconv"
    "time"
    "github.com/spf13/cobra"
    "github.com/spf13/viper"
    sdk "github.com/cosmos/cosmos-sdk/types"
    "github.com/cosmos/cosmos-sdk/client/context"
    "github.com/cosmos/cosmos-sdk/codec"
    "github.com/cosmos/cosmos-sdk/client/flags"
    "github.com/cosmos/cosmos-sdk/client/keys"
    cryptoKeys "github.com/cosmos/cosmos-sdk/crypto/keys"
    "github.com/mr-tron/base58"
    "github.com/yzhanginwa/cosmos-api/x/dbchain/internal/types"
    "github.com/tendermint/tendermint/crypto/secp256k1"
)

func GetCmdGetAccessCode(queryRoute string, cdc *codec.Codec) *cobra.Command {
    resultCmd := &cobra.Command{
        Use: "access-code",
        Short: "show access code",
        Args: cobra.MinimumNArgs(1),
        RunE: func(cmd *cobra.Command, args []string) error {
            cliCtx := context.NewCLIContext().WithCodec(cdc)
            kb, err := cryptoKeys.NewKeyring(sdk.KeyringServiceName(), viper.GetString(flags.FlagKeyringBackend), viper.GetString(flags.FlagHome), cmd.InOrStdin())
            if err != nil {
                return err
            }

            name := args[0]
            var str string

            if len(args) > 1 {
                str = args[1]
            } else {
                now := time.Now().UnixNano() / 1000000
                str = strconv.Itoa(int(now))
            }

            if out, ok := signForToken(kb, name, str); ok {
                return cliCtx.PrintOutput(types.QueryOfString(out))
            } else {
                return errors.New("Failed to parse public key")
            }
        },
    }

    // borrowed from github.com/cosmos/cosmos-sdk/client/keys/root.go
    resultCmd.PersistentFlags().String(flags.FlagKeyringBackend, flags.DefaultKeyringBackend, "Select keyring's backend (os|file|test)")
    viper.BindPFlag(flags.FlagKeyringBackend, resultCmd.Flags().Lookup(flags.FlagKeyringBackend))

    return resultCmd
}

func signForToken(kb cryptoKeys.Keybase, name string, str string) (string, bool) {
    signature, pubKey, err := kb.Sign(name, keys.DefaultKeyPass, []byte(str))
    if err != nil {
        return "", false
    }

    if pk, ok := pubKey.(secp256k1.PubKeySecp256k1); ok {
        out := base58.Encode(pk[:]) + ":" + str + ":" + base58.Encode(signature)
        return out, true
    } else {
        return "", false
    }
}
