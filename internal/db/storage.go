package db

import (
	"errors"
	"service-template/internal/config"
	"service-template/internal/db/token"
	"service-template/internal/db/users"
	"service-template/pkg/drivers/postgres"
	"service-template/pkg/drivers/redisdb"

	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/extra/bundebug"
	"github.com/uptrace/bun/extra/bunzerolog"
)

type Storage struct {
	cfg *config.Config
	log *zerolog.Logger
	pg  *bun.DB
	rdb *redis.Client

	Token token.Storage[string, *token.Subject]
	Users *users.Storage
}

func NewStorage(cfg *config.Config, log *zerolog.Logger) (*Storage, error) {
	var err error

	storage := Storage{
		cfg: cfg,
		log: log,
	}

	if storage.pg, err = postgres.NewPostgresDB(cfg.Postgres); err != nil {
		return nil, err
	}

	// Логирование SQL запросов
	if cfg.Logger.SQL.Enabled {
		if cfg.Logger.SQL.Default {
			// Стандартный логировщик
			storage.pg.AddQueryHook(bundebug.NewQueryHook(
				bundebug.WithVerbose(true),
			))
		} else {
			// Логировщик с помощью zerolog
			storage.pg.AddQueryHook(&bunzerolog.QueryHook{})
		}
	}

	if storage.rdb, err = redisdb.NewRedisDB(cfg.Redis); err != nil {
		return nil, err
	}

	storage.Token = token.NewRedisStorage[string, *token.Subject](storage.rdb, cfg.Server.Auth.AccessExpire)
	storage.Users = users.NewStorage(storage.pg)

	return &storage, nil
}

func (s *Storage) Close() error {
	var errs []error

	if s.pg != nil {
		if err := s.pg.Close(); err != nil {
			errs = append(errs, err)
		}
	}

	if s.rdb != nil {
		if err := s.rdb.Close(); err != nil {
			errs = append(errs, err)
		}
	}

	return errors.Join(errs...)
}
