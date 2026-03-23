package config

import "github.com/SaraNguyen999/setup-server/pkg/database"

type DBConfig struct {
	database.BaseConnect `yaml:",inline"`
}
