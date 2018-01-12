package app

import (
	"path"
	"reflect"
	"testing"

	"github.com/sketchdev/apricot/db"
	"github.com/sketchdev/apricot/lib"
)

func TestNewApricotFromConfigurationFile(t *testing.T) {
	type args struct {
		name string
	}
	postgres, _ := db.NewManagerFromEngine("postgres", "pg://postgres@localhost/apricot?sslmode=disable")
	postgresConfig := lib.Configuration{Engine: "postgres", Folders: []string{path.Join("testdata", "postgres", "current"), path.Join("testdata", "postgres", "release1")}, ConnectionFile: path.Join("..", "testdata", "postgres", "test.conn")}
	tests := []struct {
		name    string
		args    args
		want    Apricot
		wantErr bool
	}{
		{"ShouldBuildFromPostgresFile", args{name: path.Join("..", "testdata", "postgres", "apricot.toml")}, Apricot{Configuration: postgresConfig, DatabaseManager: postgres}, false},
		{"ShouldHandleBadTomlFile", args{name: path.Join("..", "testdata", "postgres", "bad.toml")}, Apricot{}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewApricotFromConfigurationFile(tt.args.name)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewApricotFromConfigurationFile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewApricotFromConfigurationFile() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewApricotFromConfiguration(t *testing.T) {
	type args struct {
		configuration lib.Configuration
	}
	postgres, _ := db.NewManagerFromEngine("postgres", "pg://postgres@localhost/apricot?sslmode=disable")
	postgresConfig := lib.NewConfiguration("postgres", path.Join("..", "testdata", "postgres", "test.conn"))
	badConfig := lib.Configuration{Engine: "invalid"}
	tests := []struct {
		name    string
		args    args
		want    Apricot
		wantErr bool
	}{
		{"ShouldBuildWithConfiguration", args{configuration: postgresConfig}, Apricot{Configuration: postgresConfig, DatabaseManager: postgres}, false},
		{"ShouldHandleBadEngine", args{configuration: badConfig}, Apricot{}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewApricotFromConfiguration(tt.args.configuration)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewApricotFromConfiguration() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewApricotFromConfiguration() = %v, want %v", got, tt.want)
			}
		})
	}
}
