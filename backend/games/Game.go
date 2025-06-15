package games

import (
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"os"
)

type GameMedia struct {
	Icon string `json:"icon,omitempty"`
	Hero string `json:"hero,omitempty"`
	Logo string `json:"logo,omitempty"`
}

type GameExternalID struct {
	Steam  string `json:"steam,omitempty"`
	GridDB string `json:"grid_db,omitempty"`
}

type Game struct {
	ID          string `json:"id"`
	Slug        string `json:"slug,omitempty"`
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`

	ExternalIDs GameExternalID `json:"external_ids,omitempty"`
	Media       GameMedia      `json:"media,omitempty"`
}

func NewGameFromIndexEntry(entry *GameIndexEntry) *Game {
	return &Game{
		ID:   entry.ID,
		Name: entry.Name,
	}
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
