package lib

import (
	"fmt"
	"strings"

	"github.com/fatih/color"
)

// Steps helps output nice messages
type Steps struct {
	green  func(a ...interface{}) string
	red    func(a ...interface{}) string
	yellow func(a ...interface{}) string
}

// NewSteps builds a new Steps instance
func NewSteps() Steps {
	green := color.New(color.FgGreen).SprintFunc()
	red := color.New(color.FgRed).SprintFunc()
	yellow := color.New(color.FgHiYellow).SprintFunc()
	return Steps{green, red, yellow}
}

// PrintHeader outputs the table header
func (s Steps) PrintHeader() {
	fmt.Printf("\n%-70s%7s\n", "STEP", "STATUS")
}

// PrintRule outputs a horizontal rule
func (Steps) PrintRule() {
	fmt.Println(strings.Repeat("-", 77))
}

// Start outputs the start of an entry in the table
func (s Steps) Start(title string) {
	fmt.Printf("%-70s", title)
}

func (Steps) end(status string) {
	fmt.Printf("%s\n", status)
}

// Done outputs the end of an entry in the table with a status of DONE
func (s Steps) Done() {
	s.end(fmt.Sprintf("%7s", "DONE"))
}

// Success outputs the end of an entry in the table with a status of SUCCESS
func (s Steps) Success() {
	s.end(s.green(fmt.Sprintf("%7s", "SUCCESS")))
}

// Fail outputs the end of an entry in the table with a status of FAIL
func (s Steps) Fail() {
	s.end(s.yellow(fmt.Sprintf("%7s", "FAIL")))
}

// Error outputs the end of an entry in the table with a status of ERROR
func (s Steps) Error() {
	s.end(s.red(fmt.Sprintf("%7s", "ERROR")))
}
