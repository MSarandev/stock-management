package repos

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

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

		s.logger.Error(err)

		return err
	})
}

func (s *StockRepo) UpdateOne(ctx context.Context, stock *entities.Stock) error {
	currentRecord := entities.Stock{}

	errExists := s.db.Base.NewSelect().
		Model(stock).
		Where("id = ?", stock.ID).
		Scan(ctx, &currentRecord)
	if errExists != nil {
		s.logger.Error(errExists)
		return errExists
	}

	return s.db.Base.RunInTx(ctx, &sql.TxOptions{}, func(ctx context.Context, tx bun.Tx) error {
		stock.CreatedAt = currentRecord.CreatedAt

		_, err := tx.NewUpdate().
			Model(stock).
			WherePK().
			Exec(ctx)

		s.logger.Error(err)

		return err
	})
}

func (s *StockRepo) DeleteOne(ctx context.Context, id uuid.UUID) error {
	exists, errExists := s.db.Base.NewSelect().
		Model(new(entities.Stock)).
		Where("id = ?", id).
		Exists(ctx)
	if errExists != nil {
		s.logger.Error(errExists)
		return errExists
	}

	if !exists {
		err := errors.New(fmt.Sprintf("record with id: %s doesn't exist", id))

		s.logger.Error(err)
		return err
	}

	return s.db.Base.RunInTx(ctx, &sql.TxOptions{}, func(ctx context.Context, tx bun.Tx) error {
		_, err := tx.NewDelete().
			Model(new(entities.Stock)).
			Where("id = ?", id).
			Exec(ctx)

		s.logger.Error(err)

		return err
	})
}
