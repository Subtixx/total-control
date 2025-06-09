package games

import (
	"TotalControl/backend/steam"
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"os"
)

type Game struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	SteamAppID  int    `json:"steam_appid"`
	HeaderImage string `json:"header_image,omitempty"`
}

func NewGameFromIndexEntry(entry *GameIndexEntry) *Game {
	return &Game{
		ID:         entry.ID,
		Name:       entry.Name,
		SteamAppID: entry.SteamAppID,
	}
}

func (g *Game) FetchInfoFromSteam() error {
	// https://store.steampowered.com/api/appdetails?appids=
	if g.SteamAppID == 0 {
		return nil // No Steam App ID, nothing to fetch
	}

	appInfo, err := steam.GetAppDetails(g.SteamAppID)
	if err != nil {
		return err
	}
	g.Name = appInfo.Data.Name
	g.HeaderImage = appInfo.Data.HeaderImage
	return nil
}

func (g *Game) Save(filePath string) error {
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			log.Errorf("Error closing file %s: %v", filePath, err)
		}
	}(file)

	encoder := json.NewEncoder(file)
	if err := encoder.Encode(g); err != nil {
		return err
	}
	return nil
}
