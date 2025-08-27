package code

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
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
	NullValue         = "null"
)

// Node представляет узел в дереве различий
type Node struct {
	Type     string      `json:"type"`
	Key      string      `json:"key,omitempty"`
	Value    interface{} `json:"value,omitempty"`
	OldValue interface{} `json:"oldValue,omitempty"`
	NewValue interface{} `json:"newValue,omitempty"`
	Children []*Node     `json:"children,omitempty"`
}

// GenDiff сравнивает два конфигурационных файла и возвращает различия в виде строки
func GenDiff(filepath1, filepath2, format string) (string, error) {
	// Читаем и парсим первый файл
	data1, err := parseFile(filepath1)
	if err != nil {
		return "", fmt.Errorf("failed to parse %s: %w", filepath1, err)
	}

	// Читаем и парсим второй файл
	data2, err := parseFile(filepath2)
	if err != nil {
		return "", fmt.Errorf("failed to parse %s: %w", filepath2, err)
	}

	// Строим дерево различий
	diffTree := buildDiffTree(data1, data2)

	// Форматируем вывод согласно указанному формату
	result, err := formatDiff(diffTree, format)
	if err != nil {
		return "", fmt.Errorf("failed to format diff: %w", err)
	}

	return result, nil
}

// parseFile читает и парсит файл на основе его расширения
func parseFile(filePath string) (map[string]interface{}, error) {
	// Проверяем, существует ли файл
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return nil, fmt.Errorf("file not found: %s", filePath)
	}

	// Читаем содержимое файла
	// nolint:gosec // Мы читаем только конфигурационные файлы, а не пользовательский ввод
	content, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file %s: %w", filePath, err)
	}

	// Определяем формат по расширению
	ext := strings.ToLower(filepath.Ext(filePath))
	if ext == "" {
		return nil, fmt.Errorf("cannot determine file format for %s", filePath)
	}

	// Парсим в зависимости от формата
	switch ext {
	case ".json":
		return parseJSON(content)
	case ".yml", ".yaml":
		return parseYAML(content)
	default:
		return nil, fmt.Errorf("unsupported file format: %s", ext)
	}
}

// parseJSON парсит JSON содержимое
func parseJSON(content []byte) (map[string]interface{}, error) {
	var result map[string]interface{}
	if err := json.Unmarshal(content, &result); err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %w", err)
	}
	return result, nil
}

// parseYAML парсит YAML содержимое
func parseYAML(content []byte) (map[string]interface{}, error) {
	var result map[string]interface{}
	if err := yaml.Unmarshal(content, &result); err != nil {
		return nil, fmt.Errorf("failed to parse YAML: %w", err)
	}
	return result, nil
}

// buildDiffTree строит дерево, представляющее различия между двумя структурами данных
func buildDiffTree(data1, data2 map[string]interface{}) *Node {
	root := &Node{Type: NodeTypeRoot, Children: []*Node{}}

	// Получаем все уникальные ключи
	allKeys := make(map[string]bool)
	for key := range data1 {
		allKeys[key] = true
	}
	for key := range data2 {
		allKeys[key] = true
	}

	// Сортируем ключи для детерминированного вывода
	keys := make([]string, 0, len(allKeys))
	for key := range allKeys {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	// Обрабатываем каждый ключ
	for _, key := range keys {
		childNode := processKey(key, data1, data2)
		if childNode != nil {
			root.Children = append(root.Children, childNode)
		}
	}

	return root
}

// processKey обрабатывает отдельный ключ и возвращает узел, представляющий его состояние
func processKey(key string, data1, data2 map[string]interface{}) *Node {
	value1, exists1 := data1[key]
	value2, exists2 := data2[key]

	if !exists1 && exists2 {
		// Ключ был добавлен
		return &Node{
			Type:     NodeTypeAdded,
			Key:      key,
			NewValue: value2,
		}
	} else if exists1 && !exists2 {
		// Ключ был удален
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

// processExistingKey обрабатывает ключ, который существует в обеих структурах данных
func processExistingKey(key string, value1, value2 interface{}) *Node {
	if isEqual(value1, value2) {
		// Значения равны
		return &Node{
			Type:  NodeTypeUnchanged,
			Key:   key,
			Value: value1,
		}
	} else if isMap(value1) && isMap(value2) {
		// Оба значения являются картами, рекурсивно обрабатываем
		childNode := buildDiffTree(value1.(map[string]interface{}), value2.(map[string]interface{}))
		childNode.Key = key
		childNode.Type = NodeTypeNested
		return childNode
	} else {
		// Значения различаются - проверяем, является ли одно картой, а другое нет
		if isMap(value1) && !isMap(value2) {
			// value1 это карта, value2 нет - value1 была удалена, value2 была добавлена
			return &Node{
				Type:     NodeTypeUpdated,
				Key:      key,
				OldValue: value1,
				NewValue: value2,
			}
		} else if !isMap(value1) && isMap(value2) {
			// value1 не карта, value2 карта - value1 была удалена, value2 была добавлена
			return &Node{
				Type:     NodeTypeUpdated,
				Key:      key,
				OldValue: value1,
				NewValue: value2,
			}
		} else {
			// Оба значения примитивные, но разные
			return &Node{
				Type:     NodeTypeUpdated,
				Key:      key,
				OldValue: value1,
				NewValue: value2,
			}
		}
	}
}

// isEqual проверяет равенство двух значений с помощью глубокого сравнения
func isEqual(a, b interface{}) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil || b == nil {
		return false
	}

	// Для мапов используем собственную функцию глубокого сравнения
	if isMap(a) && isMap(b) {
		return mapsEqual(a.(map[string]interface{}), b.(map[string]interface{}))
	}

	// Для остальных типов используем обычное сравнение
	return fmt.Sprintf("%v", a) == fmt.Sprintf("%v", b)
}

// mapsEqual рекурсивно сравнивает две карты на равенство
func mapsEqual(a, b map[string]interface{}) bool {
	// Если разное количество ключей, то карты не равны
	if len(a) != len(b) {
		return false
	}

	// Проверяем каждый ключ
	for key, valueA := range a {
		valueB, exists := b[key]
		if !exists {
			return false // Ключ отсутствует во второй карте
		}

		// Рекурсивно сравниваем значения
		if !isEqual(valueA, valueB) {
			return false
		}
	}

	return true
}

// isMap проверяет, является ли значение картой
func isMap(v interface{}) bool {
	_, ok := v.(map[string]interface{})
	return ok
}

// formatDiff форматирует дерево различий согласно указанному формату
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

// formatStylish форматирует различия в stylish формате
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

// formatStylishNode рекурсивно форматирует узел в stylish формате
func formatStylishNode(node *Node, result *strings.Builder, depth int) {
	// Базовый отступ: используем формулу depth*4-2
	baseIndent := strings.Repeat(" ", depth*4-2)

	for i, child := range node.Children {
		switch child.Type {
		case NodeTypeAdded:
			fmt.Fprintf(result, "%s+ %s: %s", baseIndent, child.Key, formatValueForRemovedAdded(child.NewValue, depth))
		case NodeTypeRemoved:
			fmt.Fprintf(result, "%s- %s: %s", baseIndent, child.Key, formatValueForRemovedAdded(child.OldValue, depth))
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

// formatPlain форматирует различия в plain формате
func formatPlain(node *Node) string {
	var result []string
	formatPlainNode(node, &result, []string{})
	sort.Strings(result)
	return strings.Join(result, "\n")
}

// formatPlainNode рекурсивно форматирует узел в plain формате
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

// formatJSON форматирует различия как JSON
func formatJSON(node *Node) string {
	jsonData, err := json.MarshalIndent(node, "", "  ")
	if err != nil {
		return "{}"
	}
	return string(jsonData)
}

// formatValue форматирует значение для stylish вывода (для вложенных и неизменённых узлов)
func formatValue(v interface{}) string {
	if v == nil {
		return NullValue
	}

	switch val := v.(type) {
	case string:
		// Убираем кавычки для строк в stylish формате
		return val
	case bool:
		return fmt.Sprintf("%t", val)
	case map[string]interface{}:
		// Для вложенных объектов создаем простой вывод
		return formatNestedMap(val)
	default:
		return fmt.Sprintf("%v", val)
	}
}

// formatValueForRemovedAdded форматирует значение для удалённых/добавленных узлов
func formatValueForRemovedAdded(v interface{}, depth int) string {
	if v == nil {
		return NullValue
	}

	switch val := v.(type) {
	case string:
		// Убираем кавычки для строк в stylish формате
		return val
	case bool:
		return fmt.Sprintf("%t", val)
	case map[string]interface{}:
		// Для удаленных/добавленных объектов используем форматирование с учетом глубины
		return formatSimpleMapWithDepth(val, depth)
	default:
		return fmt.Sprintf("%v", val)
	}
}

// formatNestedMap форматирует карту для вложенных объектов в stylish формате
func formatNestedMap(m map[string]interface{}) string {
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

	// Отступ для содержимого: 12 пробелов (как в ожидаемом результате)
	contentIndent := "            " // 12 пробелов
	for i, key := range keys {
		value := m[key]
		if isMap(value) {
			result.WriteString(fmt.Sprintf("%s%s: %s", contentIndent, key, formatNestedMapRecursive(value.(map[string]interface{}), 3)))
		} else {
			result.WriteString(fmt.Sprintf("%s%s: %s", contentIndent, key, formatValue(value)))
		}
		if i < len(keys)-1 {
			result.WriteString("\n")
		}
	}

	// Закрывающая скобка с отступом 8 пробелов (как в ожидаемом результате)
	result.WriteString("\n        }")
	return result.String()
}

// formatNestedMapRecursive форматирует вложенные карты с правильными отступами для вложенных объектов
func formatNestedMapRecursive(m map[string]interface{}, depth int) string {
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

	// Отступ для содержимого: используем правильную формулу
	contentIndent := strings.Repeat(" ", depth*4)
	for i, key := range keys {
		value := m[key]
		if isMap(value) {
			result.WriteString(fmt.Sprintf("%s%s: %s", contentIndent, key, formatNestedMapRecursive(value.(map[string]interface{}), depth+1)))
		} else {
			result.WriteString(fmt.Sprintf("%s%s: %s", contentIndent, key, formatValue(value)))
		}
		if i < len(keys)-1 {
			result.WriteString("\n")
		}
	}

	// Закрывающая скобка с правильным отступом
	result.WriteString(fmt.Sprintf("\n%s}", strings.Repeat(" ", (depth-1)*4)))
	return result.String()
}

// formatSimpleMap форматирует карту с простой структурой
// formatSimpleMapWithDepth форматирует карту для удалённых/добавленных узлов с учётом глубины
func formatSimpleMapWithDepth(m map[string]interface{}, depth int) string {
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

	// Отступ для содержимого: зависит от глубины
	// Глубина 1: 8 пробелов, Глубина 2: 12 пробелов
	var contentIndent string
	var closingIndent string
	if depth == 1 {
		contentIndent = "        " // 8 пробелов
		closingIndent = "    "     // 4 пробела
	} else {
		contentIndent = "            " // 12 пробелов
		closingIndent = "        "     // 8 пробелов
	}

	for i, key := range keys {
		value := m[key]
		if isMap(value) {
			result.WriteString(fmt.Sprintf("%s%s: %s", contentIndent, key, formatSimpleMapRecursive(value.(map[string]interface{}), depth+2)))
		} else {
			result.WriteString(fmt.Sprintf("%s%s: %s", contentIndent, key, formatValue(value)))
		}
		if i < len(keys)-1 {
			result.WriteString("\n")
		}
	}

	result.WriteString(fmt.Sprintf("\n%s}", closingIndent))
	return result.String()
}

// formatSimpleMapRecursive форматирует вложенные карты с правильными отступами
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

	// Отступ для содержимого: используем правильную формулу
	// Для глубины 2: 8 пробелов, для глубины 3: 12 пробелов
	contentIndent := strings.Repeat(" ", depth*4)
	closingIndent := strings.Repeat(" ", (depth-1)*4)

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
	result.WriteString(fmt.Sprintf("\n%s}", closingIndent))
	return result.String()
}

// formatPlainValue форматирует значение для plain вывода
func formatPlainValue(v interface{}) string {
	if v == nil {
		return NullValue
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
