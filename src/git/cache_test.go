package git_test

import (
	"testing"

	"github.com/git-town/git-town/src/git"
	"github.com/stretchr/testify/assert"
)

func TestStringCache(t *testing.T) {
	sc := git.StringCache{}
	assert.False(t, sc.Initialized())
	sc.Set("foo")
	assert.True(t, sc.Initialized())
	assert.Equal(t, "foo", sc.Value())
}

func TestStringSliceCache(t *testing.T) {
	ssc := git.StringSliceCache{}
	assert.False(t, ssc.Initialized())
	ssc.Set([]string{"foo"})
	assert.True(t, ssc.Initialized())
	assert.Equal(t, []string{"foo"}, ssc.Value())
}
