package app

import (
	"errors"
	"fmt"
	"time"

	"github.com/sketchdev/apricot/lib"
)

// RunUp executes the up command
func (a Apricot) RunUp() error {
	steps := lib.NewSteps()
	steps.PrintHeader()
	steps.PrintRule()

	// defer closing the database
	defer a.DatabaseManager.Close()

	// create schema table if needed
	steps.Start("Verifying schema table")
	missing, err := a.DatabaseManager.SchemaTableMissing()
	if err != nil {
		steps.Error()
		return err
	}
	if missing {
		steps.Done()
		steps.Start("Creating schema table")
		err := a.DatabaseManager.CreateSchemaTable()
		if err != nil {
			steps.Error()
			return err
		}
		steps.Success()
	}

	// look to see if any rows are not successful
	steps.Start("Validating previous migrations")
	nonSuccessfulMigrationsExist, err := a.DatabaseManager.AnyNonSuccessfulMigrations()
	if err != nil {
		steps.Error()
		return err
	}
	if nonSuccessfulMigrationsExist {
		steps.Fail()
		return errors.New("not all migrations are successful. please check the schema table to see which migrations are not successful")
	}
	steps.Done()

	// gather migration files
	steps.Start("Gathering migration files")
	files, err := lib.GatherFiles(a.Configuration.Migrations, []string{lib.UpSuffix, lib.AlwaysSuffix})
	if err != nil {
		steps.Error()
		return err
	}
	steps.Done()

	// check for duplicate migration versions in order to fail fast
	steps.Start("Finding duplicate migrations")
	duplicateMap := make(map[string]bool)
	for _, file := range files {
		if _, ok := duplicateMap[file.Version]; ok {
			steps.Fail()
			return fmt.Errorf("duplicate migration found (%s) no migrations applied", file.Version)
		}
		duplicateMap[file.Version] = true
	}
	steps.Done()

	// filter migration files to only those which haven't been applied
	steps.Start("Determining which migrations to apply")
	var migrationsToBeApplied []lib.MigrationFile
	for _, file := range files {
		if file.Type == lib.MigrationTypeAlways {
			migrationsToBeApplied = append(migrationsToBeApplied, file)
		} else {
			missing, err := a.DatabaseManager.MigrationMissing(file.Version)
			if err != nil {
				steps.Error()
				return err
			}
			if missing {
				migrationsToBeApplied = append(migrationsToBeApplied, file)
			}
		}
	}
	steps.Done()
	steps.PrintRule()

	// TODO: replace tokens in migration files

	for _, migration := range migrationsToBeApplied {
		steps.Start(migration.Filename)
		// read the migration file
		contents, err := migration.Contents()
		if err != nil {
			steps.Error()
			return err
		}
		// create the migration record in the schema table
		id, err := a.DatabaseManager.StartMigration(migration.Version, migration.Description, migration.Filename)
		if err != nil {
			steps.Error()
			return err
		}
		// apply the migration
		startTime := time.Now()
		err = a.DatabaseManager.ApplyMigration(contents)
		if err != nil {
			steps.Error()
			a.DatabaseManager.RollbackMigration(id)
			return err
		}
		// finalize the migration
		durationMs := time.Since(startTime) / time.Millisecond
		err = a.DatabaseManager.EndMigration(id, durationMs)
		if err != nil {
			steps.Error()
			return err
		}
		steps.Done()
	}

	steps.PrintRule()
	return nil
}
