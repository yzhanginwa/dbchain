package keeper

import (
	"errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/yzhanginwa/dbchain/x/dbchain/internal/types"
	"strconv"
	"strings"
	"time"
)

const (
	INHERIT     = "inherit"
	INHERITABLE = "inheritable"

)
func (k Keeper) specFindAll(ctx sdk.Context, appId uint, tableName string, user sdk.AccAddress) []uint {
	store := DbChainStore(ctx, k.storeKey)
	var result []uint

	// full table scanning
	start, end := getFieldDataIteratorStartAndEndKey(appId, tableName, "id")
	iter := store.Iterator([]byte(start), []byte(end))
	for ; iter.Valid(); iter.Next() {
		if iter.Error() != nil{
			return nil
		}
		key := iter.Key()
		id := getIdFromDataKey(key)
		if isRowFrozen(store, appId, tableName, id) {
			continue;
		}
		result = append(result, id)
	}
	return result
}


func (k Keeper) specFind(ctx sdk.Context, appId uint, tableName string, id uint, user sdk.AccAddress) (types.RowFields, error){
	return k.DoFind(ctx, appId, tableName, id)
}

func (k Keeper) specWhere(ctx sdk.Context, appId uint, tableName string, field string, operator string, value string, user sdk.AccAddress) []uint {
	//TODO: consider if the field has index and how to make use of it
	store := DbChainStore(ctx, k.storeKey)
	isInteger := k.isTypeOfInteger(ctx, appId, tableName, field)
	results := []uint{}

	start, end := getFieldDataIteratorStartAndEndKey(appId, tableName, field)
	iter := store.Iterator([]byte(start), []byte(end))
	var mold string
	for ; iter.Valid(); iter.Next() {
		if iter.Error() != nil{
			return nil
		}
		key := iter.Key()
		val := iter.Value()
		k.cdc.MustUnmarshalBinaryBare(val, &mold)

		matching := fieldValueCompare(isInteger, operator, mold, value)
		if matching {
			id := getIdFromDataKey(key)
			if isRowFrozen(store, appId, tableName, id) {
				continue;
			}
			results = append(results, id)
		}
	}
	return results
}


func specQuerierSuperHandler(ctx sdk.Context, keeper Keeper, appId uint, querierObjs [](map[string]string), owner sdk.AccAddress) ([](map[string]string), []uint, error) {
	builders := []QuerierBuilder{}
	j := -1

	for i := 0; i < len(querierObjs); i++ {
		qo := querierObjs[i]
		switch qo["method"] {
		case "table":
			//only support
			if qo["table"] != INHERIT && qo["table"] != INHERITABLE {
				return nil, nil , errors.New("only support inherit or inheritable")
			}
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
				return nil, nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "Association does not exist!")
			}
		}

		if len(builders[j].Where) == 0 {
			if 0 == j && 0 == len(builders[j].Ids) {
				builders[j].Ids = keeper.specFindAll(ctx, appId, builders[j].Table, owner)
			}
		} else {
			// to get the intersect of the result ids of all the where clauses
			intersect := []uint{}
			for index, cond := range builders[j].Where {
				ids := keeper.specWhere(ctx, appId, builders[j].Table, cond.Field, cond.Operator, cond.Value, owner)
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
			return nil, nil, err
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
				return nil, nil, err
			}
			var value string
			if bz != nil {
				keeper.cdc.MustUnmarshalBinaryBare(bz, &value)
				record[f] = value
			}
		}
		result = append(result, record)
	}
	return result, ids, nil
}

func checkResult(ctx sdk.Context, owner sdk.AccAddress, keeper Keeper, appId uint, querierObjs ,src []map[string]string) []map[string]string {
	result := make([]map[string]string, 0)
	nowTime := time.Now().UnixNano()/(1000*1000)
	tableName := ""
	for _, qo := range querierObjs {
		method , ok := qo["method"]
		if ok && method == "table" {
			tableName = qo["table"]
			break
		}
	}

	if tableName == INHERIT {
		for _, val := range src {
			if owner.String() == val["created_by"] {	//query by creator
				result = append(result, val)
			} else if owner.String() == val["receive_address"] {	//query by receiver
				receive_time, err := strconv.ParseInt(val["receive_time"],10,64)
				if err != nil {
					continue
				}
				if nowTime >= receive_time {
					result = append(result, val)
				}
			} else {
				continue
			}
		}
	} else {
		for _, val := range src {
			if owner.String() == val["created_by"] {
				result = append(result, val)
			} else {
				inheritableid := val["inheritableid"]
				id , _ := strconv.Atoi(inheritableid)
				row, err := keeper.DoFind(ctx, appId, INHERIT, uint(id))
				if err != nil || row["receive_address"] != owner.String(){
					continue
				}
				receive_time, err := strconv.ParseInt(row["receive_time"],10,64)
				if err != nil {
					continue
				}
				if nowTime < receive_time {
					continue
				}
				result = append(result, val)
			}
		}
	}

	return result
}