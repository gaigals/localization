package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
)

var locale Locale

var templateHome *Template

func initializeHandlers() {
	http.HandleFunc("/", homeHandler)
}

func initializeTemplates() error {
	funcMap := template.FuncMap{
		"GetMapValue":       GetMapValue,
		"Text":              Text,
		"TextPlural":        TextPlural,
		"TextDynamic":       TextDynamic,
		"TextPluralDynamic": TextPluralDynamic,
	}

	var err error

	templateHome, err = NewTemplate(funcMap, "templ/home.html")
	if err != nil {
		return err
	}

	return nil
}

func readTranslates(path, defaultLanguage string) (*YAMLFile, error) {
	return LoadYAMLTranslateFile(path, defaultLanguage)
}

func initializeLocale(yaml *YAMLFile) error {
	// Enable strict usage for translation (no key for lang = error)
	//locale.StrictUsage = true

	locale.AddLanguages("lv", "en")

	// Appending new translates (last param - plural).
	//locale.SetValueNoErr("lv", "hello_world", "Sveicināta, Pasaule!", "")
	//locale.SetValueNoErr("en", "hello_world", "Hello, World!", "")

	// Multi file handling.
	err := locale.AddYAMLFile(yaml)
	if err != nil {
		return err
	}

	return nil
}

func main() {
	initializeHandlers()

	err := initializeTemplates()
	if err != nil {
		log.Fatalln(err)
	}

	translates, err := readTranslates("locals/home.yaml", langLV)
	if err != nil {
		log.Fatalln(err)
	}

	for _, v := range translates.Translates {
		fmt.Println(v)
	}

	err = initializeLocale(translates)
	if err != nil {
		log.Fatalln(err)
	}

	log.Fatal(http.ListenAndServe(":3000", nil))
}
