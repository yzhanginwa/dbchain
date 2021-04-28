package keeper

import (
	"encoding/json"
	sdk "github.com/cosmos/cosmos-sdk/types"
	lua "github.com/yuin/gopher-lua"
	"strconv"
)

func getGoExportFuncNew(ctx sdk.Context, appId uint, keeper Keeper, owner sdk.AccAddress) map[string]lua.LGFunction {
	return map[string]lua.LGFunction {
		"InsertRow" : func(L *lua.LState) int {
			ParamsNum := L.GetTop()
			if ParamsNum < 1 {
				L.Push(lua.LNumber(-1))
				L.Push(lua.LString("Params Err"))
				return 2
			}
			param := L.ToString(1)
			var insertRowData ScriptInsertRow
			err := json.Unmarshal([]byte(param), &insertRowData)
			if err != nil {
				L.Push(lua.LNumber(-1))
				L.Push(lua.LString("Params unmarshal failed"))
				return 2
			}

			Id, err := keeper.Insert(ctx, appId, insertRowData.TableName, insertRowData.Fields, owner)
			if err != nil {
				L.Push(lua.LNumber(-1))
				L.Push(lua.LString(err.Error()))
			} else {
				L.Push(lua.LNumber(Id))
				L.Push(lua.LString(""))
			}
			return 2
		},

		"ForeignInsertRow" : func(L *lua.LState) int {
			ParamsNum := L.GetTop()
			if ParamsNum < 1 {
				L.Push(lua.LNumber(-1))
				L.Push(lua.LString("Params Err"))
				return 2
			}
			param := L.ToString(1)
			var foreignInsertRowData ScriptForeignInsertRow
			err := json.Unmarshal([]byte(param), &foreignInsertRowData)
			if err != nil {
				L.Push(lua.LNumber(-1))
				L.Push(lua.LString("Params unmarshal failed"))
				return 2
			}
			fields := foreignInsertRowData.Fields
			for i, v := range foreignInsertRowData.ForeignKey {
				value := L.ToString(i+1)
				fields[v] = value
			}

			Id, err := keeper.Insert(ctx, appId, foreignInsertRowData.TableName, fields, owner)
			if err != nil {
				L.Push(lua.LNumber(-1))
				L.Push(lua.LString(err.Error()))
			} else {
				L.Push(lua.LNumber(Id))
				L.Push(lua.LString(""))
			}
			return 2
		},
		"FreezeMultRow" : func(L *lua.LState) int {
			ParamsNum := L.GetTop()
			if ParamsNum < 1 {
				L.Push(lua.LString("num of param wrong"))
				return 1
			}
			param := L.ToString(1)
			var freezeMultRow ScriptFreezeMultRow
			err := json.Unmarshal([]byte(param), &freezeMultRow)
			if err != nil {
				L.Push(lua.LString("Param unmarshal failed"))
				return 1
			}

			for _,id := range freezeMultRow.Ids {
				id, err  := strconv.Atoi(id)
				if err != nil {
					continue
				}
				keeper.Freeze(ctx, appId, freezeMultRow.TableName, uint(id), owner)
			}
			L.Push(lua.LString(""))
			return 1
		},
		"FreezeMultRowByField" : func(L *lua.LState) int {
			ParamsNum := L.GetTop()
			if ParamsNum < 1 {
				L.Push(lua.LString("num of param wrong"))
				return 1
			}
			param := L.ToString(1)
			var freezeMultRow ScriptFreezeMultRowByField
			err := json.Unmarshal([]byte(param), &freezeMultRow)
			if err != nil {
				L.Push(lua.LString("Param unmarshal failed"))
				return 1
			}
			fields, values := []string{},[]string{}
			for k,v := range freezeMultRow.Fields {
				fields = append(fields, k)
				values = append(values, v)
			}
			_, ids := findByFields(keeper, ctx, appId, owner, freezeMultRow.TableName, fields, values)
			for _, id := range ids {
				keeper.Freeze(ctx, appId, freezeMultRow.TableName, id, owner)
			}
			L.Push(lua.LString(""))
			return 1
		},
		//add other functions which need to be exported
	}
}


