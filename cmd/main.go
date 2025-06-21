package main

import (
	"TotalControl/backend/scripting"
	"TotalControl/backend/steam"
	"TotalControl/backend/utils"
	"flag"
	log "github.com/sirupsen/logrus"
	"os"
)

var gamePath *string
var logPath *string
var logLevel *string

func app() {
	log.Info("Starting TotalControl...")
	// Initialize the plugin manager
	pluginManager, err := scripting.NewPluginManager("plugins")
	if err != nil {
		log.Fatalf("Failed to initialize plugin manager: %v", err)
	}

	// Iterate through loaded plugins and print their details
	for _, plugin := range pluginManager.GetLoadedPlugins() {
		log.Info("\n" + plugin.PluginInfo.String())
		isInstalled, err := plugin.DetectGameInstallation(*gamePath)
		if err != nil {
			log.Errorf("Failed to detect game installation: %v", err)
		}

		if isInstalled {
			log.Infof("Detected game installation for plugin: %s at %s", plugin.Name, *gamePath)
		} else {
			log.Warnf("Game installation not found for plugin: %s", plugin.Name)
		}
	}
	pluginManager.Shutdown()

	steamLibraries, err := steam.GetSteamLibraries()
	games := steam.GetInstalledGames(steamLibraries)
	log.Infof("Found %d installed games", len(games))
	for _, game := range games {
		log.Info(game.String())
	}

	/*
			plugin, err := scripting.LoadLuaPlugin("plugins/factorio")
			if err != nil {
				log.Fatalf("Failed to load Lua plugin: %v", err)
			}

		plugin, err := scripting.LoadLuaPluginFromZip("plugins/Factorio.tcplugin")
		if err != nil {
			log.Fatalf("Failed to load Lua plugin from zip: %v", err)
		}

		log.Infof("Loaded Lua plugin: ID=%s, Name=%s, Version=%s, EntryPoint=%s, PluginDir=%s",
			plugin.Id, plugin.Name, plugin.Version, plugin.EntryPoint, plugin.PluginDir)

		modsAvailable, err := plugin.GetMods()
		if err != nil {
			log.Fatalf("Failed to get mods: %v", err)
		}
		if len(modsAvailable) == 0 {
			log.Fatal("No mods found in Lua plugin")
		}
		for k, mod := range modsAvailable {
			log.Debugf("Mod %s: %+v", k, mod)
		}

		plugin.Shutdown()
	*/

	/*
		// Test factorio lua
		luaEngine, err := scripting.NewLuaModProviderEngine(uuid.New())
		if err != nil {
			log.Fatalf("Failed to create Lua mod provider engine: %v", err)
		}
		if err := luaEngine.LoadFile("plugins/factorio/plugin.lua"); err != nil {
			log.Fatalf("Failed to load Lua file: %v", err)
		}

		if !luaEngine.IsValid() {
			log.Fatal("Lua mod provider engine is not valid")
		}
		log.Info("Lua mod provider engine is valid")

		foundMods, err := luaEngine.GetInstalledMods()
		if err != nil {
			log.Fatalf("Failed to get mods: %v", err)
		}
		if len(foundMods) == 0 {
			log.Fatal("No mods found in Lua script")
		}
		for _, mod := range foundMods {
			log.Infof("Found mod: ID=%s, Name=%s, Author=%s, Version=%s, Enabled=%t, GameVersions=%v",
				mod.ID, mod.Name, mod.Author, mod.Version, mod.Enabled, mod.GameVersions)
		}

		modPath, err := luaEngine.GetGameModDirectory()
		if err != nil {
			log.Fatalf("Failed to get game mod directory: %v", err)
			return
		}
		log.Infof("Game mod directory: %s", modPath)
	*/
}

func packPlugin(pluginDir *string) {
	plugin, err := scripting.LoadLuaPlugin(*pluginDir)
	if err != nil {
		log.Fatalf("Failed to load Lua plugin: %v", err)
	}
	if err := plugin.Pack(); err != nil {
		log.Fatalf("Failed to pack plugin: %v", err)
	}
	log.Infof("Packed plugin '%s' into '%s.tcplugin'", plugin.Name, plugin.Name)
}

func main() {
	logPath = flag.String("log-path", "", "Path to log file (default: stdout)")
	logLevel = flag.String("log-level", "debug", "Log level (debug, info, warn, error, fatal, panic)")
	pluginDir := flag.String("plugin-dir", "", "A directory to load a plugin from")
	pluginPack := flag.Bool("plugin-pack", false, "Pack the plugin into a .tcplugin file")

	gamePath = flag.String("game-path", "", "Path to the game installation")

	flag.Parse()
	if *logPath != "" {
		f, err := os.OpenFile(*logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
		if err != nil {
			log.Fatalf("Failed to open log file: %v", err)
		}
		log.SetOutput(f)
	}

	level, err := log.ParseLevel(*logLevel)
	if err != nil {
		log.Warnf("Invalid log level '%s', defaulting to 'info'", *logLevel)
		level = log.InfoLevel
	}
	log.SetLevel(level)
	log.SetReportCaller(true)
	log.SetFormatter(&utils.CustomFormatter{})

	if gamePath == nil || *gamePath == "" {
		log.Fatal("Game path is required")
	}

	if *pluginDir != "" {
		log.Infof("Loading plugin from directory: %s", *pluginDir)
		if *pluginPack {
			packPlugin(pluginDir)
			return
		}
		return
	}

	app()
}

/*
	provider := factorio.FactorioModProvider{}
	factorioMods, err := provider.GetInstalledMods()
	if err != nil {
		panic(err)
	}
	log.Infof("Found %d Factorio mods in directory: %s", len(factorioMods), provider.GetGameModDirectory())
	for _, mod := range factorioMods {
		versions := make([]string, len(mod.GameVersions))
		for i, version := range mod.GameVersions {
			versions[i] = version.Version
		}
		log.Infof("Found mod: %s (ID: %s, Image: %d, Game Versions: %s, Author: %s)",
			mod.Name, mod.ID, len(mod.Image), versions, mod.Author)
	}

	gameIndex, err := games.NewGameIndex("data/index.json")
	if err != nil {
		println("Error loading game index:", err.Error())
		return
	}

	// Iterate through /home/subtixx/.local/share/Steam/steamapps/common/ and print all games
	steamAppsPath := "/home/subtixx/.local/share/Steam/steamapps/common/"
	files, err := os.ReadDir(steamAppsPath)
	if err != nil {
		log.Fatalf("Failed to read Steam apps directory: %v", err)
	}
	for _, file := range files {
		if file.IsDir() {
			gamePath := steamAppsPath + file.Name()
			game, err := gameIndex.DetectGame(gamePath)
			if err != nil {
			} else {
				err := game.FetchInfoFromSteam()
				if err != nil {
					log.Errorf("Failed to fetch game info from Steam for %s: %v", game.Name, err)
					continue
				}
				// Save to data/games/<game.ID>.json
				err = game.Save("data/games/" + game.ID + ".json")
				if err != nil {
					log.Errorf("Failed to save game %s: %v", game.ID, err)
					continue
				}
				log.Infof("Detected game: %s (ID: %s, Steam App ID: %d)", game.Name, game.ID, game.SteamAppID)
			}
		}
	}
*/
