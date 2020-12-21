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
    return nil
}
