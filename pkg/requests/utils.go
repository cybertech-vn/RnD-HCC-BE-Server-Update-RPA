package requests

import (
	"encoding/json"
	"fmt"
	"log"
)

func ParseBody(resp Response) (map[string]any, error) {
	var parsed map[string]any
	if body, ok := resp.Body.([]byte); ok {
		if err := json.Unmarshal(body, &parsed); err != nil {
			log.Fatal("parse error:", err)
			return nil, err
		}
		return parsed, nil
	}
	return nil, fmt.Errorf("unexpected response body type")
}
