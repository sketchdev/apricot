package commands

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path"
)

func RunInitialize() error {
	// create migrations folder if needed
	if err := os.Mkdir("migrations", os.ModeDir|os.ModePerm); err != nil {
		return errors.New(fmt.Sprintf("Failed to create migrations directory. %v\n", err))
	}
	// create apricot.toml file
	contents := []byte("engine = \"postgres\"\nmigrations = [\"migrations/current\"]\n")
	if err := ioutil.WriteFile("migrations/apricot.toml", contents, os.ModePerm); err != nil {
		return errors.New(fmt.Sprintf("Failed to create the apricot.toml file. %v\n", err))
	}
	// create current directory
	if err := os.Mkdir(path.Join("migrations", "current"), os.ModeDir|os.ModePerm); err != nil {
		fmt.Errorf("Failed to create migrations/current directory. %v\n", err)
		return errors.New(fmt.Sprintf("Failed to create migrations/current directory. %v\n", err))
	}
	return nil
}
