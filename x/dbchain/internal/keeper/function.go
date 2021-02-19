package keeper

import (
    "encoding/json"
    "errors"
    "github.com/cosmos/cosmos-sdk/codec"
    sdk "github.com/cosmos/cosmos-sdk/types"
    lua "github.com/yuin/gopher-lua"
    "github.com/yzhanginwa/dbchain/x/dbchain/internal/types"
)

// parameter is a json encoded array of string
// custom querier is also a function, but need a different store key,
// parameter t is used to distinguish type,when t == 0 ,it means function. when t == 1 ,it means querier
func (k Keeper) AddFunction(ctx sdk.Context, appId uint, functionName, parameter, body string, owner sdk.AccAddress, t int) error {
    var params = []string{}
    if err := json.Unmarshal([]byte(parameter), &params); err != nil {
        return errors.New("Parameter should be json encoded array!")
    }

    store := DbChainStore(ctx, k.storeKey)
    function := types.NewFunction()
    function.Name = functionName
    function.Parameter = params
    function.Body = body
    function.Owner = owner

    functionStoreKey := ""
    if t == 0 {
        functionStoreKey = getFunctionKey(appId, function.Name)
    } else {
        functionStoreKey = getQuerierKey(appId, function.Name)
    }

    err := store.Set([]byte(functionStoreKey), k.cdc.MustMarshalBinaryBare(function))
    if err != nil{
        return err
    }
    //store functions
    var functions []string
    functionsStoreKey := ""
    if t == 0 {
        functionsStoreKey = getFunctionsKey(appId)
    } else {
        functionsStoreKey = getQueriersKey(appId)
    }
    bz ,err := store.Get([]byte(functionsStoreKey))
    if err != nil{
        return err
    }
    if bz == nil{
        functions = append(functions,function.Name)
    }else{
        k.cdc.MustUnmarshalBinaryBare(bz,&functions)
        i := 0
        for ; i < len(functions);i++{
            if functions[i] == function.Name{
                break
            }
        }
        if i >= len(functions) {
            functions = append(functions,function.Name)
        }

    }
    //TODO is it necessary to distinguish querier from func
    defer voidLuaHandle(appId)
    return  store.Set([]byte(functionsStoreKey),k.cdc.MustMarshalBinaryBare(functions))
}

func (k Keeper) CallFunction(ctx sdk.Context, appId uint, owner sdk.AccAddress, FunctionName, Argument string) error {
    var arguments = make([]string,0)
    if err := json.Unmarshal([]byte(Argument), &arguments); err != nil {
        return errors.New("argument should be json encoded array!")
    }
    //get lua script and params
    params := make([]lua.LValue,0)
    for _,v := range arguments{
        params = append(params, lua.LString(v))
    }
    return callLuaScriptFunc(ctx, appId, owner, k, FunctionName, params)
}

// custom querier is also a function, but need a different store key,
// parameter t is used to distinguish type,when t == 0 ,it means function. when t == 1 ,it means querier
func (k Keeper) GetFunctions(ctx sdk.Context, appId uint, t int) []string {
    store := DbChainStore(ctx, k.storeKey)
    FunctionsKey := ""
    if t == 0 {
        FunctionsKey = getFunctionsKey(appId)
    } else {
        FunctionsKey = getQueriersKey(appId)
    }

    bz, err := store.Get([]byte(FunctionsKey))
    if bz == nil || err != nil{
        return []string{}
    }
    var names []string
    k.cdc.MustUnmarshalBinaryBare(bz, &names)
    return names
}

func (k Keeper) GetFunctionInfo(ctx sdk.Context, appId uint, name string, t int) types.Function {
    store := DbChainStore(ctx, k.storeKey)
    key := ""
    if t == 0 {
        key = getFunctionKey(appId, name)
    } else {
        key = getQuerierKey(appId, name)
    }

    bz, err := store.Get([]byte(key))
    if bz == nil || err != nil{
        return types.Function{}
    }
    var info types.Function
    k.cdc.MustUnmarshalBinaryBare(bz, &info)
    return info
}

func (k Keeper) DoCustomQuerier(ctx sdk.Context, appId uint, querierInfo types.Function, argument string, addr sdk.AccAddress) ([]byte, error){
    //get lua script and params
    var arguments = make([]string,0)
    if err := json.Unmarshal([]byte(argument), &arguments); err != nil {
        return nil,errors.New("argument should be json encoded array!")
    }

    body := querierInfo.Body
    params := make([]lua.LValue,0)
    for _,v := range arguments{
        params = append(params, lua.LString(v))
    }
    //point : get go function
    goExportFunc := getGoExportQueryFunc(ctx, appId, k, addr)
    L := lua.NewState(lua.Options{
        SkipOpenLibs : true, //set SkipOpenLibs true to prevent lua open libs,because this libs can call os function and operate files
    })
    defer L.Close()
    //register go function
    for name, fn := range goExportFunc {
        L.SetGlobal(name, L.NewFunction(fn))
    }
    //compile lua script
    if err := L.DoString(body); err != nil{
        return nil, err
    }
    //call lua script
    if err := L.CallByParam(lua.P{
        Fn:      L.GetGlobal(querierInfo.Name),
        NRet:    1,       //脚本返回参数个数
        Protect: true,    //这里设置为ture表示当执行脚本出现panic时，以error返回
    }, params...); err != nil{
        return nil, err
    }
    //handle return
    lRes := L.Get(1)
    if tableRes ,ok := lRes.(*lua.LTable); ok {
        res := convertLuaTableToGo(tableRes)
        bz, err := codec.MarshalJSONIndent(k.cdc, res)
        if err != nil {
            return nil, errors.New("could not marshal result to JSON")
        }
        return bz, nil
    }

    return nil, errors.New("lua return err")
}