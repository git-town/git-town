// This file defines general assertions for testing.
// Assertions take a *testing.T object and fail the test using it.

package asserts

import (
	"os"
	"testing"

	"github.com/shoenig/test/must"
)

func IsGitRepo(t *testing.T, dir string) {
	t.Helper()
	entries, err := os.ReadDir(dir)
	must.NoError(t, err)
	must.EqOp(t, ".git", entries[0].Name())
}
