package db

import (
	"reflect"
	"testing"
)

func TestNewManagerFromEngine(t *testing.T) {
	type args struct {
		name    string
		connStr string
	}
	tests := []struct {
		name    string
		args    args
		want    DatabaseManager
		wantErr bool
	}{
		{"PostgresEngine", args{"postgres", "pg://postgres@localhost/apricot?sslmode=disable"}, &postgres{connectionString: "pg://postgres@localhost/apricot?sslmode=disable"}, false},
		{"InvalidEngine", args{"invalid", ""}, nil, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewManagerFromEngine(tt.args.name, tt.args.connStr)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewManagerFromEngine() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewManagerFromEngine() = %v, want %v", got, tt.want)
			}
		})
	}
}
