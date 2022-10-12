package grpc

import (
	"flag"
	"fmt"
	"log"
	"net"
	"sync"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	pb "stocks-api/genprotos"
	"stocks-api/module/handlers"
)

// Serve is the gRPC serve wrapper.
type Serve struct {
	port    int64
	logger  *logrus.Logger
	opts    *[]grpc.ServerOption
	handler *handlers.StockHandler
	wg      *sync.WaitGroup
}

// NewServe is a wrapper constructor.
func NewServe(
	p int64,
	l *logrus.Logger,
	opts *[]grpc.ServerOption,
	handler *handlers.StockHandler,
	wg *sync.WaitGroup,
) *Serve {
	return &Serve{
		port:    p,
		logger:  l,
		opts:    opts,
		handler: handler,
		wg:      wg,
	}
}

// Serve inits the gRPC server.
func (s *Serve) Serve() {
	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", s.port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer(*s.opts...)

	pb.RegisterStockServiceServer(grpcServer, s.handler)

	s.logger.Info(fmt.Sprintf("Serving gRPC on: %d", s.port))
	grpcServer.Serve(lis)

	defer s.wg.Done()
}
