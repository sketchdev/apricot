package main

import (
	"fmt"
	"os"

	"github.com/sketchdev/apricot/commands"
	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	app     = kingpin.New("apricot", "A cross platform, cross language migration tool.")
	initCmd = app.Command("init", "Initializes apricot in the current working directory.")
)

func main() {
	switch kingpin.MustParse(app.Parse(os.Args[1:])) {
	case initCmd.FullCommand():
		if err := commands.RunInitialize(); err != nil {
			fmt.Println(err)
		} else {
			fmt.Println("Initialization complete!!")
		}
	}
}
