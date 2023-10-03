package cache_test

import (
	"testing"

	"github.com/git-town/git-town/v9/src/cache"
	"github.com/shoenig/test"
	"github.com/stretchr/testify/assert"
)

func TestBoolCache(t *testing.T) {
	t.Parallel()
	sc := cache.Bool{}
	test.False(t, sc.Initialized())
	sc.Set(true)
	test.True(t, sc.Initialized())
	test.True(t, sc.Value())
}

func TestStringCache(t *testing.T) {
	t.Parallel()
	sc := cache.String{}
	test.False(t, sc.Initialized())
	sc.Set("foo")
	test.True(t, sc.Initialized())
	test.EqOp(t, "foo", sc.Value())
}

func TestStringSliceCache(t *testing.T) {
	t.Parallel()
	ssc := cache.Strings{}
	test.False(t, ssc.Initialized())
	ssc.Set([]string{"foo"})
	test.True(t, ssc.Initialized())
	assert.Equal(t, []string{"foo"}, ssc.Value())
}
