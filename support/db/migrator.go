package db

import (
	"context"

	"github.com/sirupsen/logrus"
)

// Migrator handles everything migration related.
type Migrator struct {
	logger     *logrus.Logger
	dbInstance *Instance
}

// NewMigrator migrator constructor.
func NewMigrator(logger *logrus.Logger, db *Instance) *Migrator {
	return &Migrator{
		logger:     logger,
		dbInstance: db,
	}
}

// MigrateOne migrates a single model.
func (m *Migrator) MigrateOne(ctx context.Context, model interface{}, modelFlag string) error {
	m.logger.Infof("Starting migration for: %s", modelFlag)

	_, err := m.dbInstance.db.NewCreateTable().Model(model).Exec(ctx)
	if err != nil {
		m.logger.Error(err)
		return err
	}

	m.logger.Infof("Migration successful for: %s", modelFlag)
	return nil
}

// RollbackOne rollback a single model
func (m *Migrator) RollbackOne(ctx context.Context, model interface{}, modelFlag string) error {
	m.logger.Infof("Starting rollback for: %s", modelFlag)

	_, err := m.dbInstance.db.NewDropTable().Model(model).Exec(ctx)
	if err != nil {
		m.logger.Error(err)
		return err
	}

	m.logger.Infof("Rollback successful for: %s", modelFlag)
	return nil
}
