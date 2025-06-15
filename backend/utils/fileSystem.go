package utils

import (
	"path/filepath"
)

func GetFileNameWithoutExtension(filePath string) string {
	// Extract the file name from the file path
	fileName := filepath.Base(filePath)
	// Get the file extension
	fileExtension := filepath.Ext(fileName)
	// Remove the extension from the file name
	return fileName[:len(fileName)-len(fileExtension)]
}
