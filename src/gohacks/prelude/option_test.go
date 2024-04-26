package prelude_test

import (
	"testing"

	"github.com/git-town/git-town/v14/src/config/configdomain"
	. "github.com/git-town/git-town/v14/src/gohacks/prelude"
	"github.com/shoenig/test/must"
)

func TestOption(t *testing.T) {
	t.Parallel()
	t.Run("StringOr", func(t *testing.T) {
		t.Parallel()
		t.Run("Some(struct that implements fmt.Stringer)", func(t *testing.T) {
			t.Parallel()
			text := "my token"
			option := Some(configdomain.GitHubToken(text))
			have := option.String()
			must.EqOp(t, text, have)
		})
		t.Run("Some(struct that doesn't implement fmt.Stringer)", func(t *testing.T) {
			t.Parallel()
			type test struct {
				data bool
			}
			instance := test{
				data: true,
			}
			option := Some(instance)
			have := option.String()
			want := "&{true}"
			must.EqOp(t, want, have)
		})
		t.Run("None[int]", func(t *testing.T) {
			t.Parallel()
			option := None[int]()
			have := option.String()
			must.EqOp(t, "", have)
		})
		t.Run("None[*int]", func(t *testing.T) {
			t.Parallel()
			option := None[*int]()
			have := option.String()
			must.EqOp(t, "", have)
		})
		t.Run("None[string newtype]", func(t *testing.T) {
			t.Parallel()
			option := None[configdomain.PerennialRegex]()
			have := option.String()
			must.EqOp(t, "", have)
		})
	})
}
