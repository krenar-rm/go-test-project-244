package code

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func getFixturePath(filename string) string {
	return filepath.Join("testdata", "fixture", filename)
}

func readFixture(t *testing.T, filename string) string {
	t.Helper()
	raw, err := os.ReadFile(getFixturePath(filename))
	require.NoError(t, err)
	content := string(raw)
	// убираем завершающий перенос строки в фикстурах
	if len(content) > 0 && content[len(content)-1] == '\n' {
		content = content[:len(content)-1]
	}
	return content
}

func TestGenDiff(t *testing.T) {
	formats := []struct {
		name       string
		format     string
		resultFile string
	}{
		{"stylish", "stylish", "result_stylish.txt"},
		{"plain", "plain", "result_plain.txt"},
		{"json", "json", "result_json.json"},
	}

	inputFormats := []string{"json", "yml"}

	for _, fmt := range formats {
		expected := readFixture(t, fmt.resultFile)

		for _, ext := range inputFormats {
			name := ext + "_" + fmt.name
			t.Run(name, func(t *testing.T) {
				file1 := getFixturePath("file1." + ext)
				file2 := getFixturePath("file2." + ext)

				result, err := GenDiff(file1, file2, fmt.format)
				require.NoError(t, err)
				assert.Equal(t, expected, result)
			})
		}
	}
}

func TestGenDiffErrors(t *testing.T) {
	cases := []struct {
		name     string
		file1    string
		file2    string
		format   string
		errorMsg string
	}{
		{
			name:     "file not found",
			file1:    "nonexistent.json",
			file2:    "nonexistent.json",
			format:   "stylish",
			errorMsg: "file not found",
		},
		{
			name:     "unsupported format",
			file1:    getFixturePath("file1.json"),
			file2:    getFixturePath("file2.json"),
			format:   "unknown",
			errorMsg: "unsupported format",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := GenDiff(tc.file1, tc.file2, tc.format)
			require.Error(t, err)
			assert.Contains(t, err.Error(), tc.errorMsg)
		})
	}
}
