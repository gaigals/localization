package main

type TextMap map[string][2]string

type Language struct {
	Keyword string
	Map     TextMap
}

func (l *Language) Value(key string) string {
	return l.Map[key][0]
}

func (l *Language) ValuePlural(key string) string {
	return l.Map[key][1]
}

func (l *Language) SetValue(key, value, plural string) {
	l.Map[key] = [2]string{value, plural}
}
