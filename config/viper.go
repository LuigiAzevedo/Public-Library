package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type AppConfig struct {
	DbURL        string `mapstructure:"DB_URL"`
	DbDriver     string `mapstructure:"DB_DRIVER"`
	ServeAddress string `mapstructure:"SERVE_ADDRESS"`
}

func LoadAppConfig(path string) (AppConfig, error) {
	if path == "" {
		return AppConfig{}, fmt.Errorf("config path is empty")
	}

	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			return AppConfig{}, fmt.Errorf("config file not found: %s", path)
		}
		return AppConfig{}, fmt.Errorf("failed to read config file: %w", err)
	}

	var config AppConfig
	if err := viper.Unmarshal(&config); err != nil {
		return AppConfig{}, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return config, nil
}
