package main

import (
	"TotalControl/backend/scripting"
	"TotalControl/backend/utils"
	"flag"
	log "github.com/sirupsen/logrus"
	lua "github.com/yuin/gopher-lua"
	"os"
)

func testLuaEngine() {
	luaEngine := scripting.NewLuaModProviderEngine()
	if err := luaEngine.LoadScript(`
local mods = {
	{id = "mod1", name = "Mod One", author = "Author A", version = "1.0", enabled = true, game_versions = { {version = "1.0"} }},
	{id = "mod2", name = "Mod Two", author = "Author B", version = "1.1", enabled = false, game_versions = { {version = "1.1"} }},
}
plugin = {
	GetInstalledMods = function()
		return mods
	end,
	GetModByID = function(id)
		for _, mod in ipairs(mods) do
			if mod.id == id then
				return mod
			end
		end
	end,
	AddMod = function(mod)
		log.Info("Adding mod:", mod.name)
		return true
	end,
	RemoveMod = function(id)
		log.Info("Removing mod with ID:", id)
		return true
	end,
	UpdateMod = function(mod)
		log.Info("Updating mod:", mod.name)
		return true
	end,
	GetGameModDirectory = function()
		return "/path/to/mods"
	end,
	GetGameID = function()
		return "game123"
	end,
}
log.debug("Lua script loaded successfully. Plugin table initialized with methods for mod management.")
log.debug(operating_system.getOperatingSystem())
if operating_system.is_windows then
	log.debug("Running on Windows")
elseif operating_system.is_linux then
	log.debug("Running on Linux")
elseif operating_system.is_macos then
	log.debug("Running on MacOS")
else
	log.debug("Running on an unknown operating system")
end
`); err != nil {
		log.Fatalf("Failed to load Lua script: %v", err)
	}

	// Get global plugin table
	plugin := luaEngine.L.GetGlobal("plugin")
	if plugin.Type() != lua.LTTable {
		log.Fatal("Plugin table not found in Lua script")
	}

	// Validate plugin table
	if luaEngine.L.GetField(plugin, "GetInstalledMods").Type() != lua.LTFunction {
		log.Fatal("GetInstalledMods method not found in plugin table")
	}
	if luaEngine.L.GetField(plugin, "GetModByID").Type() != lua.LTFunction {
		log.Fatal("GetModByID method not found in plugin table")
	}
	if luaEngine.L.GetField(plugin, "AddMod").Type() != lua.LTFunction {
		log.Fatal("AddMod method not found in plugin table")
	}
	if luaEngine.L.GetField(plugin, "RemoveMod").Type() != lua.LTFunction {
		log.Fatal("RemoveMod method not found in plugin table")
	}
	if luaEngine.L.GetField(plugin, "UpdateMod").Type() != lua.LTFunction {
		log.Fatal("UpdateMod method not found in plugin table")
	}
	if luaEngine.L.GetField(plugin, "GetGameModDirectory").Type() != lua.LTFunction {
		log.Fatal("GetGameModDirectory method not found in plugin table")
	}
	if luaEngine.L.GetField(plugin, "GetGameID").Type() != lua.LTFunction {
		log.Fatal("GetGameID method not found in plugin table")
	}

	// Call GetInstalledMods method
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
}

func main() {
	logPath := flag.String("log-path", "", "Path to log file (default: stdout)")
	logLevel := flag.String("log-level", "debug", "Log level (debug, info, warn, error, fatal, panic)")
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

	// Test factorio lua
	luaEngine := scripting.NewLuaModProviderEngine()
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
