package types

import (
    "github.com/cosmos/cosmos-sdk/codec"
)

// ModuleCdc is the codec for the module
var ModuleCdc = codec.New()

func init() {
    RegisterCodec(ModuleCdc)
}

// RegisterCodec registers concrete types on the Amino codec
func RegisterCodec(cdc *codec.Codec) {
    cdc.RegisterConcrete(MsgCreateTable{}, "cosmosapi/CreateTable", nil)
    cdc.RegisterConcrete(MsgDropTable{}, "cosmosapi/DropTable", nil)
    cdc.RegisterConcrete(MsgAddColumn{}, "cosmosapi/AddColumn", nil)
    cdc.RegisterConcrete(MsgDropColumn{}, "cosmosapi/DropColumn", nil)
    cdc.RegisterConcrete(MsgRenameColumn{}, "cosmosapi/RenameColumn", nil)
    cdc.RegisterConcrete(MsgModifyOption{}, "cosmosapi/ModifyOption", nil)
    cdc.RegisterConcrete(MsgModifyFieldOption{}, "cosmosapi/ModifyFieldOption", nil)
    cdc.RegisterConcrete(MsgCreateIndex{}, "cosmosapi/CreateIndex", nil)
    cdc.RegisterConcrete(MsgInsertRow{}, "cosmosapi/InsertRow", nil)
    cdc.RegisterConcrete(MsgUpdateRow{}, "cosmosapi/UpdateRow", nil)
    cdc.RegisterConcrete(MsgDeleteRow{}, "cosmosapi/DeleteRow", nil)
    cdc.RegisterConcrete(MsgAddAdminAccount{}, "cosmosapi/AddAdminAccount", nil)
}

