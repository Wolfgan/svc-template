package services

import (
	"service-template/internal/config"
	"service-template/internal/daemon/services/auth"
	"service-template/internal/db"
)

type Interactor struct {
	Auth *auth.Service
}

func NewInteractor(cfg *config.Config, storage *db.Storage) *Interactor {
	return &Interactor{
		Auth: auth.NewService(cfg, storage),
	}
}
