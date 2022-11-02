package db

import (
	"context"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/uptrace/bun/migrate"
)

// Migrator handles everything migration related.
type Migrator struct {
	logger      *logrus.Logger
	dbInstance  *Instance
	bunMigrator *migrate.Migrator
}

// NewMigrator migrator constructor.
func NewMigrator(logger *logrus.Logger, db *Instance, migrator *migrate.Migrator) *Migrator {
	return &Migrator{
		logger:      logger,
		dbInstance:  db,
		bunMigrator: migrator,
	}
}

// GenerateMigrations generates SQL based migration files.
func (m *Migrator) GenerateMigrations(ctx context.Context, name string) error {
	files, err := m.bunMigrator.CreateSQLMigrations(ctx, name)
	if err != nil {
		return errors.WithStack(err)
	}

	for _, f := range files {
		m.logger.Infof("Generated migration file: %s \t at: %s", f.Name, f.Path)
	}

	return nil
}

// Migrate migrates all new tables.
func (m *Migrator) Migrate(ctx context.Context) error {
	group, err := m.bunMigrator.Migrate(ctx)
	if err != nil {
		m.logger.Error(err)
		return errors.WithStack(err)
	}

	if group.ID == 0 {
		m.logger.Info("No new migrations to run")
	}

	return nil
}

// Init initialises bun's migration tables.
// TODO: add an entry in the Readme for this
func (m *Migrator) Init(ctx context.Context) error {
	return m.bunMigrator.Init(ctx)
}

// Rollback rolls-back the last migration group.
func (m *Migrator) Rollback(ctx context.Context) error {
	group, err := m.bunMigrator.Rollback(ctx)
	if err != nil {
		m.logger.Error(err)
		return errors.WithStack(err)
	}

	if group.ID == 0 {
		m.logger.Info("Nothing to rollback")
	}

	return nil
}
