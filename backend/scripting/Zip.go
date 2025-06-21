package scripting

import (
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
	zipPath := utils.GetString(L, 1)
	if zipPath == "" {
		L.RaiseError("zip.readFiles: zip path cannot be empty")
		return 0
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
	luaPlugin := GetLuaPlugin(L)
	if luaPlugin == nil {
		L.RaiseError(ErrPluginNotFound)
		return 0
	}
	zipPath := utils.GetString(L, 1)
	if zipPath == "" {
		L.RaiseError("zip.readFile: zip path cannot be empty")
		return 0
	}

	fileName := utils.GetString(L, 2)
	if fileName == "" {
		L.RaiseError("zip.readFile: file name cannot be empty")
		return 0
	}
	useRegEx := true
	if L.GetTop() > 2 {
		useRegEx = utils.GetBool(L, 3)
	}

	data, err := utils.ReadFileFromZip(zipPath, fileName, useRegEx)
	if err != nil {
		L.RaiseError("Failed to read file from zip: %s", err.Error())
		return 0
	}

	L.Push(lua.LString(data))
	return 1
}
