package main

import (
	"fmt"
	"html/template"
	"net/http"
)

type TemplateData map[string]interface{}

type Template struct {
	Template *template.Template
	Data     TemplateData
}

func NewTemplate(funcMap template.FuncMap, htmlPath ...string) (*Template, error) {
	templ := Template{Template: template.New("templ")}
	_ = templ.AddFuncs(funcMap)

	err := templ.AddHTML(htmlPath...)
	if err != nil {
		return nil, err
	}

	templ.Data = make(TemplateData, 0)

	return &templ, nil
}

func (t *Template) AddHTML(htmlPath ...string) error {
	if t.Template == nil {
		return fmt.Errorf("template field is nil")
	}

	if len(htmlPath) == 0 {
		return nil
	}

	var err error

	t.Template, err = t.Template.ParseFiles(htmlPath...)
	if err != nil {
		return err
	}

	return nil
}

func (t *Template) AddFuncs(funcMap template.FuncMap) error {
	if t.Template == nil {
		return fmt.Errorf("template field is nil")
	}

	if len(funcMap) == 0 {
		return nil
	}

	t.Template = t.Template.Funcs(funcMap)
	return nil
}

func (t *Template) AddDataMap(data TemplateData) error {
	if len(data) == 0 {
		return nil
	}

	for k, v := range data {
		err := t.AddData(k, v)
		if err != nil {
			return err
		}
	}

	return nil
}

func (t *Template) AddData(key string, value interface{}) error {
	_, exist := t.Data[key]
	if exist {
		return fmt.Errorf("template Data already contains key '%s'", key)
	}

	t.Data[key] = value

	return nil
}

func (t *Template) Render(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "text/html;charset=utf-8")

	err := t.Template.ExecuteTemplate(w, "layout", t.Data)
	if err != nil {
		return err
	}

	return nil
}
