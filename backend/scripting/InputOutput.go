package scripting

import (
	"TotalControl/backend/plugins"
	"TotalControl/backend/utils"
	lua "github.com/yuin/gopher-lua"
	"os"
	"path/filepath"
)

func LuaGetFilesInDirectory(L *lua.LState) int {
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

	dir := L.ToString(1)
	patterns := make([]string, 0)

	if L.Get(2).Type() == lua.LTTable {
		patternTable := L.ToTable(2)
		patternTable.ForEach(func(key lua.LValue, value lua.LValue) {
			if value.Type() == lua.LTString {
				patterns = append(patterns, value.String())
			}
		})
	} else if L.Get(2).Type() == lua.LTString {
		patterns = append(patterns, L.ToString(2))
	} else {
		L.RaiseError("Second argument must be a string or a table of strings")
		return 0
	}

	files, err := utils.GetFilesByWildcards(dir, patterns)
	if err != nil {
		L.RaiseError("Error getting files: %v", err)
		return 0
	}

	luaFiles := L.CreateTable(len(files), 0)
	for _, file := range files {
		luaFiles.Append(lua.LString(file))
	}

	L.Push(luaFiles)
	return 1
}

func LuaGetFileName(L *lua.LState) int {
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

	filePath := L.ToString(1)
	if filePath == "" {
		L.Push(lua.LNil)
		return 1
	}

	fileName := filepath.Base(filePath)
	L.Push(lua.LString(fileName))
	return 1
}

func LuaGetFileContent(L *lua.LState) int {
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

	filePath := L.ToString(1)
	if filePath == "" {
		L.Push(lua.LNil)
		return 1
	}

	content, err := os.ReadFile(filePath)
	if err != nil {
		L.RaiseError("Error reading file '%s': %v", filePath, err)
		return 0
	}
	L.Push(lua.LString(content))
	return 1
}

func LuaExtendIoTable(l *lua.LState) {
	table := GetOrCreateTable(l, "io")
	l.SetField(table, "getFilesInDirectory", l.NewFunction(LuaGetFilesInDirectory))
	l.SetField(table, "getFileName", l.NewFunction(LuaGetFileName))
	l.SetField(table, "getFileContent", l.NewFunction(LuaGetFileContent))
}
