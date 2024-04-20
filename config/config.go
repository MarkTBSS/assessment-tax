package config

import (
	"os"
	"sync"

	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"
)

type (
	Config struct {
		Database *Database `mapstructure:"database" validate:"required"`
		Server   *Server   `mapstructure:"server" validate:"required"`
	}

	Database struct {
		DatabaseURL string `mapstructure:"DATABASE_URL" validate:"required"`
	}

	Server struct {
		Port     string `mapstructure:"PORT" validate:"required"`
		Username string `mapstructure:"ADMIN_USERNAME" validate:"required"`
		Password string `mapstructure:"ADMIN_PASSWORD" validate:"required"`
	}
)

var (
	once           sync.Once
	configInstance *Config
)

func ConfigGetting() *Config {
	once.Do(func() {
		loadEnv()
		configInstance = &Config{
			Database: &Database{
				DatabaseURL: getEnv("DATABASE_URL", ""),
			},
			Server: &Server{
				Port:     getEnv("PORT", ""),
				Username: getEnv("ADMIN_USERNAME", ""),
				Password: getEnv("ADMIN_PASSWORD", ""),
			},
		}
		validate := validator.New()
		if err := validate.Struct(configInstance); err != nil {
			panic(err)
		}
	})
	return configInstance
}

func loadEnv() {
	err := godotenv.Load()
	if err != nil {
		panic("Error loading .env file")
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
