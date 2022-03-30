package main

import (
	"fmt"
	"strings"
)

type Locale struct {
	Languages []Language
}

func (t *Locale) AddLanguages(lang ...string) {
	if len(lang) == 0 {
		return
	}

	languages := make([]Language, len(lang))

	for k, v := range lang {
		language := &languages[k]

		language.Keyword = v
		language.Map = make(TextMap, 0)
	}

	t.Languages = append(t.Languages, languages...)
}

func (t *Locale) Value(langKey, textKey string) (string, error) {
	if len(t.Languages) == 0 {
		return "", fmt.Errorf("cannot extract value, language list is empty")
	}

	hasLang, lang := t.hasLanguage(langKey)
	if !hasLang {
		return "", fmt.Errorf("language '%s' does not exist", langKey)
	}

	value := lang.Value(textKey)
	if value == "" {
		return "", fmt.Errorf("language '%s' does not contain key '%s'", langKey, textKey)
	}

	return value, nil
}

func (t *Locale) SetValue(langKey, textKey, value string) error {
	if len(t.Languages) == 0 {
		return fmt.Errorf("cannot extract value, language list is empty")
	}

	if textKey == "" {
		return fmt.Errorf("textKey is empty string")
	}

	hasLang, lang := t.hasLanguage(langKey)
	if !hasLang {
		return fmt.Errorf("language '%s' does not exist", langKey)
	}

	lang.SetValue(textKey, value)

	return nil
}

func (t *Locale) SetValueNoErr(langKey, textKey, value string) {
	_, lang := t.hasLanguage(langKey)
	lang.SetValue(textKey, value)
}

func (t *Locale) hasLanguage(langKey string) (bool, *Language) {
	for k, v := range t.Languages {
		if strings.EqualFold(langKey, v.Keyword) {
			return true, &t.Languages[k]
		}
	}

	return false, nil
}
