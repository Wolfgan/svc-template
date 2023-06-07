package token

import (
	"context"
	"encoding/json"
	"time"

	"github.com/redis/go-redis/v9"
)

type storage[k string, v *Subject] struct {
	redis      *redis.Client
	expiration time.Duration
}

func NewRedisStorage[k string, v *Subject](redis *redis.Client, expiration time.Duration) Storage[k, v] {
	return &storage[k, v]{
		redis:      redis,
		expiration: expiration,
	}
}

func (s *storage[k, v]) Set(key k, value v) error {
	return s.redis.Set(context.Background(), string(key), value, s.expiration).Err()
}

func (s *storage[k, v]) Get(key k) (v, error) {
	buf, err := s.redis.Get(context.Background(), string(key)).Bytes()
	if err != nil {
		return nil, err
	}

	var value v
	if err := json.Unmarshal(buf, value); err != nil {
		return nil, err
	}

	return value, nil
}

func (s *storage[k, v]) Del(key k) error {
	return s.redis.Del(context.Background(), string(key)).Err()
}
