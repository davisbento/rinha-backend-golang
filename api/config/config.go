package config

import (
	"fmt"
	"os"

	"github.com/spf13/viper"
)

type Config struct {
	DB_HOSTNAME string
	DB_PORT     int
	DB_USER     string
	DB_PASSWORD string
	DB_NAME     string
	REDIS_URL   string
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

		// if dont find the file, try to get the env variables
		config = Config{
			DB_HOSTNAME: os.Getenv("DB_HOSTNAME"),
			DB_PORT:     5432,
			DB_USER:     os.Getenv("DB_USER"),
			DB_PASSWORD: os.Getenv("DB_PASSWORD"),
			DB_NAME:     os.Getenv("DB_NAME"),
			REDIS_URL:   os.Getenv("REDIS_URL"),
		}
	}

	// Unmarshal the config
	if err := viper.Unmarshal(&config); err != nil {
		fmt.Printf("Error unmarshalling config: %s \n", err)
	}

	return &config
}
