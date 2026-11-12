package gohacks_test

import (
	"testing"

	"github.com/git-town/git-town/v24/internal/gohacks"
	"github.com/git-town/git-town/v24/internal/gohacks/stringss"
	. "github.com/git-town/git-town/v24/pkg/prelude"
	"github.com/shoenig/test/must"
)

type MyBool bool

func TestParseBool(t *testing.T) {
	t.Parallel()

	t.Run("ParseBool", func(t *testing.T) {
		t.Parallel()
		t.Run("valid inputs", func(t *testing.T) {
			t.Parallel()
			tests := map[stringss.Trimmed]MyBool{
				"0":        false,
				"1":        true,
				"disable":  false,
				"disabled": false,
				"enable":   true,
				"enabled":  true,
				"false":    false,
				"N":        false,
				"n":        false,
				"no":       false,
				"off":      false,
				"on":       true,
				"true":     true,
				"Y":        true,
				"y":        true,
				"YES":      true,
				"Yes":      true,
				"yes":      true,
			}
			for give, want := range tests {
				have, err := gohacks.ParseBool[MyBool](give, "test")
				must.NoError(t, err)
				must.Eq(t, want, have)
			}
		})

		t.Run("invalid inputs", func(t *testing.T) {
			t.Parallel()
			tests := []stringss.Trimmed{"", "zonk"}
			for _, give := range tests {
				_, err := gohacks.ParseBool[MyBool](give, "test")
				must.Error(t, err)
			}
		})
	})

	t.Run("ParseBoolOpt", func(t *testing.T) {
		t.Parallel()
		tests := map[stringss.Trimmed]Option[MyBool]{
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
