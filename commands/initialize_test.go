package commands

import (
	"os"
	"testing"
)

func TestInitialize(t *testing.T) {
	os.RemoveAll("migrations")
	if err := RunInitialize(); err != nil {
		t.Error(err)
	} else {
		t.Run("Should create the migrations directory", shouldCreateDirectory("migrations"))
		t.Run("Should create the migrations/current directory", shouldCreateDirectory("migrations/current"))
	}
	os.RemoveAll("migrations")
}

func shouldCreateDirectory(dir string) func(t *testing.T) {
	return func(t *testing.T) {
		if _, err := os.Stat(dir); err != nil {
			t.Error(err)
		}
	}
}
