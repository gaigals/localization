package localization

import (
	"github.com/spf13/cast"
	"regexp"
	"sort"
	"strings"
)

// regex for finding languages groups with or without priority, for example:
// "en-US,en;q=0.5"
var regexLanguagesByWeight = regexp.MustCompile(`(((([\*]{1})|(([a-z]{2}))+(-)+([a-zA-Z]{2})|([a-z]{2}))+(,?)+( ?)){1,})+(;?)+((q=[0-9\.]{1,})?)`)

// regexWeight is for finding weight/priority, for example: "q=0.5"
var regexWeight = regexp.MustCompile(`(q=[0-9\.]{1,})`)

// regexLang is used for finding languages, for example: "*", "en-US", "en", "lv"
var regexLang = regexp.MustCompile(`(([\*]{1})|([a-z]{2})+(-)+([A-Z]{2})|([a-z]{2}))`)

// PriorityGroup holds information about specific Language priority
// group.
// Each group contains specific weight (if defined) and list of languages
// inside in that priority.
type PriorityGroup struct {
	Weight    float32  // Language group weight/priority
	Languages []string // Slice of languages in current priority group.
}

// AcceptLanguages is used to read and contain information about request
// header Accept-Language data.
// To parse Request header Accept-Language use localization.ParseAcceptLanguage
// which will return AcceptLanguages structure with parsed data.
type AcceptLanguages struct {
	LangGroups []PriorityGroup // For each weight 1 language group
}

// ParseAcceptLanguage can be used to parsed passed request header Accept-Language
// into AcceptLanguage structure.
// Value MUST BE Accept-Language value.
// Returns pointer to AcceptLanguages structure or error.
//
// Information about Accept-Language ...
// Docs:
// https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Accept-Language
// Some examples:
// https://www.holisticseo.digital/technical-seo/http-header/content-negotiation/accept-language/
func ParseAcceptLanguage(value string) *AcceptLanguages {
	if value == "" {
		// Valid pointer is required for later usage (calling
		// method without panicking).
		return &AcceptLanguages{}
	}

	acceptLangs := AcceptLanguages{}

	// Split languages by weight.
	acceptLangs.splitByWeight(value)

	// Order parsed accept-languages by provided weights.
	if len(acceptLangs.LangGroups) > 0 {
		acceptLangs.orderByWeight()
	}

	return &acceptLangs
}

// FindFirstMatchingLang can be used to extract first matching language
// in AcceptLanguages by proving prioritized enabled languages slice.
// If AcceptLanguage contains '*' (asterisk) then it will be replaced with
// passed default language.
// Returns first matching language key or empty string on no match.
//
// Params:
// enabledLanguages - prioritized slice of enabled languages.
// defaultLang - default language which will be used in asterisk cases.
func (al *AcceptLanguages) FindFirstMatchingLang(enabledLanguages []string, defaultLang string) string {
	tempLangGroups := al.replaceAsterisks(defaultLang)

	for _, v := range tempLangGroups {
		for _, y := range v.Languages {
			for _, z := range enabledLanguages {
				if strings.EqualFold(y, z) {
					return z
				}
			}
		}
	}

	return ""
}

// splitByWeight split given Accept-Language string by weight.
// Returns error if something went wrong.
func (al *AcceptLanguages) splitByWeight(value string) {
	al.LangGroups = make([]PriorityGroup, 0)

	dataSlices := regexLanguagesByWeight.FindAll([]byte(value), -1)

	for _, v := range dataSlices {
		langList := PriorityGroup{}

		langList.parse(v)
		al.LangGroups = append(al.LangGroups, langList)
	}
}

// orderByWeight order generated PriorityGroups by weight (1 -> 0).
func (al *AcceptLanguages) orderByWeight() {
	sort.Slice(al.LangGroups, func(i, j int) bool {
		return al.LangGroups[i].Weight > al.LangGroups[j].Weight
	})
}

// parse is used to parse given priority group string into
// actual PriorityGroup structure/values.
func (l *PriorityGroup) parse(value []byte) {
	l.extractWeight(value)

	langs := regexLang.FindAll(value, -1)

	for _, v := range langs {
		l.Languages = append(l.Languages, string(v))
	}
}

// replaceAsterisks is used to replace all asterisks inside AcceptLanguages
// structure with provided default language.
// Returns modified slice of PriorityGroup.
func (al *AcceptLanguages) replaceAsterisks(defaultLang string) []PriorityGroup {
	langGroups := make([]PriorityGroup, len(al.LangGroups))

	for k, v := range al.LangGroups {
		langGroups[k] = v

		for x, y := range v.Languages {
			langGroups[k].Languages[x] = strings.ReplaceAll(y, "*", defaultLang)
		}
	}

	return langGroups
}

// extractWeight is used to extract current PriorityGroup weight from
// passed value bytes.
func (l *PriorityGroup) extractWeight(value []byte) {
	weight := string(regexWeight.Find(value))
	if weight == "" {
		return
	}

	weight = strings.Replace(weight, "q=", "", 1)
	l.Weight = cast.ToFloat32(weight)
}
