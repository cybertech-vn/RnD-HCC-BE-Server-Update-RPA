package main

import (
	"github.com/SaraNguyen999/setup-server/internal/fp"
	"github.com/SaraNguyen999/setup-server/pkg/requests"
)

func main() {
	url := "http://localhost:7776/api/v1/psk/" + fp.GenerateFingerprint()

	res, err := requests.Post(url, nil, nil, false)
	if err != nil {
		panic(err)
	}

	body, err := requests.ParseBody(*res)
	if err != nil {
		panic(err)
	}
	bodyMap, ok := body["data"].(map[string]any)
	if !ok {
		panic("unexpected response format")
	}
	println("Your passkey:", bodyMap["passkey"].(string))

}
