package oracle

import (
    "fmt"
    "github.com/spf13/viper"
    amino "github.com/tendermint/go-amino"
    cryptoamino "github.com/tendermint/tendermint/crypto/encoding/amino"
    rpcclient "github.com/tendermint/tendermint/rpc/client"
    "github.com/tendermint/tendermint/crypto/secp256k1"
    "github.com/cosmos/cosmos-sdk/x/auth/exported"
    sdk "github.com/cosmos/cosmos-sdk/types"
    authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
    "github.com/yzhanginwa/dbchain/x/dbchain/internal/types"
)

var (
    aminoCdc = amino.NewCodec()
)

func init () {
    aminoCdc.RegisterInterface((*sdk.Msg)(nil), nil)
    aminoCdc.RegisterInterface((*sdk.Tx)(nil), nil)
    aminoCdc.RegisterConcrete(types.MsgInsertRow{}, "dbchain/InsertRow", nil)
    cryptoamino.RegisterAmino(aminoCdc)

    //authtypes.RegisterCodec(aminoCdc)
    aminoCdc.RegisterInterface((*exported.GenesisAccount)(nil), nil)
    aminoCdc.RegisterInterface((*exported.Account)(nil), nil)
    aminoCdc.RegisterConcrete(&authtypes.BaseAccount{}, "cosmos-sdk/Account", nil)
    aminoCdc.RegisterConcrete(StdTx{}, "cosmos-sdk/StdTx", nil)
}

func buildTxAndBroadcast(msg sdk.Msg) {
    privKey, err := LoadPrivKey()
    if err != nil {
        fmt.Println("Failed to load oracle's private key!!!")
        return
    }
    oracleAccAddr := sdk.AccAddress(privKey.PubKey().Address())
    accNum, seq, err := getAccountInfo(oracleAccAddr.String())
    if err != nil {
        fmt.Println("Failed to load oracle's account info!!!")
        return
    }

    txBytes, err := buildAndSignAndBuildTxBytes(msg, accNum, seq, privKey)
    if err != nil {
        return
    }
    broadcastTxBytes(txBytes)
}

func buildAndSignAndBuildTxBytes(msg sdk.Msg, accNum uint64, seq uint64, privKey secp256k1.PrivKeySecp256k1) ([]byte, error) {
    msgs := []sdk.Msg{msg}
    stdFee := NewStdFee(200000, sdk.Coins{sdk.NewCoin("dbctoken", sdk.NewInt(1))})
    chainId := viper.GetString("chain-id")
    stdSignMsgBytes := StdSignBytes(chainId, accNum, seq, stdFee, msgs, "")

    sig, err := privKey.Sign(stdSignMsgBytes)

    if err != nil {
        fmt.Println("Oracle: Failed to sign message!!!")
        return nil, err
    }

    stdSignature := StdSignature {
        PubKey:    privKey.PubKey(),
        Signature: sig,
    }

    newStdTx := NewStdTx(msgs, stdFee, []StdSignature{stdSignature}, "")
    txBytes, err := aminoCdc.MarshalBinaryLengthPrefixed(newStdTx)
    if err != nil {
        fmt.Println("Oracle: Failed to marshal StdTx!!!")
        return nil, err
    }

    return txBytes, nil
}

func broadcastTxBytes(txBytes []byte) {
    rpc, err := rpcclient.NewHTTP("http://localhost:26657", "/websocket")
    if err != nil {
        fmt.Printf("failted to get client: %v\n", err)
        return
    }

    _, err = rpc.BroadcastTxAsync(txBytes)
    if err != nil {
        fmt.Printf("failted to broadcast transaction: %v\n", err)
        return
    }
}
