package repos

import (
	"context"

	"github.com/sirupsen/logrus"
	"stocks-api/module/entities"
	"stocks-api/support/db"
)

type StockRepo struct {
	logger *logrus.Logger
	db     *db.Instance
}

func NewStockRepo(l *logrus.Logger, db *db.Instance) *StockRepo {
	return &StockRepo{
		logger: l,
		db:     db,
	}
}

func (s *StockRepo) GetAll(ctx context.Context) ([]*entities.Stock, error) {
	var x []*entities.Stock

	s.db.Base.NewSelect().
		Model(new(entities.Stock)).
		Scan(ctx, &x)

	return x, nil
}
