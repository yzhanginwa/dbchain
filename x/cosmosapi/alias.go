package cosmosapi

import (
    "github.com/yzhanginwa/cosmos-api/x/cosmosapi/internal/keeper"
    "github.com/yzhanginwa/cosmos-api/x/cosmosapi/internal/types"
)

const (
    ModuleName = types.ModuleName
    RouterKey  = types.RouterKey
    StoreKey   = types.StoreKey
)

var (
    NewKeeper        = keeper.NewKeeper
    NewQuerier       = keeper.NewQuerier
    ModuleCdc        = types.ModuleCdc
    RegisterCodec    = types.RegisterCodec
)

type (
    Keeper          = keeper.Keeper
    MsgCreateTable  = types.MsgCreateTable
    MsgDropTable    = types.MsgDropTable
    MsgAddColumn    = types.MsgAddColumn
    MsgDropColumn   = types.MsgDropColumn
    MsgRenameColumn = types.MsgRenameColumn
    MsgModifyOption = types.MsgModifyOption
    MsgModifyFieldOption = types.MsgModifyFieldOption
    MsgCreateIndex  = types.MsgCreateIndex
    MsgInsertRow    = types.MsgInsertRow
    MsgUpdateRow    = types.MsgUpdateRow
    MsgDeleteRow    = types.MsgDeleteRow

    MsgAddAdminAccount = types.MsgAddAdminAccount

    Table           = types.Table

    GenesisState    = types.GenesisState
)
