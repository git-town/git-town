package utils

import (
	"strings"
	"time"
)

// S repeats a space n times
func S(n int) string {
	if n < 0 {
		n = 1
	}
	return strings.Repeat(" ", n)
}

// TimeNowFunc is a utility function to simply testing
// by allowing TimeNowFunc to be defined to zero time
// to remove the time domain from tests
var TimeNowFunc = func() time.Time {
	return time.Now()
}
