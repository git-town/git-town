package subshell_test

import (
	"testing"

	"github.com/git-town/git-town/v22/internal/git/gitdomain"
	"github.com/git-town/git-town/v22/internal/subshell"
	"github.com/shoenig/test/must"
)

func TestFormat(t *testing.T) {
	t.Parallel()
	tests := map[string]struct {
		branch      gitdomain.Location
		printBranch bool
		executable  string
		args        []string
	}{
		"[branch] git checkout foo":        {printBranch: true, branch: "branch", executable: "git", args: []string{"checkout", "foo"}},
		"git checkout foo":                 {printBranch: false, branch: "branch", executable: "git", args: []string{"checkout", "foo"}},
		`git config perennial-branches ""`: {printBranch: false, branch: "branch", executable: "git", args: []string{"config", "perennial-branches", ""}},
	}
	for want, give := range tests {
		have := subshell.FormatCommand(give.branch, give.printBranch, []string{}, give.executable, give.args...)
		must.EqOp(t, want, have)
	}
}
