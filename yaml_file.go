package localization

import (
	"fmt"
)

// YAMLFile is stores YAML translate file information and can be loaded in Locale.
// FileName - will provide full path of file with name (for better error messages).
// Translates - slice contains all loaded translates from YAML file.
type YAMLFile struct {
	FilePath   string
	Translates []Translate
}

// loadYAML is used to load and parse YAML file which contains translates.
// Returns YAMLFile pointer or error if something went wrong.
//
// Params:
// defaultLanguage - default language for non-list values (some_key: "value").
// path - YAML file path.
func loadYAML(defaultLanguage, path string) (*YAMLFile, error) {
	// Create new yamlContent (will for further processing).
	content := newYAMLContent(defaultLanguage)

	// Load file byte content.
	bytes, err := content.loadBytes(path)
	if err != nil {
		return nil, err
	}

	// Unmarshal file content.
	err = content.unmarshal(bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal '%s': %w", path, err)
	}

	// Apply file path for YAMLFile struct (for better error messages).
	yamlFile := YAMLFile{FilePath: path}

	// Parse YAML content.
	yamlFile.Translates, err = content.parse()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", path, err)
	}

	// Return parsed YAML file.
	return &yamlFile, nil
}

// LoadYAMLFiles can be used to load and parse one or more YAML files with
// containing translations.
// Returns []*YAMLFile or error if something went wrong.
// Returned []YAMLFile can be used as input in Locale to add translations.
//
// Params:
// path - YAML file paths.
// defaultLanguage - default language for non-list values (some_key: "value").
func LoadYAMLFiles(defaultLanguage string, path ...string) ([]*YAMLFile, error) {
	if len(path) == 0 {
		return nil, nil
	}

	yamlFiles := make([]*YAMLFile, 0)

	for _, v := range path {
		file, err := loadYAML(defaultLanguage, v)
		if err != nil {
			return nil, err
		}

		yamlFiles = append(yamlFiles, file)
	}

	return yamlFiles, nil
}
