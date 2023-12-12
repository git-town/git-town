package asserts

import (
	"os/exec"
	"testing"

	"github.com/shoenig/test/must"
)

func BranchExists(t *testing.T, dir, expectedBranch string) {
	t.Helper()
	cmd := exec.Command("git", "branch")
	cmd.Dir = dir
	output, err := cmd.CombinedOutput()
	must.NoError(t, err)
	must.StrContains(t, string(output), expectedBranch)
}
