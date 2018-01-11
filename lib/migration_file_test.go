package lib

import (
	"path"
	"testing"
)

func TestMigrationFile_Contents(t *testing.T) {
	type fields struct {
		Dirname     string
		Filename    string
		Version     string
		Description string
		Type        MigrationType
	}
	tests := []struct {
		name    string
		fields  fields
		want    string
		wantErr bool
	}{
		{"BadFile", fields{Dirname: "blah", Filename: "blah", Version: "blah", Description: "blah", Type: "blah"}, "", true},
		{"GoodFile", fields{Dirname: path.Join("..", "testdata", "postgres", "release1"), Filename: "20170108115805_contrived_select.always.sql", Version: "20170108115805", Description: "contrived_select", Type: "always"}, "select count(id) from authors;", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := MigrationFile{
				Dirname:     tt.fields.Dirname,
				Filename:    tt.fields.Filename,
				Version:     tt.fields.Version,
				Description: tt.fields.Description,
				Type:        tt.fields.Type,
			}
			got, err := f.Contents()
			if (err != nil) != tt.wantErr {
				t.Errorf("MigrationFile.Contents() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("MigrationFile.Contents() = %v, want %v", got, tt.want)
			}
		})
	}
}
