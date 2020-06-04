package helpers

import (
	"fmt"

	"github.com/fatih/color"
)

// PrintLog prints the given log message in bold.
func PrintLog(message string) {
	fmt.Println()
	_, err := color.New(color.Bold).Println(message)
	if err != nil {
		panic(err)
	}
}
