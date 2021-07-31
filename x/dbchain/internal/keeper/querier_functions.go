package keeper

import (
	"errors"
	"fmt"
	sdk "github.com/dbchaincloud/cosmos-sdk/types"
	sdkerrors "github.com/dbchaincloud/cosmos-sdk/types/errors"
	"github.com/yzhanginwa/dbchain/x/dbchain/internal/types"
	"github.com/yzhanginwa/dbchain/x/dbchain/internal/utils"
	"regexp"
	"strconv"
	"strings"
)


func (k Keeper) querierFindAll(ctx sdk.Context, appId uint, tableName string, user sdk.AccAddress) []uint {
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

	if k.isTablePublic(ctx, appId, tableName) || k.isAuditor(ctx, appId, user) {
		return result
	} else {
		return k.filterReadableIds(ctx, appId, tableName, result, user)
	}
}

func (k Keeper) queroerFind(ctx sdk.Context, appId uint, tableName string, id uint, user sdk.AccAddress) (types.RowFields, error){
	if !k.querierIsReadableId(ctx, appId, tableName, id, user) {
		return nil, errors.New(fmt.Sprintf("Failed to get fields for id %d", id))
	}

	return k.DoFind(ctx, appId, tableName, id)
}

func (k Keeper) querierIsReadableId(ctx sdk.Context, appId uint, tableName string, id uint, user sdk.AccAddress) bool{
	var ids []uint
	ids = append(ids, id)

	// if public table, return all ids
	if !k.isTablePublic(ctx, appId, tableName) && !k.isAuditor(ctx, appId, user) {
		ids = k.filterReadableIds(ctx, appId, tableName, ids, user)
		if len(ids) < 1 {
			return false
		}
	}
	return true
}
//
//func (k Keeper) specFind(ctx sdk.Context, appId uint, tableName string, id uint, user sdk.AccAddress) (types.RowFields, error){
//	return k.DoFind(ctx, appId, tableName, id)
//}

func (k Keeper) querierWhere(ctx sdk.Context, appId uint, tableName string, field string, operator string, value string, reg *regexp.Regexp, user sdk.AccAddress) []uint {
	//TODO: consider if the field has index and how to make use of it
	store := DbChainStore(ctx, k.storeKey)
	isInteger := k.isTypeOfInteger(ctx, appId, tableName, field)
	results := []uint{}
	if field == "id" && (operator ==  "==" || operator ==  "=") {
		id , err := strconv.ParseUint(value, 10, 32)
		if err != nil || !k.querierIsReadableId(ctx, appId, tableName, uint(id), user) {
			return results
		}
		results = append(results, uint(id))
		return results
	}


	if k.isIndexField(ctx, appId, tableName, field) {
		if operator ==  "==" || operator ==  "="  {
			var mold []string
			key := getIndexKey(appId, tableName, field, value)
			bz, err := store.Get([]byte(key))
			if err != nil{
				return results
			}
			if bz != nil {
				k.cdc.MustUnmarshalBinaryBare(bz, &mold)
			}
			for _, sId := range mold {
				id , err := strconv.ParseUint(sId, 10, 32)
				if err != nil {
					continue
				}
				results = append(results, uint(id))
			}
			return results
		} else {
			start, end := getIndexDataIteratorStartAndEndKey(appId, tableName, field)
			iter := store.Iterator([]byte(start), []byte(end))
			for ; iter.Valid(); iter.Next() {
				if iter.Error() != nil {
					continue
				}
				key := iter.Key()
				val := iter.Value()
				sliceKey := strings.Split(string(key),":")
				matching := fieldValueCompare(isInteger, operator, sliceKey[len(sliceKey)-1], value, reg)
				if matching {
					var mold []string
					k.cdc.MustUnmarshalBinaryBare(val, &mold)
					for _, sId := range mold {
						id , err := strconv.ParseUint(sId, 10, 32)
						if err != nil {
							continue
						}
						results = append(results, uint(id))
					}
				}
			}
			return results
		}
	}

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

		matching := fieldValueCompare(isInteger, operator, mold, value, reg)
		if matching {
			id := getIdFromDataKey(key)
			if isRowFrozen(store, appId, tableName, id) {
				continue;
			}
			results = append(results, id)
		}
	}

	if k.isTablePublic(ctx, appId, tableName) || k.isAuditor(ctx, appId, user) {
		return results
	} else {
		return k.filterReadableIds(ctx, appId, tableName, results, user)
	}
}


func customQuerierSuperHandler(ctx sdk.Context, keeper Keeper, appId uint, querierObjs [](map[string]string), owner sdk.AccAddress) (*WhereRes, []uint, error) {
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
				ids = keeper.querierFindAll(ctx, appId, builders[j].Table, owner)
			}
		} else {
			// to get the intersect of the result ids of all the where clauses and ids (if there are any)
			intersect := ids
			for index, cond := range builders[j].Where {
				tmp_ids := keeper.querierWhere(ctx, appId, builders[j].Table, cond.Field, cond.Operator, cond.Value, cond.Reg, owner)
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
