package main

import (
	"errors"
	"fmt"
	"reflect"
)

var ErrNoDynamicInput error = errors.New("no input for dynamic")
var ErrDynamicNotInt error = errors.New("plural input must be int")

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

func TextPlural(locale Locale, langKey, textKey string, isPlural bool) (string, error) {
	text, err := getPluralText(locale, langKey, textKey, isPlural)
	if err != nil {
		return "", err
	}

	return text, nil
}

func TextDynamic(locale Locale, langKey, textKey string, input ...interface{}) (string, error) {
	if len(input) == 0 {
		return "", ErrNoDynamicInput
	}

	text, err := locale.Value(langKey, textKey)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf(text, input...), nil
}

func TextPluralDynamic(locale Locale, langKey, textKey string, input ...interface{}) (string, error) {
	if len(input) == 0 {
		return "", ErrNoDynamicInput
	}

	isPlural, err := isPluralText(input...)
	if err != nil {
		return "", err
	}

	text, err := getPluralText(locale, langKey, textKey, isPlural)

	return fmt.Sprintf(text, input...), nil
}

func getPluralText(locale Locale, langKey, textKey string, isPlural bool) (string, error) {
	if isPlural {
		return locale.ValuePlural(langKey, textKey)
	}

	return locale.Value(langKey, textKey)
}

func isPluralText(input ...interface{}) (bool, error) {
	for _, v := range input {
		val := reflect.ValueOf(v)

		if val.Kind() != reflect.Int {
			return false, ErrDynamicNotInt
		}

		if val.Int() > 1 {
			return true, nil
		}
	}

	return false, nil
}
