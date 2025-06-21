package plugins

import (
	"TotalControl/backend/utils"
	"fmt"
	log "github.com/sirupsen/logrus"
)

var (
	FunctionalityModManagement      = "mod_management"
	FunctionalitySaveGameManagement = "save_game_management"
)

type PluginInfo struct {
	Id string `json:"id"`

	Author      string `json:"author"`
	Name        string `json:"name"`
	Description string `json:"description"`

	Version            string             `json:"version"`
	EntryPoint         string             `json:"entry"`
	Functionality      []string           `json:"functionality,omitempty"`
	SettingDefinitions SettingDefinitions `json:"settings"`
}

func (p *PluginInfo) IsValid() bool {
	if !utils.IsValidPackageName(p.Id) {
		p.Logger().Error("Invalid package name")
		return false
	}

	if p.Author == "" || p.Name == "" || p.Description == "" || p.Version == "" || p.EntryPoint == "" {
		return false
	}
	return true
}

func (p *PluginInfo) Logger() *log.Entry {
	return log.WithFields(log.Fields{
		"plugin": p.Id,
	})
}

// ToString returns a string representation of the PluginInfo.
func (p *PluginInfo) String() string {
	return "PluginInfo{" +
		"\n\tId: " + p.Id +
		"\n\tAuthor: " + p.Author +
		"\n\tName: " + p.Name +
		"\n\tDescription: " + p.Description +
		"\n\tVersion: " + p.Version +
		"\n\tEntryPoint: " + p.EntryPoint +
		"\n\tFunctionality: " + fmt.Sprintf("%v", p.Functionality) +
		"\n" + p.SettingDefinitions.String() +
		"\n}"
}
