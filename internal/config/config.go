package config

import (
	"github.com/spf13/viper"
)

// AppConfig holds the application configuration
type AppConfig struct {
	// LogLevel is the logging level for the application
	LogLevel string `mapstructure:"LOG_LEVEL"`

	// Database Configuration
	DBName string `mapstructure:"DB_NAME"`
	DBUser string `mapstructure:"DB_USER"`
	DBPass string `mapstructure:"DB_PASS"`
	DBHost string `mapstructure:"DB_HOST"`
	DBPort string `mapstructure:"DB_PORT"`
}

// AppConfigInstance is the instance of the application configuration
func GetConfig() *AppConfig {
	// Load .env file
	viper.SetConfigName(".env")
	viper.SetConfigType("env")
	viper.AddConfigPath(".")
	if err := viper.ReadInConfig(); err != nil {
		panic("Error reading .env file: " + err.Error())
	}
	var config AppConfig
	if err := viper.Unmarshal(&config); err != nil {
		panic("Error unmarshalling config: " + err.Error())
	}
	return &config
}
