package scripting

import (
	"TotalControl/backend/utils"
	log "github.com/sirupsen/logrus"
	lua "github.com/yuin/gopher-lua"
)

func luaGetOperatingSystem(L *lua.LState) int {
	// 1 - Windows, 2 - Linux, 3 - Mac, 0 - Unknown
	switch utils.GetOperatingSystem() {
	case utils.OperatingSystemWindows:
		L.Push(lua.LNumber(1))
		return 1
	case utils.OperatingSystemLinux:
		L.Push(lua.LNumber(2))
		return 1
	case utils.OperatingSystemMac:
		L.Push(lua.LNumber(3))
		return 1
	default:
		L.Push(lua.LNumber(0))
		return 1
	}
}

func luaIsWindows(L *lua.LState) int {
	// Check if the current operating system is Windows
	if utils.GetOperatingSystem() == utils.OperatingSystemWindows {
		L.Push(lua.LTrue)
	} else {
		L.Push(lua.LFalse)
	}
	return 1
}

func luaIsLinux(L *lua.LState) int {
	// Check if the current operating system is Linux
	if utils.GetOperatingSystem() == utils.OperatingSystemLinux {
		L.Push(lua.LTrue)
	} else {
		L.Push(lua.LFalse)
	}
	return 1
}

func luaIsMac(L *lua.LState) int {
	// Check if the current operating system is Mac
	if utils.GetOperatingSystem() == utils.OperatingSystemMac {
		L.Push(lua.LTrue)
	} else {
		L.Push(lua.LFalse)
	}
	return 1
}

func LuaExtendOsTable(l *lua.LState) {
	osTable := l.GetGlobal("os")
	tbl, ok := osTable.(*lua.LTable)
	if !ok {
		log.Warnf("os table not found, creating a new one")
		tbl = l.NewTable()
		l.SetGlobal("os", tbl)
	}
	l.SetField(tbl, "getOperatingSystem", l.NewFunction(luaGetOperatingSystem))
	l.SetField(tbl, "isWindows", l.NewFunction(luaIsWindows))
	l.SetField(tbl, "isLinux", l.NewFunction(luaIsLinux))
	l.SetField(tbl, "isMac", l.NewFunction(luaIsMac))
	l.SetField(tbl, "is_unknown", lua.LFalse)
	l.SetField(tbl, "is_windows", lua.LFalse)
	l.SetField(tbl, "is_linux", lua.LFalse)
	l.SetField(tbl, "is_mac", lua.LFalse)
	switch utils.GetOperatingSystem() {
	case utils.OperatingSystemWindows:
		l.SetField(tbl, "is_windows", lua.LTrue)
	case utils.OperatingSystemLinux:
		l.SetField(tbl, "is_linux", lua.LTrue)
	case utils.OperatingSystemMac:
		l.SetField(tbl, "is_mac", lua.LTrue)
	default:
		l.SetField(tbl, "is_unknown", lua.LTrue)
	}
}
