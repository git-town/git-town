package helpers

import (
	"strconv"
	"sync/atomic"
)

// AtomicCounter provides unique string segments in a thread-safe way.
// The zero value is an initialized instance.
type AtomicCounter struct {
	// value counts the currently executed scenario.
	value uint32
}

// NextAsString provides a globally unique text each time it is called.
func (self *AtomicCounter) NextAsString() string {
	atomic.AddUint32(&self.value, 1)
	return strconv.FormatUint(uint64(self.value), 10)
}
