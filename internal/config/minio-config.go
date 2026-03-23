package config

import (
	mini "github.com/SaraNguyen999/setup-server/pkg/database/minio"
)

type MinioConfig struct {
	mini.MinioClient `yaml:",inline"`
}
