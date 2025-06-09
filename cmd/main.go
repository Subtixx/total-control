package main

import (
	"TotalControl/backend/mods/factorio"
	"TotalControl/backend/utils"
	"flag"
	log "github.com/sirupsen/logrus"
	"os"
)

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
	log.SetFormatter(&utils.CustomFormatter{})

	provider := factorio.FactorioModProvider{}
	factorioMods, err := provider.GetMods()
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
}

/*

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
