package redisdb

import (
	"service-template/internal/config/valid"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

type Config struct {
	Addr string `json:"addr" yaml:"addr" env:"X_REDIS_ADDR"`
	DB   int    `json:"db" yaml:"db" env:"X_REDIS_DB"`
	User string `json:"user" yaml:"user" env:"X_REDIS_USER"`
	Pass string `json:"pass" yaml:"pass" env:"X_REDIS_PASS"`
}

func (cfg *Config) Validate() error {
	return validation.ValidateStruct(cfg,
		validation.Field(&cfg.Addr, validation.Required, is.URL),
		validation.Field(&cfg.DB, validation.Min(0), validation.Max(15)),
		validation.Field(&cfg.User, validation.When(cfg.User != "", validation.Match(valid.Name))),
		validation.Field(&cfg.Pass, validation.When(cfg.Pass != "", validation.Match(valid.Password))),
	)
}
