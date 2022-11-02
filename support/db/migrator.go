package db

import (
	"context"

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

// GenerateMigration generates SQL based migration files.
func (m *Migrator) GenerateMigration(ctx context.Context, name string) error {
	files, err := m.bunMigrator.CreateSQLMigrations(ctx, name)
	if err != nil {
		return err
	}

	for _, f := range files {
		m.logger.Infof("Generated migration file: %s, at: %s", f.Name, f.Path)
	}

	return nil
}
