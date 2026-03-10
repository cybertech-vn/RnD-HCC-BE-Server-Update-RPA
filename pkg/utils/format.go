package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"reflect"
)

// Format JSON data
func FormatJSON(data []byte) string {
	var prettyJSON bytes.Buffer
	if err := json.Indent(&prettyJSON, data, "", "  "); err != nil {
		panic(fmt.Errorf("failed to parse JSON: %w", err))
	}
	return prettyJSON.String()
}

// ConvertStructToArgs nhận bất kỳ struct nào và trả về []string chứa giá trị các field
func ConvertStructToArgs(input interface{}) ([]string, error) {
	v := reflect.ValueOf(input)

	// Nếu là con trỏ thì dereference
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	// Kiểm tra có phải struct không
	if v.Kind() != reflect.Struct {
		return nil, fmt.Errorf("expected a struct, got %s", v.Kind())
	}

	args := []string{}
	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)

		// Bỏ qua field không export
		if !field.CanInterface() {
			continue
		}

		val := fmt.Sprintf("%v", field.Interface())
		args = append(args, val)
	}

	return args, nil
}
