package main

import (
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
				fmt.Printf("An error occurred while migrating: \n%s\n", err)
			} else {
				fmt.Println("Migration complete.")
			}
		}
	}
}

// RunInitialize prepares a directory for apricot use
func RunInitialize() error {
	// create migrations folder if needed
	if err := os.Mkdir("migrations", os.ModeDir|os.ModePerm); err != nil {
		return fmt.Errorf("failed to create migrations directory: %v", err)
	}
	// create apricot.toml file
	contents := []byte("engine = \"postgres\"\nmigrations = [\"migrations/current\"]\n")
	if err := ioutil.WriteFile("migrations/apricot.toml", contents, os.ModePerm); err != nil {
		return fmt.Errorf("failed to create the apricot.toml file: %v", err)
	}
	// create current directory
	if err := os.Mkdir(path.Join("migrations", "current"), os.ModeDir|os.ModePerm); err != nil {
		return fmt.Errorf("failed to create migrations/current directory: %v", err)
	}
	return nil
}
