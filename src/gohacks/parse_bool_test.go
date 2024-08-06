package gohacks_test

import (
	"testing"

	. "github.com/git-town/git-town/v14/pkg/prelude"
	"github.com/git-town/git-town/v14/src/gohacks"
	"github.com/shoenig/test/must"
)

func TestParseBool(t *testing.T) {
	t.Parallel()

	t.Run("ParseBoolOpt", func(t *testing.T) {
		t.Parallel()
		tests := map[string]Option[bool]{
			"":         None[bool](),
			"yes":      Some(true),
			"Yes":      Some(true),
			"YES":      Some(true),
			"no":       Some(false),
			"on":       Some(true),
			"off":      Some(false),
			"true":     Some(true),
			"false":    Some(false),
			"enable":   Some(true),
			"enabled":  Some(true),
			"disable":  Some(false),
			"disabled": Some(false),
			"1":        Some(true),
			"0":        Some(false),
		}
		for give, want := range tests {
			have, err := gohacks.ParseBool(give, "test")
			must.NoError(t, err)
			must.Eq(t, want, have)
		}

		t.Run("invalid input", func(t *testing.T) {
			t.Parallel()
			_, err := gohacks.ParseBool("zonk", "test")
			must.Error(t, err)
		})
	})
}
