package lib

import (
	"io/ioutil"
	"sort"
	"strings"
)

// UpSuffix is the suffix for up migrations
const UpSuffix = ".up.sql"

// AlwaysSuffix is the suffix for always migrations
const AlwaysSuffix = ".always.sql"

// DownSuffix is the suffix for down migrations
const DownSuffix = ".down.sql"

// GatherFiles files in the directories supplied and returns an array of MigrationFile values if
// the files match any of the suffixes provided
func GatherFiles(dirs []string, suffixes []string) ([]MigrationFile, error) {
	var migrationFiles []MigrationFile
	for _, dirname := range dirs {
		files, err := ioutil.ReadDir(dirname)
		if err != nil {
			return nil, err
		}
		for _, file := range files {
			for _, suffix := range suffixes {
				if strings.HasSuffix(file.Name(), suffix) {
					splitted := strings.Split(file.Name()[15:], ".")
					migrationFile := MigrationFile{
						Dirname:     dirname,
						Filename:    file.Name(),
						Version:     file.Name()[0:14],
						Description: splitted[0],
						Type:        MigrationType(splitted[1]),
					}
					migrationFiles = append(migrationFiles, migrationFile)
				}
			}
		}
	}
	sort.Sort(ByVersion(migrationFiles))
	return migrationFiles, nil
}

// ByVersion is a type used for sorting files by version in ascending order
type ByVersion []MigrationFile

func (a ByVersion) Len() int {
	return len(a)
}

func (a ByVersion) Less(i, j int) bool {
	return a[i].Version < a[j].Version
}

func (a ByVersion) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}
