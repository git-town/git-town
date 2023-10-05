package gohacks_test

import (
	"testing"

	"github.com/git-town/git-town/v9/src/gohacks"
	"github.com/shoenig/test/must"
)

func TestCounter(t *testing.T) {
	t.Parallel()
	counter := gohacks.Counter{}
	counter.RegisterRun()
	counter.RegisterRun()
	must.Eq(t, 2, counter.Count())
}
