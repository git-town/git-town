package helpers

import (
	"testing"

	"github.com/davecgh/go-spew/spew"
)

func TestInspectVariable(t *testing.T) {
	t.Parallel()
	// This test is necessary to keep "github.com/davecgh/go-spew/spew" in the go.mod file.
	a := 1
	spew.Dump(a)
}
