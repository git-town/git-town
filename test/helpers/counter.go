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

// ToString provides a globally unique text each time it is called.
func (c *Counter) ToString() string {
	atomic.AddUint32(&c.value, 1)
	return strconv.Itoa(int(c.value))
}
