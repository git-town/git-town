package gohacks_test

import (
	"testing"

	"github.com/git-town/git-town/v22/internal/gohacks"
	"github.com/shoenig/test/must"
)

func TestCounter(t *testing.T) {
	t.Parallel()
	counter := gohacks.Counter(0)
	counter.Increment()
	must.EqOp(t, 1, counter)
	counter.Increment()
	must.EqOp(t, 2, counter)
}
