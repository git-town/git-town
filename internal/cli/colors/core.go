// Package colors provides definitions for ANSI colors to be used in terminal output.
package colors

import "github.com/muesli/termenv"

func Bold() termenv.Style {
	return termenv.String().Bold()
}

func BoldCyan() termenv.Style {
	return Cyan().Bold()
}

func BoldGreen() termenv.Style {
	return Green().Bold()
}

func BoldRed() termenv.Style {
	return Red().Bold()
}

func BoldUnderline() termenv.Style {
	return Bold().Underline()
}

func Cyan() termenv.Style {
	return termenv.String().Foreground(termenv.ANSICyan)
}

func Faint() termenv.Style {
	return termenv.String().Faint()
}

func FaintBold() termenv.Style {
	return termenv.String().Faint().Bold()
}

func Green() termenv.Style {
	return termenv.String().Foreground(termenv.ANSIGreen)
}

func Red() termenv.Style {
	return termenv.String().Foreground(termenv.ANSIRed)
}
