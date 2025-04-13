package config

import (
	"log"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

// DB_HOST
// DB_PORT
// DB_USER
// DB_PASSWORD
// DB_NAME
// PORT

type Config struct {
	Port       string `env:"PORT"`
	DBUrl      string `env:"DB_URL"`
	DBHost     string `env:"DB_HOST"`
	DBPort     string `env:"DB_PORT"`
	DBUser     string `env:"DB_USER"`
	DBPassword string `env:"DB_PASSWORD"`
	DBName     string `env:"DB_NAME"`
	DBTimeout  string `env:"DB_TIMEOUT" env-default:"5"`
}

func LoadConfig() *Config {
	configPath := os.Getenv("CONFIG_PATH")
	var cfg Config

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("config file not found: %s", err)
	}

	err := cleanenv.ReadConfig(configPath, &cfg)
	if os.IsNotExist(err) {
		log.Fatalf("error reading config file: %s", err)
	}

	return &cfg
}
