package scripting

import (
	"TotalControl/backend/utils"
	log "github.com/sirupsen/logrus"
	lua "github.com/yuin/gopher-lua"
	"os"
)

func luaOsGetenv(L *lua.LState) int {
	key := L.ToString(1)
	if key == "" {
		L.Push(lua.LNil)
		return 1
	}

	value, exists := os.LookupEnv(key)
	if !exists {
		L.Push(lua.LNil)
		return 1
	}

	L.Push(lua.LString(value))
	return 1
}

func luaGetOperatingSystem(L *lua.LState) int {
	// 1 - Windows, 2 - Linux, 3 - MacOS, 0 - Unknown
	switch utils.GetOperatingSystem() {
	case utils.WindowsOS:
		L.Push(lua.LNumber(1))
		return 1
	case utils.LinuxOS:
		L.Push(lua.LNumber(2))
		return 1
	case utils.MacOS:
		L.Push(lua.LNumber(3))
		return 1
	default:
		L.Push(lua.LNumber(0))
		return 1
	}
}

func luaIsWindows(L *lua.LState) int {
	// Check if the current operating system is Windows
	if utils.GetOperatingSystem() == utils.WindowsOS {
		L.Push(lua.LTrue)
	} else {
		L.Push(lua.LFalse)
	}
	return 1
}

func luaIsLinux(L *lua.LState) int {
	// Check if the current operating system is Linux
	if utils.GetOperatingSystem() == utils.LinuxOS {
		L.Push(lua.LTrue)
	} else {
		L.Push(lua.LFalse)
	}
	return 1
}

func luaIsMacOS(L *lua.LState) int {
	// Check if the current operating system is MacOS
	if utils.GetOperatingSystem() == utils.MacOS {
		L.Push(lua.LTrue)
	} else {
		L.Push(lua.LFalse)
	}
	return 1
}

func luaExtendOsTable(l *lua.LState) {
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
	l.SetField(tbl, "isMacOS", l.NewFunction(luaIsMacOS))
	l.SetField(tbl, "is_unknown", lua.LFalse)
	l.SetField(tbl, "is_windows", lua.LFalse)
	l.SetField(tbl, "is_linux", lua.LFalse)
	l.SetField(tbl, "is_macos", lua.LFalse)
	switch utils.GetOperatingSystem() {
	case utils.WindowsOS:
		l.SetField(tbl, "is_windows", lua.LTrue)
	case utils.LinuxOS:
		l.SetField(tbl, "is_linux", lua.LTrue)
	case utils.MacOS:
		l.SetField(tbl, "is_macos", lua.LTrue)
	default:
		l.SetField(tbl, "is_unknown", lua.LTrue)
	}
}
