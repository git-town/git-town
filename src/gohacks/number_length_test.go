package gohacks_test

import (
	"testing"

	"github.com/git-town/git-town/v13/src/gohacks"
	"github.com/shoenig/test/must"
)

func TestNumberLength(t *testing.T) {
	t.Parallel()
	tests := map[int]int{
		0:    1,
		1:    1,
		9:    1,
		10:   2,
		99:   2,
		100:  3,
		-1:   2,
		-9:   2,
		-10:  3,
		-99:  3,
		-100: 4,
	}
	for give, want := range tests {
		have := gohacks.NumberLength(give)
		must.EqOp(t, want, have)
	}
}
