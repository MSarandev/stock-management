package db

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"os"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
	"github.com/uptrace/bun/extra/bundebug"
)

type dbConn struct {
	Host   string `json:"host"`
	Port   int    `json:"port"`
	User   string `json:"user"`
	Pass   string `json:"pass"`
	DbName string `json:"db_name"`
}

func generateDsn(envKey string) (string, error) {
	envDb, ok := os.LookupEnv(envKey)
	if !ok {
		return "", errors.New("failed to load DB env param")
	}

	db := dbConn{}
	if err := json.Unmarshal([]byte(envDb), &db); err != nil {
		return "", err
	}

	return fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
		db.User,
		db.Pass,
		db.Host,
		db.Port,
		db.DbName,
	), nil
}

// NewConnection - initialises a db connection.
func NewConnection() (*bun.DB, error) {
	dsn, err := generateDsn("DB")
	if err != nil {
		return nil, err
	}

	conn := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(dsn)))
	db := bun.NewDB(conn, pgdialect.New())

	db.AddQueryHook(bundebug.NewQueryHook(
		bundebug.WithVerbose(true),
		bundebug.FromEnv("BUNDEBUG"),
	))

	return db, nil
}
