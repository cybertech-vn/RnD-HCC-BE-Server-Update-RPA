package config

import (
	"github.com/SaraNguyen999/setup-server/pkg/utils"
)

var (
	CONFIG Config
)

type Config struct {
	ServerConfig      ServerConfig      `yaml:"server"`
	DBConfig          DBConfig          `yaml:"database"`
	CloudServerConfig CloudServerConfig `yaml:"cloud_server"`
	MinioConfig       MinioConfig       `yaml:"minio"`
}

func init() {
	// Load default configuration or set default values
	LoadConfig("config.yml")
}

func LoadConfig(filePath string) (Config, error) {
	err := utils.ReadYAMLFile(filePath, &CONFIG)
	return CONFIG, err
}

func SaveConfig(filepath string) error {
	return utils.WriteYAMLFile(filepath, CONFIG)
}
