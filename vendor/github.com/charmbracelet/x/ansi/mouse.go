package ansi

import (
	"fmt"
)

// MouseX10 returns an escape sequence representing a mouse event in X10 mode.
// Note that this requires the terminal support X10 mouse modes.
//
//	CSI M Cb Cx Cy
//
// See: https://invisible-island.net/xterm/ctlseqs/ctlseqs.html#Mouse%20Tracking
func MouseX10(b byte, x, y int) string {
	const x10Offset = 32
	return "\x1b[M" + string(b+x10Offset) + string(byte(x)+x10Offset+1) + string(byte(y)+x10Offset+1)
}

// MouseSgr returns an escape sequence representing a mouse event in SGR mode.
//
//	CSI < Cb ; Cx ; Cy M
//	CSI < Cb ; Cx ; Cy m (release)
//
// See: https://invisible-island.net/xterm/ctlseqs/ctlseqs.html#Mouse%20Tracking
func MouseSgr(b byte, x, y int, release bool) string {
	s := "M"
	if release {
		s = "m"
	}
	if x < 0 {
		x = -x
	}
	if y < 0 {
		y = -y
	}
	return fmt.Sprintf("\x1b[<%d;%d;%d%s", b, x+1, y+1, s)
}
