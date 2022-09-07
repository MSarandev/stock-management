package repos

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"github.com/uptrace/bun"
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

func (s *StockRepo) GetOne(ctx context.Context, id uuid.UUID) (*entities.Stock, error) {
	var x entities.Stock

	s.db.Base.NewSelect().
		Model(&x).
		Where("id = ?", id).
		Scan(ctx)

	return &x, nil
}

func (s *StockRepo) InsertOne(ctx context.Context, stock *entities.Stock) error {
	return s.db.Base.RunInTx(ctx, &sql.TxOptions{}, func(ctx context.Context, tx bun.Tx) error {
		_, err := tx.NewInsert().
			Model(stock).
			Exec(ctx)

		return err
	})
}
