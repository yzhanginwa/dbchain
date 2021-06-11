package keeper

import (
	"encoding/json"
	sdk "github.com/cosmos/cosmos-sdk/types"
	lua "github.com/yuin/gopher-lua"
	"strconv"
)

type TableObj struct {
	TableName    string            `json:"table_name"`
	Value        []map[string]string `json:"value"`
}

type ExistObj struct {
	TableName  string
	Fields      map[string]string
	ctx        sdk.Context
	appId      uint
	keeper     Keeper
	addr       sdk.AccAddress
}


const LuaTableTypeName = "tableObj"
const LuaTableTypeExist = "existObj"

func registerTableType(L *lua.LState, ctx sdk.Context, appId uint, keeper Keeper, addr sdk.AccAddress) {
	mt := L.NewTypeMetatable(LuaTableTypeName)
	L.SetGlobal("tableObj", mt)
	//static attributes
	L.SetField(mt, "newObjFromJsonString", L.NewFunction(newTableObjFromJsonString))
	//methods
	L.SetField(mt, "__index", L.SetFuncs(L.NewTable(), GetTableObjMethods()))
	//exist obj
	exist := L.NewTypeMetatable(LuaTableTypeExist)
	L.SetGlobal("exist", exist)
	L.SetField(exist, "table", L.NewFunction(func(state *lua.LState) int {
		tableName := L.CheckString(1)
		table := &ExistObj{
			TableName: tableName,
			Fields: make(map[string]string, 0),
			ctx : ctx,
			appId: appId,
			keeper: keeper,
			addr: addr,
		}
		ud := L.NewUserData()
		ud.Value = table
		L.SetMetatable(ud, L.GetTypeMetatable(LuaTableTypeExist))
		L.Push(ud)
		return 1
	}))
	L.SetField(exist, "__index", L.SetFuncs(L.NewTable(), map[string]lua.LGFunction {"where": where}))

}

func newTableObjFromJsonString(L *lua.LState) int {
	js := L.CheckString(1)
	table := &TableObj{}
	err := json.Unmarshal([]byte(js), table)
	if err != nil {
		return 0
	}
	ud := L.NewUserData()
	ud.Value = table
	L.SetMetatable(ud, L.GetTypeMetatable(LuaTableTypeName))
	L.Push(ud)
	return 1
}

func CheckTable(L *lua.LState) *TableObj {
	ud := L.CheckUserData(1)
	if v, ok := ud.Value.(*TableObj); ok {
		return v
	}
	L.ArgError(1, "TableObj expected")
	return nil
}

func CheckExist(L *lua.LState) *ExistObj {
	ud := L.CheckUserData(1)
	if v, ok := ud.Value.(*ExistObj); ok {
		return v
	}
	L.ArgError(1, "TableObj expected")
	return nil
}

func GetTableObjMethods() map[string]lua.LGFunction {
	var tableObjMethods = map[string]lua.LGFunction {
		"data":           data,
		"index":           index,
		"first":          first,
		"last":           last,
	}
	return tableObjMethods
}

func data(L *lua.LState) int {
	p := CheckTable(L)
	sliceFieldTab := createLuaTable(p.Value)
	L.Push(sliceFieldTab)
	return 1
}

func index(L *lua.LState) int {
	index := L.ToInt(2)
	p := CheckTable(L)
	value := make([]map[string]string, 0)
	if len(p.Value) > index {
		value = append(value, p.Value[index])
	}
	obj := &TableObj{
		TableName: p.TableName,
		Value: value,
	}
	ud := L.NewUserData()
	ud.Value = obj
	L.SetMetatable(ud, L.GetTypeMetatable(LuaTableTypeName))
	L.Push(ud)
	return 1
}

func first(L *lua.LState) int {
	p := CheckTable(L)
	value := make([]map[string]string, 0)
	if len(p.Value) > 0 {
		value = append(value, p.Value[0])
	}
	obj := &TableObj{
		TableName: p.TableName,
		Value: value,
	}
	ud := L.NewUserData()
	ud.Value = obj
	L.SetMetatable(ud, L.GetTypeMetatable(LuaTableTypeName))
	L.Push(ud)
	return 1
}

func last(L *lua.LState) int {
	p := CheckTable(L)
	value := make([]map[string]string, 0)
	if len(p.Value) > 0 {
		value = append(value, p.Value[len(p.Value)-1])
	}
	obj := &TableObj{
		TableName: p.TableName,
		Value: value,
	}
	ud := L.NewUserData()
	ud.Value = obj
	L.SetMetatable(ud, L.GetTypeMetatable(LuaTableTypeName))
	L.Push(ud)
	return 1
}

func belongsTo(ctx sdk.Context, appId uint, keeper Keeper, addr sdk.AccAddress, L *lua.LState, TableName, ForeignKey string)  {
	index := 0
	if L.GetTop() > 1 {
		index = L.ToInt(2)
	}
	tableObj := CheckTable(L)
	if len(tableObj.Value) == 0 {
		registerUserData(ctx, appId, keeper, addr, tableObj, L, TableName, []map[string]string{})
		return
	}

	tableName := TableName
	field     := ForeignKey
	sId     := tableObj.Value[index][field]
	id , err := strconv.Atoi(sId)
	if err != nil {
		registerUserData(ctx, appId, keeper, addr, tableObj, L, tableName, []map[string]string{})
		return
	}

	fields , err := keeper.Find(ctx, appId, tableName, uint(id), addr)
	if err != nil {
		registerUserData(ctx, appId, keeper, addr, tableObj, L, tableName, []map[string]string{})
		return
	}

	registerUserData(ctx, appId, keeper, addr, tableObj, L, tableName, []map[string]string{fields})
	return
}

func hasMany(ctx sdk.Context, appId uint, keeper Keeper, addr sdk.AccAddress, L *lua.LState, TableName, ForeignKey string)  {
	index := 0
	if L.GetTop() > 1 {
		index = L.ToInt(2)
	}

	tableObj := CheckTable(L)
	if len(tableObj.Value) == 0 {
		registerUserData(ctx, appId, keeper, addr, tableObj, L, TableName, []map[string]string{})
		return
	}

	tableName := TableName
	field     := ForeignKey
	value     := tableObj.Value[index]["id"]
	ids := keeper.FindBy(ctx, appId, tableName, field, []string{value}, addr)
	if len(ids) == 0 {
		registerUserData(ctx, appId, keeper, addr, tableObj, L, tableName, []map[string]string{})
		return
	}

	values := make([]map[string]string,0)
	for _, id := range ids {
		fields , err := keeper.Find(ctx, appId, tableName, id, addr)
		if err != nil {
			continue
		}
		values = append(values, fields)
	}

	registerUserData(ctx, appId, keeper, addr, tableObj, L, tableName, values)
	return
}

func hasOne(ctx sdk.Context, appId uint, keeper Keeper, addr sdk.AccAddress, L *lua.LState, TableName, ForeignKey string) {
	index := 0
	if L.GetTop() > 1 {
		index = L.ToInt(2)
	}

	tableObj := CheckTable(L)
	if len(tableObj.Value) == 0 {
		registerUserData(ctx, appId, keeper, addr, tableObj, L, TableName, []map[string]string{})
		return
	}
	tableName := TableName
	field     := ForeignKey
	value     := tableObj.Value[index]["id"]
	ids := keeper.FindBy(ctx, appId, tableName, field, []string{value}, addr)
	if len(ids) == 0 {
		registerUserData(ctx, appId, keeper, addr, tableObj, L, tableName, []map[string]string{})
		return
	}

	fields , err := keeper.Find(ctx, appId, tableName, ids[0], addr)
	if err != nil {
		registerUserData(ctx, appId, keeper, addr, tableObj, L, tableName, []map[string]string{})
		return
	}

	registerUserData(ctx, appId, keeper, addr, tableObj, L, tableName, []map[string]string{fields})
	return
}

func where(L *lua.LState) int {
	existObj := CheckExist(L)
	key := L.ToString(2)
	val := L.ToString(3)
	existObj.Fields[key] = val

	tableName := existObj.TableName
	qo := map[string]string{
		"method": "table",
		"table": tableName,
	}
	querierObjs := []map[string]string{qo}
	for key, val := range existObj.Fields {
		qo := map[string]string{
			"method": "where",
			"operator": "==",
			"field": key,
			"value": val,
		}
		querierObjs = append(querierObjs, qo)
	}

	tableValueCallback := getGetTableValueCallback(existObj.keeper, existObj.ctx, existObj.appId, existObj.addr)
	result := tableValueCallback(querierObjs)
	if len(result) > 0 {
		ud := L.NewUserData()
		ud.Value = existObj
		L.SetMetatable(ud, L.GetTypeMetatable(LuaTableTypeExist))
		L.Push(ud)
	} else {
		L.Push(lua.LNil)
	}
	return 1
}

/////////////////////////////////
//                             //
//  helper func                //
//                             //
/////////////////////////////////

func registerUserData(ctx sdk.Context, appId uint, keeper Keeper, addr sdk.AccAddress, tableObj *TableObj, L *lua.LState, tableName string, values []map[string]string)  {

	associations := keeper.GetTableAssociations(ctx, appId, tableName)
	newTableObj := registerAssociation(ctx, appId, keeper, addr, L, tableName, associations, values)

	ud := L.NewUserData()
	ud.Value = newTableObj
	L.SetMetatable(ud, L.GetTypeMetatable(LuaTableTypeName))
	L.Push(ud)
}