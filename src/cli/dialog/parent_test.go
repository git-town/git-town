package dialog_test

import (
	"testing"

	"github.com/git-town/git-town/v13/src/cli/dialog"
	"github.com/git-town/git-town/v13/src/cli/dialog/components"
	"github.com/git-town/git-town/v13/src/config/configdomain"
	"github.com/git-town/git-town/v13/src/git/gitdomain"
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
			lineage := configdomain.Lineage{
				branch1: main,
				branch2: main,
				branch3: main,
			}
			have := dialog.ParentEntries(dialog.ParentArgs{
				Branch:          branch2,
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
			lineage := configdomain.Lineage{
				branch1:  main,
				branch1a: branch1,
				branch2:  main,
				branch2a: branch2,
				branch3:  main,
			}
			have := dialog.ParentEntries(dialog.ParentArgs{
				Branch:          branch2,
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
