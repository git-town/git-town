package gohacks_test

import (
	"testing"

	"github.com/git-town/git-town/v22/internal/gohacks"
	. "github.com/git-town/git-town/v22/pkg/prelude"
	"github.com/shoenig/test/must"
)

type MyBool bool

func TestParseBool(t *testing.T) {
	t.Parallel()

	t.Run("ParseBool", func(t *testing.T) {
		t.Parallel()
		t.Run("valid inputs", func(t *testing.T) {
			t.Parallel()
			tests := map[string]MyBool{
				"y":        true,
				"Y":        true,
				"yes":      true,
				"Yes":      true,
				"YES":      true,
				"n":        false,
				"N":        false,
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
				have, err := gohacks.ParseBool[MyBool](give, "test")
				must.NoError(t, err)
				must.Eq(t, want, have)
			}
		})

		t.Run("invalid inputs", func(t *testing.T) {
			t.Parallel()
			tests := []string{"", "zonk"}
			for _, give := range tests {
				_, err := gohacks.ParseBool[MyBool](give, "test")
				must.Error(t, err)
			}
		})
	})

	t.Run("ParseBoolOpt", func(t *testing.T) {
		t.Parallel()
		tests := map[string]Option[MyBool]{
			"":         None[MyBool](),
			"yes":      Some(MyBool(true)),
			"Yes":      Some(MyBool(true)),
			"YES":      Some(MyBool(true)),
			"no":       Some(MyBool(false)),
			"on":       Some(MyBool(true)),
			"off":      Some(MyBool(false)),
			"true":     Some(MyBool(true)),
			"false":    Some(MyBool(false)),
			"enable":   Some(MyBool(true)),
			"enabled":  Some(MyBool(true)),
			"disable":  Some(MyBool(false)),
			"disabled": Some(MyBool(false)),
			"1":        Some(MyBool(true)),
			"0":        Some(MyBool(false)),
		}
		for give, want := range tests {
			have, err := gohacks.ParseBoolOpt[MyBool](give, "test")
			must.NoError(t, err)
			must.Eq(t, want, have)
		}

		t.Run("invalid input", func(t *testing.T) {
			t.Parallel()
			_, err := gohacks.ParseBoolOpt[MyBool]("zonk", "test")
			must.Error(t, err)
		})
	})
}
