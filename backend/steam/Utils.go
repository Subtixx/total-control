package steam

import (
	"errors"
	log "github.com/sirupsen/logrus"
	"os"
	"path"
)

var blacklistedAppIDs = map[string]string{
	"1070560": "Steam Linux Runtime 1.0 (scout)",
	"1391110": "Steam Linux Runtime 2.0 (soldier)",
}

func GetAppSchemaFile(steamLibraryPath string, appId string) (*AppSchemaFile, error) {
	if appId == "" {
		return nil, errors.New("appId cannot be empty")
	}

	appManifestPath := path.Join(steamLibraryPath, "steamapps", "appmanifest_"+appId+".acf")
	if _, err := os.Stat(appManifestPath); os.IsNotExist(err) {
		log.Errorf("App manifest file does not exist for appId %s at path %s", appId, appManifestPath)
		return nil, errors.New("app manifest file does not exist")
	}

	vdfData, err := ReadVDF(appManifestPath)
	if err != nil {
		log.Errorf("Failed to read VDF data for appId %s: %v", appId, err)
		return nil, err
	}

	appSchemaFile := NewAppSchemaFile(vdfData)
	if appSchemaFile == nil {
		log.Errorf("failed to parse app schema file for appId %s", appId)
		return nil, errors.New("failed to parse app schema file")
	}

	return appSchemaFile, nil
}

func GetInstalledGames(libraries []*LibraryFolder) []*AppSchemaFile {
	var installedGames []*AppSchemaFile
	for _, library := range libraries {
		if library == nil {
			continue
		}

		for appId, manifestId := range library.Apps {
			if appId == "" {
				log.Warnf("App ID is empty in library %s", library.Path)
				continue
			}

			if manifestId == "" {
				log.Warnf("Manifest ID is empty for appId %s in library %s", appId, library.Path)
				continue
			}

			if name, exists := blacklistedAppIDs[appId]; exists {
				log.Debugf("Skipping blacklisted app: %s (%s)", appId, name)
				continue
			}

			appSchemaFile, err := GetAppSchemaFile(library.Path, appId)
			if err != nil {
				log.Errorf("Failed to get app schema file for appId %s: %v", appId, err)
				continue // Skip if we can't get the schema file
			}

			installedGames = append(installedGames, appSchemaFile)
		}
	}

	return installedGames
}

func GetInstalledGamePaths(libraries []*LibraryFolder) []string {
	var installedGames []string
	for _, library := range libraries {
		if library == nil {
			continue
		}

		for appId := range library.Apps {
			if appId == "" {
				continue // Skip empty app IDs
			}

			if name, exists := blacklistedAppIDs[appId]; exists {
				log.Debugf("Skipping blacklisted app: %s (%s)", appId, name)
				continue // Skip blacklisted app IDs
			}

			installedGames = append(installedGames, appId)
		}
	}

	return installedGames
}
