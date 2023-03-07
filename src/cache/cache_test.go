package cache_test

import (
	"testing"

	"github.com/git-town/git-town/v7/src/cache"
	"github.com/stretchr/testify/assert"
)

func TestBoolCache(t *testing.T) {
	t.Parallel()
	sc := cache.Cache[bool]{}
	assert.False(t, sc.Initialized())
	sc.Set(true)
	assert.True(t, sc.Initialized())
	assert.True(t, sc.Value())
}

func TestStringCache(t *testing.T) {
	t.Parallel()
	sc := cache.Cache[string]{}
	assert.False(t, sc.Initialized())
	sc.Set("foo")
	assert.True(t, sc.Initialized())
	assert.Equal(t, "foo", sc.Value())
}

func TestStringSliceCache(t *testing.T) {
	t.Parallel()
	ssc := cache.Cache[[]string]{}
	assert.False(t, ssc.Initialized())
	ssc.Set([]string{"foo"})
	assert.True(t, ssc.Initialized())
	assert.Equal(t, []string{"foo"}, ssc.Value())
}
