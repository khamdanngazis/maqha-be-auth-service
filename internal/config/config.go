// internal/config/config.go

package config

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/viper"
)

// DatabaseConfig holds the database configuration.
type DatabaseConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
	Debug    bool
}

// Config holds the application configuration.
type Config struct {
	Database DatabaseConfig
	AppPort  string
	GrpcPort string
}

// LoadConfig loads configuration from a specified file path, environment variables, and/or config files.
func LoadConfig(filePath string) (*Config, error) {
	if filePath != "" {
		viper.SetConfigFile(filePath)
	} else {
		viper.SetConfigName("config") // Config file name (without extension)
		viper.AddConfigPath(".")      // Look for the config file in the current directory
	}

	viper.SetConfigType("yaml") // Config file type (can be JSON, TOML, etc.)

	// Environment variables support (AUTH_DATABASE_HOST, AUTH_APPPORT, etc.)
	viper.SetEnvPrefix("AUTH")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("error reading config file: %v", err)
		}
	}

	config := &Config{}

	// Unmarshal the configuration into the Config struct
	if err := viper.Unmarshal(config); err != nil {
		return nil, fmt.Errorf("error unmarshalling config: %v", err)
	}

	return config, nil
}
