package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

func Init() {
	configDir, err := os.UserConfigDir()
	if err != nil {
		configDir = os.Getenv("HOME")
	}

	improvmxDir := filepath.Join(configDir, "improvmx")
	os.MkdirAll(improvmxDir, 0700)

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(improvmxDir)

	viper.SetEnvPrefix("IMPROVMX")
	viper.AutomaticEnv()

	viper.ReadInConfig()
}

func GetAPIKey() string {
	key := viper.GetString("api_key")
	if key == "" {
		fmt.Fprintln(os.Stderr, "Error: API key not configured. Run 'improvmx auth login' or set IMPROVMX_API_KEY.")
		os.Exit(1)
	}
	return key
}

func SaveAPIKey(key string) error {
	configDir, err := os.UserConfigDir()
	if err != nil {
		configDir = os.Getenv("HOME")
	}

	improvmxDir := filepath.Join(configDir, "improvmx")
	os.MkdirAll(improvmxDir, 0700)

	viper.Set("api_key", key)
	configPath := filepath.Join(improvmxDir, "config.yaml")
	return viper.WriteConfigAs(configPath)
}
