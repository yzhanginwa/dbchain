package cli

import (
    "errors"
    "strconv"
    "time"
    "bufio"

    "github.com/spf13/cobra"
    "github.com/cosmos/cosmos-sdk/client/context"
    "github.com/cosmos/cosmos-sdk/codec"
    "github.com/cosmos/cosmos-sdk/client/input"
    "github.com/cosmos/cosmos-sdk/client/keys"
    "github.com/mr-tron/base58"
    "github.com/yzhanginwa/cosmos-api/x/cosmosapi/internal/types"
    "github.com/tendermint/tendermint/crypto/secp256k1"
)

func GetCmdGetAccessCode(queryRoute string, cdc *codec.Codec) *cobra.Command {
    return &cobra.Command{
        Use: "access-code",
        Short: "show access code",
        Args: cobra.ExactArgs(1),
        RunE: func(cmd *cobra.Command, args []string) error {
            cliCtx := context.NewCLIContext().WithCodec(cdc)

            name := args[0]

            kb, err := keys.NewKeyBaseFromHomeFlag()
            if err != nil {
                    return err
            }
    
            buf := bufio.NewReader(cmd.InOrStdin())
            passphrase, err := input.GetPassword("Enter passphrase to decrypt your key:", buf)
            if err != nil {
                    return err
            }
    
            now := time.Now().UnixNano() / 1000000
            nowString := strconv.Itoa(int(now)) 

            signature, pubKey, err := kb.Sign(name, passphrase, []byte(nowString))

            if err != nil {
                return err
            }

            if pk, ok := pubKey.(secp256k1.PubKeySecp256k1); ok {
                out := base58.Encode(pk[:]) + ":" + nowString + ":" + base58.Encode(signature)
                return cliCtx.PrintOutput(types.QueryOfString(out))
            } else {
                return errors.New("Failed to parse public key")
            }
        },
    }
}

