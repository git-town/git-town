package subshell_test

import (
	"testing"

	"github.com/git-town/git-town/v7/src/subshell"
	"github.com/stretchr/testify/assert"
)

func TestFormat(t *testing.T) {
	t.Parallel()
	tests := map[string]struct {
		branch      string
		printBranch bool
		executable  string
		args        []string
	}{
		"[branch] git checkout foo": {printBranch: true, branch: "branch", executable: "git", args: []string{"checkout", "foo"}},
		"git checkout foo":          {printBranch: false, branch: "branch", executable: "git", args: []string{"checkout", "foo"}},
	}
	for want, give := range tests {
		have := subshell.FormatCommand(give.branch, give.printBranch, give.executable, give.args...)
		assert.Equal(t, want, have)
	}
}
