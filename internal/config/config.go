package config

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type Config struct {
	App AppConfig `yaml:"app"`
}

type AppConfig struct {
	Name    string `yaml:"name"`
	Message string `yaml:"message"`
	Timeout int    `yaml:"timeout"`
	Debug   bool   `yaml:"debug"`
}

func LoadConfig(cfgFile string) (*Config, error) {
	cfg := &Config{
		App: AppConfig{
			Name:    "World",
			Message: "Hello",
			Timeout: 30,
			Debug:   false,
		},
	}

	configPath := cfgFile
	if configPath == "" {
		configPath = findConfigFile()
	}

	if configPath == "" {
		fmt.Fprintf(os.Stderr, "Warning: No config file found, using defaults\n")
		return cfg, nil
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("error reading config file: %w", err)
	}

	if err := yaml.Unmarshal(data, cfg); err != nil {
		return nil, fmt.Errorf("error parsing config file: %w", err)
	}

	return cfg, nil
}

func findConfigFile() string {
	candidates := []string{
		"config.yaml",
		"config.yml",
	}

	homeDir, err := os.UserHomeDir()
	if err == nil {
		candidates = append(candidates,
			filepath.Join(homeDir, ".mycli", "config.yaml"),
			filepath.Join(homeDir, ".mycli", "config.yml"),
		)
	}

	for _, path := range candidates {
		if _, err := os.Stat(path); err == nil {
			return path
		}
	}

	return ""
}
