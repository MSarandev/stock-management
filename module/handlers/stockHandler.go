package handlers

import (
	"context"
	"errors"

	val "github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
	"stocks-api/module/entities"
	"stocks-api/module/services"
	"stocks-api/module/validators"
	pb "stocks-api/protos/protos/stocks"
	"stocks-api/support/db"
)

// StockService an interface to the service.
type StockService interface {
	GetAll(ctx context.Context) ([]*entities.Stock, error)
	GetOne(ctx context.Context, stockId string) (*entities.Stock, error)
	InsertOne(ctx context.Context, stock *entities.Stock) error
	UpdateOne(ctx context.Context, stock *entities.Stock, stockId string) error
	DeleteOne(ctx context.Context, stockId string) error
}

// StockHandler handles all gRPC stock requests.
type StockHandler struct {
	logger  *logrus.Logger
	service StockService
	*pb.UnimplementedStockServiceServer
}

// NewStockHandler is a constructor for a new Stock Handler.
func NewStockHandler(l *logrus.Logger, db *db.Instance, ctx context.Context) *StockHandler {
	return &StockHandler{
		logger:                          l,
		service:                         services.NewStockService(l, db, ctx),
		UnimplementedStockServiceServer: &pb.UnimplementedStockServiceServer{},
	}
}

// GetStock returns a single stock item, fetched by ID.
func (s *StockHandler) GetStock(ctx context.Context, request *pb.GetStockRequest) (*pb.GetStockResponse, error) {
	if err := validateGet(request); err != nil {
		return nil, err
	}

	stock, err := s.service.GetOne(ctx, request.GetId())
	if err != nil {
		return nil, errors.New("Failed to fetch stock item")
	}

	return &pb.GetStockResponse{
		Stock: toStockPb(stock),
	}, nil
}

// ListStocks lists all stocks available in the db.
func (s *StockHandler) ListStocks(ctx context.Context, _ *pb.ListStocksRequest) (*pb.ListStocksResponse, error) {
	stocks, err := s.service.GetAll(ctx)
	if err != nil {
		return nil, errors.New("Failed to list stocks")
	}

	return &pb.ListStocksResponse{
		Stocks: toStockListPb(stocks),
	}, nil
}

func validateGet(r *pb.GetStockRequest) error {
	return val.New().Struct(validators.GetStock{ID: r.GetId()})
}

//func (s *StockHandler) CreateStock(ctx context.Context, request *pb.CreateStockRequest) (*pb.CreateStockResponse, error) {
//	panic("implement me")
//}
//
//func (s *StockHandler) EditStock(ctx context.Context, request *pb.EditStockRequest) (*pb.EditStockResponse, error) {
//	panic("implement me")
//}
//
//func (s *StockHandler) DeleteStock(ctx context.Context, request *pb.DeleteStockRequest) (*pb.DeleteStockResponse, error) {
//	panic("implement me")
//}
