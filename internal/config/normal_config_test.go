package config_test

import (
	"testing"

	"github.com/git-town/git-town/v23/internal/config"
	"github.com/git-town/git-town/v23/internal/config/configdomain"
	"github.com/git-town/git-town/v23/internal/git/gitdomain"
	. "github.com/git-town/git-town/v23/pkg/prelude"
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

	t.Run("NewNormalConfigFromPartial uses explicit empty breadcrumb exclusions", func(t *testing.T) {
		t.Parallel()
		defaults := config.DefaultNormalConfig()
		defaults.ProposalBreadcrumbExcludeBranches = configdomain.NewProposalBreadcrumbExcludeBranches(configdomain.BranchTypePrototypeBranch)
		partial := configdomain.PartialConfig{
			ProposalBreadcrumbExcludeBranches: Some(configdomain.NewProposalBreadcrumbExcludeBranches()),
		}
		have := config.NewNormalConfigFromPartial(partial, defaults)
		want := configdomain.NewProposalBreadcrumbExcludeBranches()
		must.Eq(t, want, have.ProposalBreadcrumbExcludeBranches)
	})
}
