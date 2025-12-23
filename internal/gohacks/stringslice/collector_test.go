package stringslice_test

import (
	"testing"

	"github.com/git-town/git-town/v22/internal/gohacks/stringslice"
	"github.com/shoenig/test/must"
)

func TestCollector(t *testing.T) {
	t.Parallel()

	t.Run("AddF", func(t *testing.T) {
		t.Parallel()

		t.Run("adds formatted strings", func(t *testing.T) {
			t.Parallel()
			collector := stringslice.NewCollector()
			collector.AddF("Hello, %s!", "world")
			must.Eq(t, []string{"Hello, world!"}, collector.Result())
		})

		t.Run("multiple arguments and types", func(t *testing.T) {
			t.Parallel()
			collector := stringslice.NewCollector()
			collector.AddF("String: %s, Int: %d, Float: %.2f", "test", 42, 3.14159)
			must.Eq(t, []string{"String: test, Int: 42, Float: 3.14"}, collector.Result())
		})

		t.Run("empty format string", func(t *testing.T) {
			t.Parallel()
			collector := stringslice.NewCollector()
			collector.AddF("")
			must.Eq(t, []string{""}, collector.Result())
		})

		t.Run("no arguments", func(t *testing.T) {
			t.Parallel()
			collector := stringslice.NewCollector()
			collector.AddF("plain text")
			must.Eq(t, []string{"plain text"}, collector.Result())
		})

		t.Run("chaining with Add", func(t *testing.T) {
			t.Parallel()
			collector := stringslice.NewCollector()
			collector.Add("first")
			collector.AddF("second: %s", "value")
			collector.Add("third")
			must.Eq(t, []string{"first", "second: value", "third"}, collector.Result())
		})
	})

	t.Run("owned variable", func(t *testing.T) {
		t.Parallel()
		collector := stringslice.NewCollector()
		must.Len(t, 0, collector.Result())
		collector.Add("one")
		collector.Add("two")
		must.Eq(t, []string{"one", "two"}, collector.Result())
	})

	t.Run("works with pass by reference", func(t *testing.T) {
		t.Parallel()
		collector := stringslice.NewCollector()
		passByReference(&collector)
		must.Eq(t, []string{"external"}, collector.Result())
	})

	t.Run("works with pass by value", func(t *testing.T) {
		t.Parallel()
		collector := stringslice.NewCollector()
		passByValue(collector)
		must.Eq(t, []string{"external"}, collector.Result())
	})
}

func passByReference(collector *stringslice.Collector) {
	collector.Add("external")
}

func passByValue(collector stringslice.Collector) {
	collector.Add("external")
}
