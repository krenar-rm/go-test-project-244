package parser

import (
	"encoding/json"
	"fmt"

	"gopkg.in/yaml.v3"
)

const (
	FormatJSON = "json"
	FormatYAML = "yaml"
)

// Parse парсит данные в зависимости от указанного формата
func Parse(data []byte, format string) (map[string]interface{}, error) {
	switch format {
	case FormatJSON:
		return parseJSON(data)
	case FormatYAML:
		return parseYAML(data)
	default:
		return nil, fmt.Errorf("unsupported format: %s", format)
	}
}

func parseJSON(data []byte) (map[string]interface{}, error) {
	var result map[string]interface{}
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %w", err)
	}
	return result, nil
}

func parseYAML(data []byte) (map[string]interface{}, error) {
	var result map[string]interface{}
	if err := yaml.Unmarshal(data, &result); err != nil {
		return nil, fmt.Errorf("failed to parse YAML: %w", err)
	}
	return result, nil
}
