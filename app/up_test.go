package app

import (
	"path"
	"testing"

	_ "github.com/lib/pq"
	"github.com/sketchdev/apricot/db"
	"github.com/sketchdev/apricot/lib"
)

var engines = map[string]string{"postgres": "pg://postgres@localhost/apricot?sslmode=disable"}

func TestApricot_RunUp(t *testing.T) {
	for engine, connStr := range engines {
		t.Run(engine, testApricotRunUp(engine, connStr))
	}
}

func testApricotRunUp(engine, connStr string) func(t *testing.T) {
	return func(t *testing.T) {
		type fields struct {
			DatabaseManager       db.DatabaseManager
			Configuration         lib.Configuration
			TablesToDropAndAssert []string
		}
		databaseManager, _ := db.NewManagerFromEngine(engine, connStr)
		tests := []struct {
			name    string
			fields  fields
			wantErr bool
		}{
			{"ShouldCreateMigrations", fields{databaseManager, goodConfiguration(engine), []string{"authors", db.SchemaTableName}}, false},
			{"ShouldReturnErrorIfBadDir", fields{databaseManager, configurationWithBadFolder(engine), []string{"authors", db.SchemaTableName}}, true},
			// TODO: add a case to verify that always are always run (seed schema table)
		}
		for _, tt := range tests {
			// drop tables if needed
			for _, table := range tt.fields.TablesToDropAndAssert {
				exists, err := tt.fields.DatabaseManager.TableExists(table)
				if err != nil {
					t.Fatalf("Apricot.RunUp() failed to determine if table %s exists: %v", table, err)
				}
				if exists {
					err := tt.fields.DatabaseManager.DropTable(table)
					if err != nil {
						t.Fatalf("Apricot.RunUp() failed to drop table: %s", table)
					}
				}
			}
			// run test
			t.Run(tt.name, func(t *testing.T) {
				a := Apricot{
					DatabaseManager: tt.fields.DatabaseManager,
					Configuration:   tt.fields.Configuration,
				}
				if err := a.RunUp(); (err != nil) != tt.wantErr {
					t.Errorf("Apricot.RunUp() error = %v, wantErr %v", err, tt.wantErr)
				}
			})
			// assert tables
			if !tt.wantErr {
				for _, table := range tt.fields.TablesToDropAndAssert {
					missing, err := tt.fields.DatabaseManager.TableMissing(table)
					if err != nil {
						t.Error(err)
					}
					if missing {
						t.Errorf("Apricot.RunUp() fail = %s table not found", table)
					}
				}
			}
		}
	}
}

func goodConfiguration(engine string) lib.Configuration {
	configuration := lib.NewConfiguration(engine, path.Join("..", "testdata", engine, "test.conn"))
	configuration.Folders = []string{path.Join("..", "testdata", engine, "current"), path.Join("..", "testdata", engine, "release1")}
	return configuration
}

func configurationWithBadFolder(engine string) lib.Configuration {
	configuration := goodConfiguration(engine)
	configuration.Folders = []string{path.Join("..", "testdata", engine, "current"), path.Join("..", "testdata", engine, "baddir")}
	return configuration
}
