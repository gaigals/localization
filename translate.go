package main

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"reflect"
)

type YAMLFile struct {
	FileName   string
	Translates []Translate
}

type Translate struct {
	Key      string
	Language string
	Value    string
	Plural   string
}

type YAMLContent struct {
	DefaultLanguage string
	Data            map[string]interface{}
}

func LoadTranslates(path, defaultLanguage string) (*YAMLFile, error) {
	bytes, err := loadFileContent(path)
	if err != nil {
		return nil, err
	}

	yamlContent := newYAMLContent(defaultLanguage)

	err = yamlContent.unmarshal(bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal '%s': %w", path, err)
	}

	yamlFile := YAMLFile{}

	yamlFile.Translates, err = yamlContent.parse()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", path, err)
	}

	return &yamlFile, nil
}

func loadFileContent(path string) ([]byte, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open path '%s', error: %w",
			path, err)
	}

	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		_ = file.Close()
		return nil, fmt.Errorf("failed to cast '%s' as bytes, error: %w",
			path, err)
	}

	_ = file.Close()

	return bytes, nil
}

func newYAMLContent(defaultLanguage string) YAMLContent {
	return YAMLContent{DefaultLanguage: defaultLanguage}
}

func (c *YAMLContent) unmarshal(bytes []byte) error {
	err := yaml.Unmarshal(bytes, &c.Data)
	if err != nil {
		return err
	}

	return nil
}

func (c *YAMLContent) parse() ([]Translate, error) {
	if len(c.Data) == 0 {
		return nil, nil
	}

	translates := make([]Translate, 0)

	for k, v := range c.Data {
		dataType, dataValue := c.getReflectData(v)

		content, err := c.readContent(k, dataType, dataValue)
		if err != nil {
			return nil, err
		}

		translates = append(translates, content...)
	}

	return translates, nil
}

func (c *YAMLContent) getReflectData(data interface{}) (reflect.Type, reflect.Value) {
	return reflect.TypeOf(data), reflect.ValueOf(data)
}

func (c *YAMLContent) string(dataValue reflect.Value) string {
	return dataValue.String()
}

func (c *YAMLContent) extractMapRange(dataValue reflect.Value) *reflect.MapIter {
	return dataValue.MapRange()
}

func (c *YAMLContent) readContent(key string, dataType reflect.Type, dataValue reflect.Value) ([]Translate, error) {
	kind := dataType.Kind()

	if kind == reflect.String {
		translate := Translate{
			Key:      key,
			Language: c.DefaultLanguage,
			Value:    dataValue.String(),
		}

		return []Translate{translate}, nil
	}

	if kind == reflect.Slice {
		languages, err := c.getSliceChilds(dataValue)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", err)
		}

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

	if kind == reflect.Map {
		return c.getMapLangContents(key, dataValue)
	}

	return nil, fmt.Errorf("%s: unsupported type=%v", key, kind)
}

func (c *YAMLContent) getMapLangContents(key string, dataValue reflect.Value) ([]Translate, error) {
	translates := make([]Translate, 0)
	mapRange := dataValue.MapRange()

	for mapRange.Next() {
		lang := reflect.ValueOf(mapRange.Key().Interface())
		value := reflect.ValueOf(mapRange.Value().Interface())

		if lang.Kind() != reflect.String {
			return nil, fmt.Errorf("'%s'>'%s' must be string", key, lang)
		}

		if value.Kind() == reflect.String {
			translates = append(translates, Translate{Key: key, Language: lang.String(), Value: value.String()})
			continue
		}

		if value.Kind() != reflect.Slice {
			return nil, fmt.Errorf("'%s'>'%s' value must be string or list", key, lang.String())
		}

		plurals, err := c.getSliceChilds(value)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", key, err)
		}

		translates = append(translates, Translate{Key: key, Language: lang.String(), Value: plurals[0].String(), Plural: plurals[1].String()})
	}

	return translates, nil
}

func (c *YAMLContent) getSliceChilds(dataValue reflect.Value) ([]reflect.Value, error) {
	indexCount := dataValue.Len()
	arr := make([]reflect.Value, indexCount)

	for k := range arr {
		arr[k] = reflect.ValueOf(dataValue.Index(k).Interface())

		//if arr[k].Kind() != reflect.String {
		//	return nil, fmt.Errorf("not a string")
		//}
	}

	return arr, nil
}
