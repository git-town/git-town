package cache_test

import (
	"testing"

	"github.com/git-town/git-town/v21/internal/gohacks/cache"
	"github.com/shoenig/test/must"
)

func TestStringSliceCache(t *testing.T) {
	t.Parallel()
	cache := cache.Cache[[]string]{}
	must.True(t, cache.Value().IsNone())
	data := []string{"foo"}
	cache.Set(data)
	value, hasValue := cache.Value().Get()
	must.True(t, hasValue)
	must.Eq(t, data, value)
}
