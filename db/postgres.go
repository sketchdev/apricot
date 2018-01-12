package db

import (
	"database/sql"
	"fmt"
	"time"

	// used to import the postgres database driver
	_ "github.com/lib/pq"
)

type postgres struct {
	db *sql.DB
}

func (p postgres) Username() string {
	return "postgres"
}

func (p postgres) connectionString() string {
	return "user=postgres dbname=apricot sslmode=disable" // TODO: read from config file
}

func (p *postgres) connect() error {
	if p.db != nil {
		return nil
	}
	db, err := sql.Open("postgres", p.connectionString())
	if err != nil {
		return err
	}
	if err := db.Ping(); err != nil {
		return err
	}
	p.db = db
	return nil
}

func (p postgres) Close() {
	if p.db != nil {
		p.db.Close()
	}
}

func (p postgres) CreateSchemaTable() error {
	err := p.connect()
	if err != nil {
		return err
	}
	stmt := fmt.Sprintf(`
		create table %s (
			id serial primary key not null,
			version varchar(14) not null,
			description varchar(255) not null,
			filename varchar(1024) not null,
			username varchar(100) not null,
			created_at timestamp default now() not null,
			duration int null,
			state varchar(7) default 'unknown' not null
		)
	`, SchemaTableName)
	if _, err := p.db.Exec(stmt); err != nil {
		return err
	}
	return nil
}

func (p *postgres) SchemaTableMissing() (bool, error) {
	return p.TableMissing(SchemaTableName)
}

func (p postgres) DropTable(name string) error {
	err := p.connect()
	if err != nil {
		return err
	}
	stmt := fmt.Sprintf("drop table if exists %s", name)
	if _, err := p.db.Exec(stmt); err != nil {
		return err
	}
	return nil
}

func (p postgres) TableExists(name string) (bool, error) {
	err := p.connect()
	if err != nil {
		return false, err
	}
	row := p.db.QueryRow("select exists(select 1 from pg_tables where tablename = $1);", name)
	var exists bool
	err = row.Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}

func (p postgres) TableMissing(name string) (bool, error) {
	exists, err := p.TableExists(name)
	if err != nil {
		return true, err
	}
	return !exists, nil
}

func (p postgres) AnyNonSuccessfulMigrations() (bool, error) {
	err := p.connect()
	if err != nil {
		return false, err
	}
	row := p.db.QueryRow(fmt.Sprintf("select count(id) as row_count from %s where state != 'success'", SchemaTableName))
	var rowCount int
	err = row.Scan(&rowCount)
	if err != nil {
		return rowCount > 0, err
	}
	return rowCount > 0, nil
}

func (p postgres) MigrationMissing(version string) (bool, error) {
	err := p.connect()
	if err != nil {
		return false, err
	}
	row := p.db.QueryRow(fmt.Sprintf("select count(id) as row_count from %s where version = $1", SchemaTableName), version)
	var rowCount int
	err = row.Scan(&rowCount)
	if err != nil {
		return rowCount == 0, err
	}
	return rowCount == 0, nil
}

func (p postgres) StartMigration(version string, description string, filename string) (int, error) {
	err := p.connect()
	if err != nil {
		return 0, err
	}
	tx, err := p.db.Begin()
	if err != nil {
		return 0, err
	}
	row := tx.QueryRow(fmt.Sprintf("insert into %s (version, description, filename, username) values ($1, $2, $3, $4) returning id", SchemaTableName), version, description, filename, p.Username())
	var id int
	err = row.Scan(&id)
	if err != nil {
		tx.Rollback()
		return 0, err
	}
	tx.Commit()
	return id, nil
}

func (p postgres) ApplyMigration(contents string) error {
	err := p.connect()
	if err != nil {
		return err
	}
	tx, err := p.db.Begin()
	if err != nil {
		return err
	}
	_, err = tx.Exec(contents)
	if err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}

func (p postgres) RollbackMigration(id int) error {
	err := p.connect()
	if err != nil {
		return err
	}
	tx, err := p.db.Begin()
	if err != nil {
		return err
	}
	_, err = tx.Exec(fmt.Sprintf("delete from %s where id = $1", SchemaTableName), id)
	if err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}

func (p postgres) EndMigration(id int, duration time.Duration) error {
	err := p.connect()
	if err != nil {
		return err
	}
	tx, err := p.db.Begin()
	if err != nil {
		return err
	}
	_, err = tx.Exec(fmt.Sprintf("update %s set state = $1, duration = $2 where id = $3", SchemaTableName), "success", duration, id)
	if err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}
