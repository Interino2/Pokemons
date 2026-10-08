package config

import (
	"errors"
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	AppPort    string
}

func Load() (Config, error) {
	return load(".env")
}

func load(filename string) (Config, error) {
	if err := godotenv.Load(filename); err != nil && !errors.Is(err, os.ErrNotExist) {
		return Config{}, fmt.Errorf("cargando configuración desde .env: %w", err)
	}

	return Config{
		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnv("DB_PORT", "5432"),
		DBUser:     getEnv("DB_USER", "postgres"),
		DBPassword: os.Getenv("DB_PASSWORD"),
		DBName:     getEnv("DB_NAME", "pokemon_db"),
		AppPort:    getEnv("APP_PORT", "8080"),
	}, nil
}

func getEnv(key string, defaultValue string) string {
	value := os.Getenv(key)

	if value == "" {
		return defaultValue
	}

	return value
}
