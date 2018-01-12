package lib

import (
	"fmt"
	"strings"

	"github.com/fatih/color"
)

type Steps struct {
	green  func(a ...interface{}) string
	red    func(a ...interface{}) string
	yellow func(a ...interface{}) string
}

func NewSteps() Steps {
	green := color.New(color.FgGreen).SprintFunc()
	red := color.New(color.FgRed).SprintFunc()
	yellow := color.New(color.FgHiYellow).SprintFunc()
	return Steps{green, red, yellow}
}

func (s Steps) PrintHeader() {
	fmt.Printf("%-70s%7s\n", "STEP", "STATUS")
}

func (Steps) PrintRule() {
	fmt.Println(strings.Repeat("-", 77))
}

func (s Steps) Start(title string) {
	fmt.Printf("%-70s", title)
}

func (Steps) end(status string) {
	fmt.Printf("%s\n", status)
}

func (s Steps) Done() {
	s.end(fmt.Sprintf("%7s", "DONE"))
}

func (s Steps) Success() {
	s.end(s.green(fmt.Sprintf("%7s", "SUCCESS")))
}

func (s Steps) Fail() {
	s.end(s.yellow(fmt.Sprintf("%7s", "FAIL")))
}

func (s Steps) Error() {
	s.end(s.red(fmt.Sprintf("%7s", "ERROR")))
}
