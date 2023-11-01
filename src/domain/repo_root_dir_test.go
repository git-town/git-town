package domain_test

import (
	"testing"

	"github.com/git-town/git-town/v10/src/domain"
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
			rootDir := domain.NewRepoRootDir(give)
			have := rootDir.IsEmpty()
			must.EqOp(t, want, have)
		}
	})
}
