package scripting

import (
	"TotalControl/backend/mods"
	"fmt"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	lua "github.com/yuin/gopher-lua"
)

// Deprecated: Use LuaPlugin instead
type LuaModProviderEngine struct {
	LuaEngine
	plugin *lua.LTable
}

// Deprecated: Use LuaPlugin instead
func NewLuaModProviderEngine(luaEngineId uuid.UUID) (*LuaModProviderEngine, error) {
	luaEngine := &LuaModProviderEngine{
		LuaEngine: LuaEngine{
			L:    lua.NewState(),
			uuid: luaEngineId,
		},
	}
	err := luaEngine.Setup()
	if err != nil {
		log.Errorf("Failed to setup Lua engine: %v", err)
		return nil, err
	}
	return luaEngine, nil
}

// Deprecated: Use LuaPlugin instead
func (p *LuaModProviderEngine) GetPlugin() (*lua.LTable, error) {
	if p.L == nil {
		return nil, fmt.Errorf("lua state is not initialized")
	}
	plugin := p.L.GetGlobal("plugin")
	if plugin == lua.LNil {
		return nil, fmt.Errorf("plugin global is not set in Lua state")
	}

	if plugin.Type() != lua.LTTable {
		return nil, fmt.Errorf("plugin global is not a Lua table, got %s", plugin.Type().String())
	}

	if !p.IsValid() {
		return nil, fmt.Errorf("plugin is not valid, missing required methods")
	}
	return plugin.(*lua.LTable), nil
}

// Deprecated: Use LuaPlugin instead
func (p *LuaModProviderEngine) IsValid() bool {
	plugin := p.L.GetGlobal("plugin")
	if plugin.Type() != lua.LTTable {
		log.Fatal("plugin global is not a Lua table, got " + plugin.Type().String())
		return false
	}

	pluginMethods := []string{
		"GetInstalledMods", "GetInstalledModByID",
		"GetMods", "GetModByID",
		"AddMod", "RemoveMod", "UpdateMod", "GetGameModDirectory", "GetGameID",
	}
	for _, method := range pluginMethods {
		if !p.HasMethod(plugin, method) {
			log.Fatal(fmt.Sprintf("Method %s not found in plugin table", method))
			return false
		}
	}
	return true
}

// Deprecated: Use LuaPlugin instead
func (p *LuaModProviderEngine) GetInstalledMods() ([]mods.Mod, error) {
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

	val, err := p.Call(plugin, "GetInstalledMods")
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
		mod, err := mods.NewModFromLuaTable(modTable)
		if err != nil {
			log.Errorf("failed to create mod from Lua table: %v", err)
			return // Skip this mod if it cannot be created
		}
		foundMods = append(foundMods, *mod)
	})

	return foundMods, nil
}

// Deprecated: Use LuaPlugin instead
func (p *LuaModProviderEngine) GetModByID(id string) (*mods.Mod, error) {
	plugin := p.L.GetGlobal("plugin")
	if plugin.Type() != lua.LTTable {
		return nil, fmt.Errorf("plugin global is not a Lua table, got %s", plugin.Type().String())
	}

	luaID := lua.LString(id)
	val, err := p.Call(plugin, "GetModByID", luaID)
	if err != nil {
		return nil, fmt.Errorf("failed to get mod by ID: %v", err)
	}

	if val.Type() != lua.LTTable {
		return nil, fmt.Errorf("expected Lua table for mod, got %s", val.Type().String())
	}

	modTable := val.(*lua.LTable)
	mod, err := mods.NewModFromLuaTable(modTable)
	if err != nil {
		return nil, fmt.Errorf("failed to create mod from Lua table: %v", err)
	}
	return mod, nil
}

// Deprecated: Use LuaPlugin instead
func (p *LuaModProviderEngine) GetGameModDirectory() (string, error) {
	plugin, err := p.GetPlugin()
	if err != nil {
		return "", fmt.Errorf("failed to get plugin: %v", err)
	}

	val, err := p.Call(plugin, "GetGameModDirectory")
	if err != nil {
		return "", fmt.Errorf("failed to get game mod directory: %v", err)
	}

	if val.Type() != lua.LTString {
		return "", fmt.Errorf("expected Lua string for game mod directory, got %s", val.Type().String())
	}

	return val.String(), nil
}
