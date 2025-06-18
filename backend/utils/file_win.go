//go:build windows

package utils

import (
	"os"
	"path/filepath"
)

var ExecutableExtensions = []string{".exe"}

func IsExecutable(filePath string) bool {
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		return false
	}

	fileExt := filepath.Ext(fileInfo.Name())
	for _, ext := range ExecutableExtensions {
		if fileExt == ext {
			return true
		}
	}

	return false
}
