package scripting

import (
	"TotalControl/backend/utils"
	lua "github.com/yuin/gopher-lua"
)

func RegisterPluginManagerObject(L *lua.LState) {
	table := L.NewTable()
	table.RawSetString("isLoaded", L.NewFunction(LuaIsPluginLoaded))
	table.RawSetString("isEnabled", L.NewFunction(LuaIsPluginEnabled))
	L.SetGlobal("plugin", table)
}

func LuaIsPluginLoaded(L *lua.LState) int {
	luaPlugin := GetLuaPlugin(L)
	if luaPlugin == nil {
		L.RaiseError(ErrPluginNotFound)
		return 0
	}

	pluginName := utils.GetString(L, 1)
	if pluginName == "" {
		L.Push(lua.LBool(false))
		return 1
	}

	L.Push(lua.LBool(false))
	return 1
}

func LuaIsPluginEnabled(L *lua.LState) int {
	luaPlugin := GetLuaPlugin(L)
	if luaPlugin == nil {
		L.RaiseError(ErrPluginNotFound)
		return 0
	}

	pluginName := utils.GetString(L, 1)
	if pluginName == "" {
		L.Push(lua.LBool(false))
		return 1
	}

	L.Push(lua.LBool(false))
	return 1
}
