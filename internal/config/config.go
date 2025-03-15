package config

import (
	"encoding/json"
	"os"
)

type Config struct {
	Log LogConfig `json:"log"`
	Bot BotConfig `json:"bot"`
}

type LogConfig struct {
	Level  string `json:"level"`
	Stdout bool   `json:"stdout"`
}

type BotConfig struct {
	Token string `json:"token"`
}

func LoadConfig() (*Config, error) {
	bytes, err := os.ReadFile("/app/config/config.json")
	if err != nil {
		return nil, err
	}

	var config Config
	err = json.Unmarshal(bytes, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}
