package config

import (
	"github.com/SaraNguyen999/setup-server/pkg/utils"
)

var (
	CONFIG Config
)

type Config struct {
	ServerConfig ServerConfig `yaml:"server"`
}

func LoadConfig(filePath string) (Config, error) {
	err := utils.ReadYAMLFile(filePath, &CONFIG)
	return CONFIG, err
}

func SaveConfig(filepath string) error {
	return utils.WriteYAMLFile(filepath, CONFIG)
}
