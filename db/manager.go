package db

import (
	"errors"
)

const SchemaTableName = "apricot_migrations"

type DatabaseManager interface {
	Connect() error
	Close()
	CreateSchemaTable() error
	DropSchemaTable() error
	SchemaTableExists() bool
	SchemaTableMissing() bool
	DropTable(string) error
	TableExists(string) bool
	TableMissing(string) bool
}

func NewManagerFromEngine(name string) (DatabaseManager, error) {
	switch name {
	case "postgres":
		return &postgres{}, nil
	}
	return nil, errors.New("invalid database engine")
}
