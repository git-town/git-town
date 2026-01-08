package prelude_test

import (
	"testing"

	"github.com/git-town/git-town/v22/internal/config/configdomain"
	"github.com/git-town/git-town/v22/internal/forge/forgedomain"
	. "github.com/git-town/git-town/v22/pkg/prelude"
	"github.com/shoenig/test/must"
)

func TestOption(t *testing.T) {
	t.Parallel()

	t.Run("Equal", func(t *testing.T) {
		t.Parallel()
		t.Run("int", func(t *testing.T) {
			t.Parallel()
			t.Run("equal", func(t *testing.T) {
				t.Parallel()
				must.True(t, Some(1).Equal(Some(1)))
			})
			t.Run("not equal", func(t *testing.T) {
				t.Parallel()
				must.False(t, Some(1).Equal(Some(2)))
			})
			t.Run("Some and None", func(t *testing.T) {
				t.Parallel()
				must.False(t, Some(1).Equal(None[int]()))
			})
			t.Run("None and Some", func(t *testing.T) {
				t.Parallel()
				must.False(t, None[int]().Equal(Some(1)))
			})
			t.Run("None and None", func(t *testing.T) {
				t.Parallel()
				must.True(t, None[int]().Equal(None[int]()))
			})
			t.Run("Some(Default) and None", func(t *testing.T) {
				t.Parallel()
				must.False(t, Some(0).Equal(None[int]()))
			})
		})
	})

	t.Run("EqualSome", func(t *testing.T) {
		t.Parallel()
		t.Run("int", func(t *testing.T) {
			t.Parallel()
			t.Run("equal", func(t *testing.T) {
				t.Parallel()
				must.True(t, Some(1).EqualSome(1))
			})
			t.Run("not equal", func(t *testing.T) {
				t.Parallel()
				must.False(t, Some(1).EqualSome(2))
			})
			t.Run("None", func(t *testing.T) {
				t.Parallel()
				must.False(t, None[int]().EqualSome(1))
			})
		})
	})

	t.Run("MarshalJSON", func(t *testing.T) {
		t.Parallel()
		t.Run("Some", func(t *testing.T) {
			t.Parallel()
			value := Some(12)
			json, err := value.MarshalJSON()
			must.NoError(t, err)
			must.Eq(t, "12", string(json))
		})
		t.Run("None", func(t *testing.T) {
			t.Parallel()
			value := None[int]()
			json, err := value.MarshalJSON()
			must.NoError(t, err)
			must.Eq(t, "null", string(json))
		})
	})

	t.Run("NewOption", func(t *testing.T) {
		t.Parallel()
		tests := map[string]Option[string]{
			"":    None[string](),
			"foo": Some("foo"),
		}
		for give, want := range tests {
			have := NewOption(give)
			must.Eq(t, want, have)
		}
	})

	t.Run("Or", func(t *testing.T) {
		t.Parallel()
		t.Run("none or none = none", func(t *testing.T) {
			t.Parallel()
			option := None[int]()
			other := None[int]()
			have := option.Or(other)
			want := None[int]()
			must.Eq(t, want, have)
		})
		t.Run("none or some = some", func(t *testing.T) {
			t.Parallel()
			option := None[int]()
			other := Some(2)
			have := option.Or(other)
			want := Some(2)
			must.Eq(t, want, have)
		})
		t.Run("some or none = some", func(t *testing.T) {
			t.Parallel()
			option := Some(1)
			other := None[int]()
			have := option.Or(other)
			want := Some(1)
			must.Eq(t, want, have)
		})
		t.Run("some1 or some2 = some1", func(t *testing.T) {
			t.Parallel()
			option := Some(1)
			other := Some(2)
			have := option.Or(other)
			want := Some(1)
			must.Eq(t, want, have)
		})
	})

	t.Run("String", func(t *testing.T) {
		t.Parallel()
		t.Run("Some(struct that implements fmt.Stringer)", func(t *testing.T) {
			t.Parallel()
			text := "my token"
			option := Some(forgedomain.GithubToken(text))
			have := option.String()
			must.EqOp(t, "Some(my token)", have)
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
			want := "Some({true})"
			must.EqOp(t, want, have)
		})
		t.Run("None[int]", func(t *testing.T) {
			t.Parallel()
			option := None[int]()
			have := option.String()
			must.EqOp(t, "None", have)
		})
		t.Run("None[*int]", func(t *testing.T) {
			t.Parallel()
			option := None[*int]()
			have := option.String()
			must.EqOp(t, "None", have)
		})
		t.Run("None[string newtype]", func(t *testing.T) {
			t.Parallel()
			option := None[configdomain.VerifiedRegex]()
			have := option.String()
			must.EqOp(t, "None", have)
		})
	})

	t.Run("UnmarshalJSON", func(t *testing.T) {
		t.Run("Some", func(t *testing.T) {
			t.Parallel()
			json := "12"
			value := None[int]()
			err := value.UnmarshalJSON([]byte(json))
			must.NoError(t, err)
			must.True(t, value.EqualSome(12))
		})
		t.Run("None", func(t *testing.T) {
			t.Parallel()
			json := "null"
			value := None[int]()
			err := value.UnmarshalJSON([]byte(json))
			must.NoError(t, err)
			must.Eq(t, None[int](), value)
		})
	})
}
