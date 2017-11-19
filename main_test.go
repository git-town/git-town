package main

// This file is mandatory as otherwise the filebeat.test binary is not generated correctly.
import (
	"testing"
)

// Test started when the test binary is started. Only calls main.
func TestSystem(t *testing.T) {
	main()
}
