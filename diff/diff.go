package diff

import (
	"fmt"
	"sort"
)

const (
	TypeRoot      = "root"
	TypeAdded     = "added"
	TypeRemoved   = "removed"
	TypeUpdated   = "updated"
	TypeUnchanged = "unchanged"
	TypeNested    = "nested"
)

// Node представляет узел в дереве различий
type Node struct {
	Type     string
	Key      string
	Value    interface{}
	OldValue interface{}
	NewValue interface{}
	Children []*Node
}

// Build строит дерево, представляющее различия между двумя структурами данных
func Build(data1, data2 map[string]interface{}) *Node {
	root := &Node{Type: TypeRoot, Children: []*Node{}}

	keys := uniqueSortedKeys(data1, data2)

	for _, key := range keys {
		if child := processKey(key, data1, data2); child != nil {
			root.Children = append(root.Children, child)
		}
	}

	return root
}

func uniqueSortedKeys(data1, data2 map[string]interface{}) []string {
	allKeys := make(map[string]bool)
	for key := range data1 {
		allKeys[key] = true
	}
	for key := range data2 {
		allKeys[key] = true
	}

	keys := make([]string, 0, len(allKeys))
	for key := range allKeys {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	return keys
}

func processKey(key string, data1, data2 map[string]interface{}) *Node {
	value1, exists1 := data1[key]
	value2, exists2 := data2[key]

	switch {
	case !exists1 && exists2:
		return &Node{Type: TypeAdded, Key: key, NewValue: value2}
	case exists1 && !exists2:
		return &Node{Type: TypeRemoved, Key: key, OldValue: value1}
	case exists1 && exists2:
		return processExistingKey(key, value1, value2)
	}

	return nil
}

func processExistingKey(key string, value1, value2 interface{}) *Node {
	if isEqual(value1, value2) {
		return &Node{Type: TypeUnchanged, Key: key, Value: value1}
	}

	if isMap(value1) && isMap(value2) {
		child := Build(value1.(map[string]interface{}), value2.(map[string]interface{}))
		child.Key = key
		child.Type = TypeNested
		return child
	}

	return &Node{Type: TypeUpdated, Key: key, OldValue: value1, NewValue: value2}
}

func isEqual(a, b interface{}) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil || b == nil {
		return false
	}

	if isMap(a) && isMap(b) {
		return mapsEqual(a.(map[string]interface{}), b.(map[string]interface{}))
	}

	return fmt.Sprintf("%v", a) == fmt.Sprintf("%v", b)
}

func mapsEqual(a, b map[string]interface{}) bool {
	if len(a) != len(b) {
		return false
	}

	for key, valueA := range a {
		valueB, exists := b[key]
		if !exists {
			return false
		}
		if !isEqual(valueA, valueB) {
			return false
		}
	}

	return true
}

func isMap(v interface{}) bool {
	_, ok := v.(map[string]interface{})
	return ok
}
