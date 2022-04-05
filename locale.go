package main

import (
	"fmt"
	"strings"
)

type Locale struct {
	Languages   []Language
	StrictUsage bool // Is other language usage allowed if key does not exist for given lang.
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
	hasLang, lang := l.hasLanguage(langKey)
	if !hasLang {
		return "", fmt.Errorf("language '%s' does not exist", langKey)
	}

	langList := l.buildLanguageList(lang)
	text := ""

	for k := range langList {
		text = langList[k].Value(textKey)

		if text != "" {
			return text, nil
		}
	}

	if l.StrictUsage {
		return "", fmt.Errorf("language '%s' does not contain key '%s'", langKey, textKey)
	}

	return "", fmt.Errorf("non of languages contain key '%s'", textKey)
}

func (l *Locale) ValuePlural(langKey, textKey string) (string, error) {
	hasLang, lang := l.hasLanguage(langKey)
	if !hasLang {
		return "", fmt.Errorf("language '%s' does not exist", langKey)
	}

	langList := l.buildLanguageList(lang)
	text := ""

	for k := range langList {
		text = langList[k].ValuePlural(textKey)

		if text != "" {
			return text, nil
		}
	}

	if l.StrictUsage {
		return "", fmt.Errorf("language '%s' does not contain key '%s'", langKey, textKey)
	}

	return "", fmt.Errorf("non of languages contain key '%s'", textKey)
}

func (l *Locale) SetValue(langKey, textKey, value, plural string) error {
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
	lang.SetValue(textKey, value, plural)

	return nil
}

func (l *Locale) SetValueNoErr(langKey, textKey, value, plural string) {
	_, lang := l.hasLanguage(langKey)
	lang.SetValue(textKey, value, plural)
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
		err := l.SetValue(v.Language, v.Key, v.Value, v.Plural)
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

func (l *Locale) buildLanguageList(language *Language) []*Language {
	if l.StrictUsage {
		return []*Language{language}
	}

	keys := make([]*Language, len(l.Languages))
	keys[0] = language

	idx := 1

	currentLangKey := language.Keyword

	for k, v := range l.Languages {
		if v.Keyword != currentLangKey {
			keys[idx] = &l.Languages[k]
			idx++
		}
	}

	return keys
}
