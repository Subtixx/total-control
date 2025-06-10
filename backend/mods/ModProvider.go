package mods

import (
	"errors"
	lua "github.com/yuin/gopher-lua"
)

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

func NewModFromLuaTable(modTable *lua.LTable) (*Mod, error) {
	if modTable == nil {
		return nil, nil // No mod data provided
	}
	if modTable.Type() != lua.LTTable {
		return nil, errors.New("modTable is not a Lua table")
	}
	if modTable.RawGetString("id").Type() != lua.LTString {
		return nil, errors.New("modTable is missing 'id' field or it is not a string")
	}
	if modTable.RawGetString("name").Type() != lua.LTString {
		return nil, errors.New("modTable is missing 'name' field or it is not a string")
	}
	if modTable.RawGetString("version").Type() != lua.LTString {
		return nil, errors.New("modTable is missing 'version' field or it is not a string")
	}
	if modTable.RawGetString("enabled").Type() != lua.LTBool {
		return nil, errors.New("modTable is missing 'enabled' field or it is not a boolean")
	}
	if modTable.RawGetString("game_id").Type() != lua.LTString {
		return nil, errors.New("modTable is missing 'game_id' field or it is not a string")
	}

	var gameVersions []GameVersion
	gameVersionsTable := modTable.RawGetString("game_versions")
	if gameVersionsTable.Type() == lua.LTTable {
		gameVersionsTable.(*lua.LTable).ForEach(func(_, value lua.LValue) {
			if value.Type() == lua.LTTable {
				versionTable := value.(*lua.LTable)
				if versionTable.RawGetString("version").Type() == lua.LTString {
					gameVersions = append(gameVersions, GameVersion{
						Version: versionTable.RawGetString("version").String(),
					})
				}
			}
		})
	}

	mod := &Mod{
		ID:           modTable.RawGetString("id").String(),
		Name:         modTable.RawGetString("name").String(),
		Description:  modTable.RawGetString("description").String(),
		Author:       modTable.RawGetString("author").String(),
		Version:      modTable.RawGetString("version").String(),
		Enabled:      bool(modTable.RawGetString("enabled").(lua.LBool)),
		Dependencies: nil,
		DownloadURL:  "",
		IconURL:      "",
		HeaderImage:  "",
		GameVersions: gameVersions,
		Image:        nil,
		GameID:       modTable.RawGetString("game_id").String(),
	}

	return mod, nil
}
