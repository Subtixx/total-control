package steam

import (
	"TotalControl/backend/utils"
	log "github.com/sirupsen/logrus"
	"strconv"
)

/*
"AppState"
{
	"appid"		"92800"
	"Universe"		"1"
	"name"		"SpaceChem"
	"StateFlags"		"4"
	"installdir"		"SpaceChem"
	"lastupdated"		"1750372431"
	"LastPlayed"		"0"
	"SizeOnDisk"		"306849016"
	"StagingSize"		"0"
	"buildid"		"5484159"
	"LastOwner"		"76561198018743176"
	"DownloadType"		"1"
	"UpdateResult"		"0"
	"BytesToDownload"		"100657840"
	"BytesDownloaded"		"100657840"
	"BytesToStage"		"306849016"
	"BytesStaged"		"306849016"
	"TargetBuildID"		"5484159"
	"AutoUpdateBehavior"		"0"
	"AllowOtherDownloadsWhileRunning"		"0"
	"ScheduledAutoUpdate"		"0"
	"InstalledDepots"
	{
		"92804"
		{
			"manifest"		"6451263791585779898"
			"size"		"306849016"
		}
	}
	"UserConfig"
	{
		"language"		"english"
	}
	"MountedConfig"
	{
		"language"		"english"
	}
}
*/

type InstalledDepots struct {
	DepotID  string `json:"depotid"`
	Manifest string `json:"manifest"`
	Size     int64  `json:"size"`
}

type UserConfig struct {
	Language string `json:"language"`
}

type MountedConfig struct {
	Language string `json:"language"`
}

type AppState struct {
	AppID                           string                      `json:"appid"`
	Universe                        string                      `json:"universe"`
	Name                            string                      `json:"name"`
	StateFlags                      int                         `json:"stateflags"`
	InstallDir                      string                      `json:"installdir"`
	LastUpdated                     int64                       `json:"lastupdated"`
	LastPlayed                      int64                       `json:"lastplayed"`
	SizeOnDisk                      int64                       `json:"sizeondisk"`
	StagingSize                     int64                       `json:"stagingsize"`
	BuildID                         string                      `json:"buildid"`
	LastOwner                       string                      `json:"lastowner"`
	DownloadType                    int                         `json:"downloadtype"`
	UpdateResult                    int                         `json:"updateresult"`
	BytesToDownload                 int64                       `json:"bytestodownload"`
	BytesDownloaded                 int64                       `json:"bytesdownloaded"`
	BytesToStage                    int64                       `json:"bytestostage"`
	BytesStaged                     int64                       `json:"bytesstaged"`
	TargetBuildID                   string                      `json:"targetbuildid"`
	AutoUpdateBehavior              int                         `json:"autoupdatebehavior"`
	AllowOtherDownloadsWhileRunning int                         `json:"allowotherdownloadswilerunning"`
	ScheduledAutoUpdate             int                         `json:"scheduledautoupdate"`
	InstalledDepots                 map[string]*InstalledDepots `json:"installeddepots"`
	UserConfig                      *UserConfig                 `json:"userconfig"`
	MountedConfig                   *MountedConfig              `json:"mountedconfig"`
}

type AppSchemaFile struct {
	AppState AppState `json:"AppState"`
}

func NewInstalledDepots(data map[string]interface{}) *InstalledDepots {
	depotID, ok := data["DepotID"].(string)
	if !ok {
		depotID = ""
	}

	manifest, ok := data["Manifest"].(string)
	if !ok {
		manifest = ""
	}

	return &InstalledDepots{
		DepotID:  depotID,
		Manifest: manifest,
		Size:     utils.StringToInt64(data["Size"]),
	}
}

func NewMountedConfig(data map[string]interface{}) *MountedConfig {
	mountedConfig, ok := data["MountedConfig"].(map[string]interface{})
	if !ok {
		return nil
	}

	language, ok := mountedConfig["language"].(string)
	if !ok {
		language = ""
	}

	return &MountedConfig{
		Language: language,
	}
}

func NewUserConfig(data map[string]interface{}) *UserConfig {
	userConfig, ok := data["UserConfig"].(map[string]interface{})
	if !ok {
		return nil
	}

	language, ok := userConfig["language"].(string)
	if !ok {
		language = ""
	}

	return &UserConfig{
		Language: language,
	}
}

func NewAppSchemaFile(data map[string]interface{}) *AppSchemaFile {
	appState := data["AppState"].(map[string]interface{})
	installedDepots := make(map[string]*InstalledDepots)

	for key, value := range appState["InstalledDepots"].(map[string]interface{}) {
		depot, ok := value.(map[string]interface{})
		if !ok {
			log.Warnf("Invalid depot data for key %s: %v", key, value)
			continue
		}
		installedDepots[key] = NewInstalledDepots(depot)
	}

	targetBuildID, ok := appState["TargetBuildID"].(string)
	if !ok {
		targetBuildID = ""
	}

	return &AppSchemaFile{
		AppState: AppState{
			AppID:                           appState["appid"].(string),
			Universe:                        appState["Universe"].(string),
			Name:                            appState["name"].(string),
			StateFlags:                      utils.StringToInt(appState["StateFlags"]),
			InstallDir:                      appState["installdir"].(string),
			LastUpdated:                     utils.StringToInt64(appState["lastupdated"]),
			LastPlayed:                      utils.StringToInt64(appState["LastPlayed"]),
			SizeOnDisk:                      utils.StringToInt64(appState["SizeOnDisk"]),
			StagingSize:                     utils.StringToInt64(appState["StagingSize"]),
			BuildID:                         appState["buildid"].(string),
			LastOwner:                       appState["LastOwner"].(string),
			DownloadType:                    utils.StringToInt(appState["DownloadType"]),
			UpdateResult:                    utils.StringToInt(appState["UpdateResult"]),
			BytesToDownload:                 utils.StringToInt64(appState["BytesToDownload"]),
			BytesDownloaded:                 utils.StringToInt64(appState["BytesDownloaded"]),
			BytesToStage:                    utils.StringToInt64(appState["BytesToStage"]),
			BytesStaged:                     utils.StringToInt64(appState["BytesStaged"]),
			TargetBuildID:                   targetBuildID,
			AutoUpdateBehavior:              utils.StringToInt(appState["AutoUpdateBehavior"]),
			AllowOtherDownloadsWhileRunning: utils.StringToInt(appState["AllowOtherDownloadsWhileRunning"]),
			ScheduledAutoUpdate:             utils.StringToInt(appState["ScheduledAutoUpdate"]),
			InstalledDepots:                 installedDepots,
			UserConfig:                      NewUserConfig(appState),
			MountedConfig:                   NewMountedConfig(appState),
		},
	}
}

func (appState *AppState) String() string {
	return "AppState{\n" +
		"\t\tAppID: " + appState.AppID + ",\n" +
		"\t\tUniverse: " + appState.Universe + ",\n" +
		"\t\tName: " + appState.Name + ",\n" +
		"\t\tStateFlags: " + strconv.Itoa(appState.StateFlags) + ",\n" +
		"\t\tInstallDir: " + appState.InstallDir + ",\n" +
		"\t\tLastUpdated: " + strconv.FormatInt(appState.LastUpdated, 10) + ",\n" +
		"\t\tLastPlayed: " + strconv.FormatInt(appState.LastPlayed, 10) + ",\n" +
		"\t\tSizeOnDisk: " + strconv.FormatInt(appState.SizeOnDisk, 10) + ",\n" +
		"\t\tStagingSize: " + strconv.FormatInt(appState.StagingSize, 10) + ",\n" +
		"\t\tBuildID: " + appState.BuildID + ",\n" +
		"\t\tLastOwner: " + appState.LastOwner + ",\n" +
		"\t\tDownloadType: " + strconv.Itoa(appState.DownloadType) + ",\n" +
		"\t\tUpdateResult: " + strconv.Itoa(appState.UpdateResult) + ",\n" +
		"\t\tBytesToDownload: " + strconv.FormatInt(appState.BytesToDownload, 10) + ",\n" +
		"\t\tBytesDownloaded: " + strconv.FormatInt(appState.BytesDownloaded, 10) + ",\n" +
		"\t\tBytesToStage: " + strconv.FormatInt(appState.BytesToStage, 10) + ",\n" +
		"\t\tBytesStaged: " + strconv.FormatInt(appState.BytesStaged, 10) + ",\n" +
		"\t\tTargetBuildID: " + appState.TargetBuildID + ",\n" +
		"\t\tAutoUpdateBehavior: " + strconv.Itoa(appState.AutoUpdateBehavior) + ",\n" +
		"\t\tAllowOtherDownloadsWhileRunning: " + strconv.Itoa(appState.AllowOtherDownloadsWhileRunning) + ",\n" +
		"\t\tScheduledAutoUpdate: " + strconv.Itoa(appState.ScheduledAutoUpdate) + ",\n" +
		"\t\tInstalledDepots\n" +
		"\t\tUserConfig: {\n" +
		"\t\t\tLanguage: " + appState.UserConfig.Language + ",\n" +
		"\t\t},\n" +
		"\t\tMountedConfig: {\n" +
		"\t\t\tLanguage: " + appState.MountedConfig.Language + ",\n" +
		"\t\t},\n" +
		"\t}"
}

func (appSchemaFile *AppSchemaFile) String() string {
	return "AppSchemaFile{\n" +
		"\t" + appSchemaFile.AppState.String() + "\n" +
		"}"
}
