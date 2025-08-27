package go_test_project_2

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGenDiff_JSON_Simple(t *testing.T) {
	// Create temporary test files
	file1 := createTempFile(t, `{"host":"hexlet.io","timeout":50,"proxy":"123.234.53.22","follow":false}`)
	file2 := createTempFile(t, `{"timeout":20,"verbose":true,"host":"hexlet.io"}`)
	defer os.Remove(file1)
	defer os.Remove(file2)

	// Test stylish format
	result, err := GenDiff(file1, file2, "stylish")
	require.NoError(t, err)

	// Check that all expected changes are present
	assert.Contains(t, result, "- follow: false")
	assert.Contains(t, result, "- proxy: \"123.234.53.22\"")
	assert.Contains(t, result, "- timeout: 50")
	assert.Contains(t, result, "+ timeout: 20")
	assert.Contains(t, result, "+ verbose: true")
	assert.Contains(t, result, "host: \"hexlet.io\"")

	// Test plain format
	result, err = GenDiff(file1, file2, "plain")
	require.NoError(t, err)
	assert.Contains(t, result, "Property 'follow' was removed")
	assert.Contains(t, result, "Property 'proxy' was removed")
	assert.Contains(t, result, "Property 'timeout' was updated")
	assert.Contains(t, result, "Property 'verbose' was added")

	// Test JSON format
	result, err = GenDiff(file1, file2, "json")
	require.NoError(t, err)
	// Verify JSON structure
	assert.Contains(t, result, "type")
	assert.Contains(t, result, "children")
}

func TestGenDiff_YAML_Simple(t *testing.T) {
	// Create temporary test files
	file1 := createTempYAMLFile(t, `host: hexlet.io
timeout: 50
proxy: 123.234.53.22
follow: false`)
	file2 := createTempYAMLFile(t, `timeout: 20
verbose: true
host: hexlet.io`)
	defer os.Remove(file1)
	defer os.Remove(file2)

	result, err := GenDiff(file1, file2, "stylish")
	require.NoError(t, err)
	assert.Contains(t, result, "follow: false")
	assert.Contains(t, result, "verbose: true")
}

func TestGenDiff_FileNotFound(t *testing.T) {
	_, err := GenDiff("nonexistent1.json", "nonexistent2.json", "stylish")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "file not found")
}

func TestGenDiff_UnsupportedFormat(t *testing.T) {
	file1 := createTempFile(t, `{"test": "data"}`)
	file2 := createTempFile(t, `{"test": "data"}`)
	defer os.Remove(file1)
	defer os.Remove(file2)

	_, err := GenDiff(file1, file2, "unsupported")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "unsupported format")
}

func TestGenDiff_InvalidJSON(t *testing.T) {
	file1 := createTempFile(t, `{"invalid": json}`)
	file2 := createTempFile(t, `{"test": "data"}`)
	defer os.Remove(file1)
	defer os.Remove(file2)

	_, err := GenDiff(file1, file2, "stylish")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to parse JSON")
}

func TestGenDiff_EmptyFiles(t *testing.T) {
	file1 := createTempFile(t, `{}`)
	file2 := createTempFile(t, `{}`)
	defer os.Remove(file1)
	defer os.Remove(file2)

	result, err := GenDiff(file1, file2, "stylish")
	require.NoError(t, err)
	assert.Equal(t, "{\n}", result)
}

func TestGenDiff_NestedStructures(t *testing.T) {
	file1 := createTempFile(t, `{
  "common": {
    "setting1": "Value 1",
    "setting2": 200,
    "setting3": true
  },
  "group1": {
    "baz": "bas",
    "foo": "bar"
  }
}`)
	file2 := createTempFile(t, `{
  "common": {
    "setting1": "Value 1",
    "setting2": 200,
    "setting3": false
  },
  "group1": {
    "baz": "bas",
    "foo": "bar"
  },
  "group2": {
    "abc": 12345
  }
}`)
	defer os.Remove(file1)
	defer os.Remove(file2)

	result, err := GenDiff(file1, file2, "stylish")
	require.NoError(t, err)
	assert.Contains(t, result, "setting3: true")
	assert.Contains(t, result, "setting3: false")
	assert.Contains(t, result, "group2:")
}

// Helper function to create temporary files
func createTempFile(t *testing.T, content string) string {
	// Create temporary file with .json extension for JSON content
	tmpfile, err := os.CreateTemp("", "gendiff_test_*.json")
	require.NoError(t, err)

	_, err = tmpfile.WriteString(content)
	require.NoError(t, err)

	err = tmpfile.Close()
	require.NoError(t, err)

	return tmpfile.Name()
}

// Helper function to create temporary YAML files
func createTempYAMLFile(t *testing.T, content string) string {
	// Create temporary file with .yml extension for YAML content
	tmpfile, err := os.CreateTemp("", "gendiff_test_*.yml")
	require.NoError(t, err)

	_, err = tmpfile.WriteString(content)
	require.NoError(t, err)

	err = tmpfile.Close()
	require.NoError(t, err)

	return tmpfile.Name()
}
