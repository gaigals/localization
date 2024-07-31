package localization

import (
	"fmt"
	"os"
	"reflect"
	"strings"
	"testing"
)

func createTempFile(dirPath, fileName, content string) error {
	file, err := os.Create(fmt.Sprintf("%s/%s", dirPath, fileName))
	if err != nil {
		return err
	}

	_, err = file.WriteString(content)

	return err
}

func TestNewYAMLContent(t *testing.T) {
	testCases := []struct {
		lang string
	}{
		{"en"},
		{"lv"},
		{""},
	}

	for k, v := range testCases {
		content := newYAMLContent(v.lang)

		if !strings.EqualFold(v.lang, content.defaultLanguage) {
			t.Fatalf("unexpected language, index=%d, expected=%s, actual=%s",
				k, v.lang, content.defaultLanguage)
		}
	}
}

func TestYAMLContent_loadBytes(t *testing.T) {
	tempDir := t.TempDir()

	testCases := []struct {
		fileName        string
		fileContent     string
		createFile      bool
		failureExpected bool
	}{
		// No errors - file exist and content matches.
		{"loadBytes_0.yaml", "some random text\n", true, false},
		// Error - file does not exist.
		{"rand.yaml", "", false, true},
	}

	for k, v := range testCases {
		filePath := fmt.Sprintf("%s/%s", tempDir, v.fileName)

		if v.createFile {
			err := createTempFile(tempDir, v.fileName, v.fileContent)
			if err != nil {
				t.Fatalf("unexpected error, index=%d, error: %s", k, err)
			}
		}

		content := yamlContent{}

		bytes, err := content.loadBytes(filePath)
		if err != nil && !v.failureExpected {
			t.Fatalf("unexpected error, index=%d, error: %s", k, err)
		}
		if err == nil && v.failureExpected {
			t.Fatalf("expected error, index=%d", k)
		}

		if !strings.EqualFold(string(bytes), v.fileContent) {
			t.Fatalf("unexpected file content, index=%d, expected=%s, actual=%s",
				k, v.fileContent, string(bytes))
		}
	}
}

func TestYAMLContent_unmarshal(t *testing.T) {
	tempDir := t.TempDir()

	testCases := []struct {
		fileName            string
		fileContent         string
		expectedMapAsString string
		createFile          bool
		failureExpected     bool
	}{
		{ // No errors - file exist and content matches.
			"unmarshal_0.yaml",
			"key0: \"text\"\n",
			"map[key0:text]",
			true, false,
		},
		{ // Error - unknown YAML format (yaml.Unmarshal error).
			"unmarshal_1.yaml",
			"some random text\n",
			"map[]",
			true, true,
		},
		{ // Error - no file (passed bytes are nil).
			"",
			"",
			"map[]",
			false, true,
		},
	}

	for k, v := range testCases {
		filePath := fmt.Sprintf("%s/%s", tempDir, v.fileName)

		var err error
		var bytes []byte = nil
		content := newYAMLContent("en")

		if v.createFile {
			err = createTempFile(tempDir, v.fileName, v.fileContent)
			if err != nil {
				t.Fatalf("unexpected error, index=%d, error: %s", k, err)
			}
			bytes, err = content.loadBytes(filePath)
			if err != nil {
				t.Fatalf("unexpected error, index=%d, error: %s", k, err)
			}
		}

		err = content.unmarshal(bytes)
		if err != nil && !v.failureExpected {
			t.Fatalf("unexpected error, index=%d, error: %s", k, err)
		}
		if err == nil && v.failureExpected {
			t.Fatalf("expected error, index=%d", k)
		}

		mapAsString := fmt.Sprintf("%v", content.Data)

		if !strings.EqualFold(mapAsString, v.expectedMapAsString) {
			t.Fatalf("unexpected Data content, index=%d, expected=%s, actual=%s",
				k, v.expectedMapAsString, mapAsString)
		}
	}
}

func TestYAMLContent_parse(t *testing.T) {
	tempDir := t.TempDir()

	testCases := []struct {
		fileName        string
		fileContent     string
		defaultLang     string
		unmarshalData   bool
		failureExpected bool
		expected        []Translate
	}{
		{ // No errors - file exist and content matches.
			"unmarshal_0.yaml",
			"key0: \"text\"\n",
			"en",
			true,
			false,
			[]Translate{
				{Key: "key0", Language: "en", Value: "text", Plural: ""},
			},
		},
		{ // No errors - file exist and content matches.
			"unmarshal_1.yaml",
			"key0: \"text\"\nkey1: \"text_1\"\n",
			"lv",
			true,
			false,
			[]Translate{
				{Key: "key0", Language: "lv", Value: "text", Plural: ""},
				{Key: "key1", Language: "lv", Value: "text_1", Plural: ""},
			},
		},
		{ // No errors - file exist and content matches.
			"unmarshal_2.yaml",
			"key0:\n  - en: \"some_text\"\n",
			"en",
			true,
			false,
			[]Translate{
				{Key: "key0", Language: "en", Value: "some_text", Plural: ""},
			},
		},
		{ // No errors - file exist and content matches.
			"unmarshal_3.yaml",
			"key0:\n  - lv: \"some_text\"\n",
			"en",
			true,
			false,
			[]Translate{
				{Key: "key0", Language: "lv", Value: "some_text", Plural: ""},
			},
		},
		{ // No errors - file exist and content matches.
			"unmarshal_4.yaml",
			"key0:\n  - lv: \"some_text\"\n  - en: \"some_other_text\"\n",
			"en",
			true,
			false,
			[]Translate{
				{Key: "key0", Language: "lv", Value: "some_text", Plural: ""},
				{Key: "key0", Language: "en", Value: "some_other_text", Plural: ""},
			},
		},
		{ // No errors - file exist and content matches.
			"unmarshal_5.yaml",
			"key0:\n  - lv:\n    - \"non_plural\"\n    - \"plural\"\n",
			"en",
			true,
			false,
			[]Translate{
				{Key: "key0", Language: "lv", Value: "non_plural", Plural: "plural"},
			},
		},
		{ // No errors - file exist and content matches.
			"unmarshal_6.yaml",
			"key0:\n  - lv:\n    - \"non_plural\"\n    - \"plural\"\n" +
				"  - en:\n    - \"en_non_plural\"\n    - \"en_plural\"\n",
			"en",
			true,
			false,
			[]Translate{
				{Key: "key0", Language: "lv", Value: "non_plural", Plural: "plural"},
				{Key: "key0", Language: "en", Value: "en_non_plural", Plural: "en_plural"},
			},
		},
		{ // No error and no data - yamlContent.Data is nil
			"unmarshal_7.yaml",
			"key: \"some_text\"",
			"en",
			false,
			false,
			nil,
		},
		{ // No errors - file exist and content matches.
			"unmarshal_8.yaml",
			"key0:\n  - lv:\n    - \"non_plural\"\n",
			"en",
			true,
			false,
			[]Translate{
				{Key: "key0", Language: "lv", Value: "non_plural", Plural: ""},
			},
		},

		{ // Error - unsupported type - int
			"unmarshal_9.yaml",
			"key: 1",
			"en",
			true,
			true,
			nil,
		},
		{ // Error - unsupported map key type - int
			"unmarshal_10.yaml",
			"key0:\n  - 1: \"some_text\"\n",
			"en",
			true,
			true,
			nil,
		},
		{ // Error - map value must be string or slice/list
			"unmarshal_11.yaml",
			"key0:\n  - en: 1.00",
			"en",
			true,
			true,
			nil,
		},
		{ // Error - more than 2 entries in plural list
			"unmarshal_12.yaml",
			"key0:\n  - lv:\n    - \"\"\n    - \"\"\n    - \"\"\n",
			"en",
			true,
			true,
			nil,
		},
	}

	for k, v := range testCases {
		filePath := fmt.Sprintf("%s/%s", tempDir, v.fileName)

		var err error
		var bytes []byte = nil
		content := newYAMLContent(v.defaultLang)

		err = createTempFile(tempDir, v.fileName, v.fileContent)
		if err != nil {
			t.Fatalf("unexpected error, index=%d, error: %s", k, err)
		}
		bytes, err = content.loadBytes(filePath)
		if err != nil {
			t.Fatalf("unexpected error, index=%d, error: %s", k, err)
		}

		if v.unmarshalData {
			err = content.unmarshal(bytes)
			if err != nil {
				t.Fatalf("unexpected error, index=%d, error: %s", k, err)
			}
		}

		translates, err := content.parse()
		if err != nil && !v.failureExpected {
			t.Fatalf("unexpected error, index=%d, error: %s", k, err)
		}
		if err == nil && v.failureExpected {
			t.Fatalf("expected error, index=%d", k)
		}

		if !reflect.DeepEqual(translates, v.expected) {
			t.Fatalf("unexpected results, index=%d, input:\n%s\nexpected:\n%+v\nactual:\n%+v\n\n",
				k, v.fileContent, v.expected, translates)
		}
	}
}
