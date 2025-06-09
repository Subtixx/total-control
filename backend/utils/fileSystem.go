package utils

import (
	"os"
	"path/filepath"
)

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
