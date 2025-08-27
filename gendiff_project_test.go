package code

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGenDiff_GendiffProject(t *testing.T) {
	// Тестируем с файлами из go_gendiff_project
	file1 := "testdata/fixture/file1.json"
	file2 := "testdata/fixture/file2.json"

	// Проверяем, что файлы существуют
	_, err := os.Stat(file1)
	require.NoError(t, err, "file1.json should exist")

	_, err = os.Stat(file2)
	require.NoError(t, err, "file2.json should exist")

	// Тестируем stylish формат
	result, err := GenDiff(file1, file2, "stylish")
	require.NoError(t, err)

	// Проверяем основные элементы (без кавычек для stylish формата)
	assert.Contains(t, result, "+ follow: false")
	assert.Contains(t, result, "- setting2: 200")
	assert.Contains(t, result, "+ setting3:")
	assert.Contains(t, result, "key: value")
	assert.Contains(t, result, "language: js")

	// Проверяем структуру
	assert.Contains(t, result, "common: {")
	assert.Contains(t, result, "group1: {")
	assert.Contains(t, result, "group4: {")

	t.Logf("Stylish output:\n%s", result)
}
