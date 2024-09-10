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
		t.Run("single regex", func(t *testing.T) {
			t.Parallel()
			regexes, err := regexes.NewRegexes([]string{"^kg-"})
			must.NoError(t, err)
			must.True(t, regexes.Matches("kg-one"))
			must.True(t, regexes.Matches("kg-two"))
			must.False(t, regexes.Matches("other"))
		})
		t.Run("multiple regexes", func(t *testing.T) {
			t.Parallel()
			regexes, err := regexes.NewRegexes([]string{"^kg-", "main"})
			must.NoError(t, err)
			must.True(t, regexes.Matches("kg-one"))
			must.True(t, regexes.Matches("kg-two"))
			must.True(t, regexes.Matches("main"))
			must.False(t, regexes.Matches("other"))
		})
	})
}
