package keeper

import (
    "errors"
    "encoding/json"
    sdk "github.com/cosmos/cosmos-sdk/types"
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