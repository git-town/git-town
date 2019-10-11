package command

import (
	"strings"

	"github.com/Originate/exit"
	"github.com/fatih/color"
)

var debug bool
var count int

// SetDebug sets whether or not we are in debug mode
func SetDebug(value bool) {
	debug = value
}

func logRun(argv ...string) {
	if debug {
		count++
		_, err := color.New(color.FgBlue).Printf("DEBUG (%d): %s\n", count, strings.Join(argv, " "))
		exit.If(err)
	}
}
