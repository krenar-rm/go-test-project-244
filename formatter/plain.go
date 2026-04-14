package formatter

import (
	"fmt"
	"sort"
	"strings"

	"code/diff"
)

func renderPlain(node *diff.Node) string {
	var result []string
	renderPlainNode(node, &result, []string{})
	sort.Strings(result)
	return strings.Join(result, "\n")
}

func renderPlainNode(node *diff.Node, result *[]string, path []string) {
	for _, child := range node.Children {
		currentPath := append(path, child.Key)
		pathStr := strings.Join(currentPath, ".")

		switch child.Type {
		case diff.TypeAdded:
			*result = append(*result, fmt.Sprintf("Property '%s' was added with value: %s", pathStr, formatPlainValue(child.NewValue)))
		case diff.TypeRemoved:
			*result = append(*result, fmt.Sprintf("Property '%s' was removed", pathStr))
		case diff.TypeUpdated:
			*result = append(*result, fmt.Sprintf("Property '%s' was updated. From %s to %s", pathStr, formatPlainValue(child.OldValue), formatPlainValue(child.NewValue)))
		case diff.TypeNested:
			renderPlainNode(child, result, currentPath)
		}
	}
}

func formatPlainValue(v interface{}) string {
	if v == nil {
		return nullValue
	}
	switch val := v.(type) {
	case string:
		return fmt.Sprintf("'%s'", val)
	case bool:
		return fmt.Sprintf("%t", val)
	case map[string]interface{}:
		return "[complex value]"
	default:
		return fmt.Sprintf("%v", val)
	}
}
