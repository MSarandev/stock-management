package handlers

import (
	"time"

	"google.golang.org/protobuf/types/known/timestamppb"
	pb "stocks-api/genprotos"
	"stocks-api/module/entities"
)

func toStockPb(stock *entities.Stock) *pb.SingleStock {
	return &pb.SingleStock{
		Id:       stock.ID.String(),
		Name:     stock.Name,
		Quantity: stock.Quantity,
		CreatedAt: &timestamppb.Timestamp{
			Seconds: stock.CreatedAt.Unix(),
			Nanos:   int32(stock.CreatedAt.Nanosecond()),
		},
		UpdatedAt: &timestamppb.Timestamp{
			Seconds: stock.UpdatedAt.Unix(),
			Nanos:   int32(stock.UpdatedAt.Nanosecond()),
		},
		HCreatedAt: stock.CreatedAt.Format(time.RFC3339),
		HUpdatedAt: stock.UpdatedAt.Format(time.RFC3339),
	}
}

func toStockListPb(stocks []*entities.Stock) []*pb.SingleStock {
	response := make([]*pb.SingleStock, 0, len(stocks))

	for _, s := range stocks {
		response = append(response, toStockPb(s))
	}

	return response
}

func fromCreatePb(req *pb.CreateStockRequest) *entities.Stock {
	return &entities.Stock{
		Name:     req.GetStock().GetName(),
		Quantity: req.GetStock().GetQuantity(),
	}
}

func fromEditPb(req *pb.EditStockRequest) *entities.Stock {
	return &entities.Stock{
		Name:     req.GetStock().GetName(),
		Quantity: req.GetStock().GetQuantity(),
	}
}
