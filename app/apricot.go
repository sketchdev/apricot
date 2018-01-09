package app

import (
	"github.com/sketchdev/apricot/db"
	"github.com/sketchdev/apricot/lib"
)

type Apricot struct {
	DatabaseManager db.DatabaseManager
	Configuration   lib.Configuration
}

func NewApricotFromConfigurationFile(name string) (Apricot, error) {
	if configuration, err := lib.NewConfigurationFromFile(name); err != nil {
		return Apricot{}, err
	} else {
		return NewApricotFromConfiguration(configuration)
	}
}

func NewApricotFromConfiguration(configuration lib.Configuration) (Apricot, error) {
	if databaseManager, err := db.NewManagerFromEngine(configuration.Engine); err != nil {
		return Apricot{}, err
	} else {
		apricot := Apricot{}
		apricot.Configuration = configuration
		apricot.DatabaseManager = databaseManager
		return apricot, nil
	}
}
