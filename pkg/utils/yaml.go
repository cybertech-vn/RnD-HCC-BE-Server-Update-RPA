package utils

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v2"
)

func WriteYAMLFile(filename string, data interface{}) error {
	yamlData, err := yaml.Marshal(data)
	if err != nil {
		return err
	}
	err = os.WriteFile(filename, yamlData, 0644)
	if err != nil {
		return err
	}

	return nil
}

func ReadYAMLFile(filename string, result interface{}) error {
	yamlData, err := os.ReadFile(filename)
	if err != nil {
		return err
	}

	// Expand environment variables trong YAML content
	// Hỗ trợ ${VAR_NAME} hoặc ${VAR_NAME:-default}
	expandedData := os.ExpandEnv(string(yamlData))

	err = yaml.Unmarshal([]byte(expandedData), result)
	if err != nil {
		return err
	}

	return nil
}

func ConvertToStruct(data interface{}, result interface{}) error {
	// Chuyển đổi map sang YAML
	yamlData, err := yaml.Marshal(data)
	if err != nil {
		return fmt.Errorf("error marshaling map to yaml: %v", err)
	}

	// Chuyển đổi YAML sang cấu trúc
	if err := yaml.Unmarshal(yamlData, result); err != nil {
		return fmt.Errorf("error unmarshaling yaml to struct: %v", err)
	}
	return nil
}
