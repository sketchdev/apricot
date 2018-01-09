package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path"

	"github.com/sketchdev/apricot/app"
	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	rootCmd = kingpin.New("apricot", "A cross platform, cross language migration tool.")
	initCmd = rootCmd.Command("init", "Initializes apricot in the current working directory.")
	upCmd   = rootCmd.Command("up", "Migrate unapplied changes to the database.")
)

func main() {
	switch kingpin.MustParse(rootCmd.Parse(os.Args[1:])) {
	case initCmd.FullCommand():
		if err := RunInitialize(); err != nil {
			fmt.Errorf("error: %s", err)
		} else {
			fmt.Println("Initialization complete.")
		}
	case upCmd.FullCommand():
		if apricot, err := app.NewApricotFromConfigurationFile("apricot.toml"); err != nil {
			fmt.Errorf("error: %s", err)
		} else {
			if err := apricot.RunUp(); err != nil {
				fmt.Errorf("error: %s", err)
			} else {
				fmt.Println("Migration complete.")
			}
		}
	}
}

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
