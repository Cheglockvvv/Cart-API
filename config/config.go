package config

import (
	"fmt"
	"github.com/caarlos0/env"
)

type APIConfig struct {
	Port       string `env:"API_PORT"`
	Migrations string `env:"API_MIGRATIONS_LOCATION"`
}

type DBConfig struct {
	User     string `env:"DB_USER"`
	Password string `env:"DB_PASSWORD"`
	Host     string `env:"DB_HOST"`
	Port     string `env:"DB_PORT"`
	DBName   string `env:"DB_NAME"`
	SSLMode  string `env:"DB_SSL_MODE"`
}

type Cart struct {
	API APIConfig
	DB  DBConfig
}

func LoadEnv() (*Cart, error) {
	var dbConfig DBConfig
	err := env.Parse(&dbConfig)
	if err != nil {
		return nil, fmt.Errorf("env.Parse: %w", err)
	}

	var apiConfig APIConfig
	err = env.Parse(&apiConfig)
	if err != nil {
		return nil, fmt.Errorf("env.Parse: %w", err)
	}

	config := Cart{
		API: apiConfig,
		DB:  dbConfig,
	}

	return &config, nil
}
