package main

import (
	"fmt"
	"strings"
)

type Locale struct {
	Languages    []Language
	AllowDefault bool
}

func (l *Locale) AddLanguages(lang ...string) {
	if len(lang) == 0 {
		return
	}

	languages := make([]Language, len(lang))

	for k, v := range lang {
		language := &languages[k]

		language.Keyword = v
		language.Map = make(TextMap, 0)
	}

	l.Languages = append(l.Languages, languages...)
}

func (l *Locale) Value(langKey, textKey string) (string, error) {
	if len(l.Languages) == 0 {
		return "", fmt.Errorf("cannot extract value, language list is empty")
	}

	hasLang, lang := l.hasLanguage(langKey)
	if !hasLang {
		return "", fmt.Errorf("language '%s' does not exist", langKey)
	}

	value := lang.Value(textKey)
	if value == "" && !l.AllowDefault {
		return "", fmt.Errorf("language '%s' does not contain key '%s'", langKey, textKey)
	}
	if value == "" {
		value = l.Languages[0].Value(textKey)
		if value == "" {
			return "", fmt.Errorf("language '%s' (backup) does not contain key '%s'",
				langKey, textKey)
		}
	}

	return value, nil
}

func (l *Locale) SetValue(langKey, textKey, value string) error {
	if len(l.Languages) == 0 {
		return fmt.Errorf("cannot extract value, language list is empty")
	}

	if textKey == "" {
		return fmt.Errorf("textKey is empty string")
	}

	hasLang, lang := l.hasLanguage(langKey)
	if !hasLang {
		return fmt.Errorf("language '%s' does not exist", langKey)
	}
	lang.SetValue(textKey, value)

	return nil
}

func (l *Locale) SetValueNoErr(langKey, textKey, value string) {
	_, lang := l.hasLanguage(langKey)
	lang.SetValue(textKey, value)
}

func (l *Locale) hasLanguage(langKey string) (bool, *Language) {
	for k, v := range l.Languages {
		if strings.EqualFold(langKey, v.Keyword) {
			return true, &l.Languages[k]
		}
	}

	return false, nil
}

func (l *Locale) AddTranslate(translates ...Translate) error {
	if len(translates) == 0 {
		return nil
	}

	for k, v := range translates {
		err := locale.SetValue(v.Language, v.Key, v.Value)
		if err != nil {
			return fmt.Errorf("translate index=%d: %w",
				k, err)
		}
	}

	return nil
}

func (l *Locale) AddYAMLFile(files ...*YAMLFile) error {
	if len(files) == 0 {
		return nil
	}

	for _, v := range files {
		err := l.AddTranslate(v.Translates...)
		if err != nil {
			return fmt.Errorf("'%s': %w", v.FileName, err)
		}
	}

	return nil
}
