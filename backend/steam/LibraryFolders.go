package steam

import (
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"os"
	"path"
)

type LibraryFolder struct {
	Path                   string
	Label                  string
	ContentId              string
	TotalSize              string
	UpdateCleanBytesTally  string
	TimeLastUpdateVerified string
	Apps                   map[string]string
}

func GetSteamLibraries() ([]*LibraryFolder, error) {
	steamInstall, err := FindSteamInstallation()
	if err != nil {
		return nil, err
	}

	libraryFoldersPath := path.Join(steamInstall, "steamapps", "libraryfolders.vdf")
	if _, err := os.Stat(libraryFoldersPath); os.IsNotExist(err) {
		return nil, err
	}

	libraries, err := ReadLibraryFolders(libraryFoldersPath)
	if err != nil {
		return nil, err
	}

	return libraries, nil
}

func ReadLibraryFolders(filePath string) ([]*LibraryFolder, error) {
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		log.Errorf("Failed to find libraryfolders.vdf: %s", err)
		return nil, err
	}

	vdf, err := ReadVDF(filePath)
	if err != nil {
		return nil, err
	}

	libraries, ok := vdf["libraryfolders"].(map[string]interface{})
	if !ok {
		return nil, errors.New("Invalid libraryfolders format")
	}

	var result []*LibraryFolder
	for key, value := range libraries {
		if key == "version" {
			continue // Skip version key
		}
		libData, ok := value.(map[string]interface{})
		if !ok {
			log.Warnf("Skipping invalid library folder entry: %s", key)
			continue
		}

		library := &LibraryFolder{
			Path:                   libData["path"].(string),
			Label:                  libData["label"].(string),
			ContentId:              libData["contentid"].(string),
			TotalSize:              libData["totalsize"].(string),
			UpdateCleanBytesTally:  libData["update_clean_bytes_tally"].(string),
			TimeLastUpdateVerified: libData["time_last_update_verified"].(string),
			Apps:                   make(map[string]string),
		}
		apps, ok := libData["apps"].(map[string]interface{})
		if ok {
			for appId, appData := range apps {
				appDataStr, ok := appData.(string)
				if !ok {
					log.Warnf("Skipping invalid app data for ID %s: %v", appId, appData)
					continue
				}
				library.Apps[appId] = appDataStr
			}
		}
		result = append(result, library)
	}

	return result, nil
}
