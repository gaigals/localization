package main

import (
	"fmt"
	"net/http"
)

func homeHandler(w http.ResponseWriter, r *http.Request) {
	err := templateHome.AddDataMap(TemplateData{"Lang": "en", "Locale": locale})
	if err != nil {
		fmt.Println(err)
		return
	}

	err = templateHome.Render(w)
	if err != nil {
		fmt.Println(err)
		return
	}
}
