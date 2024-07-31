package localization

import "fmt"

// TextMap is typedef for map which contains non-plural and plural
// translation. First element will always be non-plural.
type TextMap map[string][2]string

// Language contains specific language plural and non-plural translation.
type Language struct {
	Map     TextMap
	Keyword string
}

// ValueNoErr can be used to extract non-plural translation from language
// by providing translation keyword/key.
// No error check included (does not check if key exists).
func (l *Language) ValueNoErr(key string) string {
	return l.Map[key][0]
}

// Value can be used to extract non-plural translation from language
// by providing translation keyword/key.
// Returns error if key does not exist.
func (l *Language) Value(key string) (string, error) {
	_, exist := l.Map[key]
	if !exist {
		return "", fmt.Errorf("key '%s' does not exist", key)
	}

	return l.ValueNoErr(key), nil
}

// ValuePluralNoErr can be used to extract plural translation from language
// by providing translation keyword/key.
// No error check included (does not check if key exists).
func (l *Language) ValuePluralNoErr(key string) string {
	return l.Map[key][1]
}

// ValuePlural can be used to extract plural translation from language
// by providing translation keyword/key.
// Returns error if key does not exist.
func (l *Language) ValuePlural(key string) (string, error) {
	_, exist := l.Map[key]
	if !exist {
		return "", fmt.Errorf("key '%s' does not exist", key)
	}

	return l.ValuePluralNoErr(key), nil
}

// SetValue can be used to set non-plural and plural translation for language
// by providing translation keyword/key, non-plural value (value) and plural value (plural).
func (l *Language) SetValue(key, value, plural string) {
	// Do not allow empty key assignment.
	if key == "" {
		return
	}

	l.Map[key] = [2]string{value, plural}
}
