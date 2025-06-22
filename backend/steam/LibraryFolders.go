package steam

import (
	"TotalControl/backend/utils"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"os"
	"path"
	"time"
)

/*
Example file:
"libraryfolders"
{
	"0"
	{
		"path"		"/home/subtixx/.local/share/Steam"
		"label"		""
		"contentid"		"3750600032654308393"
		"totalsize"		"0"
		"update_clean_bytes_tally"		"2147789714"
		"time_last_update_verified"		"1750370904"
		"apps"
		{
			"92800"		"306849016"
			"105600"		"715107635"
			"228980"		"192183527"
			"300570"		"1872223447"
			"370360"		"84039298"
			"413150"		"708901690"
			"427520"		"1961969117"
			"504210"		"458845960"
			"558990"		"517220955"
			"716490"		"720514770"
			"858820"		"20099896615"
			"948420"		"134689365"
			"1070560"		"13004"
			"1106840"		"15108361510"
			"1168880"		"195989343"
			"1391110"		"655991598"
			"1493710"		"1391081609"
			"1511780"		"220588560"
			"1628350"		"781520723"
			"1677110"		"0"
			"1683150"		"8898329381"
		}
	}
}

*/

type LibraryApp struct {
	AppId      string `json:"appid"`
	ManifestId string `json:"manifestid"`
}

func (app *LibraryApp) IsBlacklisted() bool {
	if app == nil || app.AppId == "" {
		return false
	}

	if _, exists := blacklistedAppIDs[app.AppId]; exists {
		log.Warnf("App %s (%s) is blacklisted", app.AppId, blacklistedAppIDs[app.AppId])
		return true
	}

	return false
}

func (app *LibraryApp) IsValid() bool {
	if app == nil {
		return false
	}

	if app.AppId == "" {
		log.Warn("Invalid library app: AppId is empty")
		return false
	}

	if app.ManifestId == "" {
		log.Warn("Invalid library app: ManifestId is empty")
		return false
	}

	return true
}

type LibraryFolder struct {
	Path                   string
	Label                  string
	ContentId              string
	TotalSize              int64
	UpdateCleanBytesTally  int64
	TimeLastUpdateVerified time.Time
	Apps                   map[string]*LibraryApp
}

func (library *LibraryFolder) GetAppManifestPath(appId string) string {
	return path.Join(library.Path, "steamapps", "appmanifest_"+appId+".acf")
}

func IsValidLibraryObject(obj map[string]interface{}) bool {
	if obj == nil {
		return false
	}

	if !utils.ValidateStringFromMap("path", obj) {
		log.Warn("Invalid library object: invalid Path")
		return false
	}

	if !utils.ValidateStringFromMap("label", obj) {
		log.Warn("Invalid library object: invalid Label")
		return false
	}

	if !utils.ValidateStringFromMap("contentid", obj) {
		log.Warn("Invalid library object: invalid ContentId")
		return false
	}

	if !utils.ValidateIntFromMap("totalsize", obj) {
		log.Warn("Invalid library object: invalid TotalSize")
		return false
	}

	if !utils.ValidateIntFromMap("update_clean_bytes_tally", obj) {
		log.Warn("Invalid library object: invalid UpdateCleanBytesTally")
		return false
	}

	if !utils.ValidateIntFromMap("time_last_update_verified", obj) {
		log.Warn("Invalid library object: invalid TimeLastUpdateVerified")
		return false
	}

	// TODO: Apps

	return true
}

func NewLibraryFolder(data map[string]interface{}) (*LibraryFolder, error) {
	if data == nil {
		return nil, errors.New("data cannot be nil")
	}

	if !IsValidLibraryObject(data) {
		return nil, errors.New("invalid library folder data")
	}

	timeLastUpdateVerifiedTime := time.Unix(utils.GetIntFromMap("time_last_update_verified", data), 0)

	library := &LibraryFolder{
		Path:                   utils.GetStringFromMap("path", data),
		Label:                  utils.GetStringFromMap("label", data),
		ContentId:              utils.GetStringFromMap("contentid", data),
		TotalSize:              utils.GetIntFromMap("totalsize", data),
		UpdateCleanBytesTally:  utils.GetIntFromMap("update_clean_bytes_tally", data),
		TimeLastUpdateVerified: timeLastUpdateVerifiedTime,
		Apps:                   make(map[string]*LibraryApp),
	}
	apps, ok := data["apps"].(map[string]interface{})
	if ok {
		for appId, manifestId := range apps {
			manifestIdStr, ok := manifestId.(string)
			if !ok {
				log.Warnf("Invalid manifest ID for app %s: %v", appId, manifestId)
				continue
			}

			libraryApp := &LibraryApp{
				AppId:      appId,
				ManifestId: manifestIdStr,
			}
			if libraryApp.IsBlacklisted() {
				log.Debugf("Skipping blacklisted app: %s (%s)", libraryApp.AppId, blacklistedAppIDs[libraryApp.AppId])
				continue
			}

			library.Apps[appId] = libraryApp
		}
	}

	return library, nil
}

func GetSteamLibraries(steamInstallPath string) ([]*LibraryFolder, error) {
	libraryFoldersPath := path.Join(steamInstallPath, "steamapps", "libraryfolders.vdf")
	if _, err := os.Stat(libraryFoldersPath); os.IsNotExist(err) {
		return nil, err
	}

	libraries, err := ReadLibraryFolders(libraryFoldersPath)
	if err != nil {
		return nil, err
	}

	return libraries, nil
}

func ReadLibraryFolders(filePath string) ([]*LibraryFolder, error) {
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		log.Errorf("Failed to find libraryfolders.vdf: %s", err)
		return nil, err
	}

	vdf, err := ReadVDF(filePath)
	if err != nil {
		return nil, err
	}

	libraries, ok := vdf["libraryfolders"].(map[string]interface{})
	if !ok {
		return nil, errors.New("Invalid libraryfolders format")
	}

	var result []*LibraryFolder
	for key, value := range libraries {
		libData, ok := value.(map[string]interface{})
		if !ok {
			log.Warnf("Skipping invalid library folder entry: %s %+v", key, value)
			continue
		}

		library, err := NewLibraryFolder(libData)
		if err != nil {
			log.Warnf("Skipping invalid library folder entry %s -> %+v: %v", key, value, err)
			continue
		}

		result = append(result, library)
	}

	return result, nil
}
