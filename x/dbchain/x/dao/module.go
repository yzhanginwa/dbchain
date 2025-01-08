package dao

import (
    "github.com/dbchaincloud/cosmos-sdk/codec"
    "github.com/dbchaincloud/cosmos-sdk/types/module"
    "github.com/yzhanginwa/dbchain/x/dao/internal/keeper"
    "github.com/yzhanginwa/dbchain/x/dao/internal/types"
)

type AppModuleBasic struct{}

func (AppModuleBasic) Name() string {
    return types.ModuleName
}

func (AppModuleBasic) RegisterCodec(cdc *codec.Codec) {
    // Register module codec
}

func (AppModuleBasic) DefaultGenesis() json.RawMessage {
    return nil
}

func (AppModuleBasic) ValidateGenesis(bz json.RawMessage) error {
    return nil
}

type AppModule struct {
    AppModuleBasic
    keeper keeper.Keeper
}

func NewAppModule(k keeper.Keeper) AppModule {
    return AppModule{
        AppModuleBasic: AppModuleBasic{},
        keeper:         k,
    }
}

func (AppModule) Name() string {
    return types.ModuleName
}

func (am AppModule) RegisterInvariants(ir sdk.InvariantRegistry) {}

func (am AppModule) Route() string {
    return types.RouterKey
}

func (am AppModule) NewHandler() sdk.Handler {
    return nil
}

func (am AppModule) QuerierRoute() string {
    return types.ModuleName
}

func (am AppModule) NewQuerierHandler() sdk.Querier {
    return nil
}

func (am AppModule) InitGenesis(ctx sdk.Context, data json.RawMessage) []abci.ValidatorUpdate {
    return []abci.ValidatorUpdate{}
}

func (am AppModule) ExportGenesis(ctx sdk.Context) json.RawMessage {
    return nil
}
