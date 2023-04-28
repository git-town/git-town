package asserts

import (
	"os/exec"
	"testing"

	"github.com/stretchr/testify/assert"
)

func BranchExists(t *testing.T, dir, expectedBranch string) {
	t.Helper()
	cmd := exec.Command("git", "branch")
	cmd.Dir = dir
	output, err := cmd.CombinedOutput()
	assert.Nilf(t, err, "cannot run 'git branch' in %q", dir)
	assert.Contains(t, string(output), expectedBranch, "doesn't have Git branch")
}
