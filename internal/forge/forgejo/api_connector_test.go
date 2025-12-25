package forgejo_test

import (
	"testing"

	forgejoSDK "codeberg.org/mvdkleijn/forgejo-sdk/forgejo/v2"
	"github.com/git-town/git-town/v22/internal/forge/forgejo"
	"github.com/shoenig/test/must"
)

func TestFilterPullRequests(t *testing.T) {
	t.Parallel()
	give := []*forgejoSDK.PullRequest{
		// matching branch
		{
			Head: &forgejoSDK.PRBranchInfo{Name: "branch"},
			Base: &forgejoSDK.PRBranchInfo{Name: "target"},
		},
		// branch with different name
		{
			Head: &forgejoSDK.PRBranchInfo{Name: "other"},
			Base: &forgejoSDK.PRBranchInfo{Name: "target"},
		},
		// branch with different target
		{
			Head: &forgejoSDK.PRBranchInfo{Name: "branch"},
			Base: &forgejoSDK.PRBranchInfo{Name: "other"},
		},
	}
	want := []*forgejoSDK.PullRequest{
		{
			Head: &forgejoSDK.PRBranchInfo{Name: "branch"},
			Base: &forgejoSDK.PRBranchInfo{Name: "target"},
		},
	}
	have := forgejo.FilterPullRequests(give, "branch", "target")
	must.Eq(t, want, have)
}
