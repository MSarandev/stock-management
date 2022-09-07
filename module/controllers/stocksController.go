package controllers

import (
	"context"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"sync"

	val "github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"stocks-api/module/entities"
	"stocks-api/module/services"
	"stocks-api/support/db"
)

type OpType string

const (
	insert OpType = "INSERT"
	update OpType = "UPDATE"
)

type StockService interface {
	GetAll(ctx context.Context) ([]*entities.Stock, error)
	GetOne(ctx context.Context, stockId string) (*entities.Stock, error)
	InsertOne(ctx context.Context, stock *entities.Stock) error
	UpdateOne(ctx context.Context, stock *entities.Stock, stockId string) error
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

func (s *StockController) GetOne(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	res, errGet := s.service.GetOne(s.ctx, vars["id"])
	if errGet != nil {
		w.Write([]byte(errGet.Error()))
		return
	}

	json.NewEncoder(w).Encode(res)
}

func (s *StockController) InsertOne(w http.ResponseWriter, r *http.Request) {
	stock, errParse := reqToStock(r)
	if errParse != nil {
		w.Write([]byte(errParse.Error()))
		return
	}

	errValidation := validate(stock, insert)
	if errValidation != nil {
		w.Write([]byte(errValidation.Error()))
		return
	}

	var m sync.Mutex
	m.Lock()

	if err := s.service.InsertOne(s.ctx, stock); err != nil {
		w.Write([]byte(err.Error()))
		m.Unlock()
		return
	}

	m.Unlock()
}

func (s *StockController) UpdateOne(w http.ResponseWriter, r *http.Request) {
	stock, errParse := reqToStock(r)
	if errParse != nil {
		w.Write([]byte(errParse.Error()))
		return
	}

	vars := mux.Vars(r)

	errValidation := validate(stock, update)
	if errValidation != nil {
		w.Write([]byte(errValidation.Error()))
		return
	}

	var m sync.Mutex
	m.Lock()

	if err := s.service.UpdateOne(s.ctx, stock, vars["id"]); err != nil {
		w.Write([]byte(err.Error()))
		m.Unlock()
		return
	}

	m.Unlock()
}

func (s *StockController) DeleteOne(http.ResponseWriter, *http.Request) {
	return
}

func validate(input *entities.Stock, op OpType) error {
	vl := val.New()

	if op == insert {
		return vl.Struct(InsertStock{
			ID:   input.ID.String(),
			Name: input.Name,
		})
	}

	if op == update {
		return vl.Struct(UpdateStock{
			Name:     input.Name,
			Quantity: int(input.Quantity),
		})
	}

	return errors.New("unexpected operation type received")
}

func reqToStock(r *http.Request) (*entities.Stock, error) {
	reqBody, errRead := ioutil.ReadAll(r.Body)
	if errRead != nil {
		return nil, errRead
	}

	stock := entities.Stock{}
	if err := json.Unmarshal(reqBody, &stock); err != nil {
		return nil, err
	}

	return &stock, nil
}
