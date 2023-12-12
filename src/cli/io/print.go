package io

import (
	"fmt"

	"github.com/fatih/color"
)

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
