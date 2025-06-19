package main

import (
	"TotalControl/backend/games"
	"TotalControl/backend/plugins"
	"TotalControl/backend/scripting"
	"context"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/runtime"
	"strings"
)

// App struct
type App struct {
	ctx context.Context
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called at application startup
func (a *App) startup(ctx context.Context) {
	// Perform your setup here
	a.ctx = ctx
}

func (a *App) onSecondInstanceLaunch(secondInstanceData options.SecondInstanceData) {
	secondInstanceArgs := secondInstanceData.Args

	println("user opened second instance", strings.Join(secondInstanceData.Args, ","))
	println("user opened second from", secondInstanceData.WorkingDirectory)
	runtime.WindowUnminimise(a.ctx)
	runtime.Show(a.ctx)
	go runtime.EventsEmit(a.ctx, "launchArgs", secondInstanceArgs)
}

// domReady is called after the front-end dom has been loaded
func (a *App) domReady(ctx context.Context) {
	// Add your action here
}

// beforeClose is called when the application is about to quit,
// either by clicking the window close button or calling runtime.Quit.
// Returning true will cause the application to continue,
// false will continue shutdown as normal.
func (a *App) beforeClose(ctx context.Context) (prevent bool) {
	return false
}

// shutdown is called at application termination
func (a *App) shutdown(ctx context.Context) {
	// Perform your teardown here
	// 在此处做一些资源释放的操作
}

func (a *App) GetInstalledPlugins() []*scripting.LuaPlugin {
	pluginManager := GetPluginManager()
	if pluginManager == nil {
		return nil
	}

	loadedPlugins := pluginManager.GetLoadedPlugins()
	if loadedPlugins == nil {
		return nil
	}
	installedPlugins := make([]*scripting.LuaPlugin, 0, len(loadedPlugins))
	for _, plugin := range loadedPlugins {
		if plugin == nil {
			continue
		}
		installedPlugins = append(installedPlugins, plugin)
	}
	return installedPlugins
}

func (a *App) GetAvailablePlugins() []*plugins.PluginRepositoryInfo {
	pluginManager := GetPluginManager()
	if pluginManager == nil {
		return nil
	}

	pluginRepository := pluginManager.GetPluginRepository(plugins.DefaultPluginRepositoryId)
	if pluginRepository == nil {
		return nil
	}
	availablePlugins := make([]*plugins.PluginRepositoryInfo, 0, len(pluginRepository.Plugins))
	for _, plugin := range pluginRepository.Plugins {
		if plugin == nil {
			continue
		}
		availablePlugins = append(availablePlugins, plugin)
	}
	return availablePlugins
}

func (a *App) GetPluginRepositories() []*plugins.PluginRepository {
	pluginManager := GetPluginManager()
	if pluginManager == nil {
		return nil
	}

	return pluginManager.GetPluginRepositories()
}

func (a *App) GetInstalledGames() []*games.Game {
	return make([]*games.Game, 0)
}
