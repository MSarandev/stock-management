package main

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"stocks-api/support/db"
)

func init() {
	if err := godotenv.Load(".env"); err != nil {
		log.Print("No .env file found")
	}
}

func main() {
	logger := logrus.New()

	conn, err := db.NewConnection()
	if err != nil {
		logger.Fatal(err)
	}

	instance := db.NewInstance(conn, logger)

	ok := instance.Health()

	logger.Log(logrus.InfoLevel, "DB conn: ", ok)
}
