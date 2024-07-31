package localization

// Translate represents each language entry in YAML translate file.
type Translate struct {
	Key      string // Key/ID for translation.
	Language string // Language keyword ("lv", "en" etc.).
	Value    string // Translation ("some text").
	Plural   string // Translation in plural.
}
