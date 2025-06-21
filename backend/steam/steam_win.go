//go:build windows

package steam

import (
	"errors"
	"golang.org/x/sys/windows/registry"
	"os"
	"path/filepath"
)

func FindSteamInstallation() (string, error) {
	key, err := registry.OpenKey(registry.CURRENT_USER, `Software\Valve\Steam`, registry.QUERY_VALUE)
	if err == nil {
		defer key.Close()
		if path, _, err := key.GetStringValue("SteamPath"); err == nil {
			return path, nil
		}
	}

	paths := []string{
		`C:\Program Files (x86)\Steam`,
		`C:\Program Files\Steam`,
	}
	for _, p := range paths {
		if _, err := os.Stat(filepath.Join(p, "steam.exe")); err == nil {
			return p, nil
		}
	}
	return "", errors.New("steam installation not found")
}
