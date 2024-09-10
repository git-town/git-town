package regexes_test

import (
	"testing"

	"github.com/git-town/git-town/v16/internal/cmd/regexes"
	"github.com/shoenig/test/must"
)

func TestRegexes(t *testing.T) {
	t.Parallel()
	t.Run("Matches", func(t *testing.T) {
		t.Parallel()
		t.Run("no regexes defined", func(t *testing.T) {
			t.Parallel()
			regexes, err := regexes.NewRegexes([]string{})
			must.NoError(t, err)
			must.True(t, regexes.Matches("foo"))
			must.True(t, regexes.Matches("bar"))
		})
	})
}
