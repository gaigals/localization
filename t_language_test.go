package localization

import (
	"reflect"
	"strings"
	"testing"
)

func createTestLanguage(langKey string, translations TextMap) Language {
	return Language{Keyword: langKey, Map: translations}
}

func TestLanguage_ValueNoErrAndValuePluralNoErr(t *testing.T) {
	translations := TextMap{
		"key0": [2]string{"1 item", "2 items"},
		"key1": [2]string{"1 item", ""},
		"key2": [2]string{"", "2 items"},
	}

	language := createTestLanguage("en", translations)

	testCases := []struct {
		key             string
		searchNonPlural bool
		expected        string
	}{
		{"key0", true, "1 item"},
		{"key0", false, "2 items"},
		{"key1", true, "1 item"},
		{"key1", false, ""},
		{"key2", true, ""},
		{"key2", false, "2 items"},
		{"non_existing_key", true, ""},
		{"non_existing_key", false, ""},
		{"", true, ""},
		{"", false, ""},
	}

	for k, v := range testCases {
		text := ""

		if v.searchNonPlural {
			text = language.ValueNoErr(v.key)
		} else {
			text = language.ValuePluralNoErr(v.key)
		}

		if !strings.EqualFold(v.expected, text) {
			t.Fatalf("unexpected result, index=%d, expected=%s, actual=%s",
				k, v.expected, text)
		}
	}
}

func TestLanguage_SetValue(t *testing.T) {
	testCases := []struct {
		key      string
		value    string
		plural   string
		expected TextMap
	}{
		{"key0", "1 item", "2 items", TextMap{"key0": [2]string{"1 item", "2 items"}}},
		{"key0", "", "", TextMap{"key0": [2]string{"", ""}}},
		{"", "", "", TextMap{}}, // Check if empty key assignment is forbidden.
	}

	for k, v := range testCases {
		language := createTestLanguage("en", TextMap{})
		language.SetValue(v.key, v.value, v.plural)

		if !reflect.DeepEqual(v.expected, language.Map) {
			t.Fatalf("unexpected result, index=%d, expected=%v, actual=%v",
				k, v.expected, language.Map)
		}
	}
}
