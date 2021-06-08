package keeper

import (
	"errors"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	lua "github.com/yuin/gopher-lua"
	"github.com/yzhanginwa/dbchain/x/dbchain/internal/super_script"
	"strings"
)
//Different databases correspond to different handles
var luaFuncHandles = make(map[uint]*lua.LState)
var luaQuerierHandles = make(map[uint]*lua.LState)

const (
	FuncHandleType = iota
	QueryHandleType
)

func getAppLuaHandle(appId uint, handleType int)  *lua.LState {
	var luaHandles = getLuaHandles(handleType)

	l , ok := luaHandles[appId]
	if !ok {
		var luaHandle = lua.NewState(lua.Options{
			SkipOpenLibs: true, //set SkipOpenLibs true to prevent lua open libs,because this libs can call os function and operate files
			RegistrySize: 256, //Save function return value， default value is 256*20, it is too large
		})
		openBase(luaHandle)
		luaHandles[appId] = luaHandle
		if handleType == QueryHandleType {
			registerTableType(luaHandle)
		}
		return luaHandle
	}
	return l
}
func getRegisterLuaFunc(ctx sdk.Context, keeper Keeper, appId uint, funcName string, handleType int) lua.LValue{
	var luaHandles = getLuaHandles(handleType)

	l := luaHandles[appId]

	Fn := l.GetGlobal(funcName)
	if Fn.String() == "nil" || Fn == nil {
		err := registerLuaFunc(ctx, appId, funcName, keeper, l, handleType)
		if err != nil {
			return nil
		}
		Fn = l.GetGlobal(funcName)
	}
	return Fn
}

func registerLuaFunc(ctx sdk.Context, appId uint, luaFunc string, keeper Keeper, l *lua.LState, handleType int) error {
	functionInfo := keeper.GetFunctionInfo(ctx, appId, luaFunc, handleType)
	if functionInfo.Body == "" {
		return errors.New("this func is unRegister")
	}
	luaScript := restructureLuaScript(functionInfo.Body)
	err := l.DoString(luaScript)
	if err != nil {
		return err
	}
	return nil
}

func restructureLuaScript(funcBody string) string {
	p := super_script.NewPreprocessor(strings.NewReader(funcBody))
	p.Process()
	newScript := p.Reconstruct()
	return newScript
}

//when add new func or modify func, luaHandle needs to be delete
func voidLuaHandle(appId uint, handleType int) {
	var luaHandles = getLuaHandles(handleType)

	l, ok := luaHandles[appId]
	if ok {
		l.Close()
	}
	delete(luaHandles, appId)
}

func callLuaScriptFunc(ctx sdk.Context, appId uint, owner sdk.AccAddress, keeper Keeper, funcName string, params []lua.LValue, handleType int) error{
	l := getAppLuaHandle(appId, handleType)
	//point : get go function
	goExportFunc := getGoExportFunc(ctx, appId, keeper, owner)
	//register go function
	for name, fn := range goExportFunc{
		l.SetGlobal(name, l.NewFunction(fn))
	}
	//register new go function
	goExportFunc = getGoExportFuncNew(ctx, appId, keeper, owner)
	for name, fn := range goExportFunc{
		l.SetGlobal(name, l.NewFunction(fn))
	}
	//call lua script
	fn := getRegisterLuaFunc(ctx, keeper, appId, funcName, handleType)
	if fn == nil || fn.String() == "nil" {
		return errors.New("this func has not been registered")
	}
	if err := l.CallByParam(lua.P{
		Fn:      fn,
		NRet:    1,       //脚本返回参数个数
		Protect: true,    //这里设置为ture表示当执行脚本出现panic时，以error返回
	}, params...); err != nil{
		return err
	}
	//handle return
	defer l.Pop(l.GetTop())
	if k := l.GetTop(); k == 1 {
		strErr := l.Get(1).String()
		if strErr != "" && strErr != "nil"{
			return errors.New(strErr)
		}
		return nil
	}

	return errors.New("lua return err")
}


func callLuaScriptQuerierFunc(ctx sdk.Context, appId uint, owner sdk.AccAddress, keeper Keeper, querierName string, params []lua.LValue, handleType int) ([]byte, error){
	l := getAppLuaHandle(appId, handleType)
	//point : get go function
	goExportFunc := getGoExportQueryFunc(ctx, appId, keeper, owner)
	//register go function
	for name, fn := range goExportFunc{
		l.SetGlobal(name, l.NewFunction(fn))
	}
	//call lua script
	fn := getRegisterLuaFunc(ctx, keeper, appId, querierName, handleType)
	if fn == nil || fn.String() == "nil" {
		return nil, errors.New("this func has not been registered")
	}
	if err := l.CallByParam(lua.P{
		Fn:      fn,
		NRet:    1,       //脚本返回参数个数
		Protect: true,    //这里设置为ture表示当执行脚本出现panic时，以error返回
	}, params...); err != nil{
		return nil, err
	}
	//handle return
	defer l.Pop(l.GetTop())
	lRes := l.Get(1)
	if tableRes ,ok := lRes.(*lua.LTable); ok {
		res := convertLuaTableToGo(tableRes)
		bz, err := codec.MarshalJSONIndent(keeper.cdc, res)
		if err != nil {
			return nil, errors.New("could not marshal result to JSON")
		}
		return bz, nil
	}

	return nil, errors.New("lua return err")
}


func getLuaHandles(handleType int)  map[uint]*lua.LState {
	var luaHandles  map[uint]*lua.LState
	if handleType == FuncHandleType {
		luaHandles = luaFuncHandles
	} else if handleType == QueryHandleType {
		luaHandles = luaQuerierHandles
	}
	return luaHandles
}

func openBase(L *lua.LState) {
	global := L.Get(lua.GlobalsIndex).(*lua.LTable)
	global.RawSetString("pairs", L.NewClosure(basePairs, L.NewFunction(pairsaux)))
}

func basePairs(L *lua.LState) int {
	tb := L.CheckTable(1)
	L.Push(L.Get(lua.UpvalueIndex(1)))
	L.Push(tb)
	L.Push(lua.LNil)
	return 3
}

func pairsaux(L *lua.LState) int {
	tb := L.CheckTable(1)
	key, value := tb.Next(L.Get(2))
	if key == lua.LNil {
		return 0
	} else {
		L.Pop(1)
		L.Push(key)
		L.Push(key)
		L.Push(value)
		return 2
	}
}