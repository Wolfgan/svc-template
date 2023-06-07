package config

import (
	"path/filepath"
	"service-template/internal/config/logger"
	"service-template/internal/config/server"
	"service-template/pkg/drivers/postgres"
	"service-template/pkg/drivers/redisdb"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Server   *server.Config   `json:"server" yaml:"server"`
	Logger   *logger.Config   `json:"logger" yaml:"logger"`
	Redis    *redisdb.Config  `json:"redis" yaml:"redis"`
	Postgres *postgres.Config `json:"postgres" yaml:"postgres"`
}

// New создает новую конфигурацию и загружает значения из файла.
// Если заданы переменные окружения, тогда они будут иметь приоритет.
func New(filename string) (*Config, error) {
	cfg := Config{
		Postgres: &postgres.Config{},
	}

	path, err := filepath.Abs(filename)
	if err != nil {
		return nil, err
	}

	if err := cleanenv.ReadConfig(path, &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

// Validate проверяет заполнение обязательных полей.
func (cfg *Config) Validate() error {
	return validation.ValidateStruct(cfg,
		validation.Field(&cfg.Server),
		validation.Field(&cfg.Logger),
		validation.Field(&cfg.Redis),
		validation.Field(&cfg.Postgres),
	)
}
