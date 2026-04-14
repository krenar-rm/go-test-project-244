package formatter

import (
	"encoding/json"
	"fmt"

	"code/diff"
)

// jsonNode используется для сериализации дерева в JSON-формат,
// совместимый с ожидаемой структурой вывода.
// Указатели на interface{} нужны чтобы отличать "поле отсутствует" от "значение = 0/false/пустая строка"
type jsonNode struct {
	Key      string       `json:"key"`
	Type     string       `json:"type"`
	Value1   *interface{} `json:"value1,omitempty"`
	Value2   *interface{} `json:"value2,omitempty"`
	Children []jsonNode   `json:"children,omitempty"`
}

func renderJSON(node *diff.Node) (string, error) {
	jn := convertToJSON(node)
	data, err := json.MarshalIndent(jn, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to marshal JSON: %w", err)
	}
	return string(data), nil
}

func convertToJSON(node *diff.Node) jsonNode {
	jn := jsonNode{
		Key:  node.Key,
		Type: mapNodeType(node.Type),
	}

	switch node.Type {
	case diff.TypeUnchanged:
		jn.Value1 = &node.Value
	case diff.TypeAdded:
		jn.Value2 = &node.NewValue
	case diff.TypeRemoved:
		jn.Value1 = &node.OldValue
	case diff.TypeUpdated:
		if node.OldValue != nil {
			jn.Value1 = &node.OldValue
		}
		if node.NewValue != nil {
			jn.Value2 = &node.NewValue
		}
	case diff.TypeRoot, diff.TypeNested:
		if len(node.Children) > 0 {
			jn.Children = make([]jsonNode, 0, len(node.Children))
			for _, child := range node.Children {
				jn.Children = append(jn.Children, convertToJSON(child))
			}
		}
	}

	return jn
}

func mapNodeType(t string) string {
	switch t {
	case diff.TypeRemoved:
		return "deleted"
	case diff.TypeUpdated:
		return "changed"
	default:
		return t
	}
}
