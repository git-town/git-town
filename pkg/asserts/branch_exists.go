package asserts

import (
	"context"
	"os/exec"
	"testing"

	"github.com/shoenig/test/must"
)

func BranchExists(t *testing.T, dir, expectedBranch string) {
	t.Helper()
	cmd := exec.CommandContext(context.Background(), "git", "rev-parse", "--verify", "-q", "refs/heads/"+expectedBranch) //nolint:gosec
	cmd.Dir = dir
	err := cmd.Run()
	must.NoError(t, err)
}
