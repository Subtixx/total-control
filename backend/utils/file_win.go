//go:build windows

package utils

import (
	"os"
)

func IsExecutable(file *os.File) bool {
	executableExtensions := []string{".exe", ".bat", ".cmd", ".com", ".ps1"}
	for _, ext := range executableExtensions {
		if file.Name() == ext || file.Name()[len(file.Name())-len(ext):] == ext {
			return true
		}
	}

	return false
}
