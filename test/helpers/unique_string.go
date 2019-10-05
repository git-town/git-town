package helpers

import (
	"strconv"
	"sync/atomic"
)

// scenarioCounter counts the currently executed scenario
var scenarioCounter uint32

// UniqueString provides a globally unique number.
func UniqueString() string {
	return strconv.Itoa(int(atomic.AddUint32(&scenarioCounter, 1)))
}
