package plugins

import (
	"fmt"
	"github.com/google/uuid"
)

var (
	CapabilityFileSystem = "filesystem"
	CapabilityNetwork    = "network"

	FunctionalityModManagement      = "mod_management"
	FunctionalitySaveGameManagement = "save_game_management"

	ErrFileSystemAccessDenied = fmt.Sprintf("filesystem access is not allowed for this plugin")
	ErrNetworkAccessDenied    = fmt.Sprintf("network access is not allowed for this plugin")
)

type PluginInfo struct {
	Id uuid.UUID `json:"id"`

	Author      string `json:"author"`
	Name        string `json:"name"`
	Description string `json:"description"`

	Version       string   `json:"version"`
	EntryPoint    string   `json:"entry"`
	Capabilities  []string `json:"capabilities,omitempty"`
	Functionality []string `json:"functionality,omitempty"`
}

func (p *PluginInfo) Can(capability string) bool {
	if p.Capabilities == nil {
		return false
	}

	for _, canCapability := range p.Capabilities {
		if canCapability == capability {
			return true
		}
	}
	return false
}

func (p *PluginInfo) CanAccessFileSystem() bool {
	return p.Can(CapabilityFileSystem)
}

func (p *PluginInfo) CanAccessNetwork() bool {
	return p.Can(CapabilityNetwork)
}

func (p *PluginInfo) IsValid() bool {
	if p.Id == uuid.Nil {
		return false
	}
	if p.Author == "" || p.Name == "" || p.Description == "" || p.Version == "" || p.EntryPoint == "" {
		return false
	}
	return true
}

// ToString returns a string representation of the PluginInfo.
func (p *PluginInfo) ToString() string {
	return "PluginInfo{" +
		"\n\tId: " + p.Id.String() +
		"\n\tAuthor: " + p.Author +
		"\n\tName: " + p.Name +
		"\n\tDescription: " + p.Description +
		"\n\tVersion: " + p.Version +
		"\n\tEntryPoint: " + p.EntryPoint +
		"\n\tCapabilities: " + fmt.Sprintf("%v", p.Capabilities) +
		"\n}"
}
