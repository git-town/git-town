package cache_test

import (
	"testing"

	"github.com/git-town/git-town/v22/internal/gohacks/cache"
	"github.com/shoenig/test/must"
)

func TestStringSliceCache(t *testing.T) {
	t.Parallel()
	cache := cache.Cache[[]string]{}
	_, hasValue := cache.Get()
	must.False(t, hasValue)
	data := []string{"foo"}
	cache.Set(data)
	value, hasValue := cache.Get()
	must.True(t, hasValue)
	must.Eq(t, data, value)
}
