package localization

import (
	"strings"
	"testing"
)

func TestText(t *testing.T) {
	locale0, _ := NewLocale(false, "lv", "en")
	locale0.SetValueNoErr("lv", "key0", "non-plural", "plural")
	locale0.SetValueNoErr("en", "key0", "en-non-plural", "en-plural")
	locale0.SetValueNoErr("en", "key1", "en-non-plural-1", "en-plural-1")

	testCases := []struct {
		langKey         string
		textKey         string
		expected        string
		strictUsage     bool
		failureExpected bool
	}{
		{"en", "key0", "en-non-plural", false, false},
		{"lv", "key0", "non-plural", false, false},
		{"lv", "key1", "en-non-plural-1", false, false},
		{"en", "key1", "en-non-plural-1", true, false},
		// Error - StrictUsage on and "lv" does not contain key "key1"
		{"lv", "key1", "", true, true},
		// Error - non existing language
		{"ee", "key0", "", false, true},
	}

	for k, v := range testCases {
		locale0.StrictUsage = v.strictUsage

		text, err := Text(*locale0, v.langKey, v.textKey)
		if err != nil && !v.failureExpected {
			t.Fatalf("unexpected failure, index=%d, error: %s", k, err)
		}
		if err == nil && v.failureExpected {
			t.Fatalf("expected failure, index=%d", k)
		}

		if !strings.EqualFold(v.expected, text) {
			t.Fatalf("unexpected result, index=%d, expected=%s, actual=%s",
				k, v.expected, text)
		}
	}
}

func TestTextPlural(t *testing.T) {
	locale0, _ := NewLocale(false, "lv", "en")
	locale0.SetValueNoErr("lv", "key0", "non-plural", "plural")
	locale0.SetValueNoErr("en", "key0", "en-non-plural", "en-plural")
	locale0.SetValueNoErr("en", "key1", "en-non-plural-1", "en-plural-1")

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
		{"en", "key0", "en-non-plural", false, false, false},
		{"lv", "key1", "en-plural-1", true, false, false},
		// Error - StrictUsage on and "lv" does not contain key "key1"
		{"lv", "key1", "", true, true, true},
		// Error - non existing language
		{"ee", "key0", "", true, false, true},
	}

	for k, v := range testCases {
		locale0.StrictUsage = v.strictUsage

		text, err := TextPlural(*locale0, v.langKey, v.textKey, v.askPlural)
		if err != nil && !v.failureExpected {
			t.Fatalf("unexpected failure, index=%d, error: %s", k, err)
		}
		if err == nil && v.failureExpected {
			t.Fatalf("expected failure, index=%d", k)
		}

		if !strings.EqualFold(v.expected, text) {
			t.Fatalf("unexpected result, index=%d, expected=%s, actual=%s",
				k, v.expected, text)
		}
	}
}

func TestTextf(t *testing.T) {
	locale0, _ := NewLocale(false, "lv", "en")
	locale0.SetValueNoErr("en", "key0", "%d %s", "")
	locale0.SetValueNoErr("en", "key1", "en-non-plural-1", "en-plural-1")

	testCases := []struct {
		langKey         string
		textKey         string
		dynamicInput    []interface{}
		expected        string
		strictUsage     bool
		failureExpected bool
	}{
		{"en", "key0", []interface{}{1, "text"}, "1 text", false, false},
		// Valid but broken output (no input passed).
		{"en", "key0", []interface{}{}, "%!d(MISSING) %!s(MISSING)", false, false},
		// Error - StrictUsage on and "lv" does not contain key "key1"
		{"lv", "key1", []interface{}{"something"}, "", true, true},
		// Error - non existing language
		{"ee", "key0", []interface{}{"something"}, "", true, true},
	}

	for k, v := range testCases {
		locale0.StrictUsage = v.strictUsage

		text, err := Textf(*locale0, v.langKey, v.textKey, v.dynamicInput...)
		if err != nil && !v.failureExpected {
			t.Fatalf("unexpected failure, index=%d, error: %s", k, err)
		}
		if err == nil && v.failureExpected {
			t.Fatalf("expected failure, index=%d", k)
		}

		if !strings.EqualFold(v.expected, text) {
			t.Fatalf("unexpected result, index=%d, expected=%s, actual=%s",
				k, v.expected, text)
		}
	}
}

func TestTextPluralf(t *testing.T) {
	locale0, _ := NewLocale(false, "lv", "en")
	locale0.SetValueNoErr("lv", "key0", "%d non-plural", "%d plural")
	locale0.SetValueNoErr("en", "key0", "%d en-non-plural", "%d en-plural")
	locale0.SetValueNoErr("en", "key1", "en-non-plural-1", "en-plural-1")

	testCases := []struct {
		langKey         string
		textKey         string
		dynamicInput    []interface{}
		expected        string
		askPlural       bool
		strictUsage     bool
		failureExpected bool
	}{
		{"lv", "key0", []interface{}{1}, "1 non-plural", false, false, false},
		{"en", "key0", []interface{}{2}, "2 en-plural", true, false, false},
		// Valid but broken output - no input passed.
		{"en", "key0", []interface{}{}, "%!d(MISSING) en-non-plural", false, false, false},
		// Error - StrictUsage on and "lv" does not contain key "key1"
		{"lv", "key1", []interface{}{"something"}, "", false, true, true},
		// Error - non existing language
		{"ee", "key0", []interface{}{"something"}, "", false, true, true},
	}

	for k, v := range testCases {
		locale0.StrictUsage = v.strictUsage

		text, err := TextPluralf(*locale0, v.langKey, v.textKey, v.askPlural, v.dynamicInput...)
		if err != nil && !v.failureExpected {
			t.Fatalf("unexpected failure, index=%d, error: %s", k, err)
		}
		if err == nil && v.failureExpected {
			t.Fatalf("expected failure, index=%d", k)
		}

		if !strings.EqualFold(v.expected, text) {
			t.Fatalf("unexpected result, index=%d, expected=%s, actual=%s",
				k, v.expected, text)
		}
	}
}

func TestTextPluralIntf(t *testing.T) {
	locale0, _ := NewLocale(false, "lv", "en")
	locale0.SetValueNoErr("lv", "key0", "%d lieta %s", "%d lietas %s")
	locale0.SetValueNoErr("en", "key0", "%d item %s", "%d items %s")
	locale0.SetValueNoErr("en", "key1", "en-non-plural-1", "en-plural-1")

	testCases := []struct {
		langKey         string
		textKey         string
		dynamicInput    []interface{}
		expected        string
		strictUsage     bool
		failureExpected bool
	}{
		{"lv", "key0", []interface{}{1, ":("}, "1 lieta :(", false, false},
		{"en", "key0", []interface{}{2, ":)"}, "2 items :)", false, false},
		// Error - first input element must be int type.
		{"en", "key0", []interface{}{"text", 1}, "", false, true},
		// Error - no dynamic input
		{"en", "key0", []interface{}{}, "", false, true},
		// Error - StrictUsage on and "lv" does not contain key "key1"
		{"lv", "key1", []interface{}{1}, "", true, true},
		// Error - non existing language
		{"ee", "key0", []interface{}{1}, "", true, true},
	}

	for k, v := range testCases {
		locale0.StrictUsage = v.strictUsage

		text, err := TextPluralIntf(*locale0, v.langKey, v.textKey, v.dynamicInput...)
		if err != nil && !v.failureExpected {
			t.Fatalf("unexpected failure, index=%d, error: %s", k, err)
		}
		if err == nil && v.failureExpected {
			t.Fatalf("expected failure, index=%d", k)
		}

		if !strings.EqualFold(v.expected, text) {
			t.Fatalf("unexpected result, index=%d, expected=%s, actual=%s",
				k, v.expected, text)
		}
	}
}
