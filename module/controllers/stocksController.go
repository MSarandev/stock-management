package controllers

import (
	"context"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"

	val "github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"stocks-api/module/entities"
	"stocks-api/module/entities/filters"
	"stocks-api/module/services"
	"stocks-api/module/validators"
	"stocks-api/support/db"
)

// OpType defines the operation type.
type OpType string

const (
	insert OpType = "INSERT"
	update OpType = "UPDATE"
)

// A contract to the StockService for high level logic operations.
type StockService interface {
	GetAll(ctx context.Context, filters *filters.Pagination) ([]*entities.Stock, error)
	GetOne(ctx context.Context, stockId string) (*entities.Stock, error)
	InsertOne(ctx context.Context, stock *entities.Stock) error
	UpdateOne(ctx context.Context, stock *entities.Stock, stockId string) error
	DeleteOne(ctx context.Context, stockId string) error
	Count(ctx context.Context) (int, error)
}

// StockController handles the business logic when an endpoint is hit.
type StockController struct {
	logger  *logrus.Logger
	db      *db.Instance
	service StockService
	ctx     context.Context
}

// NewStockController a constructor for the StockController.
func NewStockController(l *logrus.Logger, db *db.Instance, ctx context.Context) *StockController {
	return &StockController{
		logger:  l,
		db:      db,
		service: services.NewStockService(l, db, ctx),
		ctx:     ctx,
	}
}

// GetAll returns all records in the database.
func (s *StockController) GetAll(w http.ResponseWriter, req *http.Request) {
	pagination, errParse := parsePagination(req)
	if errParse != nil {
		s.logger.Errorf("Failed to parse pagination: %s", errParse.Error())

		json.NewEncoder(w).Encode(errors.New("Failed to parse pagination"))
	}

	res, errGet := s.service.GetAll(s.ctx, pagination)
	if errGet != nil {
		json.NewEncoder(w).Encode(errGet)
	}

	count, errCount := s.service.Count(s.ctx)
	if errCount != nil {
		s.logger.Error(errCount)
		json.NewEncoder(w).Encode(errCount)
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"stocks":      res,
		"total_count": count,
	})
}

// GetOne returns a single record in the database.
func (s *StockController) GetOne(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	res, errGet := s.service.GetOne(s.ctx, vars["id"])
	if errGet != nil {
		w.Write([]byte(errGet.Error()))
		return
	}

	json.NewEncoder(w).Encode(res)
}

// InsertOne adds a new record to the database.
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

	if err := s.service.InsertOne(s.ctx, stock); err != nil {
		w.Write([]byte(err.Error()))
		return
	}
}

// UpdateOne updates a single record in the database.
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

	if err := s.service.UpdateOne(s.ctx, stock, vars["id"]); err != nil {
		w.Write([]byte(err.Error()))
		return
	}
}

// DeleteOne deletes a single record in the database.
func (s *StockController) DeleteOne(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	if err := s.service.DeleteOne(s.ctx, vars["id"]); err != nil {
		w.Write([]byte(err.Error()))
		return
	}
}

func validate(input *entities.Stock, op OpType) error {
	vl := val.New()

	if op == insert {
		return vl.Struct(validators.InsertStock{
			ID:   input.ID.String(),
			Name: input.Name,
		})
	}

	if op == update {
		return vl.Struct(validators.UpdateStock{
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

func parsePagination(r *http.Request) (*filters.Pagination, error) {
	reqBody, errRead := ioutil.ReadAll(r.Body)
	if errRead != nil {
		return nil, errRead
	}

	var unm map[string]*filters.Pagination

	if err := json.Unmarshal(reqBody, &unm); err != nil {
		return nil, err
	}

	return unm["pagination"], nil
}
