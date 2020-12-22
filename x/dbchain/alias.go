package dbchain

import (
    "github.com/yzhanginwa/dbchain/x/dbchain/internal/keeper"
    "github.com/yzhanginwa/dbchain/x/dbchain/internal/types"
)

const (
    ModuleName = types.ModuleName
    RouterKey  = types.RouterKey
    StoreKey   = types.StoreKey
    OracleHome = types.OracleHome
    CLIHome    = types.CLIHome
    NodeHome   = types.NodeHome
)

var (
    NewKeeper        = keeper.NewKeeper
    NewQuerier       = keeper.NewQuerier
    ModuleCdc        = types.ModuleCdc
    RegisterCodec    = types.RegisterCodec
)

type (
    Keeper          = keeper.Keeper
    MsgCreateApplication = types.MsgCreateApplication
    MsgCreateSysDatabase = types.MsgCreateSysDatabase
    MsgSetSchemaStatus   = types.MsgSetSchemaStatus
    MsgSetDatabasePermission = types.MsgSetDatabasePermission
    MsgModifyDatabaseUser    = types.MsgModifyDatabaseUser
    MsgAddFunction  = types.MsgAddFunction
    MsgCallFunction = types.MsgCallFunction
    MsgAddCustomQuerier = types.MsgAddCustomQuerier
    MsgCreateTable  = types.MsgCreateTable
    MsgDropTable    = types.MsgDropTable
    MsgAddColumn    = types.MsgAddColumn
    MsgDropColumn   = types.MsgDropColumn
    MsgRenameColumn = types.MsgRenameColumn
    MsgModifyOption = types.MsgModifyOption
    MsgAddInsertFilter    = types.MsgAddInsertFilter
    MsgDropInsertFilter   = types.MsgDropInsertFilter
    MsgAddTrigger   = types.MsgAddTrigger
    MsgDropTrigger  = types.MsgDropTrigger
    MsgSetTableMemo = types.MsgSetTableMemo
    MsgModifyColumnOption = types.MsgModifyColumnOption
    MsgSetColumnMemo      = types.MsgSetColumnMemo
    MsgCreateIndex  = types.MsgCreateIndex
    MsgDropIndex    = types.MsgDropIndex
    MsgInsertRow    = types.MsgInsertRow
    MsgUpdateRow    = types.MsgUpdateRow
    MsgDeleteRow    = types.MsgDeleteRow
    MsgFreezeRow    = types.MsgFreezeRow

    MsgAddFriend    = types.MsgAddFriend
    MsgDropFriend   = types.MsgDropFriend
    MsgRespondFriend   = types.MsgRespondFriend
    MsgModifyGroup     = types.MsgModifyGroup
    MsgSetGroupMemo    = types.MsgSetGroupMemo
    MsgModifyGroupMember  = types.MsgModifyGroupMember

    Table           = types.Table

    GenesisState    = types.GenesisState
)
