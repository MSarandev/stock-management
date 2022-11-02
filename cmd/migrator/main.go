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

func functionsRegistrar() [4]fnFlags {
	ctx := context.Background()

	migrate := &cobra.Command{
		// TODO: implement this
		Use:   "migrate",
		Short: "Migrate applies migrations",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("In migrate")
		},
	}

	rollback := &cobra.Command{
		// TODO: implement this
		Use:   "rollback",
		Short: "Rollback handles rollbacks",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("In rollback")
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

			return migrator.GenerateMigration(ctx, name)
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

	return [4]fnFlags{
		{
			fn:    generate,
			flag:  _generateFlag,
			usage: "generation name",
		},
		{
			fn:    migrate,
			flag:  "",
			usage: "",
		},
		{
			fn:    status,
			flag:  "",
			usage: "",
		},
		{
			fn:    rollback,
			flag:  "",
			usage: "",
		},
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

	migrator := migrate.NewMigrator(
		instance.Base,
		migrate.NewMigrations(migrate.WithMigrationsDirectory("migrations")),
	)

	return db.NewMigrator(logger, instance, migrator), nil
}
