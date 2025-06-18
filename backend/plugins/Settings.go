package plugins

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
)

var (
	ErrSettingNotFound    = fmt.Errorf("setting not found")
	ErrSettingValueNotSet = fmt.Errorf("setting value not set")
	ErrSettingInvalid     = fmt.Errorf("setting value is invalid")
	ErrValueTypeMismatch  = fmt.Errorf("value type does not match setting type")
)

type SettingDefinition struct {
	Type        string `json:"type"`
	Description string `json:"description"`
	Default     string `json:"default"`
}

type SettingDefinitions map[string]SettingDefinition
type Setting map[string]interface{}

func (s *SettingDefinition) Accepts(value interface{}) bool {
	// Check if value type matches setting type
	switch s.Type {
	case "string":
		_, ok := value.(string)
		return ok
	case "number", "int":
		_, ok := value.(int)
		return ok
	case "boolean", "bool":
		_, ok := value.(bool)
		return ok
	case "array":
		_, ok := value.([]interface{})
		return ok
	case "object":
		_, ok := value.(map[string]interface{})
		return ok
	default:
		return false // Unsupported type
	}
}

func (p *Plugin) Set(key string, value interface{}) error {
	setting, exists := p.SettingDefinitions[key]
	// Check if valid setting
	if !exists {
		return ErrSettingNotFound
	}

	if !setting.Accepts(value) {
		return ErrValueTypeMismatch
	}

	p.Settings[key] = value
	return nil
}

func (p *Plugin) Get(key string) (interface{}, error) {
	if p.Settings == nil {
		return nil, fmt.Errorf("settings not initialized")
	}

	_, exists := p.SettingDefinitions[key]
	if !exists {
		return nil, ErrSettingNotFound
	}
	value, exists := p.Settings[key]
	if !exists {
		return nil, ErrSettingValueNotSet
	}

	return value, nil
}

func (p *Plugin) Has(key string) bool {
	_, exists := p.Settings[key]
	return exists
}

func (p *Plugin) IsValid(key string, value interface{}) (bool, error) {
	setting, exists := p.SettingDefinitions[key]
	if !exists {
		return false, ErrSettingNotFound
	}

	if !setting.Accepts(value) {
		return false, ErrValueTypeMismatch
	}

	return true, nil
}

func (p *Plugin) Save() error {
	// If no values, use all default values
	if len(p.Settings) == 0 {
		p.Settings = make(Setting)
		for key, setting := range p.SettingDefinitions {
			p.Settings[key] = setting.Default
		}
	}

	settingsFilePath := path.Join(p.GetPluginAppDataPath(), "settings.json")
	file, err := os.Create(settingsFilePath)
	if err != nil {
		return err
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			fmt.Printf("Failed to close file: %v\n", err)
			return
		}
	}(file)

	encoder := json.NewEncoder(file)
	return encoder.Encode(p.Settings)
}

func (s *SettingDefinitions) String() string {
	if s == nil {
		return "SettingDefinitions{nil}"
	}
	result := "\tSettingDefinitions {"
	for key, setting := range *s {
		defaultValue := setting.Default
		if defaultValue == "" {
			defaultValue = "null"
		}
		result += "\n\t\t" + key + ": {Type: " + setting.Type + ", Description: " + setting.Description + ", Default: " + defaultValue + "},"
	}
	result += "\n\t}"
	return result
}
