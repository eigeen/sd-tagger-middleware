package config

import (
	"fmt"
	"github.com/spf13/viper"
)

func LoadConfig() {
	viper.SetConfigName("config")
	viper.SetConfigType("toml")
	viper.AddConfigPath("configs")
	err := viper.ReadInConfig()
	if err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			panic("config file not found: require configs/config.toml")
		} else {
			panic(fmt.Errorf("fatal error config file: %w", err))
		}
	}

	viper.SetDefault("override.model", "")
	viper.SetDefault("override.threshold", float64(0))
}
