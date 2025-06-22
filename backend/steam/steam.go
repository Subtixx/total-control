package steam

import (
	"errors"
	log "github.com/sirupsen/logrus"
	"os"
	"path"
)

var (
	errNoLibraryFolders = errors.New("no library folders found")
)

type Steam struct {
	InstallPath    string
	LibraryFolders []*LibraryFolder
	AppSchemas     map[string]*AppSchemaFile
}

func (s *Steam) GetLibraryFoldersFilePath() string {
	return path.Join(s.InstallPath, "steamapps", "libraryfolders.vdf")
}

func NewSteam() *Steam {
	steam := &Steam{
		LibraryFolders: make([]*LibraryFolder, 0),
		AppSchemas:     make(map[string]*AppSchemaFile),
	}

	if err := steam.Initialize(); err != nil {
		panic("Failed to initialize Steam: " + err.Error())
	}

	return steam
}

func (s *Steam) Initialize() error {
	steamInstallPath, err := FindSteamInstallation()
	if err != nil {
		return err
	}
	s.InstallPath = steamInstallPath

	if err := s.loadLibraryFolders(); err != nil {
		return err
	}
	if err := s.loadAppSchemas(); err != nil {
		return err
	}
	return nil
}

func (s *Steam) loadLibraryFolders() error {
	libraries, err := ReadLibraryFolders(s.GetLibraryFoldersFilePath())
	if err != nil {
		return err
	}

	s.LibraryFolders = libraries
	return nil
}

func (s *Steam) loadAppSchemas() error {
	if len(s.LibraryFolders) == 0 {
		return errNoLibraryFolders
	}

	for _, library := range s.LibraryFolders {
		for _, app := range library.Apps {
			appSchemaPath := library.GetAppManifestPath(app.AppId)
			if _, err := os.Stat(appSchemaPath); os.IsNotExist(err) {
				log.Warnf("App schema file does not exist: %s", appSchemaPath)
				continue
			}

			appSchema, err := ReadAppSchemaFile(library, appSchemaPath)
			if err != nil {
				return err
			}

			s.AppSchemas[app.AppId] = appSchema
		}
	}
	return nil
}
