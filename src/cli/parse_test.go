package cli_test

import (
	"testing"

	"github.com/git-town/git-town/v7/src/cli"
	"github.com/stretchr/testify/assert"
)

func TestParseBool(t *testing.T) {
	t.Parallel()
	t.Run("valid input", func(t *testing.T) {
		t.Parallel()
		tests := map[string]bool{
			"yes":   true,
			"no":    false,
			"on":    true,
			"off":   false,
			"true":  true,
			"false": false,
			"1":     true,
			"0":     false,
		}
		for give, want := range tests {
			have, err := cli.ParseBool(give)
			assert.Nil(t, err)
			assert.Equal(t, want, have)
		}
	})

	t.Run("case insensitive", func(t *testing.T) {
		t.Parallel()
		for _, give := range []string{"yes", "Yes", "YES"} {
			have, err := cli.ParseBool(give)
			assert.Nil(t, err)
			assert.Equal(t, true, have)
		}
	})

	t.Run("invalid input", func(t *testing.T) {
		t.Parallel()
		_, err := cli.ParseBool("zonk")
		assert.Error(t, err)
	})
}
