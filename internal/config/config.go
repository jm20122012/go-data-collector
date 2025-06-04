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
	DBPass string `mapstructure:"DB_PASSWORD"`
	DBHost string `mapstructure:"DB_HOST"`
	DBPort string `mapstructure:"DB_PORT"`

	// Avtech Device Configuration
	AvtechURL string `mapstructure:"AVTECH_URL"`

	// MQTT Configuration
	MQTTBroker string `mapstructure:"MQTT_BROKER"`
	MQTTPort   int    `mapstructure:"MQTT_PORT"`

	// Ambient Weather Configuration
	AmbientAPIKey  string `mapstructure:"AMBIENT_API_KEY"`
	AmbientAppKey  string `mapstructure:"AMBIENT_APP_KEY"`
	AmbientURLFull string `mapstructure:"AMBIENT_URL_FULL"`

	// Feature Toggles
	EnableAvtechCollector  bool `mapstructure:"ENABLE_AVTECH_COLLECTOR"`
	EnableMQTTCollector    bool `mapstructure:"ENABLE_MQTT_COLLECTOR"`
	EnableAmbientCollector bool `mapstructure:"ENABLE_AMBIENT_COLLECTOR"`
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
