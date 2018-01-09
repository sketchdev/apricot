package db

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

type postgres struct {
	db *sql.DB
}

func (p *postgres) Connect() error {
	if db, err := sql.Open("postgres", p.connectionString()); err != nil {
		return err
	} else {
		if err := db.Ping(); err != nil {
			return err
		}
		p.db = db
		return nil
	}
}

func (p *postgres) Close() {
	p.db.Close()
}

func (p postgres) connectionString() string {
	return "user=postgres dbname=apricot sslmode=disable" // TODO: read from config file
}

func (p postgres) CreateSchemaTable() error {
	// TODO: add more columns
	stmt := fmt.Sprintf("create table %s (id SERIAL PRIMARY KEY NOT NULL, version varchar(14))", SchemaTableName)
	if _, err := p.db.Exec(stmt); err != nil {
		return err
	}
	return nil
}

func (p *postgres) DropSchemaTable() error {
	return p.DropTable(SchemaTableName)
}

func (p *postgres) SchemaTableExists() bool {
	return p.TableExists(SchemaTableName)
}

func (p *postgres) SchemaTableMissing() bool {
	return p.TableMissing(SchemaTableName)
}

func (p postgres) DropTable(name string) error {
	stmt := fmt.Sprintf("drop table if exists %s", name)
	if _, err := p.db.Exec(stmt); err != nil {
		return err
	}
	return nil
}

func (p postgres) TableExists(name string) bool {
	stmt := fmt.Sprintf("select count(id) as total_rows from %s", name)
	row := p.db.QueryRow(stmt)
	var totalRows int
	if err := row.Scan(&totalRows); err != nil {
		fmt.Errorf("error: %v", err)
		return false
	}
	return true
}

func (p postgres) TableMissing(name string) bool {
	return !p.TableExists(name)
}
