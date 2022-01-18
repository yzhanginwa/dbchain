package app

import (
    "crypto/sha256"
    "encoding/hex"
    "encoding/json"
    "fmt"
    "github.com/dbchaincloud/cosmos-sdk/server"
    "github.com/dbchaincloud/cosmos-sdk/x/auth/types"
    "github.com/spf13/viper"
    "github.com/yzhanginwa/dbchain/address"
    "os"
    "reflect"
    "strings"
    "unsafe"

    abci "github.com/dbchaincloud/tendermint/abci/types"
    "github.com/dbchaincloud/tendermint/libs/log"
    tmos "github.com/dbchaincloud/tendermint/libs/os"
    tmtypes "github.com/dbchaincloud/tendermint/types"
    dbm "github.com/tendermint/tm-db"
    "github.com/yzhanginwa/dbchain/x/bank"
    qch "github.com/yzhanginwa/dbchain/x/dbchain/query_cache_helper"

    bam "github.com/dbchaincloud/cosmos-sdk/baseapp"
    "github.com/dbchaincloud/cosmos-sdk/codec"
    sdk "github.com/dbchaincloud/cosmos-sdk/types"
    "github.com/dbchaincloud/cosmos-sdk/types/module"
    "github.com/dbchaincloud/cosmos-sdk/version"
    "github.com/dbchaincloud/cosmos-sdk/x/auth"
    distr "github.com/dbchaincloud/cosmos-sdk/x/distribution"
    "github.com/dbchaincloud/cosmos-sdk/x/genutil"
    "github.com/dbchaincloud/cosmos-sdk/x/params"
    "github.com/dbchaincloud/cosmos-sdk/x/slashing"
    "github.com/dbchaincloud/cosmos-sdk/x/staking"
    "github.com/dbchaincloud/cosmos-sdk/x/supply"

    "github.com/yzhanginwa/dbchain/x/dbchain"
)

const (
    appName     = "dbchain"
    OracleHome  = dbchain.OracleHome
    CLIHome     = dbchain.CLIHome
    NodeHome    = dbchain.NodeHome
    //define format of addr
    // PrefixValidator is the prefix for validator keys
    PrefixValidator = "val"
    // PrefixConsensus is the prefix for consensus keys
    PrefixConsensus = "cons"
    // PrefixPublic is the prefix for public keys
    PrefixPublic = "pub"
    // PrefixOperator is the prefix for operator keys
    PrefixOperator = "oper"

    // Bech32PrefixAccAddr defines the Bech32 prefix of an account's address
    Bech32PrefixAccAddr = address.Bech32MainPrefix
    // Bech32PrefixAccPub defines the Bech32 prefix of an account's public key
    Bech32PrefixAccPub = address.Bech32MainPrefix + PrefixPublic
    // Bech32PrefixValAddr defines the Bech32 prefix of a validator's operator address
    Bech32PrefixValAddr = address.Bech32MainPrefix + PrefixValidator + PrefixOperator
    // Bech32PrefixValPub defines the Bech32 prefix of a validator's operator public key
    Bech32PrefixValPub = address.Bech32MainPrefix + PrefixValidator + PrefixOperator + PrefixPublic
    // Bech32PrefixConsAddr defines the Bech32 prefix of a consensus node address
    Bech32PrefixConsAddr = address.Bech32MainPrefix + PrefixValidator + PrefixConsensus
    // Bech32PrefixConsPub defines the Bech32 prefix of a consensus node public key
    Bech32PrefixConsPub = address.Bech32MainPrefix + PrefixValidator + PrefixConsensus + PrefixPublic
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
    //set minGasPrices
    miniGasPrice := viper.GetString(server.FlagMinGasPrices)
    if miniGasPrice != "" {
        baseAppOptions = append(baseAppOptions, bam.SetMinGasPrices(miniGasPrice))
    }

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
        keys[dbchain.StoreKey],
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

func (app *dbChainApp) DeliverTx(req abci.RequestDeliverTx) abci.ResponseDeliverTx {

    resp :=  app.BaseApp.DeliverTx(req)

    //get gas price
    txDecoder := auth.DefaultTxDecoder(app.cdc)
    tx, err := txDecoder(req.Tx)
    if err != nil {
        return resp
    }
    stdTx , ok := tx.(auth.StdTx)
    if !ok {
        return resp
    }
    var gasPrices sdk.DecCoins
    if !stdTx.Fee.Amount.IsZero() {
        gasPrices = stdTx.Fee.GasPrices()
    }



    //calc Fees
    requiredFees := make(sdk.Coins, 0)
    glDecWanted := sdk.NewDec(resp.GasWanted)
    glDecUsed := sdk.NewDec(resp.GasUsed)
    for _, gp := range gasPrices {
        feeWanted := gp.Amount.Mul(glDecWanted)
        feeUsed := gp.Amount.Mul(glDecUsed)
        if feeUsed.GTE(feeWanted) {
            continue
        }
        coinWanted := sdk.NewCoin(gp.Denom, feeWanted.Ceil().RoundInt())
        coinUsed := sdk.NewCoin(gp.Denom, feeUsed.Ceil().RoundInt())
        if coinUsed.IsGTE(coinWanted) {
            continue
        }
        requiredFees = append(requiredFees, coinWanted.Sub(coinUsed)) //coinWanted.Sub(coinUsed)
    }
    defer func() {
        // pkg of reflect may be panic
        recover()
    }()
    //get ctx by pkg of reflect because deliverState is private
    baseApp := reflect.ValueOf(app.BaseApp).Elem()
    var ctx sdk.Context
    for i := 0; i < baseApp.NumField(); i++ {
        if baseApp.Type().Field(i).Name == "deliverState" {
            p := baseApp.Field(i).Pointer()
            np := (*state)(unsafe.Pointer(p))
            ctx = np.Context()
            break
        }
    }

    if ctx.IsZero() {
        return resp
    }

    //return extra fees
    if len(requiredFees) > 0 {
        feePayer := stdTx.FeePayer()
        err = app.supplyKeeper.SendCoinsFromModuleToAccount(ctx, types.FeeCollectorName,feePayer, requiredFees)
        if err != nil {
            fmt.Println(err)
        }
    }

    //set account txs hash
    var hash = sha256.Sum256(req.Tx)
    hexHash := hex.EncodeToString(hash[:])
    usedGas := app.SaveAddrTx(ctx, resp, stdTx, gasPrices, hexHash)
    resp.Info = usedGas
    return resp
}

func (app *dbChainApp) SaveAddrTx(ctx sdk.Context ,resp abci.ResponseDeliverTx, stdTx auth.StdTx, gasPrices sdk.DecCoins, txHash string)  string {
    //set account txs hash
    var addr sdk.AccAddress
    if len(stdTx.Msgs) > 0 {
        signers := stdTx.Msgs[0].GetSigners()
        if len(signers) > 0 {
            addr = signers[0]
        }
    }

    data := map[string]interface{} {
        "userAccountAddress" : addr.String(),
        "txHash" : txHash,
        "txTime" : ctx.BlockHeader().Time.Local().Format("2006-01-02 15:04:05"),
        "state" : checkTxStatus(resp.Log),
        "blockHeight" : ctx.BlockHeader().Height,
    }


    //calc usedFees
    usedFees := make([]string, len(gasPrices))//
    glDecWanted := sdk.NewDec(resp.GasWanted)
    glDecUsed := sdk.NewDec(resp.GasUsed)
    //usedGas := sdk.NewDec(resp.GasUsed)
    for i, gp := range gasPrices {
        feeWanted :=gp.Amount.Mul(glDecWanted)
        feeUsed := gp.Amount.Mul(glDecUsed)

        if feeUsed.GTE(feeWanted) {
            coin := sdk.NewCoin(gp.Denom, feeWanted.Ceil().RoundInt())
            usedFees[i] = coin.String()
        } else {
            coin := sdk.NewCoin(gp.Denom, feeUsed.Ceil().RoundInt())
            usedFees[i] = coin.String()
        }
    }
    usedFeesString := strings.Join(usedFees,",")

    data["gas"] = usedFeesString
    app.dbChainKeeper.SaveAddrTxs(ctx, addr, data)
    feeOfTxCost := "fee of tx cost : "
    if usedFeesString == "" {
        return feeOfTxCost + "null"
    }
    return feeOfTxCost + usedFeesString
}

func (app *dbChainApp) EndBlocker(ctx sdk.Context, req abci.RequestEndBlock) abci.ResponseEndBlock {
    ret := app.mm.EndBlock(ctx, req)
    qch.NotifyTableExpiration("", "")       // to notify querier cache to invalidate accumulated tables
    return ret
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

// 2 : success
// 3 : fail
func checkTxStatus(log string) int {
    data := make([]interface{}, 0)
    err := json.Unmarshal([]byte(log), &data)
    if err != nil {
        //tx fail
        // "insufficient funds: insufficient account funds; 18dbctoken \u003c 10000dbctoken: failed to execute message; message index: 0"
        return 3
    } else {
        //tx success
        //[{"msg_index":0,"log":"","events":[{"type":"message","attributes":[{"key":"action","value":"send"},{"key":"sender","value":"cosmos1n6yqmysvcz0cpnd52427ldjmjj493pk2uhpcu3"},{"key":"module","value":"bank"}]},{"type":"transfer","attributes":[{"key":"recipient","value":"cosmos156p5rmhpd3l709ygg7t80fu96fm4mrtsqhftvx"},{"key":"sender","value":"cosmos1n6yqmysvcz0cpnd52427ldjmjj493pk2uhpcu3"},{"key":"amount","value":"1dbctoken"}]}]}]
        return 2
    }
}
