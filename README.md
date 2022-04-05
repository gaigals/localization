# Localization example

__Valodu pievienošana:__
```go
var locale Locale

locale.AddLanguages("lv", "en")

// Enable this is translation is required for every language key
// By default, if key is missing then localizator will try to find key in other languages.
locale.StrictUsage = true
```

__Manuālā teksta/tulkojumu pievienošana:__
```go

// func SetValueNoErr(langKey, key, text, pluralText string)
// func SetValue(langKey, key, text, pluralText string) error

// Set translate without error (checks are done by caller application).
locale.SetValueNoErr("lv", "hello_world", "Sveicināta, Pasaule!", "")

// Set translate and return error if something goes wrong.
err := locale.SetValue("en", "hello_world", "Hello, World!", "")
if err != nil {
	// Error handling
}
```

__Tulkojuma ielādēšana no `YAML` faila:__
```go
// LoadYAMLTranslateFile(filePath, defaultLanguage) (*YAMLFile, error)

// Load targeted YAML file
translates, err := LoadYAMLTranslateFile(path, "lv")
if err != nil {
    // Error handling
}
```

__Tulkojumu pievienošana izmantojot `*YAMLFile`:__
```go
// func AddYAMLFile(yamlFiles ...*YAMLFile) error

// Add 1 or more yaml file translates
err := locale.AddYAMLFile(yamlFile)
if err != nil {
    // Error handling
}
```


#### Templates:

Piejamās template funckijas:
```go
// Text returns translation by using provided langKey and textKey.
func Text(locale Locale, langKey, textKey string) (string, error)

// TextPlural returns plural or non-plural translation by provided langKey and textKey (hard-coded plural text).
func TextPlural(locale Locale, langKey, textKey string, isPlural bool) (string, error)

// TextDynamic returns translation by using provided langKey and textKey, and applies dynamic input.
func TextDynamic(locale Locale, langKey, textKey string, input ...interface{}) (string, error)

// TextDynamic returns plural or non-plural translation by using provided langKey and textKey, and applies dynamic input.
func TextPluralDynamic(locale Locale, langKey, textKey string, input ...interface{}) (string, error)
```


Pilns template piemērs:
```html
{{ $Page := . }}

<!--Load page locale-->
{{ $Locale := (GetMapValue $Page "Locale") }}
<!--Get language keyword ("en" or "lv")-->
{{ $Lang := (GetMapValue $Page "Lang") }}

<!-- ... -->

<!--Basic translation loading-->
{{ Text $Locale $Lang "hello_world" }}

<!--Loading hard-coded plural translation-->
{{ TextPlural $Locale $Lang "plural_text_hard_coded" false }}

<!--Loading plural dynamic text, param must be int -->
{{ TextPluralDynamic $Locale $Lang "plural_text_dynamic" 1 }}
{{ TextPluralDynamic $Locale $Lang "plural_text_dynamic" 3 }}

<!--Dynamic text loading-->
{{ TextDynamic $Locale $Lang "dynamic_hello" "Jānis Bērziņš" }}
{{ TextDynamic $Locale $Lang "dynamic_bill_sum" 7.23 }}
```