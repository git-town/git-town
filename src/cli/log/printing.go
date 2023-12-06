package log

import (
	"fmt"

	"github.com/fatih/color"
)

// The Printing logger logs activities of a particular component on the CLI.
// TODO: move this to the print package?
type Printing struct{}

func (self Printing) Failed(failure error) {
	_, err := color.New(color.Bold, color.FgRed).Printf("FAILED: %v\n", failure)
	if err != nil {
		fmt.Printf("FAILED: %v\n", err)
	}
}

func (self Printing) Start(template string, data ...interface{}) {
	fmt.Println()
	_, err := color.New(color.Bold).Printf(template, data...)
	if err != nil {
		fmt.Printf(template, data...)
	}
}

func (self Printing) Success() {
	_, err := color.New(color.Bold, color.FgGreen).Printf("ok\n")
	if err != nil {
		fmt.Println("ok")
	}
}
