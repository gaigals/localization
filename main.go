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

func readTranslates(path string) ([]Translate, error) {
	return LoadTranslatesXML(path)
}

func initializeLocale(translates []Translate) error {
	locale.AddLanguages("lv", "en")
	//locale.SetValueNoErr("lv", "hello_world", "Sveicināta, Pasaule!")
	//locale.SetValueNoErr("en", "hello_world", "Hello, World!")

	for _, v := range translates {
		err := locale.SetValue(v.Language, v.Key, v.Value)
		if err != nil {
			return err
		}
	}

	return nil
}

func main() {
	initializeHandlers()

	err := initializeTemplates()
	if err != nil {
		log.Fatalln(err)
	}

	translates, err := readTranslates("locals/home.xml")
	if err != nil {
		log.Fatalln(err)
	}

	err = initializeLocale(translates)
	if err != nil {
		log.Fatalln(err)
	}

	log.Fatal(http.ListenAndServe(":3000", nil))
}
