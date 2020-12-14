package keeper

import (
    "fmt"
    "strings"
    "strconv"
    "github.com/mr-tron/base58"
    "encoding/json"
    "github.com/cosmos/cosmos-sdk/codec"
    sdk "github.com/cosmos/cosmos-sdk/types"
    sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
    abci "github.com/tendermint/tendermint/abci/types"
    "github.com/yzhanginwa/dbchain/x/dbchain/internal/utils"
)

type Condition struct {
    Field string
    Operator string
    Value string
}

type QuerierBuilder struct {
    Table string
    Ids []uint
    Select []string
    Where []Condition
    Order []string
    Limit int
    Last bool
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
    builders := []QuerierBuilder{}
    j := -1

    for i := 0; i < len(querierObjs); i++ {
        qo := querierObjs[i]
        switch qo["method"] {
        case "table":
            builders = append(builders, QuerierBuilder{})
            j += 1
            builders[j].Table = qo["table"]
        case "select":
            fields := strings.Split(qo["fields"], ",")
            builders[j].Select = fields
        case "find":
           id, err := strconv.Atoi(qo["id"])
           if err != nil {
               return nil, err
           }
           builders[j].Ids = []uint{uint(id)}
        case "first":
           builders[j].Limit = 1
        case "last":
           builders[j].Last = true
        case "where":
           cond := Condition{
               Field: qo["field"],
               Operator: qo["operator"],
               Value: qo["value"],
           }
           builders[j].Where = append(builders[j].Where, cond)
        }
    }

    ids := []uint{}
    for j = 0; j < len(builders); j++ {
        if len(ids) == 0 && j > 0 {
            break
        }

        if j != 0 {
            preTable := builders[j-1].Table
            curTable := builders[j].Table
            fld1 := curTable + "_id"
            fld2 := preTable + "_id"
            if keeper.HasField(ctx, appId, preTable, fld1) {
                builders[j].Ids = getIdsFromLeftToRight(ctx, keeper, appId, preTable, ids, fld1)
            } else if keeper.HasField(ctx, appId, curTable, fld2) {
                builders[j].Ids = getIdsFromRightToLeft(ctx, keeper, appId, ids, curTable, fld2, owner)
            } else {
                return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "Association does not exist!")
            }
        }

        if len(builders[j].Where) == 0 {
            if 0 == j && 0 == len(builders[j].Ids) {
                builders[j].Ids = keeper.FindAll(ctx, appId, builders[j].Table, owner)
            }
        } else {
            // to get the intersect of the result ids of all the where clauses
            intersect := []uint{}
            for index, cond := range builders[j].Where {
                ids := keeper.Where(ctx, appId, builders[j].Table, cond.Field, cond.Operator, cond.Value, owner)
                if index == 0 {
                    intersect = ids
                } else {
                    new_intersect := []uint{}
                    for _, a := range intersect {
                        for _, b := range ids {
                            if a == b {
                                new_intersect = append(new_intersect, a)
                            }
                        }
                    }
                    intersect = new_intersect
                }
            }
            builders[j].Ids = intersect
        }

        if builders[j].Last {
            length := len(builders[j].Ids)
            if length > 0 {
                ids = builders[j].Ids[length-1:]
            }
        } else {
            if builders[j].Limit == 0 || builders[j].Limit >= len(builders[j].Ids) {
                ids = builders[j].Ids
            } else {
                ids = builders[j].Ids[:(builders[j].Limit)]
            }
        }
    }

    j -= 1
    if len(builders[j].Select) == 0 {
        table, err := keeper.GetTable(ctx, appId, builders[j].Table)
        if err != nil {
            return nil, err
        }
        builders[j].Select = table.Fields
    }

    store := DbChainStore(ctx, keeper.storeKey)
    var result = [](map[string]string){}
    for _, id := range ids {
        record := map[string]string{}
        for _, f := range builders[j].Select {
            key := getDataKeyBytes(appId, builders[j].Table, f, id)
            bz, err := store.Get(key)
            if err != nil{
                return nil, err
            }
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

func getIdsFromLeftToRight(ctx sdk.Context, keeper Keeper, appId uint, preTable string, ids []uint, field string) []uint {
    store := DbChainStore(ctx, keeper.storeKey)
    var result []uint
    for i:= 0; i < len(ids); i++ {
        key := getDataKeyBytes(appId, preTable, field, ids[i])
        bz, err := store.Get(key)
        if err != nil{
            return nil
        }
        var value string
        if bz != nil {
            keeper.cdc.MustUnmarshalBinaryBare(bz, &value)
            ptrId, err := strconv.Atoi(value)
            if err == nil {
                result = append(result, uint(ptrId))
            }
        }
    }
    return result
}

func getIdsFromRightToLeft(ctx sdk.Context, keeper Keeper, appId uint, ids []uint, curTable string, field string, owner sdk.AccAddress) []uint {
    var values []string
    for i := 0; i < len(ids); i ++ {
        values = append(values, fmt.Sprintf("%d", ids[i]))
    }
    return keeper.FindBy(ctx, appId, curTable, field,  values, owner)
}
