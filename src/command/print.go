package command

import (
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/git-town/git-town/src/util"
)

// ExitWithErrorMessage prints the given error message and terminates the application.
func ExitWithErrorMessage(messages ...string) {
	PrintError(messages...)
	os.Exit(1)
}

// Printf prints the given text using fmt.Printf
// in a way where colors work on Windows.
func Printf(format string, a ...interface{}) {
	_, err := fmt.Fprintf(color.Output, format, a...)
	if err != nil {
		panic(err)
	}
}

// Println prints the given text using fmt.Println
// in a way where colors work on Windows.
func Println(a ...interface{}) {
	_, err := fmt.Fprintln(color.Output, a...)
	if err != nil {
		panic(err)
	}
}

// PrintlnColor prints using the given color function.
// If that doesn't work, it falls back to printing without color.
func PrintlnColor(color *color.Color, messages ...interface{}) {
	_, err := color.Println(messages...)
	if err != nil {
		fmt.Println(messages...)
	}
}

// PrintError prints the given error message to the console.
func PrintError(messages ...string) {
	errHeaderFmt := color.New(color.Bold).Add(color.FgRed)
	errMessageFmt := color.New(color.FgRed)
	fmt.Println()
	PrintlnColor(errHeaderFmt, "  Error")
	for _, message := range messages {
		PrintlnColor(errMessageFmt, "  "+message)
	}
	fmt.Println()
}

// PrintLabelAndValue prints the label bolded and underlined
// the value indented on the next line
// followed by an empty line.
func PrintLabelAndValue(label, value string) {
	labelFmt := color.New(color.Bold).Add(color.Underline)
	PrintlnColor(labelFmt, label+":")
	Println(util.Indent(value, 1))
	fmt.Println()
}
