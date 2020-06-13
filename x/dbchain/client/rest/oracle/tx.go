package oracle

import (
    "fmt"
    sdk "github.com/cosmos/cosmos-sdk/types"
    authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
    amino "github.com/tendermint/go-amino"
    cryptoamino "github.com/tendermint/tendermint/crypto/encoding/amino"

    rpcclient "github.com/tendermint/tendermint/rpc/client"
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
    authtypes.RegisterCodec(aminoCdc)
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
   
    //fmt.Printf("\nAccount number: %d, seq: %d\n\n", accNum, seq)

    msgs := []sdk.Msg{msg}
    stdFee := authtypes.NewStdFee(200000, sdk.Coins{sdk.NewCoin("dbctoken", sdk.NewInt(1))})

    stdSignMsg := authtypes.StdSignMsg{
        ChainID:       "testnet",
        AccountNumber: accNum,
        Sequence:      seq,
        Memo:          "",
        Msgs:          msgs,
        Fee:           stdFee,
    }

    sig, err := privKey.Sign(stdSignMsg.Bytes())
    if err != nil {
        fmt.Println("Oracle: Failed to sign message!!!")
        return
    }

    stdSignature :=authtypes.StdSignature{
        PubKey:    privKey.PubKey(),
        Signature: sig,
    }

    newStdTx := authtypes.NewStdTx(msgs, stdFee, []authtypes.StdSignature{stdSignature}, "")

    encoder := authtypes.DefaultTxEncoder(aminoCdc)
    txBytes, err := encoder(newStdTx)
    if err != nil {
        fmt.Println("Oracle: Failed to marshal StdTx!!!")
        return
    }

    //cliCtx.BroadcastTxAsync(txBytes)
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
