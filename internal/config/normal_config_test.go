package config_test

import (
	"testing"

	"github.com/git-town/git-town/v24/internal/config"
	"github.com/git-town/git-town/v24/internal/git/gitdomain"
	. "github.com/git-town/git-town/v24/pkg/prelude"
	"github.com/shoenig/test/must"
)

func TestNormalConfig(t *testing.T) {
	t.Parallel()

	t.Run("Author", func(t *testing.T) {
		t.Parallel()
		config := config.NormalConfig{
			GitUserEmail: Some(gitdomain.GitUserEmail("email")),
			GitUserName:  Some(gitdomain.GitUserName("name")),
		}
		have := config.Author().GetOrPanic()
		want := gitdomain.Author("name <email>")
		must.EqOp(t, want, have)
	})
}
