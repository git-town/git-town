package gohacks_test

import (
	"testing"

	"github.com/git-town/git-town/v14/src/gohacks"
	"github.com/shoenig/test/must"
)

func TestCounter(t *testing.T) {
	t.Parallel()

	t.Run("owned variable", func(t *testing.T) {
		t.Parallel()
		counter := gohacks.NewCounter()
		counter.Register()
		counter.Register()
		must.Eq(t, 2, counter.Count())
	})

	t.Run("pass by reference works", func(t *testing.T) {
		t.Parallel()
		counter := gohacks.NewCounter()
		passByReference(&counter)
		must.Eq(t, 2, counter.Count())
	})

	t.Run("pass by value works", func(t *testing.T) {
		t.Parallel()
		counter := gohacks.NewCounter()
		passByValue(counter)
		must.Eq(t, 2, counter.Count())
	})
}

func passByReference(counter *gohacks.Counter) {
	counter.Register()
	counter.Register()
}

func passByValue(counter gohacks.Counter) {
	counter.Register()
	counter.Register()
}
