package gitea_test

import (
	"testing"

	giteasdk "code.gitea.io/sdk/gitea"
	"github.com/git-town/git-town/v24/internal/forge/gitea"
	"github.com/shoenig/test/must"
)

func TestFilterGiteaPullRequests(t *testing.T) {
	t.Parallel()
	give := []*giteasdk.PullRequest{
		// matching branch
		{
			Base: &giteasdk.PRBranchInfo{
				Name: "target",
			},
			Head: &giteasdk.PRBranchInfo{
				Name: "branch",
			},
		},
		// branch with different name
		{
			Base: &giteasdk.PRBranchInfo{
				Name: "target",
			},
			Head: &giteasdk.PRBranchInfo{
				Name: "other",
			},
		},
		// branch with different target
		{
			Base: &giteasdk.PRBranchInfo{
				Name: "other",
			},
			Head: &giteasdk.PRBranchInfo{
				Name: "branch",
			},
		},
	}
	want := []*giteasdk.PullRequest{
		{
			Base: &giteasdk.PRBranchInfo{
				Name: "target",
			},
			Head: &giteasdk.PRBranchInfo{
				Name: "branch",
			},
		},
	}
	have := gitea.FilterPullRequests(give, "branch", "target")
	must.Eq(t, want, have)
}
