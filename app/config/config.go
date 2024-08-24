package config

import (
	"fmt"
	"os"
)

type APIConfig struct {
	Port string
}

type DBConfig struct {
	User     string
	Password string
	Host     string
	Port     string
	DBName   string
	SSLMode  string
}

type Cart struct {
	API APIConfig
	DB  DBConfig
}

func LoadEnv() (*Cart, error) {

	requiredEnvs := []string{
		"DB_HOST", "DB_USER", "DB_PASSWORD", "DB_NAME", "DB_PORT", "API_PORT", "DB_SSLMODE"}

	for _, env := range requiredEnvs {
		if os.Getenv(env) == "" {
			return nil, fmt.Errorf("missing environment variable %s", env)
		}
	}

	config := Cart{
		API: APIConfig{
			Port: os.Getenv("API_PORT"),
		},
		DB: DBConfig{
			Host:     os.Getenv("DB_HOST"),
			User:     os.Getenv("DB_USER"),
			Password: os.Getenv("DB_PASSWORD"),
			DBName:   os.Getenv("DB_NAME"),
			Port:     os.Getenv("DB_PORT"),
			SSLMode:  os.Getenv("DB_SSLMODE"),
		},
	}

	return &config, nil
}
