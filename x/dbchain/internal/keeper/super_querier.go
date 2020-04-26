package keeper

import (
    "github.com/mr-tron/base58"
    "encoding/json"
    "github.com/cosmos/cosmos-sdk/codec"
    sdk "github.com/cosmos/cosmos-sdk/types"
    sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
    abci "github.com/tendermint/tendermint/abci/types"
    "github.com/yzhanginwa/dbchain/x/dbchain/internal/utils"
)

type Conditional struct {
    Left string
    Operator string
    Right string
}

type Logical struct {
    Left Conditional
    Operator string
    Right Conditional
}

type QuerierBuilder struct {
    Table string
    Select []string
    Where []Conditional
    Order []string
    Limit uint
}



/////////////////////
//                 //
// dbChain Querier //
//                 //
/////////////////////

func queryQuerier(ctx sdk.Context, path []string, req abci.RequestQuery, keeper Keeper) ([]byte, error) {
    accessCode:= path[0]
    addr, err := utils.VerifyAccessCode(accessCode)
    if err != nil {
        return []byte{}, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "Access code is not valid!")
    }

    appId, err := keeper.GetDatabaseId(ctx, path[1])
    if err != nil {
        return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "Invalid app code")
    }

    querierObjJson, err := base58.Decode(path[2])
    if err != nil {
        return []byte{}, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "Querier object json string base58 encoding error!")
    }

    var querierObjs [](map[string]string)

    if err := json.Unmarshal(querierObjJson, &querierObjs); err != nil {
        return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest,"Failed to parse querier objects!")
    }
   
    result, err := querierSuperHandler(ctx, keeper, appId, querierObjs, addr)
    if err != nil {
        return nil, err
    }
        
    res, err := codec.MarshalJSONIndent(keeper.cdc, result)
    if err != nil {
        panic("could not marshal result to JSON")
    }

    return res, nil
}

func querierSuperHandler(ctx sdk.Context, keeper Keeper, appId uint, querierObjs [](map[string]string), owner sdk.AccAddress) ([](map[string]string), error) {
    builder := QuerierBuilder{}
    
    for i := 0; i < len(querierObjs); i++ {
        qo := querierObjs[i]
        switch qo["method"] {
        case "table":
            if builder.Table != "" {
                return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "Querier can only have one table command!")
            }
            builder.Table = qo["table"]
        case "select" :
            fields := []string{}
            if err := json.Unmarshal([]byte(qo["fields"]), &fields); err != nil {
                return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest,"Failed to parse fields!")
            }
            builder.Select = fields
        }
    }

    if len(builder.Select) == 0 {
        table, err := keeper.GetTable(ctx, appId, builder.Table)    
        if err != nil {
            return nil, err
        }
        builder.Select = table.Fields
    }

    store := ctx.KVStore(keeper.storeKey)
    ids := keeper.FindAll(ctx, appId, builder.Table, owner) 
    var result = [](map[string]string){}
    for _, id := range ids {
        record := map[string]string{}
        for _, f := range builder.Select {
            key := getDataKeyBytes(appId, builder.Table, f, id)
            bz := store.Get(key)
            var value string
            if bz != nil {
                keeper.cdc.MustUnmarshalBinaryBare(bz, &value)
                record[f] = value
            }
        }
        result = append(result, record)
    }
    return result, nil
}

//////////////////
//              //
// helper funcs //
//              //
//////////////////

