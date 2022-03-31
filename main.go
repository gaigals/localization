package main

import (
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
		"GetMapValue": GetMapValue,
		"Text":        Text,
	}

	var err error

	templateHome, err = NewTemplate(funcMap, "templ/home.html")
	if err != nil {
		return err
	}

	return nil
}

func readTranslates(path, defaultLanguage string) (*YAMLFile, error) {
	return LoadTranslates(path, defaultLanguage)
}

func initializeLocale(yaml *YAMLFile) error {
	locale.AddLanguages("lv", "en")

	// Appending new translates.
	//locale.SetValueNoErr("lv", "hello_world", "Sveicināta, Pasaule!")
	//locale.SetValueNoErr("en", "hello_world", "Hello, World!")

	// 1 file manual handling.
	//err := locale.AddTranslate(xml.Translates...)
	//if err != nil {
	//	return err
	//}

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

	err = initializeLocale(translates)
	if err != nil {
		log.Fatalln(err)
	}

	log.Fatal(http.ListenAndServe(":3000", nil))
}
