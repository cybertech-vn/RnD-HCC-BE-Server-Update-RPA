package config

import (
	"fmt"
	"os"
	"path/filepath"

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
	filePath := resolveConfigPath()
	if _, err := LoadConfig(filePath); err != nil {
		panic(fmt.Errorf("load config failed (path=%s): %w", filePath, err))
	}
}

func LoadConfig(filePath string) (Config, error) {
	err := utils.ReadYAMLFile(filePath, &CONFIG)
	return CONFIG, err
}

func SaveConfig(filepath string) error {
	return utils.WriteYAMLFile(filepath, CONFIG)
}

func resolveConfigPath() string {
	// 1) Explicit override via env
	if p := os.Getenv("CONFIG_PATH"); p != "" {
		// Try as-is first
		if fileExists(p) {
			return p
		}
		// If relative, also try relative to executable dir
		if !filepath.IsAbs(p) {
			if exeDir, err := executableDir(); err == nil {
				alt := filepath.Join(exeDir, p)
				if fileExists(alt) {
					return alt
				}
			}
		}
		// Keep original (will error with a clear message on LoadConfig)
		return p
	}

	// 2) Default discovery:
	// - cwd
	// - executable dir
	// - parent of executable dir (useful when running ./build/server.exe)
	candidates := []string{
		"config.local.yml",
		"config.yml",
	}

	if exeDir, err := executableDir(); err == nil {
		candidates = append(
			candidates,
			filepath.Join(exeDir, "config.local.yml"),
			filepath.Join(exeDir, "config.yml"),
			filepath.Join(exeDir, "..", "config.local.yml"),
			filepath.Join(exeDir, "..", "config.yml"),
		)
	}

	for _, c := range candidates {
		if fileExists(c) {
			return c
		}
	}

	// Fallback to default name (will produce a readable error)
	return "config.yml"
}

func executableDir() (string, error) {
	exe, err := os.Executable()
	if err != nil {
		return "", err
	}
	exe, err = filepath.EvalSymlinks(exe)
	if err != nil {
		return "", err
	}
	return filepath.Dir(exe), nil
}

func fileExists(path string) bool {
	info, err := os.Stat(path)
	return err == nil && !info.IsDir()
}
