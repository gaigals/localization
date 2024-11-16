# Localization

Localization is API written in GoLang and provides easy-to-use tools for handling localization in
target API by loading YAML translate files or by manually adding translates.
There is no limit of how many languages are used, and if there is need, it is also possible to use only 1 target
language for displaying text by using keywords.

## YAML files

YAML file probably is the best way of handling localization texts. YAML file parser supports basic translations,
non-plural and plural translations.

__Examples:__
One-lines (default language must be provided for loader).
```yaml
# If provided default language in loader is "en" then parser will parse "key0" as "en" translation.
# This can be useful for writing short one-lines if only 1 translation is required or for
# splitting translations in different files, for example, each language in different file.
key0: "my super text"
```
or can be written as:
```yaml
# Nothing changes but language is determined in YAML not in loader.
key0:
  - en: "my super text"
```

Basic translations:
```yaml
# Basic translation example, "en" and "lv" translations for keyword "key0".
# Order does not matter.
key0:
  - en: "my super text"
  - lv: "mans super teksts"
```
or (less readable but not restricted):
```yaml
key0:
  - en: "my super text"

key0:
  - lv: "mans super teksts"
```

Dynamic text:
```yaml
# Translation by using GoLang supported formatting symbols.
# There is no restrictions for custom formatting but it will require custom text
# formatter.
key0:
  - en: "Hello, %s"
  - lv: "Sveiks, %s"
```

Writing plural translations:
```yaml
# Way of providing plural text. It is NOT required to write plural version for both
# languages.
key0:
  - en:
      - "not plural"
      - "plural"
  - lv:
      - "nav daudzskaitlis"
      - "daudzskaitlis"
```

Plural dynamic translations:
```yaml
key0:
  - en:
      - "%d item"
      - "%d items"
  - lv:
      - "%d lieta"
      - "%d lietas"
```

## Locale structure

ALl translations are contained in `Locale` structure. Every language keyword must be
unique or multiple `Locale`s must be used. Structure looks like this:

```go
type Locale struct {
	Languages   []Language // List of initialized languages.
	StrictUsage bool       // Is it allowed to use other language keys as backup.
}
```

`Locale.StrictUsage` - when FALSE, `Locale` `Value()` or `ValuePlural()` func calls creates language priority
list which is used whenever target language does not contain key. In this case, `Locale` will
loop over generated list until match is found or else error gets returned.
If this order matter to you, then initialize languages in desired order.

__For example:__
Initialized order `[]string{"lv", "en", "lt"}` on  `Locale.Value("lv", "key")` will result in `lv -> en -> lt`
but calling `Locale.Value("en", "key")` will result in `en -> lv -> lt`

To disable this functionality, set `Locale.StrictUsage` as true.


Can be initialized with:
```go
// NewLocale can be used to initialize new Locale structure with provided languages.
func NewLocale(strictUsage bool, lang ...string) (*Locale, error)
```


`Locale` public methods:
```go
// GlobalYAMLLoad can be used to load yaml files directly into locale.
// Examples: "locales/*", "locales/*.yml", "locales/**/*.yml"
func (l *Locale) GlobalYAMLLoad(defaultLang, pattern string) error {

// AddLanguages can be used to add new languages to Locale.
func (l *Locale) AddLanguages(lang ...string) error

// LoadYAMLFile can be used to load and parse multiple YAML files with
// containing translations and directly load them into current Locale.
func (l *Locale) LoadYAMLFile(defaultLanguage string, filePath ...string) error

// AddYAMLFile can be used to add 1 or more YAMLFile's translations to current Locale.
func (l *Locale) AddYAMLFile(files ...*YAMLFile) error

// AddTranslate can be used to add 1 or more translations to current Locale.
func (l *Locale) AddTranslate(translates ...Translate) error

// SetValue can be used to set translation plural and non-plural values for target
// language.
func (l *Locale) SetValue(langKey, textKey, value, plural string) error

// SetValueNoErr can be used to set translation plural and non-plural values for target
// language without returning any errors.
func (l *Locale) SetValueNoErr(langKey, textKey, value, plural string)

// GetLanguage can be used to get language with specific keyword ("en", "lv" etc).
func (l *Locale) GetLanguage(langKey string) (*Language, error)
```

## Language structure

`Language`  contain all targeted language plural and non-plural translation, and required
methods for populating and for accessing values. Each `Language` MUST have unique keyword (for example, "en", "lv" etc.)
All initialized languages are hold by `Locale`.

Note: Use `Language` methods only if it's required as they do not provide extra checks, for that functionality use
`Locale` methods as they do more error checking.


```go
type TextMap map[string][2]string

type Language struct {
	Keyword string  // Language keyword ("en", "lv" etc.)
	Map     TextMap // Translations (1. Non-Plural 2. Plural)
}

```

`Language` public methods:
```go
// ValueNoErr can be used to extract non-plural translation from language
// by providing translation keyword/key but does not return error if key
// does not exist (returns empty string).
func (l *Language) ValueNoErr(key string) string

// Value can be used to extract non-plural translation from language
// by providing translation keyword/key.
func (l *Language) Value(key string) (string, error)


// ValuePluralNoErr can be used to extract plural translation from language
// by providing translation keyword/key  but does not return error if key
// does not exist (returns empty string).
func (l *Language) ValuePluralNoErr(key string) string

// ValuePlural can be used to extract plural translation from language
// by providing translation keyword/key.
func (l *Language) ValuePlural(key string) (string, error)

// SetValue can be used to set non-plural and plural translation for language
// by providing translation keyword/key, non-plural value (value) and plural value (plural).
func (l *Language) SetValue(key, value, plural string)
```


## Translate structure

`Translate` structure is used to hold translation information and is mainly used by YAML file loader
to represent each translation entry.
`Translate` also can be used to manually add translates to `Locale` by using `Locale.AddTranslate()`.

```go
type Translate struct {
	Key      string // Key/ID for translation.
	Language string // Language keyword ("lv", "en" etc.).
	Value    string // Translation ("some text").
	Plural   string // Translation in plural.
}
```


## YAMLFile structure

`YAMLFile` is used to load and contain YAML file translates. YAML file or multiple files can be directly loaded in `Locale`
with `Locale.LoadYAMLFile()` or `localization.LoadYAMLFiles()` methods. Use this struct if additional setup,
or more control is needed.


```go
type YAMLFile struct {
    FileName   string               // Target YAML file path
    Translates []Translate          // YAML file loaded translates.
}
````


## Other public functions:

```go
// LoadYAMLFiles can be used to load and parse one or more YAML files with
// containing translations.
func LoadYAMLFiles(defaultLanguage string, path ...string) ([]*YAMLFile, error)

// Text returns translation by using provided langKey and textKey.
func Text(locale Locale, langKey, textKey string) (string, error)

// TextPlural returns plural or non-plural translation by provided langKey and textKey (hard-coded plural text).
func TextPlural(locale Locale, langKey, textKey string, isPlural bool) (string, error)

// Textf can be used to extract value from provided language and Locale,
// and additionally applies string formatting before returning result.
func Textf(locale Locale, langKey, textKey string, input ...interface{}) (string, error)

// TextPluralf can be used to extract value from provided language and Locale,
// and will return non-plural or plural value determined by passed bool. Final
// result will be formatted by using passed input.
func TextPluralf(locale Locale, langKey, textKey string, bool, input ...interface{}) (string, error)

// TextPluralIntf can be used to extract value from provided language and Locale,
// and will return non-plural or plural value determined by first int param.
// If first int param is more than "1" then plural value will be returned, additionally applies
// string formatting before returning result.
func TextPluralIntf(locale Locale, langKey, textKey string, input ...interface{}) (string, error)
```


## Setup

To add translations, new `Locale` structure is required with initialized languages. After
it's done, you can manually add translations by direct func calls or by loading YAML files.
If target application requires key reuse, for example, identical pages with same keys but different values,
then create new `Locale` and YAML file for each page and load those file into target `Locale`.


### Locale setup:

Single `Locale` creation:
```go
// Create new Locale with StrictUsage off and add languages "en", "lv".
locale, err := localization.NewLocale(false, "en", "lv")
if err != nil {
	log.Fatalf(err)
}
```
or do it manually:
```go
// Manually create Locale with StrictUsage off.
locale := localization.Locale{StrictUsage: false}

// Add languages
err := localization.AddLanguages("en", "lv")
if err != nil {
    log.Fatalf(err)
}
```

### Loading translations from YAML file:

To simply load YAML files with given pattern directly in created `Locale` with:
```go
err = createdLocale.GlobalYAMLLoad("lv", "locales/*.yml")
if err != nil {
    log.Fatalf(err)
}
``


But if you want manually load every single yaml file, you can do it in this way:
```go
// Load 2 YAML translates and set "en" as default language.
// Default language is required of for loader set target language
// when YAML file uses one-lines (for example: key0: "some text").
err := locale.LoadYAMLFiles("en", "translate_0.yaml", "translate_1.yaml")
if err != nil {
    log.Fatalf(err)
}
```

__Other ways to do it (longer way):__
```go
// These steps can be useful if you want to do some validation or additional parsing
// before adding translations to the target Locale.

// Load multiple YAML files.
yamlFiles, err := localization.LoadYAMLFiles("en", "translate_0.yaml", "translate_1.yaml")
if err != nil {
    log.Fatalf(err)
}
```

Adding []*YAMLFile to Locale:
```go
err := locale.AddYAMLFile(yamlFiles...)
if err != nil {
    log.Fatalf(err)
}
```


### Adding translations manually (without YAML):

Straight ahead adding values:
```go
// SetValue will do additional checks and will return error if something is not correct.
err := locale.SetValue("en", "key_hello", "non plural text", "plural text but not required")
if err != nil {
    log.Fatalf(err)
}

// If target language is initialized and you are confident about input then you can use
// this method (does not return errors):
locale.SetValueNoErr("en", "key_hello", "non plural text", "plural text but not required")
```

Using `localization.Translate`:
```go
// Create as many Translates you want.
translates := []localization.Translate{
	{ Key: "key_hello", Language: "en", Value: "Hello, World!", Plural: "" },
        { Key: "key_hello", Language: "lv", Value: "Sveika, Pasaule!", Plural: "" },
        { Key: "other_key", Language: "en", Value: "%d item", Plural: "%d items" },
}

err := locale.AddTranslate(translates...)
if err != nil {
    log.Fatalf(err)
}
```



### Loading translations by ready to use Template functions

localization API provides ready-to-use functions for loading translations from templates
(can be used for other purposes). All you have to do - add desired functions to the
template func map.

__Examples:__

Simply loading text value (only non-plural):
```go
// Load "key_hello" english translation.
text, err := localization.Text(locale, "en", "key_hello")
if err != nil {
    log.Fatalf(err)
}
```

Loading plural or non-plural text:
```go
// Load "key_hello" english plural or non-plural translation.
// true - return plural.
// false - return non-plural.
text, err := localization.TextPlural(locale, "en", "key_hello", true)
if err != nil {
    log.Fatalf(err)
}
```

Loading dynamic text (only non-plural):
```go
// Load "key_hello" dynamic english translation.
// Example:
// key_hello: "Hello, %s!"
// Result: "Hello, John!"
text, err := localization.Textf(locale, "en", "key_hello", "John")
if err != nil {
    log.Fatalf(err)
}
```

Loading plural or non-plural dynamic text:
```go
// Load "key0" english plural or non-plural dynamic translation.
// true - return plural.
// false - return non-plural.
// Example:
// key0:
//  - en
//      - "%s %d %s :("
//      - "%s %d %s :)"
// Result: "John has 5 items :)"
text, err := localization.TextPluralf(locale, "en", "key0", true,  "John has", 5, "items")
if err != nil {
    log.Fatalf(err)
}
```

Loading plural or non-plural dynamic controlled by 1st `int` value:
```go
// Load "key0" english plural or non-plural dynamic translation controlled
// by first int parameter.
// If 1st int parameter is greater than 1 then returned value will be plural.
// Example:
// key0:
//  - en
//      - "John has %d item %s"
//      - "John has %d items %s"
// Result: "John has 1 item :)"
text, err := localization.TextPluralIntf(locale, "en", "key0",  1, ":)")
if err != nil {
    log.Fatalf(err)
}
```

Example of HTML file:
```html
{{ $Page := . }}

<!--Load page locale-->
{{ $Locale := $Page.Locale }}
<!--Get language keyword ("en" or "lv")-->
{{ $Lang := $Page.Lang }}

<!-- ... -->

<!--Basic translation loading-->
{{ Text $Locale $Lang "hello_world" }}

<!--Loading hard-coded plural translation-->
{{ TextPlural $Locale $Lang "plural_text_hard_coded" false }}

<!--Loading plural dynamic text, controlled by bool -->
{{ TextPluralf $Locale $Lang true "plural_text_dynamic" 1 }}
{{ TextPluralf $Locale $Lang true "plural_text_dynamic" 3 }}

<!--Loading plural dynamic text, param must be int -->
{{ TextPluralIntf $Locale $Lang "plural_text_dynamic_int" 1 }}
{{ TextPluralIntf $Locale $Lang "plural_text_dynamic_int" 3 }}

<!--Dynamic text loading-->
{{ Textf $Locale $Lang "dynamic_hello" "Jānis Bērziņš" }}
{{ Textf $Locale $Lang "dynamic_bill_sum" 7.23 }}
```



If you have single, global `Locale` then you can wrap provided functions and your own version.

For example:

```go
func MyTextFunc(langKey, textKey string) (string, error) {
	return localization.Text(myGlobalLocale, langKey, textKey)
}
```


Or you can build your own processing functions by using API provided public `Locale` methods.



### Manually reading Locale values

__Examples:__

Non-plural value extraction:
```go
// Get non-plural "lv" translation with key "hello".
myNonPluralValue, err := locale.Value("lv", "hello")
if err != nil {
	log.Fatalf(err)
}
```

Plural value extraction:
```go
// Get plural "lv" translation with key "hello".
myPluralValue, err := locale.ValuePlural("lv", "hello")
if err != nil {
	log.Fatalf(err)
}
```
