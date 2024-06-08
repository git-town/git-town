package dialog_test

import (
	"testing"

	"github.com/git-town/git-town/v14/src/cli/dialog"
	"github.com/git-town/git-town/v14/src/cli/dialog/components"
	"github.com/git-town/git-town/v14/src/config/configdomain"
	"github.com/git-town/git-town/v14/src/git/gitdomain"
	"github.com/shoenig/test/must"
)

func TestParent(t *testing.T) {
	t.Parallel()

	t.Run("ParentEntries", func(t *testing.T) {
		t.Parallel()
		t.Run("omits the branch for which to select the parent", func(t *testing.T) {
			t.Parallel()
			branch1 := gitdomain.NewLocalBranchName("branch-1")
			branch2 := gitdomain.NewLocalBranchName("branch-2")
			branch3 := gitdomain.NewLocalBranchName("branch-3")
			main := gitdomain.NewLocalBranchName("main")
			localBranches := gitdomain.LocalBranchNames{branch1, branch2, branch3, main}
			lineage := configdomain.NewLineage()
			lineage.Add(branch1, main)
			lineage.Add(branch2, main)
			lineage.Add(branch3, main)
			have := dialog.ParentCandidateNames(dialog.ParentArgs{
				Branch:          branch2,
				DefaultChoice:   main,
				DialogTestInput: components.TestInput{},
				Lineage:         lineage,
				LocalBranches:   localBranches,
				MainBranch:      main,
			})
			want := gitdomain.LocalBranchNames{dialog.PerennialBranchOption, main, branch1, branch3}
			must.Eq(t, want, have)
		})
		t.Run("omits all descendents of the branch for which to select the parent", func(t *testing.T) {
			t.Parallel()
			branch1 := gitdomain.NewLocalBranchName("branch-1")
			branch1a := gitdomain.NewLocalBranchName("branch-1a")
			branch2 := gitdomain.NewLocalBranchName("branch-2")
			branch2a := gitdomain.NewLocalBranchName("branch-2a")
			branch3 := gitdomain.NewLocalBranchName("branch-3")
			main := gitdomain.NewLocalBranchName("main")
			localBranches := gitdomain.LocalBranchNames{branch1, branch1a, branch2, branch2a, branch3, main}
			lineage := configdomain.NewLineage()
			lineage.Add(branch1, main)
			lineage.Add(branch1a, branch1)
			lineage.Add(branch2, main)
			lineage.Add(branch2a, branch2)
			lineage.Add(branch3, main)
			have := dialog.ParentCandidateNames(dialog.ParentArgs{
				Branch:          branch2,
				DefaultChoice:   main,
				DialogTestInput: components.TestInput{},
				Lineage:         lineage,
				LocalBranches:   localBranches,
				MainBranch:      main,
			})
			want := gitdomain.LocalBranchNames{dialog.PerennialBranchOption, main, branch1, branch1a, branch3}
			must.Eq(t, want, have)
		})
	})
}
