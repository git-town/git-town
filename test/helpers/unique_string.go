package helpers

import (
	"strconv"
	"sync/atomic"
)

// Counter provides unique string segments.
// The zero value is an initialized instance.
type Counter struct {
	// value counts the currently executed scenario.
	value uint32
}

// UniqueString provides a globally unique number.
func (us *Counter) ToString() string {
	return strconv.Itoa(int(atomic.AddUint32(&us.value, 1)))
}
