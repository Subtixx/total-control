//go:build !windows

package utils

import (
	"os"
	"path/filepath"
)

func IsHidden(path string) bool {
	if path == "" {
		return false
	}

	fileInfo, err := os.Stat(path)
	if err != nil {
		return false // If we can't stat the file, assume it's not hidden
	}

	if !fileInfo.IsDir() {
		return false // Only directories can be hidden in this context
	}

	base := filepath.Base(path)
	return base[0] == '.'
}
