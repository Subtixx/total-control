package plugins

import "github.com/google/uuid"

type PluginInfo struct {
	Id uuid.UUID `json:"id"`

	Author      string `json:"author"`
	Name        string `json:"name"`
	Description string `json:"description"`

	Version    string `json:"version"`
	EntryPoint string `json:"entry"`
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
		"\n}"
}
