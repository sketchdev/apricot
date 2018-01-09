package app

import (
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
		conn := openConnection(engine, t)
		defer conn.Close()
		cleanTable(conn, t)
		apricot := buildApricot(engine, t)
		apricot.RunUp()
		t.Run("ShouldCreateSchemaTable", assertTableExists(conn, db.SchemaTableName))
		//t.Run("ShouldCreateAuthors", assertTableExists(conn, "authors"))
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
	apricot, err := NewApricotFromConfiguration(lib.NewConfiguration(engine))
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

func cleanTable(conn db.DatabaseManager, t *testing.T) {
	t.Helper()
	if err := conn.DropTable(db.SchemaTableName); err != nil {
		t.Fatal(err)
	}
}
