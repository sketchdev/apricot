package lib

import (
	"path"
	"reflect"
	"testing"
)

func TestNewConfiguration(t *testing.T) {
	type args struct {
		engine  string
		connStr string
	}
	tests := []struct {
		name string
		args args
		want Configuration
	}{
		{"WithEngine", args{engine: "postgres"}, Configuration{Engine: "postgres", Folders: []string{path.Join("migrations", "current")}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewConfiguration(tt.args.engine, tt.args.connStr); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewConfiguration() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewConfigurationFromFile(t *testing.T) {
	type args struct {
		name string
	}
	tests := []struct {
		name    string
		args    args
		want    Configuration
		wantErr bool
	}{
		{"BadFile", args{name: "badfile"}, Configuration{}, true},
		{"GoodFile", args{name: path.Join("..", "testdata", "postgres", "apricot.toml")}, Configuration{Engine: "postgres", Folders: []string{"testdata/postgres/current", "testdata/postgres/release1"}, ConnectionFile: path.Join("..", "testdata", "postgres", "test.conn")}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewConfigurationFromFile(tt.args.name)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewConfigurationFromFile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewConfigurationFromFile() = %v, want %v", got, tt.want)
			}
		})
	}
}
