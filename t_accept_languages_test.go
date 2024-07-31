package localization

import (
	"reflect"
	"testing"
)

func TestParseAcceptLanguage(t *testing.T) {
	testCases := []struct {
		value    string
		expected *AcceptLanguages
	}{
		{
			"en-US,en;q=0.5",
			&AcceptLanguages{
				[]PriorityGroup{
					{0.5, []string{"en-US", "en"}},
				},
			},
		},
		{
			"*",
			&AcceptLanguages{
				[]PriorityGroup{
					{.0, []string{"*"}},
				},
			},
		},
		{
			"fr-CH, fr;q=0.9, en;q=0.8, de;q=0.7, *;q=0",
			&AcceptLanguages{
				[]PriorityGroup{
					{0.9, []string{"fr-CH", "fr"}},
					{0.8, []string{"en"}},
					{0.7, []string{"de"}},
					{.0, []string{"*"}},
				},
			},
		},
		{
			"*;q=0",
			&AcceptLanguages{
				[]PriorityGroup{
					{.0, []string{"*"}},
				},
			},
		},
		{
			"*",
			&AcceptLanguages{
				[]PriorityGroup{
					{.0, []string{"*"}},
				},
			},
		},
		{
			"", &AcceptLanguages{nil},
		},
	}

	for k, v := range testCases {
		lang := ParseAcceptLanguage(v.value)

		if !reflect.DeepEqual(v.expected, lang) {
			t.Fatalf("unexpected result, index=%d, input=%s, expected=%v, actual=%v",
				k, v.value, v.expected, lang)
		}
	}
}

func TestAcceptLanguages_FindFirstMatchingLang(t *testing.T) {
	testCases := []struct {
		enabledLangs    []string
		defaultLang     string
		acceptLanguages AcceptLanguages
		expected        string
	}{
		{
			[]string{"lv", "en"}, "lv",
			AcceptLanguages{
				[]PriorityGroup{
					{0.5, []string{"en-US", "en"}},
				},
			},
			"en",
		},
		{
			[]string{"lv", "en"}, "lv",
			AcceptLanguages{
				[]PriorityGroup{
					{0.5, []string{"ee", "lt"}},
					{0.4, []string{"*"}},
					{.0, []string{"es"}},
				},
			},
			"lv",
		},
		{
			[]string{"lv", "en"}, "en",
			AcceptLanguages{
				[]PriorityGroup{
					{0.5, []string{"ee", "lt"}},
					{0.4, []string{"*"}},
					{.0, []string{"lv"}},
				},
			},
			"en",
		},
		{
			[]string{"lv", "en"}, "lv",
			AcceptLanguages{
				[]PriorityGroup{
					{0.5, []string{"en-US", "en"}},
					{0.4, []string{"ee"}},
					{.0, []string{"lt"}},
				},
			},
			"en",
		},
		{
			[]string{"de", "fr"}, "lv",
			AcceptLanguages{
				[]PriorityGroup{
					{0.5, []string{"en-US", "en"}},
					{0.4, []string{"ee"}},
					{.0, []string{"lt"}},
				},
			},
			"",
		},
	}

	for k, v := range testCases {
		lang := v.acceptLanguages.FindFirstMatchingLang(v.enabledLangs, v.defaultLang)

		if v.expected != lang {
			t.Fatalf("unexpected result, index=%d, expected=%s, actual=%s",
				k, v.expected, lang)
		}
	}

}
