package games

type GameIndexEntry struct {
	ID          string                     `json:"id"`
	Name        string                     `json:"name"`
	SteamAppID  int                        `json:"steam_appid"`
	Executables []GameIndexExecutableEntry `json:"executables"`
}
