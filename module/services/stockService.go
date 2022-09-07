package services

import (
	"context"

	"github.com/sirupsen/logrus"
	"stocks-api/module/entities"
	"stocks-api/module/repos"
	"stocks-api/support/db"
)

type StockStore interface {
	GetAll(ctx context.Context) ([]*entities.Stock, error)
}

type StockService struct {
	repo    StockStore
	logger  *logrus.Logger
	db      *db.Instance
	Context context.Context
}

func NewStockService(l *logrus.Logger, db *db.Instance, ctx context.Context) *StockService {
	return &StockService{
		repo:    repos.NewStockRepo(l, db),
		db:      db,
		Context: ctx,
	}
}

func (s *StockService) GetAll(ctx context.Context) ([]*entities.Stock, error) {
	return s.repo.GetAll(ctx)
}
