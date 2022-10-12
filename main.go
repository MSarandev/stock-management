package main

import (
	"context"
	"errors"
	"log"
	"os"
	"strconv"
	"sync"

	grpc_logrus "github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	rpc "google.golang.org/grpc"
	"stocks-api/module/controllers"
	"stocks-api/module/handlers"
	"stocks-api/support/db"
	"stocks-api/support/grpc"
	"stocks-api/support/http"
)

func init() {
	if err := godotenv.Load(".env"); err != nil {
		log.Print("No .env file found")
	}
}

func main() {
	logger := logrus.New()

	ctx := context.Background()

	wg := &sync.WaitGroup{}

	db := prepDB(logger)
	s := prepServer(logger, db, ctx, wg)

	g, err := prepGrpc(logger, db, ctx, wg)
	if err != nil {
		logger.Warningf("gRPC server failed to start: %s", err)
	}

	wg.Add(2)

	go s.Serve()
	go g.Serve()

	wg.Wait()
}

// prepDB prepare the database.
func prepDB(l *logrus.Logger) *db.Instance {
	conn, err := db.NewConnection()
	if err != nil {
		l.Fatal(err)
	}

	instance := db.NewInstance(conn, l)
	l.Log(logrus.InfoLevel, "DB conn: ", instance.Health())

	return instance
}

// prepServer prepare the HTTP server.
func prepServer(l *logrus.Logger, db *db.Instance, ctx context.Context, wg *sync.WaitGroup) *http.Serve {
	controller := controllers.NewStockController(l, db, ctx)

	return http.NewServe(controller, l, wg)
}

// prepGrpc prepare the gRPC server.
func prepGrpc(l *logrus.Logger, db *db.Instance, ctx context.Context, wg *sync.WaitGroup) (*grpc.Serve, error) {
	sPort, ok := os.LookupEnv("GRPC_PORT")
	if !ok {
		return nil, errors.New("failed to start gRPC server, missing port")
	}

	port := 0
	var err error

	if port, err = strconv.Atoi(sPort); err != nil {
		return nil, errors.New("failed to start gRPC server, faulty port")
	}

	logEntry := logrus.NewEntry(l)

	opts := []rpc.ServerOption{
		rpc.StreamInterceptor(grpc_logrus.StreamServerInterceptor(logEntry)),
		rpc.UnaryInterceptor(grpc_logrus.UnaryServerInterceptor(logEntry)),
	}

	handler := handlers.NewStockHandler(l, db, ctx)

	return grpc.NewServe(int64(port), l, &opts, handler, wg), nil
}
