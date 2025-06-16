package scripting

import (
	"TotalControl/backend/plugins"
	"TotalControl/backend/utils"
	lua "github.com/yuin/gopher-lua"
)

func LuaRegisterZipObject(l *lua.LState) {
	table := l.NewTable()
	l.SetGlobal("zip", table)
	l.SetField(table, "readFile", l.NewFunction(LuaReadFileFromZip))
	l.SetField(table, "readFiles", l.NewFunction(LuaReadFilesFromZip))
}

func LuaReadFilesFromZip(L *lua.LState) int {
	luaPlugin := GetLuaPlugin(L)
	if luaPlugin == nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(ErrPluginNotFound))
		return 2
	}

	if !luaPlugin.CanAccessFileSystem() {
		L.Push(lua.LNil)
		L.Push(lua.LString(plugins.ErrFileSystemAccessDenied))
		return 2
	}

	zipPath := L.ToString(1)
	if zipPath == "" {
		L.Push(lua.LNil)
		return 1
	}

	files, err := utils.ReadFilesFromZip(zipPath)
	if err != nil {
		L.RaiseError("Failed to read files from zip: %s", err.Error())
		return 0
	}

	resultTable := L.NewTable()
	for name, content := range files {
		resultTable.RawSetString(name, lua.LString(content))
	}
	L.Push(resultTable)
	return 1
}

func LuaReadFileFromZip(L *lua.LState) int {
	zipPath := L.ToString(1)
	fileName := L.ToString(2)
	useRegEx := true
	if L.GetTop() > 2 {
		if L.Get(3).Type() != lua.LTBool {
			L.RaiseError("Expected boolean for useRegEx, got %s", L.Get(3).Type().String())
			return 0
		}
		useRegEx = L.ToBool(3)
	}

	if zipPath == "" || fileName == "" {
		L.Push(lua.LNil)
		return 1
	}

	data, err := utils.ReadFileFromZip(zipPath, fileName, useRegEx)
	if err != nil {
		L.RaiseError("Failed to read file from zip: %s", err.Error())
		return 0
	}

	L.Push(lua.LString(data))
	return 1
}
