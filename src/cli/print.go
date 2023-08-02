package cli

import (
	"fmt"

	"github.com/fatih/color"
)

// BoolSetting provides a human-readable serialization for bool values.
func BoolSetting(value bool) string {
	if value {
		return "yes"
	}
	return "no"
}

// Printf prints the given text using fmt.Printf
// in a way where colors work on Windows.
func Printf(format string, a ...interface{}) {
	_, err := fmt.Fprintf(color.Output, format, a...)
	if err != nil {
		fmt.Printf(format, a...)
	}
}

// Println prints the given text using fmt.Println
// in a way where colors work on Windows.
func Println(a ...interface{}) {
	_, err := fmt.Fprintln(color.Output, a...)
	if err != nil {
		fmt.Println(a...)
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

func PrintEntry(label, value string) {
	Printf("  %s: %s\n", label, value)
}

// PrintError prints the given error message to the console.
func PrintError(err error) {
	PrintlnColor(color.New(color.Bold).Add(color.FgRed), "\nError:", err.Error(), "\n")
}

func PrintHeader(text string) {
	boldUnderline := color.New(color.Bold).Add(color.Underline)
	PrintlnColor(boldUnderline, text+":")
}

// PrintLabelAndValue prints the label bolded and underlined
// the value indented on the next line
// followed by an empty line.
func PrintLabelAndValue(label, value string) {
	PrintHeader(label)
	Println(Indent(value))
	fmt.Println()
}

func StringSetting(text string) string {
	if text != "" {
		return text
	}
	return "(not set)"
}

// PrintingLog logs activities of a particular component on the CLI.
type PrintingLog struct{}

func (l PrintingLog) Start(template string, messages ...interface{}) {
	fmt.Println()
	_, err := color.New(color.Bold).Printf(template, messages...)
	if err != nil {
		fmt.Printf(template, messages...)
	}
}

func (l PrintingLog) Success() {
	_, err := color.New(color.Bold, color.FgGreen).Printf("ok\n")
	if err != nil {
		fmt.Println("ok")
	}
}

func (l PrintingLog) Failed(failure error) {
	_, err := color.New(color.Bold, color.FgRed).Printf("FAILED: %v\n", failure)
	if err != nil {
		fmt.Printf("FAILED: %v\n", err)
	}
}

type SilentLog struct{}

func (p SilentLog) Start(string, ...interface{}) {}
func (p SilentLog) Success()                     {}
func (p SilentLog) Failed(error)                 {}
