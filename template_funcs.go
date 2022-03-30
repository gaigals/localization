package main

import (
	"fmt"
)

func GetMapValue(data TemplateData, key string) (interface{}, error) {
	value, exist := data[key]
	if !exist {
		return nil, fmt.Errorf("passed data key does not exist")
	}

	return value, nil
}

func Text(locale Locale, langKey, textKey string) (string, error) {
	return locale.Value(langKey, textKey)
}
