package cursor

import "io"

var height int

// Bottom moves the cursor to the bottom of the terminal.
// This is done by calculating how many lines were moved by Up and Down.
func Bottom() {
	if height > 0 {
		Down(height)
		StartOfLine()
		height = 0
	}
}

// StartOfLine moves the cursor to the start of the current line.
func StartOfLine() {
	HorizontalAbsolute(0)
}

// StartOfLineDown moves the cursor down by n lines, then moves to cursor to the start of the line.
func StartOfLineDown(n int) {
	Down(n)
	StartOfLine()
}

// StartOfLineUp moves the cursor up by n lines, then moves to cursor to the start of the line.
func StartOfLineUp(n int) {
	Up(n)
	StartOfLine()
}

// UpAndClear moves the cursor up by n lines, then clears the line.
func UpAndClear(n int) {
	Up(n)
	ClearLine()
}

// DownAndClear moves the cursor down by n lines, then clears the line.
func DownAndClear(n int) {
	Down(n)
	ClearLine()
}

// Move moves the cursor relative by x and y.
func Move(x, y int) {
	if x > 0 {
		Right(x)
	} else if x < 0 {
		x *= -1
		Left(x)
	}

	if y > 0 {
		Up(y)
	} else if y < 0 {
		y *= -1
		Down(y)
	}
}

// ClearLinesUp clears n lines upwards from the current position and moves the cursor.
func ClearLinesUp(n int) {
	for i := 0; i < n; i++ {
		UpAndClear(1)
	}
}

// ClearLinesDown clears n lines downwards from the current position and moves the cursor.
func ClearLinesDown(n int) {
	for i := 0; i < n; i++ {
		DownAndClear(1)
	}
}

// Writer is an expanded io.Writer interface with a file descriptor.
type Writer interface {
	io.Writer
	Fd() uintptr
}
