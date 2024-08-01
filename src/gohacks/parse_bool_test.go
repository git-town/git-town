package gohacks_test

import (
	"testing"

	"github.com/git-town/git-town/v14/src/gohacks"
	. "github.com/git-town/git-town/v14/src/gohacks/prelude"
	"github.com/shoenig/test/must"
)

func TestParseBool(t *testing.T) {
	t.Parallel()

	t.Run("ParseBool", func(t *testing.T) {
		t.Parallel()
		tests := map[string]bool{
			"yes":      true,
			"Yes":      true,
			"YES":      true,
			"no":       false,
			"on":       true,
			"off":      false,
			"true":     true,
			"false":    false,
			"enable":   true,
			"enabled":  true,
			"disable":  false,
			"disabled": false,
			"1":        true,
			"0":        false,
		}
		for give, want := range tests {
			have, err := gohacks.ParseBool(give)
			must.NoError(t, err)
			must.EqOp(t, want, have)
		}

		t.Run("invalid input", func(t *testing.T) {
			t.Parallel()
			_, err := gohacks.ParseBool("zonk")
			must.Error(t, err)
		})
	})

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
			have, err := gohacks.ParseBoolOpt(give, "test")
			must.NoError(t, err)
			must.Eq(t, want, have)
		}

		t.Run("invalid input", func(t *testing.T) {
			t.Parallel()
			_, err := gohacks.ParseBoolOpt("zonk", "test")
			must.Error(t, err)
		})
	})
}
