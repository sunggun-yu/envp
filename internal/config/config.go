package config

import (
	"os"
	"path/filepath"
)

const (
	configRoot     string = ".config/prw"
	configFileName string = "config.yaml"
)

type ProxyWrapperConfig struct {
	Default string `mapstructure:"default"`
}

type ProxyProfile struct {
	Description string `mapstructure:"desc"`
	Host        string `mapstructure:"host"`
	NoProxy     string `mapstructure:"noproxy"`
}

// TODO: working on it
func ConfigFile() (string, error) {
	// read config file
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	configPath := filepath.Join(home, configRoot)
	configFile := filepath.Join(configPath, configFileName)

	// create config file if not exists
	if _, err := os.Stat(configPath); err != nil {
		if os.IsNotExist(err) {
			os.MkdirAll(configPath, 0755)
		} else {
			return "", err
		}
	}
	return configFile, nil
}
