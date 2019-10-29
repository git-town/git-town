package cfmt

import (
	"fmt"

	"github.com/fatih/color"
)

// Print prints the given text using fmt.Print
// in a way where colors work on Windows
func Print(a ...interface{}) {
	_, err := fmt.Fprint(color.Output, a...)
	if err != nil {
		panic(err)
	}
}

// Printf prints the given text using fmt.Printf
// in a way where colors work on Windows
func Printf(format string, a ...interface{}) {
	_, err := fmt.Fprintf(color.Output, format, a...)
	if err != nil {
		panic(err)
	}
}

// Println prints the given text using fmt.Println
// in a way where colors work on Windows
func Println(a ...interface{}) {
	_, err := fmt.Fprintln(color.Output, a...)
	if err != nil {
		panic(err)
	}
}
