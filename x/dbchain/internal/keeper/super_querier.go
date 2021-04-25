package keeper

import (
    "fmt"
    "regexp"
    "strings"
    "strconv"
    "sort"
    "github.com/mr-tron/base58"
    "encoding/json"
    "github.com/cosmos/cosmos-sdk/codec"
    sdk "github.com/cosmos/cosmos-sdk/types"
    sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
    abci "github.com/tendermint/tendermint/abci/types"
    "github.com/yzhanginwa/dbchain/x/dbchain/internal/utils"
    "github.com/yzhanginwa/dbchain/x/dbchain/internal/types"
)

type Condition struct {
    Field string
    Operator string
    Value string
    Reg   *regexp.Regexp
}

type OrderBy struct {
    Field string
    Direction string
}

type QuerierBuilder struct {
    Table string
    Ids []uint
    Select []string
    Where []Condition
    Order OrderBy
    Offset int
    Limit int
    Last bool
}

type WhereRes struct {
    Data []map[string]string `json:"data"`
    Count string             `json:"count"`
}

//////////////////////////////////
//                              //
// implement the ByValue sorter //
//                              //
//////////////////////////////////

type idAndValues struct {
    Id uint
    IsInt bool
    StringValue string
    IntValue int
}

type ByValue []idAndValues

func (a ByValue) Len() int           { return len(a) }
func (a ByValue) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByValue) Less(i, j int) bool {
    if a[i].IsInt {
        return a[i].IntValue < a[j].IntValue
    } else {
        return a[i].StringValue < a[j].StringValue
    }
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

    result, _, err := querierSuperHandler(ctx, keeper, appId, querierObjs, addr)
    if err != nil {
        return nil, err
    }

    var res []byte
    if isCountQuery(querierObjs) {
        temp := map[string]string {
            "count" : result.Count,
        }
        res, err = codec.MarshalJSONIndent(keeper.cdc, temp)
    } else {
        res, err = codec.MarshalJSONIndent(keeper.cdc, result.Data)
    }

    if err != nil {
        panic("could not marshal result to JSON")
    }
    return res, nil
}

func querierSuperHandler(ctx sdk.Context, keeper Keeper, appId uint, querierObjs [](map[string]string), owner sdk.AccAddress) (*WhereRes, []uint, error) {
    whereRes := &WhereRes{}

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
                return nil, nil, err
            }
            builders[j].Ids = []uint{uint(id)}
        case "first":
            builders[j].Limit = 1
        case "limit", "offset":
            val, err := strconv.Atoi(qo["value"])
            if err != nil || val < 0 {
                return nil, nil, err
            }
            if qo["method"] == "limit" {
                builders[j].Limit = val
            } else {
                builders[j].Offset = val
            }

        case "last":
            builders[j].Last = true
        case "where":
            cond := Condition{
                Field: qo["field"],
                Operator: qo["operator"],
                Value: qo["value"],
            }
            if qo["operator"] == "like" {
                reg, err := utils.DealFuzzyQueryString(qo["value"])
                if err != nil {
                    return nil, nil, err
                }
                cond.Reg = reg
            }
            builders[j].Where = append(builders[j].Where, cond)
        case "order":
            builders[j].Order = OrderBy{
                Field: qo["field"],
                Direction: qo["direction"],
           }
        }
    }

    ids := []uint{}
    count := 0
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
                return nil, nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "Association does not exist!")
            }
        }

        if len(builders[j].Ids) > 0 {
            ids = builders[j].Ids
        }

        if len(builders[j].Where) == 0 {
            if 0 == j && 0 == len(ids) {
                ids = keeper.FindAll(ctx, appId, builders[j].Table, owner)
            }
        } else {
            // to get the intersect of the result ids of all the where clauses and ids (if there are any)
            intersect := ids
            for index, cond := range builders[j].Where {
                tmp_ids := keeper.Where(ctx, appId, builders[j].Table, cond.Field, cond.Operator, cond.Value, cond.Reg, owner)
                if index == 0 && len(intersect) == 0 {
                    intersect = tmp_ids
                } else {
                    new_intersect := []uint{}
                    for _, a := range intersect {
                        for _, b := range tmp_ids {
                            if a == b {
                                new_intersect = append(new_intersect, a)
                            }
                        }
                    }
                    intersect = new_intersect
                }
            }
            ids = intersect
        }

        count += len(ids)   // this might need more consideration

        if builders[j].Order.Field != "" {
            ids = sortIdsOnOrder(ctx, keeper, appId, ids, builders[j].Table, builders[j].Order.Field, builders[j].Order.Direction)
        }

        if builders[j].Last {
            length := len(ids)
            if length > 0 {
                ids = ids[length-1:]
            }
        } else if builders[j].Offset == 0 {
            if builders[j].Limit > 0 && builders[j].Limit < len(ids) {
                ids = ids[:(builders[j].Limit)]
            }
        } else {
            if builders[j].Offset >= len(ids) {
                ids = []uint{}
            } else if builders[j].Limit == 0 || builders[j].Limit + builders[j].Offset >= len(builders[j].Ids) {
                ids = ids[builders[j].Offset :]
            } else {
                ids = ids[builders[j].Offset : builders[j].Offset + (builders[j].Limit)]
            }
        }
    }

    if isCountQuery(querierObjs) {
        whereRes.Count = strconv.Itoa(count)
        return whereRes, nil, nil
    }

    j -= 1
    if len(builders[j].Select) == 0 {
        table, err := keeper.GetTable(ctx, appId, builders[j].Table)
        if err != nil {
            return nil, nil, err
        }
        builders[j].Select = table.Fields
    }

    store := DbChainStore(ctx, keeper.storeKey)
    var validId = make([]uint,0)
    var result = [](map[string]string){}
    for _, id := range ids {
        record := map[string]string{}
        for _, f := range builders[j].Select {
            key := getDataKeyBytes(appId, builders[j].Table, f, id)
            bz, err := store.Get(key)
            if err != nil{
                return nil, nil, err
            }
            var value string
            if bz != nil {
                keeper.cdc.MustUnmarshalBinaryBare(bz, &value)
                record[f] = value
            }
        }
        if len(record) > 0 {
            validId = append(validId, id)
            result = append(result, record)
        }
    }
    whereRes.Data = result
    return whereRes, validId, nil
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

func isCountQuery (querierObjs []map[string]string) bool {
    for _, m := range querierObjs {
        if m["method"] == "count" {
            return true
        }
    }
    return false
}

func isColumnTypeOfInteger(ctx sdk.Context, keeper Keeper, appId uint, tableName, fieldName string) bool {
    fieldDataType, _ := keeper.GetColumnDataType(ctx, appId, tableName, fieldName)
    return (fieldDataType == string(types.FLDTYP_INT))
}

func sortIdsOnOrder(ctx sdk.Context, keeper Keeper, appId uint, ids []uint, tableName, fieldName, direction string) []uint {
    store := DbChainStore(ctx, keeper.storeKey)
    records := []idAndValues{}
    isInt := isColumnTypeOfInteger(ctx, keeper, appId, tableName, fieldName)

    for _, id := range ids {
        key := getDataKeyBytes(appId, tableName, fieldName, id)
        bz, err := store.Get(key)
        if err != nil {
            return []uint{}
        }
        var value string
        if bz != nil {
            record := idAndValues{Id: id, IsInt: false}
            keeper.cdc.MustUnmarshalBinaryBare(bz, &value)
            if isInt {
                intValue, err := strconv.Atoi(value)
                if err != nil { continue }
                record.IntValue = intValue
                record.IsInt = true
            } else {
                record.StringValue = value
            }
            records = append(records, record)
        }
    }

    sort.Sort(ByValue(records))
    result := []uint{}
    for _, record := range records {
        result = append(result, record.Id)
    }

    if direction == "desc" {
        for i, j := 0, len(result)-1; i < j; i, j = i+1, j-1 {
            result[i], result[j] = result[j], result[i]
        }
    }
    return result
}
