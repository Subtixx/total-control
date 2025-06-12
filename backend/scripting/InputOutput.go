package scripting

import (
	"TotalControl/backend/utils"
	log "github.com/sirupsen/logrus"
	lua "github.com/yuin/gopher-lua"
	"os"
	"path/filepath"
)

func LuaGetFilesInDirectory(L *lua.LState) int {
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
func LuaReadFilesFromZip(L *lua.LState) int {
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

func luaExtendIoTable(l *lua.LState) {
	ioTable := l.GetGlobal("io")
	tbl, ok := ioTable.(*lua.LTable)
	if !ok {
		log.Warnf("io table not found, creating a new one")
		tbl = l.NewTable()
		l.SetGlobal("io", tbl)
	}
	l.SetField(tbl, "readFileFromZip", l.NewFunction(LuaReadFileFromZip))
	l.SetField(tbl, "readFilesFromZip", l.NewFunction(LuaReadFilesFromZip))
	l.SetField(tbl, "getFilesInDirectory", l.NewFunction(LuaGetFilesInDirectory))
	l.SetField(tbl, "getFileName", l.NewFunction(LuaGetFileName))
	l.SetField(tbl, "getFileContent", l.NewFunction(LuaGetFileContent))
}
