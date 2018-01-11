package app

import (
	"path"
	"testing"

	_ "github.com/lib/pq"
	"github.com/sketchdev/apricot/db"
	"github.com/sketchdev/apricot/lib"
)

var engines = []string{"postgres"}

func TestUp(t *testing.T) {
	for _, engine := range engines {
		t.Run(engine, testUp(engine))
	}
}

func testUp(engine string) func(t *testing.T) {
	return func(t *testing.T) {
		// setup
		conn := openConnection(engine, t)
		defer conn.Close()
		// clean
		dropTable("authors", conn, t)
		dropTable(db.SchemaTableName, conn, t)
		// run
		apricot := buildApricot(engine, t)
		if err := apricot.RunUp(); err != nil {
			t.Fatal(err)
		}
		// assert
		t.Run("ShouldCreateSchemaTable", assertTableExists(conn, db.SchemaTableName))
		t.Run("ShouldCreateAuthors", assertTableExists(conn, "authors"))
		// TODO: assert the migration records
		//t.Run("ShouldCreateMigrationRecord1", assertTableExists(conn, "authors"))
		//t.Run("ShouldCreateMigrationRecord2", assertTableExists(conn, "authors"))
	}
}

func assertTableExists(conn db.DatabaseManager, name string) func(t *testing.T) {
	return func(t *testing.T) {
		if conn.TableMissing(name) {
			t.Fatal("table not found")
		}
	}
}

func buildApricot(engine string, t *testing.T) Apricot {
	t.Helper()
	configuration := lib.NewConfiguration(engine)
	configuration.Migrations = []string{path.Join("..", "testdata", engine, "release1")}
	apricot, err := NewApricotFromConfiguration(configuration)
	if err != nil {
		t.Fatal(err)
	}
	return apricot
}

func openConnection(engine string, t *testing.T) db.DatabaseManager {
	t.Helper()
	conn, err := db.NewManagerFromEngine(engine)
	if err != nil {
		t.Fatal(err)
	}
	err = conn.Connect()
	if err != nil {
		t.Fatal(err)
	}
	return conn
}

func dropTable(name string, conn db.DatabaseManager, t *testing.T) {
	t.Helper()
	if err := conn.DropTable(name); err != nil {
		t.Fatal(err)
	}
}
