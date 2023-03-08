package cmd_test

import (
	"testing"

	"github.com/git-town/git-town/v7/src/cmd"
	"github.com/stretchr/testify/assert"
)

func TestNewCompletionType(t *testing.T) {
	t.Parallel()
	t.Run("valid input", func(t *testing.T) {
		t.Parallel()
		tests := map[string]cmd.CompletionType{
			"bash":       cmd.CompletionTypeBash,
			"Bash":       cmd.CompletionTypeBash,
			"BASH":       cmd.CompletionTypeBash,
			"zsh":        cmd.CompletionTypeZsh,
			"fish":       cmd.CompletionTypeFish,
			"powershell": cmd.CompletionTypePowershell,
		}
		for give, want := range tests {
			have, err := cmd.NewCompletionType(give)
			assert.Nil(t, err)
			assert.Equal(t, want, have)
		}
	})
	t.Run("invalid input", func(t *testing.T) {
		_, err := cmd.NewCompletionType("zonk")
		assert.Error(t, err)
	})
}
