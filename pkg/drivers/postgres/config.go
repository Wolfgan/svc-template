package postgres

import (
	"service-template/internal/config/valid"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

type Config struct {
	Addr string `json:"addr" yaml:"addr" env:"X_DB_ADDR"`
	Name string `json:"name" yaml:"name" env:"X_DB_NAME"`
	User string `json:"user" yaml:"user" env:"X_DB_USER"`
	Pass string `json:"pass" yaml:"pass" env:"X_DB_PASS"`
}

func (cfg *Config) Validate() error {
	return validation.ValidateStruct(cfg,
		validation.Field(&cfg.Addr, validation.Required, is.URL),
		validation.Field(&cfg.Name, validation.Required),
		validation.Field(&cfg.User, validation.Required, validation.Match(valid.Name)),
		validation.Field(&cfg.Pass, validation.Required, validation.Match(valid.Password)),
	)
}
