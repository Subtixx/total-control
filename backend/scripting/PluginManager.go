package scripting

import (
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
	"os"
	"path/filepath"
	"strings"
)

var (
	ErrorAlreadyLoaded = fmt.Errorf("plugin already loaded")
	ErrorNotFound      = fmt.Errorf("plugin not found")
	ErrorInvalidPlugin = fmt.Errorf("invalid plugin format")
)

type PluginManager struct {
	Plugins map[string]*LuaPlugin
}

func (pm *PluginManager) GetLoadedPlugins() map[string]*LuaPlugin {
	return pm.Plugins
}

func (pm *PluginManager) Logger() *log.Entry {
	return log.WithFields(log.Fields{
		"prefix": "PM",
	})
}

func NewPluginManager(pluginDir string) (*PluginManager, error) {
	pm := &PluginManager{
		Plugins: make(map[string]*LuaPlugin),
	}
	err := pm.LoadPlugins(pluginDir)
	if err != nil {
		return nil, err
	}
	return pm, nil
}

func (pm *PluginManager) LoadPlugins(path string) error {
	// Iterate through the directory and load plugins
	if info, err := os.Stat(path); os.IsNotExist(err) || !info.IsDir() {
		return fmt.Errorf("plugin directory does not exist: %s", path)
	}
	files, err := os.ReadDir(path)
	if err != nil {
		return err
	}

	for _, file := range files {
		if !file.IsDir() &&
			!strings.HasSuffix(file.Name(), ".tcplugin") {
			continue // Skip non-directory files that are not tcplugin
		}

		if strings.HasPrefix(file.Name(), ".") {
			continue // Skip hidden files and directories
		}

		dirPath := filepath.Join(path, file.Name())
		if err := pm.loadPlugin(dirPath); !errors.Is(err, ErrorAlreadyLoaded) {
			if err != nil {
				return fmt.Errorf("failed to load plugin %s: %w", dirPath, err)
			}
		} else {
			pm.Logger().Warnf("Plugin %s is already loaded, skipping", file.Name())
		}
	}
	pm.Logger().Infof("Loaded %d plugins from '%s'", len(pm.Plugins), path)

	return nil
}

func (pm *PluginManager) loadPlugin(file string) error {
	luaPlugin, err := LoadPlugin(file)
	if err != nil {
		if errors.Is(err, ErrorNotFound) || errors.Is(err, ErrorInvalidPlugin) {
			return fmt.Errorf("failed to load plugin %s: %w", file, err)
		}
		return fmt.Errorf("error loading plugin %s: %w", file, err)
	}

	if _, exists := pm.Plugins[luaPlugin.Id.String()]; exists {
		return fmt.Errorf("%w: %s (%s)", ErrorAlreadyLoaded, luaPlugin.Name, luaPlugin.Id.String())
	}

	pm.Plugins[luaPlugin.Id.String()] = luaPlugin
	return nil
}

func (pm *PluginManager) Shutdown() {
	for _, plugin := range pm.Plugins {
		if err := plugin.Shutdown(); err != nil {
			pm.Logger().Errorf("Failed to shutdown plugin %s: %v", plugin.Name, err)
		}
	}
	pm.Plugins = make(map[string]*LuaPlugin) // Clear loaded plugins
	pm.Logger().Debug("All plugins have been shut down and cleared from memory.")
}
