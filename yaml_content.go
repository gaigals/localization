package localization

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io"
	"os"
	"reflect"
)

// yamlContent is literally holds YAML translate file content and is used
// for content parsing.
//
// Params:
// defaultLanguage - default language for non-list values (some_key: "value").
// Data - unmarshalled YAML file content as map.
type yamlContent struct {
	defaultLanguage string
	Data            map[string]interface{}
}

// newYAMLContent constructs new yamlContent struct with default language field.
func newYAMLContent(defaultLanguage string) yamlContent {
	return yamlContent{defaultLanguage: defaultLanguage}
}

// loadBytes is used to load given file byte content.
// Returns byte slice or error if something went wrong.
func (c *yamlContent) loadBytes(path string) ([]byte, error) {
	// Ignore warning about "Potential file inclusion via variable".
	// #nosec G304
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open path '%s', error: %w",
			path, err)
	}

	bytes, err := io.ReadAll(file)
	if err != nil {
		_ = file.Close()
		return nil, fmt.Errorf("failed to cast '%s' as bytes, error: %w",
			path, err)
	}

	_ = file.Close()

	return bytes, nil
}

// unmarshal is used to unmarshal passed YAML file bytes as a map.
// Final results will be applied to yamlContent.Data field.
// Returns error if something went wrong.
func (c *yamlContent) unmarshal(bytes []byte) error {
	if bytes == nil {
		return fmt.Errorf("unmarshal failure, bytes slice is nil")
	}

	err := yaml.Unmarshal(bytes, &c.Data)
	if err != nil {
		return err
	}

	return nil
}

// parse goes over unmarshalled map and parses content as Translate slice.
// Returns Translate slice or error if something went wrong.
func (c *yamlContent) parse() ([]Translate, error) {
	// If yamlContent.Data is empty then return, nothing to do.
	if len(c.Data) == 0 {
		return nil, nil
	}

	// Create translates slice.
	translates := make([]Translate, 0)

	// Loop over all map keys.
	for k, v := range c.Data {
		// Get key->value reflect type and value.
		dataType, dataValue := c.getReflectData(v)

		// Read key->value content
		content, err := c.readContent(k, dataType, dataValue)
		if err != nil {
			return nil, err
		}

		// Append translates to the slice.
		translates = append(translates, content...)
	}

	// Done, return parsed translates.
	return translates, nil
}

// readContent reads provided map value and extracts translations by using reflection.
// Returns Translate slice and error if something went wrong.
//
// Examples how YAML gets unmarshalled as map[string]interface{}:
// key: "text" <- key:"text" (map[string]string)
//
// key:
//		- en: "text" <-- key:[map[en:"text"]] (map[string][]map[string]string)
//
// key:
//		- en: "text"
//		- lv: "other" <-- key:[map[en:"text"], map[lv:"other"]] (map[string][]map[string]string)
//
// key:
//		- en:
// 			- "non-plural"
//      	- "plural" <-- key:[map[en: ["non-plural", "plural"] ]] (map[string][]map[string][]string)
func (c *yamlContent) readContent(key string, dataType reflect.Type, dataValue reflect.Value) ([]Translate, error) {
	// Get data kind/type
	kind := dataType.Kind()

	// Is value a string then build translate and return.
	if kind == reflect.String {
		return []Translate{c.buildTranslateString(key, dataValue)}, nil
	}

	// Is value slice/list.
	if kind == reflect.Slice {
		// Extract slice values and return extracted/parsed values.
		// buildTranslateSlice() will call this function for every slice elements
		// as long as all values gets extracted.
		return c.buildTranslateSlice(key, dataValue)
	}

	// Is value map.
	if kind == reflect.Map {
		// Extract all map values and return extracted/parsed values.
		// buildTranslateMap() will call this function for every map element as
		// long as all map values gets extracted.
		return c.buildTranslateMap(key, dataValue)
	}

	return nil, fmt.Errorf("%s: unsupported type=%v", key, kind)
}

// buildTranslateString takes passed translate key, string value
// and builds Translate.
// Returns constructed Translate.
//
// Params:
// key - original yamlContent.Data map key (for error messages).
// dataValue - target string reflect.Value.
func (c *yamlContent) buildTranslateString(key string, dataValue reflect.Value) Translate {
	return Translate{
		Key:      key,
		Language: c.defaultLanguage,
		Value:    dataValue.String(),
	}
}

// buildTranslateSlice extracts slice values anc builds Translate slice from them.
// Returns Translate slice or error if something went wrong.
//
// Params:
// key - original yamlContent.Data map key (for error messages).
// dataValue - target slice reflect.Value.
func (c *yamlContent) buildTranslateSlice(key string, dataValue reflect.Value) ([]Translate, error) {
	languages := c.getSliceChilds(dataValue)
	translates := make([]Translate, 0)

	for k := range languages {
		langType := reflect.TypeOf(languages[k].Interface())

		results, err := c.readContent(key, langType, languages[k])
		if err != nil {
			return nil, err
		}

		translates = append(translates, results...)
	}

	return translates, nil
}

// buildTranslateMap extracts map contained translates.
// Returns extracted Translate slice or error if something went wrong.
//
// Params:
// key - original yamlContent.Data map key (for error messages).
// dataValue - target map reflect.Value.
func (c *yamlContent) buildTranslateMap(key string, dataValue reflect.Value) ([]Translate, error) {
	translates := make([]Translate, 0)

	// Get map range.
	mapRange := dataValue.MapRange()

	// Loop over all map keys.
	for mapRange.Next() {
		// Get map key reflect value.
		mapKey := reflect.ValueOf(mapRange.Key().Interface())
		// Get map value reflect value.
		mapValue := reflect.ValueOf(mapRange.Value().Interface())

		// If map key type is not string then return error.
		// This check is required to avoid panic.
		if mapKey.Kind() != reflect.String {
			return nil, fmt.Errorf("'%s' > '%v' must be string", key, mapKey)
		}

		// If map value type is string then build Translate and append it to final slice and
		// continue with next map key.
		if mapValue.Kind() == reflect.String {
			translates = append(translates, Translate{Key: key, Language: mapKey.String(), Value: mapValue.String()})
			continue
		}

		// If map is not string then it MUST be slice, if not, then return error.
		// This check is required to avoid panic.
		if mapValue.Kind() != reflect.Slice {
			return nil, fmt.Errorf("'%s' > '%v' value must be string or list", key, mapKey)
		}

		// Extract slice values (plurals in this case).
		plurals := c.getSliceChilds(mapValue)

		if len(plurals) > 2 {
			return nil, fmt.Errorf("'%s' > '%v': contains more than 2 plural entries", key, mapKey)
		}

		if len(plurals) == 1 {
			// Build translate from extracted plurals and append it to final slice.
			translates = append(translates, Translate{Key: key, Language: mapKey.String(), Value: plurals[0].String()})
			continue
		}

		// Build translate from extracted plurals and append it to final slice.
		translates = append(translates, Translate{Key: key, Language: mapKey.String(), Value: plurals[0].String(),
			Plural: plurals[1].String()})
	}

	return translates, nil
}

// getSliceChilds extract slice elements.
// Returns slice elements reflect values as slice or error if something went wrong.
func (c *yamlContent) getSliceChilds(dataValue reflect.Value) []reflect.Value {
	// Get length of slice.
	elementCount := dataValue.Len()
	// Create slice for each element.
	arr := make([]reflect.Value, elementCount)

	// Get reflect value for each element and apply it to the final slice.
	for k := range arr {
		arr[k] = reflect.ValueOf(dataValue.Index(k).Interface())
	}

	// Return results.
	return arr
}

// getReflectData is helper method for casting passed data interface as reflect type and value.
func (c *yamlContent) getReflectData(data interface{}) (reflect.Type, reflect.Value) {
	return reflect.TypeOf(data), reflect.ValueOf(data)
}
