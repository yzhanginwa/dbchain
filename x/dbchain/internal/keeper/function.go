package keeper

import (
    "encoding/json"
    "errors"
    sdk "github.com/cosmos/cosmos-sdk/types"
    lua "github.com/yuin/gopher-lua"
    "github.com/yzhanginwa/dbchain/x/dbchain/internal/types"
)

// parameter is a json encoded array of string
func (k Keeper) AddFunction(ctx sdk.Context, appId uint, functionName, parameter, body string, owner sdk.AccAddress) error {
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

    err := store.Set([]byte(getFunctionKey(appId, function.Name)), k.cdc.MustMarshalBinaryBare(function))
    if err != nil{
        return err
    }
    //store functions
    var functions []string
    bz ,err := store.Get([]byte(getFunctionsKey(appId)))
    if err != nil{
        return err
    }
    if bz == nil{
        functions = append(functions,function.Name)
    }else{
        k.cdc.MustUnmarshalBinaryBare(bz,&functions)
        functions = append(functions,function.Name)
    }
    return  store.Set([]byte(getFunctionsKey(appId)),k.cdc.MustMarshalBinaryBare(functions))
}

func (k Keeper) CallFunction(ctx sdk.Context, appId uint, owner sdk.AccAddress, FunctionName, Argument string)error{
    functionInfo := k.GetFunctionInfo(ctx, appId, FunctionName)
    var arguments = make([]string,0)
    if err := json.Unmarshal([]byte(Argument), &arguments); err != nil {
        return errors.New("argument should be json encoded array!")
    }
    //get lua script and params
    body := functionInfo.Body
    params := make([]lua.LValue,0)
    for _,v := range arguments{
        params = append(params, lua.LString(v))
    }
    //point : get go function
    goExportFunc := getGoExportFunc(ctx, appId, k, owner)
    L := lua.NewState(lua.Options{
        SkipOpenLibs : true, //set SkipOpenLibs true to prevent lua open libs,because this libs can call os function and operate files
    })
    defer L.Close()
    //register go function
    for name, fn := range goExportFunc{
        L.SetGlobal(name, L.NewFunction(fn))
    }
    //compile lua script
    if err := L.DoString(body); err != nil{
        return err
    }
    //call lua script
    if err := L.CallByParam(lua.P{
        Fn:      L.GetGlobal(FunctionName),
        NRet:    0,       //脚本返回参数个数，暂时不用返回参数
        Protect: true,    //这里设置为ture表示当执行脚本出现panic时，以error返回
    }, params...); err != nil{
        return err
    }
    return nil
}

func (k Keeper) GetFunctions(ctx sdk.Context, appId uint) []string {
    store := DbChainStore(ctx, k.storeKey)
    FunctionsKey := getFunctionsKey(appId)
    bz, err := store.Get([]byte(FunctionsKey))
    if bz == nil || err != nil{
        return []string{}
    }
    var functionNames []string
    k.cdc.MustUnmarshalBinaryBare(bz, &functionNames)
    return functionNames
}

func (k Keeper) GetFunctionInfo(ctx sdk.Context, appId uint, functionName string) types.Function {
    store := DbChainStore(ctx, k.storeKey)
    FunctionKey := getFunctionKey(appId, functionName)
    bz, err := store.Get([]byte(FunctionKey))
    if bz == nil || err != nil{
        return types.Function{}
    }
    var functionInfo types.Function
    k.cdc.MustUnmarshalBinaryBare(bz, &functionInfo)
    return functionInfo
}