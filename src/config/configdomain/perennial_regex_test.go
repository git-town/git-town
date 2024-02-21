package configdomain_test

import (
	"testing"

	"github.com/git-town/git-town/v12/src/config/configdomain"
	"github.com/shoenig/test/must"
)

func TestPerennialRegex(t *testing.T) {
	t.Parallel()
	t.Run("empty regex matches nothing", func(t *testing.T) {
		t.Parallel()
		perennialRegex := configdomain.PerennialRegex("")
		must.False(t, perennialRegex.MatchBranch(""))
		must.False(t, perennialRegex.MatchBranch("foo"))
	})
	t.Run("only characters", func(t *testing.T) {
		t.Parallel()
	})
	t.Run("with wildcards", func(t *testing.T) {
		t.Parallel()
	})
}
