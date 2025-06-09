package factorio

import (
	"TotalControl/backend/mods"
	"TotalControl/backend/utils"
	"encoding/json"
	"errors"
	log "github.com/sirupsen/logrus"
	"os"
	"path/filepath"
	"strings"
)

type FactorioModProvider struct {
}

func (modProvider *FactorioModProvider) ReadModInfo(modPath string) (*ModInfo, error) {
	if filepath.Ext(modPath) != ".zip" {
		return nil, errors.New("mod file must be a .zip file")
	}
	// Factorio mods typically have a control.lua or info.json file inside the zip.
	modZip, err := utils.NewFile(modPath)
	if err != nil {
		return nil, err
	}
	files, err := modZip.GetZipFileContents([]string{"info.json", "thumbnail.png"})
	if err != nil {
		return nil, err
	}
	var infoFile *utils.File
	var ok bool
	if infoFile, ok = files["info.json"]; !ok {
		return nil, errors.New("info.json not found in " + modPath + ", is this a valid Factorio mod?")
	}

	var mod ModInfo
	if err := json.Unmarshal(infoFile.Contents, &mod); err != nil {
		return nil, err
	}

	if thumbnail, ok := files["thumbnail.png"]; ok {
		mod.Image = thumbnail.Contents
	}

	return &mod, nil
}

func (modProvider *FactorioModProvider) GetGameModDirectory() string {
	// Factorio mods are typically stored in the user's Factorio directory.
	// This is usually located at:
	// - Linux: ~/.factorio/mods/
	// - Windows: C:\Users\<Username>\AppData\Roaming\Factorio\mods\
	// - macOS: ~/Library/Application Support/factorio/mods/
	if os.Getenv("FACTORIO_MODS_DIR") != "" {
		return os.Getenv("FACTORIO_MODS_DIR")
	}

	if utils.GetOperatingSystem() == utils.WindowsOS {
		return os.Getenv("APPDATA") + "\\Factorio\\mods"
	} else if utils.GetOperatingSystem() == utils.MacOS {
		return os.Getenv("HOME") + "/Library/Application Support/factorio/mods"
	} else if utils.GetOperatingSystem() == utils.LinuxOS {
		return os.Getenv("HOME") + "/.factorio/mods"
	}
	panic("Unsupported OS: " + utils.GetOperatingSystem())
}

func (modProvider *FactorioModProvider) GetModFile(modID string) (string, error) {
	// The gamemods contain versions.... Shortcuts-ick_2.0.7.zip so we need to remove the version part.
	modDirectory := modProvider.GetGameModDirectory()
	filesInDir, err := os.ReadDir(modDirectory)
	if err != nil {
		return "", err
	}
	for _, file := range filesInDir {
		if file.IsDir() {
			continue // Skip directories
		}
		if file.Name() == modID || file.Name() == modID+".zip" {
			return filepath.Join(modDirectory, file.Name()), nil
		}
		// Check if the file name starts with the mod ID and ends with .zip
		if len(file.Name()) > len(modID)+4 && file.Name()[:len(modID)] == modID && file.Name()[len(file.Name())-4:] == ".zip" {
			return filepath.Join(modDirectory, file.Name()), nil
		}
	}
	return "", errors.New("mod file not found: " + modID)
}

func (modProvider *FactorioModProvider) GetMods() ([]mods.Mod, error) {
	modDirectory := modProvider.GetGameModDirectory()
	files, err := os.ReadDir(modDirectory)
	if err != nil {
		return nil, err
	}

	// If we have a mod-list.json file, we should read it and return the mods from there.
	modListFile := filepath.Join(modDirectory, "mod-list.json")
	if _, err := os.Stat(modListFile); err == nil {
		modListData, err := os.ReadFile(modListFile)
		if err != nil {
			return nil, err
		}
		var modList struct {
			Mods []struct {
				Name    string `json:"name"`
				Enabled bool   `json:"enabled"`
			} `json:"mods"`
		}
		if err := json.Unmarshal(modListData, &modList); err != nil {
			return nil, err
		}
		var foundMods []mods.Mod
		for _, mod := range modList.Mods {
			// Check if the zip file exists
			modPath, err := modProvider.GetModFile(mod.Name)
			if _, err := os.Stat(modPath); os.IsNotExist(err) {
				log.Debugf("Mod %s not found in %s", mod.Name, modDirectory)
				continue
			}
			factorioMod, err := modProvider.ReadModInfo(modPath)
			if err != nil {
				return nil, err
			}
			versions := []mods.GameVersion{
				{
					Version: factorioMod.FactorioVersion,
				},
			}
			foundMods = append(foundMods, mods.Mod{
				ID:           strings.ToLower(factorioMod.Name),
				GameID:       modProvider.GetGameID(),
				Name:         factorioMod.Title,
				Description:  factorioMod.Description,
				Version:      factorioMod.Version,
				Author:       factorioMod.Author,
				Image:        factorioMod.Image,
				GameVersions: versions,
				Enabled:      mod.Enabled,
			})
		}
		return foundMods, nil
	}

	var foundMods []mods.Mod
	for _, file := range files {
		fileExtension := filepath.Ext(file.Name())
		if file.IsDir() || fileExtension != ".zip" {
			continue // Skip directories and non-zip files
		}
		mod := mods.Mod{
			ID:      file.Name(),
			GameID:  modProvider.GetGameID(),
			Enabled: true,
		}
		// Here we could potentially read the mod's metadata from its contents,
		// but for now, we just use the file name as the mod ID.
		foundMods = append(foundMods, mod)
	}
	return foundMods, nil
}

func (modProvider *FactorioModProvider) GetModByID(id string) (*mods.Mod, error) {
	//TODO implement me
	panic("implement me")
}

func (modProvider *FactorioModProvider) AddMod(mod mods.Mod) error {
	//TODO implement me
	panic("implement me")
}

func (modProvider *FactorioModProvider) RemoveMod(id string) error {
	//TODO implement me
	panic("implement me")
}

func (modProvider *FactorioModProvider) UpdateMod(mod mods.Mod) error {
	//TODO implement me
	panic("implement me")
}

func (modProvider *FactorioModProvider) ListGameMods() ([]mods.Mod, error) {
	//TODO implement me
	panic("implement me")
}

func (modProvider *FactorioModProvider) GetGameID() string {
	return "factorio"
}
