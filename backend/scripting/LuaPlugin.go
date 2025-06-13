package scripting

import (
	"TotalControl/backend/mods"
	"TotalControl/backend/utils"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	lua "github.com/yuin/gopher-lua"
	"path/filepath"
)

type PluginInfo struct {
	Id         uuid.UUID `json:"id"`
	Name       string    `json:"name"`
	Version    string    `json:"version"`
	EntryPoint string    `json:"entry"`
}

type Plugin struct {
	PluginInfo

	PluginDir string `json:"-"`
}

type LuaPlugin struct {
	Plugin
	LuaEngine
	// -------------------------------
	plugin *lua.LTable
	// -------------------------------
	getMods             *lua.LFunction
	getInstalledMods    *lua.LFunction
	getModByID          *lua.LFunction
	addMod              *lua.LFunction
	removeMod           *lua.LFunction
	updateMod           *lua.LFunction
	getGameModDirectory *lua.LFunction
	getGameID           *lua.LFunction
}

func (p *LuaPlugin) Initialize() error {
	if p.plugin == nil {
		return fmt.Errorf("plugin table is not initialized")
	}

	p.getMods = p.L.GetField(p.plugin, "GetMods").(*lua.LFunction)
	if p.getMods == nil {
		return fmt.Errorf("GetMods function not found in plugin table")
	}

	p.getInstalledMods = p.L.GetField(p.plugin, "GetInstalledMods").(*lua.LFunction)
	if p.getInstalledMods == nil {
		return fmt.Errorf("GetInstalledMods function not found in plugin table")
	}

	p.getModByID = p.L.GetField(p.plugin, "GetModByID").(*lua.LFunction)
	if p.getModByID == nil {
		return fmt.Errorf("GetModByID function not found in plugin table")
	}

	p.addMod = p.L.GetField(p.plugin, "AddMod").(*lua.LFunction)
	if p.addMod == nil {
		return fmt.Errorf("AddMod function not found in plugin table")
	}

	p.removeMod = p.L.GetField(p.plugin, "RemoveMod").(*lua.LFunction)
	if p.removeMod == nil {
		return fmt.Errorf("RemoveMod function not found in plugin table")
	}

	p.updateMod = p.L.GetField(p.plugin, "UpdateMod").(*lua.LFunction)
	if p.updateMod == nil {
		return fmt.Errorf("UpdateMod function not found in plugin table")
	}

	p.getGameModDirectory = p.L.GetField(p.plugin, "GetGameModDirectory").(*lua.LFunction)
	if p.getGameModDirectory == nil {
		return fmt.Errorf("GetGameModDirectory function not found in plugin table")
	}

	p.getGameID = p.L.GetField(p.plugin, "GetGameID").(*lua.LFunction)
	if p.getGameID == nil {
		return fmt.Errorf("GetGameID function not found in plugin table")
	}

	log.Debugf("Lua plugin %s initialized with ID %s", p.Name, p.Id.String())
	return nil
}

func (p *LuaPlugin) GetMods() (map[string]interface{}, error) {
	if p.getMods == nil {
		return nil, fmt.Errorf("GetMods function is not initialized")
	}

	callFunc, err := p.CallFunc(p.plugin, p.getMods)
	if err != nil {
		return nil, err
	}
	if callFunc.Type() != lua.LTTable {
		return nil, fmt.Errorf("expected Lua table for mods, got %s", callFunc.Type().String())
	}
	modsTable := callFunc.(*lua.LTable)
	foundMods := make(map[string]interface{})
	modsTable.ForEach(func(key lua.LValue, value lua.LValue) {
		if value.Type() != lua.LTTable {
			log.Warnf("Unexpected value type %s for mod %s, expected table", value.Type().String(), key.String())
			return
		}
		value.(*lua.LTable).ForEach(func(k lua.LValue, v lua.LValue) {
			log.Debugf("Mod %s.%s = %s", key.String(), k.String(), v.String())
		})
		mod, err := mods.NewModFromLuaTable(value.(*lua.LTable))
		if err != nil {
			log.Warnf("Failed to create mod from Lua table for key %s: %v", key.String(), err)
			return
		}
		foundMods[key.String()] = mod
	})
	log.Debugf("Retrieved %d mods from Lua plugin %s", len(foundMods), p.Name)
	return foundMods, nil
}

// LoadLuaPluginFromZip Loads a plugin using a zip with the custom extension ".tcplugin".
func LoadLuaPluginFromZip(pluginZipPath string) (*LuaPlugin, error) {
	files, err := utils.ReadFilesFromZip(pluginZipPath)
	if err != nil {
		return nil, err
	}

	var plugin LuaPlugin
	if err := json.Unmarshal(files["info.json"], &plugin); err != nil {
		return nil, fmt.Errorf("failed to unmarshal plugin info: %w", err)
	}

	if plugin.Id == uuid.Nil {
		return nil, fmt.Errorf("plugin ID is not set or is invalid")
	}

	if plugin.EntryPoint == "" {
		return nil, fmt.Errorf("plugin entry point is not set")
	}

	plugin.PluginDir = pluginZipPath
	plugin.LuaEngine = LuaEngine{
		L:    lua.NewState(),
		uuid: plugin.Id,
	}

	if err := plugin.Setup(); err != nil {
		return nil, fmt.Errorf("failed to setup Lua plugin: %w", err)
	}

	scriptFile, ok := files[plugin.EntryPoint]
	if !ok {
		return nil, fmt.Errorf("plugin entry point %s not found in zip", plugin.EntryPoint)
	}
	luaPlugin, err := loadPluginScript(plugin.L, string(scriptFile))
	if err != nil {
		return nil, fmt.Errorf("failed to load Lua plugin script: %w", err)
	}

	plugin.plugin = luaPlugin

	if err := plugin.Initialize(); err != nil {
		return nil, fmt.Errorf("failed to initialize Lua plugin: %w", err)
	}

	return &plugin, nil
}

func LoadLuaPlugin(pluginDir string) (*LuaPlugin, error) {
	infoFile := filepath.Join(pluginDir, "info.json")
	pluginInfo, err := utils.ReadFile(infoFile)
	if err != nil {
		return nil, err
	}

	var plugin LuaPlugin
	if err := json.Unmarshal(pluginInfo, &plugin); err != nil {
		return nil, err
	}

	if plugin.Id == uuid.Nil {
		return nil, fmt.Errorf("plugin ID is not set or is invalid")
	}

	if plugin.EntryPoint == "" || utils.FileExists(filepath.Join(pluginDir, plugin.EntryPoint)) == false {
		return nil, fmt.Errorf("plugin entry point is not set or does not exist")
	}
	plugin.PluginDir = pluginDir

	plugin.LuaEngine = LuaEngine{
		L:    lua.NewState(),
		uuid: plugin.Id,
	}

	if err := plugin.Setup(); err != nil {
		return nil, err
	}

	scriptPath := filepath.Join(pluginDir, plugin.EntryPoint)
	luaPlugin, err := loadPluginScriptFile(plugin.L, scriptPath)
	if err != nil {
		return nil, err
	}
	plugin.plugin = luaPlugin

	return &plugin, nil
}

func loadPluginScript(l *lua.LState, scriptContent string) (*lua.LTable, error) {
	if err := l.DoString(scriptContent); err != nil {
		return nil, fmt.Errorf("failed to load Lua script: %w", err)
	}

	val := l.Get(-1)
	l.Pop(1)
	if val.Type() != lua.LTTable {
		return nil, fmt.Errorf("expected Lua table for plugin, got %s", val.Type().String())
	}

	return val.(*lua.LTable), nil
}

func loadPluginScriptFile(l *lua.LState, scriptPath string) (*lua.LTable, error) {
	if err := l.DoFile(scriptPath); err != nil {
		return nil, fmt.Errorf("failed to load Lua script %s: %w", scriptPath, err)
	}

	val := l.Get(-1)
	l.Pop(1)
	if val.Type() != lua.LTTable {
		return nil, fmt.Errorf("expected Lua table for plugin, got %s", val.Type().String())
	}

	return val.(*lua.LTable), nil
}
