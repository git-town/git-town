package cmd_test

import (
	"testing"

	"github.com/git-town/git-town/v14/src/cmd"
	"github.com/shoenig/test/must"
)

func TestNewCompletionType(t *testing.T) {
	t.Parallel()

	t.Run("recognizes shells", func(t *testing.T) {
		t.Parallel()
		tests := map[string]cmd.CompletionType{
			"bash":       cmd.CompletionTypeBash,
			"zsh":        cmd.CompletionTypeZsh,
			"fish":       cmd.CompletionTypeFish,
			"powershell": cmd.CompletionTypePowershell,
		}
		for give, want := range tests {
			have, err := cmd.NewCompletionType(give)
			must.NoError(t, err)
			must.EqOp(t, want, have)
		}
	})

	t.Run("case insensitive", func(t *testing.T) {
		t.Parallel()
		for _, give := range []string{"bash", "Bash", "BASH"} {
			have, err := cmd.NewCompletionType(give)
			must.NoError(t, err)
			must.EqOp(t, cmd.CompletionTypeBash, have)
		}
	})

	t.Run("invalid input", func(t *testing.T) {
		t.Parallel()
		_, err := cmd.NewCompletionType("zonk")
		must.Error(t, err)
	})
}
