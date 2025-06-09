package factorio

type ModInfo struct {
	Name            string   `json:"name"`
	Version         string   `json:"version"`
	FactorioVersion string   `json:"factorio_version"`
	Title           string   `json:"title"`
	Author          string   `json:"author"`
	Contact         string   `json:"contact"`
	Homepage        string   `json:"homepage"`
	Description     string   `json:"description"`
	Dependencies    []string `json:"dependencies"`

	Image []byte `json:"-"`
}
