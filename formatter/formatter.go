package formatter

import (
	"fmt"
	"sort"
	"strings"

	"code/diff"
)

const (
	FormatStylish = "stylish"
	FormatPlain   = "plain"
	FormatJSON    = "json"
)

// Format форматирует дерево различий согласно указанному формату
func Format(tree *diff.Node, format string) (string, error) {
	switch strings.ToLower(format) {
	case FormatStylish:
		return renderStylish(tree), nil
	case FormatPlain:
		return renderPlain(tree), nil
	case FormatJSON:
		return renderJSON(tree)
	default:
		return "", fmt.Errorf("unsupported format: %s", format)
	}
}

func getSortedKeys(m map[string]interface{}) []string {
	keys := make([]string, 0, len(m))
	for key := range m {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	return keys
}
