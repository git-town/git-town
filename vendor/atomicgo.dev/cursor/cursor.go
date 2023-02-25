//go:build !windows
// +build !windows

package cursor

import (
	"fmt"
	"os"
)

var target Writer = os.Stdout

// SetTarget allows for any arbitrary io.Writer to be used
// for cursor movement (will not work on Windows).
func SetTarget(w Writer) {
	target = w
}

// Up moves the cursor n lines up relative to the current position.
func Up(n int) {
	fmt.Fprintf(target, "\x1b[%dA", n)
	height += n
}

// Down moves the cursor n lines down relative to the current position.
func Down(n int) {
	fmt.Fprintf(target, "\x1b[%dB", n)
	if height-n <= 0 {
		height = 0
	} else {
		height -= n
	}
}

// Right moves the cursor n characters to the right relative to the current position.
func Right(n int) {
	fmt.Fprintf(target, "\x1b[%dC", n)
}

// Left moves the cursor n characters to the left relative to the current position.
func Left(n int) {
	fmt.Fprintf(target, "\x1b[%dD", n)
}

// HorizontalAbsolute moves the cursor to n horizontally.
// The position n is absolute to the start of the line.
func HorizontalAbsolute(n int) {
	n += 1 // Moves the line to the character after n
	fmt.Fprintf(target, "\x1b[%dG", n)
}

// Show the cursor if it was hidden previously.
// Don't forget to show the cursor at least at the end of your application.
// Otherwise the user might have a terminal with a permanently hidden cursor, until they reopen the terminal.
func Show() {
	fmt.Fprint(target, "\x1b[?25h")
}

// Hide the cursor.
// Don't forget to show the cursor at least at the end of your application with Show.
// Otherwise the user might have a terminal with a permanently hidden cursor, until they reopen the terminal.
func Hide() {
	fmt.Fprintf(target, "\x1b[?25l")
}

// ClearLine clears the current line and moves the cursor to it's start position.
func ClearLine() {
	fmt.Fprintf(target, "\x1b[2K")
}
