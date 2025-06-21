package steam

import log "github.com/sirupsen/logrus"

var blacklistedAppIDs = map[string]string{
	"1070560": "Steam Linux Runtime 1.0 (scout)",
	"1391110": "Steam Linux Runtime 2.0 (soldier)",
}

func GetInstalledGames() []string {
	libraries, err := GetSteamLibraries()
	if err != nil {
		return nil
	}

	// Map all library paths to their respective games
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
