package users

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"service-template/internal/model"

	"github.com/uptrace/bun"
)

var (
	ErrNotExists = fmt.Errorf("user not exists")
)

type Storage struct {
	db *bun.DB
}

func NewStorage(db *bun.DB) *Storage {
	return &Storage{
		db: db,
	}
}

func (s *Storage) Create(ctx context.Context, user *model.User) (*model.User, error) {
	if _, err := s.db.NewInsert().Model(user).Returning("*").Exec(ctx); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *Storage) Update(user *model.User) error {
	return nil
}

func (s *Storage) Get(ctx context.Context, user *model.User) (*model.User, error) {
	query := s.db.NewSelect().Model(user)

	if user.Email != "" {
		query = query.Where("email = ?", user.Email)
	}

	if user.Phone != "" {
		query = query.Where("phone = ?", user.Phone)
	}

	if err := query.Scan(ctx); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotExists
		}

		return nil, err
	}

	return user, nil
}
func (s *Storage) Exists(ctx context.Context, user *model.User) (bool, error) {
	query := s.db.NewSelect().Model(user)

	if user.Email != "" {
		query = query.Where("email = ?", user.Email)
	}

	if user.Phone != "" {
		query = query.Where("phone = ?", user.Phone)
	}

	return query.Exists(ctx)
}

func (s *Storage) GetByID(id int64) (*model.User, error) {
	return nil, nil
}

func (s *Storage) Delete(id int64) error {
	return nil
}

func (s *Storage) List() ([]*model.User, error) {
	return nil, nil
}
