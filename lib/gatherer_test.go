package lib

import (
	"path"
	"reflect"
	"testing"
)

type args struct {
	dirs     []string
	suffixes []string
}

func TestGatherFiles(t *testing.T) {
	tests := []struct {
		name    string
		args    args
		want    []MigrationFile
		wantErr bool
	}{
		{"UpAndAlways", upAndAlways(), assertUpAndAlways(), false},
		{"BadDir", badDir(), nil, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GatherFiles(tt.args.dirs, tt.args.suffixes)
			if (err != nil) != tt.wantErr {
				t.Errorf("GatherFiles() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GatherFiles() = %v, want %v", got, tt.want)
			}
		})
	}
}

func assertUpAndAlways() []MigrationFile {
	return []MigrationFile{
		{
			Dirname:     "../testdata/postgres/release1",
			Filename:    "20170108115705_create_authors_table.up.sql",
			Description: "create_authors_table",
			Type:        MigrationTypeUp,
			Version:     "20170108115705",
		},
		{
			Dirname:     "../testdata/postgres/release1",
			Filename:    "20170108115805_contrived_select.always.sql",
			Description: "contrived_select",
			Type:        MigrationTypeAlways,
			Version:     "20170108115805",
		},
		{
			Dirname:     "../testdata/postgres/current",
			Filename:    "20170109115705_alter_authors_table.up.sql",
			Description: "alter_authors_table",
			Type:        MigrationTypeUp,
			Version:     "20170109115705",
		},
	}
}

func upAndAlways() args {
	current := path.Join("..", "testdata", "postgres", "current")
	release1 := path.Join("..", "testdata", "postgres", "release1")
	return args{
		dirs:     []string{current, release1},
		suffixes: []string{UpSuffix, AlwaysSuffix},
	}
}

func badDir() args {
	current := "blah"
	return args{
		dirs:     []string{current},
		suffixes: []string{UpSuffix, AlwaysSuffix},
	}
}
