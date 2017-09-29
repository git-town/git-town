package cfmt

import (
	"fmt"

	"github.com/fatih/color"
)

// Print prints the given text using fmt.Printf
// in a way where colors work on Windows
func Print(a ...interface{}) (int, error) {
	return fmt.Fprint(color.Output, a...)
}

// Printf prints the given text using fmt.Printf
// in a way where colors work on Windows
func Printf(format string, a ...interface{}) (int, error) {
	return fmt.Fprintf(color.Output, format, a...)
}

// Println prints the given text using fmt.Println
// in a way where colors work on Windows
func Println(a ...interface{}) (n int, err error) {
	return fmt.Fprintln(color.Output, a...)
}
