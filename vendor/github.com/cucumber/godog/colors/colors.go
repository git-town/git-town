package colors

import (
	"fmt"
	"strings"
)

const ansiEscape = "\x1b"

// a color code type
type color int

// some ansi colors
const (
	black color = iota + 30
	red
	green
	yellow
	blue    // unused
	magenta // unused
	cyan
	white
)

func colorize(s interface{}, c color) string {
	return fmt.Sprintf("%s[%dm%v%s[0m", ansiEscape, c, s, ansiEscape)
}

// ColorFunc is a helper type to create colorized strings.
type ColorFunc func(interface{}) string

// Bold will accept a ColorFunc and return a new ColorFunc
// that will make the string bold.
func Bold(fn ColorFunc) ColorFunc {
	return ColorFunc(func(input interface{}) string {
		return strings.Replace(fn(input), ansiEscape+"[", ansiEscape+"[1;", 1)
	})
}

// Green will accept an interface and return a colorized green string.
func Green(s interface{}) string {
	return colorize(s, green)
}

// Red will accept an interface and return a colorized green string.
func Red(s interface{}) string {
	return colorize(s, red)
}

// Cyan will accept an interface and return a colorized green string.
func Cyan(s interface{}) string {
	return colorize(s, cyan)
}

// Black will accept an interface and return a colorized green string.
func Black(s interface{}) string {
	return colorize(s, black)
}

// Yellow will accept an interface and return a colorized green string.
func Yellow(s interface{}) string {
	return colorize(s, yellow)
}

// White will accept an interface and return a colorized green string.
func White(s interface{}) string {
	return colorize(s, white)
}
