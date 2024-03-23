package gitdomain_test

import (
	"testing"

	"github.com/git-town/git-town/v13/src/git/gitdomain"
	"github.com/shoenig/test/must"
)

func TestRepoRootDir(t *testing.T) {
	t.Parallel()

	t.Run("IsEmpty", func(t *testing.T) {
		t.Parallel()
		tests := map[string]bool{
			"content": false,
			"":        true,
		}
		for give, want := range tests {
			rootDir := gitdomain.NewRepoRootDir(give)
			have := rootDir.IsEmpty()
			must.EqOp(t, want, have)
		}
	})
}
