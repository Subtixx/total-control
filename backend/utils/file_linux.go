//go:build !windows

package utils

import (
	"os"
)

func IsExecutable(filePath string) bool {
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		return false
	}

	// Check if the file is a regular file and has execute permissions
	return fileInfo.Mode()&0111 != 0 && !fileInfo.IsDir()
}
