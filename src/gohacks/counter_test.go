package gohacks_test

import (
	"testing"

	"github.com/git-town/git-town/v14/src/gohacks"
	"github.com/shoenig/test/must"
)

func TestCounter(t *testing.T) {
	t.Parallel()
	counter := gohacks.NewCounter(0)
	counter.Inc()
	must.EqOp(t, 1, counter)
}
