package prelude_test

import (
	"testing"

	"github.com/git-town/git-town/v14/src/config/configdomain"
	. "github.com/git-town/git-town/v14/src/gohacks/prelude"
	"github.com/shoenig/test/must"
)

func TestOption(t *testing.T) {
	t.Parallel()

	t.Run("comparison", func(t *testing.T) {
		t.Parallel()
		t.Run("int", func(t *testing.T) {
			t.Parallel()
			t.Run("equal", func(t *testing.T) {
				a := Some(1)
				b := Some(1)
				must.True(t, a.Equal(b))
			})
			t.Run("not equal", func(t *testing.T) {
				a := Some(1)
				b := Some(2)
				must.False(t, a.Equal(b))
			})
			t.Run("Some and None", func(t *testing.T) {
				a := Some(1)
				b := None[int]()
				must.False(t, a.Equal(b))
			})
			t.Run("None and Some", func(t *testing.T) {
				a := Some(1)
				b := None[int]()
				must.False(t, a.Equal(b))
			})
			t.Run("None and None", func(t *testing.T) {
				a := None[int]()
				b := None[int]()
				must.True(t, a.Equal(b))
			})
			t.Run("Some(Default) and None", func(t *testing.T) {
				a := Some(0)
				b := None[int]()
				must.False(t, a.Equal(b))
			})
		})
	})

	t.Run("String", func(t *testing.T) {
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
