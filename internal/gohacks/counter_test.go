package gohacks_test

import (
	"testing"

	"github.com/git-town/git-town/v16/internal/gohacks"
	"github.com/shoenig/test/must"
)

func TestCounter(t *testing.T) {
	t.Parallel()
	counter := gohacks.Counter(0)
	counter.Inc()
	must.EqOp(t, 1, counter)
	counter.Inc()
	must.EqOp(t, 2, counter)
}
