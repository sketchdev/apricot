package db

import (
	"errors"
	"time"
)

// SchemaTableName is the default schema table name
const SchemaTableName = "apricot_migrations"

// DatabaseManager is the interface implemented by each specific database
type DatabaseManager interface {
	Username() string
	Close()
	CreateSchemaTable() error
	SchemaTableMissing() (bool, error)
	DropTable(string) error
	TableExists(string) (bool, error)
	TableMissing(string) (bool, error)
	AnyNonSuccessfulMigrations() (bool, error)
	MigrationMissing(version string) (bool, error)
	StartMigration(version string, description string, filename string) (int, error)
	ApplyMigration(contents string) error
	RollbackMigration(id int) error
	EndMigration(id int, duration time.Duration) error
}

// NewManagerFromEngine is a factory method to produce DatabaseManager types by a database engine name
func NewManagerFromEngine(name string) (DatabaseManager, error) {
	switch name {
	case "postgres":
		return &postgres{}, nil
	}
	return nil, errors.New("invalid database engine")
}
