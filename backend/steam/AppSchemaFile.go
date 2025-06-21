package steam

import (
	"TotalControl/backend/utils"
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
	AppID                           string                     `json:"appid"`
	Universe                        string                     `json:"universe"`
	Name                            string                     `json:"name"`
	StateFlags                      int                        `json:"stateflags"`
	InstallDir                      string                     `json:"installdir"`
	LastUpdated                     int64                      `json:"lastupdated"`
	LastPlayed                      int64                      `json:"lastplayed"`
	SizeOnDisk                      int64                      `json:"sizeondisk"`
	StagingSize                     int64                      `json:"stagingsize"`
	BuildID                         string                     `json:"buildid"`
	LastOwner                       string                     `json:"lastowner"`
	DownloadType                    int                        `json:"downloadtype"`
	UpdateResult                    int                        `json:"updateresult"`
	BytesToDownload                 int64                      `json:"bytestodownload"`
	BytesDownloaded                 int64                      `json:"bytesdownloaded"`
	BytesToStage                    int64                      `json:"bytestostage"`
	BytesStaged                     int64                      `json:"bytesstaged"`
	TargetBuildID                   string                     `json:"targetbuildid"`
	AutoUpdateBehavior              int                        `json:"autoupdatebehavior"`
	AllowOtherDownloadsWhileRunning int                        `json:"allowotherdownloadswilerunning"`
	ScheduledAutoUpdate             int                        `json:"scheduledautoupdate"`
	InstalledDepots                 map[string]InstalledDepots `json:"installeddepots"`
	UserConfig                      UserConfig                 `json:"userconfig"`
	MountedConfig                   MountedConfig              `json:"mountedconfig"`
}

type AppSchemaFile struct {
	AppState AppState `json:"AppState"`
}

func NewAppSchemaFile(data map[string]interface{}) *AppSchemaFile {
	appState := data["AppState"].(map[string]interface{})
	installedDepots := make(map[string]InstalledDepots)

	for key, value := range appState["InstalledDepots"].(map[string]interface{}) {
		depot := value.(map[string]interface{})
		installedDepots[key] = InstalledDepots{
			DepotID:  key,
			Manifest: depot["manifest"].(string),
			Size:     utils.StringToInt64(depot["size"].(string)),
		}
	}

	return &AppSchemaFile{
		AppState: AppState{
			AppID:                           appState["appid"].(string),
			Universe:                        appState["Universe"].(string),
			Name:                            appState["name"].(string),
			StateFlags:                      utils.StringToInt(appState["StateFlags"].(string)),
			InstallDir:                      appState["installdir"].(string),
			LastUpdated:                     utils.StringToInt64(appState["lastupdated"].(string)),
			LastPlayed:                      utils.StringToInt64(appState["LastPlayed"].(string)),
			SizeOnDisk:                      utils.StringToInt64(appState["SizeOnDisk"].(string)),
			StagingSize:                     utils.StringToInt64(appState["StagingSize"].(string)),
			BuildID:                         appState["buildid"].(string),
			LastOwner:                       appState["LastOwner"].(string),
			DownloadType:                    utils.StringToInt(appState["DownloadType"].(string)),
			UpdateResult:                    utils.StringToInt(appState["UpdateResult"].(string)),
			BytesToDownload:                 utils.StringToInt64(appState["BytesToDownload"].(string)),
			BytesDownloaded:                 utils.StringToInt64(appState["BytesDownloaded"].(string)),
			BytesToStage:                    utils.StringToInt64(appState["BytesToStage"].(string)),
			BytesStaged:                     utils.StringToInt64(appState["BytesStaged"].(string)),
			TargetBuildID:                   appState["TargetBuildID"].(string),
			AutoUpdateBehavior:              utils.StringToInt(appState["AutoUpdateBehavior"].(string)),
			AllowOtherDownloadsWhileRunning: utils.StringToInt(appState["AllowOtherDownloadsWhileRunning"].(string)),
			ScheduledAutoUpdate:             utils.StringToInt(appState["ScheduledAutoUpdate"].(string)),
			InstalledDepots:                 installedDepots,
			UserConfig: UserConfig{
				Language: appState["UserConfig"].(map[string]interface{})["language"].(string),
			},
			MountedConfig: MountedConfig{
				Language: appState["MountedConfig"].(map[string]interface{})["language"].(string),
			},
		},
	}
}

func (appSchemaFile *AppSchemaFile) String() string {
	return appSchemaFile.AppState.Name + " (" + appSchemaFile.AppState.AppID + ")"
}
