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
	"strings"
)

func getGoExportFunc(ctx sdk.Context, appId uint, keeper Keeper, owner sdk.AccAddress) map[string]lua.LGFunction {
	return map[string]lua.LGFunction{
		"LCT": func(L *lua.LState) int {
			tableName := L.ToString(1)
			sFieldName := L.ToString(2)
			fieldNames := strings.Split(sFieldName,",")
			keeper.CreateTable(ctx, appId, owner, tableName, fieldNames)
			return 1
		},
		//add other functions which need to be exported
	}
}
