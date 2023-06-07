package logger

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/rs/zerolog"
)

type Config struct {
	Level      zerolog.Level `json:"level" yaml:"level" env:"X_LOG_LEVEL"`
	Console    bool          `json:"console" yaml:"console" env:"X_LOG_CONSOLE"`
	WithCaller bool          `json:"with_caller" yaml:"with_caller" env:"X_LOG_WITH_CALLER"`
	SQL        struct {
		Enabled bool `json:"enabled" yaml:"enabled" env:"X_LOG_SQL_ENABLE"`
		Default bool `json:"default" yaml:"default" env:"X_LOG_SQL_DEFAULT"`
	} `json:"sql" yaml:"sql"`
}

func (cfg Config) Validate() error {
	return validation.ValidateStruct(&cfg,
		validation.Field(&cfg.Level, validation.Min(zerolog.TraceLevel), validation.Max(zerolog.Disabled)),
	)
}
