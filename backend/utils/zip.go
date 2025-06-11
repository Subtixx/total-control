package utils

import (
	"archive/zip"
	"bytes"
	"fmt"
	lua "github.com/yuin/gopher-lua"
	"io"
	"regexp"
)

// Lua API
func LuaReadFilesFromZip(L *lua.LState) int {
	zipPath := L.ToString(1)
	if zipPath == "" {
		L.Push(lua.LNil)
		return 1
	}

	files, err := ReadFilesFromZip(zipPath)
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

	data, err := ReadFileFromZip(zipPath, fileName, useRegEx)
	if err != nil {
		L.RaiseError("Failed to read file from zip: %s", err.Error())
		return 0
	}

	L.Push(lua.LString(data))
	return 1
}

// ReadFilesFromZip reads all files from a zip archive and returns a map of filenames to their contents.
func ReadFilesFromZip(zipPath string) (map[string][]byte, error) {
	result := make(map[string][]byte)

	r, err := zip.OpenReader(zipPath)
	if err != nil {
		return nil, err
	}

	for _, f := range r.File {
		rc, err := f.Open()
		if err != nil {
			return nil, err
		}
		var buf bytes.Buffer
		if _, err := io.Copy(&buf, rc); err != nil {
			return nil, err
		}
		result[f.Name] = buf.Bytes()
		if err := rc.Close(); err != nil {
			return nil, fmt.Errorf("failed to close file %s: %w", f.Name, err)
		}
	}
	if err := r.Close(); err != nil {
		return nil, fmt.Errorf("failed to close zip reader: %w", err)
	}
	return result, nil
}

// ReadFileFromZip reads a specific file from a zip archive and returns its contents.
func ReadFileFromZip(zipPath, fileName string, useRegEx bool) ([]byte, error) {
	r, err := zip.OpenReader(zipPath)
	if err != nil {
		return nil, err
	}

	for _, f := range r.File {
		if useRegEx {
			re, err := regexp.Compile(fileName)
			if err != nil {
				return nil, fmt.Errorf("failed to compile regex for file name %s: %w", fileName, err)
			}
			if !re.MatchString(f.Name) {
				continue
			}
		} else if f.Name != fileName {
			continue
		}
		rc, err := f.Open()
		if err != nil {
			return nil, err
		}

		var buf bytes.Buffer
		if _, err := io.Copy(&buf, rc); err != nil {
			return nil, err
		}
		if err := rc.Close(); err != nil {
			return nil, fmt.Errorf("failed to close file %s: %w", f.Name, err)
		}
		if err := r.Close(); err != nil {
			return nil, fmt.Errorf("failed to close zip reader: %w", err)
		}
		return buf.Bytes(), nil
	}
	return nil, fmt.Errorf("file %s not found in zip archive", fileName)
}
