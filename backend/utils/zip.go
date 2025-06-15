package utils

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io"
	"regexp"
)

func processSingleFile(file *zip.File, fileName string, useRegEx bool) ([]byte, error) {
	if useRegEx {
		re, err := regexp.Compile(fileName)
		if err != nil {
			return nil, fmt.Errorf("failed to compile regex for file name %s: %w", fileName, err)
		}
		if !re.MatchString(file.Name) {
			return nil, fmt.Errorf("file %s does not match the specified regex", fileName)
		}
	} else if file.Name != fileName {
		return nil, fmt.Errorf("file %s does not match the specified name", fileName)
	}

	rc, err := file.Open()
	if err != nil {
		return nil, fmt.Errorf("failed to open file %s: %w", file.Name, err)
	}

	var buf bytes.Buffer
	if _, err := io.Copy(&buf, rc); err != nil {
		return nil, fmt.Errorf("failed to copy contents of file %s: %w", file.Name, err)
	}
	if err := rc.Close(); err != nil {
		return nil, fmt.Errorf("failed to close file %s: %w", file.Name, err)
	}

	content := buf.Bytes()
	return content, nil
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
		return processSingleFile(f, fileName, useRegEx)
	}
	return nil, fmt.Errorf("file %s not found in zip archive", fileName)
}
