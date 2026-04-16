package formatter

import (
	"fmt"
	"sort"
	"strings"

	"code/diff"
)

const nullValue = "null"

func renderStylish(node *diff.Node) string {
	var result strings.Builder
	result.WriteString("{\n")
	renderStylishNode(node, &result, 1)
	if len(node.Children) > 0 {
		result.WriteString("\n")
	}
	result.WriteString("}")
	return result.String()
}

func renderStylishNode(node *diff.Node, result *strings.Builder, depth int) {
	baseIndent := strings.Repeat(" ", depth*4-2)

	for i, child := range node.Children {
		switch child.Type {
		case diff.TypeAdded:
			fmt.Fprintf(result, "%s+ %s: %s", baseIndent, child.Key, formatValueForRemovedAdded(child.NewValue, depth))
		case diff.TypeRemoved:
			fmt.Fprintf(result, "%s- %s: %s", baseIndent, child.Key, formatValueForRemovedAdded(child.OldValue, depth))
		case diff.TypeUpdated:
			fmt.Fprintf(result, "%s- %s: %s\n%s+ %s: %s",
				baseIndent, child.Key, formatValue(child.OldValue),
				baseIndent, child.Key, formatValue(child.NewValue))
		case diff.TypeUnchanged:
			fmt.Fprintf(result, "%s  %s: %s", baseIndent, child.Key, formatValue(child.Value))
		case diff.TypeNested:
			fmt.Fprintf(result, "%s  %s: {\n", baseIndent, child.Key)
			renderStylishNode(child, result, depth+1)
			fmt.Fprintf(result, "\n%s  }", baseIndent)
		}

		if i < len(node.Children)-1 {
			result.WriteString("\n")
		}
	}
}

func formatValue(v interface{}) string {
	if v == nil {
		return nullValue
	}

	if m, ok := v.(map[string]interface{}); ok {
		return formatNestedMap(m)
	}

	return formatPrimitiveValue(v)
}

func formatValueForRemovedAdded(v interface{}, depth int) string {
	if v == nil {
		return nullValue
	}

	if m, ok := v.(map[string]interface{}); ok {
		return formatSimpleMapWithDepth(m, depth)
	}

	return formatPrimitiveValue(v)
}

func formatPrimitiveValue(v interface{}) string {
	switch val := v.(type) {
	case string:
		return val
	case bool:
		return fmt.Sprintf("%t", val)
	default:
		return fmt.Sprintf("%v", val)
	}
}

func formatNestedMap(m map[string]interface{}) string {
	if len(m) == 0 {
		return "{}"
	}

	var result strings.Builder
	result.WriteString("{\n")

	keys := getSortedKeys(m)

	contentIndent := "            "
	for i, key := range keys {
		value := m[key]
		if _, ok := value.(map[string]interface{}); ok {
			fmt.Fprintf(&result, "%s%s: %s", contentIndent, key, formatNestedMapRecursive(value.(map[string]interface{}), 3))
		} else {
			fmt.Fprintf(&result, "%s%s: %s", contentIndent, key, formatValue(value))
		}
		if i < len(keys)-1 {
			result.WriteString("\n")
		}
	}

	result.WriteString("\n        }")
	return result.String()
}

func formatNestedMapRecursive(m map[string]interface{}, depth int) string {
	if len(m) == 0 {
		return "{}"
	}

	var result strings.Builder
	result.WriteString("{\n")

	keys := getSortedKeys(m)

	contentIndent := strings.Repeat(" ", depth*4)
	for i, key := range keys {
		value := m[key]
		if _, ok := value.(map[string]interface{}); ok {
			fmt.Fprintf(&result, "%s%s: %s", contentIndent, key, formatNestedMapRecursive(value.(map[string]interface{}), depth+1))
		} else {
			fmt.Fprintf(&result, "%s%s: %s", contentIndent, key, formatValue(value))
		}
		if i < len(keys)-1 {
			result.WriteString("\n")
		}
	}

	fmt.Fprintf(&result, "\n%s}", strings.Repeat(" ", (depth-1)*4))
	return result.String()
}

func formatSimpleMapWithDepth(m map[string]interface{}, depth int) string {
	if len(m) == 0 {
		return "{}"
	}

	var result strings.Builder
	result.WriteString("{\n")

	keys := getSortedKeys(m)

	contentIndent := strings.Repeat(" ", (depth+1)*4)
	closingIndent := strings.Repeat(" ", depth*4)

	for i, key := range keys {
		value := m[key]
		if _, ok := value.(map[string]interface{}); ok {
			fmt.Fprintf(&result, "%s%s: %s", contentIndent, key, formatSimpleMapRecursive(value.(map[string]interface{}), depth+2))
		} else {
			fmt.Fprintf(&result, "%s%s: %s", contentIndent, key, formatValue(value))
		}
		if i < len(keys)-1 {
			result.WriteString("\n")
		}
	}

	fmt.Fprintf(&result, "\n%s}", closingIndent)
	return result.String()
}

func formatSimpleMapRecursive(m map[string]interface{}, depth int) string {
	if len(m) == 0 {
		return "{}"
	}

	var result strings.Builder
	result.WriteString("{\n")

	keys := getSortedKeys(m)

	contentIndent := strings.Repeat(" ", depth*4)
	closingIndent := strings.Repeat(" ", (depth-1)*4)

	for i, key := range keys {
		value := m[key]
		if _, ok := value.(map[string]interface{}); ok {
			fmt.Fprintf(&result, "%s%s: %s", contentIndent, key, formatSimpleMapRecursive(value.(map[string]interface{}), depth+1))
		} else {
			fmt.Fprintf(&result, "%s%s: %s", contentIndent, key, formatValue(value))
		}
		if i < len(keys)-1 {
			result.WriteString("\n")
		}
	}

	fmt.Fprintf(&result, "\n%s}", closingIndent)
	return result.String()
}

func getSortedKeys(m map[string]interface{}) []string {
	keys := make([]string, 0, len(m))
	for key := range m {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	return keys
}
