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
	sdk "github.com/cosmos/cosmos-sdk/types"
	lua "github.com/yuin/gopher-lua"
	"github.com/yzhanginwa/dbchain/x/dbchain/internal/types"
	"strings"
)

func getGoExportFunc(ctx sdk.Context, appId uint, keeper Keeper, owner sdk.AccAddress) map[string]lua.LGFunction {
	return map[string]lua.LGFunction{
		"LCT": func(L *lua.LState) int {
			tableName := L.ToString(1)
			sFieldName := L.ToString(2)
			fieldNames := strings.Split(sFieldName, ",")
			keeper.CreateTable(ctx, appId, owner, tableName, fieldNames)
			L.Push(lua.LString(""))
			return 1
		},

		"Insert": func(L *lua.LState) int {
			ParamsNum := L.GetTop()
			if ParamsNum >= 2 && ParamsNum%2 == 0 { //Normal inserttab,fields
				tableName := L.ToString(1)
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
		"MultInsert": func(L *lua.LState) int{ //往同一张表插入多条数据
			ParamsNum := L.GetTop()
			if ParamsNum < 2 {
				L.Push(lua.LNumber(-1))
				L.Push(lua.LString("num of param wrong"))
				return 2
			}
			tableName := L.ToString(1)
			count := 0
			for i := 2 ; i <= ParamsNum; i++{
				sFieldAndValues := L.ToString(i)
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
			for i := 2 ; i <= ParamsNum; i++ {
				id := L.ToInt(i)
				keeper.Freeze(ctx, appId, tableName, uint(id), owner)
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
					"method": "equal",
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

func getFieldValueMap(ctx sdk.Context, appId uint, keep Keeper, tableName string, s string) (types.RowFields, error) {
	tbFields, err := keep.getTableFields(ctx, appId, tableName)
	if err != nil {
		return nil, err
	}

	values := strings.Split(s, ",")
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
