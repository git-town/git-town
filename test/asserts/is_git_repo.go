// This file defines general assertions for testing.
// Assertions take a *testing.T object and fail the test using it.

package asserts

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func IsGitRepo(t *testing.T, dir string) {
	t.Helper()
	FolderExists(t, dir)
	entries, err := os.ReadDir(dir)
	assert.Nilf(t, err, "cannot list directory %q", dir)
	assert.Equal(t, ".git", entries[0].Name())
}
