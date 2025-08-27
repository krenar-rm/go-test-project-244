package code

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"strings"

	"gopkg.in/yaml.v3"
)

const (
	NodeTypeRoot      = "root"
	NodeTypeAdded     = "added"
	NodeTypeRemoved   = "removed"
	NodeTypeUpdated   = "updated"
	NodeTypeUnchanged = "unchanged"
	NodeTypeNested    = "nested"
)

// Node represents a node in the diff tree
type Node struct {
	Type     string      `json:"type"`
	Key      string      `json:"key,omitempty"`
	Value    interface{} `json:"value,omitempty"`
	OldValue interface{} `json:"oldValue,omitempty"`
	NewValue interface{} `json:"newValue,omitempty"`
	Children []*Node     `json:"children,omitempty"`
}

// GenDiff compares two configuration files and returns the difference as a string
func GenDiff(filepath1, filepath2, format string) (string, error) {
	// Read and parse first file
	data1, err := parseFile(filepath1)
	if err != nil {
		return "", fmt.Errorf("failed to parse %s: %w", filepath1, err)
	}

	// Read and parse second file
	data2, err := parseFile(filepath2)
	if err != nil {
		return "", fmt.Errorf("failed to parse %s: %w", filepath2, err)
	}

	// Build diff tree
	diffTree := buildDiffTree(data1, data2)

	// Format output according to specified format
	result, err := formatDiff(diffTree, format)
	if err != nil {
		return "", fmt.Errorf("failed to format diff: %w", err)
	}

	return result, nil
}

// parseFile reads and parses a file based on its extension
func parseFile(filePath string) (map[string]interface{}, error) {
	// Check if file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return nil, fmt.Errorf("file not found: %s", filePath)
	}

	// Read file content
	// nolint:gosec // We only read configuration files, not user input
	content, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file %s: %w", filePath, err)
	}

	// Determine format by extension
	ext := strings.ToLower(filepath.Ext(filePath))
	if ext == "" {
		return nil, fmt.Errorf("cannot determine file format for %s", filePath)
	}

	// Parse based on format
	switch ext {
	case ".json":
		return parseJSON(content)
	case ".yml", ".yaml":
		return parseYAML(content)
	default:
		return nil, fmt.Errorf("unsupported file format: %s", ext)
	}
}

// parseJSON parses JSON content
func parseJSON(content []byte) (map[string]interface{}, error) {
	var result map[string]interface{}
	if err := json.Unmarshal(content, &result); err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %w", err)
	}
	return result, nil
}

// parseYAML parses YAML content
func parseYAML(content []byte) (map[string]interface{}, error) {
	var result map[string]interface{}
	if err := yaml.Unmarshal(content, &result); err != nil {
		return nil, fmt.Errorf("failed to parse YAML: %w", err)
	}
	return result, nil
}

// buildDiffTree builds a tree representing differences between two data structures
func buildDiffTree(data1, data2 map[string]interface{}) *Node {
	root := &Node{Type: NodeTypeRoot, Children: []*Node{}}

	// Get all unique keys
	allKeys := make(map[string]bool)
	for key := range data1 {
		allKeys[key] = true
	}
	for key := range data2 {
		allKeys[key] = true
	}

	// Sort keys for deterministic output
	keys := make([]string, 0, len(allKeys))
	for key := range allKeys {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	// Process each key
	for _, key := range keys {
		childNode := processKey(key, data1, data2)
		if childNode != nil {
			root.Children = append(root.Children, childNode)
		}
	}

	return root
}

// processKey processes a single key and returns a node representing its state
func processKey(key string, data1, data2 map[string]interface{}) *Node {
	value1, exists1 := data1[key]
	value2, exists2 := data2[key]

	if !exists1 && exists2 {
		// Key was added
		return &Node{
			Type:     NodeTypeAdded,
			Key:      key,
			NewValue: value2,
		}
	} else if exists1 && !exists2 {
		// Key was removed
		return &Node{
			Type:     NodeTypeRemoved,
			Key:      key,
			OldValue: value1,
		}
	} else if exists1 && exists2 {
		return processExistingKey(key, value1, value2)
	}

	return nil
}

// processExistingKey processes a key that exists in both data structures
func processExistingKey(key string, value1, value2 interface{}) *Node {
	if isEqual(value1, value2) {
		// Values are equal
		return &Node{
			Type:  NodeTypeUnchanged,
			Key:   key,
			Value: value1,
		}
	} else if isMap(value1) && isMap(value2) {
		// Both values are maps, recurse
		childNode := buildDiffTree(value1.(map[string]interface{}), value2.(map[string]interface{}))
		childNode.Key = key
		childNode.Type = NodeTypeNested
		return childNode
	} else {
		// Values are different - check if one is map and other is not
		if isMap(value1) && !isMap(value2) {
			// value1 is map, value2 is not - value1 was removed, value2 was added
			return &Node{
				Type:     NodeTypeUpdated,
				Key:      key,
				OldValue: value1,
				NewValue: value2,
			}
		} else if !isMap(value1) && isMap(value2) {
			// value1 is not map, value2 is map - value1 was removed, value2 was added
			return &Node{
				Type:     NodeTypeUpdated,
				Key:      key,
				OldValue: value1,
				NewValue: value2,
			}
		} else {
			// Both are primitive values, but different
			return &Node{
				Type:     NodeTypeUpdated,
				Key:      key,
				OldValue: value1,
				NewValue: value2,
			}
		}
	}
}

// isEqual checks if two values are equal using deep comparison
func isEqual(a, b interface{}) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil || b == nil {
		return false
	}

	// Для мапов используем глубокое сравнение
	if isMap(a) && isMap(b) {
		return reflect.DeepEqual(a, b)
	}

	// Для остальных типов используем обычное сравнение
	return fmt.Sprintf("%v", a) == fmt.Sprintf("%v", b)
}

// isMap checks if a value is a map
func isMap(v interface{}) bool {
	_, ok := v.(map[string]interface{})
	return ok
}

// formatDiff formats the diff tree according to the specified format
func formatDiff(diffTree *Node, format string) (string, error) {
	switch strings.ToLower(format) {
	case "stylish":
		return formatStylish(diffTree), nil
	case "plain":
		return formatPlain(diffTree), nil
	case "json":
		return formatJSON(diffTree), nil
	default:
		return "", fmt.Errorf("unsupported format: %s", format)
	}
}

// formatStylish formats the diff in stylish format
func formatStylish(node *Node) string {
	var result strings.Builder
	result.WriteString("{\n")
	formatStylishNode(node, &result, 1)
	// Убираем лишний перенос строки, если нет дочерних элементов
	if len(node.Children) > 0 {
		result.WriteString("\n")
	}
	result.WriteString("}")
	return result.String()
}

// formatStylishNode recursively formats a node in stylish format
func formatStylishNode(node *Node, result *strings.Builder, depth int) {
	// Базовый отступ: используем формулу depth*4-2
	baseIndent := strings.Repeat(" ", depth*4-2)

	for i, child := range node.Children {
		switch child.Type {
		case NodeTypeAdded:
			fmt.Fprintf(result, "%s+ %s: %s", baseIndent, child.Key, formatValue(child.NewValue))
		case NodeTypeRemoved:
			fmt.Fprintf(result, "%s- %s: %s", baseIndent, child.Key, formatValue(child.OldValue))
		case NodeTypeUpdated:
			fmt.Fprintf(result, "%s- %s: %s\n%s+ %s: %s",
				baseIndent, child.Key, formatValue(child.OldValue),
				baseIndent, child.Key, formatValue(child.NewValue))
		case NodeTypeUnchanged:
			fmt.Fprintf(result, "%s  %s: %s", baseIndent, child.Key, formatValue(child.Value))
		case NodeTypeNested:
			fmt.Fprintf(result, "%s  %s: {\n", baseIndent, child.Key)
			formatStylishNode(child, result, depth+1)
			fmt.Fprintf(result, "\n%s  }", baseIndent)
		}

		// Добавляем перенос строки между элементами, кроме последнего
		if i < len(node.Children)-1 {
			result.WriteString("\n")
		}
	}
}

// formatPlain formats the diff in plain format
func formatPlain(node *Node) string {
	var result []string
	formatPlainNode(node, &result, []string{})
	sort.Strings(result)
	return strings.Join(result, "\n")
}

// formatPlainNode recursively formats a node in plain format
func formatPlainNode(node *Node, result *[]string, path []string) {
	for _, child := range node.Children {
		currentPath := append(path, child.Key)
		pathStr := strings.Join(currentPath, ".")

		switch child.Type {
		case NodeTypeAdded:
			*result = append(*result, fmt.Sprintf("Property '%s' was added with value: %s", pathStr, formatPlainValue(child.NewValue)))
		case NodeTypeRemoved:
			*result = append(*result, fmt.Sprintf("Property '%s' was removed", pathStr))
		case NodeTypeUpdated:
			*result = append(*result, fmt.Sprintf("Property '%s' was updated. From %s to %s", pathStr, formatPlainValue(child.OldValue), formatPlainValue(child.NewValue)))
		case NodeTypeNested:
			formatPlainNode(child, result, currentPath)
		}
	}
}

// formatJSON formats the diff as JSON
func formatJSON(node *Node) string {
	jsonData, err := json.MarshalIndent(node, "", "  ")
	if err != nil {
		return "{}"
	}
	return string(jsonData)
}

// formatValue formats a value for stylish output
func formatValue(v interface{}) string {
	if v == nil {
		return "null"
	}

	switch val := v.(type) {
	case string:
		// Убираем кавычки для строк в stylish формате
		return val
	case bool:
		return fmt.Sprintf("%t", val)
	case map[string]interface{}:
		// Для вложенных объектов создаем простой вывод
		return formatSimpleMap(val)
	default:
		return fmt.Sprintf("%v", val)
	}
}

// formatSimpleMap formats a map with simple structure
func formatSimpleMap(m map[string]interface{}) string {
	if len(m) == 0 {
		return "{}"
	}

	var result strings.Builder
	result.WriteString("{\n")

	// Сортируем ключи для детерминированного вывода
	keys := make([]string, 0, len(m))
	for key := range m {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	// Отступ для содержимого: используем формулу (depth+1)*4-2
	contentIndent := strings.Repeat(" ", 6) // (1+1)*4-2 = 6
	for i, key := range keys {
		value := m[key]
		if isMap(value) {
			result.WriteString(fmt.Sprintf("%s%s: %s", contentIndent, key, formatSimpleMapRecursive(value.(map[string]interface{}), 1)))
		} else {
			result.WriteString(fmt.Sprintf("%s%s: %s", contentIndent, key, formatValue(value)))
		}
		if i < len(keys)-1 {
			result.WriteString("\n")
		}
	}

	result.WriteString("\n    }")
	return result.String()
}

// formatSimpleMapRecursive formats nested maps with proper indentation
func formatSimpleMapRecursive(m map[string]interface{}, depth int) string {
	if len(m) == 0 {
		return "{}"
	}

	var result strings.Builder
	result.WriteString("{\n")

	// Сортируем ключи для детерминированного вывода
	keys := make([]string, 0, len(m))
	for key := range m {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	// Отступ для содержимого: используем формулу (depth+1)*4-2
	contentIndent := strings.Repeat(" ", (depth+1)*4-2)
	for i, key := range keys {
		value := m[key]
		if isMap(value) {
			result.WriteString(fmt.Sprintf("%s%s: %s", contentIndent, key, formatSimpleMapRecursive(value.(map[string]interface{}), depth+1)))
		} else {
			result.WriteString(fmt.Sprintf("%s%s: %s", contentIndent, key, formatValue(value)))
		}
		if i < len(keys)-1 {
			result.WriteString("\n")
		}
	}

	// Закрывающая скобка с правильным отступом
	result.WriteString(fmt.Sprintf("\n%s}", strings.Repeat(" ", depth*4-2)))
	return result.String()
}

// formatPlainValue formats a value for plain output
func formatPlainValue(v interface{}) string {
	if v == nil {
		return "null"
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
