package main

import (
	"errors"
	"fmt"
	"net/http"
)

const cookieLocaleName = "locale"

func CreateLocaleCookie(w http.ResponseWriter, lang string) {
	cookie := http.Cookie{
		Name:     cookieLocaleName,
		Value:    lang,
		Path:     "/",
		Domain:   "localhost",
		MaxAge:   0,
		Secure:   false,
		HttpOnly: true,
	}

	http.SetCookie(w, &cookie)
}

func ReadLocaleCookie(r *http.Request) (string, error) {
	cookie, err := r.Cookie(cookieLocaleName)
	if err != nil {
		return "", fmt.Errorf("failed to acquire cookie '%s', error: %w",
			cookieLocaleName, err)
	}

	return cookie.Value, err
}

func GetCookieLanguage(r *http.Request, defaultLang string) (string, error) {
	lang, err := ReadLocaleCookie(r)

	// If cookie is not in present then return default language
	if errors.Is(err, http.ErrNoCookie) {
		return defaultLang, nil
	}

	// Handle unexpected error.
	if err != nil {
		return "", err
	}

	return lang, nil
}
