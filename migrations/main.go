package migrations

import (
	"embed"
)

//go:embed *.sql
var sqlMigrations embed.FS

// MigrationFilesScan returns all .sql files in the migration dir.
func MigrationFilesScan() embed.FS {
	return sqlMigrations
}
