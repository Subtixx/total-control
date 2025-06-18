package utils

import (
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"io"
	"os"
	"strings"
)

func FileExists(filePath string) bool {
	_, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		return false // File does not exist
	}
	return true // File exists
}

func FileTouch(filePath string, content ...string) error {
	if FileExists(filePath) {
		return nil // File already exists, no need to touch it
	}

	file, err := os.Create(filePath)
	if err != nil {
		return err // Error creating the file
	}

	if len(content) > 0 {
		if _, err := file.Write([]byte(content[0])); err != nil {
			err := file.Close()
			if err != nil {
				return err
			}
		}
	}

	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			log.Errorf("Error closing file %s: %v", filePath, err)
		}
	}(file)

	return nil
}

func ReadFile(filePath string) ([]byte, error) {
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return nil, fmt.Errorf("file %s does not exist", filePath)
	}

	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			log.Errorf("Error closing file %s: %v", filePath, err)
		}
	}(file)

	data, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func ReadJSONFile(filePath string) (map[string]interface{}, error) {
	data, err := ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	var jsonData map[string]interface{}
	if err := json.Unmarshal(data, &jsonData); err != nil {
		return nil, err
	}
	return jsonData, nil
}

func SanitizeFileName(fileName string) string {
	// Replace all non-alphanumeric characters with underscores
	sanitized := strings.Map(func(r rune) rune {
		if r >= 'a' && r <= 'z' || r >= 'A' && r <= 'Z' || r >= '0' && r <= '9' || r == '_' {
			return r
		}
		return '_'
	}, fileName)

	// Remove leading and trailing underscores
	sanitized = strings.Trim(sanitized, "_")

	return sanitized
}
