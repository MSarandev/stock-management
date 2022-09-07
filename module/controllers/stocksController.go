package controllers

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/sirupsen/logrus"
	"stocks-api/module/entities"
	"stocks-api/module/services"
	"stocks-api/support/db"
)

type StockService interface {
	GetAll(ctx context.Context) ([]*entities.Stock, error)
}

type StockController struct {
	logger  *logrus.Logger
	db      *db.Instance
	service StockService
	ctx     context.Context
}

func NewStockController(l *logrus.Logger, db *db.Instance, ctx context.Context) *StockController {
	return &StockController{
		logger:  l,
		db:      db,
		service: services.NewStockService(l, db, ctx),
		ctx:     ctx,
	}
}

func (s *StockController) GetAll(w http.ResponseWriter, _ *http.Request) {
	res, errGet := s.service.GetAll(s.ctx)
	if errGet != nil {
		json.NewEncoder(w).Encode(errGet)
	}

	json.NewEncoder(w).Encode(res)
}

func (s *StockController) GetOne(http.ResponseWriter, *http.Request) {
	return
}

func (s *StockController) InsertOne(http.ResponseWriter, *http.Request) {

}

func (s *StockController) UpdateOne(http.ResponseWriter, *http.Request) {
	return
}

func (s *StockController) DeleteOne(http.ResponseWriter, *http.Request) {
	return
}
