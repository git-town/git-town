package command

import (
	"strings"

	"github.com/Originate/exit"
	"github.com/fatih/color"
)

var debug bool
var debugCount int
var debugFmt *color.Color

// SetDebug sets whether or not we are in debug mode
func SetDebug(value bool) {
	debug = value
}

func logRun(c *Command) {
	if debug {
		debugCount++
		_, err := debugFmt.Printf("DEBUG (%d): %s\n", debugCount, strings.Join(append([]string{c.name}, c.args...), " "))
		exit.If(err)
	}
}

func init() {
	debugFmt = color.New(color.FgBlue)
}
