package keeper

import (
	"encoding/json"
	storeTypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	lua "github.com/yuin/gopher-lua"
	"strconv"
	"time"
)

func getGoExportFuncNew(ctx sdk.Context, appId uint, keeper Keeper, owner sdk.AccAddress) map[string]lua.LGFunction {
	return map[string]lua.LGFunction {
		"InsertRow" : func(L *lua.LState) int {
			//params : 1. tableName, 2. fields
			ParamsNum := L.GetTop()
			if ParamsNum < 2 {
				L.Push(lua.LString("-1"))
				L.Push(lua.LString("Params Err"))
				return 2
			}
			tableName := L.CheckString(1)
			fields := luaTableToGoMap(L.CheckTable(2))

			Id, err := keeper.Insert(ctx, appId, tableName, fields, owner)
			if err != nil {
				L.Push(lua.LString("-1"))
				L.Push(lua.LString(err.Error()))
			} else {
				L.Push(lua.LString(Id))
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
			if ParamsNum < 2 {
				L.Push(lua.LString("num of param wrong"))
				return 1
			}
			tableName := L.CheckString(1)
			Ids := luaArrayToGoArray(L.CheckTable(2))

			for _,id := range Ids {
				id, err  := strconv.Atoi(id)
				if err != nil {
					continue
				}
				keeper.Freeze(ctx, appId, tableName, uint(id), owner)
			}
			L.Push(lua.LString(""))
			return 1
		},
		"FreezeMultRowByField" : func(L *lua.LState) int {
			ParamsNum := L.GetTop()
			if ParamsNum < 2 {
				L.Push(lua.LString("num of param wrong"))
				return 2
			}
			tableName := L.CheckString(1)
			Fields := luaTableToGoMap(L.CheckTable(2))
			fields, values := []string{},[]string{}
			for k,v := range Fields {
				fields = append(fields, k)
				values = append(values, v)
			}
			_, ids := findByFields(keeper, ctx, appId, owner, tableName, fields, values)
			for _, id := range ids {
				keeper.Freeze(ctx, appId, tableName, id, owner)
			}
			L.Push(lua.LString(""))
			return 1
		},
		"itemExists" : func(L *lua.LState) int {
			paramsNum := L.GetTop()
			if paramsNum < 2 {
				L.Push(lua.LFalse)
				return 1
			}
			validTables := getValidTables(L)
			target := L.ToString(2)
			if validTables[target] {
				L.Push(lua.LTrue)
			} else {
				L.Push(lua.LFalse)
			}
			return 1
		},
		//add other functions which need to be exported
	}
}

func getGoExportQueryFuncNew(ctx sdk.Context, appId uint, keeper Keeper, addr sdk.AccAddress) map[string]lua.LGFunction {
	return map[string]lua.LGFunction {
		"findRowsQuerier": func(L *lua.LState) int {

			validResult := make([]map[string]string,0)
			paramsNum := L.GetTop()
			if paramsNum < 2 {
				ud := setUserData(ctx, appId, keeper, addr, "", validResult, L)
				L.Push(ud)
				return 1
			}
			//get valid Tables
			tableName := L.ToString(1)
			var querierObjs  = luaArrayTableToGoArrayMap(L.CheckTable(2))

			querierTableName := ""
			for _, qo := range querierObjs {
				method , ok := qo["method"]
				if ok && method == "table" {
					querierTableName = qo["table"]
					break
				}
			}
			if tableName != querierTableName{
				ud := setUserData(ctx, appId, keeper, addr, tableName, validResult, L)
				L.Push(ud)
				return 1
			}

			//query
			result, _, err := customQuerierSuperHandler(ctx, keeper, appId, querierObjs, addr)
			if err != nil {
				ud := setUserData(ctx, appId, keeper, addr, tableName, validResult, L)
				L.Push(ud)
				return 1
			}

			//check
			checkTimeField := ""
			if L.GetTop() > 2 {
				checkTimeField = L.ToString(3)
			}

			if checkTimeField != "" {
				validResult = checkTime(checkTimeField, result.Data)
			} else {
				validResult = result.Data
			}
			ud := setUserData(ctx, appId, keeper, addr, tableName, validResult, L)
			L.Push(ud)
			return 1
		},
		"findRow" : func(L *lua.LState) int {
			res := make([]map[string]string, 0)
			paramsNum := L.GetTop()
			if paramsNum < 2 {
				ud := setUserData(ctx, appId, keeper, addr, "", res, L)
				L.Push(ud)
				return 1
			}

			tableName := L.ToString(1)
			sId := L.ToString(2)
			Id , err := strconv.Atoi(sId)
			if err != nil {
				ud := setUserData(ctx, appId, keeper, addr, "", res, L)
				L.Push(ud)
				return 1
			}
			checkField := ""
			if L.GetTop() > 2 {
				checkField = L.ToString(3)
			}

			fields, err := keeper.queroerFind(ctx, appId, tableName, uint(Id), addr)
			if err != nil {
				ud := setUserData(ctx, appId, keeper, addr, tableName, res, L)
				L.Push(ud)
				return 1
			}
			if checkField != "" {
				checkRes := checkTime(checkField, []map[string]string{fields})
				if len(checkRes) > 0 {
					res = append(res, fields)
				}
			} else {
				res = append(res, fields)
			}

			ud := setUserData(ctx, appId, keeper, addr, tableName, res, L)
			L.Push(ud)
			return 1
		},
		"findRowsBy" : func(L *lua.LState) int {
			res := make([]map[string]string, 0)
			paramsNum := L.GetTop()
			if paramsNum < 2 {
				ud := setUserData(ctx, appId, keeper, addr, "", res, L)
				L.Push(ud)
				return 1
			}

			tableName := L.ToString(1)

			var Fields  = luaTableToGoMap(L.CheckTable(2))

			checkTimeField := ""
			if L.GetTop() > 2 {
				checkTimeField = L.ToString(3)
			}
			querierObjs := makeWhereEqualQuerierObjs(tableName, Fields)
			//query
			result, _, err := customQuerierSuperHandler(ctx, keeper, appId, querierObjs, addr)
			if err != nil {
				ud := setUserData(ctx, appId, keeper, addr, tableName, res, L)
				L.Push(ud)
				return 1
			}

			validResult := make([]map[string]string,0)
			if checkTimeField != "" {
				validResult = checkTime(checkTimeField, result.Data)
			} else {
				validResult = result.Data
			}
			ud := setUserData(ctx, appId, keeper, addr, tableName, validResult, L)
			L.Push(ud)
			return 1
		},
		"itemExists" : func(L *lua.LState) int {
			paramsNum := L.GetTop()
			if paramsNum < 2 {
				L.Push(lua.LFalse)
				return 1
			}
			validTables := getValidTables(L)
			target := L.ToString(2)
			if validTables[target] {
				L.Push(lua.LTrue)
			} else {
				L.Push(lua.LFalse)
			}
			return 1
		},
	}
}

func getGoExportToolFunc(ctx sdk.Context) map[string]lua.LGFunction {
	return map[string]lua.LGFunction {
		"jsonStringToArray" : func(L *lua.LState) int {
			params := L.CheckString(1)
			var array []string
			err := json.Unmarshal([]byte(params), &array)
			if err != nil {
				L.Push(lua.LNil)
				return 1
			}
			table := L.NewTable()
			for i, v := range array {
				table.RawSetInt(i+1, lua.LString(v))
			}
			L.Push(table)
			return 1
		},
		"jsonStringToMap": func(L *lua.LState) int {
			params := L.CheckString(1)
			var data = make(map[string]string, 0)
			err := json.Unmarshal([]byte(params), &data)
			if err != nil {
				L.Push(lua.LNil)
				return 1
			}
			table := L.NewTable()
			for k, v := range data {
				table.RawSetString(k, lua.LString(v))
			}
			L.Push(table)
			return 1
		},
		"jsonStringToArrayMap" : func(L *lua.LState) int {
			params := L.CheckString(1)
			var data []map[string]string
			err := json.Unmarshal([]byte(params), &data)
			if err != nil {
				L.Push(lua.LNil)
				return 1
			}
			table := L.NewTable()
			for i, m := range data {
				temp := L.NewTable()
				for k, v := range m {
					temp.RawSetString(k, lua.LString(v))
				}
				table.RawSetInt(i+1, temp)
			}
			L.Push(table)
			return 1
		},
		"scriptConsumeGas" : func(L *lua.LState) int {
			gasNum := L.ToInt64(1)
			if gasNum == 0 {
				gasNum = 1000
			}
			gas := storeTypes.Gas(int64(gasNum))
			ctx.GasMeter().ConsumeGas(gas,"script consume")
			return 0
		},
	}
}


func checkTime(checkTimeField string ,src []map[string]string) []map[string]string {
	result := make([]map[string]string, 0)
	nowTime := time.Now().UnixNano()/(1000*1000)

	if checkTimeField != "" {
		for _, row := range src {
			canReadTime, err := strconv.ParseInt(row[checkTimeField],10,64)
			if err != nil {
				continue
			}
			if nowTime < canReadTime {
				continue
			}
			result = append(result, row)
		}
		return result
	}
	return src
}



///////////////////////////////////////
//                                   //
//           help func               //
//                                   //
///////////////////////////////////////

func getValidTables(L *lua.LState) map[string]bool {
	//get valid Tables
	lTables := L.ToTable(1)
	validTables := make(map[string]bool)
	lTables.ForEach(func(lId lua.LValue, val lua.LValue) {
		table := val.(lua.LString)
		validTables[table.String()] = true
	})
	return validTables
}

func makeWhereEqualQuerierObjs(tableName string, fields map[string]string) []map[string]string {

	querierObjs := []map[string]string{}
	var ent map[string]string
	ent = map[string]string{
		"method": "table",
		"table":  tableName,
	}
	querierObjs = append(querierObjs, ent)

	for field, val := range fields {
		ent = map[string]string{
			"method":   "where",
			"field":    field,
			"value":    val,
			"operator": "==",
		}
		querierObjs = append(querierObjs, ent)
	}
	return querierObjs
}

func luaArrayToGoArray(table *lua.LTable) []string{
	array := make([]string, 0)
	table.ForEach(func(index lua.LValue, val lua.LValue) {
		element := val.String()
		array = append(array, element)
	})
	return array
}

func luaTableToGoMap(table *lua.LTable) map[string]string{
	data := make(map[string]string, 0)
	table.ForEach(func(key lua.LValue, val lua.LValue) {
		mapKey := key.String()
		element := val.String()
		data[mapKey] = element
	})
	return data
}

func luaArrayTableToGoArrayMap(arrayTable *lua.LTable) []map[string]string{
	data := make([]map[string]string, 0)
	arrayTable.ForEach(func(index lua.LValue, table lua.LValue) {
		temp := make(map[string]string)
		nTable := table.(*lua.LTable)
		nTable.ForEach(func(key lua.LValue, val lua.LValue) {
			temp[key.String()] = val.String()
		})
		data = append(data, temp)
	})
	return data
}