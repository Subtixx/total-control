package scripting

import (
	"TotalControl/backend/mods"
	"TotalControl/backend/plugins"
	"TotalControl/backend/utils"
	"context"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	lua "github.com/yuin/gopher-lua"
	"net/http"
	"os"
	"path/filepath"
	"runtime/debug"
	"strings"
	"time"
)

var (
	ErrPluginNotFound = fmt.Sprintf("LuaPlugin not found in context")
)

type LuaPlugin struct {
	plugins.Plugin
	LuaEngine

	httpClient *http.Client
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

func LoadPlugin(filePath string) (*LuaPlugin, error) {
	if info, err := os.Stat(filePath); err == nil && info.IsDir() {
		// Check if the directory contains an info.json file
		infoJson := filepath.Join(filePath, "info.json")
		if _, err := os.Stat(infoJson); err != nil {
			return nil, fmt.Errorf("info.json not found in plugin directory: %s", filePath)
		}
		luaPlugin, err := LoadLuaPlugin(filePath)
		if err != nil {
			return nil, err
		}
		return luaPlugin, nil
	} else if strings.HasSuffix(filePath, ".tcplugin") {
		luaPlugin, err := LoadLuaPluginFromZip(filePath)
		if err != nil {
			return nil, fmt.Errorf("failed to load Lua plugin from zip: %w", err)
		}
		if luaPlugin == nil {
			return nil, fmt.Errorf("failed to load Lua plugin from zip: plugin is nil")
		}
		return luaPlugin, nil
	}

	return nil, fmt.Errorf("invalid plugin path: %s", filePath)
}

// LoadLuaPluginFromZip Loads a plugin using a zip with the custom extension ".tcplugin".
func LoadLuaPluginFromZip(pluginZipPath string) (*LuaPlugin, error) {
	files, err := utils.ReadFilesFromZip(pluginZipPath)
	if err != nil {
		return nil, err
	}

	var plugin *LuaPlugin
	if err := json.Unmarshal(files["info.json"], &plugin); err != nil {
		return nil, fmt.Errorf("failed to unmarshal plugin info: %w", err)
	}

	if plugin.Id == uuid.Nil {
		return nil, fmt.Errorf("plugin ID is not set or is invalid")
	}

	if plugin.EntryPoint == "" {
		return nil, fmt.Errorf("plugin entry point is not set")
	}

	plugin.httpClient = &http.Client{
		Timeout: 5 * time.Second,
	}
	plugin.PluginDir = pluginZipPath
	plugin.IsPacked = true
	plugin.LuaEngine = LuaEngine{
		L:    lua.NewState(),
		uuid: plugin.Id,
	}
	luaEngineOptions := LuaEngineSetupOptions{
		Context:       context.Background(),
		ContextValues: map[string]interface{}{"plugin": plugin},
	}

	if err := plugin.Setup(luaEngineOptions); err != nil {
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
		plugin.Logger().Errorf("Failed to initialize Lua plugin: %v", err)
		plugin.L.Close()
		return nil, fmt.Errorf("failed to initialize Lua plugin: %w", err)
	}

	return plugin, nil
}

func LoadLuaPlugin(pluginDir string) (*LuaPlugin, error) {
	infoFile := filepath.Join(pluginDir, "info.json")
	pluginInfo, err := utils.ReadFile(infoFile)
	if err != nil {
		return nil, err
	}

	var plugin *LuaPlugin
	if err := json.Unmarshal(pluginInfo, &plugin); err != nil {
		return nil, err
	}

	if plugin.Id == uuid.Nil {
		return nil, fmt.Errorf("plugin ID is not set or is invalid")
	}

	if plugin.EntryPoint == "" || utils.FileExists(filepath.Join(pluginDir, plugin.EntryPoint)) == false {
		return nil, fmt.Errorf("plugin entry point is not set or does not exist")
	}

	plugin.httpClient = &http.Client{
		Timeout: 5 * time.Second,
	}
	plugin.PluginDir = pluginDir
	plugin.IsPacked = false
	plugin.LuaEngine = LuaEngine{
		L:    lua.NewState(),
		uuid: plugin.Id,
	}
	luaEngineOptions := LuaEngineSetupOptions{
		Context:       context.Background(),
		ContextValues: map[string]interface{}{"plugin": plugin},
	}

	if err := plugin.Setup(luaEngineOptions); err != nil {
		return nil, err
	}

	scriptPath := filepath.Join(pluginDir, plugin.EntryPoint)
	luaPlugin, err := loadPluginScriptFile(plugin.L, scriptPath)
	if err != nil {
		return nil, err
	}
	plugin.plugin = luaPlugin

	return plugin, nil
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

func (p *LuaPlugin) Logger() *log.Entry {
	return log.WithFields(log.Fields{
		"lua":    true,
		"plugin": p.Id.String(),
	})
}

// ToString returns a string representation of the LuaPlugin.
func (p *LuaPlugin) String() string {
	return "LuaPlugin{" +
		"\n\t" + p.Plugin.String() +
		"\n\t" + p.LuaEngine.String() +
		"\n}"
}

func (p *LuaPlugin) Initialize() error {
	if p.plugin == nil {
		p.Logger().Error("plugin table is not initialized")
		return fmt.Errorf("plugin table is not initialized")
	}

	p.getMods = p.L.GetField(p.plugin, "GetMods").(*lua.LFunction)
	if p.getMods == nil {
		p.Logger().Error("GetMods function not found in plugin table")
		return fmt.Errorf("GetMods function not found in plugin table")
	}

	p.getInstalledMods = p.L.GetField(p.plugin, "GetInstalledMods").(*lua.LFunction)
	if p.getInstalledMods == nil {
		p.Logger().Error("GetInstalledMods function not found in plugin table")
		return fmt.Errorf("GetInstalledMods function not found in plugin table")
	}

	p.getModByID = p.L.GetField(p.plugin, "GetModByID").(*lua.LFunction)
	if p.getModByID == nil {
		p.Logger().Error("GetModByID function not found in plugin table")
		return fmt.Errorf("GetModByID function not found in plugin table")
	}

	p.addMod = p.L.GetField(p.plugin, "AddMod").(*lua.LFunction)
	if p.addMod == nil {
		p.Logger().Error("AddMod function not found in plugin table")
		return fmt.Errorf("AddMod function not found in plugin table")
	}

	p.removeMod = p.L.GetField(p.plugin, "RemoveMod").(*lua.LFunction)
	if p.removeMod == nil {
		p.Logger().Error("RemoveMod function not found in plugin table")
		return fmt.Errorf("RemoveMod function not found in plugin table")
	}

	p.updateMod = p.L.GetField(p.plugin, "UpdateMod").(*lua.LFunction)
	if p.updateMod == nil {
		p.Logger().Error("UpdateMod function not found in plugin table")
		return fmt.Errorf("UpdateMod function not found in plugin table")
	}

	p.getGameModDirectory = p.L.GetField(p.plugin, "GetGameModDirectory").(*lua.LFunction)
	if p.getGameModDirectory == nil {
		p.Logger().Error("GetGameModDirectory function not found in plugin table")
		return fmt.Errorf("GetGameModDirectory function not found in plugin table")
	}

	p.getGameID = p.L.GetField(p.plugin, "GetGameID").(*lua.LFunction)
	if p.getGameID == nil {
		p.Logger().Error("GetGameID function not found in plugin table")
		return fmt.Errorf("GetGameID function not found in plugin table")
	}

	p.Logger().Debugf("Lua plugin %s initialized with ID %s", p.Name, p.Id.String())
	return nil
}

func (p *LuaPlugin) Shutdown() error {
	err := p.LuaEngine.Shutdown()
	if err != nil {
		return fmt.Errorf("failed to shutdown Lua engine: %w", err)
	}
	err = p.Plugin.Shutdown()
	if err != nil {
		return fmt.Errorf("failed to shutdown plugin: %w", err)
	}

	p.Logger().Infof("Lua plugin %s with ID %s has been shut down", p.Name, p.Id.String())
	return nil
}

func (p *LuaPlugin) GetMods() (map[string]interface{}, error) {
	if p.getMods == nil {
		p.Logger().Error("GetMods function is not initialized")
		return nil, fmt.Errorf("GetMods function is not initialized")
	}

	callFunc, err := p.CallFunc(p.plugin, p.getMods)
	if err != nil {
		return nil, err
	}
	if callFunc.Type() != lua.LTTable {
		p.Logger().Errorf("expected Lua table for mods, got %s", callFunc.Type().String())
		return nil, fmt.Errorf("expected Lua table for mods, got %s", callFunc.Type().String())
	}
	modsTable := callFunc.(*lua.LTable)
	foundMods := make(map[string]interface{})
	modsTable.ForEach(func(key lua.LValue, value lua.LValue) {
		if value.Type() != lua.LTTable {
			p.Logger().Warnf("Unexpected value type %s for mod %s, expected table", value.Type().String(), key.String())
			return
		}
		value.(*lua.LTable).ForEach(func(k lua.LValue, v lua.LValue) {
			p.Logger().Debugf("Mod %s.%s = %s", key.String(), k.String(), v.String())
		})
		mod, err := mods.NewModFromLuaTable(value.(*lua.LTable))
		if err != nil {
			p.Logger().Warnf("Failed to create mod from Lua table for key %s: %v", key.String(), err)
			return
		}
		foundMods[key.String()] = mod
	})
	p.Logger().Debugf("Retrieved %d mods from Lua plugin %s", len(foundMods), p.Name)
	return foundMods, nil
}

func (p *LuaPlugin) GetHttpClient() *http.Client {
	if p.httpClient == nil {
		p.httpClient = &http.Client{
			Timeout: 5 * time.Second,
		}
	}
	return p.httpClient
}

func GetLuaHttpClient(L *lua.LState) *http.Client {
	if L == nil {
		log.Error("Lua state is nil, cannot get HTTP client")
		return nil
	}

	luaPlugin := GetLuaPlugin(L)
	if luaPlugin != nil {
		return luaPlugin.GetHttpClient()
	}
	return http.DefaultClient
}

func GetLuaPlugin(L *lua.LState) *LuaPlugin {
	if L == nil {
		log.Error("Lua state is nil, cannot get LuaPlugin")
		return nil
	}

	if L.Context() == nil {
		return nil
	}

	if v := L.Context().Value("plugin"); v != nil {
		if plugin, ok := v.(*LuaPlugin); ok {
			return plugin
		}
		log.Errorf("context value 'plugin' is not a LuaPlugin, got %T", v)
		debug.PrintStack()
	} else {
		log.Error("context value 'plugin' is nil")
	}
	return nil
}
