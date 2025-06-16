package scripting

import (
	lua "github.com/yuin/gopher-lua"
)

func LuaRegisterCapabilitiesObject(L *lua.LState) int {
	luaPlugin := GetLuaPlugin(L)
	if luaPlugin == nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(ErrPluginNotFound))
		return 2
	}

	capabilities := L.NewTable()
	capabilities.RawSetString("canAccessFileSystem", lua.LBool(luaPlugin.CanAccessFileSystem()))
	capabilities.RawSetString("canAccessNetwork", lua.LBool(luaPlugin.CanAccessNetwork()))

	L.Push(capabilities)
	return 1
}

func LuaCheckCan(L *lua.LState, capability string) bool {
	luaPlugin := GetLuaPlugin(L)
	if luaPlugin != nil {
		if !luaPlugin.Can(capability) {
			return false
		}
	}

	return true
}
