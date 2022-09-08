package services

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"stocks-api/module/entities"
	"stocks-api/module/repos"
	"stocks-api/support/db"
)

// StockStore a contract to the Stock Repo.
type StockStore interface {
	GetAll(ctx context.Context) ([]*entities.Stock, error)
	GetOne(ctx context.Context, id uuid.UUID) (*entities.Stock, error)
	InsertOne(ctx context.Context, stock *entities.Stock) error
	UpdateOne(ctx context.Context, stock *entities.Stock) error
	DeleteOne(ctx context.Context, id uuid.UUID) error
}

// StockService provides high level logic.
type StockService struct {
	repo    StockStore
	logger  *logrus.Logger
	db      *db.Instance
	Context context.Context
}

// NewStockService a constructor for the Stock Service.
func NewStockService(l *logrus.Logger, db *db.Instance, ctx context.Context) *StockService {
	return &StockService{
		repo:    repos.NewStockRepo(l, db),
		db:      db,
		Context: ctx,
	}
}

// GetAll returns all records in the db.
func (s *StockService) GetAll(ctx context.Context) ([]*entities.Stock, error) {
	return s.repo.GetAll(ctx)
}

// GetOne returns a single record in the db.
func (s *StockService) GetOne(ctx context.Context, stockId string) (*entities.Stock, error) {
	id, errParse := uuid.Parse(stockId)
	if errParse != nil {
		return nil, errParse
	}

	return s.repo.GetOne(ctx, id)
}

// InsertOne adds a new record in the db.
func (s *StockService) InsertOne(ctx context.Context, stock *entities.Stock) error {
	if err := checkQuantity(stock); err != nil {
		return err
	}

	return s.repo.InsertOne(ctx, stock)
}

// UpdateOne updates a single record in the db.
func (s *StockService) UpdateOne(ctx context.Context, stock *entities.Stock, stockId string) error {
	id, errParse := uuid.Parse(stockId)
	if errParse != nil {
		return errParse
	}

	if err := checkQuantity(stock); err != nil {
		return err
	}

	stock.ID = id

	return s.repo.UpdateOne(ctx, stock)
}

// DeleteOne removes a record from the db.
func (s *StockService) DeleteOne(ctx context.Context, stockId string) error {
	id, errParse := uuid.Parse(stockId)
	if errParse != nil {
		return errParse
	}

	return s.repo.DeleteOne(ctx, id)
}

// checkQuantity a custom quantity check, due to the unique way Go handles zero values.
func checkQuantity(s *entities.Stock) error {
	if s.Quantity < 0 {
		return errors.New("Input quantity is less than 0")
	}

	return nil
}
