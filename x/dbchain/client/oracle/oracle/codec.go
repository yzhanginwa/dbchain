package oracle

import (
    amino "github.com/tendermint/go-amino"
    cryptoamino "github.com/dbchaincloud/tendermint/crypto/encoding/amino"
    sdk "github.com/dbchaincloud/cosmos-sdk/types"
    "github.com/dbchaincloud/cosmos-sdk/x/auth/exported"
    authtypes "github.com/dbchaincloud/cosmos-sdk/x/auth/types"
    "github.com/yzhanginwa/dbchain/x/dbchain/internal/types"
)

var aminoCdc = amino.NewCodec()

func init () {
    aminoCdc.RegisterInterface((*sdk.Msg)(nil), nil)
    aminoCdc.RegisterInterface((*sdk.Tx)(nil), nil)
    aminoCdc.RegisterInterface((*UniversalMsg)(nil), nil)
    aminoCdc.RegisterConcrete(types.MsgInsertRow{}, "dbchain/InsertRow", nil)
    aminoCdc.RegisterConcrete(types.MsgUpdateTotalTx{}, "dbchain/UpdateTotalTx", nil)
    aminoCdc.RegisterConcrete(types.MsgUpdateTxStatistic{}, "dbchain/UpdateTxStatistic", nil)
    aminoCdc.RegisterConcrete(types.MsgFreezeRow{}, "dbchain/FreezeRow", nil)
    aminoCdc.RegisterConcrete(types.MsgSaveUserPrivateKey{}, "dbchain/SaveUserPrivateKey", nil)
    aminoCdc.RegisterConcrete(types.MsgCallFunction{}, "dbchain/CallFunction", nil)
    cryptoamino.RegisterAmino(aminoCdc)

    //authtypes.RegisterCodec(aminoCdc)
    aminoCdc.RegisterInterface((*exported.GenesisAccount)(nil), nil)
    aminoCdc.RegisterInterface((*exported.Account)(nil), nil)
    aminoCdc.RegisterConcrete(&authtypes.BaseAccount{}, "cosmos-sdk/Account", nil)
    aminoCdc.RegisterConcrete(StdTx{}, "cosmos-sdk/StdTx", nil)

    aminoCdc.RegisterConcrete(MsgSend{}, "cosmos-sdk/MsgSend", nil)
}

