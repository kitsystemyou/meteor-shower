package config

import (
	"fmt"
	"os"

	"github.com/spf13/viper"
)

type Config struct {
	App AppConfig `mapstructure:"app"`
}

type AppConfig struct {
	Name    string `mapstructure:"name"`
	Message string `mapstructure:"message"`
	Timeout int    `mapstructure:"timeout"`
	Debug   bool   `mapstructure:"debug"`
}

func LoadConfig(cfgFile string) (*Config, error) {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		viper.SetConfigName("config")
		viper.SetConfigType("yaml")
		viper.AddConfigPath(".")
		viper.AddConfigPath("$HOME/.mycli")
	}

	viper.SetDefault("app.name", "World")
	viper.SetDefault("app.message", "Hello")
	viper.SetDefault("app.timeout", 30)
	viper.SetDefault("app.debug", false)

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("error reading config file: %w", err)
		}
		fmt.Fprintf(os.Stderr, "Warning: No config file found, using defaults\n")
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("unable to decode config: %w", err)
	}

	return &cfg, nil
}
