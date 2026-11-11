package forgejo_test

import (
	"testing"

	forgejoSDK "codeberg.org/mvdkleijn/forgejo-sdk/forgejo/v2"
	"github.com/git-town/git-town/v24/internal/forge/forgejo"
	"github.com/shoenig/test/must"
)

func TestFilterPullRequests(t *testing.T) {
	t.Parallel()
	give := []*forgejoSDK.PullRequest{
		// matching branch
		{
			Base: &forgejoSDK.PRBranchInfo{Name: "target"},
			Head: &forgejoSDK.PRBranchInfo{Name: "branch"},
		},
		// branch with different name
		{
			Base: &forgejoSDK.PRBranchInfo{Name: "target"},
			Head: &forgejoSDK.PRBranchInfo{Name: "other"},
		},
		// branch with different target
		{
			Base: &forgejoSDK.PRBranchInfo{Name: "other"},
			Head: &forgejoSDK.PRBranchInfo{Name: "branch"},
		},
	}
	want := []*forgejoSDK.PullRequest{
		{
			Base: &forgejoSDK.PRBranchInfo{Name: "target"},
			Head: &forgejoSDK.PRBranchInfo{Name: "branch"},
		},
	}
	have := forgejo.FilterPullRequests(give, "branch", "target")
	must.Eq(t, want, have)
}
