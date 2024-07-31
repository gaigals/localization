package localization

import (
	"errors"
	"fmt"
	"reflect"
)

// ErrDynamicNotInt gets triggered when dynamic-plural based function first
// dynamic parameter is not int type.
var ErrDynamicNotInt error = errors.New("plural dynamic input must be int")

// Text can be used to extract text value from provided language and Locale.
// Returns textKey value or error if something went wrong.
//
// Params:
// locale - target Locale.
// langKey - target language keyword ("en", "lv" etc).
// textKey - text keyword/id.
func Text(locale Locale, langKey, textKey string) (string, error) {
	return locale.Value(langKey, textKey)
}

// TextPlural can be used to extract value from provided language and Locale,
// and will return non-plural or plural value determined by isPlural param.
//
// Params:
// locale - target Locale.
// langKey - target language keyword ("en", "lv" etc).
// textKey - text keyword/id.
// isPlural - true if desired result must be in plural or false if desired result
// must be non-plural value.
func TextPlural(locale Locale, langKey, textKey string, isPlural bool) (string, error) {
	// Extract plural or non-plural value determined by isPlural.
	text, err := extractPluralText(locale, langKey, textKey, isPlural)
	if err != nil {
		return "", err
	}

	return text, nil
}

// Textf can be used to extract value from provided language and Locale,
// and additionally applies string formatting before returning result.
// For example:
// Targeted key value ----> key_hello: "Hello, %s"
// With func call ----> Textf(myLocale, "en", "key_hello", "John")
// Will return ----> "Hello, John"
//
// Params:
// locale - target Locale.
// langKey - target language keyword ("en", "lv" etc).
// textKey - text keyword/id.
// input - dynamic input, requires as many params as required by targeted string format.
func Textf(locale Locale, langKey, textKey string, input ...interface{}) (string, error) {
	// Extract non-plural value.
	text, err := locale.Value(langKey, textKey)
	if err != nil {
		return "", err
	}

	// Format string and return.
	return fmt.Sprintf(text, input...), nil
}

// TextPluralf can be used to extract value from provided language and Locale,
// and will return non-plural or plural value determined by passed bool. Final
// result will be formatted by using passed input.
// If first int param is more than "1" then plural value will be returned.
// For example:
// Targeted key value:
//	key_item:
//		- en:
//			- "%d item %s"
// 			- "%d items %s"
// With func call ----> TextPluralf(myLocale, "en", "key_item", 1, ":)")
// Will return ----> "1 item :)"
// OR
// With func call ----> TextPluralf(myLocale, "en", "key_item", 2, ":)")
// Will return ----> "2 items :)"
//
// Params:
// locale - target Locale.
// langKey - target language keyword ("en", "lv" etc).
// textKey - text keyword/id.
// isPlural - true if desired result must be in plural or false if desired result
// must be non-plural value.
// input - dynamic input, requires as many params as required by targeted string format.
func TextPluralf(locale Locale, langKey, textKey string, isPlural bool, input ...interface{}) (string, error) {
	// Extract plural or non-plural value determined by isPlural.
	text, err := extractPluralText(locale, langKey, textKey, isPlural)
	if err != nil {
		return "", err
	}

	// Format string and return.
	return fmt.Sprintf(text, input...), nil
}

// TextPluralIntf can be used to extract value from provided language and Locale,
// and will return non-plural or plural value determined by first int param.
// If first int param is more than "1" then plural value will be returned, additionally applies
// string formatting before returning result.
// For example:
// Targeted key value:
//	key_item:
//		- en:
//			- "%d item %s"
// 			- "%d items %s"
// With func call ----> TextPluralIntf(myLocale, "en", "key_item", 1, ":)")
// Will return ----> "1 item :)"
// OR
// With func call ----> TextPluralIntf(myLocale, "en", "key_item", 2, ":)")
// Will return ----> "2 items :)"
//
// Params:
// locale - target Locale.
// langKey - target language keyword ("en", "lv" etc).
// textKey - text keyword/id.
// input - dynamic input, requires as many params as required by targeted string format.
// First input param must be int.
func TextPluralIntf(locale Locale, langKey, textKey string, input ...interface{}) (string, error) {
	// Check if first input element is more than 1.
	isPlural, err := isFirstElementPlural(input...)
	if err != nil {
		return "", err
	}

	// Extract plural or non-plural value determined by isPlural.
	text, err := extractPluralText(locale, langKey, textKey, isPlural)
	if err != nil {
		return "", err
	}

	// Format string and return.
	return fmt.Sprintf(text, input...), nil
}

// extractPluralText extracts plural or non-plural value from target locale determined
// by isPlural.
// Returns error if something went wrong.
func extractPluralText(locale Locale, langKey, textKey string, isPlural bool) (string, error) {
	if isPlural {
		return locale.ValuePlural(langKey, textKey)
	}

	return locale.Value(langKey, textKey)
}

// isFirstElementPlural checks if passed input first element is int type and
// if int value is greater than 1.
// Returns true if int value is greater than 1 or error first element is not int.
func isFirstElementPlural(input ...interface{}) (bool, error) {
	if len(input) == 0 {
		return false, ErrDynamicNotInt
	}

	// Check if first element is int.
	val := reflect.ValueOf(input[0])
	if val.Kind() != reflect.Int {
		return false, ErrDynamicNotInt
	}

	// Check if int is greater than 1 (plural).
	if val.Int() > 1 {
		return true, nil
	}

	// Target value is non-plural.
	return false, nil
}
