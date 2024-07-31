package localization

import (
	"fmt"
	"reflect"
	"strings"
	"testing"
)

func TestLocale_buildPrioritizedLanguageList(t *testing.T) {
	locale0, _ := NewLocale(false, "en", "lv", "lt", "ee")
	langEN, _ := locale0.GetLanguage("en")
	langLV, _ := locale0.GetLanguage("lv")
	langLT, _ := locale0.GetLanguage("lt")
	langEE, _ := locale0.GetLanguage("ee")

	testCases := []struct {
		inputLang   *Language
		strictUsage bool
		expected    []*Language
	}{
		// Check if only input gets returned on StrictUsage=true
		{langEN, true, []*Language{langEN}},
		// Check if correctly ordered slice gets returned (langEN must be first).
		{langEN, false, []*Language{langEN, langLV, langLT, langEE}},
		// Check if only input gets returned on StrictUsage=true
		{langLV, true, []*Language{langLV}},
		// Check if correctly ordered slice gets returned (langLV must be first).
		{langLV, false, []*Language{langLV, langEN, langLT, langEE}},
		// Check if correctly ordered slice gets returned (langLT must be first).
		{langLT, false, []*Language{langLT, langEN, langLV, langEE}},
	}

	for k, v := range testCases {
		locale0.StrictUsage = v.strictUsage

		list := locale0.buildPrioritizedLanguageList(v.inputLang)

		// Validate each pointer value. Each slice element pointer must match with
		// v.expected slice pointer values.
		for x := range list {
			if list[x] != v.expected[x] {
				t.Fatalf("unexpected pointer, index=%d.%d, expected=%p, actual=%p",
					k, x, v.expected[k], list[k])
			}
		}
	}
}

func TestLocale_AddLanguages(t *testing.T) {
	testCases := []struct {
		existingLangs   []string
		input           []string
		expected        []string
		failureExpected bool
	}{
		// No error
		{[]string{}, []string{"en"}, []string{"en"}, false},
		// No error
		{[]string{}, []string{"en", "lv"}, []string{"en", "lv"}, false},
		// No error
		{[]string{}, []string{"en", "lv", "lt", "ee"}, []string{"en", "lv", "lt", "ee"}, false},
		// No error
		{[]string{}, []string{}, []string{}, false},
		// Error - 'lv' redefined in params
		{[]string{}, []string{"lv", "lv"}, []string{}, true},
		// Error - 'ee' redefined in params
		{[]string{}, []string{"en", "ee", "ee"}, []string{}, true},
		// Error - 'lv' already initialized in language list.
		{[]string{"lv"}, []string{"lv"}, []string{"lv"}, true},
		// No error
		{[]string{"lv"}, []string{"en"}, []string{"lv", "en"}, false},
		// No error
		{[]string{"lv", "en"}, []string{"ee", "lt"}, []string{"lv", "en", "ee", "lt"}, false},
	}

	for k, v := range testCases {
		locale0 := Locale{}
		if len(v.existingLangs) != 0 {
			_ = locale0.AddLanguages(v.existingLangs...)
		}

		err := locale0.AddLanguages(v.input...)
		if err != nil && !v.failureExpected {
			t.Fatalf("unexpected failure, index=%d, error: %s", k, err)
		}
		if err == nil && v.failureExpected {
			t.Fatalf("expected failure, index=%d", k)
		}

		if len(locale0.Languages) != len(v.expected) {
			t.Fatalf("unexpected language length for locale, index=%d, expected=%d, actual=%d",
				k, len(locale0.Languages), len(v.expected))
		}

		for x, y := range locale0.Languages {
			if !strings.EqualFold(y.Keyword, v.expected[x]) {
				t.Fatalf("unexpected keyword, index=%d.%d, expected=%s, actual=%s",
					k, x, v.expected[x], y.Keyword)
			}
		}
	}
}

func TestNewLocale(t *testing.T) {
	testCases := []struct {
		strictUsage     bool
		langs           []string
		expected        *Locale
		failureExpected bool
	}{
		{
			true,
			[]string{"lv"},
			&Locale{StrictUsage: true, Languages: []Language{{Keyword: "lv"}}},
			false,
		},
		{
			false,
			[]string{"lv"},
			&Locale{StrictUsage: false, Languages: []Language{{Keyword: "lv"}}},
			false,
		},
		{
			true,
			[]string{},
			&Locale{StrictUsage: true},
			false,
		},
		{
			true,
			[]string{"lv", "lv"},
			nil,
			true,
		},
	}

	for k, v := range testCases {
		locale0, err := NewLocale(v.strictUsage, v.langs...)
		if err != nil && !v.failureExpected {
			t.Fatalf("unexpected failure, index=%d, error: %s", k, err)
		}
		if err == nil && v.failureExpected {
			t.Fatalf("expected failure, index=%d", k)
		}

		if err != nil && locale0 != nil {
			t.Fatalf("returned *Locale should be nill")
		}

		if err != nil {
			continue
		}

		// if !reflect.DeepEqual(*v.expected, locale0) {
		//	 t.Fatalf("unexpected result, index=%d, expected=%+v, actual=%+v",
		//		k, *v.expected, *locale0)
		// }

		// ^^^^^^^^^^^^^^^^^^^^^^^
		// Dos not work :(

		if v.expected.StrictUsage != locale0.StrictUsage {
			t.Fatalf("unexpected Locale.StrictUsage value, index=%d, expected=%v, actual=%v",
				k, v.expected.StrictUsage, locale0.StrictUsage)
		}

		if len(v.expected.Languages) != len(locale0.Languages) {
			t.Fatalf("unexpected Language length, index=%d, expected=%d, actual=%d",
				k, len(v.expected.Languages), len(locale0.Languages))
		}

		for x, y := range v.expected.Languages {
			if len(y.Map) != len(locale0.Languages[x].Map) {
				t.Fatalf("unexpected Language map length, index=%d.%d, expected=%d, actual=%d",
					k, x, len(y.Map), len(locale0.Languages[x].Map))
			}
		}
	}
}

func TestLocale_ValueAndValuePlural(t *testing.T) {
	locale0, _ := NewLocale(false, "lv", "en")
	locale0.SetValueNoErr("lv", "key0", "non-plural", "plural")
	locale0.SetValueNoErr("en", "key0", "en-non-plural", "en-plural")
	locale0.SetValueNoErr("lv", "key1", "", "plural")
	locale0.SetValueNoErr("lv", "key2", "non-plural", "")
	locale0.SetValueNoErr("lv", "key3", "", "")
	locale0.SetValueNoErr("en", "key4", "en-non-plural", "en-plural")

	testCases := []struct {
		langKey         string
		textKey         string
		expected        string
		askPlural       bool
		strictUsage     bool
		failureExpected bool
	}{
		{"lv", "key0", "non-plural", false, false, false},
		{"lv", "key0", "plural", true, false, false},
		{"lv", "key1", "", false, false, false},
		{"lv", "key1", "plural", true, false, false},
		{"lv", "key2", "non-plural", false, false, false},
		{"lv", "key2", "", true, false, false},
		{"lv", "key3", "", false, false, false},
		{"lv", "key3", "", true, false, false},
		// Strict usage - false. EN translations used as backup.
		{"lv", "key4", "en-non-plural", false, false, false},
		{"lv", "key4", "en-plural", true, false, false},
		// Strict usage - true. Other language key backups are restricted, expected error - key does not exist.
		{"lv", "key4", "", false, true, true},
		{"lv", "key4", "", true, true, true},
		// Error - unknown language keyword
		{"ee", "key0", "", false, false, true},
		// Error - key does not exist in any language.
		{"lv", "non_existing_key", "", false, false, true},
		{"lv", "non_existing_key", "", true, false, true},
	}

	for k, v := range testCases {
		locale0.StrictUsage = v.strictUsage

		var err error
		text := ""

		if v.askPlural {
			text, err = locale0.ValuePlural(v.langKey, v.textKey)
		} else {
			text, err = locale0.Value(v.langKey, v.textKey)
		}

		if err != nil && !v.failureExpected {
			t.Fatalf("unexpected error, index=%d, error: %s", k, err)
		}
		if err == nil && v.failureExpected {
			t.Fatalf("expected error, index=%d", k)
		}

		if !strings.EqualFold(v.expected, text) {
			t.Fatalf("unexpected result, index=%d, expected=%s, actual=%s", k, v.expected, text)
		}
	}

	// Reset locale
	locale0 = &Locale{}

	// Check if error gets returned when languages aren't initialized.
	_, err := locale0.Value("lv", "key")
	if err == nil {
		t.Fatalf("expected error")
	}

	// Check if error gets returned when languages aren't initialized.
	_, err = locale0.ValuePlural("lv", "key")
	if err == nil {
		t.Fatalf("expected error")
	}
}

func TestLocale_SetValue(t *testing.T) {
	testCases := []struct {
		activeLangs     []string
		langKey         string
		textKey         string
		value           string
		plural          string
		failureExpected bool
	}{
		{[]string{"lv", "en"}, "lv", "key0", "non-plural", "plural", false},
		{[]string{"lv", "en"}, "lv", "key1", "", "plural", false},
		{[]string{"lv", "en"}, "lv", "key2", "non-plural", "", false},
		// Error - lang key does not exist
		{[]string{"lv", "en"}, "ee", "key3", "non-plural", "", true},
		// Error - text key is empty string
		{[]string{"lv", "en"}, "ee", "", "non-plural", "", true},
		// Error - no languages initialized
		{[]string{}, "lv", "key0", "non-plural", "", true},
	}

	for k, v := range testCases {
		locale0 := &Locale{StrictUsage: true}

		if len(v.activeLangs) > 0 {
			locale0, _ = NewLocale(true, v.activeLangs...)
		}

		err := locale0.SetValue(v.langKey, v.textKey, v.value, v.plural)

		if err != nil && !v.failureExpected {
			t.Fatalf("unexpected error, index=%d, error: %s", k, err)
		}
		if err == nil && v.failureExpected {
			t.Fatalf("expected error, index=%d", k)
		}

		if v.failureExpected {
			continue
		}

		val, err := locale0.Value(v.langKey, v.textKey)
		if err != nil {
			t.Fatalf("unexpected error, index=%d, erorr: %s", k, err)
		}
		if !strings.EqualFold(v.value, val) {
			t.Fatalf(
				"unexpected non-plural value, index=%d, expected=%s, actual=%s",
				k,
				v.value,
				val,
			)
		}

		val, err = locale0.ValuePlural(v.langKey, v.textKey)
		if err != nil {
			t.Fatalf("unexpected error, index=%d, erorr: %s", k, err)
		}
		if !strings.EqualFold(v.plural, val) {
			t.Fatalf("unexpected plural value, index=%d, expected=%s, actual=%s", k, v.plural, val)
		}
	}
}

func TestLocale_AddTranslate(t *testing.T) {
	locale0, _ := NewLocale(true, "lv", "en")

	testCases := []struct {
		input           []Translate
		failureExpected bool
	}{
		{
			[]Translate{
				{Key: "key0", Language: "lv", Value: "non-plural", Plural: "plural"},
			},
			false,
		},
		{
			nil,
			false,
		},
		{
			[]Translate{
				{Key: "key0", Language: "lv", Value: "non-plural", Plural: "plural"},
				{Key: "key1", Language: "lv", Value: "non-plural_1", Plural: "plural_1"},
				{Key: "key0", Language: "en", Value: "en_non-plural_1", Plural: "en_plural_1"},
				{Key: "key1", Language: "en", Value: "en_non-plural_1", Plural: "en_plural_1"},
			},
			false,
		},
		{ // Error - language key does not exist.
			[]Translate{
				{Key: "key5", Language: "lt", Value: "non-plural", Plural: "plural"},
			},
			true,
		},
	}

	for k, v := range testCases {
		err := locale0.AddTranslate(v.input...)

		if err != nil && !v.failureExpected {
			t.Fatalf("unexpected error, index=%d, error: %s", k, err)
		}
		if err == nil && v.failureExpected {
			t.Fatalf("expected error, index=%d", k)
		}

		if v.failureExpected {
			continue
		}

		for x, y := range v.input {
			text, err := locale0.Value(y.Language, y.Key)
			if err != nil {
				t.Fatalf("unexpected error, index=%d.%d, erorr: %s", k, x, err)
			}
			if !strings.EqualFold(y.Value, text) {
				t.Fatalf(
					"unexpected non-plural value, index=%d.%d, expected=%s, actual=%s",
					k,
					x,
					y.Value,
					text,
				)
			}

			text, err = locale0.ValuePlural(y.Language, y.Key)
			if err != nil {
				t.Fatalf("unexpected error, index=%d.%d, erorr: %s", k, x, err)
			}
			if !strings.EqualFold(y.Plural, text) {
				t.Fatalf(
					"unexpected plural value, index=%d.%d, expected=%s, actual=%s",
					k,
					x,
					y.Plural,
					text,
				)
			}
		}
	}
}

func TestLocale_AddYAMLFile(t *testing.T) {
	locale0, _ := NewLocale(false, "lv", "en")

	testCases := []struct {
		input           []*YAMLFile
		failureExpected bool
	}{
		{
			[]*YAMLFile{
				{FilePath: "test.yaml", Translates: []Translate{
					{Key: "key0", Language: "lv", Value: "non_plural", Plural: "plural"},
					{Key: "key0", Language: "en", Value: "non_plural", Plural: "plural"},
				}},
			},
			false,
		},
		{
			nil,
			false,
		},
		{
			[]*YAMLFile{nil},
			true,
		},
		{
			[]*YAMLFile{
				{FilePath: "test.yaml", Translates: []Translate{
					{Key: "key1", Language: "lt", Value: "non_plural", Plural: "plural"},
					{Key: "key1", Language: "ee", Value: "non_plural", Plural: "plural"},
				}},
			},
			true,
		},
	}

	for k, v := range testCases {
		err := locale0.AddYAMLFile(v.input...)

		if err != nil && !v.failureExpected {
			t.Fatalf("unexpected error, index=%d, error: %s", k, err)
		}
		if err == nil && v.failureExpected {
			t.Fatalf("expected error, index=%d", k)
		}

		if v.failureExpected {
			continue
		}

		for x, y := range v.input {
			for z, w := range y.Translates {
				text, err := locale0.Value(w.Language, w.Key)
				if err != nil {
					t.Fatalf("unexpected error, index=%d.%d.%d, erorr: %s", k, x, z, err)
				}
				if !strings.EqualFold(w.Value, text) {
					t.Fatalf(
						"unexpected non-plural value, index=%d.%d.%d, expected=%s, actual=%s",
						k,
						x,
						z,
						w.Value,
						text,
					)
				}

				text, err = locale0.ValuePlural(w.Language, w.Key)
				if err != nil {
					t.Fatalf("unexpected error, index=%d.%d.%d, erorr: %s", k, x, z, err)
				}
				if !strings.EqualFold(w.Plural, text) {
					t.Fatalf(
						"unexpected plural value, index=%d.%d.%d, expected=%s, actual=%s",
						k,
						x,
						z,
						w.Plural,
						text,
					)
				}
			}
		}
	}
}

func TestLocale_LoadYAMLFile(t *testing.T) {
	tempDir := t.TempDir()

	testCases := []struct {
		fileNames       []string
		fileContents    []string
		defaultLang     string
		createFiles     bool
		failureExpected bool
		expected        Locale
	}{
		{ // No errors - file exist and content matches.
			[]string{"file_0.yaml"},
			[]string{"key0: \"non_plural\"\n"},
			"en",
			true,
			false,
			Locale{
				[]Language{
					{TextMap{}, "lv"},
					{TextMap{"key0": [2]string{"non_plural", ""}}, "en"},
				}, true,
			},
		},
		{ // No errors - file exist and content matches.
			[]string{"file_0.yaml"},
			[]string{"key0: \"non_plural\"\n"},
			"lv",
			true,
			false,
			Locale{
				[]Language{
					{TextMap{"key0": [2]string{"non_plural", ""}}, "lv"},
					{TextMap{}, "en"},
				}, true,
			},
		},
		{
			[]string{"file_0.yaml"},
			[]string{
				"key0:\n  - lv:\n    - \"non_plural\"\n    - \"plural\"\n" +
					"  - en:\n    - \"en_non_plural\"\n    - \"en_plural\"\n",
			},
			"lv",
			true,
			false,
			Locale{
				[]Language{
					{TextMap{"key0": [2]string{"non_plural", "plural"}}, "lv"},
					{TextMap{"key0": [2]string{"en_non_plural", "en_plural"}}, "en"},
				}, true,
			},
		},
		{ // No errors - file exist and content matches.
			[]string{"file_0.yaml", "file_1.yaml"},
			[]string{"key0: \"non_plural\"\n", "key1: \"non_plural_1\"\n"},
			"en",
			true,
			false,
			Locale{
				[]Language{
					{TextMap{}, "lv"},
					{
						TextMap{
							"key0": [2]string{"non_plural", ""},
							"key1": [2]string{"non_plural_1", ""},
						}, "en",
					},
				}, true,
			},
		},
		{ // No errors - file exist and content matches.
			[]string{"file_0.yaml"},
			[]string{},
			"lv",
			false,
			true,
			Locale{},
		},
		{ // No errors - file exist and content matches.
			[]string{"file_0.yaml"},
			[]string{"key0: \"non_plural\"\n"},
			"ee",
			true,
			true,
			Locale{},
		},
	}

	for k, v := range testCases {
		filePaths := make([]string, len(v.fileNames))

		if v.createFiles {
			for x, y := range v.fileNames {
				err := createTempFile(tempDir, y, v.fileContents[x])
				if err != nil {
					t.Fatalf("unexpected error, index=%d, error: %s", k, err)
				}

				filePaths[x] = fmt.Sprintf("%s/%s", tempDir, y)
			}
		}

		locale0, _ := NewLocale(true, "lv", "en")

		err := locale0.LoadYAMLFile(v.defaultLang, filePaths...)
		if err != nil && !v.failureExpected {
			t.Fatalf("unexpected error, index=%d, error: %s", k, err)
		}
		if err == nil && v.failureExpected {
			t.Fatalf("expected error, index=%d, error: %s", k, err)
		}

		if v.failureExpected {
			continue
		}

		if !reflect.DeepEqual(v.expected, *locale0) {
			t.Fatalf("unexpected result, index=%d, expected=%+v, actual=%+v",
				k, v.expected, *locale0)
		}
	}
}

func TestLocale_HasLanguage(t *testing.T) {
	testCases := []struct {
		langList []string
		langKey  string
		expected bool
	}{
		{[]string{"lv", "en"}, "lv", true},
		{[]string{"lv"}, "lv", true},
		{[]string{"en"}, "en", true},
		{[]string{"ee", "lt"}, "lv", false},
		{nil, "lv", false},
		{[]string{}, "lv", false},
	}

	for k, v := range testCases {
		locale, err := NewLocale(true, v.langList...)
		if err != nil {
			t.Fatalf("index=%d, unexpected error: %s", k, err)
		}

		hasLang := locale.HasLanguage(v.langKey)
		if v.expected != hasLang {
			t.Fatalf("unexpected result, index=%d, expected=%v, actual=%v",
				k, v.expected, hasLang)
		}
	}
}

func TestLocale_ValueNoErr_ValuePluralNoErr(t *testing.T) {
	testCases := []struct {
		data             []Language
		initializedLangs []string
		langKey          string
		key              string
		plural           bool
		strictUsage      bool
		expected         string
	}{
		{
			[]Language{
				{TextMap{"key0": [2]string{"non_plural", ""}}, "lv"},
			},
			[]string{"lv"},
			"lv", "key0", false, false, "non_plural",
		},
		{
			[]Language{
				{TextMap{"key0": [2]string{"non_plural", ""}}, "lv"},
				{TextMap{"key0": [2]string{"non_plural_en", ""}}, "en"},
			},
			[]string{"lv", "en"},
			"en", "key0", false, false, "non_plural_en",
		},
		{
			[]Language{
				{TextMap{"key0": [2]string{"non_plural", "plural_lv"}}, "lv"},
				{TextMap{"key0": [2]string{"non_plural_en", "plural_en"}}, "en"},
			},
			[]string{"lv", "en"},
			"en", "key0", true, false, "plural_en",
		},
		{
			// Expected no results -> strict usage - True and key does not exist
			[]Language{
				{TextMap{"key1": [2]string{"non_plural", ""}}, "lv"}, // Different key
				{TextMap{"key0": [2]string{"non_plural_en", ""}}, "en"},
			},
			[]string{"lv", "en"},
			"lv", "key0", false, true, "",
		},
		{
			// Expected result in EN -> lv key does not exist and strict usage is - False.
			[]Language{
				{TextMap{"key1": [2]string{"non_plural", ""}}, "lv"}, // Different key
				{TextMap{"key0": [2]string{"non_plural_en", ""}}, "en"},
			},
			[]string{"lv", "en"},
			"lv", "key0", false, false, "non_plural_en",
		},
		{
			// Expected no result -> lang does not exist
			[]Language{
				{TextMap{"key0": [2]string{"non_plural", ""}}, "lv"},
				{TextMap{"key0": [2]string{"non_plural_en", ""}}, "en"},
			},
			[]string{"lv", "en"},
			"lt", "key0", false, false, "",
		},
		{
			// Expected no result -> no translations available.
			[]Language{},
			[]string{"lv", "en"},
			"lv", "key0", false, false, "",
		},
		{
			// Expected no result -> no translations available.
			nil,
			[]string{"lv", "en"},
			"lv", "key0", false, false, "",
		},
	}

	for k, v := range testCases {
		locale, err := NewLocale(v.strictUsage, v.initializedLangs...)
		if err != nil {
			t.Fatalf("index=%d, unexpected error: %s", k, err)
		}

		locale.Languages = v.data

		value := ""

		if !v.plural {
			value = locale.ValueNoErr(v.langKey, v.key)
		} else {
			value = locale.ValuePluralNoErr(v.langKey, v.key)
		}

		if value != v.expected {
			t.Fatalf("unexpected result, index=%d, expected=%s, actual=%s",
				k, v.expected, value)
		}
	}
}

func TestLocale_EnabledLanguages(t *testing.T) {
	testCases := []struct {
		enabledLangs []string
	}{
		{[]string{"lv", "en"}},
		{[]string{"lv", "en", "ee", "lt", "en-US"}},
		{[]string{"lv"}},
		{[]string{""}},
		{nil},
	}

	for k, v := range testCases {
		locale, err := NewLocale(true, v.enabledLangs...)
		if err != nil {
			t.Fatalf("index=%d, unexpected error: %s", k, err)
		}

		langs := locale.EnabledLanguages()

		if !reflect.DeepEqual(v.enabledLangs, langs) {
			t.Fatalf("unexpected results, index=%d, expected=%v, actual=%v",
				k, v.enabledLangs, langs)
		}
	}
}
