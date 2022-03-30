package main

type TextMap map[string]string

type Language struct {
	Keyword string
	Map     TextMap
}

func (l *Language) Value(key string) string {
	return l.Map[key]
}

func (l *Language) SetValue(key, value string) {
	l.Map[key] = value
}
