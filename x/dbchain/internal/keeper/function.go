package keeper

import (
    "encoding/json"
    "errors"
    "github.com/cosmos/cosmos-sdk/codec"
    sdk "github.com/cosmos/cosmos-sdk/types"
    "github.com/mr-tron/base58"
    lua "github.com/yuin/gopher-lua"
    "github.com/yuin/gopher-lua/parse"
    "github.com/yzhanginwa/dbchain/x/dbchain/internal/super_script"
    "github.com/yzhanginwa/dbchain/x/dbchain/internal/types"
    "strings"
)

// parameter is a json encoded array of string
// custom querier is also a function, but need a different store key,
// parameter t is used to distinguish type,when t == 0 ,it means function. when t == 1 ,it means querier
func (k Keeper) AddFunction(ctx sdk.Context, appId uint, functionName, description, body string, owner sdk.AccAddress, t int) error {

    if err := checkLuaSyntax(body); err != nil {
        return err
    }
    store := DbChainStore(ctx, k.storeKey)
    function := types.NewFunction()
    function.Name = getFuncNameFromBody(ctx, k, appId, owner, body)
    function.Description = description
    function.Body = body
    function.Owner = owner

    functionStoreKey := ""
    if t == FuncHandleType {
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
    if t == FuncHandleType {
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
    defer voidLuaHandle(appId, t)
    return  store.Set([]byte(functionsStoreKey),k.cdc.MustMarshalBinaryBare(functions))
}

func getFuncNameFromBody(ctx sdk.Context, keeper Keeper, appId uint, owner sdk.AccAddress, body string)string {
    L := lua.NewState()
    defer L.Close()
    p := super_script.NewPreprocessor(strings.NewReader(body))
    p.Process()
    newScript := p.Reconstruct()

    fn,_ := L.LoadString(newScript)
    funcName := fn.Proto.Constants[0].String()
    return funcName
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
    return callLuaScriptFunc(ctx, appId, owner, k, FunctionName, params, FuncHandleType)
}


func (k Keeper) DropFunction(ctx sdk.Context, appId uint, owner sdk.AccAddress, FunctionName string, t int) error {
    store := DbChainStore(ctx, k.storeKey)
    functionStoreKey := ""
    if t == FuncHandleType {
        functionStoreKey = getFunctionKey(appId, FunctionName)
    } else {
        functionStoreKey = getQuerierKey(appId, FunctionName)
    }
    err := store.Delete([]byte(functionStoreKey))
    if err != nil {
        return err
    }
    //update Functions set
    FunctionsKey := ""
    if t == FuncHandleType {
        FunctionsKey = getFunctionsKey(appId)
    } else {
        FunctionsKey = getQueriersKey(appId)
    }

    bz, err := store.Get([]byte(FunctionsKey))
    if err != nil{
        return err
    }
    if bz == nil {
        return nil
    }
    var names []string
    k.cdc.MustUnmarshalBinaryBare(bz, &names)
    for index , val := range names {
        if FunctionName == val {
            names = append(names[:index],names[index+1: ]... )
            break
        }
    }
    bz = k.cdc.MustMarshalBinaryBare(names)
    err = store.Set([]byte(FunctionsKey), bz)
    if err != nil {
        return err
    }
    return nil
}

// custom querier is also a function, but need a different store key,
// parameter t is used to distinguish type,when t == 0 ,it means function. when t == 1 ,it means querier
func (k Keeper) GetFunctions(ctx sdk.Context, appId uint, t int) []string {
    store := DbChainStore(ctx, k.storeKey)
    FunctionsKey := ""
    if t == FuncHandleType {
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
    if t == FuncHandleType {
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
    bArgument, err := base58.Decode(argument)
    if err != nil {
        return nil,errors.New("call DoCustomQuerier err :" + err.Error())
    }
    arguments := strings.Split( string(bArgument), "/")
    params := make([]lua.LValue,0)
    for _,v := range arguments{
        param, err := base58.Decode(v)
        if err != nil {
            return nil,errors.New("call DoCustomQuerier err :" + err.Error())
        }
        params = append(params, lua.LString(string(param)))
    }

    return callLuaScriptQuerierFunc(ctx, appId, addr, k, querierInfo.Name, params, QueryHandleType)
}

func (k Keeper) DoDynamicScript(ctx sdk.Context, appId uint, script string, addr sdk.AccAddress) ([]byte, error){
    var l = lua.NewState(lua.Options{
        SkipOpenLibs: true, //set SkipOpenLibs true to prevent lua open libs,because this libs can call os function and operate files
        RegistrySize: 256, //Save function return value， default value is 256*20, it is too large
    })
    openBase(l)
    registerTableType(l, ctx , appId, k, addr)
    //point : get go function
    goExportFunc := getGoExportQueryFunc(ctx, appId, k, addr)
    goExportFuncNew := getGoExportQueryFuncNew(ctx, appId, k, addr)
    goExportToolFunc := getGoExportToolFunc(ctx)
    //register go function
    for name, fn := range goExportFunc{
        l.SetGlobal(name, l.NewFunction(fn))
    }
    for name, fn := range goExportFuncNew{
        l.SetGlobal(name, l.NewFunction(fn))
    }
    for name, fn := range goExportToolFunc{
        l.SetGlobal(name, l.NewFunction(fn))
    }
    newScript := restructureLuaScript(script)
    if err := l.DoString(newScript); err != nil{
       return nil, err
    }

    //handle return
    defer l.Pop(l.GetTop())
    lRes := l.Get(1)
    if tableRes ,ok := lRes.(*lua.LTable); ok {
        res := convertLuaTableToGo(tableRes)
        bz, err := codec.MarshalJSONIndent(k.cdc, res)
        if err != nil {
            return nil, errors.New("could not marshal result to JSON")
        }
        return bz, nil
    }

    switch lRes.Type() {
    case lua.LTNil:
        return []byte("[]"), nil
    case lua.LTString:
        errString := lRes.(lua.LString)
        return nil, errors.New(errString.String())
    default:
        return []byte("[]") , nil
    }

}

func checkLuaSyntax(script string) error {
    p := super_script.NewPreprocessor(strings.NewReader(script))
    p.Process()
    if !p.Success {
    	return  errors.New("Script syntax error")
    }
    newScript := p.Reconstruct()
    err := compileAndCheckLuaScript(newScript)
    return err
}

func compileAndCheckLuaScript(script string) error {
    name := "<string>"
    chunk, err := parse.Parse(strings.NewReader(script), name)
    if err != nil {
        return err
    }
    _, err = lua.Compile(chunk, name)
    if err != nil {
        return err
    }
    return nil
}