package cache_test

import (
	"testing"

	"github.com/git-town/git-town/v17/internal/gohacks/cache"
	"github.com/shoenig/test/must"
)

func TestStringSliceCache(t *testing.T) {
	t.Parallel()
	ssc := cache.Cache[[]string]{}
	must.False(t, ssc.Initialized())
	ssc.Set(&[]string{"foo"})
	must.True(t, ssc.Initialized())
	must.Eq(t, []string{"foo"}, *ssc.Value())
}
