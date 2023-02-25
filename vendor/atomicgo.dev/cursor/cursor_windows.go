package cursor

import (
	"os"
	"syscall"
	"unsafe"
)

var target Writer = os.Stdout

// SetTarget allows for any arbitrary Writer to be used
func SetTarget(w Writer) {
	target = w
}

// Up moves the cursor n lines up relative to the current position.
func Up(n int) {
	move(0, -n)
	height += n
}

// Down moves the cursor n lines down relative to the current position.
func Down(n int) {
	move(0, n)
	if height-n <= 0 {
		height = 0
	} else {
		height -= n
	}
}

// Right moves the cursor n characters to the right relative to the current position.
func Right(n int) {
	move(n, 0)
}

// Left moves the cursor n characters to the left relative to the current position.
func Left(n int) {
	move(-n, 0)
}

func move(x int, y int) {
	handle := syscall.Handle(target.Fd())

	var csbi consoleScreenBufferInfo
	_, _, _ = procGetConsoleScreenBufferInfo.Call(uintptr(handle), uintptr(unsafe.Pointer(&csbi)))

	var cursor coord
	cursor.x = csbi.cursorPosition.x + short(x)
	cursor.y = csbi.cursorPosition.y + short(y)

	_, _, _ = procSetConsoleCursorPosition.Call(uintptr(handle), uintptr(*(*int32)(unsafe.Pointer(&cursor))))
}

// HorizontalAbsolute moves the cursor to n horizontally.
// The position n is absolute to the start of the line.
func HorizontalAbsolute(n int) {
	handle := syscall.Handle(target.Fd())

	var csbi consoleScreenBufferInfo
	_, _, _ = procGetConsoleScreenBufferInfo.Call(uintptr(handle), uintptr(unsafe.Pointer(&csbi)))

	var cursor coord
	cursor.x = short(n)
	cursor.y = csbi.cursorPosition.y

	if csbi.size.x < cursor.x {
		cursor.x = csbi.size.x
	}

	_, _, _ = procSetConsoleCursorPosition.Call(uintptr(handle), uintptr(*(*int32)(unsafe.Pointer(&cursor))))
}

// Show the cursor if it was hidden previously.
// Don't forget to show the cursor at least at the end of your application.
// Otherwise the user might have a terminal with a permanently hidden cursor, until he reopens the terminal.
func Show() {
	handle := syscall.Handle(target.Fd())

	var cci consoleCursorInfo
	_, _, _ = procGetConsoleCursorInfo.Call(uintptr(handle), uintptr(unsafe.Pointer(&cci)))
	cci.visible = 1

	_, _, _ = procSetConsoleCursorInfo.Call(uintptr(handle), uintptr(unsafe.Pointer(&cci)))
}

// Hide the cursor.
// Don't forget to show the cursor at least at the end of your application with Show.
// Otherwise the user might have a terminal with a permanently hidden cursor, until he reopens the terminal.
func Hide() {
	handle := syscall.Handle(target.Fd())

	var cci consoleCursorInfo
	_, _, _ = procGetConsoleCursorInfo.Call(uintptr(handle), uintptr(unsafe.Pointer(&cci)))
	cci.visible = 0

	_, _, _ = procSetConsoleCursorInfo.Call(uintptr(handle), uintptr(unsafe.Pointer(&cci)))
}

// ClearLine clears the current line and moves the cursor to it's start position.
func ClearLine() {
	handle := syscall.Handle(target.Fd())

	var csbi consoleScreenBufferInfo
	_, _, _ = procGetConsoleScreenBufferInfo.Call(uintptr(handle), uintptr(unsafe.Pointer(&csbi)))

	var w uint32
	var x short
	cursor := csbi.cursorPosition
	x = csbi.size.x
	_, _, _ = procFillConsoleOutputCharacter.Call(uintptr(handle), uintptr(' '), uintptr(x), uintptr(*(*int32)(unsafe.Pointer(&cursor))), uintptr(unsafe.Pointer(&w)))
}
