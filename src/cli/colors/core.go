// Package colors provides definitions for ANSI colors to be used in terminal output.
package colors

import "github.com/muesli/termenv"

func Bold() termenv.Style {
	return termenv.String().Bold()
}

func Green() termenv.Style {
	return termenv.String().Foreground(termenv.ANSIGreen)
}

func Red() termenv.Style {
	return termenv.String().Foreground(termenv.ANSIRed)
}
