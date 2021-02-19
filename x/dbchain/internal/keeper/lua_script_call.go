package keeper

import (
	"errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	lua "github.com/yuin/gopher-lua"
)
//Different databases correspond to different handles
var luaHandles = make(map[uint]*lua.LState)

func getAppLuaHandle(appId uint)  *lua.LState {
	l , ok := luaHandles[appId]
	if !ok {
		var luaHandle = lua.NewState(lua.Options{
			SkipOpenLibs: true, //set SkipOpenLibs true to prevent lua open libs,because this libs can call os function and operate files
			RegistrySize: 256, //Save function return value， default value is 256*20, it is too large
		})
		luaHandles[appId] = luaHandle
		return luaHandle
	}
	return l
}
func getRegisterLuaFunc(ctx sdk.Context, keeper Keeper, appId uint, funcName string) lua.LValue{
	l := luaHandles[appId]

	Fn := l.GetGlobal(funcName)
	if Fn.String() == "nil" || Fn == nil {
		err := registerLuaFunc(ctx, appId, funcName, keeper, l)
		if err != nil {
			return nil
		}
		Fn = l.GetGlobal(funcName)
	}
	return Fn
}

func registerLuaFunc(ctx sdk.Context, appId uint, luaFunc string, keeper Keeper, l *lua.LState) error {
	functionInfo := keeper.GetFunctionInfo(ctx, appId, luaFunc, 0)
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
	//TODO new format script needs to be restructure
	return funcBody
}

//when add new func or modify func, luaHandle needs to be delete
func voidLuaHandle(appId uint) {
	l, ok := luaHandles[appId]
	if ok {
		l.Close()
	}
	delete(luaHandles, appId)
}

func callLuaScriptFunc(ctx sdk.Context, appId uint, owner sdk.AccAddress, keeper Keeper, funcName string, params []lua.LValue) error{
	l := getAppLuaHandle(appId)
	//point : get go function
	goExportFunc := getGoExportFunc(ctx, appId, keeper, owner)
	//register go function
	for name, fn := range goExportFunc{
		l.SetGlobal(name, l.NewFunction(fn))
	}
	//call lua script
	fn := getRegisterLuaFunc(ctx, keeper, appId, funcName)
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


