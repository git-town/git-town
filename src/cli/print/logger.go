package print

import (
	"fmt"

	"github.com/fatih/color"
)

// The Logger logger logs activities of a particular component on the CLI.
type Logger struct{}

func (l Logger) Failed(failure error) {
	// TODO: use termenv colors here
	_, err := color.New(color.Bold, color.FgRed).Printf("FAILED: %v\n", failure)
	if err != nil {
		fmt.Printf("FAILED: %v\n", err)
	}
}

func (l Logger) Start(template string, data ...interface{}) {
	fmt.Println()
	_, err := color.New(color.Bold).Printf(template, data...)
	if err != nil {
		fmt.Printf(template, data...)
	}
}

func (l Logger) Success() {
	_, err := color.New(color.Bold, color.FgGreen).Printf("ok\n")
	if err != nil {
		fmt.Println("ok")
	}
}

// The silent logger acts as a stand-in for loggers when no logging is desired.
type NoLogger struct{}

func (n NoLogger) Failed(error)                 {}
func (n NoLogger) Start(string, ...interface{}) {}
func (n NoLogger) Success()                     {}
