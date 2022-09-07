package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"stocks-api/module/controllers"
	"stocks-api/support/server"
)

func init() {
	if err := godotenv.Load(".env"); err != nil {
		log.Print("No .env file found")
	}
}

func main() {
	logger := logrus.New()

	port := os.Getenv("SERVER_PORT")
	address := os.Getenv("SERVER_ADDRESS")
	if port == "" || address == "" {
		logger.Fatal(errors.New("failed to load env param"))
	}

	stockController := controllers.NewStockController(logger)

	server := server.NewServe(stockController)

	logger.Info(fmt.Sprintf("Serving on: %s:%s", address, port))
	log.Fatal(http.ListenAndServe(fmt.Sprintf("%s:%s", address, port), server.Server))
}
