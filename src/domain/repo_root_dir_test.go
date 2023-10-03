package domain_test

import (
	"testing"

	"github.com/git-town/git-town/v9/src/domain"
	"github.com/shoenig/test"
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
			test.EqOp(t, want, have)
		}
	})
}
