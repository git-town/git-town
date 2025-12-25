package gitea_test

import (
	"testing"

	giteasdk "code.gitea.io/sdk/gitea"
	"github.com/git-town/git-town/v22/internal/forge/gitea"
	"github.com/shoenig/test/must"
)

func TestFilterGiteaPullRequests(t *testing.T) {
	t.Parallel()
	give := []*giteasdk.PullRequest{
		// matching branch
		{
			Head: &giteasdk.PRBranchInfo{
				Name: "branch",
			},
			Base: &giteasdk.PRBranchInfo{
				Name: "target",
			},
		},
		// branch with different name
		{
			Head: &giteasdk.PRBranchInfo{
				Name: "other",
			},
			Base: &giteasdk.PRBranchInfo{
				Name: "target",
			},
		},
		// branch with different target
		{
			Head: &giteasdk.PRBranchInfo{
				Name: "branch",
			},
			Base: &giteasdk.PRBranchInfo{
				Name: "other",
			},
		},
	}
	want := []*giteasdk.PullRequest{
		{
			Head: &giteasdk.PRBranchInfo{
				Name: "branch",
			},
			Base: &giteasdk.PRBranchInfo{
				Name: "target",
			},
		},
	}
	have := gitea.FilterPullRequests(give, "branch", "target")
	must.Eq(t, want, have)
}
