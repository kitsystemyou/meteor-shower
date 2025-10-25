package config

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type Config struct {
	LoadTest LoadTestConfig `yaml:"loadtest"`
}

type LoadTestConfig struct {
	Domain      string     `yaml:"domain"`
	Endpoints   []Endpoint `yaml:"endpoints"`
	RPS         int        `yaml:"rps"`
	Concurrency int        `yaml:"concurrency"`
	Duration    int        `yaml:"duration"`
	Output      string     `yaml:"output"`
}

type Endpoint struct {
	Path   string  `yaml:"path"`
	Weight float64 `yaml:"weight"`
}

func LoadConfig(cfgFile string) (*Config, error) {
	cfg := &Config{
		LoadTest: LoadTestConfig{
			Domain: "http://localhost:8080",
			Endpoints: []Endpoint{
				{Path: "/", Weight: 1.0},
			},
			RPS:         10,
			Concurrency: 1,
			Duration:    10,
			Output:      "html",
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
