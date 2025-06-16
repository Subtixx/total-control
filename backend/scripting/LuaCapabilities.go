package scripting

import (
	lua "github.com/yuin/gopher-lua"
)

func LuaRegisterCapabilitiesObject(L *lua.LState) int {
	capabilities := L.NewTable()
	capabilities.RawSetString("canAccessFileSystem", L.NewFunction(LuaCheckCanAccessFileSystem))
	capabilities.RawSetString("canAccessNetwork", L.NewFunction(LuaCheckCanAccessNetwork))
	capabilities.RawSetString("can", L.NewFunction(LuaCheckCan))

	L.SetGlobal("capabilities", capabilities)
	return 1
}

func LuaCheckCanAccessFileSystem(L *lua.LState) int {
	luaPlugin := GetLuaPlugin(L)
	if luaPlugin == nil {
		L.Push(lua.LTrue)
		return 1
	}

	if luaPlugin.CanAccessFileSystem() {
		L.Push(lua.LTrue)
		return 1
	}

	L.Push(lua.LFalse)
	return 1
}

func LuaCheckCanAccessNetwork(L *lua.LState) int {
	luaPlugin := GetLuaPlugin(L)
	if luaPlugin == nil {
		L.Push(lua.LTrue)
		return 1
	}

	if luaPlugin.CanAccessNetwork() {
		L.Push(lua.LTrue)
		return 1
	}

	L.Push(lua.LFalse)
	return 1
}

func LuaCheckCan(L *lua.LState) int {
	capability := L.ToString(1)
	if CheckCan(L, capability) {
		L.Push(lua.LTrue)
	} else {
		L.Push(lua.LFalse)
	}
	return 1
}

func CheckCan(L *lua.LState, capability string) bool {
	luaPlugin := GetLuaPlugin(L)
	if luaPlugin != nil {
		if !luaPlugin.Can(capability) {
			return false
		}
	}

	return true
}
