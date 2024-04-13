package gohacks_test

import (
	"testing"

	"github.com/git-town/git-town/v14/src/gohacks"
	"github.com/shoenig/test/must"
)

func TestCounter(t *testing.T) {
	t.Parallel()
	counter := gohacks.Counter{}
	counter.Register()
	counter.Register()
	must.Eq(t, 2, counter.Count())
}
