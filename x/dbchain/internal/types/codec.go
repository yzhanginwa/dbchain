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
    cdc.RegisterConcrete(MsgCreateApplication{}, "dbchain/CreateApplication", nil)
    cdc.RegisterConcrete(MsgCreateSysDatabase{}, "dbchain/CreateSysDatabase", nil)
    cdc.RegisterConcrete(MsgSetSchemaStatus{}, "dbchain/MsgSetSchemaStatus", nil)
    cdc.RegisterConcrete(MsgAddDatabaseUser{}, "dbchain/AddDatabaseUser", nil)
    cdc.RegisterConcrete(MsgCreateTable{}, "dbchain/CreateTable", nil)
    cdc.RegisterConcrete(MsgDropTable{}, "dbchain/DropTable", nil)
    cdc.RegisterConcrete(MsgAddColumn{}, "dbchain/AddColumn", nil)
    cdc.RegisterConcrete(MsgDropColumn{}, "dbchain/DropColumn", nil)
    cdc.RegisterConcrete(MsgRenameColumn{}, "dbchain/RenameColumn", nil)
    cdc.RegisterConcrete(MsgModifyOption{}, "dbchain/ModifyOption", nil)
    cdc.RegisterConcrete(MsgModifyColumnOption{}, "dbchain/ModifyColumnOption", nil)
    cdc.RegisterConcrete(MsgCreateIndex{}, "dbchain/CreateIndex", nil)
    cdc.RegisterConcrete(MsgDropIndex{}, "dbchain/DropIndex", nil)
    cdc.RegisterConcrete(MsgAddInsertFilter{}, "dbchain/AddInsertFilter", nil)
    cdc.RegisterConcrete(MsgDropInsertFilter{}, "dbchain/DropInsertFilter", nil)
    cdc.RegisterConcrete(MsgAddTrigger{}, "dbchain/AddTrigger", nil)
    cdc.RegisterConcrete(MsgDropTrigger{}, "dbchain/DropTrigger", nil)
    cdc.RegisterConcrete(MsgInsertRow{}, "dbchain/InsertRow", nil)
    cdc.RegisterConcrete(MsgUpdateRow{}, "dbchain/UpdateRow", nil)
    cdc.RegisterConcrete(MsgDeleteRow{}, "dbchain/DeleteRow", nil)
    cdc.RegisterConcrete(MsgFreezeRow{}, "dbchain/FreezeRow", nil)
    cdc.RegisterConcrete(MsgModifyGroup{}, "dbchain/ModifyGroup", nil)
    cdc.RegisterConcrete(MsgSetGroupMemo{}, "dbchain/SetGroupMemo", nil)
    cdc.RegisterConcrete(MsgModifyGroupMember{}, "dbchain/ModifyGroupMember", nil)
    cdc.RegisterConcrete(MsgAddFriend{}, "dbchain/AddFriend", nil)
    cdc.RegisterConcrete(MsgRespondFriend{}, "dbchain/RespondFriend", nil)
}

