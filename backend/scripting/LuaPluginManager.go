package scripting

import (
	log "github.com/sirupsen/logrus"
	lua "github.com/yuin/gopher-lua"
)

func RegisterPluginManagerObject(L *lua.LState) {
	table := L.NewTable()
	table.RawSetString("isLoaded", L.NewFunction(LuaIsPluginLoaded))
	table.RawSetString("isEnabled", L.NewFunction(LuaIsPluginEnabled))
	L.SetGlobal("plugin", table)
}

func LuaIsPluginLoaded(L *lua.LState) int {
	if L.GetTop() != 1 {
		L.Push(lua.LBool(false))
		return 1
	}
	pluginName := L.ToString(1)
	log.Debugf("Checking if plugin '%s' is loaded", pluginName)
	L.Push(lua.LBool(false))
	return 1
}

func LuaIsPluginEnabled(L *lua.LState) int {
	if L.GetTop() != 1 {
		L.Push(lua.LBool(false))
		return 1
	}
	pluginName := L.ToString(1)
	log.Debugf("Checking if plugin '%s' is enabled", pluginName)
	L.Push(lua.LBool(false))
	return 1
}
