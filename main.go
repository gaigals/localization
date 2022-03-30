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

func initializeLocale() error {
	locale.AddLanguages("lv", "en")
	locale.SetValueNoErr("lv", "hello_world", "Sveicināta, Pasaule!")
	locale.SetValueNoErr("en", "hello_world", "Hello, World!")

	return nil
}

func main() {
	initializeHandlers()

	err := initializeTemplates()
	if err != nil {
		log.Fatalln(err)
	}

	err = initializeLocale()
	if err != nil {
		log.Fatalln(err)
	}

	log.Fatal(http.ListenAndServe(":3000", nil))
}
