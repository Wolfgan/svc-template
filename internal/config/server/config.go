package server

import (
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

type Config struct {
	Port string `json:"port" yaml:"port" env:"X_SRV_PORT"`
	Auth Auth   `json:"auth" yaml:"auth"`
}

func (cfg Config) Validate() error {
	return validation.ValidateStruct(&cfg,
		validation.Field(&cfg.Port, validation.Required, is.Port),
		validation.Field(&cfg.Auth),
	)
}

type Auth struct {
	TokenSecret   string        `json:"token_secret" yaml:"token_secret" env:"X_TOKEN_SECRET"`
	AccessExpire  time.Duration `json:"access_expire" yaml:"access_expire" env:"X_ACCESS_EXPIRE"`
	RefreshExpire time.Duration `json:"refresh_expire" yaml:"refresh_expire" env:"X_REFRESH_EXPIRE"`
}

func (auth Auth) Validate() error {
	return validation.ValidateStruct(&auth,
		validation.Field(&auth.TokenSecret, validation.Required),
		validation.Field(&auth.AccessExpire, validation.Required),
		validation.Field(&auth.RefreshExpire, validation.Required),
	)
}
