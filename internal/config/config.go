package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

type Config struct {
	DbURL           string `json:"db_url"`
	CurrentUserName string `json:"current_user_name"`
}

const configFileName = ".gatorconfig.json"

func Read() (Config, error) {
	configFilePath, err := getConfigFilePath()
	if err != nil {
		return Config{}, fmt.Errorf("Error getting config file path: %s", err)
	}

	data, err := os.ReadFile(configFilePath)
	if err != nil {
		return Config{}, fmt.Errorf("Error reading file: %s", err)
	}

	config := Config{}
	err = json.Unmarshal(data, &config)
	if err != nil {
		return Config{}, fmt.Errorf("Error decoding bytes to json: %s", err)
	}
	return config, nil
}

func getConfigFilePath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("Error getting user home directory: %s", err)
	}
	configFilePath := filepath.Join(homeDir, configFileName)
	return configFilePath, nil
}

func (config *Config) SetUser(user string) error {
	config.CurrentUserName = user
	return write(*config)
}

func write(config Config) error {
	data, err := json.Marshal(config)
	if err != nil {
		return fmt.Errorf("Error encoding config struct to json: %s", err)
	}
	configFilePath, err := getConfigFilePath()
	if err != nil {
		return fmt.Errorf("Error getting config file path: %s", err)
	}
	err = os.WriteFile(configFilePath, data, 0600)
	if err != nil {
		return fmt.Errorf("Error writing to config file: %s", err)
	}
	return nil
}
