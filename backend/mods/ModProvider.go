package mods

import (
	"errors"
	lua "github.com/yuin/gopher-lua"
)

type GameVersion struct {
	Version    string `json:"version"`
	ModVersion string `json:"mod_version"`
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

func NewGameVersionsFromLuaTable(gameVersionTable *lua.LTable) ([]GameVersion, error) {
	if gameVersionTable == nil {
		return nil, nil // No game versions provided
	}
	if gameVersionTable.Type() != lua.LTTable {
		return nil, errors.New("gameVersionTable is not a Lua table")
	}

	var gameVersions []GameVersion
	gameVersionTable.ForEach(func(key lua.LValue, value lua.LValue) {
		if key.Type() == lua.LTString && value.Type() == lua.LTString {
			gameVersions = append(gameVersions, GameVersion{
				Version:    key.String(),
				ModVersion: value.String(),
			})
		}
	})

	return gameVersions, nil
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

	if modTable.RawGetString("game_id").Type() != lua.LTString {
		return nil, errors.New("modTable is missing 'game_id' field or it is not a string")
	}

	var gameVersions []GameVersion
	gameVersionsTable := modTable.RawGetString("game_versions")
	if gameVersionsTable.Type() == lua.LTTable {
		// game_versions = {[gameVersion] = modVersion}
		var err error
		gameVersions, err = NewGameVersionsFromLuaTable(gameVersionsTable.(*lua.LTable))
		if err != nil {
			return nil, err
		}
	}

	enabled := false
	if modTable.RawGetString("enabled").Type() == lua.LTBool {
		enabled = bool(modTable.RawGetString("enabled").(lua.LBool))
	}

	mod := &Mod{
		ID:           modTable.RawGetString("id").String(),
		Name:         modTable.RawGetString("name").String(),
		Description:  modTable.RawGetString("description").String(),
		Author:       modTable.RawGetString("author").String(),
		Version:      modTable.RawGetString("version").String(),
		Enabled:      enabled,
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
