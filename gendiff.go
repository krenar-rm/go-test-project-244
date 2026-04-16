package code

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"code/diff"
	"code/formatter"
	"code/parser"
)

// GenDiff сравнивает два конфигурационных файла и возвращает различия в виде строки
func GenDiff(filepath1, filepath2, format string) (string, error) {
	data1, err := readFile(filepath1)
	if err != nil {
		return "", err
	}

	data2, err := readFile(filepath2)
	if err != nil {
		return "", err
	}

	ext1 := formatFromExtension(filepath1)
	ext2 := formatFromExtension(filepath2)

	parsed1, err := parser.Parse(data1, ext1)
	if err != nil {
		return "", fmt.Errorf("failed to parse %s: %w", filepath1, err)
	}

	parsed2, err := parser.Parse(data2, ext2)
	if err != nil {
		return "", fmt.Errorf("failed to parse %s: %w", filepath2, err)
	}

	diffTree := diff.Build(parsed1, parsed2)

	result, err := formatter.Format(diffTree, format)
	if err != nil {
		return "", fmt.Errorf("failed to format diff: %w", err)
	}

	return result, nil
}

func readFile(path string) ([]byte, error) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil, fmt.Errorf("file not found: %s", path)
	}

	// nolint:gosec
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read file %s: %w", path, err)
	}

	return data, nil
}

func formatFromExtension(path string) string {
	ext := strings.TrimPrefix(strings.ToLower(filepath.Ext(path)), ".")
	if ext == "yml" {
		return parser.FormatYAML
	}
	return ext
}
