package main

import (
	"fmt"
	"net/http"
)

const langEN = "en"
const langLV = "lv"

func parseLocaleForm(w http.ResponseWriter, r *http.Request) error {
	err := r.ParseForm()
	if err != nil {
		return fmt.Errorf("failed to parse locale form, error: %w", err)
	}

	lang := r.Form.Get("lang")

	if lang != langEN && lang != langLV {
		return fmt.Errorf("unsupported language '%s'", lang)
	}

	CreateLocaleCookie(w, lang)

	return nil
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		err := parseLocaleForm(w, r)
		if err != nil {
			fmt.Println(err)
		}

		http.Redirect(w, r, "/", http.StatusFound)
	}

	// GET

	lang, err := GetCookieLanguage(r, langEN)
	if err != nil {
		fmt.Println(err)
		return
	}

	err = templateHome.AddDataMap(TemplateData{"Lang": lang, "Locale": locale})
	if err != nil {
		fmt.Println("here", err)
		return
	}

	err = templateHome.Render(w)
	if err != nil {
		fmt.Println(err)
		return
	}
}
