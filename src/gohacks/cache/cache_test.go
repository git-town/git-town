package cache_test

import (
	"testing"

	"github.com/git-town/git-town/v14/src/gohacks/cache"
	"github.com/shoenig/test/must"
)

func TestBoolCache(t *testing.T) {
	t.Parallel()
	sc := cache.Bool{}
	must.False(t, sc.Initialized())
	sc.Set(true)
	must.True(t, sc.Initialized())
	must.True(t, sc.Value())
}

func TestStringCache(t *testing.T) {
	t.Parallel()
	sc := cache.String{}
	must.False(t, sc.Initialized())
	sc.Set("foo")
	must.True(t, sc.Initialized())
	must.EqOp(t, "foo", sc.Value())
}

func TestStringSliceCache(t *testing.T) {
	t.Parallel()
	ssc := cache.Strings{}
	must.False(t, ssc.Initialized())
	ssc.Set([]string{"foo"})
	must.True(t, ssc.Initialized())
	must.Eq(t, []string{"foo"}, ssc.Value())
}
