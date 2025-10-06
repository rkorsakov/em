package config

import (
	"fmt"
	"log"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Server   ServerConfig   `yaml:"server" env-prefix:"SERVER_"`
	Database DatabaseConfig `yaml:"database" env-prefix:"DB_"`
}

type ServerConfig struct {
	Port string `yaml:"port" env:"PORT" env-default:"8080"`
	Host string `yaml:"host" env:"HOST" env-default:"localhost"`
}

type DatabaseConfig struct {
	Host     string `yaml:"host" env:"HOST" env-default:"localhost"`
	Port     string `yaml:"port" env:"PORT" env-default:"5432"`
	User     string `yaml:"user" env:"USER" env-default:"postgres"`
	Password string `yaml:"password" env:"PASSWORD" env-default:"password"`
	Name     string `yaml:"name" env:"NAME" env-default:"subscriptions"`
	SSLMode  string `yaml:"sslmode" env:"SSLMODE" env-default:"disable"`
}

func Load() (*Config, error) {
	var cfg Config

	err := cleanenv.ReadConfig("config/config.yaml", &cfg)
	if err != nil {
		log.Printf("Warning: config.yaml not found, using environment variables and defaults: %v", err)
	}

	err = cleanenv.ReadEnv(&cfg)
	if err != nil {
		return nil, fmt.Errorf("error reading environment variables: %w", err)
	}

	log.Printf("Configuration loaded successfully: Server=%s:%s, DB=%s@%s:%s",
		cfg.Server.Host, cfg.Server.Port,
		cfg.Database.User, cfg.Database.Host, cfg.Database.Port)

	return &cfg, nil
}
