package main

import (
	"context"
	"log"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"stocks-api/module/controllers"
	"stocks-api/support/db"
	"stocks-api/support/server"
)

func init() {
	if err := godotenv.Load(".env"); err != nil {
		log.Print("No .env file found")
	}
}

func main() {
	logger := logrus.New()

	db := prepDB(logger)
	s := prepServer(logger, db, context.Background())

	s.Serve()
}

func prepDB(l *logrus.Logger) *db.Instance {
	conn, err := db.NewConnection()
	if err != nil {
		l.Fatal(err)
	}

	instance := db.NewInstance(conn, l)
	l.Log(logrus.InfoLevel, "DB conn: ", instance.Health())

	return instance
}

func prepServer(l *logrus.Logger, db *db.Instance, ctx context.Context) *server.Serve {
	controller := controllers.NewStockController(l, db, ctx)

	return server.NewServe(controller, l)
}
