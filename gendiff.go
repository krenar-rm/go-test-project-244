package code

import (
	"fmt"

	"code/diff"
	"code/formatter"
	"code/parser"
)

// GenDiff сравнивает два конфигурационных файла и возвращает различия в виде строки
func GenDiff(filepath1, filepath2, format string) (string, error) {
	data1, err := parser.Parse(filepath1)
	if err != nil {
		return "", fmt.Errorf("failed to parse %s: %w", filepath1, err)
	}

	data2, err := parser.Parse(filepath2)
	if err != nil {
		return "", fmt.Errorf("failed to parse %s: %w", filepath2, err)
	}

	diffTree := diff.Build(data1, data2)

	result, err := formatter.Format(diffTree, format)
	if err != nil {
		return "", fmt.Errorf("failed to format diff: %w", err)
	}

	return result, nil
}
