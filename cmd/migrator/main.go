package main

import (
	"context"

	"github.com/sirupsen/logrus"
	"stocks-api/module/entities"
	"stocks-api/support/db"
)

func main() {
	logger := logrus.New()

	conn, err := db.NewConnection()
	if err != nil {
		logger.Fatal(err)
	}

	instance := db.NewInstance(conn, logger)

	ok := instance.Health()

	logger.Log(logrus.InfoLevel, "DB conn: ", ok)

	migrator := db.NewMigrator(logger, instance)

	migrator.MigrateOne(context.Background(), &entities.Stock{}, "stock")
}
