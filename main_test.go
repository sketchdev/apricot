package main

import (
	"os"
	"testing"
)

func TestInitialize(t *testing.T) {
	os.RemoveAll("migrations")
	if err := RunInitialize(); err != nil {
		t.Error(err)
	} else {
		assertDirectoryExists("migrations", t)
		assertDirectoryExists("migrations/current", t)
	}
	os.RemoveAll("migrations")
}

func assertDirectoryExists(dir string, t *testing.T) {
	if _, err := os.Stat(dir); err != nil {
		t.Error(err)
	}
}
