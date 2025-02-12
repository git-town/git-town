package configdomain_test

import (
	"testing"

	"github.com/git-town/git-town/v18/internal/config/configdomain"
	"github.com/git-town/git-town/v18/internal/git/gitdomain"
	"github.com/shoenig/test/must"
)

func TestValidatedConfig(t *testing.T) {
	t.Parallel()

	t.Run("Author", func(t *testing.T) {
		t.Parallel()
		config := configdomain.ValidatedConfigData{
			GitUserName:  "name",
			GitUserEmail: "email",
		}
		have := config.Author()
		want := gitdomain.Author("name <email>")
		must.EqOp(t, want, have)
	})

	t.Run("IsMainBranch", func(t *testing.T) {
		t.Parallel()
		config := configdomain.ValidatedConfigData{
			MainBranch: "main",
		}
		must.False(t, config.IsMainBranch("feature"))
		must.True(t, config.IsMainBranch("main"))
		must.False(t, config.IsMainBranch("peren1"))
		must.False(t, config.IsMainBranch("peren2"))
	})
}
