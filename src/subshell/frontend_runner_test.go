package subshell_test

import (
	"testing"

	"github.com/git-town/git-town/v9/src/domain"
	"github.com/git-town/git-town/v9/src/subshell"
	"github.com/stretchr/testify/assert"
)

func TestFormat(t *testing.T) {
	t.Parallel()
	tests := map[string]struct {
		branch     domain.LocalBranchName
		omitBranch bool
		executable string
		args       []string
	}{
		"[branch] git checkout foo": {omitBranch: false, branch: domain.NewLocalBranchName("branch"), executable: "git", args: []string{"checkout", "foo"}},
		"git checkout foo":          {omitBranch: true, branch: domain.NewLocalBranchName("branch"), executable: "git", args: []string{"checkout", "foo"}},
	}
	for want, give := range tests {
		have := subshell.FormatCommand(give.branch, give.omitBranch, give.executable, give.args...)
		assert.Equal(t, want, have)
	}
}
