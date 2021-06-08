/*
这个文件主要用于导出go函数，用来给lua脚本调用
要导出的go函数必须申明为func(L *lua.LState) int 格式
在脚本中调用带参go导出函数时，通过L.ToString(n)来获取参数，不同的类型可以用不同的函数获取，如L.ToInt(),参数n表示获取函数的第几个参数，如果需要将数据传出给lua脚本，使用L.Push()函数
下面是简单的示例
func Double(L *lua.LState) int {
	lv := L.ToInt(1)            // get argument
	L.Push(lua.LNumber(lv * 2)) // push result
	return 1                    // number of results
}
*/
package keeper

import (
	"encoding/json"
	"errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	lua "github.com/yuin/gopher-lua"
	"github.com/yzhanginwa/dbchain/x/dbchain/internal/types"
	"strconv"
	"strings"
	"time"
)

const (
	tablePrefix = "tableName__"
	foreignPrefix = "foreignKeyName__"
)

func getGoExportFunc(ctx sdk.Context, appId uint, keeper Keeper, owner sdk.AccAddress) map[string]lua.LGFunction {
	return map[string]lua.LGFunction{
		"Insert": func(L *lua.LState) int {
			ParamsNum := L.GetTop()
			if ParamsNum >= 2 && ParamsNum%2 == 0 { //Normal inserttab,fields
				tableName := L.ToString(1)
				if strings.HasPrefix(tableName, tablePrefix){
					tableName = strings.TrimPrefix(tableName, tablePrefix)
				}
				sFieldAndValues := L.ToString(2)
				fieldAndValues, err := getFieldValueMap(ctx, appId, keeper, tableName, sFieldAndValues)
				if err != nil {
					L.Push(lua.LNumber(-1))
					L.Push(lua.LString(err.Error()))
					return 2
				}
				if ParamsNum > 2 { //此时表示有外键插入，可以有多个外键插入，格式为foreigntab，foreignid 循环
					for i := 3; i < ParamsNum; i+=2{
						fTableName := L.ToString(i)
						fId := L.ToString(i+1)
						fKey := strings.ToLower(fTableName)
						fieldAndValues[fKey] = fId
					}

				}
				Id, err := keeper.Insert(ctx, appId, tableName, fieldAndValues, owner)
				if err != nil {
					L.Push(lua.LNumber(-1))
					L.Push(lua.LString(err.Error()))
					return 2
				}
				L.Push(lua.LNumber(Id))
				L.Push(lua.LString(""))
			} else {
				L.Push(lua.LNumber(-1))
				L.Push(lua.LString("num of param wrong"))
			}
			return 2
		},
		//插入一条数据，取得ID，然后这个id作为关联值，往下一张表中插入多个值
		"ForeignMultInsert" : func(L *lua.LState) int {
			//至少4个参数，分别为 id,tab1，foreignKey1，val, tab2，foreignKey2,val2...
			ParamsNum := L.GetTop()
			if ParamsNum < 4 {
				L.Push(lua.LNumber(-1))
				L.Push(lua.LString("num of param wrong"))
				return 2
			}
			foreignKeyId := L.ToString(1)
			tableName := L.ToString(2)
			if strings.HasPrefix(tableName, tablePrefix) {
				tableName = strings.TrimPrefix(tableName, tablePrefix)
			}
			foreignKeyName := L.ToString(3)
			if strings.HasPrefix(foreignKeyName, foreignPrefix) {
				foreignKeyName = strings.TrimPrefix(foreignKeyName, foreignPrefix)
			}
			count := 0
			for i := 4; i <= ParamsNum; i++ {
				sFieldAndValues := L.ToString(i)
				//change table
				if strings.HasPrefix(sFieldAndValues, tablePrefix) {
					tableName = strings.TrimPrefix(sFieldAndValues, tablePrefix)
					continue
				} else if strings.HasPrefix(sFieldAndValues,foreignPrefix) {
					//change foreign key
					foreignKeyName = strings.TrimPrefix(sFieldAndValues, foreignPrefix)
					continue
				}
				//insert value
				fieldAndValues, err := getFieldValueMap(ctx, appId, keeper, tableName, sFieldAndValues)
				if err != nil {
					L.Push(lua.LNumber(-1))
					L.Push(lua.LString(err.Error()))
					return 2
				}
				foreignKeyName = strings.ToLower(foreignKeyName)
				fieldAndValues[foreignKeyName] = foreignKeyId
				_, err = keeper.Insert(ctx, appId, tableName, fieldAndValues, owner)
				if err != nil {
					L.Push(lua.LNumber(-1))
					L.Push(lua.LString(err.Error()))
					return 2
				}
				count++
			}
			L.Push(lua.LNumber(count))
			L.Push(lua.LString(""))
			return 2
		},
		"MultInsert": func(L *lua.LState) int{ //往同一张表插入多条数据
			ParamsNum := L.GetTop()
			if ParamsNum < 2 {
				L.Push(lua.LNumber(-1))
				L.Push(lua.LString("num of param wrong"))
				return 2
			}
			//By default, the first parameter is the table name
			tableName := L.ToString(1)
			if strings.HasPrefix(tableName,tablePrefix) {
				tableName = strings.TrimPrefix(tableName,tablePrefix)
			}
			count := 0
			for i := 2 ; i <= ParamsNum; i++{
				sFieldAndValues := L.ToString(i)
				//change table
				if strings.HasPrefix(sFieldAndValues,tablePrefix) {
					tableName = strings.TrimPrefix(sFieldAndValues, tablePrefix)
					continue
				}
				fieldAndValues, err := getFieldValueMap(ctx, appId, keeper, tableName, sFieldAndValues)
				if err != nil {
					L.Push(lua.LNumber(count))
					L.Push(lua.LString(err.Error()))
					return 2
				}
				_, err = keeper.Insert(ctx, appId, tableName, fieldAndValues, owner)
				if err != nil {
					L.Push(lua.LNumber(count))
					L.Push(lua.LString(err.Error()))
					return 2
				}
				count++
			}
			L.Push(lua.LNumber(count))
			L.Push(lua.LString(""))
			return 2
		},
		"MultFreeze" : func(L *lua.LState) int {
			ParamsNum := L.GetTop()
			if ParamsNum < 2 {
				L.Push(lua.LString("num of param wrong"))
				return 1
			}
			tableName := L.ToString(1)
			//By default, the first parameter is the table name
			if strings.HasPrefix(tableName,tablePrefix) {
				tableName = strings.TrimPrefix(tableName,tablePrefix)
			}
			for i := 2 ; i <= ParamsNum; i++ {
				isTable := L.ToString(i)
				if strings.HasPrefix(isTable,tablePrefix) {
					tableName = strings.TrimPrefix(isTable,tablePrefix)
					continue
				}
				id, err  := strconv.Atoi(isTable)
				if err != nil {
					continue
				}
				keeper.Freeze(ctx, appId, tableName, uint(id), owner)
			}
			L.Push(lua.LString(""))
			return 1
		},
		"MultFreezeBy" : func(L *lua.LState) int {
			ParamsNum := L.GetTop()
			if ParamsNum < 3 {
				L.Push(lua.LString("num of param wrong"))
				return 1
			}
			tableName := L.ToString(1)
			//By default, the first parameter is the table name
			if strings.HasPrefix(tableName,tablePrefix) {
				tableName = strings.TrimPrefix(tableName,tablePrefix)
			}
			for i := 2; i <= ParamsNum; i++ {
				strField := L.ToString(i)
				if strings.HasPrefix(strField,tablePrefix) {
					tableName = strings.TrimPrefix(strField,tablePrefix)
					continue
				}
				fields := strings.Split(strField,",")
				i++
				var values []string
				strValue := L.ToString(i)
				err := json.Unmarshal([]byte(strValue), &values)
				if err != nil {
					L.Push(lua.LString("val of field err"))
					return 1
				}
				_, ids := findByFields(keeper, ctx, appId, owner, tableName, fields, values)
				for _, id := range ids {
					keeper.Freeze(ctx, appId, tableName, id, owner)
				}
			}

			L.Push(lua.LString(""))
			return 1
		},
		"MultFreezeExpiration" : func(L *lua.LState) int {
			ParamsNum := L.GetTop()
			if ParamsNum < 2 {
				L.Push(lua.LString("num of param wrong"))
				return 1
			}
			tableName := L.ToString(1)
			//By default, the first parameter is the table name
			if strings.HasPrefix(tableName,tablePrefix) {
				tableName = strings.TrimPrefix(tableName,tablePrefix)
			}
			days := L.ToInt(2)
			if days <= 0 {
				days = 0
			}
			t := time.Now()
			ids := keeper.FindAll(ctx, appId, tableName, owner)
			for index, id := range ids {
				res , err := keeper.Find(ctx, appId, tableName, id, owner)
				if err != nil  || res["created_by"] != owner.String(){
					continue
				}
				timeStamp := res["created_at"]
				dt, _ := time.ParseDuration(timeStamp + "ms")
				nt := time.Unix(dt.Milliseconds()/1000, 0)
				nt.Add(time.Hour * 24 * time.Duration(days))
				if t.Unix() > nt.Unix() {
					keeper.Freeze(ctx, appId, tableName, ids[index], owner)
				}
			}
			L.Push(lua.LString(""))
			return 1
		},
		//删除主表数据，同时删除从表数据
		"RelationDelete" : func(L *lua.LState) int {
			param := L.ToString(1)
			params := strings.Split(param, ",")
			if len(params) < 2  {
				L.Push(lua.LString("param err"))
				return 1
			}
			//删除主表
			tableName := params[0]
			if strings.HasPrefix(tableName,tablePrefix) {
				tableName = strings.TrimPrefix(tableName,tablePrefix)
			}
			id, err := strconv.Atoi(params[1])
			if err != nil{
				L.Push(lua.LString("id err"))
				return 1
			}
			keeper.Freeze(ctx, appId, tableName, uint(id), owner)
			tParams := splitParams(params[2:])
			for _, tParam := range tParams {
				var querierObjs []map[string]string
				if len(tParam) == 2 {
					querierObjs = []map[string]string{
						map[string]string{
							"method": "table",
							"table":  tParam[0],
						},
						map[string]string{
							"method":   "where",
							"field":    tParam[1],
							"value":    params[1],
							"operator": "==",
						},
					}
				} else if len(tParam) == 3 {
					querierObjs = []map[string]string{
						map[string]string{
							"method": "table",
							"table":  tParam[0],
						},
						map[string]string{
							"method":   "where",
							"field":    tParam[1],
							"value":    tableName,
							"operator": "==",
						},
						map[string]string{
							"method":   "where",
							"field":    tParam[2],
							"value":    params[1],
							"operator": "==",
						},

					}
				}
				_, ids, _ :=querierSuperHandler(ctx, keeper, appId, querierObjs, owner)
				for _, dId := range ids {
					keeper.Freeze(ctx, appId, tParam[0], dId, owner)
				}
			}
			L.Push(lua.LString(""))
			return 1
		},
		"fieldIn" : func(L *lua.LState) int {
			ParamsNum := L.GetTop()
			if ParamsNum < 2 {
				L.Push(lua.LBool(false))
				return 1
			}

			src := L.ToString(1)
			for i := 2; i <= ParamsNum; i++ {
				dst := L.ToString(i)
				if src == dst{
					L.Push(lua.LBool(true))
					return 1
				}
			}
			L.Push(lua.LBool(false))
			return 1
		},
		"exist" : func(L *lua.LState) int {
			ParamsNum := L.GetTop()
			if ParamsNum < 4 || (ParamsNum - 1)%3 != 0{
				L.Push(lua.LBool(false))
				return 1
			}
			tableName := L.ToString(1)
			qo := map[string]string{
				"method": "table",
				"table": tableName,
			}
			querierObjs := []map[string]string{qo}
			for i := 2; i < ParamsNum; i += 3 {
				field := L.ToString(i)
				op  := L.ToString(i+1)
				value := L.ToString(i+2)
				if op != "==" {
					continue
				}
				qo := map[string]string{
					"method": "where",
					"operator": "==",
					"field": field,
					"value": value,
				}
				querierObjs = append(querierObjs, qo)
			}

			tableValueCallback := getGetTableValueCallback(keeper, ctx, appId, owner)
			result := tableValueCallback(querierObjs)
			if len(result) > 0 {
				L.Push(lua.LBool(true))
			} else {
				L.Push(lua.LBool(false))
			}
			return 1
		},
		//add other functions which need to be exported
	}
}

func findByFields( keeper Keeper, ctx sdk.Context, appId uint, owner sdk.AccAddress, tableName string, fields, values []string,) ([]map[string]string, []uint) {

	qo := map[string]string{
		"method": "table",
		"table": tableName,
	}
	querierObjs := []map[string]string{qo}
	for index, field := range fields {
		field := field
		value := values[index]
		qo := map[string]string{
			"method": "where",
			"operator": "==",
			"field": field,
			"value": value,
		}
		querierObjs = append(querierObjs, qo)
	}

	qq := map[string]string{
		"method": "select",
		"fields": "id",
	}
	newQuerierObjs := append(querierObjs, qq)
	res, ids, err := querierSuperHandler(ctx, keeper, appId, newQuerierObjs, owner)
	if err != nil {
		return nil, nil
	}
	return res.Data, ids
}
func getGoExportFilterFunc(ctx sdk.Context, appId uint, keeper Keeper, owner sdk.AccAddress) map[string]lua.LGFunction {
	return map[string]lua.LGFunction{
		"Insert": func(L *lua.LState) int {
			ParamsNum := L.GetTop()
			if ParamsNum >= 2 && ParamsNum%2 == 0 { //Normal inserttab,fields
				tableName := L.ToString(1)
				if strings.HasPrefix(tableName, tablePrefix){
					tableName = strings.TrimPrefix(tableName, tablePrefix)
				}
				sFieldAndValues := L.ToString(2)
				fieldAndValues, err := getFieldValueMap(ctx, appId, keeper, tableName, sFieldAndValues)
				if err != nil {
					L.Push(lua.LNumber(-1))
					L.Push(lua.LString(err.Error()))
					return 2
				}
				if ParamsNum > 2 { //此时表示有外键插入，可以有多个外键插入，格式为foreigntab，foreignid 循环
					for i := 3; i < ParamsNum; i+=2{
						fTableName := L.ToString(i)
						fId := L.ToString(i+1)
						fKey := strings.ToLower(fTableName)
						fieldAndValues[fKey] = fId
					}

				}
				Write := getInsertCallback(keeper, ctx, appId, owner)
				Write(tableName, fieldAndValues)//keeper.Insert(ctx, appId, tableName, fieldAndValues, owner)
				L.Push(lua.LNumber(1))
				L.Push(lua.LString(""))
			} else {
				L.Push(lua.LNumber(-1))
				L.Push(lua.LString("num of param wrong"))
			}
			return 2
		},
		"fieldIn" : func(L *lua.LState) int {
			ParamsNum := L.GetTop()
			if ParamsNum < 2 {
				L.Push(lua.LBool(false))
				return 1
			}

			src := L.ToString(1)
			for i := 2; i <= ParamsNum; i++ {
				dst := L.ToString(i)
				if src == dst{
					L.Push(lua.LBool(true))
					return 1
				}
			}
			L.Push(lua.LBool(false))
			return 1
		},
		"exist" : func(L *lua.LState) int {
			ParamsNum := L.GetTop()
			if ParamsNum < 4 || (ParamsNum - 1)%3 != 0{
				L.Push(lua.LBool(false))
				return 1
			}
			tableName := L.ToString(1)
			qo := map[string]string{
				"method": "table",
				"table": tableName,
			}
			querierObjs := []map[string]string{qo}
			for i := 2; i < ParamsNum; i += 3 {
				field := L.ToString(i)
				op  := L.ToString(i+1)
				value := L.ToString(i+2)
				if op != "==" {
					continue
				}
				qo := map[string]string{
					"method": "where",
					"operator": "==",
					"field": field,
					"value": value,
				}
				querierObjs = append(querierObjs, qo)
			}

			tableValueCallback := getGetTableValueCallback(keeper, ctx, appId, owner)
			result := tableValueCallback(querierObjs)
			if len(result) > 0 {
				L.Push(lua.LBool(true))
			} else {
				L.Push(lua.LBool(false))
			}
			return 1
		},
		"valitatecode" : func(L *lua.LState) int {
			code := L.ToString(1)
			if len(code) != 6 {
				L.Push(lua.LBool(false))
				return 1
			}
			base := "wzAG2dsfbkrEinDKPamBpQ6WtUuHLNceyRVXZ78h3TCJSY5qxjvM14F"
			for _,v := range code {
				index := strings.IndexByte(base,byte(v))
				if index == -1 {
					L.Push(lua.LBool(false))
					return 1
				}
			}
			L.Push(lua.LBool(true))
			return 1
		},
	}
}
func getGoExportQueryFunc(ctx sdk.Context, appId uint, keeper Keeper, addr sdk.AccAddress) map[string]lua.LGFunction {
	return map[string]lua.LGFunction {
		"findRow" : func(L *lua.LState) int {
			tableName := L.ToString(1)
			id := L.ToInt(2)
			fields, err := keeper.Find(ctx, appId, tableName, uint(id), addr)
			if err != nil {
				setLuaFuncRes(L, createLuaTable(false), lua.LString(""))
				return 2
			}
			ud := setUserData(ctx, appId, keeper, addr, tableName, []map[string]string{fields}, L)
			L.Push(ud)
			return 1
		},
		"findRows" : func(L *lua.LState) int {
			tableName := L.ToString(1)
			ids := L.ToTable(2)
			res := make([]types.RowFields, 0)
			ids.ForEach(func(lId lua.LValue, val lua.LValue) {
				id := val.(lua.LNumber)
				fields, err := keeper.Find(ctx, appId, tableName, uint(id), addr)
				if err != nil {
					return
				}
				res = append(res, fields)
			})
			sliceFieldTab := createLuaTable(res)
			setLuaFuncRes(L, sliceFieldTab, lua.LString(""))
			return 2
		},
		"findIdsBy" : func(L *lua.LState) int {
			tableName := L.ToString(1)
			fieldName := L.ToString(2)
			value := L.ToString(3)
			ids := keeper.FindBy(ctx, appId, tableName, fieldName, []string{value}, addr)
			idsTable := createLuaTable(ids)
			setLuaFuncRes(L, idsTable,lua.LString(""))
			return 2
		},
		"findAllIds" : func(L *lua.LState) int {
			tableName := L.ToString(1)
			ids := keeper.FindAll(ctx, appId, tableName, addr)
			idsTable := createLuaTable(ids)
			setLuaFuncRes(L, idsTable,lua.LString(""))
			return 2
		},
		"specFindRowsBy" : func(L *lua.LState) int {
			res := make([]types.RowFields, 0)
			querierObjJson := L.ToString(1)
			var querierObjs [](map[string]string)
			if err := json.Unmarshal([]byte(querierObjJson), &querierObjs); err != nil {
				sliceFieldTab := createLuaTable(res)
				setLuaFuncRes(L, sliceFieldTab, lua.LString(err.Error()))
				return 2
			}
			result, _, err := specQuerierSuperHandler(ctx, keeper, appId, querierObjs, addr, false)
			if err != nil {
				sliceFieldTab := createLuaTable(res)
				setLuaFuncRes(L, sliceFieldTab, lua.LString(err.Error()))
				return 2
			}
			result = checkResult(ctx, addr, keeper, appId, querierObjs,result)
			for _,val := range result {
				res = append(res, val)
			}
			sliceFieldTab := createLuaTable(res)
			setLuaFuncRes(L, sliceFieldTab, lua.LString(""))
			return 2
		},
		"queryOracle" : func(L *lua.LState) int {
			res := make([]types.RowFields, 0)
			querierObjJson := L.ToString(1)
			appId, _ := keeper.GetDatabaseId(ctx, "0000000001")
			var querierObjs [](map[string]string)
			if err := json.Unmarshal([]byte(querierObjJson), &querierObjs); err != nil {
				sliceFieldTab := createLuaTable(res)
				setLuaFuncRes(L, sliceFieldTab, lua.LString(err.Error()))
				return 2
			}
			datas, _, err := specQuerierSuperHandler(ctx, keeper, appId, querierObjs, addr, true)
			if err != nil {
				sliceFieldTab := createLuaTable(res)
				setLuaFuncRes(L, sliceFieldTab, lua.LString(err.Error()))
				return 2
			}

			for _, data := range datas{
				if data["owner"] == addr.String() {
					res = append(res, data)
				}
			}
			sliceFieldTab := createLuaTable(res)
			setLuaFuncRes(L, sliceFieldTab, lua.LString(""))
			return 2
		},
	}
}
func getFieldValueMap(ctx sdk.Context, appId uint, keep Keeper, tableName string, s string) (types.RowFields, error) {
	tbFields, err := keep.getTableFields(ctx, appId, tableName)
	if err != nil {
		return nil, err
	}

	//use json format
	values := make([]string,0)
	err = json.Unmarshal([]byte(s),&values)
	if err != nil {
		return nil, errors.New("func getFieldValueMap err, unmarshal params err")
	}
	rowFields := make(types.RowFields)
	/*
	表的前三个字段固定 由系统创建 id create_by create_at,如果有外键，外键为第一个字段，否则会添加数据出错
	*/
	tbFields = tbFields[3:]
	for i := 0; i < len(tbFields); i++ {
		if i < len(values) {
			field := tbFields[i]
			rowFields[field] = values[i]
		}
	}

	return rowFields, nil
}

func createLuaTable(src interface{}) *lua.LTable{
	L := lua.NewState()
	defer L.Close()
	tb := L.NewTable()
	switch src.(type) {
	case types.RowFields:
		ns := src.(types.RowFields)
		for key, val := range ns {
			tb.RawSetString(key,lua.LString(val))
		}
	case []uint:
		ns := src.([]uint)
		for id, val := range ns {
			tb.RawSetInt(id+1,lua.LNumber(val))
		}
	case []types.RowFields:
		ns := src.([]types.RowFields)
		for index, m := range ns {
			subTb := L.NewTable()
			for key, val := range m {
				subTb.RawSetString(key,lua.LString(val))
			}
			tb.RawSetInt(index+1,subTb)
		}
	case []map[string]string:
		ns := src.([]map[string]string)
		for index, m := range ns {
			subTb := L.NewTable()
			for key, val := range m {
				subTb.RawSetString(key,lua.LString(val))
			}
			tb.RawSetInt(index+1,subTb)
		}
	default:
		return tb
	}
	return tb
}

func convertLuaTableToGo(table *lua.LTable) interface{}{
	resMap := make(map[string]string)
	resSlice := make([]uint, 0)
	resSliceMap := make([]map[string]string, 0)
	table.ForEach(func(key lua.LValue, val lua.LValue) {
		nKey,ok  := key.(lua.LString)
		if ok { //it means the format of this table is map[string]string
			nVal := val.(lua.LString)
			resMap[nKey.String()] = nVal.String()
		} else  {
			nVal, ok  := val.(lua.LNumber)
			if ok { //map[int]int
				resSlice = append(resSlice, uint(nVal))
			} else { //[]map[string]string
				temp := make(map[string]string)
				tVal := val.(*lua.LTable)
				tVal.ForEach(func(k lua.LValue, v lua.LValue) {
					temp[k.String()] = v.String()
				})
				resSliceMap = append(resSliceMap, temp)
			}

		}
	})
	if len(resMap) > 0 {
		return resMap
	} else if len(resSlice) > 0{
		return resSlice
	} else {
		return resSliceMap
	}
}

func setLuaFuncRes(L *lua.LState, value, err lua.LValue){
	L.Push(value)
	L.Push(err)
}

/////////////////////////////////
//                             //
//  helper func                //
//                             //
/////////////////////////////////

func splitParams(src []string) [][]string {
	var res [][]string
	for i := 0; i < len(src); i++ {
		if strings.HasPrefix(src[i],tablePrefix) {
			temp := make([]string, 0)
			temp = append(temp, strings.TrimPrefix(src[i],tablePrefix))
			i++
			for ; i < len(src); i++{
				if strings.HasPrefix(src[i],tablePrefix) {
					i--
					break
				}
				temp = append(temp, strings.TrimPrefix(src[i],foreignPrefix))
			}
			res = append(res, temp)
		}
	}
	return res
}

func setUserData(ctx sdk.Context, appId uint, keeper Keeper, addr sdk.AccAddress, tableName string, values []map[string]string, L *lua.LState) *lua.LUserData{
	associations := keeper.GetTableAssociations(ctx, appId, tableName)

	tableObj := registerAssociation(ctx, appId, keeper, addr, L,tableName, associations, values)
	ud := L.NewUserData()
	ud.Value = tableObj
	L.SetMetatable(ud, L.GetTypeMetatable(LuaTableTypeName))
	return ud
}

func registerAssociation(ctx sdk.Context, appId uint, keeper Keeper, addr sdk.AccAddress, L *lua.LState, tableName string , associations []types.Association, value []map[string]string)  *TableObj {

	table := &TableObj{
		TableName: tableName,
		Value: value,
	}
	tableObjMethods := GetTableObjMethods()

	if associations == nil {
		return table
	}

	for _, association := range associations {

		TableName := association.AssociationTable
		ForeignKey := association.ForeignKey
		MethodName := association.Method
		switch association.AssociationMode {
		case "has_one":
			tableObjMethods[MethodName] = func(state *lua.LState) int {
				hasOne(ctx, appId, keeper, addr, state, TableName, ForeignKey)
				return 1
			}
		case "has_many":
			tableObjMethods[MethodName] = func(state *lua.LState) int {
				hasMany(ctx, appId, keeper, addr, state, TableName, ForeignKey)
				return 1
			}
		case "belongs_to":
			tableObjMethods[MethodName] = func(state *lua.LState) int {
				belongsTo(ctx, appId, keeper, addr, state, TableName, ForeignKey)
				return 1
			}
		}
	}

	mt := L.GetTypeMetatable(LuaTableTypeName)
	L.SetField(mt, "__index", L.SetFuncs(L.NewTable(), tableObjMethods))

	return table
}
