package scripting

import (
	"TotalControl/backend/mods"
	"fmt"
	log "github.com/sirupsen/logrus"
	lua "github.com/yuin/gopher-lua"
)

type LuaModProviderEngine struct {
	LuaEngine
}

func NewLuaModProviderEngine() *LuaModProviderEngine {
	luaEngine := &LuaModProviderEngine{
		LuaEngine: LuaEngine{
			L: lua.NewState(),
		},
	}
	luaEngine.Setup()
	return luaEngine
}

func (p *LuaModProviderEngine) IsValid() bool {
	plugin := p.L.GetGlobal("plugin")
	if plugin.Type() != lua.LTTable {
		log.Fatal("plugin global is not a Lua table, got " + plugin.Type().String())
		return false
	}

	pluginMethods := []string{
		"GetMods", "GetModByID", "AddMod", "RemoveMod", "UpdateMod", "GetGameModDirectory", "GetGameID",
	}
	for _, method := range pluginMethods {
		if !p.HasMethod(plugin, method) {
			log.Fatal(fmt.Sprintf("Method %s not found in plugin table", method))
			return false
		}
	}
	return true
}

func (p *LuaModProviderEngine) GetMods() ([]mods.Mod, error) {
	plugin := p.L.GetGlobal("plugin")
	if plugin.Type() != lua.LTTable {
		return nil, fmt.Errorf("plugin global is not a Lua table, got %s", plugin.Type().String())
	}

	luaGameId, err := p.Call(plugin, "GetGameID")
	if err != nil {
		return nil, fmt.Errorf("failed to get game ID: %v", err)
	}
	if luaGameId.Type() != lua.LTString {
		return nil, fmt.Errorf("expected Lua string for game ID, got %s", luaGameId.Type().String())
	}

	val, err := p.Call(plugin, "GetMods")
	if err != nil {
		return nil, fmt.Errorf("failed to get mods: %v", err)
	}

	if val.Type() != lua.LTTable {
		return nil, fmt.Errorf("expected Lua table, got %s", val.Type().String())
	}

	var foundMods []mods.Mod
	luaTable := val.(*lua.LTable)
	luaTable.ForEach(func(_, value lua.LValue) {
		if value.Type() != lua.LTTable {
			return // Skip non-table values
		}

		modTable := value.(*lua.LTable)
		if modTable.RawGetString("id").Type() != lua.LTString ||
			modTable.RawGetString("name").Type() != lua.LTString ||
			modTable.RawGetString("enabled").Type() != lua.LTBool ||
			modTable.RawGetString("author").Type() != lua.LTString ||
			modTable.RawGetString("version").Type() != lua.LTString ||
			modTable.RawGetString("game_versions").Type() != lua.LTTable {
			log.Warn("Mod table is missing required fields or fields are not strings")
			return // Skip if any required field is missing or not a string
		}

		var gameVersions []mods.GameVersion
		gameVersionsTable := modTable.RawGetString("game_versions")
		if gameVersionsTable.Type() == lua.LTTable {
			gameVersionsTable.(*lua.LTable).ForEach(func(_, value lua.LValue) {
				if value.Type() == lua.LTTable {
					versionTable := value.(*lua.LTable)
					if versionTable.RawGetString("version").Type() == lua.LTString {
						gameVersions = append(gameVersions, mods.GameVersion{
							Version: versionTable.RawGetString("version").String(),
						})
					} else {
						log.Warn("Game version table is missing 'version' field or it is not a string")
					}
				} else {
					log.Warn("Game versions entry is not a table")
				}
			})
		} else {
			log.Warn("Game versions is not a Lua table")
		}

		mod := mods.Mod{
			ID:           modTable.RawGetString("id").String(),
			Name:         modTable.RawGetString("name").String(),
			Author:       modTable.RawGetString("author").String(),
			Version:      modTable.RawGetString("version").String(),
			Enabled:      bool(modTable.RawGetString("enabled").(lua.LBool)),
			GameVersions: gameVersions,
			GameID:       gameID,
		}
		foundMods = append(foundMods, mod)
	})

	return foundMods, nil
}

func (p *LuaModProviderEngine) GetModByID(id string) (*mods.Mod, error) {
	plugin := p.L.GetGlobal("plugin")
	if plugin.Type() != lua.LTTable {
		return nil, fmt.Errorf("plugin global is not a Lua table, got %s", plugin.Type().String())
	}

	val, err := p.Call(plugin, "GetModByID", id)
	if err != nil {
		return nil, fmt.Errorf("failed to get mod by ID: %v", err)
	}

	if val.Type() != lua.LTTable {
		return nil, fmt.Errorf("expected Lua table for mod, got %s", val.Type().String())
	}

	modTable := val.(*lua.LTable)
	if modTable.RawGetString("id").Type() != lua.LTString ||
		modTable.RawGetString("name").Type() != lua.LTString ||
		modTable.RawGetString("enabled").Type() != lua.LTBool ||
		modTable.RawGetString("author").Type() != lua.LTString ||
		modTable.RawGetString("version").Type() != lua.LTString ||
		modTable.RawGetString("game_versions").Type() != lua.LTTable {
		return nil, fmt.Errorf("mod table is missing required fields or fields are not strings")
	}

	var gameVersions []mods.GameVersion
	gameVersionsTable := modTable.RawGetString("game_versions")
	if gameVersionsTable.Type() == lua.LTTable {
		gameVersionsTable.(*lua.LTable).ForEach(func(_, value lua.LValue) {
			if value.Type() == lua.LTTable {
				versionTable := value.(*lua.LTable)
				if versionTable.RawGetString("version").Type() == lua.LTString {
					gameVersions = append(gameVersions, mods.GameVersion{
						Version: versionTable.RawGetString("version").String(),
					})
				} else {
					log.Warn("Game version table is missing 'version' field or it is not a string")
				}
			} else {
				log.Warn("Game versions entry is not a table")
			}
		})
	} else {
		log.Warn("Game versions is not a Lua table")
	}

	mod := &mods.Mod{
		ID:           modTable.RawGetString("id").String(),
		Name:         modTable.RawGetString("name").String(),
		Description:  modTable.RawGetString("description").String(),
		Author:       modTable.RawGetString("author").String(),
		Version:      modTable.RawGetString("version").String(),
		Enabled:      bool(modTable.RawGetString("enabled").(lua.LBool)),
		Dependencies: nil, // Assuming dependencies are not handled in this Lua mod provider
		DownloadURL:  "",
		IconURL:      "",
		HeaderImage:  "",
		GameVersions: gameVersions,
		Image:        nil,
		GameID:       "",
	}
	if modTable.RawGetString("game_id").Type() == lua.LTString {
		mod.GameID = modTable.RawGetString("game_id").String()
	} else {
		log.Warn("Game ID field is missing or not a string, using empty string")
	}
	if modTable.RawGetString("image").Type() == lua.LTString {
		mod.Image = &mods.Image{
			URL: modTable.RawGetString("image").String(),
		}
	} else {
		log.Warn("Image field is missing or not a string, using nil")
	}
	if modTable.RawGetString("download_url").Type() == lua.LTString {
		mod.DownloadURL = modTable.RawGetString("download_url").String()
	} else {
		log.Warn("Download URL field is missing or not a string, using empty string")
	}
	if modTable.RawGetString("header_image").Type() == lua.LTString {
		mod.HeaderImage = modTable.RawGetString("header_image").String()
	} else {
		log.Warn("Header image field is missing or not a string, using empty string")
	}
	if modTable.RawGetString("dependencies").Type() == lua.LTTable {
		dependenciesTable := modTable.RawGetString("dependencies").(*lua.LTable)
		dependenciesTable.ForEach(func(_, depValue lua.LValue) {
			if depValue.Type() == lua.LTString {
				mod.Dependencies = append(mod.Dependencies, depValue.String())
			} else {
				log.Warn("Dependency entry is not a string, skipping")
			}
		})
	} else {
		log.Warn("Dependencies field is missing or not a table, using empty slice")
	}
	return mod, nil
}
