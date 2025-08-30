package helpers_test

import (
	"testing"

	"github.com/kr/pretty"
)

// this function exists only to use pretty somewhere so that the half-intelligent Go module system
// doesn't auto-remove it from go.mod.
func TestPretty(t *testing.T) {
	t.Parallel()
	a := 1
	pretty.Ldiff(t, a, a)
}
