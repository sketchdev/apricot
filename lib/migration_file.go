package lib

import (
	"io/ioutil"
	"path"
)

// MigrationType is used to contain types of migrations (up, down, always, invalid)
type MigrationType string

const (
	// MigrationTypeUp represents "up" migrations
	MigrationTypeUp = "up"
	// MigrationTypeDown represents "down" migrations
	MigrationTypeDown = "down"
	// MigrationTypeAlways represents "always" migrations
	MigrationTypeAlways = "always"
)

// MigrationFile represents the parsed data from an indiviual migration file
type MigrationFile struct {
	Dirname     string
	Filename    string
	Version     string
	Description string
	Type        MigrationType
}

// Contents reads then entire contents of the MigrationFile
func (f MigrationFile) Contents() (string, error) {
	data, err := ioutil.ReadFile(path.Join(f.Dirname, f.Filename))
	if err != nil {
		return "", err
	}
	return string(data), nil
}
