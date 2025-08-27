package code

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenDiff_Integration(t *testing.T) {
	getFixturePath := func(filename string) string {
		return filepath.Join("testdata", "fixture", filename)
	}

	readExpected := func(filename string) string {
		raw, err := os.ReadFile(getFixturePath(filename))
		assert.NoError(t, err)
		result := string(raw)
		return result[:len(result)-1] // Trim trailing newline
	}

	expectedStylish := readExpected("result_stylish.txt")

	// Тестируем JSON и YAML входные форматы
	inputFormats := []string{"json", "yml"}

	for _, inputFormat := range inputFormats {
		file1 := getFixturePath("file1." + inputFormat)
		file2 := getFixturePath("file2." + inputFormat)

		t.Run(inputFormat+"_stylish", func(t *testing.T) {
			result, err := GenDiff(file1, file2, "stylish")
			assert.NoError(t, err)
			assert.Equal(t, expectedStylish, result)
		})
	}
}
