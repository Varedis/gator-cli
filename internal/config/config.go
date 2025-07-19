// Package config implements utility functions for manipulating config files for the Gator CLI
package config

import (
	"encoding/json"
	"os"
	"path"
)

const CONFIG_FILE_NAME = ".gatorconfig.json"

type Config struct {
	DbURL           string `json:"db_url"`
	CurrentUserName string `json:"current_user_name"`
}

func (cfg *Config) SetUser(name string) error {
	cfg.CurrentUserName = name
	return write(*cfg)
}

func Read() (Config, error) {
	filePath, err := getConfigFilePath()
	if err != nil {
		return Config{}, err
	}

	data, err := os.ReadFile(filePath)
	if err != nil {
		return Config{}, err
	}

	var config Config
	if err := json.Unmarshal(data, &config); err != nil {
		return Config{}, err
	}
	return config, nil
}

func getConfigFilePath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return path.Join(homeDir, CONFIG_FILE_NAME), nil
}

func write(cfg Config) error {
	jsonData, err := json.Marshal(cfg)
	if err != nil {
		return err
	}

	filePath, err := getConfigFilePath()
	if err != nil {
		return err
	}

	err = os.WriteFile(filePath, jsonData, 0644)
	if err != nil {
		return err
	}

	return nil
}
