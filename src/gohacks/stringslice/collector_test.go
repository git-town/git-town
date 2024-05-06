package stringslice_test

import (
	"testing"

	"github.com/git-town/git-town/v14/src/gohacks/stringslice"
	"github.com/shoenig/test/must"
)

func TestCollector(t *testing.T) {
	t.Parallel()
	t.Run("owned variable", func(t *testing.T) {
		collector := stringslice.NewCollector()
		must.Eq(t, []string{}, collector.Result())
		collector.Add("one")
		collector.Add("two")
		must.Eq(t, []string{"one", "two"}, collector.Result())
	})
	t.Run("works with pass by value", func(t *testing.T) {
		t.Parallel()
		collector := stringslice.NewCollector()
		passByValue(collector)
		must.Eq(t, []string{"external"}, collector.Result())
	})
	t.Run("works with pass by reference", func(t *testing.T) {
		t.Parallel()
		collector := stringslice.NewCollector()
		passByReference(&collector)
		must.Eq(t, []string{"external"}, collector.Result())
	})
}

func passByReference(collector *stringslice.Collector) {
	collector.Add("external")
}

func passByValue(collector stringslice.Collector) {
	collector.Add("external")
}
