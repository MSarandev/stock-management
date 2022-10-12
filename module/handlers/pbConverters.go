package handlers

import (
	"google.golang.org/protobuf/types/known/timestamppb"
	"stocks-api/module/entities"
	pb "stocks-api/protos/protos/stocks"
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
	}
}

func toStockListPb(stocks []*entities.Stock) []*pb.SingleStock {
	response := make([]*pb.SingleStock, 0, len(stocks))

	for _, s := range stocks {
		response = append(response, toStockPb(s))
	}

	return response
}
