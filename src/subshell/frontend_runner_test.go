package subshell_test

import (
	"testing"

	"github.com/git-town/git-town/v14/src/git/gitdomain"
	"github.com/git-town/git-town/v14/src/subshell"
	"github.com/shoenig/test/must"
)

func TestFormat(t *testing.T) {
	t.Parallel()
	tests := map[string]struct {
		branch     gitdomain.LocalBranchName
		omitBranch bool
		executable string
		args       []string
	}{
		"[branch] git checkout foo":        {omitBranch: false, branch: gitdomain.NewLocalBranchName("branch"), executable: "git", args: []string{"checkout", "foo"}},
		"git checkout foo":                 {omitBranch: true, branch: gitdomain.NewLocalBranchName("branch"), executable: "git", args: []string{"checkout", "foo"}},
		`git config perennial-branches ""`: {omitBranch: true, branch: gitdomain.NewLocalBranchName("branch"), executable: "git", args: []string{"config", "perennial-branches", ""}},
	}
	for want, give := range tests {
		have := subshell.FormatCommand(give.branch, give.omitBranch, give.executable, give.args...)
		must.EqOp(t, want, have)
	}
}
