package app

import (
	"github.com/sketchdev/apricot/db"
	"github.com/sketchdev/apricot/lib"
)

// Apricot is the primary type which controls migrations
type Apricot struct {
	DatabaseManager db.DatabaseManager
	Configuration   lib.Configuration
}

// NewApricotFromConfigurationFile returns a new Apricot based on the file provided
func NewApricotFromConfigurationFile(name string) (Apricot, error) {
	configuration, err := lib.NewConfigurationFromFile(name)
	if err != nil {
		return Apricot{}, err
	}
	return NewApricotFromConfiguration(configuration)
}

// NewApricotFromConfiguration returns a new Apricot based on the configuration provided
func NewApricotFromConfiguration(configuration lib.Configuration) (Apricot, error) {
	connStr, err := configuration.ConnectionString()
	if err != nil {
		return Apricot{}, err
	}
	databaseManager, err := db.NewManagerFromEngine(configuration.Engine, connStr)
	if err != nil {
		return Apricot{}, err
	}
	apricot := Apricot{}
	apricot.Configuration = configuration
	apricot.DatabaseManager = databaseManager
	return apricot, nil
}
