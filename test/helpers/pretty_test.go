package helpers_test

import (
	"testing"

	"github.com/kr/pretty"
)

// this function exists only to use pretty somewhere so that the half-intelligent Go module system
// doesn't auto-remove it from go.mod.
func TestPretty(t *testing.T) {
	a := 1
	b := 2
	pretty.Ldiff(t, a, b)
}
