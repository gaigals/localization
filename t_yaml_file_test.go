package localization

import (
	"fmt"
	"reflect"
	"strings"
	"testing"
)

func applyTemDirLocation(isFailureExpected bool, tempFilePath []string, yamlFiles ...*YAMLFile) {
	if isFailureExpected {
		return
	}

	if len(yamlFiles) == 0 {
		return
	}

	for k, v := range yamlFiles {
		if v == nil {
			continue
		}

		v.FilePath = tempFilePath[k]
	}
}

func TestLoadYAML(t *testing.T) {
	tempDir := t.TempDir()

	testCases := []struct {
		fileName        string
		fileContent     string
		defaultLang     string
		createFile      bool
		failureExpected bool
		expected        *YAMLFile
	}{
		{ // No errors - file exist and content matches.
			"file_0.yaml",
			"key0: \"text\"\n",
			"en",
			true,
			false,
			&YAMLFile{
				FilePath: "file_0.yaml", Translates: []Translate{
					{Key: "key0", Language: "en", Value: "text", Plural: ""},
				},
			},
		},
		{ // No Error - empty file exist but no errors.
			"file_1.yaml",
			"",
			"en",
			true,
			false,
			&YAMLFile{
				FilePath: "file_0.yaml", Translates: nil,
			},
		},

		{ // Error - file does not exist.
			"non_existing_file.yaml",
			"key0: \"text\"\n",
			"en",
			false,
			true,
			nil,
		},
		{ // Error - file does not exist.
			"",
			"key0: \"text\"\n",
			"en",
			false,
			true,
			nil,
		},
		{ // Error - unsupported value type - int
			"file_0.yaml",
			"key0: 1\n",
			"en",
			true,
			true,
			nil,
		},
	}

	for k, v := range testCases {
		filePath := fmt.Sprintf("%s/%s", tempDir, v.fileName)

		if v.createFile {
			err := createTempFile(tempDir, v.fileName, v.fileContent)
			if err != nil {
				t.Fatalf("unexpected error, index=%d, error: %s", k, err)
			}
		}

		yamlFile, err := loadYAML(v.defaultLang, filePath)
		if err != nil && !v.failureExpected {
			t.Fatalf("unexpected error, index=%d, error: %s", k, err)
		}
		if err == nil && v.failureExpected {
			t.Fatalf("expected error, index=%d, error: %s", k, err)
		}

		if yamlFile != nil && !strings.HasSuffix(yamlFile.FilePath, v.fileName) {
			t.Fatalf("unexpected YAMLFile name ending, index=%d, expected=%s, actual=%s",
				k, filePath, yamlFile.FilePath)
		}

		applyTemDirLocation(v.failureExpected, []string{filePath}, v.expected)

		if !reflect.DeepEqual(yamlFile, v.expected) {
			t.Fatalf("unexpected results, index=%d, input:\n%s\nexpected:\n%+v\nactual:\n%+v\n\n",
				k, v.fileContent, v.expected, yamlFile)
		}
	}
}

func TestLoadYAMLFiles(t *testing.T) {
	tempDir := t.TempDir()

	testCases := []struct {
		fileNames       []string
		fileContents    []string
		defaultLang     string
		createFiles     bool
		failureExpected bool
		expected        []*YAMLFile
	}{
		{ // No errors - file exist and content matches.
			[]string{"file_0.yaml"},
			[]string{"key0: \"text\"\n"},
			"en",
			true,
			false,
			[]*YAMLFile{
				&YAMLFile{
					FilePath: "file_0.yaml", Translates: []Translate{
						{Key: "key0", Language: "en", Value: "text", Plural: ""},
					},
				},
			},
		},
		{ // No errors - file exist and content matches.
			[]string{"file_1.yaml", "file_2.yaml"},
			[]string{"key0: \"text1\"\n", "key1: \"text2\"\n"},
			"en",
			true,
			false,
			[]*YAMLFile{
				&YAMLFile{
					FilePath: "file_1.yaml", Translates: []Translate{
						{Key: "key0", Language: "en", Value: "text1", Plural: ""},
					},
				},
				&YAMLFile{
					FilePath: "file_2.yaml", Translates: []Translate{
						{Key: "key1", Language: "en", Value: "text2", Plural: ""},
					},
				},
			},
		},
		{ // No errors - no files passed and no results ...
			[]string{},
			[]string{},
			"en",
			false,
			false,
			nil,
		},

		{ // Error - unknown YAML format in second file.
			[]string{"file_3.yaml", "file_4.yaml"},
			[]string{"key0: \"text1\"\n", "something"},
			"en",
			true,
			true,
			nil,
		},
	}

	for k, v := range testCases {
		filePaths := make([]string, len(v.fileNames))

		if v.createFiles {
			for x, y := range v.fileNames {
				err := createTempFile(tempDir, y, v.fileContents[x])
				if err != nil {
					t.Fatalf("unexpected error, index=%d, error: %s", k, err)
				}

				filePaths[x] = fmt.Sprintf("%s/%s", tempDir, y)
			}

		}

		yamlFiles, err := LoadYAMLFiles(v.defaultLang, filePaths...)
		if err != nil && !v.failureExpected {
			t.Fatalf("unexpected error, index=%d, error: %s", k, err)
		}
		if err == nil && v.failureExpected {
			t.Fatalf("expected error, index=%d, error: %s", k, err)
		}

		for x, y := range yamlFiles {
			if !strings.HasSuffix(y.FilePath, v.fileNames[x]) {
				t.Fatalf("unexpected YAMLFile name ending, index=%d, expected=%s, actual=%s",
					k, filePaths[x], y.FilePath)
			}
		}

		applyTemDirLocation(v.failureExpected, filePaths, v.expected...)

		if !reflect.DeepEqual(yamlFiles, v.expected) {
			t.Fatalf("unexpected results, index=%d, input:\n%v\nexpected:\n%+v\nactual:\n%+v\n\n",
				k, v.fileContents, v.expected, yamlFiles)
		}
	}
}
