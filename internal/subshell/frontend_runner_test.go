package subshell_test

import (
	"testing"

	"github.com/git-town/git-town/v15/internal/git/gitdomain"
	"github.com/git-town/git-town/v15/internal/subshell"
	"github.com/shoenig/test/must"
)

func TestFormat(t *testing.T) {
	t.Parallel()
	tests := map[string]struct {
		branch      gitdomain.LocalBranchName
		printBranch bool
		executable  string
		args        []string
	}{
		"[branch] git checkout foo":        {printBranch: true, branch: gitdomain.NewLocalBranchName("branch"), executable: "git", args: []string{"checkout", "foo"}},
		"git checkout foo":                 {printBranch: false, branch: gitdomain.NewLocalBranchName("branch"), executable: "git", args: []string{"checkout", "foo"}},
		`git config perennial-branches ""`: {printBranch: false, branch: gitdomain.NewLocalBranchName("branch"), executable: "git", args: []string{"config", "perennial-branches", ""}},
	}
	for want, give := range tests {
		have := subshell.FormatCommand(give.branch, give.printBranch, give.executable, give.args...)
		must.EqOp(t, want, have)
	}
}
