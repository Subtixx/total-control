//go:build !windows

package utils

import (
	"os"
)

func IsExecutable(file *os.File) bool {
	// On Linux, we check the file permissions for executables
	fileInfo, err := file.Stat()
	if err != nil {
		return false // If we can't get file info, assume it's not executable
	}

	// Check if the file is a regular file and has execute permissions
	return fileInfo.Mode()&0111 != 0 && !fileInfo.IsDir()
}
