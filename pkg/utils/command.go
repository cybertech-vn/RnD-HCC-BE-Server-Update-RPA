package utils

import (
	"os"
	"strings"
)

func GetArgValue(key string) string {
	for _, arg := range os.Args {
		if strings.HasPrefix(arg, key+"=") {
			return strings.TrimPrefix(arg, key+"=")
		}
	}
	return ""
}
