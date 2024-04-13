package stringslice_test

import (
	"testing"

	"github.com/git-town/git-town/v14/src/gohacks/stringslice"
	"github.com/shoenig/test/must"
)

func TestCollector(t *testing.T) {
	t.Parallel()
	collector := stringslice.Collector{}
	must.Eq(t, []string{}, collector.Result())
	collector.Add("one")
	collector.Add("two")
	must.Eq(t, []string{"one", "two"}, collector.Result())
}
