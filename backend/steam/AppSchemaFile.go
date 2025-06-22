package steam

import (
	"TotalControl/backend/utils"
	"errors"
	log "github.com/sirupsen/logrus"
	lua "github.com/yuin/gopher-lua"
	"path"
	"strconv"
	"time"
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

var ErrInvalidData = errors.New("invalid data format")

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
	LastUpdated                     time.Time                   `json:"lastupdated"`
	LastPlayed                      time.Time                   `json:"lastplayed"`
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
	Library  *LibraryFolder `json:"-"`
	AppState AppState       `json:"AppState"`
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

func NewAppSchemaFile(libraryFolder *LibraryFolder, data map[string]interface{}) *AppSchemaFile {
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

	lastUpdatedTime := utils.GetTimeFromMap("lastupdated", appState)
	lastPlayedTime := utils.GetTimeFromMap("LastPlayed", appState)

	return &AppSchemaFile{
		Library: libraryFolder,
		AppState: AppState{
			AppID:                           appState["appid"].(string),
			Universe:                        appState["Universe"].(string),
			Name:                            appState["name"].(string),
			StateFlags:                      utils.StringToInt(appState["StateFlags"]),
			InstallDir:                      appState["installdir"].(string),
			LastUpdated:                     lastUpdatedTime,
			LastPlayed:                      lastPlayedTime,
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

func ReadAppSchemaFile(libraryFolder *LibraryFolder, filePath string) (*AppSchemaFile, error) {
	data, err := ReadVDF(filePath)
	if err != nil {
		log.Errorf("Failed to read AppSchemaFile: %s", err)
		return nil, err
	}

	appSchemaFile := NewAppSchemaFile(libraryFolder, data)
	if appSchemaFile == nil {
		log.Error("Failed to create AppSchemaFile from data")
		return nil, ErrInvalidData
	}

	return appSchemaFile, nil
}

func (appState *AppState) String() string {
	return "AppState{\n" +
		"\t\tAppID: " + appState.AppID + ",\n" +
		"\t\tUniverse: " + appState.Universe + ",\n" +
		"\t\tName: " + appState.Name + ",\n" +
		"\t\tStateFlags: " + strconv.Itoa(appState.StateFlags) + ",\n" +
		"\t\tInstallDir: " + appState.InstallDir + ",\n" +
		"\t\tLastUpdated: " + strconv.FormatInt(appState.LastUpdated.Unix(), 10) + ",\n" +
		"\t\tLastPlayed: " + strconv.FormatInt(appState.LastPlayed.Unix(), 10) + ",\n" +
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

func (appSchemaFile *AppSchemaFile) GetAppFullInstallPath() string {
	if appSchemaFile.Library == nil {
		log.Warn("LibraryFolder is nil, cannot get full install path")
		return ""
	}

	return path.Join(appSchemaFile.Library.Path, "steamapps", "common", appSchemaFile.AppState.InstallDir)
}

func (appSchemaFile *AppSchemaFile) ToLuaTable(L *lua.LState) *lua.LTable {
	appTable := L.NewTable()
	L.SetField(appTable, "app_id", utils.ToLuaValue(L, appSchemaFile.AppState.AppID))
	L.SetField(appTable, "universe", utils.ToLuaValue(L, appSchemaFile.AppState.Universe))
	L.SetField(appTable, "name", utils.ToLuaValue(L, appSchemaFile.AppState.Name))
	L.SetField(appTable, "state_flags", utils.ToLuaValue(L, appSchemaFile.AppState.StateFlags))
	L.SetField(appTable, "install_dir", utils.ToLuaValue(L, appSchemaFile.AppState.InstallDir))
	L.SetField(appTable, "last_updated", utils.ToLuaValue(L, appSchemaFile.AppState.LastUpdated.Unix()))
	L.SetField(appTable, "last_played", utils.ToLuaValue(L, appSchemaFile.AppState.LastPlayed.Unix()))
	L.SetField(appTable, "size_on_disk", utils.ToLuaValue(L, appSchemaFile.AppState.SizeOnDisk))
	L.SetField(appTable, "staging_size", utils.ToLuaValue(L, appSchemaFile.AppState.StagingSize))
	L.SetField(appTable, "build_id", utils.ToLuaValue(L, appSchemaFile.AppState.BuildID))
	L.SetField(appTable, "last_owner", utils.ToLuaValue(L, appSchemaFile.AppState.LastOwner))
	L.SetField(appTable, "download_type", utils.ToLuaValue(L, appSchemaFile.AppState.DownloadType))
	L.SetField(appTable, "update_result", utils.ToLuaValue(L, appSchemaFile.AppState.UpdateResult))
	L.SetField(appTable, "bytes_to_download", utils.ToLuaValue(L, appSchemaFile.AppState.BytesToDownload))
	L.SetField(appTable, "bytes_downloaded", utils.ToLuaValue(L, appSchemaFile.AppState.BytesDownloaded))
	L.SetField(appTable, "bytes_to_stage", utils.ToLuaValue(L, appSchemaFile.AppState.BytesToStage))
	L.SetField(appTable, "bytes_staged", utils.ToLuaValue(L, appSchemaFile.AppState.BytesStaged))
	L.SetField(appTable, "target_build_id", utils.ToLuaValue(L, appSchemaFile.AppState.TargetBuildID))
	L.SetField(appTable, "auto_update_behavior", utils.ToLuaValue(L, appSchemaFile.AppState.AutoUpdateBehavior))
	L.SetField(appTable, "allow_other_downloads_while_running", utils.ToLuaValue(L, appSchemaFile.AppState.AllowOtherDownloadsWhileRunning))
	L.SetField(appTable, "scheduled_auto_update", utils.ToLuaValue(L, appSchemaFile.AppState.ScheduledAutoUpdate))
	installedDepotsTable := L.NewTable()
	for depotID, depot := range appSchemaFile.AppState.InstalledDepots {
		depotTable := L.NewTable()
		L.SetField(depotTable, "depot_id", utils.ToLuaValue(L, depot.DepotID))
		L.SetField(depotTable, "manifest", utils.ToLuaValue(L, depot.Manifest))
		L.SetField(depotTable, "size", utils.ToLuaValue(L, depot.Size))
		L.SetField(installedDepotsTable, depotID, depotTable)
	}
	L.SetField(appTable, "installed_depots", installedDepotsTable)
	L.SetField(appTable, "user_config", utils.ToLuaValue(L, appSchemaFile.AppState.UserConfig.Language))
	L.SetField(appTable, "mounted_config", utils.ToLuaValue(L, appSchemaFile.AppState.MountedConfig.Language))
	if appSchemaFile.Library != nil {
		libraryTable := appSchemaFile.Library.ToLuaTable(L)
		L.SetField(appTable, "library", libraryTable)
	} else {
		log.Warn("LibraryFolder is nil, cannot set library field in Lua table")
		L.SetField(appTable, "library", lua.LNil)
	}
	return appTable
}
