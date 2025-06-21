//go:build linux

package steam

import (
	"errors"
	"os"
)

func FindSteamInstallation() (string, error) {
	if path := os.Getenv("STEAM_ROOT"); path != "" {
		return path, nil
	}

	paths := []string{
		"$HOME/.steam/steam",
		"$HOME/.local/share/Steam",
		"$HOME/.steam/root",
	}

	for _, path := range paths {
		expanded := os.ExpandEnv(path)
		if stat, err := os.Stat(expanded); err == nil && stat.IsDir() {
			return expanded, nil
		}
	}
	return "", errors.New("steam installation not found")
}
