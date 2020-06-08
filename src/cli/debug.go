package cli

import (
	"strings"

	"github.com/fatih/color"
)

var debug bool
var count int

// SetDebug sets whether or not we are in debug mode.
func SetDebug(value bool) {
	debug = value
}

// LogRun debugs the given executed command on the command line.
func LogRun(cmd string, args ...string) {
	if debug {
		count++
		_, err := color.New(color.FgBlue).Printf("DEBUG (%d): %s %s\n", count, cmd, strings.Join(args, " "))
		if err != nil {
			panic(err)
		}
	}
}
