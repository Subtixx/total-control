package games

import (
	"TotalControl/backend/utils"
	"encoding/json"
	"errors"
	log "github.com/sirupsen/logrus"
	"os"
	"strconv"
	"strings"
)

type GameIndex struct {
	Entries []GameIndexEntry `json:"entries"`
}

func NewGameIndex(filePath string) (*GameIndex, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, errors.New("failed to read game index file: " + err.Error())
	}
	var index GameIndex
	err = json.Unmarshal(data, &index)
	if err != nil {
		return nil, errors.New("failed to parse game index file: " + err.Error())
	}
	return &index, nil
}

func (g *GameIndex) GetGameByID(id string) (*GameIndexEntry, error) {
	for _, entry := range g.Entries {
		if entry.ID == id {
			return &entry, nil
		}
	}
	return nil, errors.New("game not found with ID: " + id)
}

func (g *GameIndex) GetGameEntryByName(name string) (*GameIndexEntry, error) {
	for _, entry := range g.Entries {
		if entry.Name == name {
			return &entry, nil
		}
	}
	return nil, errors.New("game not found with name: " + name)
}

func (g *GameIndex) GetGameEntryBySteamAppID(steamAppID int) (*GameIndexEntry, error) {
	for _, entry := range g.Entries {
		if entry.SteamAppID == steamAppID {
			return &entry, nil
		}
	}
	return nil, errors.New("game not found with Steam App ID: " + strconv.Itoa(steamAppID))
}

func (g *GameIndex) GetGameEntryByExecutablePath(executablePath string) (*GameIndexEntry, error) {
	for _, entry := range g.Entries {
		for _, executable := range entry.Executables {
			if strings.ToLower(executable.Path) == strings.ToLower(executablePath) {
				return &entry, nil
			}
		}
	}
	return nil, errors.New("game not found with executable path: " + executablePath)
}

func (g *GameIndex) DetectGameFromSteamAppID(steamAppID int) (*Game, error) {
	gameEntry, err := g.GetGameEntryBySteamAppID(steamAppID)
	if err == nil {
		return NewGameFromIndexEntry(gameEntry), nil
	}
	return nil, errors.New("game not found for Steam App ID: " + strconv.Itoa(steamAppID))
}

func (g *GameIndex) DetectGameFromExecutableName(filePath string) (*Game, error) {
	// On windows its easy with .exe but on linux we have to check the file name and the ELF header
	// Get all exe files in the directory and subdirectories
	files, err := utils.GetFilesByWildcards(filePath, []string{"*.exe", "*.elf", "*.bin", "*.out", "*?"})
	if err != nil {
		return nil, errors.New("failed to read directory: " + err.Error())
	}

	for _, file := range files {
		file, err := utils.NewFile(file)
		if err != nil {
			continue // Skip files that cannot be read
		}

		// Check if the file has a known executable extension
		if file.IsExecutable() {
			// This is a Windows executable, we can check the game index
			gameEntry, err := g.GetGameEntryByExecutablePath(file.FileName)
			if err == nil {
				return NewGameFromIndexEntry(gameEntry), nil
			}
		}
	}
	return nil, errors.New("no game detected from executable name in " + filePath)
}

func (g *GameIndex) DetectGame(filePath string) (*Game, error) {
	// Check if the path is a directory
	info, err := os.Stat(filePath)
	if err != nil {
		return nil, errors.New("path does not exist: " + filePath)
	}

	if !info.IsDir() {
		return nil, errors.New("path is not a directory: " + filePath)
	}

	// First try steam_appid.txt
	if _, err := os.Stat(filePath + "/steam_appid.txt"); err == nil {
		data, err := os.ReadFile(filePath + "/steam_appid.txt")
		if err != nil {
			return nil, errors.New("failed to read steam_appid.txt: " + err.Error())
		}
		steamAppIDStr := string(data)
		steamAppID, err := strconv.Atoi(steamAppIDStr)
		if err == nil {
			log.Debugf("Using steam_appid.txt to detect game: %s", steamAppIDStr)
			return g.DetectGameFromSteamAppID(steamAppID)
		}
	}

	game, err := g.DetectGameFromExecutableName(filePath)
	if err == nil {
		log.Debugf("Detected game from executable name in directory: %s", filePath)
		return game, nil
	}

	return nil, errors.New("no game detected at " + filePath)
}
