package app

import (
    "encoding/json"
    "errors"
    "fmt"
    sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
    "os"
    "reflect"
    "unsafe"

    abci "github.com/tendermint/tendermint/abci/types"
    "github.com/tendermint/tendermint/libs/log"
    tmos "github.com/tendermint/tendermint/libs/os"
    tmtypes "github.com/tendermint/tendermint/types"
    dbm "github.com/tendermint/tm-db"

    bam "github.com/cosmos/cosmos-sdk/baseapp"
    "github.com/cosmos/cosmos-sdk/codec"
    sdk "github.com/cosmos/cosmos-sdk/types"
    "github.com/cosmos/cosmos-sdk/types/module"
    "github.com/cosmos/cosmos-sdk/version"
    "github.com/cosmos/cosmos-sdk/x/auth"
    distr "github.com/cosmos/cosmos-sdk/x/distribution"
    "github.com/cosmos/cosmos-sdk/x/genutil"
    "github.com/cosmos/cosmos-sdk/x/params"
    "github.com/cosmos/cosmos-sdk/x/slashing"
    "github.com/cosmos/cosmos-sdk/x/staking"
    "github.com/cosmos/cosmos-sdk/x/supply"
    "github.com/yzhanginwa/dbchain/x/bank"

    "github.com/yzhanginwa/dbchain/x/dbchain"
)

const (
	appName     = "dbchain"
	OracleHome  = dbchain.OracleHome
	CLIHome     = dbchain.CLIHome
	NodeHome    = dbchain.NodeHome
	dailyHeight = 17280
	days        = 265 * 2
	ValidBlockHeight = dailyHeight * days
)

var (
    // default OracleHome directories for the application oracle
    DefaultOracleHome = os.ExpandEnv(OracleHome)

    // default home directories for the application CLI
    DefaultCLIHome = os.ExpandEnv(CLIHome)

    // DefaultNodeHome sets the folder where the applcation data and configuration will be stored
    DefaultNodeHome = os.ExpandEnv(NodeHome)

    // NewBasicManager is in charge of setting up basic module elemnets
    ModuleBasics = module.NewBasicManager(
        genutil.AppModuleBasic{},
        auth.AppModuleBasic{},
        bank.AppModuleBasic{},
        staking.AppModuleBasic{},
        distr.AppModuleBasic{},
        params.AppModuleBasic{},
        slashing.AppModuleBasic{},
        supply.AppModuleBasic{},

        dbchain.AppModule{},
    )
    // account permissions
    maccPerms = map[string][]string{
        auth.FeeCollectorName:     nil,
        distr.ModuleName:          nil,
        staking.BondedPoolName:    {supply.Burner, supply.Staking},
        staking.NotBondedPoolName: {supply.Burner, supply.Staking},
    }
)

// MakeCodec generates the necessary codecs for Amino
func MakeCodec() *codec.Codec {
    var cdc = codec.New()
    ModuleBasics.RegisterCodec(cdc)
    sdk.RegisterCodec(cdc)
    codec.RegisterCrypto(cdc)
    return cdc
}

type dbChainApp struct {
    *bam.BaseApp
    cdc *codec.Codec

    // keys to access the substores
    keys  map[string]*sdk.KVStoreKey
    tkeys map[string]*sdk.TransientStoreKey

    // Keepers
    accountKeeper  auth.AccountKeeper
    bankKeeper     bank.Keeper
    stakingKeeper  staking.Keeper
    slashingKeeper slashing.Keeper
    distrKeeper    distr.Keeper
    supplyKeeper   supply.Keeper
    paramsKeeper   params.Keeper
    dbChainKeeper dbchain.Keeper

    // Module Manager
    mm *module.Manager
}

// NewNameServiceApp is a constructor function for nameServiceApp
func NewDbChainApp(
    logger log.Logger, db dbm.DB, baseAppOptions ...func(*bam.BaseApp),
) *dbChainApp {

    // First define the top level codec that will be shared by the different modules
    cdc := MakeCodec()

    // BaseApp handles interactions with Tendermint through the ABCI protocol
    bApp := bam.NewBaseApp(appName, logger, db, auth.DefaultTxDecoder(cdc), baseAppOptions...)

    bApp.SetAppVersion(version.Version)

    keys := sdk.NewKVStoreKeys(bam.MainStoreKey, auth.StoreKey, staking.StoreKey,
        supply.StoreKey, distr.StoreKey, slashing.StoreKey, params.StoreKey, dbchain.StoreKey)

    tkeys := sdk.NewTransientStoreKeys(staking.TStoreKey, params.TStoreKey)

    // Here you initialize your application with the store keys it requires
    var app = &dbChainApp{
        BaseApp: bApp,
        cdc:     cdc,
        keys:    keys,
        tkeys:   tkeys,
    }

    // The ParamsKeeper handles parameter storage for the application
    app.paramsKeeper = params.NewKeeper(app.cdc, keys[params.StoreKey], tkeys[params.TStoreKey])
    // Set specific supspaces
    authSubspace := app.paramsKeeper.Subspace(auth.DefaultParamspace)
    bankSupspace := app.paramsKeeper.Subspace(bank.DefaultParamspace)
    stakingSubspace := app.paramsKeeper.Subspace(staking.DefaultParamspace)
    distrSubspace := app.paramsKeeper.Subspace(distr.DefaultParamspace)
    slashingSubspace := app.paramsKeeper.Subspace(slashing.DefaultParamspace)

    // The AccountKeeper handles address -> account lookups
    app.accountKeeper = auth.NewAccountKeeper(
        app.cdc,
        keys[auth.StoreKey],
        authSubspace,
        auth.ProtoBaseAccount,
    )

    // The BankKeeper allows you perform sdk.Coins interactions
    app.bankKeeper = bank.NewBaseKeeper(
        app.accountKeeper,
        bankSupspace,
        app.ModuleAccountAddrs(),
    )

    // The SupplyKeeper collects transaction fees and renders them to the fee distribution module
    app.supplyKeeper = supply.NewKeeper(
        app.cdc,
        keys[supply.StoreKey],
        app.accountKeeper,
        app.bankKeeper,
        maccPerms,
    )

    // The staking keeper
    stakingKeeper := staking.NewKeeper(
        app.cdc,
        keys[staking.StoreKey],
        app.supplyKeeper,
        stakingSubspace,
    )

    app.distrKeeper = distr.NewKeeper(
        app.cdc,
        keys[distr.StoreKey],
        distrSubspace,
        &stakingKeeper,
        app.supplyKeeper,
        auth.FeeCollectorName,
        app.ModuleAccountAddrs(),
    )

    app.slashingKeeper = slashing.NewKeeper(
        app.cdc,
        keys[slashing.StoreKey],
        &stakingKeeper,
        slashingSubspace,
    )

    // register the staking hooks
    // NOTE: stakingKeeper above is passed by reference, so that it will contain these hooks
    app.stakingKeeper = *stakingKeeper.SetHooks(
        staking.NewMultiStakingHooks(
            app.distrKeeper.Hooks(),
            app.slashingKeeper.Hooks()),
    )

    // The DbChainKeeper is the Keeper from module dbchain
    // It handles interactions with the namestore
    app.dbChainKeeper = dbchain.NewKeeper(
        app.bankKeeper,
        app.accountKeeper,
        keys[dbchain.StoreKey],
        app.cdc,
    )

    app.mm = module.NewManager(
        genutil.NewAppModule(app.accountKeeper, app.stakingKeeper, app.BaseApp.DeliverTx),
        auth.NewAppModule(app.accountKeeper),
        bank.NewAppModule(app.bankKeeper, app.accountKeeper),
        dbchain.NewAppModule(app.dbChainKeeper, app.bankKeeper),
        supply.NewAppModule(app.supplyKeeper, app.accountKeeper),
        distr.NewAppModule(app.distrKeeper, app.accountKeeper, app.supplyKeeper, app.stakingKeeper),
        slashing.NewAppModule(app.slashingKeeper, app.accountKeeper, app.stakingKeeper),
        staking.NewAppModule(app.stakingKeeper, app.accountKeeper, app.supplyKeeper),
    )

    app.mm.SetOrderBeginBlockers(distr.ModuleName, slashing.ModuleName, dbchain.ModuleName)
    app.mm.SetOrderEndBlockers(staking.ModuleName)

    // Sets the order of Genesis - Order matters, genutil is to always come last
    // NOTE: The genutils moodule must occur after staking so that pools are
    // properly initialized with tokens from genesis accounts.
    app.mm.SetOrderInitGenesis(
        distr.ModuleName,
        staking.ModuleName,
        auth.ModuleName,
        bank.ModuleName,
        slashing.ModuleName,
        dbchain.ModuleName,
        supply.ModuleName,
        genutil.ModuleName,
    )

    // register all module routes and module queriers
    app.mm.RegisterRoutes(app.Router(), app.QueryRouter())

    // The initChainer handles translating the genesis.json file into initial state for the network
    app.SetInitChainer(app.InitChainer)
    app.SetBeginBlocker(app.BeginBlocker)
    app.SetEndBlocker(app.EndBlocker)

    // The AnteHandler handles signature verification and transaction pre-processing
    app.SetAnteHandler(
        auth.NewAnteHandler(
            app.accountKeeper,
            app.supplyKeeper,
            auth.DefaultSigVerificationGasConsumer,
        ),
    )

    // initialize stores
    app.MountKVStores(keys)
    app.MountTransientStores(tkeys)

    err := app.LoadLatestVersion(app.keys[bam.MainStoreKey])
    if err != nil {
        tmos.Exit(err.Error())
    }

    return app
}

// GenesisState represents chain state at the start of the chain. Any initial state (account balances) are stored here.
type GenesisState map[string]json.RawMessage

func NewDefaultGenesisState() GenesisState {
    return ModuleBasics.DefaultGenesis()
}

func (app *dbChainApp) InitChainer(ctx sdk.Context, req abci.RequestInitChain) abci.ResponseInitChain {
    var genesisState GenesisState

    err := app.cdc.UnmarshalJSON(req.AppStateBytes, &genesisState)
    if err != nil {
        panic(err)
    }

    return app.mm.InitGenesis(ctx, genesisState)
}

func (app *dbChainApp) BeginBlocker(ctx sdk.Context, req abci.RequestBeginBlock) abci.ResponseBeginBlock {
    return app.mm.BeginBlock(ctx, req)
}
func (app *dbChainApp) EndBlocker(ctx sdk.Context, req abci.RequestEndBlock) abci.ResponseEndBlock {
    return app.mm.EndBlock(ctx, req)
}
func (app *dbChainApp) LoadHeight(height int64) error {
    return app.LoadVersion(height, app.keys[bam.MainStoreKey])
}

// ModuleAccountAddrs returns all the app's module account addresses.
func (app *dbChainApp) ModuleAccountAddrs() map[string]bool {
    modAccAddrs := make(map[string]bool)
    for acc := range maccPerms {
        modAccAddrs[supply.NewModuleAddress(acc).String()] = true
    }

    return modAccAddrs
}

//_________________________________________________________

func (app *dbChainApp) ExportAppStateAndValidators(forZeroHeight bool, jailWhiteList []string,
) (appState json.RawMessage, validators []tmtypes.GenesisValidator, err error) {

    // as if they could withdraw from the start of the next block
    ctx := app.NewContext(true, abci.Header{Height: app.LastBlockHeight()})

    genState := app.mm.ExportGenesis(ctx)
    appState, err = codec.MarshalJSONIndent(app.cdc, genState)
    if err != nil {
        return nil, nil, err
    }

    validators = staking.WriteValidators(ctx, app.stakingKeeper)

    return appState, validators, nil
}

func (app *dbChainApp) CheckTx(req abci.RequestCheckTx) abci.ResponseCheckTx {
    resp :=  app.BaseApp.CheckTx(req)
    if ValidBlockHeight <= 0 {
        return  resp
    }
    var ctx sdk.Context
    baseApp := reflect.ValueOf(app.BaseApp).Elem()
    for i := 0; i < baseApp.NumField(); i++ {
        if baseApp.Type().Field(i).Name == "checkState" {
            p := baseApp.Field(i).Pointer()
            np := (*state)(unsafe.Pointer(p))
            ctx = np.Context()
            break
        }
    }

    if ctx.IsZero() {
        return resp
    }

    current := ctx.BlockHeader().Height
    if current > ValidBlockHeight {
        err := errors.New(fmt.Sprintf("current block height is %d", current))
        ctx.Logger().Debug(err.Error())
        return sdkerrors.ResponseCheckTx(err, uint64(resp.GasWanted), uint64(resp.GasUsed), false)
    }
    return resp
}


///////////////////////////
//                       //
//     help struct       //
//                       //
///////////////////////////

type state struct {
    ms  sdk.CacheMultiStore
    ctx sdk.Context
}

func (st *state) CacheMultiStore() sdk.CacheMultiStore {
    return st.ms.CacheMultiStore()
}

func (st *state) Context() sdk.Context {
    return st.ctx
}
