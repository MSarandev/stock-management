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

const envParams string = "{\n  \"host\": \"DB_HOST\",\n  \"port\": \"DB_PORT\",\n  \"user\": \"DB_USER\",\n  \"pass\": \"DB_PASS\",\n  \"db_name\": \"DB_NAME\"\n}"

type dbConn struct {
	Host   string `json:"host"`
	Port   string `json:"port"`
	User   string `json:"user"`
	Pass   string `json:"pass"`
	DbName string `json:"db_name"`
}

func unwrapEnv() (*dbConn, error) {
	params := map[string]string{}
	if err := json.Unmarshal([]byte(envParams), &params); err != nil {
		return nil, err
	}

	for key, param := range params {
		val, ok := os.LookupEnv(param)
		if !ok {
			return nil, errors.New(fmt.Sprintf("failed to load '%s' env param", param))
		}

		params[key] = val
	}

	m, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}

	c := dbConn{}

	if err := json.Unmarshal(m, &c); err != nil {
		return nil, err
	}

	return &c, nil
}

func generateDsn() (string, error) {
	db, err := unwrapEnv()
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		db.User,
		db.Pass,
		db.Host,
		db.Port,
		db.DbName,
	), nil
}

// NewConnection - initialises a db connection.
func NewConnection() (*bun.DB, error) {
	dsn, err := generateDsn()
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
