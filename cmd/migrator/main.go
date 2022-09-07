package main

import (
	"context"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"stocks-api/module/entities"
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

	flag := os.Args[1]

	instance := db.NewInstance(conn, logger)
	logger.Log(logrus.InfoLevel, "DB conn: ", instance.Health())

	migrator := db.NewMigrator(logger, instance)

	if flag == "migrate" {
		migrate(migrator)
	}

	if flag == "rollback" {
		rollback(migrator)
	}
}

func migrate(m *db.Migrator) {
	m.MigrateOne(context.Background(), &entities.Stock{}, "stock")
}

func rollback(m *db.Migrator) {
	m.RollbackOne(context.Background(), &entities.Stock{}, "stock")
}
