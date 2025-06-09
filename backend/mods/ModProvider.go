package mods

type GameVersion struct {
	Version string `json:"version"` // The version of the game this mod is compatible with
}
type Mod struct {
	ID           string        `json:"id"`
	Name         string        `json:"name"`
	Description  string        `json:"description,omitempty"`
	Author       string        `json:"author,omitempty"`
	Version      string        `json:"version,omitempty"`
	Enabled      bool          `json:"enabled"`
	Dependencies []string      `json:"dependencies,omitempty"`
	DownloadURL  string        `json:"download_url,omitempty"`
	IconURL      string        `json:"icon_url,omitempty"`
	HeaderImage  string        `json:"header_image,omitempty"`
	GameVersions []GameVersion `json:"game_versions,omitempty"` // List of game versions this mod is compatible with

	Image []byte `json:"-"`

	GameID string `json:"game_id,omitempty"` // ID of the game this mod belongs to
}

type ModProvider interface {
	GetMods() ([]Mod, error)
	GetModByID(id string) (*Mod, error)

	GetGameModDirectory() string // Returns the directory where mods for this game are stored

	AddMod(mod Mod) error
	RemoveMod(id string) error
	UpdateMod(mod Mod) error
	ListGameMods() ([]Mod, error)

	GetGameID() string // Returns the ID of the game this provider is for
}
