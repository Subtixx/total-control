package plugins

import (
	"TotalControl/backend/utils"
	"encoding/json"
	"errors"
	log "github.com/sirupsen/logrus"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

var (
	ErrInvalidRepository     = errors.New("invalid plugin repository: id and url must not be empty")
	ErrFetchRepositoryFailed = errors.New("failed to fetch plugin repository: status code not OK")
)

const DefaultPluginRepositoryId = "1767017e-590a-40af-a18b-b036d744a766"
const DefaultPluginRepositoryUrl = "https://raw.githubusercontent.com/subtixx/total-control/main/plugins/repository.json"

type PluginRepositoryInfo struct {
	Name      string `json:"name"`
	Version   string `json:"version"`
	Changelog string `json:"changelog"`
}

type PluginRepository struct {
	Id      string                  `json:"id"`
	Url     string                  `json:"url"`
	Plugins []*PluginRepositoryInfo `json:"plugins"`
}

func NewDefaultPluginRepository() (*PluginRepository, error) {
	return NewPluginRepository(DefaultPluginRepositoryId, DefaultPluginRepositoryUrl)
}

func NewPluginRepository(id string, url string) (*PluginRepository, error) {
	if id == "" || url == "" {
		return nil, ErrInvalidRepository
	}

	if _, err := os.Stat(filepath.Join(utils.GetAppDataPath(), "repositories", id+".json")); err == nil {
		log.Infof("Plugin repository '%s' already exists, loading from file", id)
		file, err := os.ReadFile(filepath.Join(utils.GetAppDataPath(), "repositories", id+".json"))
		if err != nil {
			return nil, err
		}
		var repository PluginRepository
		err = json.Unmarshal(file, &repository)
		if err != nil {
			return nil, err
		}
		return &repository, nil
	}

	response, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	if response.StatusCode != http.StatusOK {
		return nil, ErrFetchRepositoryFailed
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Errorf("Failed to close response body: %v", err)
		}
	}(response.Body)
	jsonData, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	var repository PluginRepository
	repository.Id = id
	repository.Url = url
	err = json.Unmarshal(jsonData, &repository)
	if err != nil {
		return nil, err
	}

	// Write the repository to a file
	repoDir := filepath.Join(utils.GetAppDataPath(), "repositories")
	if _, err := os.Stat(repoDir); os.IsNotExist(err) {
		if err := os.MkdirAll(repoDir, 0755); err != nil {
			log.Errorf("Failed to create repository directory: %v", err)
			return nil, err
		}
	}
	repoFile := filepath.Join(repoDir, id+".json")
	file, err := os.Create(repoFile)
	if err != nil {
		log.Errorf("Failed to create repository file: %v", err)
		return nil, err
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			log.Errorf("Failed to close repository file: %v", err)
		}
	}(file)
	_, err = file.Write(jsonData)
	if err != nil {
		log.Errorf("Failed to write repository file: %v", err)
		return nil, err
	}
	log.Infof("Plugin repository '%s' created successfully at %s", id, repoFile)
	return &repository, nil
}
