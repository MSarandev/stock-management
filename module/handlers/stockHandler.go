package handlers

import (
	"context"
	"errors"

	val "github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
	pb "stocks-api/genprotos"
	"stocks-api/module/entities"
	"stocks-api/module/entities/filters"
	"stocks-api/module/services"
	"stocks-api/module/validators"
	"stocks-api/support/db"
)

// StockService an interface to the service.
type StockService interface {
	GetAll(ctx context.Context, pagination *filters.Pagination) ([]*entities.Stock, error)
	GetOne(ctx context.Context, stockId string) (*entities.Stock, error)
	InsertOne(ctx context.Context, stock *entities.Stock) error
	UpdateOne(ctx context.Context, stock *entities.Stock, stockId string) error
	DeleteOne(ctx context.Context, stockId string) error
	Count(ctx context.Context) (int, error)
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
func (s *StockHandler) ListStocks(ctx context.Context, req *pb.ListStocksRequest) (*pb.ListStocksResponse, error) {
	if req.GetPagination() == nil {
		return nil, errors.New("Pagination is required")
	}

	// TODO: validation
	pagination := &filters.Pagination{
		Page:         int(req.GetPagination().GetPage()),
		ItemsPerPage: int(req.GetPagination().GetItemsPerPage()),
	}

	stocks, err := s.service.GetAll(ctx, pagination)
	if err != nil {
		s.logger.Error(err)
		return nil, errors.New("Failed to list stocks")
	}

	count, err := s.service.Count(ctx)
	if err != nil {
		s.logger.Error(err)
		return nil, errors.New("Failed to get count")
	}

	return &pb.ListStocksResponse{
		Stocks:     toStockListPb(stocks),
		TotalCount: int64(count),
	}, nil
}

// CreateStock creates a new stock item
func (s *StockHandler) CreateStock(ctx context.Context, request *pb.CreateStockRequest) (*pb.CreateStockResponse, error) {
	if err := validateCreate(request); err != nil {
		return nil, errors.New("Failed to create stock")
	}

	stock := fromCreatePb(request)

	if err := s.service.InsertOne(ctx, stock); err != nil {
		return nil, err
	}

	return &pb.CreateStockResponse{}, nil
}

// EditStock modifies an existing stock item
func (s *StockHandler) EditStock(ctx context.Context, request *pb.EditStockRequest) (*pb.EditStockResponse, error) {
	if err := validateUpdate(request); err != nil {
		return nil, errors.New("Failed to update stock")
	}

	stock := fromEditPb(request)

	if err := s.service.UpdateOne(ctx, stock, request.GetStock().GetId()); err != nil {
		return nil, errors.New("Failed to update stock")
	}

	return &pb.EditStockResponse{}, nil
}

// DeleteStock removes a given stock item.
func (s *StockHandler) DeleteStock(ctx context.Context, request *pb.DeleteStockRequest) (*pb.DeleteStockResponse, error) {
	if request.GetId() == "" {
		return nil, errors.New("StockID is required")
	}

	if err := s.service.DeleteOne(ctx, request.GetId()); err != nil {
		return nil, err
	}

	return &pb.DeleteStockResponse{}, nil
}

func validateGet(r *pb.GetStockRequest) error {
	return val.New().Struct(validators.GetStock{ID: r.GetId()})
}

func validateCreate(r *pb.CreateStockRequest) error {
	return val.New().Struct(validators.InsertStock{
		Name:     r.GetStock().GetName(),
		Quantity: int(r.GetStock().GetQuantity()),
	})
}

func validateUpdate(r *pb.EditStockRequest) error {
	return val.New().Struct(validators.UpdateStock{
		Name:     r.GetStock().GetName(),
		Quantity: int(r.GetStock().GetQuantity()),
	})
}
