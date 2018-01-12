package app

import (
	"path"
	"testing"

	_ "github.com/lib/pq"
	"github.com/sketchdev/apricot/db"
	"github.com/sketchdev/apricot/lib"
)

//var engines = []string{"postgres"}
//
//func TestUp(t *testing.T) {
//	for _, engine := range engines {
//		t.Run(engine, testUp(engine))
//	}
//}
//
//func testUp(engine string) func(t *testing.T) {
//	return func(t *testing.T) {
//		// setup
//		conn := openConnection(engine, t)
//		defer conn.Close()
//		// clean
//		dropTable("authors", conn, t)
//		dropTable(db.SchemaTableName, conn, t)
//		// run
//		apricot := buildApricot(engine, t)
//		if err := apricot.RunUp(); err != nil {
//			t.Fatal(err)
//		}
//		// assert
//		t.Run("ShouldCreateSchemaTable", assertTableExists(conn, db.SchemaTableName))
//		t.Run("ShouldCreateAuthors", assertTableExists(conn, "authors"))
//		// TODO: assert the migration records
//		//t.Run("ShouldCreateMigrationRecord1", assertTableExists(conn, "authors"))
//		//t.Run("ShouldCreateMigrationRecord2", assertTableExists(conn, "authors"))
//	}
//}
//
//func assertTableExists(conn db.DatabaseManager, name string) func(t *testing.T) {
//	return func(t *testing.T) {
//		if conn.TableMissing(name) {
//			t.Fatal("table not found")
//		}
//	}
//}
//
//func buildApricot(engine string, t *testing.T) Apricot {
//	t.Helper()
//	configuration := lib.NewConfiguration(engine)
//	configuration.Migrations = []string{path.Join("..", "testdata", engine, "release1")}
//	apricot, err := NewApricotFromConfiguration(configuration)
//	if err != nil {
//		t.Fatal(err)
//	}
//	return apricot
//}
//
//func openConnection(engine string, t *testing.T) db.DatabaseManager {
//	t.Helper()
//	conn, err := db.NewManagerFromEngine(engine)
//	if err != nil {
//		t.Fatal(err)
//	}
//	err = conn.Connect()
//	if err != nil {
//		t.Fatal(err)
//	}
//	return conn
//}
//
//func dropTable(name string, conn db.DatabaseManager, t *testing.T) {
//	t.Helper()
//	if err := conn.DropTable(name); err != nil {
//		t.Fatal(err)
//	}
//}

var engines = []string{"postgres"}

func TestApricot_RunUp(t *testing.T) {
	for _, engine := range engines {
		t.Run(engine, testApricotRunUp(engine))
	}
}

func testApricotRunUp(engine string) func(t *testing.T) {
	return func(t *testing.T) {
		type fields struct {
			DatabaseManager db.DatabaseManager
			Configuration   lib.Configuration
			TablesToDrop    []string
		}
		databaseManager, _ := db.NewManagerFromEngine(engine)
		configuration := lib.NewConfiguration(engine)
		configuration.Migrations = []string{path.Join("..", "testdata", engine, "current"), path.Join("..", "testdata", engine, "release1")}
		tests := []struct {
			name    string
			fields  fields
			wantErr bool
		}{
			{"ShouldCreateMigrations", fields{databaseManager, configuration, []string{"authors", db.SchemaTableName}}, false},
		}
		for _, tt := range tests {
			// drop tables if needed
			for _, table := range tt.fields.TablesToDrop {
				exists, err := tt.fields.DatabaseManager.TableExists(table)
				if err != nil {
					t.Fatalf("failed to determine if table %s exists: %v", table, err)
				}
				if exists {
					err := tt.fields.DatabaseManager.DropTable(table)
					if err != nil {
						t.Fatalf("failed to drop table: %s", table)
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
		}
	}
}
