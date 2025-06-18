package scripting

import (
	"TotalControl/backend/utils"
	lua "github.com/yuin/gopher-lua"
)

func LuaRegisterSettingsObject(L *lua.LState) {
	settingsTable := L.NewTable()
	settingsTable.RawSetString("get", L.NewFunction(luaSettingsGet))
	settingsTable.RawSetString("set", L.NewFunction(luaSettingsSet))
	L.SetGlobal("settings", settingsTable)
}

func luaSettingsGet(L *lua.LState) int {
	luaPlugin := GetLuaPlugin(L)
	if luaPlugin == nil {
		L.Push(lua.LFalse)
		L.Push(lua.LString("Plugin not found"))
		return 2
	}

	if L.GetTop() != 1 {
		L.Push(lua.LFalse)
		L.Push(lua.LString("Expected 1 argument: key string"))
		return 2
	}

	key := L.ToString(1)
	value, err := luaPlugin.Get(key)
	if err != nil {
		L.Push(lua.LFalse)
		L.Push(lua.LString(err.Error()))
		return 2
	}

	L.Push(lua.LTrue)
	L.Push(utils.ToLuaValue(L, value))
	return 2
}

func luaSettingsSet(L *lua.LState) int {
	luaPlugin := GetLuaPlugin(L)
	if luaPlugin == nil {
		L.Push(lua.LFalse)
		L.Push(lua.LString("Plugin not found"))
		return 2
	}

	if L.GetTop() != 2 {
		L.Push(lua.LFalse)
		L.Push(lua.LString("Expected 2 arguments: key string, value any"))
		return 2
	}

	key := L.ToString(1)
	value := utils.FromLuaValue(L, L.Get(2))
	isValid, err := luaPlugin.IsValid(key, value)
	if !isValid {
		if err != nil {
			L.Push(lua.LFalse)
			L.Push(lua.LString(err.Error()))
			return 2
		}
		L.Push(lua.LFalse)
		L.Push(lua.LString("Invalid value type for setting: " + key))
		return 2
	}

	err = luaPlugin.Set(key, value)
	if err != nil {
		L.Push(lua.LFalse)
		L.Push(lua.LString(err.Error()))
		return 2
	}

	L.Push(lua.LTrue)
	return 1
}
