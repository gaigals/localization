package localization

import (
	"fmt"
	"path/filepath"
	"strings"
)

// Locale contains all initialized languages and can be used for handling
// localization/translation. Each Locale.Language represents translation
// layer and each layer contains initialized translation keywords/keys.
type Locale struct {
	Languages   []Language // List of initialized languages.
	StrictUsage bool       // Is other language usage allowed if key does not exist for given lang.
}

// NewLocale can be used to initialize new Locale structure with provided languages.
// Params:
// strictUsage - is translation key usage restricted to specific language or other
// languages can be used as backup. Set TRUE to restrict keyword usage.
// lang - list of languages keywords to initialize ("en", "lv" etc).
func NewLocale(strictUsage bool, lang ...string) (*Locale, error) {
	locale := Locale{StrictUsage: strictUsage}

	err := locale.addLanguages(lang...)
	if err != nil {
		return nil, err
	}

	return &locale, nil
}

// GlobalYAMLLoad loads YAML files with given pattern.
// Examples:
// "file.yml", "*.yaml", "path/*", "**/*.yml"
func (l *Locale) GlobalYAMLLoad(defaultLang, pattern string) error {
	// Use filepath.Glob to find matching files
	files, err := filepath.Glob(pattern)
	if err != nil {
		return fmt.Errorf("locale: failed to match pattern: %w", err)
	}

	fmt.Println(files)
	err = l.LoadYAMLFile(defaultLang, files...)
	if err != nil {
		return fmt.Errorf("locale: yaml load error: %w", err)
	}

	return nil
}

// addLanguages can be used to add new languages to Locale.
// Returns error if something went wrong.
// Params:
// lang - list of languages keywords to initialize ("en", "lv" etc).
func (l *Locale) addLanguages(lang ...string) error {
	if len(lang) == 0 {
		return nil
	}

	languages := make([]Language, len(lang))

	for k, v := range lang {
		// Check if language keyword already is initialized (important to avoid future bugs).
		// err == nil if language with keyword exists.
		_, err := l.GetLanguage(v)
		if err == nil {
			return fmt.Errorf("language '%s' already exists", v)
		}

		// Check if lang param element is unique (important to avoid future bugs).
		// Will return error if param 'lang' contains []string{"en", "lv", "lv"}.
		count := l.countLangSliceEntries(v, lang)
		if count != 1 {
			return fmt.Errorf("language '%s' redefined in passed lang parameter", v)
		}

		language := &languages[k]

		language.Keyword = v
		language.Map = make(TextMap, 0)
	}

	l.Languages = append(l.Languages, languages...)

	return nil
}

// Value can be used to extract non-plural translation from target language
// by providing translation keyword/key. If Locale.StrictUsage is FALSE then
// other languages will be used as backup for searching keyword/key.
// Returns non-plural value or error if langKey does not exist, or key does not exist.
//
// Params:
// langKey - target language keyword ("en", "lv" etc).
// textKey - translation keyword/key.
func (l *Locale) Value(langKey, textKey string) (string, error) {
	return l._value(langKey, textKey, false)
}

// ValueNoErr can be used to extract non-plural translation from target language
// by providing translation keyword/key. If Locale.StrictUsage is FALSE then
// other languages will be used as backup for searching keyword/key.
// Does NOT return error if key or language key does not exist.
// Returns value as string or empty string if langKey or textKey does not exist.
//
// Params:
// langKey - target language keyword ("en", "lv" etc).
// textKey - translation keyword/key.
func (l *Locale) ValueNoErr(langKey, textKey string) string {
	value, _ := l._value(langKey, textKey, false)
	return value
}

// ValuePlural can be used to extract plural translation from target language
// by providing translation keyword/key. If Locale.StrictUsage is FALSE then
// other languages will be used as backup for searching keyword/key.
// Returns plural value or error if langKey does not exist, or key does not exist.
//
// Params:
// langKey - target language keyword ("en", "lv" etc).
// textKey - translation keyword/key.
func (l *Locale) ValuePlural(langKey, textKey string) (string, error) {
	return l._value(langKey, textKey, true)
}

// ValuePluralNoErr can be used to extract plural translation from target language
// by providing translation keyword/key. If Locale.StrictUsage is FALSE then
// other languages will be used as backup for searching keyword/key.
// Does NOT return error if key or language key does not exist.
// Returns plural value as string or empty string if langKey or textKey
// does not exist.
//
// Params:
// langKey - target language keyword ("en", "lv" etc).
// textKey - translation keyword/key.
func (l *Locale) ValuePluralNoErr(langKey, textKey string) string {
	value, _ := l._value(langKey, textKey, true)
	return value
}

// _value is helper method which can be used to extract plural translation from
// target language by providing translation keyword/key.
// If Locale.StrictUsage is FALSE then
// other languages will be used as backup for searching keyword/key.
// Returns plural or non-plural value or error if langKey does not exist,
// or key does not exist.
//
// Params:
// langKey - target language keyword ("en", "lv" etc).
// textKey - translation keyword/key.
// isPlural - find plural translation value.
func (l *Locale) _value(langKey, textKey string, isPlural bool) (string, error) {
	lang, err := l.GetLanguage(langKey)
	if err != nil {
		return "", err
	}

	langList := l.buildPrioritizedLanguageList(lang)
	text := ""

	for k := range langList {
		if !isPlural {
			text, err = langList[k].Value(textKey)
		} else {
			text, err = langList[k].ValuePlural(textKey)
		}

		if err == nil {
			return text, nil
		}
	}

	if l.StrictUsage {
		return "", fmt.Errorf("language '%s' does not contain key '%s'", langKey, textKey)
	}

	return "", fmt.Errorf("non of languages contain key '%s'", textKey)
}

// SetValue can be used to set translation plural and non-plural values for target
// language.
// Returns error if something went wrong.
// Params:
// langKey - target language keyword ("en", "lv" etc).
// textKey - translation keyword/key.
// value - non-plural value.
// plural - plural value.
func (l *Locale) SetValue(langKey, textKey, value, plural string) error {
	if len(l.Languages) == 0 {
		return fmt.Errorf("language list is empty")
	}

	lang, err := l.GetLanguage(langKey)
	if err != nil {
		return err
	}

	lang.SetValue(textKey, value, plural)

	return nil
}

// SetValueNoErr can be used to set translation plural and non-plural values for target
// language. Works exactly like Locale.SetValue but without any error checking and expects
// successful outcome. Use this if all error checking is done by caller API.
//
// Params:
// langKey - target language keyword ("en", "lv" etc).
// textKey - translation keyword/key.
// value - non-plural value.
// plural - plural value.
func (l *Locale) SetValueNoErr(langKey, textKey, value, plural string) {
	lang, _ := l.GetLanguage(langKey)
	lang.SetValue(textKey, value, plural)
}

// GetLanguage can be used to get language with specific keyword ("en", "lv" etc).
// Returns pointer to target language or error if language with provided keyword
// does not exist
func (l *Locale) GetLanguage(langKey string) (*Language, error) {
	for k, v := range l.Languages {
		if strings.EqualFold(langKey, v.Keyword) {
			return &l.Languages[k], nil
		}
	}

	return nil, fmt.Errorf("language '%s' does not exist", langKey)
}

// AddTranslate can be used to add 1 or more translations to current Locale.
// Returns error if something went wrong.
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

// AddYAMLFile can be used to add 1 or more YAMLFile's translations to current Locale.
// Returns error if something went wrong.
func (l *Locale) AddYAMLFile(files ...*YAMLFile) error {
	if len(files) == 0 {
		return nil
	}

	for k, v := range files {
		// Check if current YAMLFile is not nil.
		if v == nil {
			return fmt.Errorf("YAMLFile with index=%d is nil", k)
		}

		err := l.AddTranslate(v.Translates...)
		if err != nil {
			return fmt.Errorf("'%s': %w", v.FilePath, err)
		}
	}

	return nil
}

// LoadYAMLFile can be used to load and parse multiple YAML files with
// containing translations and directly load them into current Locale.
// Requires previous language initialization (Locale.AddLanguages()) before
// YAML file loading.
// Returns error if something went wrong.
//
// Params:
// path - YAML file paths.
// defaultLanguage - default language for non-list values (some_key: "value").
func (l *Locale) LoadYAMLFile(defaultLanguage string, filePath ...string) error {
	// Load/parse provided YAML files.
	yamlFiles, err := LoadYAMLFiles(defaultLanguage, filePath...)
	if err != nil {
		return err
	}

	// Add YAML file translations to the current Locale.
	return l.AddYAMLFile(yamlFiles...)
}

// buildPrioritizedLanguageList is used to build prioritized list of Languages
// for searching plural and non-plural values.
// If Locale.StrictUsage is TRUE then method will return slice of passed language as
// other language backups are restricted.
// If Locale.StrictUsage is FALSE then method will return slice of prioritized languages.
// Either case, passed language will always be as first element in returned slice.
func (l *Locale) buildPrioritizedLanguageList(language *Language) []*Language {
	// If StrictUsage is enabled then return slice with passed language (no backups
	// allowed).
	if l.StrictUsage {
		return []*Language{language}
	}

	// Create new slice.
	langs := make([]*Language, len(l.Languages))
	// Add passed language as first element in slice (first priority).
	langs[0] = language

	// Set slice index. Start from 1 (0 is already reserved).
	idx := 1

	for k, v := range l.Languages {
		// If this is same language as passed one then continue.
		if v.Keyword == language.Keyword {
			continue
		}

		// Apply current language to slice.
		langs[idx] = &l.Languages[k]
		idx++
	}

	return langs
}

// countLangSliceEntries is used as helper for Locale.AddLanguages() to check if
// provided lang keyword repeats in given slice.
// Returns count of keyword in given slice.
func (l *Locale) countLangSliceEntries(keyword string, slice []string) int {
	counter := 0

	for _, y := range slice {
		if keyword != y {
			continue
		}

		counter++
	}

	return counter
}

// HasLanguage can be used to check if Locale contains provided
// language.
// Returns true if language is initialized.
func (l *Locale) HasLanguage(langKey string) bool {
	_, err := l.GetLanguage(langKey)
	return err == nil
}

// EnabledLanguages can be used to get list/slice of enabled
// languages in current Locale structure.
func (l *Locale) EnabledLanguages() []string {
	// If there are no initialized languages then return nil rather
	// than empty, valid slice pointer.
	if len(l.Languages) == 0 {
		return nil
	}

	langs := make([]string, 0)

	for _, v := range l.Languages {
		langs = append(langs, v.Keyword)
	}

	return langs
}
