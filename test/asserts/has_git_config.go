package asserts

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func HasGitConfiguration(t *testing.T, dir string) {
	t.Helper()
	entries, err := os.ReadDir(dir)
	assert.Nilf(t, err, "cannot list directory %q", dir)
	for e := range entries {
		if entries[e].Name() == ".gitconfig" {
			return
		}
	}
	t.Fatalf(".gitconfig not found in %q", dir)
}
