package plugins

import (
	"TotalControl/backend/utils"
	"archive/zip"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	"os"
	"path"
	"strings"
)

type Plugin struct {
	PluginInfo

	PluginDir string
	IsPacked  bool
	Enabled   bool
	Settings  Setting
}

func (p *Plugin) Shutdown() error {
	err := p.Save()
	if err != nil {
		return err
	}
	p.Logger().Debug("Plugin shutdown successful")
	return nil
}

func (p *Plugin) Pack() error {
	if p.IsPacked {
		log.Warnf("Plugin %s is already packed, skipping packing process.", p.Name)
		return nil // Already packed
	}

	// Create the zip file in the plugin directory
	zipPath := p.PluginDir + "/" + utils.SanitizeFileName(p.Name) + ".tcplugin"
	zipFile, err := os.Create(zipPath)
	if err != nil {
		return err
	}
	defer func(zipFile *os.File) {
		err := zipFile.Close()
		if err != nil {
			log.Errorf("Failed to close zip file: %v", err)
			return
		}
	}(zipFile)

	zipWriter := zip.NewWriter(zipFile)
	defer func(zipWriter *zip.Writer) {
		err := zipWriter.Close()
		if err != nil {
			log.Errorf("Failed to close zip writer: %v", err)
			return
		}
	}(zipWriter)

	// Write info.json
	infoJsonPath := p.PluginDir + "/info.json"
	infoJsonData, err := os.ReadFile(infoJsonPath)
	if err != nil {
		return err
	}
	infoWriter, err := zipWriter.Create("info.json")
	if err != nil {
		return err
	}
	_, err = infoWriter.Write(infoJsonData)
	if err != nil {
		return err
	}

	// Copy all *.lua files into the zip
	files, err := os.ReadDir(p.PluginDir)
	if err != nil {
		return err
	}
	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".lua") {
			luaPath := p.PluginDir + "/" + file.Name()
			luaData, err := os.ReadFile(luaPath)
			if err != nil {
				return err
			}
			luaWriter, err := zipWriter.Create(file.Name())
			if err != nil {
				return err
			}
			_, err = luaWriter.Write(luaData)
			if err != nil {
				return err
			}
		}
	}

	p.IsPacked = true
	return nil
}

func GetPluginAppDataPath(pluginId uuid.UUID) string {
	appDataPath := utils.GetAppDataPath()
	return path.Join(appDataPath, "plugins", pluginId.String())
}

func (p *Plugin) GetPluginAppDataPath() string {
	appDataPath := utils.GetAppDataPath()
	return path.Join(appDataPath, "plugins", p.Id.String())
}

func (p *Plugin) Logger() *log.Entry {
	return log.WithFields(log.Fields{
		"prefix": "PLG",
		"plugin": p.Id.String(),
	})
}

func (p *Plugin) String() string {
	return "Plugin{" +
		"\n\tPluginDir: " + p.PluginDir +
		"\n\tIsPacked: " + utils.BoolToString(p.IsPacked) +
		"\n\tEnabled: " + utils.BoolToString(p.Enabled) +
		"\n\tInfo: " + p.PluginInfo.String() +
		"\n}"
}
