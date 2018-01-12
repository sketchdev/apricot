package db

import (
	"reflect"
	"testing"
)

func TestNewManagerFromEngine(t *testing.T) {
	type args struct {
		name string
	}
	tests := []struct {
		name    string
		args    args
		want    DatabaseManager
		wantErr bool
	}{
		{"PostgresEngine", args{name: "postgres"}, &postgres{}, false},
		{"InvalidEngine", args{name: "invalid"}, nil, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewManagerFromEngine(tt.args.name)
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
