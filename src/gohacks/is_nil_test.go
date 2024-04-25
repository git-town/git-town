package gohacks_test

import (
	"testing"

	"github.com/git-town/git-town/v14/src/config/configdomain"
	"github.com/git-town/git-town/v14/src/gohacks"
	"github.com/shoenig/test/must"
)

func TestIsNil(t *testing.T) {
	t.Parallel()

	t.Run("direct nil given", func(t *testing.T) {
		t.Parallel()
		must.True(t, gohacks.IsNil(nil))
	})

	t.Run("pointer to nil interface", func(t *testing.T) {
		t.Parallel()
		var give *configdomain.GitHubToken
		must.True(t, gohacks.IsNil(give))
	})

	t.Run("non-nil given", func(t *testing.T) {
		t.Parallel()
		give := configdomain.NewGitHubTokenRef("foo")
		must.False(t, gohacks.IsNil(give))
	})
}
