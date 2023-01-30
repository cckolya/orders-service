package config

import (
	"encoding/json"
	"github.com/go-playground/validator/v10"
	"os"
	"time"
)

const configPath = "./config/config.json"

type Config struct {
	Postgres Postgres `validate:"required"`
	Handler  Handler  `validate:"required"`
}

type Postgres struct {
	Host     string `validate:"required"`
	Port     string `validate:"required"`
	User     string `validate:"required"`
	Password string `validate:"required"`
	DBName   string `validate:"required"`
	SSLMode  string `validate:"required"`
	Settings struct {
		MaxOpenConns    int           `validate:"required"`
		ConnMaxLifeTime time.Duration `validate:"required"`
		MaxIdleConns    int           `validate:"required"`
		MaxIdleLifeTime time.Duration `validate:"required"`
	} `validate:"required"`
}

type Handler struct {
	Url string `validate:"required"`
}

func LoadConfig() (*Config, error) {
	file, err := os.Open(configPath)
	if err != nil {
		return nil, err
	}

	var cfg *Config
	err = json.NewDecoder(file).Decode(&cfg)
	if err != nil {
		return nil, err
	}

	err = validator.New().Struct(cfg)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}
