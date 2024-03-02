package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	DB_HOSTNAME string
	DB_PORT     int
	DB_USER     string
	DB_PASSWORD string
	DB_NAME     string
}

func NewConfig() *Config {
	var config Config

	// Set up Viper
	viper.SetConfigFile(".env")
	viper.SetConfigType("env")
	viper.AutomaticEnv()

	// Read in the config file
	if err := viper.ReadInConfig(); err != nil {
		fmt.Printf("Error reading config file: %s \n", err)
	}

	// Unmarshal the config
	if err := viper.Unmarshal(&config); err != nil {
		fmt.Printf("Error unmarshalling config: %s \n", err)
	}

	return &config
}
