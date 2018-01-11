package app

import (
	"errors"
	"fmt"
	"time"

	"github.com/sketchdev/apricot/lib"
)

// RunUp executes the up command
func (a Apricot) RunUp() error {
	// connect to the database and defer its closing
	if err := a.DatabaseManager.Connect(); err != nil {
		return err
	}
	defer a.DatabaseManager.Close()

	// create schema table if needed
	if a.DatabaseManager.SchemaTableMissing() {
		if err := a.DatabaseManager.CreateSchemaTable(); err != nil {
			return err
		}
	}

	// look to see if any rows are not successful
	nonSuccessfulMigrationsExist, err := a.DatabaseManager.AnyNonSuccessfulMigrations()
	if err != nil {
		return err
	}
	if nonSuccessfulMigrationsExist {
		return errors.New("not all migrations are successful. please check the schema table to see which migrations are not successful")
	}

	// gather migration files
	files, err := lib.GatherFiles(a.Configuration.Migrations, []string{lib.UpSuffix, lib.AlwaysSuffix})
	if err != nil {
		return err
	}

	// check for duplicate migration versions in order to fail fast
	duplicateMap := make(map[string]bool)
	for _, file := range files {
		if _, ok := duplicateMap[file.Version]; ok {
			return fmt.Errorf("duplicate migration found (%s) no migrations applied", file.Version)
		}
		duplicateMap[file.Version] = true
	}

	// filter migration files to only those which haven't been applied
	var migrationsToBeApplied []lib.MigrationFile
	for _, file := range files {
		if file.Type == lib.MigrationTypeAlways {
			migrationsToBeApplied = append(migrationsToBeApplied, file)
		} else {
			missing, err := a.DatabaseManager.MigrationMissing(file.Version)
			if err != nil {
				return err
			}
			if missing {
				migrationsToBeApplied = append(migrationsToBeApplied, file)
			}
		}
	}

	// TODO: replace tokens in migration files

	for _, migration := range migrationsToBeApplied {
		// read the migration file
		contents, err := migration.Contents()
		if err != nil {
			return err
		}
		// create the migration record in the schema table
		id, err := a.DatabaseManager.StartMigration(migration.Version, migration.Description, migration.Filename)
		if err != nil {
			return err
		}
		// apply the migration
		startTime := time.Now()
		err = a.DatabaseManager.ApplyMigration(contents)
		if err != nil {
			a.DatabaseManager.RollbackMigration(id)
			return err
		}
		// finalize the migration
		durationMs := time.Since(startTime) / time.Millisecond
		err = a.DatabaseManager.EndMigration(id, durationMs)
		if err != nil {
			return err
		}
	}
	return nil
}
