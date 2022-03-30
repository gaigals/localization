package main

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"os"
)

type XMLFile struct {
	Translates []Translate `xml:"value"`
}

type Translate struct {
	Key      string `xml:"key,attr"`
	Language string `xml:"lang,attr"`
	Value    string `xml:",chardata"`
}

func LoadTranslatesXML(path string) ([]Translate, error) {
	xmlFile, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open xmlFile with path '%s', error: %w",
			path, err)
	}

	bytes, err := ioutil.ReadAll(xmlFile)
	if err != nil {
		_ = xmlFile.Close()
		return nil, fmt.Errorf("failed to cast xml file content as bytes, error: %w",
			err)
	}

	_ = xmlFile.Close()

	var translates XMLFile

	err = xml.Unmarshal(bytes, &translates)
	if err != nil {
		return nil,
			fmt.Errorf("failed to unmarshal xml file with path '%s', error: %w", path, err)
	}

	return translates.Translates, nil
}
