package subshell_test

import (
	"testing"

	"github.com/git-town/git-town/v24/internal/git/gitdomain"
	"github.com/git-town/git-town/v24/internal/subshell"
	"github.com/shoenig/test/must"
)

func TestFormat(t *testing.T) {
	t.Parallel()
	tests := map[string]struct {
		args        []string
		branch      gitdomain.Location
		executable  string
		printBranch bool
	}{
		"[branch] git checkout foo":        {args: []string{"checkout", "foo"}, branch: "branch", executable: "git", printBranch: true},
		"git checkout foo":                 {args: []string{"checkout", "foo"}, branch: "branch", executable: "git", printBranch: false},
		`git config perennial-branches ""`: {args: []string{"config", "perennial-branches", ""}, branch: "branch", executable: "git", printBranch: false},
	}
	for want, give := range tests {
		have := subshell.FormatCommand(give.branch, give.printBranch, []string{}, give.executable, give.args...)
		must.EqOp(t, want, have)
	}
}
