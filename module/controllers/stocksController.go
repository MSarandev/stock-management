package controllers

import (
	"net/http"

	"github.com/sirupsen/logrus"
)

type StockController struct {
	logger *logrus.Logger
}

func NewStockController(l *logrus.Logger) *StockController {
	return &StockController{logger: l}
}

func (s *StockController) GetAll(http.ResponseWriter, *http.Request) {
	return
}

func (s *StockController) GetOne(http.ResponseWriter, *http.Request) {
	return
}

func (s *StockController) UpdateOne(http.ResponseWriter, *http.Request) {
	return
}

func (s *StockController) DeleteOne(http.ResponseWriter, *http.Request) {
	return
}
