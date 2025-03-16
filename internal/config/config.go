package config

import (
	"fmt"
	"os"
	"strconv"
)

type Config struct {
	Log      LogConfig
	Bot      BotConfig
	Postgres PostgresConfig
}

type LogConfig struct {
	Level  string
	Stdout bool
}

type BotConfig struct {
	Token string
}

type PostgresConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
}

func LoadConfig() (*Config, error) {
	config := &Config{}

	// Логирование
	config.Log.Level = os.Getenv("LOG_LEVEL")
	if config.Log.Level == "" {
		return nil, fmt.Errorf("LOG_LEVEL is required")
	}

	stdoutStr := os.Getenv("LOG_STDOUT")
	if stdoutStr == "" {
		return nil, fmt.Errorf("LOG_STDOUT is required")
	}
	stdout, err := strconv.ParseBool(stdoutStr)
	if err != nil {
		return nil, fmt.Errorf("invalid LOG_STDOUT value: %v", err)
	}
	config.Log.Stdout = stdout

	// Бот
	config.Bot.Token = os.Getenv("BOT_TOKEN")
	if config.Bot.Token == "" {
		return nil, fmt.Errorf("BOT_TOKEN is required")
	}

	// PostgreSQL
	config.Postgres.Host = os.Getenv("DB_HOST")
	if config.Postgres.Host == "" {
		return nil, fmt.Errorf("DB_HOST is required")
	}

	portStr := os.Getenv("DB_PORT")
	if portStr == "" {
		return nil, fmt.Errorf("DB_PORT is required")
	}
	port, err := strconv.Atoi(portStr)
	if err != nil {
		return nil, fmt.Errorf("invalid DB_PORT value: %v", err)
	}
	config.Postgres.Port = port

	config.Postgres.User = os.Getenv("DB_USER")
	if config.Postgres.User == "" {
		return nil, fmt.Errorf("DB_USER is required")
	}

	config.Postgres.Password = os.Getenv("DB_PASSWORD")
	if config.Postgres.Password == "" {
		return nil, fmt.Errorf("DB_PASSWORD is required")
	}

	config.Postgres.DBName = os.Getenv("DB_NAME")
	if config.Postgres.DBName == "" {
		return nil, fmt.Errorf("DB_NAME is required")
	}

	return config, nil
}
