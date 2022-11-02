package main

import (
	"context"
	"fmt"
	"log"

	"github.com/joho/godotenv"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/uptrace/bun/migrate"
	migrations2 "stocks-api/migrations"
	"stocks-api/support/db"
)

type fnFlags struct {
	fn    *cobra.Command
	flag  string
	usage string
}

const (
	_generateFlag = "name"
)

func init() {
	if err := godotenv.Load(".env"); err != nil {
		log.Print("No .env file found")
	}
}

func main() {
	var rootCmd = &cobra.Command{}

	// Flags registration.
	for _, fn := range functionsRegistrar() {
		if fn.flag != "" {
			fn.fn.PersistentFlags().String(fn.flag, "", fn.usage)
		}

		rootCmd.AddCommand(fn.fn)
	}

	rootCmd.Execute()
}

func functionsRegistrar() [5]fnFlags {
	ctx := context.Background()

	init := &cobra.Command{
		Use:   "init",
		Short: "Initialises the bun_migration tables",
		RunE: func(cmd *cobra.Command, args []string) error {
			migrator, err := spinUpDb()
			if err != nil {
				return errors.WithStack(err)
			}

			return migrator.Init(ctx)
		},
	}

	migrate := &cobra.Command{
		Use:   "migrate",
		Short: "Migrate applies migrations",
		RunE: func(cmd *cobra.Command, args []string) error {
			migrator, err := spinUpDb()
			if err != nil {
				return errors.WithStack(err)
			}

			return migrator.Migrate(ctx)
		},
	}

	rollback := &cobra.Command{
		Use:   "rollback",
		Short: "Rollback handles rollbacks",
		RunE: func(cmd *cobra.Command, args []string) error {
			migrator, err := spinUpDb()
			if err != nil {
				return errors.WithStack(err)
			}

			return migrator.Rollback(ctx)
		},
	}

	generate := &cobra.Command{
		Use:     "generate",
		Example: fmt.Sprintf("generate --%s <migration-name>", _generateFlag),
		Short:   "Generates SQL migrations files",
		RunE: func(cmd *cobra.Command, args []string) error {
			name, err := cmd.Flags().GetString(_generateFlag)
			if err != nil {
				return errors.WithStack(err)
			}

			migrator, err := spinUpDb()
			if err != nil {
				return errors.WithStack(err)
			}

			return migrator.GenerateMigrations(ctx, name)
		},
	}

	status := &cobra.Command{
		// TODO: implement this
		Use:   "status",
		Short: "Returns the current status of the migrations",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("In status check")
		},
	}

	return [5]fnFlags{
		{fn: init},
		{fn: generate, flag: _generateFlag, usage: "generation name"},
		{fn: migrate},
		{fn: status},
		{fn: rollback},
	}
}

func spinUpDb() (*db.Migrator, error) {
	logger := logrus.New()

	conn, err := db.NewConnection()
	if err != nil {
		return nil, err
	}

	instance := db.NewInstance(conn, logger)
	logger.Log(logrus.InfoLevel, "DB conn: ", instance.Health())

	migrations := migrate.NewMigrations(migrate.WithMigrationsDirectory("migrations"))
	migrations.Discover(migrations2.MigrationFilesScan())

	migrator := migrate.NewMigrator(instance.Base, migrations)

	return db.NewMigrator(logger, instance, migrator), nil
}
