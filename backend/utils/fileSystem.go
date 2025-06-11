package utils

import (
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

	files, err := GetFilesByWildcards(dir, patterns)
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

func GetFilesByWildcards(dir string, patterns []string) ([]string, error) {
	var matchedFiles []string

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			// Skip directories
			return nil
		}

		matched := false
		for _, pattern := range patterns {
			// We need a special pattern for a filename without an extension
			if pattern == "*?" {
				re := filepath.Ext(info.Name())
				if re == "" {
					matched = true
				} else {
					continue // Skip this pattern if the file has an extension
				}
			} else {
				if matched, err = filepath.Match(pattern, info.Name()); err != nil {
					return err
				}
			}
			if matched {
				matchedFiles = append(matchedFiles, path)
				break // No need to check other patterns if one matches
			}
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return matchedFiles, nil
}

// GetFilesByWildcard Function to get all files in a directory and its subdirectories and filtering them by a wildcard pattern
func GetFilesByWildcard(dir string, pattern string) ([]string, error) {
	return GetFilesByWildcards(dir, []string{pattern})
}
