package utils

import (
	"archive/zip"
	log "github.com/sirupsen/logrus"
	"io"
	"os"
	"path/filepath"
	"strings"
)

type File struct {
	FileName      string
	FilePath      string
	FileExtension string
	Contents      []byte // Contents of the file, if needed

	FileInfo os.FileInfo
}

func NewFile(filePath string) (*File, error) {
	// Extract the file name from the file path
	fileName := filePath
	if lastSlash := strings.LastIndex(filePath, "/"); lastSlash != -1 {
		fileName = filePath[lastSlash+1:]
	}

	// Get file info
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		if !os.IsNotExist(err) {
			log.Errorf("Error getting file info for %s: %v", filePath, err)
		}
		return nil, err
	}

	// Get the file extension
	fileExtension := filepath.Ext(fileName)

	return &File{
		FileName:      fileName,
		FilePath:      filePath,
		FileExtension: fileExtension,
		FileInfo:      fileInfo,
	}, nil
}

func (file *File) IsHidden() bool {
	return strings.HasPrefix(file.FileName, ".")
}

func (file *File) ToString() string {
	return file.FileName + " (" + file.FilePath + ")"
}

func (file *File) IsExecutable() bool {
	if file.FileInfo.IsDir() {
		return false // Directories are not executable
	}

	// Check if the file has a known executable extension
	executableExtensions := []string{".exe", ".bat", ".sh", ".bin", ".out"}
	for _, ext := range executableExtensions {
		if strings.EqualFold(file.FileExtension, ext) {
			return true
		}
	}

	// On Unix-like systems, check the file permissions
	if file.FileInfo.Mode()&0111 != 0 {
		return true // The file is executable
	}

	return false // Not an executable file
}

func (file *File) GetZipFileContents(filePath []string) (map[string]*File, error) {
	reader, err := os.Open(file.FilePath)
	if err != nil {
		return nil, err
	}
	defer func(reader *os.File) {
		err := reader.Close()
		if err != nil {
			log.Errorf("Error closing file %s: %v", file.FilePath, err)
		}
	}(reader)

	stat, err := reader.Stat()
	if err != nil {
		return nil, err
	}

	zipReader, err := zip.NewReader(reader, stat.Size())
	if err != nil {
		return nil, err
	}

	files := make(map[string]*File)
	for _, zipFile := range zipReader.File {
		if zipFile.FileInfo().IsDir() {
			continue // Skip directories
		}
		// Check if the file matches any of the specified paths
		matched := ""
		for _, path := range filePath {
			if strings.HasSuffix(zipFile.Name, path) && files[path] == nil {
				matched = path
				break
			}
		}
		if matched == "" {
			continue // Skip files that do not match the specified paths
		}
		log.Debugf("Found file in zip: %s", zipFile.Name)
		fileReader, err := zipFile.Open()
		if err != nil {
			return nil, err
		}

		fileContents, err := io.ReadAll(fileReader)
		if err != nil {
			return nil, err
		}
		if err := fileReader.Close(); err != nil {
			log.Errorf("Error closing file reader for %s: %v", zipFile.Name, err)
		}

		files[matched] = &File{
			FileName:      zipFile.Name,
			FilePath:      zipFile.Name,
			FileExtension: filepath.Ext(zipFile.Name),
			Contents:      fileContents,
		}
	}
	if len(files) == 0 {
		return nil, os.ErrNotExist // File not found in the zip archive
	}
	return files, nil
}

func (file *File) GetZipFileContent(filePath string) (*File, error) {
	files, err := file.GetZipFileContents([]string{filePath})
	if err != nil {
		return nil, err
	}
	if len(files) == 0 {
		return nil, os.ErrNotExist // File not found in the zip archive
	}
	return files[filePath], nil
}
